package relay

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/mux"
	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/republic-go/dark"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
)

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
	r.Methods("GET").Path("/orders/{orderID}").Handler(RecoveryHandler(HandleGetOrder()))
	r.Methods("DELETE").Path("/orders/{orderID}").Handler(RecoveryHandler(DeleteOrderHandler(relay.multiAddress, *relay.darkPools)))
	return r
}

// SendOrderToDarkOcean will fragment and send orders to the dark ocean
func SendOrderToDarkOcean(order order.Order, traderMultiAddress *identity.MultiAddress, pools dark.Pools) {
	fmt.Println(order.Type)

	sendOrder(order, pools, *traderMultiAddress)
}

// SendOrderFragmentsToDarkOcean will send order fragments to the dark ocean
func SendOrderFragmentsToDarkOcean(order OrderFragments, traderMultiAddress *identity.MultiAddress, pools dark.Pools) {
	// TODO: (Check) Validate that there are enough fragments for each node in the pool (number of fragments should be atleast 2/3 * size of pool)
	// TODO: Integrate Dark Pool ID here ?
	valid := false
	fragments := order.FragmentSet[0].Fragment
	for i := range pools {
		if len(fragments) >= 2/3*pools[i].Size() {
			valid = true
			sendSharesToDarkPool(pools[i], *traderMultiAddress, fragments)
		}
	}
	if !valid {
		fmt.Errorf("cannot send fragments to pools")
	}
}

// CancelOrder will cancel orders that aren't confirmed or settled in the dark ocean
func CancelOrder(order order.ID, traderMultiAddress *identity.MultiAddress, pools dark.Pools) {
	fmt.Println(order.String())

	// For every Dark Pool
	for i := range pools {

		// Cancel orders for all nodes in the pool
		pools[i].ForAll(func(n *dark.Node) {

			// Get multiaddress
			multiaddress, err := getMultiAddress(n.ID.Address(), *traderMultiAddress)
			if err != nil {
				log.Fatalf("cannot read multi-address: %v", err)
			}

			// Create a client
			client, err := rpc.NewClient(multiaddress, *traderMultiAddress)
			if err != nil {
				log.Fatalf("cannot connect to client: %v", err)
			}

			// Close order
			err = client.CancelOrder(&rpc.OrderSignature{})
			if err != nil {
				log.Println(err)
				log.Printf("%sCoudln't cancel order to %v%s\n", red, base58.Encode(n.ID), reset)
				return
			}
		})
	}
}

func sendOrder(openOrder order.Order, pools dark.Pools, traderMultiAddress identity.MultiAddress) {
	// Buy or sell
	if openOrder.Parity == order.ParityBuy {
		log.Println("sending buy order : ", base58.Encode(openOrder.ID))
	} else {
		log.Println("sending sell order : ", base58.Encode(openOrder.ID))
	}

	var wg sync.WaitGroup
	wg.Add(len(pools))
	// For every dark pool
	for i := range pools {
		go func(darkPool *dark.Pool) {
			defer wg.Done()
			// Split order into (number of nodes in each pool) * 2/3 fragments
			shares, err := openOrder.Split(int64(darkPool.Size()), int64(darkPool.Size()*2/3), &prime)
			if err != nil {
				log.Println("cannot split orders: ", err)
				return
			}
			sendSharesToDarkPool(darkPool, traderMultiAddress, shares)
		}(pools[i])
	}
	wg.Wait()
}

// Send the shares across all nodes within the Dark Pool
func sendSharesToDarkPool(pool *dark.Pool, multi identity.MultiAddress, shares []*order.Fragment) {

	i := 1
	pool.For(func(n *dark.Node) {
		go func(n *dark.Node, multi identity.MultiAddress, share *order.Fragment) {
			// Get multiaddress
			multiaddress, err := getMultiAddress(n.ID.Address(), multi)
			if err != nil {
				log.Fatalf("cannot read multi-address: %v", err)
			}

			// Create a client
			client, err := rpc.NewClient(multiaddress, multi)
			if err != nil {
				log.Fatalf("cannot connect to client: %v", err)
			}

			// Send fragment to node
			err = client.OpenOrder(&rpc.OrderSignature{}, rpc.SerializeOrderFragment(shares[i]))
			if err != nil {
				log.Println(err)
				log.Printf("%sCoudln't send order fragment to %v%s\n", red, base58.Encode(n.ID), reset)
				return
			}
		}(n, multi, shares[i])
		i++
	})
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
