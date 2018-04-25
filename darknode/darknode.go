package darknode

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/darkocean"
	"github.com/republicprotocol/republic-go/delta"
	"github.com/republicprotocol/republic-go/dispatch"
	ethclient "github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
	"github.com/republicprotocol/republic-go/hyperdrive"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/relay"
	"github.com/republicprotocol/republic-go/rpc"
	"github.com/republicprotocol/republic-go/smpc"
	"google.golang.org/grpc"
)

// Darknodes is an alias.
type Darknodes []Darknode

type Darknode struct {
	Config *Config
	Logger *logger.Logger

	id           identity.ID
	address      identity.Address
	multiAddress identity.MultiAddress
	orderbook    orderbook.Orderbook
	crypter      crypto.Crypter

	darknodeRegistry   contracts.DarkNodeRegistry
	hyperdriveContract contracts.HyperdriveContract
	txsToBeFinalized   chan hyperdrive.TxWithTimestamp

	orderFragments chan order.Fragment
	rpc            *rpc.RPC
	smpc           smpc.Smpc
	relay          relay.Relay
}

// NewDarknode returns a new Darknode.
func NewDarknode(multiAddr identity.MultiAddress, config *Config) (Darknode, error) {
	node := Darknode{
		Config: config,
		Logger: logger.StdoutLogger,
	}

	// Get identity information from the Config
	key, err := identity.NewKeyPairFromPrivateKey(node.Config.EcdsaKey.PrivateKey)
	if err != nil {
		return node, fmt.Errorf("cannot get ID from private key: %v", err)
	}
	node.id = key.ID()
	node.address = key.Address()
	node.multiAddress = multiAddr
	node.orderbook = orderbook.NewOrderbook(4)

	// Open a connection to the Ethereum network
	transactOpts := bind.NewKeyedTransactor(config.EcdsaKey.PrivateKey)
	ethclient, err := ethclient.Connect(
		config.Ethereum.URI,
		config.Ethereum.Network,
		config.Ethereum.RepublicTokenAddress,
		config.Ethereum.DarknodeRegistryAddress,
		config.Ethereum.HyperdriveAddress,
	)
	if err != nil {
		return node, err
	}

	// Create bindings to the DarknodeRegistry and Hyperdrive
	darknodeRegistry, err := contracts.NewDarkNodeRegistry(context.Background(), ethclient, transactOpts, &bind.CallOpts{})
	if err != nil {
		return Darknode{}, err
	}
	node.darknodeRegistry = darknodeRegistry
	hyperdriveContract, err  := contracts.NewHyperdriveContract(context.Background(), ethclient, transactOpts, &bind.CallOpts{})
	if err != nil {
		return Darknode{}, err
	}
	node.hyperdriveContract = hyperdriveContract
	node.txsToBeFinalized = make(chan hyperdrive.TxWithTimestamp)

	// FIXME: Use a production Crypter implementation
	weakCrypter := crypto.NewWeakCrypter()
	node.crypter = &weakCrypter

	node.orderFragments = make(chan order.Fragment, 1)
	node.rpc = rpc.NewRPC(node.crypter, node.multiAddress, &node.orderbook)
	node.rpc.OnOpenOrder(func(sig []byte, orderFragment order.Fragment) error {
		node.orderFragments <- orderFragment
		return nil
	})
	node.relay = relay.NewRelay(relay.Config{}, darkocean.Pools{}, darknodeRegistry, &node.orderbook, node.rpc.RelayerClient(), node.rpc.SmpcerClient(), node.rpc.SwarmerClient())

	return node, nil
}

// Bootstrap the Darknode into the swarm network. The Darknode will query all
// reachable nodes for itself, updating its dht.DHT as it connects to other
// nodes. Calls to Darknode.Bootstrap are not blocking, and return a channel of
// errors encountered. Users should not call Darknode.Bootstrap until the
// Darknode is registered, and the its registration is approved.
func (node *Darknode) Bootstrap(ctx context.Context) <-chan error {
	return node.rpc.SwarmerClient().Bootstrap(ctx, node.Config.BootstrapMultiAddresses, -1)
}

// Serve the Darknode services until the done channel is closed.
func (node *Darknode) Serve(done <-chan struct{}) <-chan error {
	errs := make(chan error, 1)

	go func() {
		defer close(errs)

		// Wait until registration is approved
		if err := node.darknodeRegistry.WaitUntilRegistration(node.ID()[:]); err != nil {
			errs <- err
			return
		}

		server := grpc.NewServer()
		listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", node.Config.Host, node.Config.Port))
		if err != nil {
			errs <- err
			return
		}

		node.rpc.Relayer().Register(server)
		node.rpc.Smpcer().Register(server)
		node.rpc.Swarmer().Register(server)

		go func() {
			if err := server.Serve(listener); err != nil {
				errs <- err
				return
			}
		}()
		go func() {
			<-done
			server.Stop()
		}()

		// Bootstrap into the network for 10 seconds maximum
		time.Sleep(time.Second)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		dispatch.Pipe(done, node.Bootstrap(ctx), errs)

		// Run epochs
		dispatch.Pipe(done, node.RunEpochs(done), errs)
	}()

	return errs
}

// RunEpochs will watch for changes to the Ocean and run the secure
// multi-party computation with new Pools. Stops when the done channel is
// closed, and will attempt to recover from errors encountered while
// interacting with the Ocean.
func (node *Darknode) RunEpochs(done <-chan struct{}) <-chan error {
	errs := make(chan error, 1)

	go func() {
		// Maintain multiple done channels so that multiple epochs can be running
		// in parallel
		var prevDone chan struct{}
		var currDone chan struct{}
		defer func() {
			if prevDone != nil {
				close(prevDone)
			}
			if currDone != nil {
				close(currDone)
			}
		}()

		// Looping until the done channel is closed will recover from errors
		// returned by watching the Ocean
		for {
			select {
			case <-done:
				return
			default:
			}

			// Start watching epochs
			epochs, epochErrs := RunEpochWatcher(done, node.darknodeRegistry)
			go dispatch.Pipe(done, epochErrs, errs)

			for quit := false; !quit; {
				select {

				case <-done:
					return

				case err, ok := <-errs:
					if !ok {
						quit = true
						break
					}
					node.Logger.Network(logger.Error, err.Error())

				case epoch, ok := <-epochs:
					if !ok {
						quit = true
						break
					}
					if prevDone != nil {
						close(prevDone)
					}
					prevDone = currDone
					currDone = make(chan struct{})

					darknodeIDs, err := node.darknodeRegistry.GetAllNodes()
					if err != nil {
						// FIXME: Do not skip the epoch. Retry with a backoff.
						errs <- err
						continue
					}

					darkOcean := darkocean.NewDarkOcean(epoch.Blockhash, darknodeIDs)
					deltas, deltaErrs := node.RunEpochProcess(currDone, darkOcean)
					go dispatch.Pipe(done, deltaErrs, errs)
					go func() {
						for delta := range deltas {
							if delta.IsMatch(smpc.Prime) {
								node.Logger.OrderMatch(logger.Info, delta.ID.String(), delta.BuyOrderID.String(), delta.SellOrderID.String())
								err = node.OrderMatchToHyperdrive(delta)
								if err != nil {
									node.Logger.Compute(logger.Error, err.Error())
								}
							}
						}
					}()
				}
			}
		}
	}()

	return errs
}

// ID returns the ID of the Darknode.
func (node *Darknode) ID() identity.ID {
	return node.id
}

// Address returns the Address of the Darknode.
func (node *Darknode) Address() identity.Address {
	return node.address
}

// MultiAddress returns the MultiAddress of the Darknode.
func (node *Darknode) MultiAddress() identity.MultiAddress {
	return node.multiAddress
}

// OnOpenOrder implements the rpc.RelayDelegate interface.
func (node *Darknode) OnOpenOrder(from identity.MultiAddress, orderFragment *order.Fragment) {
	node.orderFragments <- *orderFragment
	entry := orderbook.NewEntry(order.Order{
		ID: orderFragment.OrderID,
	}, order.Open)
	err := node.orderbook.Open(entry)
	if err != nil {
		node.Logger.Compute(logger.Error, err.Error())
	}
}

// OrderMatchToHyperdrive converts an order match into a hyperdrive.Tx and
// forwards it to the Hyperdrive.
func (node *Darknode) OrderMatchToHyperdrive(delta delta.Delta) error {

	// Defensively check that the smpc.Delta is actually a match
	if !delta.IsMatch(smpc.Prime) {
		return errors.New("delta is not an order match")
	}

	// Update the buy/sell orders in the orderbook
	entryBuy := orderbook.NewEntry(order.Order{
		ID: delta.BuyOrderID,
	}, order.Unconfirmed)
	node.orderbook.Match(entryBuy)

	entrySell := orderbook.NewEntry(order.Order{
		ID: delta.SellOrderID,
	}, order.Unconfirmed)
	node.orderbook.Match(entrySell)

	// Convert an order match into a Tx
	tx := hyperdrive.NewTxFromByteSlices(delta.SellOrderID, delta.BuyOrderID)

	_, err := node.hyperdriveContract.SendTx(tx)
	if err != nil {
		return fmt.Errorf("fail to send tx to hyperdrive contract , %s", err)
	}

	node.txsToBeFinalized <- hyperdrive.NewTxWithTimestamp(tx, time.Now())

	return nil
}

func (node *Darknode) WatchForHyperdriveContract(done <-chan struct{}, depth uint64) <-chan error {
	errs := make(chan error, 1)

	go func() {
		defer close(errs)

		watchingList := map[string]hyperdrive.TxWithTimestamp{}

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case tx := <-node.txsToBeFinalized:
				watchingList[string(tx.Tx.Hash)] = tx
			case <-ticker.C:
				for key, tx := range watchingList {
					if time.Now().Before(tx.Timestamp.Add(5 * time.Minute)) {
						finalized := true
						for _, nonce := range tx.Nonces {
							dep, err := node.hyperdriveContract.GetDepth(nonce)
							if err != nil {
								for _, nonce := range tx.Nonces {
									entry := orderbook.Entry{
										Order: order.Order{
											ID: order.ID(nonce),
										},
										Status: order.Unconfirmed,
									}
									node.orderbook.Release(entry)
								}
								finalized = false
								delete(watchingList, key)
								break
							}
							if dep < depth {
								finalized = false
								break
							}
						}

						if finalized {
							node.Logger.Info(fmt.Sprintf("%v has been confirmed by hyperdrive ", tx.Nonces))
							for _, nonce := range tx.Nonces {
								entry := orderbook.Entry{
									Order: order.Order{
										ID: order.ID(nonce),
									},
									Status: order.Confirmed,
								}
								node.orderbook.Confirm(entry)
							}
							delete(watchingList, key)
						}
					} else {
						for _, nonce := range tx.Nonces {
							entry := orderbook.Entry{
								Order: order.Order{
									ID: order.ID(nonce),
								},
								Status: order.Unconfirmed,
							}
							node.orderbook.Release(entry)
						}
						delete(watchingList, key)
					}
				}
			}
		}
	}()

	return errs
}
