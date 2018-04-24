package relay

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/republic-go/darkocean"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
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
	DarkPools darkocean.Pools
	Registrar contracts.DarkNodeRegistry

	orderbook     *orderbook.Orderbook
	relayer       relayer.Relayer
	relayerClient relayer.Client
	swarmerClient swarmer.Client
	smpcerClient  smpcer.Client
}

// NewRelay returns a new Relay object
func NewRelay(config Config, pools darkocean.Pools, registrar contracts.DarkNodeRegistry, orderbook *orderbook.Orderbook, relayerClient relayer.Client, swarmerClient swarmer.Client, smpcerClient smpcer.Client) Relay {
	return Relay{
		Config:    config,
		DarkPools: pools,
		Registrar: registrar,

		orderbook:     orderbook,
		relayer:       relayer.NewRelayer(orderbook),
		relayerClient: relayerClient,
		swarmerClient: swarmerClient,
		smpcerClient:  smpcerClient,
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
func (relay *Relay) Sync(ctx context.Context, epoch [32]byte, peers int) {
	relay.relayerClient.Sync(ctx, relay.orderbook, epoch, peers)
}

// SendOrderToDarkOcean will fragment and send orders to the dark ocean
func (relay *Relay) SendOrderToDarkOcean(openOrder order.Order) error {
	errCh := make(chan error, len(relay.DarkPools))

	multiSignature, err := relay.Config.KeyPair.Sign(&openOrder)
	if err != nil {
		return err
	}
	openOrder.Signature = multiSignature

	go func() {
		defer close(errCh)

		var wg sync.WaitGroup
		wg.Add(len(relay.DarkPools))

		for i := range relay.DarkPools {
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
			}(relay.DarkPools[i])
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

	if len(relay.DarkPools) > 0 && errNum == len(relay.DarkPools) {
		return fmt.Errorf("could not send order to any dark pool: %v", err)
	}
	return nil
}

// SendOrderFragmentsToDarkOcean will send order fragments to the dark ocean
func (relay *Relay) SendOrderFragmentsToDarkOcean(order OrderFragments) error {
	valid := false
	for poolIndex := range relay.DarkPools {
		fragments := order.DarkPools[GeneratePoolID(relay.DarkPools[poolIndex])]
		if fragments != nil && isSafeToSend(len(fragments), relay.DarkPools[poolIndex].Size()) {
			if err := relay.sendSharesToDarkPool(relay.DarkPools[poolIndex], fragments); err == nil {
				valid = true
			}
		}
	}
	if !valid && len(relay.DarkPools) > 0 {
		return fmt.Errorf("cannot send fragments to pools: number of fragments do not match pool size")
	}
	return nil
}

// CancelOrder will cancel orders that aren't confirmed or settled in the dark ocean
func (relay *Relay) CancelOrder(cancelOrder order.ID) error {
	errCh := make(chan error, len(relay.DarkPools))

	orderSignature, err := relay.Config.KeyPair.Sign(&cancelOrder)
	if err != nil {
		return err
	}

	go func() {
		defer close(errCh)
		// For every Dark Pool
		for i := range relay.DarkPools {
			// Cancel orders for all nodes in the pool
			var wg sync.WaitGroup
			wg.Add(relay.DarkPools[i].Size())
			for _, node := range relay.DarkPools[i].Addresses() {
				defer wg.Done()
				// Get multiaddress
				multiaddress, err := getMultiAddress(node, relay)
				if err != nil {
					errCh <- fmt.Errorf("cannot read multi-address: %v", err)
					return
				}

				// Create a client
				client, err := rpc.NewClient(context.Background(), multiaddress, relay.Config.MultiAddress, relay.Config.KeyPair)
				if err != nil {
					errCh <- fmt.Errorf("cannot connect to client: %v", err)
					return
				}

				// TODO: Send signature here?
				cancelOrderRequest := &rpc.CancelOrderRequest{
					From: &rpc.MultiAddress{
						Signature:    []byte{},
						MultiAddress: relay.Config.MultiAddress.String(),
					},
					OrderFragmentId: &rpc.OrderFragmentId{
						Signature:       orderSignature,
						OrderFragmentId: cancelOrder,
					},
				}

				// Cancel order
				err = client.CancelOrder(context.Background(), cancelOrderRequest)
				if err != nil {
					errCh <- fmt.Errorf("cannot cancel order to %v", base58.Encode(node.ID()))
					return
				}
			}
			wg.Wait()
		}
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

	if len(relay.DarkPools) > 0 && errNum == len(relay.DarkPools) {
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

					relaySignature, err := relay.Config.KeyPair.Sign(relay.Config.MultiAddress)
					if err != nil {
						errCh <- err
						return
					}

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
