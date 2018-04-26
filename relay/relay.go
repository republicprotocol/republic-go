package relay

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/republicprotocol/republic-go/darkocean"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/rpc/relayer"
	"github.com/republicprotocol/republic-go/rpc/smpcer"
	"github.com/republicprotocol/republic-go/rpc/swarmer"
	"github.com/republicprotocol/republic-go/stackint"
	"google.golang.org/grpc"
)

var prime, _ = stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")

type Config struct {
	KeyPair      identity.KeyPair
	MultiAddress identity.MultiAddress
	Token        string
}

type Relay struct {
	Config

	registry dnr.DarknodeRegistry

	relayer   relayer.Relayer
	orderbook *orderbook.Orderbook

	relayerClient *relayer.Client
	smpcerClient  *smpcer.Client
	swarmerClient *swarmer.Client
}

// NewRelay returns a new Relay object
func NewRelay(config Config, registry dnr.DarknodeRegistry, orderbook *orderbook.Orderbook, relayerClient *relayer.Client, smpcerClient *smpcer.Client, swarmerClient *swarmer.Client) Relay {
	return Relay{
		Config: config,

		registry: registry,

		relayer:   relayer.NewRelayer(relayerClient, orderbook),
		orderbook: orderbook,

		relayerClient: relayerClient,
		smpcerClient:  smpcerClient,
		swarmerClient: swarmerClient,
	}
}

// ListenAndServe on a binding address and port for HTTP requests.
func (relay *Relay) ListenAndServe(bind, port string) error {
	r := mux.NewRouter().StrictSlash(true)
	r.Methods("POST").Path("/orders").Handler(RecoveryHandler(relay.AuthorizationHandler(relay.OpenOrdersHandler())))
	r.Methods("GET").Path("/orders").Handler(RecoveryHandler(relay.AuthorizationHandler(relay.GetOrdersHandler())))
	r.Methods("GET").Path("/orders/{orderID}").Handler(RecoveryHandler(relay.AuthorizationHandler(relay.GetOrderHandler())))
	r.Methods("DELETE").Path("/orders/{orderID}").Handler(RecoveryHandler(relay.AuthorizationHandler(relay.CancelOrderHandler())))
	if err := http.ListenAndServe(fmt.Sprintf("%v:%v", bind, port), r); err != nil {
		return fmt.Errorf("cannot listen and serve: %v", err)
	}
	return nil
}

// Register the Relay service to a grpc.Server.
func (relay *Relay) Register(server *grpc.Server) {
	relay.relayer.Register(server)
}

// Sync the Relay with other random Relay in the network.
func (relay *Relay) Sync(ctx context.Context, peers int) {
	relay.relayerClient.Sync(ctx, relay.orderbook, peers)
}

// SendOrderToDarkOcean will fragment and send orders to the dark ocean
func (relay *Relay) SendOrderToDarkOcean(openOrder order.Order) error {
	// Construct a new DarkOcean and store it in the relay to get the current
	// pool layout.
	epoch, err := relay.registry.CurrentEpoch()
	if err != nil {
		return err
	}
	darknodeIDs, err := relay.registry.GetAllNodes()
	if err != nil {
		return err
	}
	darkOcean := darkocean.NewDarkOcean(epoch.Blockhash, darknodeIDs)
	relay.ocean = &darkOcean
	darkPools := darkOcean.Pools()

	errCh := make(chan error, len(darkPools))

	multiSignature, err := relay.Config.KeyPair.Sign(&openOrder)
	if err != nil {
		return err
	}
	openOrder.Signature = multiSignature

	go func() {
		defer close(errCh)

		var wg sync.WaitGroup
		wg.Add(len(darkPools))

		for i := range darkPools {
			go func(darkPool darkocean.Pool) {
				defer wg.Done()
				// Split order into (number of nodes in each pool) * 2/3 fragments
				shares, err := openOrder.Split(int64(darkPool.Size()), int64(darkPool.Size()*2/3), &prime)
				if err != nil {
					errCh <- err
					return
				}
				if err := relay.sendSharesToDarkPool(darkPool, shares); err != nil {
					errCh <- err
					return
				}
			}(darkPools[i])
		}

		wg.Wait()
	}()

	var errNum int
	// var err error
	for errLocal := range errCh {
		if errLocal != nil {
			errNum++
			if err == nil {
				err = errLocal
			}
		}
	}

	if len(darkPools) > 0 && errNum == len(darkPools) {
		return fmt.Errorf("could not send order to any dark pool: %v", err)
	}
	return nil
}

// SendOrderFragmentsToDarkOcean will send order fragments to the dark ocean
func (relay *Relay) SendOrderFragmentsToDarkOcean(order OrderFragments) error {
	valid := false
	for poolIndex := range relay.ocean.Pools() {
		fragments := order.DarkPools[GeneratePoolID(relay.ocean.Pools()[poolIndex])]
		if fragments != nil && isSafeToSend(len(fragments), relay.ocean.Pools()[poolIndex].Size()) {
			if err := relay.sendSharesToDarkPool(relay.ocean.Pools()[poolIndex], fragments); err == nil {
				valid = true
			}
		}
	}
	if !valid && len(relay.ocean.Pools()) > 0 {
		return fmt.Errorf("cannot send fragments to pools: number of fragments do not match pool size")
	}
	return nil
}

// CancelOrder will cancel orders that aren't confirmed or settled in the dark ocean
func (relay *Relay) CancelOrder(cancelOrder order.ID) error {
	errs := make(chan error, len(relay.ocean.Pools()))
	go func() {
		defer close(errs)

		dispatch.CoForAll(relay.ocean.Pools(), func(i int) {
			addrs := relay.ocean.Pools()[i].Addresses()

			dispatch.CoForAll(addrs, func(j int) {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				multiAddr, err := relay.swarmerClient.Query(ctx, addrs[j], 1)
				if err != nil {
					errs <- err
					return
				}

				// FIXME: Encryption
				if err := relay.smpcerClient.CloseOrder(ctx, multiAddr, cancelOrder); err != nil {
					errs <- err
					return
				}
			})
		})
	}()

	// Store the first error that occurred, if any
	var err error
	var errNum int
	for errLocal := range errs {
		if errLocal != nil {
			errNum++
			if err == nil {
				err = errLocal
			}
		}
	}

	// FIXME: Error if at least one dark pool could not cancel the order
	if len(relay.ocean.Pools()) > 0 && errNum == len(relay.ocean.Pools()) {
		return fmt.Errorf("could not cancel order to any dark pool: %v", err)
	}
	return nil
}

// Send the shares across all nodes within the Dark Pool
func (relay *Relay) sendSharesToDarkPool(pool darkocean.Pool, shares []*order.Fragment) error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)

		shareIndex := 0

		var wg sync.WaitGroup
		wg.Add(pool.Size())

		// pool.For(func(n darknode.Darknode) {
		for _, node := range pool.Addresses() {
			if shareIndex < len(shares) {
				go func(nodeAddress identity.Address, share *order.Fragment) {
					defer wg.Done()

					shareSignature, err := relay.Config.KeyPair.Sign(share)
					if err != nil {
						errCh <- err
						return
					}
					share.Signature = shareSignature

					ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
					defer cancel()

					nodeMultiAddr, err := relay.swarmerClient.Query(ctx, nodeAddress, 1)
					if err != nil {
						errCh <- err
						return
					}

					// FIXME: Encryption
					if err := relay.smpcerClient.OpenOrder(ctx, nodeMultiAddr, *share); err != nil {
						errCh <- err
						return
					}
				}(node, shares[shareIndex])
			}
			shareIndex++
		}
		wg.Wait()
	}()

	var errNum int
	var err error
	for errLocal := range errCh {
		if errLocal != nil {
			errNum++
			if err == nil {
				err = errLocal
			}
		}
	}
	// check if atleast 2/3 nodes of the pool has recieved the order fragments
	if pool.Size() > 0 && errNum > ((1/3)*pool.Size()) {
		return fmt.Errorf("could not send orders to %v nodes (out of %v nodes) in pool %v", errNum, pool.Size(), GeneratePoolID(pool))
	}
	return nil
}

// GeneratePoolID will generate crypto hash for a pool
func GeneratePoolID(pool darkocean.Pool) string {
	poolID := pool.ID()
	return string(poolID[:])
}

// isSafeToSend checks if there are enough fragments to successfully complete sending an order
func isSafeToSend(fragmentCount, poolSize int) bool {
	return fragmentCount >= 2/3*poolSize && fragmentCount <= poolSize
}
