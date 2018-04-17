package relay

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/republic-go/dark"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/stackint"
)

var prime, _ = stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")

// Relay consists of configuration values (?)
type Relay struct {
	multiAddress   identity.MultiAddress
	darkPools      dark.Pools
	bootstrapNodes []string
}

// NewRelay returns a new Relay object
func NewRelay(multi identity.MultiAddress, pools dark.Pools, bootstrapNodes []string) Relay {
	return Relay{
		multiAddress:   multi,
		darkPools:      pools,
		bootstrapNodes: bootstrapNodes,
	}
}

// NewRouter prepares Relay to handle HTTP requests
func NewRouter(relay Relay) *mux.Router {
	orderbook := orderbook.NewOrderbook(100)
	r := mux.NewRouter().StrictSlash(true)
	r.Methods("POST").Path("/orders").Handler(RecoveryHandler(OpenOrdersHandler(relay.multiAddress, relay.darkPools)))
	r.Methods("GET").Path("/orders").Handler(RecoveryHandler(GetOrdersHandler(&orderbook)))
	r.Methods("GET").Path("/orders/{orderID}").Handler(RecoveryHandler(GetOrderHandler(&orderbook, "")))
	r.Methods("DELETE").Path("/orders/{orderID}").Handler(RecoveryHandler(CancelOrderHandler(relay.multiAddress, relay.darkPools)))
	return r
}

// SendOrderToDarkOcean will fragment and send orders to the dark ocean
func SendOrderToDarkOcean(openOrder order.Order, traderMultiAddress identity.MultiAddress, pools dark.Pools, bootstrapNodes []string) error {
	errCh := make(chan error, len(pools))

	go func() {
		defer close(errCh)

		var wg sync.WaitGroup
		wg.Add(len(pools))

		for i := range pools {
			go func(darkPool *dark.Pool) {
				defer wg.Done()
				// Split order into (number of nodes in each pool) * 2/3 fragments
				shares, err := openOrder.Split(int64(darkPool.Size()), int64(darkPool.Size()*2/3), &prime)
				if err != nil {
					errCh <- err
					return
				}
				if err := sendSharesToDarkPool(darkPool, traderMultiAddress, shares, bootstrapNodes); err != nil {
					errCh <- err
					return
				}
			}(pools[i])
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

	if len(pools) > 0 && errNum == len(pools) {
		return fmt.Errorf("could not send order to any dark pool: %v", err)
	}
	return nil
}

// SendOrderFragmentsToDarkOcean will send order fragments to the dark ocean
func SendOrderFragmentsToDarkOcean(order OrderFragments, traderMultiAddress identity.MultiAddress, pools dark.Pools, bootstrapNodes []string) error {
	valid := false
	for poolIndex := range pools {
		fragments := order.DarkPools[GeneratePoolID(pools[poolIndex])]
		if fragments != nil && isSafeToSend(len(fragments), pools[poolIndex].Size()) {
			if err := sendSharesToDarkPool(pools[poolIndex], traderMultiAddress, fragments, bootstrapNodes); err == nil {
				valid = true
			}
		}
	}
	if !valid && len(pools) > 0 {
		return fmt.Errorf("cannot send fragments to pools: number of fragments do not match pool size")
	}
	return nil
}

func CancelOrder(order order.ID, traderMultiAddress identity.MultiAddress, pools dark.Pools, bootstrapNodes []string) error {
	errCh := make(chan error, len(pools))
	go func() {
		defer close(errCh)
		// For every Dark Pool
		for i := range pools {
			// Cancel orders for all nodes in the pool
			var wg sync.WaitGroup
			wg.Add(pools[i].Size())
			pools[i].ForAll(func(n *dark.Node) {
				defer wg.Done()
				// Get multiaddress
				multiaddress, err := getMultiAddress(n.ID.Address(), traderMultiAddress, bootstrapNodes)
				if err != nil {
					errCh <- fmt.Errorf("cannot read multi-address: %v", err)
					return
				}

				// Create a client
				client, err := rpc.NewClient(multiaddress, traderMultiAddress)
				if err != nil {
					errCh <- fmt.Errorf("cannot connect to client: %v", err)
					return
				}

				// Close order
				err = client.CancelOrder(&rpc.OrderSignature{})
				if err != nil {
					errCh <- fmt.Errorf("cannot cancel order to %v", base58.Encode(n.ID))
					return
				}
			})
			wg.Wait()
		}
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
	if len(pools) > 0 && errNum == len(pools) {
		return fmt.Errorf("could not cancel order to any dark pool: %v", err)
	}
	return nil
}

// Send the shares across all nodes within the Dark Pool
func sendSharesToDarkPool(pool *dark.Pool, multi identity.MultiAddress, shares []*order.Fragment, bootstrapNodes []string) error {

	errCh := make(chan error)

	go func() {
		defer close(errCh)

		shareIndex := 0

		var wg sync.WaitGroup
		wg.Add(pool.Size())

		pool.For(func(n *dark.Node) {
			if shareIndex < len(shares) {
				go func(n *dark.Node, multi identity.MultiAddress, share *order.Fragment) {
					defer wg.Done()
					// Get multiaddress
					multiaddress, err := getMultiAddress(n.ID.Address(), multi, bootstrapNodes)
					if err != nil {
						errCh <- fmt.Errorf("cannot read multi-address: %v", err)
						return
					}

					// Create a client
					client, err := rpc.NewClient(multiaddress, multi)
					if err != nil {
						errCh <- fmt.Errorf("cannot connect to client: %v", err)
						return
					}

					// Send fragment to node
					err = client.OpenOrder(&rpc.OrderSignature{}, rpc.SerializeOrderFragment(share))
					if err != nil {
						errCh <- fmt.Errorf("cannot send order fragment: %v", err)
						return
					}
				}(n, multi, shares[shareIndex])
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
func getMultiAddress(address identity.Address, traderMultiAddress identity.MultiAddress, bootstrapNodes []string) (identity.MultiAddress, error) {
	serializedTarget := rpc.SerializeAddress(address)
	for _, peer := range bootstrapNodes {

		bootStrapMultiAddress, err := identity.NewMultiAddressFromString(peer)
		if err != nil {
			return traderMultiAddress, err
		}

		client, err := rpc.NewClient(bootStrapMultiAddress, traderMultiAddress)
		if err != nil {
			return traderMultiAddress, fmt.Errorf("cannot establish connection with bootstrap node: %v", err)
		}

		candidates, err := client.QueryPeersDeep(serializedTarget)
		if err != nil {
			return traderMultiAddress, err
		}

		for candidate := range candidates {
			deserializedCandidate, err := rpc.DeserializeMultiAddress(candidate)
			if err != nil {
				return traderMultiAddress, err
			}
			if address == deserializedCandidate.Address() {
				return deserializedCandidate, nil
			}
		}
	}
	return traderMultiAddress, nil
}

// GeneratePoolID will generate crypto hash for a pool
func GeneratePoolID(pool *dark.Pool) string {
	var id identity.ID
	pool.For(func(n *dark.Node) {
		id = append(id, []byte(n.ID.String())...)
	})
	id = crypto.Keccak256(id)
	return id.String()
}

// isSafeToSend checks if there are enough fragments to successfully complete sending an order
func isSafeToSend(fragmentCount, poolSize int) bool {
	return fragmentCount >= 2/3*poolSize && fragmentCount <= poolSize
}
