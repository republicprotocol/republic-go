package relay_test

import (
	"context"
	"errors"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/darknodetest"
	. "github.com/republicprotocol/republic-go/relay"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/republic-go/darknode"
	"github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
	"github.com/republicprotocol/republic-go/ethereum/ganache"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/rpc"
	"github.com/republicprotocol/republic-go/rpc/dht"
	"github.com/republicprotocol/republic-go/rpc/relayer"
	"github.com/republicprotocol/republic-go/rpc/smpcer"
	"github.com/republicprotocol/republic-go/rpc/swarmer"
	"github.com/republicprotocol/republic-go/stackint"
	"google.golang.org/grpc"
)

var Prime, _ = stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")
var darknodeRegistry contracts.DarknodeRegistry
var traderMulti string

const (
	GanacheRPC                 = "http://localhost:8545"
	NumberOfDarkNodes          = 5
	NumberOfBootstrapDarkNodes = 5
	NumberOfOrders             = 1
)

var _ = Describe("Relay", func() {

	var conn client.Connection
	var darknodes darknode.Darknodes
	var bootstrapNodes []string

	BeforeSuite(func() {
		var err error

		// Connect to Ganache
		conn, err = ganache.Connect("http://localhost:8545")
		Expect(err).ShouldNot(HaveOccurred())
		darknodeRegistry, err = contracts.NewDarknodeRegistry(context.Background(), conn, ganache.GenesisTransactor(), &bind.CallOpts{})
		Expect(err).ShouldNot(HaveOccurred())
		darknodeRegistry.SetGasLimit(1000000)

		// Create DarkNodes and contexts/cancels for running them
		darknodes, err = NewDarknodes(NumberOfDarkNodes, NumberOfBootstrapDarkNodes)
		Expect(err).ShouldNot(HaveOccurred())

		// Populate bootstrap nodes
		bootstrapNodes = make([]string, 5)
		for i := 0; i < NumberOfBootstrapDarkNodes; i++ {
			bootstrapNodes[i] = darknodes[i].MultiAddress().String()
		}

		// Assign trader address with the first dark node multiaddress
		traderMulti = darknodes[0].MultiAddress().String()

		// Register the Darknodes and trigger an epoch to accept their
		// registrations
		err = RegisterDarknodes(darknodes, conn, darknodeRegistry)
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterSuite(func() {
		var err error

		// Deregister the DarkNodes
		err = DeregisterDarknodes(darknodes, conn, darknodeRegistry)
		Expect(err).ShouldNot(HaveOccurred())

		// Refund the DarkNodes
		err = RefundDarknodes(darknodes, conn, darknodeRegistry)
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("running the relay", func() {
		It("should", func() {
			config := Config{
				KeyPair:      identity.KeyPair{},
				MultiAddress: identity.MultiAddress{},
				Token:        "",
			}
			pools := darknode.NewOcean(darknodeRegistry).GetPools()
			orderbook := orderorderbook.NewOrderbook(100)

			// Initialise DHT using registered nodes
			dht := dht.NewDHT(identity.Address{}, 100)
			for i := 0; i < len(darknodes); i++ {
				dht.UpdateMultiAddress(darknodes[i].MultiAddress)
			}

			connPool := client.NewConnPool(100)
			relayerClient := relayer.NewClient(dht, connPool)
			swarmerClient := swarmer.NewClient(config.MultiAddress, dht, connPool)
			smpcerClient := smpcer.NewClient(config.MultiAddress, connPool)

			relay := NewRelay(config, pools, darknodeRegistry, orderbook, relayerClient, swarmerClient, smpcerClient)

			server := grpc.NewServer()
			relay.Register(server)
			relay.Sync(context.Background(), make([]byte, 32), 3)

			// TODO: Send orders to selected nodes
		})
	})

	Context("storing and updating orders", func() {
		It("should store entry in local orderbook", func() {
			book := orderbook.NewOrderbook(100)
			block := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: []byte{},
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("ID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}

			err := StoreEntryInOrderbook(&block, [32]byte{}, &book)
			Ω(err).ShouldNot(HaveOccurred())

			// Check to see if orderbook is as expected
			blocks := book.Blocks()
			Ω(len(blocks)).Should(Equal(1))
			Ω(blocks[0].Status).Should(Equal(order.Open))
		})

		It("should store multiple entries in local orderbook", func() {
			book := orderbook.NewOrderbook(100)
			fstBlock := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: []byte{},
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("fstID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}
			sndBlock := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: []byte{},
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("sndID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}

			err := StoreEntryInOrderbook(&fstBlock, [32]byte{}, &book)
			Ω(err).ShouldNot(HaveOccurred())
			err = StoreEntryInOrderbook(&sndBlock, [32]byte{}, &book)
			Ω(err).ShouldNot(HaveOccurred())

			// Check to see if orderbook is as expected
			blocks := book.Blocks()
			Ω(len(blocks)).Should(Equal(2))
			Ω(blocks[0].Status).Should(Equal(order.Open))
			Ω(blocks[1].Status).Should(Equal(order.Open))
		})

		It("should update entries with a higher status", func() {
			book := orderbook.NewOrderbook(100)
			openBlock := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: []byte{},
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("ID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}
			confirmedBlock := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: []byte{},
				OrderBlock: &rpc.SyncBlock_Confirmed{
					Confirmed: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("ID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}

			err := StoreEntryInOrderbook(&openBlock, [32]byte{}, &book)
			Ω(err).ShouldNot(HaveOccurred())
			err = StoreEntryInOrderbook(&confirmedBlock, [32]byte{}, &book)
			Ω(err).ShouldNot(HaveOccurred())

			// Check to see if orderbook is as expected
			blocks := book.Blocks()
			Ω(len(blocks)).Should(Equal(1))
			Ω(blocks[0].Status).Should(Equal(order.Confirmed))
		})

		It("should not update entries with a lesser status", func() {
			book := orderbook.NewOrderbook(100)
			confirmedBlock := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: []byte{},
				OrderBlock: &rpc.SyncBlock_Confirmed{
					Confirmed: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("ID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}
			openBlock := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: []byte{},
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("ID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}

			err := StoreEntryInOrderbook(&confirmedBlock, [32]byte{}, &book)
			Ω(err).ShouldNot(HaveOccurred())
			err = StoreEntryInOrderbook(&openBlock, [32]byte{}, &book)
			Ω(err).Should(HaveOccurred())

			// Check to see if orderbook is as expected
			blocks := book.Blocks()
			Ω(len(blocks)).Should(Equal(1))
			Ω(blocks[0].Status).Should(Equal(order.Confirmed))
		})
	})

	Context("forwarding orders", func() {
		It("should forward orders read from the connection", func() {
			// Construct channels
			blocks, errs := make(chan *rpc.SyncBlock), make(chan error)
			defer close(blocks)
			defer close(errs)

			connections := int32(1)
			book := orderbook.NewOrderbook(100)
			fstBlock := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: make([]byte, 32),
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("fstID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}
			sndBlock := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: make([]byte, 32),
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("sndID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				blocks <- &fstBlock
				blocks <- &sndBlock
				errs <- errors.New("connection lost")
			}()

			Ω(len(book.Blocks())).Should(Equal(0))
			err := ForwardMessagesToOrderbook(blocks, errs, &connections, &book)
			Ω(err).Should(HaveOccurred())
			Ω(len(book.Blocks())).Should(Equal(2))
			wg.Wait()
		})

		It("should forward orders from multiple connections", func() {
			// Construct channels
			fstBlocks, fstErrs := make(chan *rpc.SyncBlock), make(chan error)
			sndBlocks, sndErrs := make(chan *rpc.SyncBlock), make(chan error)
			defer close(fstBlocks)
			defer close(fstErrs)
			defer close(sndBlocks)
			defer close(sndErrs)

			connections := int32(2)
			book := orderbook.NewOrderbook(100)
			fstBlock := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: make([]byte, 32),
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("fstID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}
			sndBlock := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: make([]byte, 32),
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("sndID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				fstBlocks <- &fstBlock
				fstErrs <- errors.New("connection lost")
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				sndBlocks <- &sndBlock
				sndErrs <- errors.New("connection lost")
			}()

			Ω(len(book.Blocks())).Should(Equal(0))
			err := ForwardMessagesToOrderbook(fstBlocks, fstErrs, &connections, &book)
			Ω(err).Should(HaveOccurred())
			err = ForwardMessagesToOrderbook(sndBlocks, sndErrs, &connections, &book)
			Ω(err).Should(HaveOccurred())
			Ω(len(book.Blocks())).Should(Equal(2))
			wg.Wait()
		})

		It("should not forward orders with an invalid epoch hash", func() {
			// Construct channels
			blocks, errs := make(chan *rpc.SyncBlock), make(chan error)
			defer close(blocks)
			defer close(errs)

			connections := int32(1)
			book := orderbook.NewOrderbook(100)
			block := rpc.SyncBlock{
				Signature: []byte{},
				Timestamp: 0,
				EpochHash: []byte{},
				OrderBlock: &rpc.SyncBlock_Open{
					Open: &rpc.Order{
						Id: &rpc.OrderId{
							OrderId: []byte("fstID"),
						},
						Type:   0,
						Parity: 0,
						Expiry: 0,
					},
				},
			}

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				blocks <- &block
			}()

			Ω(len(book.Blocks())).Should(Equal(0))
			err := ForwardMessagesToOrderbook(blocks, errs, &connections, &book)
			Ω(err).Should(HaveOccurred())
			Ω(len(book.Blocks())).Should(Equal(0))
			wg.Wait()
		})
	})

	Context("when sending full orders", func() {

		It("should not return an error", func() {
			pools, trader := getPoolsAndTrader(darknodeRegistry)

			keyPair, err := identity.NewKeyPair()
			Expect(err).ShouldNot(HaveOccurred())
			config := Config{
				KeyPair:        keyPair,
				MultiAddress:   trader,
				Token:          "",
				BootstrapNodes: bootstrapNodes,
				BindAddress:    "127.0.0.1:8000",
			}

			relayNode := NewRelay(config, pools, orderbook.NewOrderbook(100), darknodeRegistry)
			sendOrder := getFullOrder()

			err = SendOrderToDarkOcean(sendOrder, relayNode)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("when sending fragmented orders that do not have sufficient fragments", func() {

		It("should return an error", func() {
			pools, trader := getPoolsAndTrader(darknodeRegistry)

			keyPair, err := identity.NewKeyPair()
			Expect(err).ShouldNot(HaveOccurred())
			config := Config{
				KeyPair:        keyPair,
				MultiAddress:   trader,
				Token:          "",
				BootstrapNodes: bootstrapNodes,
				BindAddress:    "127.0.0.1:8000",
			}

			relayNode := NewRelay(config, pools, orderbook.NewOrderbook(100), darknodeRegistry)
			sendOrder := getFragmentedOrder()

			err = SendOrderFragmentsToDarkOcean(sendOrder, relayNode)
			Ω(err).Should(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("number of fragments do not match pool size"))
		})
	})

	Context("when sending fragmented orders that have sufficient fragments for atleast one dark pool", func() {

		It("should not return an error", func() {
			pools, trader := getPoolsAndTrader(darknodeRegistry)

			keyPair, err := identity.NewKeyPair()
			Expect(err).ShouldNot(HaveOccurred())
			config := Config{
				KeyPair:        keyPair,
				MultiAddress:   trader,
				Token:          "",
				BootstrapNodes: bootstrapNodes,
				BindAddress:    "127.0.0.1:8000",
			}

			relayNode := NewRelay(config, pools, orderbook.NewOrderbook(100), darknodeRegistry)

			sendOrder, err := generateFragmentedOrderForDarkPool(pools[0])
			Ω(err).ShouldNot(HaveOccurred())

			err = SendOrderFragmentsToDarkOcean(sendOrder, relayNode)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("when canceling orders", func() {

		It("should not return an error", func() {
			pools, trader := getPoolsAndTrader(darknodeRegistry)

			keyPair, err := identity.NewKeyPair()
			Expect(err).ShouldNot(HaveOccurred())
			config := Config{
				KeyPair:        keyPair,
				MultiAddress:   trader,
				Token:          "",
				BootstrapNodes: bootstrapNodes,
				BindAddress:    "127.0.0.1:8000",
			}

			relayNode := NewRelay(config, pools, orderbook.NewOrderbook(100), darknodeRegistry)
			orderID := []byte("vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM=")

			err = CancelOrder(orderID, relayNode)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})
})

// getPools return dark pools from a mock dnr
func getPools(dnr contracts.DarknodeRegistry) darknode.Pools {
	darknodeIDs, err := dnr.GetAllNodes()
	Ω(err).ShouldNot(HaveOccurred())

	epoch, err := dnr.CurrentEpoch()
	Ω(err).ShouldNot(HaveOccurred())
	darkOcean := darknode.NewDarkOcean(epoch.Blockhash, darknodeIDs)

	return darkOcean.Pools()
}

func getFullOrder() order.Order {
	fullOrder := order.Order{}

	defaultStackVal, err := stackint.FromString("179761232312312")
	Expect(err).ShouldNot(HaveOccurred())

	fullOrder.ID = []byte("vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM=")
	fullOrder.Type = 2
	fullOrder.Parity = 1
	fullOrder.Expiry = time.Time{}
	fullOrder.FstCode = order.CurrencyCodeETH
	fullOrder.SndCode = order.CurrencyCodeBTC
	fullOrder.Price = defaultStackVal
	fullOrder.MaxVolume = defaultStackVal
	fullOrder.MinVolume = defaultStackVal
	fullOrder.Nonce = defaultStackVal
	return fullOrder
}

func getFragmentedOrder() OrderFragments {
	defaultStackVal, _ := stackint.FromString("179761232312312")

	fragmentedOrder := OrderFragments{}
	fragmentSet := map[string][]*order.Fragment{}
	fragments := []*order.Fragment{}

	var err error
	fragments, err = order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, defaultStackVal, defaultStackVal, defaultStackVal, defaultStackVal).Split(2, 1, &Prime)
	Ω(err).ShouldNot(HaveOccurred())

	fragmentSet["vrZhWU3VV9LRIM="] = fragments

	fragmentedOrder.ID = []byte("vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM=")
	fragmentedOrder.Type = 2
	fragmentedOrder.Parity = 1
	fragmentedOrder.Expiry = time.Time{}
	fragmentedOrder.DarkPools = fragmentSet

	return fragmentedOrder
}

func generateFragmentedOrderForDarkPool(pool darknode.Pool) (OrderFragments, error) {
	sendOrder := getFullOrder()
	fragments, err := sendOrder.Split(int64(pool.Size()), int64(pool.Size()*2/3), &Prime)
	if err != nil {
		return OrderFragments{}, err
	}
	fragmentSet := map[string][]*order.Fragment{}
	fragmentOrder := getFragmentedOrder()
	fragmentSet[GeneratePoolID(pool)] = fragments
	fragmentOrder.DarkPools = fragmentSet
	return fragmentOrder, nil
}

func getPoolsAndTrader(darknodeRegistry contracts.DarknodeRegistry) (darknode.Pools, identity.MultiAddress) {
	trader, err := identity.NewMultiAddressFromString(traderMulti)
	Ω(err).ShouldNot(HaveOccurred())

	return getPools(darknodeRegistry), trader
}
