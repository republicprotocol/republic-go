package relay

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/republic-go/darknode"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/rpc"
	"github.com/republicprotocol/republic-go/stackint"
)

var prime, _ = stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")

// Relay consists of configuration values (?)
type Relay struct {
	keyPair        identity.KeyPair
	multiAddress   identity.MultiAddress
	darkPools      darknode.Pools
	token          string
	bootstrapNodes []string
	orderbook      orderbook.Orderbook
}

// NewRelay returns a new Relay object
func NewRelay(keyPair identity.KeyPair, multi identity.MultiAddress, pools darknode.Pools, authToken string, bootstrapNodes []string, book orderbook.Orderbook) Relay {
	return Relay{
		keyPair:        keyPair,
		multiAddress:   multi,
		darkPools:      pools,
		token:          authToken,
		bootstrapNodes: bootstrapNodes,
		orderbook:      book,
	}
}

// NewRouter prepares Relay to handle HTTP requests
func NewRouter(relay Relay) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Methods("POST").Path("/orders").Handler(RecoveryHandler(AuthorizationHandler(OpenOrdersHandler(relay), relay.token)))
	r.Methods("GET").Path("/orders").Handler(RecoveryHandler(AuthorizationHandler(GetOrdersHandler(&relay.orderbook), relay.token)))
	r.Methods("GET").Path("/orders/{orderID}").Handler(RecoveryHandler(AuthorizationHandler(GetOrderHandler(&relay.orderbook, ""), relay.token)))
	r.Methods("DELETE").Path("/orders/{orderID}").Handler(RecoveryHandler(AuthorizationHandler(CancelOrderHandler(relay), relay.token)))
	return r
}

// SendOrderToDarkOcean will fragment and send orders to the dark ocean
func SendOrderToDarkOcean(openOrder order.Order, relayConfig Relay) error {
	errCh := make(chan error, len(relayConfig.darkPools))

	multiSignature, err := relayConfig.keyPair.Sign(&openOrder)
	if err != nil {
		return err
	}
	openOrder.Signature = multiSignature

	go func() {
		defer close(errCh)

		var wg sync.WaitGroup
		wg.Add(len(relayConfig.darkPools))

		for i := range relayConfig.darkPools {
			go func(darkPool *darknode.Pool) {
				defer wg.Done()
				// Split order into (number of nodes in each pool) * 2/3 fragments
				shares, err := openOrder.Split(int64(darkPool.Size()), int64(darkPool.Size()*2/3), &prime)
				if err != nil {
					errCh <- err
					return
				}
				if err := sendSharesToDarkPool(darkPool, shares, relayConfig); err != nil {
					errCh <- err
					return
				}
			}(relayConfig.darkPools[i])
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

	if len(relayConfig.darkPools) > 0 && errNum == len(relayConfig.darkPools) {
		return fmt.Errorf("could not send order to any dark pool: %v", err)
	}
	return nil
}

// SendOrderFragmentsToDarkOcean will send order fragments to the dark ocean
func SendOrderFragmentsToDarkOcean(order OrderFragments, relayConfig Relay) error {
	valid := false
	for poolIndex := range relayConfig.darkPools {
		fragments := order.DarkPools[GeneratePoolID(relayConfig.darkPools[poolIndex])]
		if fragments != nil && isSafeToSend(len(fragments), relayConfig.darkPools[poolIndex].Size()) {
			if err := sendSharesToDarkPool(relayConfig.darkPools[poolIndex], fragments, relayConfig); err == nil {
				valid = true
			}
		}
	}
	if !valid && len(relayConfig.darkPools) > 0 {
		return fmt.Errorf("cannot send fragments to pools: number of fragments do not match pool size")
	}
	return nil
}

// CancelOrder will cancel orders that aren't confirmed or settled in the dark ocean
func CancelOrder(cancelOrder order.ID, relayConfig Relay) error {
	errCh := make(chan error, len(relayConfig.darkPools))

	orderSignature, err := relayConfig.keyPair.Sign(&cancelOrder)
	if err != nil {
		return err
	}

	go func() {
		defer close(errCh)
		// For every Dark Pool
		for i := range relayConfig.darkPools {
			// Cancel orders for all nodes in the pool
			var wg sync.WaitGroup
			wg.Add(relayConfig.darkPools[i].Size())
			relayConfig.darkPools[i].ForAll(func(n *darknode.Node) {
				defer wg.Done()
				// Get multiaddress
				multiaddress, err := getMultiAddress(n.ID.Address(), relayConfig)
				if err != nil {
					errCh <- fmt.Errorf("cannot read multi-address: %v", err)
					return
				}

				// Create a client
				client, err := rpc.NewClient(context.Background(), multiaddress, relayConfig.multiAddress, relayConfig.keyPair)
				if err != nil {
					errCh <- fmt.Errorf("cannot connect to client: %v", err)
					return
				}

				// TODO: Send signature here?
				cancelOrderRequest := &rpc.CancelOrderRequest{
					From: &rpc.MultiAddress{
						Signature:    []byte{},
						MultiAddress: relayConfig.multiAddress.String(),
					},
					OrderFragmentId: &rpc.OrderFragmentId{
						Signature:       orderSignature,
						OrderFragmentId: cancelOrder,
					},
				}

				// Cancel order
				err = client.CancelOrder(context.Background(), cancelOrderRequest)
				if err != nil {
					errCh <- fmt.Errorf("cannot cancel order to %v", base58.Encode(n.ID))
					return
				}
			})
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

	if len(relayConfig.darkPools) > 0 && errNum == len(relayConfig.darkPools) {
		return fmt.Errorf("could not cancel order to any dark pool: %v", err)
	}
	return nil
}

// Send the shares across all nodes within the Dark Pool
func sendSharesToDarkPool(pool *darknode.Pool, shares []*order.Fragment, relayConfig Relay) error {

	errCh := make(chan error)

	go func() {
		defer close(errCh)

		shareIndex := 0

		var wg sync.WaitGroup
		wg.Add(pool.Size())

		pool.For(func(n *darknode.Node) {
			if shareIndex < len(shares) {
				go func(n *darknode.Node, share *order.Fragment, relayConfig Relay) {
					defer wg.Done()

					shareSignature, err := relayConfig.keyPair.Sign(share)
					if err != nil {
						errCh <- err
					}
					share.Signature = shareSignature

					relaySignature, err := relayConfig.keyPair.Sign(relayConfig.multiAddress)
					if err != nil {
						errCh <- err
					}

					// Get multiaddress
					multiaddress, err := getMultiAddress(n.ID.Address(), relayConfig)
					if err != nil {
						errCh <- fmt.Errorf("cannot read multi-address: %v", err)
						return
					}

					// Create a client
					client, err := rpc.NewClient(context.Background(), multiaddress, relayConfig.multiAddress, relayConfig.keyPair)
					if err != nil {
						errCh <- fmt.Errorf("cannot connect to client: %v", err)
						return
					}

					orderRequest := &rpc.OpenOrderRequest{
						From: &rpc.MultiAddress{
							Signature:    relaySignature,
							MultiAddress: relayConfig.multiAddress.String(),
						},
						OrderFragment: rpc.MarshalOrderFragment(share),
					}

					// Send fragment to node
					err = client.OpenOrder(context.Background(), orderRequest)
					if err != nil {
						errCh <- fmt.Errorf("cannot send order fragment: %v", err)
						return
					}
				}(n, shares[shareIndex], relayConfig)
			}
			shareIndex++
		})
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

// Function to obtain multiaddress of a node by sending requests to bootstrap nodes
func getMultiAddress(address identity.Address, relayConfig Relay) (identity.MultiAddress, error) {
	serializedTarget := rpc.MarshalAddress(address)
	for _, peer := range relayConfig.bootstrapNodes {

		bootStrapMultiAddress, err := identity.NewMultiAddressFromString(peer)
		if err != nil {
			return (identity.MultiAddress{}), err
		}

		client, err := rpc.NewClient(context.Background(), bootStrapMultiAddress, relayConfig.multiAddress, relayConfig.keyPair)
		if err != nil {
			return (identity.MultiAddress{}), fmt.Errorf("cannot establish connection with bootstrap node: %v", err)
		}

		candidates, errs := client.QueryPeersDeep(context.Background(), serializedTarget)

		// TODO: duplicated from swarm.go
		continuing := true
		for continuing {
			select {
			case err := <-errs:
				if err != nil {
					log.Println(fmt.Errorf(fmt.Sprintf("cannot deepen query: %v", err)))
				}
				continuing = false
			case marshaledPeer, ok := <-candidates:
				if !ok {
					continuing = false
					break
				}

				peer, _, err := rpc.UnmarshalMultiAddress(marshaledPeer)
				if err != nil {
					log.Println(fmt.Errorf(fmt.Sprintf("cannot deserialize multiaddress: %v", err)))
					continue
				}
				if address == peer.Address() {
					return peer, nil
				}
			}
		}
	}
	return (identity.MultiAddress{}), nil
}

// GeneratePoolID will generate crypto hash for a pool
func GeneratePoolID(pool *darknode.Pool) string {
	var id identity.ID
	pool.For(func(n *darknode.Node) {
		id = append(id, []byte(n.ID.String())...)
	})
	id = crypto.Keccak256(id)
	return id.String()
}

// isSafeToSend checks if there are enough fragments to successfully complete sending an order
func isSafeToSend(fragmentCount, poolSize int) bool {
	return fragmentCount >= 2/3*poolSize && fragmentCount <= poolSize
}
