package relay

import (
	"context"
	"fmt"
	"log"
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

// Relays is an alias
type Relays []Relay

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
	darkPools, err := getPools(&relay.registry)
	if err != nil {
		return err
	}
	errCh := make(chan error, len(darkPools))

	openOrderSignature, err := relay.Config.Keystore.Sign(&openOrder)
	if err != nil {
		return err
	}
	openOrder.Signature = openOrderSignature

	go func() {
		defer close(errCh)

		var wg sync.WaitGroup
		wg.Add(len(darkPools))

		for i := range darkPools {
			go func(darkPool darkocean.Pool) {
				defer wg.Done()
				// Split order into (number of nodes in each pool) * 2/3 fragments
				// todo : change back to
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
	darkPools, err := getPools(&relay.registry)
	if err != nil {
		return err
	}
	valid := false
	for poolIndex := range darkPools {
		fragments := order.DarkPools[GeneratePoolID(darkPools[poolIndex])]
		if fragments != nil && isSafeToSend(len(fragments), darkPools[poolIndex].Size()) {
			if err := relay.sendSharesToDarkPool(darkPools[poolIndex], fragments); err == nil {
				valid = true
			}
		}
	}
	if !valid && len(darkPools) > 0 {
		return fmt.Errorf("cannot send fragments to pools: number of fragments do not match pool size")
	}
	return nil
}

// CancelOrder will cancel orders that aren't confirmed or settled in the dark ocean
func (relay *Relay) CancelOrder(cancelOrder order.ID) error {
	darkPools, err := getPools(&relay.registry)
	if err != nil {
		return err
	}
	errs := make(chan error, len(darkPools))
	go func() {
		defer close(errs)

		dispatch.CoForAll(darkPools, func(i int) {
			addrs := darkPools[i].Addresses()

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
	if len(darkPools) > 0 && errNum == len(darkPools) {
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

					relay.swarmerClient.Address()

					ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
					defer cancel()

					nodeMultiAddr, err := relay.swarmerClient.Query(ctx, nodeAddress, -1)
					if err != nil {
						log.Printf("fail to query peers %s", err.Error())
						errCh <- err
						return
					}

					if err := relay.smpcerClient.OpenOrder(ctx, nodeMultiAddr, *share); err != nil {
						log.Printf("fail to open order %s", err.Error())
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
			log.Println(errLocal)
			errNum++
			if err == nil {
				err = errLocal
			}
		}
	}
	// Check if at least 2/3 of the nodes in the specified pool have recieved
	// the order fragments.
	if pool.Size() > 0 && errNum > pool.Size()/3 {
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

func getPools(registry *dnr.DarknodeRegistry) (darkocean.Pools, error) {
	// Construct a new DarkOcean using the registry and use it to retrieve the
	// current pool layout.
	epoch, err := registry.CurrentEpoch()
	if err != nil {
		return darkocean.Pools{}, err
	}
	darknodeIDs, err := registry.GetAllNodes()
	if err != nil {
		return darkocean.Pools{}, err
	}
	darkOcean := darkocean.NewDarkOcean(epoch.Blockhash, darknodeIDs)
	return darkOcean.Pools(), nil
}
