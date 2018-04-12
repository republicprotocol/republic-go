package relay

import (
	"fmt"
	"log"
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
	multiAddress *identity.MultiAddress
	darkPools    *dark.Pools
}

// NewRouter prepares Relay to handle HTTP requests
func NewRouter() *mux.Router {
	// relay := Relay{}
	orderBook := orderbook.NewOrderBook(100)
	r := mux.NewRouter().StrictSlash(true)
	// r.Methods("POST").Path("/orders").Handler(RecoveryHandler(PostOrdersHandler(*relay.multiAddress, *relay.darkPools)))
	r.Methods("GET").Path("/orders").Handler(RecoveryHandler(GetOrdersHandler(orderBook)))
	r.Methods("GET").Path("/orders/{orderID}").Handler(RecoveryHandler(GetOrderHandler(orderBook)))
	// r.Methods("DELETE").Path("/orders/{orderID}").Handler(RecoveryHandler(DeleteOrderHandler(relay.multiAddress, *relay.darkPools)))
	return r
}

// SendOrderToDarkOcean will fragment and send orders to the dark ocean
func SendOrderToDarkOcean(order order.Order, traderMultiAddress *identity.MultiAddress, pools dark.Pools) error {
	return sendOrder(order, pools, *traderMultiAddress)
}

// SendOrderFragmentsToDarkOcean will send order fragments to the dark ocean
func SendOrderFragmentsToDarkOcean(order Fragments, traderMultiAddress *identity.MultiAddress, pools dark.Pools) error {
	valid := false
	for poolIndex := range pools {
		fragments := order.DarkPools[GeneratePoolID(pools[poolIndex])]
		if fragments != nil && isSafeToSend(len(fragments), pools[poolIndex].Size()) {
			if err := sendSharesToDarkPool(pools[poolIndex], *traderMultiAddress, fragments); err == nil {
				valid = true
			}
		}
	}
	if !valid {
		return fmt.Errorf("cannot send fragments to pools: number of fragments do not match pool size")
	}
	return nil
}

// CancelOrder will cancel orders that aren't confirmed or settled in the dark ocean
func CancelOrder(order order.ID, traderMultiAddress *identity.MultiAddress, pools dark.Pools) error {
	//TODO: Handle errors
	// For every Dark Pool
	for i := range pools {
		// Cancel orders for all nodes in the pool
		var wg sync.WaitGroup
		wg.Add(pools[i].Size())
		pools[i].ForAll(func(n *dark.Node) {
			defer wg.Done()
			// Get multiaddress
			multiaddress, err := getMultiAddress(n.ID.Address(), *traderMultiAddress)
			if err != nil {
				log.Printf("cannot read multi-address: %v", err)
			}

			// Create a client
			client, err := rpc.NewClient(multiaddress, *traderMultiAddress)
			if err != nil {
				log.Printf("cannot connect to client: %v", err)
			}

			// Close order
			err = client.CancelOrder(&rpc.OrderSignature{})
			if err != nil {
				log.Printf("cannot cancel order to %v", base58.Encode(n.ID))
			}
		})
		wg.Wait()
	}
	return nil
}

func sendOrder(openOrder order.Order, pools dark.Pools, traderMultiAddress identity.MultiAddress) error {
	errCh := make(chan error)

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
					log.Printf("cannot split orders: %v", err)
					errCh <- err
					return
				}
				if err := sendSharesToDarkPool(darkPool, traderMultiAddress, shares); err != nil {
					log.Println(err)
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

// Send the shares across all nodes within the Dark Pool
func sendSharesToDarkPool(pool *dark.Pool, multi identity.MultiAddress, shares []*order.Fragment) error {

	shareIndex := 0
	poolID := GeneratePoolID(pool)
	errorCh := make(chan error)
	// canSend := true
	pool.For(func(n *dark.Node) {
		if shareIndex >= len(shares) {
			go func(n *dark.Node, multi identity.MultiAddress, share *order.Fragment) {
				// Get multiaddress
				multiaddress, err := getMultiAddress(n.ID.Address(), multi)
				if err != nil {
					log.Printf("cannot read multi-address: %v", err)
				}

				// Create a client
				client, err := rpc.NewClient(multiaddress, multi)
				if err != nil {
					log.Printf("cannot connect to client: %v", err)
				}

				// Send fragment to node
				err = client.OpenOrder(&rpc.OrderSignature{}, rpc.SerializeOrderFragment(share))
				if err != nil {
					log.Println(err)
					log.Printf("cannot send order fragment to %v%s\n", base58.Encode(n.ID), reset)
					errorCh <- err
					return
				}
			}(n, multi, shares[shareIndex])
		}
		shareIndex++
	})

	// check if atleast 2/3 nodes of the pool has recieved the order fragments
	if pool.Size() > 0 && len(errorCh) > ((1/3)*pool.Size()) {
		return fmt.Errorf("could not send orders to %v nodes (out of %v nodes) in pool %v", len(errorCh), pool.Size(), poolID)
	}
	return nil
}

// Function to obtain multiaddress of a node by sending requests to bootstrap nodes
func getMultiAddress(address identity.Address, traderMultiAddress identity.MultiAddress) (identity.MultiAddress, error) {
	BootstrapMultiAddresses := []string{
		"/ip4/52.77.88.84/tcp/18514/republic/8MGzXN7M1ucxvtumVjQ7Ybb7xQ8TUw",
		"/ip4/52.79.194.108/tcp/18514/republic/8MGBUdoFFd8VsfAG5bQSAptyjKuutE",
		"/ip4/52.59.176.141/tcp/18514/republic/8MHmrykz65HimBPYaVgm8bTSpRUoXA",
		"/ip4/52.21.44.236/tcp/18514/republic/8MKFT9CDQQru1hYqnaojXqCQU2Mmuk",
		"/ip4/52.41.118.171/tcp/18514/republic/8MGb8k337pp2GSh6yG8iv2GK6FbNHN",
	}

	serializedTarget := rpc.SerializeAddress(address)
	for _, peer := range BootstrapMultiAddresses {

		bootStrapMultiAddress, err := identity.NewMultiAddressFromString(peer)
		if err != nil {
			return traderMultiAddress, err
		}

		client, err := rpc.NewClient(bootStrapMultiAddress, traderMultiAddress)
		if err != nil {
			return traderMultiAddress, err
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
				fmt.Println("Found the target : ", address)
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

func isSafeToSend(fragmentCount, poolSize int) bool {
	return fragmentCount >= 2/3*poolSize && fragmentCount <= poolSize
}
