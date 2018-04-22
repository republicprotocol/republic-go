package relay

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/darknode"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/rpc"
	"github.com/republicprotocol/republic-go/stackint"
)

var prime, _ = stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")

type Relay struct {
	Config
	DarkPools darknode.Pools
	Orderbook orderbook.Orderbook
	Registrar contracts.DarkNodeRegistry
}

type Config struct {
	KeyPair        identity.KeyPair
	MultiAddress   identity.MultiAddress
	Token          string
	BootstrapNodes []string
	BindAddress    string
}

// NewRelay returns a new Relay object
func NewRelay(config Config, pools darknode.Pools, book orderbook.Orderbook, registrar contracts.DarkNodeRegistry) Relay {
	return Relay{
		Config:    config,
		DarkPools: pools,
		Orderbook: book,
		Registrar: registrar,
	}
}

// NewRouter prepares Relay to handle HTTP requests
func NewRouter(relay Relay) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Methods("POST").Path("/orders").Handler(RecoveryHandler(AuthorizationHandler(OpenOrdersHandler(relay), relay.Config.Token)))
	r.Methods("GET").Path("/orders").Handler(RecoveryHandler(AuthorizationHandler(GetOrdersHandler(&relay.Orderbook), relay.Config.Token)))
	r.Methods("GET").Path("/orders/{orderID}").Handler(RecoveryHandler(AuthorizationHandler(GetOrderHandler(&relay.Orderbook, ""), relay.Config.Token)))
	r.Methods("DELETE").Path("/orders/{orderID}").Handler(RecoveryHandler(AuthorizationHandler(CancelOrderHandler(relay), relay.Config.Token)))
	return r
}

// RunRelay starts the relay
func RunRelay(config Config, pools darknode.Pools, book orderbook.Orderbook, registrar contracts.DarkNodeRegistry) {
	relayNode := NewRelay(config, pools, book, registrar)
	r := NewRouter(relayNode)
	if err := http.ListenAndServe(config.BindAddress, r); err != nil {
		fmt.Println(fmt.Errorf("could not start router: %s", err))
		return
	}

	// Handle orderbook synchronization
	multi, err := identity.NewMultiAddressFromString("/ip4/0.0.0.0/tcp/18415/republic/8MGzNX7M1ucyvtumVj7QYbb7wQ8UTx")
	if err != nil {
		fmt.Println(fmt.Errorf("could not generate multiaddress: %s", err))
		return
	}
	clientPool := rpc.NewClientPool(multi, config.KeyPair).WithTimeout(10 * time.Second).WithTimeoutBackoff(0)
	go SynchronizeOrderbook(&book, clientPool, registrar)
}

// SendOrderToDarkOcean will fragment and send orders to the dark ocean
func SendOrderToDarkOcean(openOrder order.Order, relay Relay) error {
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
			go func(darkPool darknode.Pool) {
				defer wg.Done()
				// Split order into (number of nodes in each pool) * 2/3 fragments
				shares, err := openOrder.Split(int64(darkPool.Size()), int64(darkPool.Size()*2/3), &prime)
				if err != nil {
					errCh <- err
					return
				}
				if err := sendSharesToDarkPool(darkPool, shares, relay); err != nil {
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
func SendOrderFragmentsToDarkOcean(order OrderFragments, relay Relay) error {
	valid := false
	for poolIndex := range relay.DarkPools {
		fragments := order.DarkPools[GeneratePoolID(relay.DarkPools[poolIndex])]
		if fragments != nil && isSafeToSend(len(fragments), relay.DarkPools[poolIndex].Size()) {
			if err := sendSharesToDarkPool(relay.DarkPools[poolIndex], fragments, relay); err == nil {
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
func CancelOrder(cancelOrder order.ID, relay Relay) error {
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

// SynchronizeOrderbook synchronizes the orderbook within the relay using 3
// randomly selected nodes.
func SynchronizeOrderbook(book *orderbook.Orderbook, clientPool *rpc.ClientPool, registrar contracts.DarkNodeRegistry) {
	nodes, err := registrar.GetAllNodes()
	if err != nil {
		fmt.Println(fmt.Errorf("could not retrieve nodes: %s", err))
	}
	connections := int32(0)
	context, cancel := context.WithCancel(context.Background())
	defer cancel() // TODO: Check this
	for {
		// If there are at least 3 connections, we try again in 10 seconds.
		if atomic.LoadInt32(&connections) >= 3 {
			time.Sleep(10 * time.Second)
			break
		}

		// Select a node in a random position and increment the number of
		// connected nodes.
		randIndex := rand.Intn(len(nodes))
		multiaddressString := nodes[randIndex]
		atomic.AddInt32(&connections, 1)

		// Retrieve the multiaddress of the selected node.
		multi, err := identity.NewMultiAddressFromString(string(multiaddressString))
		if err != nil {
			fmt.Println(fmt.Errorf("unable to convert \"%s\" to multiaddress: %s", multiaddressString, err))
		}

		// Check for any messages received from this node and forward them to
		// orderbook stored in the relay.
		blocks, errs := clientPool.Sync(context, multi)
		go func() {
			if err := ForwardMessagesToOrderbook(blocks, errs, &connections, book); err != nil {
				fmt.Println(err)
			}
		}()
	}
}

// ForwardMessagesToOrderbook forwards any messages we receive from the channels
// and stores them in the relay orderbook.
func ForwardMessagesToOrderbook(blocks <-chan *rpc.SyncBlock, errs <-chan error, connections *int32, book *orderbook.Orderbook) error {
	// When the function ends we decrement the total number of connections.
	defer atomic.AddInt32(connections, -1)
	for {
		select {
		case err, ok := <-errs:
			// Output an error and end the connection.
			if !ok || err != nil {
				return fmt.Errorf("error when trying to sync client pool: %s", err)
			}
		case block, ok := <-blocks:
			if !ok {
				return fmt.Errorf("error when trying to sync client pool")
			}

			// The epoch hash we retrieve is stored in a dynamic sized byte
			// array, so we must copy this to one of a fixed length in order
			// to include it in the order entry.
			var epochHash [32]byte
			if len(block.EpochHash) == 32 {
				copy(epochHash[:], block.EpochHash[:32])
			} else {
				return fmt.Errorf("epoch hash is required to be exactly 32 bytes (%d)", len(block.EpochHash))
			}

			// Store this entry in the relay orderbook.
			if err := StoreEntryInOrderbook(block, epochHash, book); err != nil {
				return err
			}
		}
	}
}

// StoreEntryInOrderbook stores any messages we receive from the channels into the
// orderbook provided.
func StoreEntryInOrderbook(block *rpc.SyncBlock, epochHash [32]byte, book *orderbook.Orderbook) error {
	// Check the status of the order message received and call the
	// corresponding function from the orderbook.
	switch block.OrderBlock.(type) {
	case *rpc.SyncBlock_Open:
		ord := rpc.UnmarshalOrder(block.GetOpen())
		entry := orderbook.NewEntry(ord, order.Open, epochHash)
		if err := book.Open(entry); err != nil {
			return fmt.Errorf("error when synchronizing order: %s", err)
		}
	case *rpc.SyncBlock_Confirmed:
		ord := rpc.UnmarshalOrder(block.GetConfirmed())
		entry := orderbook.NewEntry(ord, order.Confirmed, epochHash)
		if err := book.Confirm(entry); err != nil {
			return fmt.Errorf("error when synchronizing order: %s", err)
		}
	case *rpc.SyncBlock_Unconfirmed:
		ord := rpc.UnmarshalOrder(block.GetUnconfirmed())
		entry := orderbook.NewEntry(ord, order.Unconfirmed, epochHash)
		if err := book.Match(entry); err != nil {
			return fmt.Errorf("error when synchronizing order: %s", err)
		}
	case *rpc.SyncBlock_Canceled:
		ord := rpc.UnmarshalOrder(block.GetCanceled())
		entry := orderbook.NewEntry(ord, order.Canceled, epochHash)
		if err := book.Release(entry); err != nil {
			return fmt.Errorf("error when synchronizing order: %s", err)
		}
	case *rpc.SyncBlock_Settled:
		ord := rpc.UnmarshalOrder(block.GetSettled())
		entry := orderbook.NewEntry(ord, order.Settled, epochHash)
		if err := book.Settle(entry); err != nil {
			return fmt.Errorf("error when synchronizing order: %s", err)
		}
	default:
		return fmt.Errorf("unknown order status, %t", block.OrderBlock)
	}

	return nil
}

// Send the shares across all nodes within the Dark Pool
func sendSharesToDarkPool(pool darknode.Pool, shares []*order.Fragment, relay Relay) error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)

		shareIndex := 0

		var wg sync.WaitGroup
		wg.Add(pool.Size())

		// pool.For(func(n darknode.Darknode) {
		for _, node := range pool.Addresses() {
			if shareIndex < len(shares) {
				go func(nodeAddress identity.Address, share *order.Fragment, relay Relay) {
					defer wg.Done()

					shareSignature, err := relay.Config.KeyPair.Sign(share)
					if err != nil {
						errCh <- err
					}
					share.Signature = shareSignature

					relaySignature, err := relay.Config.KeyPair.Sign(relay.Config.MultiAddress)
					if err != nil {
						errCh <- err
					}

					// Get multiaddress
					multiaddress, err := getMultiAddress(nodeAddress, relay)
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

					nodePublicKey, err := relay.Registrar.GetPublicKey([]byte(node.String()))
					if err != nil {
						errCh <- fmt.Errorf("cannot get public key: %v", err)
						return
					}

					nodeRsaKey, err := crypto.BytesToPublicKey(nodePublicKey)
					if err != nil {
						errCh <- fmt.Errorf("cannot convert public key to rsa public key: %v", err)
						return
					}

					orderFragment, err := rpc.MarshalOrderFragment(nodeRsaKey, share)
					if err != nil {
						errCh <- fmt.Errorf("cannot marshal order fragment: %v", err)
						return
					}

					orderRequest := &rpc.OpenOrderRequest{
						From: &rpc.MultiAddress{
							Signature:    relaySignature,
							MultiAddress: relay.Config.MultiAddress.String(),
						},
						OrderFragment: orderFragment,
					}

					// Send fragment to node
					err = client.OpenOrder(context.Background(), orderRequest)
					if err != nil {
						errCh <- fmt.Errorf("cannot send order fragment: %v", err)
						return
					}
				}(node, shares[shareIndex], relay)
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

// TODO: Change this function to not Query and fetch multiaddresses from DHT
// Function to obtain multiaddress of a node by sending requests to bootstrap nodes
func getMultiAddress(address identity.Address, relay Relay) (identity.MultiAddress, error) {

	// serializedTarget := rpc.MarshalAddress(address)
	
	// for _, peer := range relay.Config.BootstrapNodes {

	// 	bootStrapMultiAddress, err := identity.NewMultiAddressFromString(peer)
	// 	if err != nil {
	// 		return (identity.MultiAddress{}), err
	// 	}

	// 	client, err := rpc.NewClient(context.Background(), bootStrapMultiAddress, relay.Config.MultiAddress, relay.Config.KeyPair)
	// 	if err != nil {
	// 		return (identity.MultiAddress{}), fmt.Errorf("cannot establish connection with bootstrap node: %v", err)
	// 	}

	// 	candidates, errs := client.QueryPeersDeep(context.Background(), serializedTarget)
	// 	for err := range errs {
	// 		log.Println(fmt.Sprintf("error finding node: %s", err.Error()))
	// 		continue
	// 	}
	// 	for candidate := range candidates {
	// 		unmarshalCandidate, _, err := rpc.UnmarshalMultiAddress(candidate)
	// 		if err != nil {
	// 			log.Println(fmt.Sprintf("cannot deserialize multiaddress: %s", err.Error()))
	// 			continue
	// 		}
	// 		if address == unmarshalCandidate.Address() {
	// 			log.Println("found address")
	// 			return unmarshalCandidate, nil
	// 		}
	// 	}
	// }
	return (identity.MultiAddress{}), nil
}

// GeneratePoolID will generate crypto hash for a pool
func GeneratePoolID(pool darknode.Pool) string {
	poolID := pool.ID()
	return string(poolID[:])
}

// isSafeToSend checks if there are enough fragments to successfully complete sending an order
func isSafeToSend(fragmentCount, poolSize int) bool {
	return fragmentCount >= 2/3*poolSize && fragmentCount <= poolSize
}
