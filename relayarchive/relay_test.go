package relay_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/darknode"
	. "github.com/republicprotocol/republic-go/relay"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/stackint"
)

var Prime, _ = stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")

var _ = Describe("Relay", func() {

	// 	Context("running the relay", func() {
	// 		It("should", func() {
	// 			config := Config{
	// 				KeyPair:      identity.KeyPair{},
	// 				MultiAddress: identity.MultiAddress{},
	// 				Token:        "",
	// 			}
	// 			pools := getPools(darknodeRegistry)
	// 			book := orderbook.NewOrderbook(100)

	// 			// Initialise DHT using registered nodes
	// 			address, _, _ := identity.NewAddress()
	// 			dht := dht.NewDHT(address, 100)
	// 			for i := 0; i < len(darknodes); i++ {
	// 				dht.UpdateMultiAddress(darknodes[i].MultiAddress())
	// 			}

	// 			crypter := crypto.NewWeakCrypter()
	// 			connPool := client.NewConnPool(100)
	// 			relayerClient := relayer.NewClient(&crypter, &dht, &connPool)
	// 			swarmerClient := swarmer.NewClient(&crypter, config.MultiAddress, &dht, &connPool)
	// 			smpcerClient := smpcer.NewClient(&crypter, config.MultiAddress, &connPool)

	// 			relay := NewRelay(config, darknodeRegistry, &book, &relayerClient, &smpcerClient, &swarmerClient)

	// 			server := grpc.NewServer()
	// 			relay.Register(server)
	// 			relay.Sync(context.Background(), 3)

	// 			// TODO: Send orders to selected nodes
	// 		})
	// 	})

	// 	Context("storing and updating orders", func() {
	// 		It("should store entry in local orderbook", func() {
	// 			book := orderbook.NewOrderbook(100)
	// 			block := rpc.SyncBlock{
	// 				Signature: []byte{},
	// 				Timestamp: 0,
	// 				EpochHash: []byte{},
	// 				OrderBlock: &rpc.SyncBlock_Open{
	// 					Open: &rpc.Order{
	// 						Id: &rpc.OrderId{
	// 							OrderId: []byte("ID"),
	// 						},
	// 						Type:   0,
	// 						Parity: 0,
	// 						Expiry: 0,
	// 					},
	// 				},
	// 			}

	// 			err := StoreEntryInOrderbook(&block, &book)
	// 			Ω(err).ShouldNot(HaveOccurred())

	// 			// Check to see if orderbook is as expected
	// 			blocks := book.Blocks()
	// 			Ω(len(blocks)).Should(Equal(1))
	// 			Ω(blocks[0].Status).Should(Equal(order.Open))
	// 		})

	// 		It("should store multiple entries in local orderbook", func() {
	// 			book := orderbook.NewOrderbook(100)
	// 			fstBlock := rpc.SyncBlock{
	// 				Signature: []byte{},
	// 				Timestamp: 0,
	// 				EpochHash: []byte{},
	// 				OrderBlock: &rpc.SyncBlock_Open{
	// 					Open: &rpc.Order{
	// 						Id: &rpc.OrderId{
	// 							OrderId: []byte("fstID"),
	// 						},
	// 						Type:   0,
	// 						Parity: 0,
	// 						Expiry: 0,
	// 					},
	// 				},
	// 			}
	// 			sndBlock := rpc.SyncBlock{
	// 				Signature: []byte{},
	// 				Timestamp: 0,
	// 				EpochHash: []byte{},
	// 				OrderBlock: &rpc.SyncBlock_Open{
	// 					Open: &rpc.Order{
	// 						Id: &rpc.OrderId{
	// 							OrderId: []byte("sndID"),
	// 						},
	// 						Type:   0,
	// 						Parity: 0,
	// 						Expiry: 0,
	// 					},
	// 				},
	// 			}

	// 			err := StoreEntryInOrderbook(&fstBlock, &book)
	// 			Ω(err).ShouldNot(HaveOccurred())
	// 			err = StoreEntryInOrderbook(&sndBlock, &book)
	// 			Ω(err).ShouldNot(HaveOccurred())

	// 			// Check to see if orderbook is as expected
	// 			blocks := book.Blocks()
	// 			Ω(len(blocks)).Should(Equal(2))
	// 			Ω(blocks[0].Status).Should(Equal(order.Open))
	// 			Ω(blocks[1].Status).Should(Equal(order.Open))
	// 		})

	// 		It("should update entries with a higher status", func() {
	// 			book := orderbook.NewOrderbook(100)
	// 			openBlock := rpc.SyncBlock{
	// 				Signature: []byte{},
	// 				Timestamp: 0,
	// 				EpochHash: []byte{},
	// 				OrderBlock: &rpc.SyncBlock_Open{
	// 					Open: &rpc.Order{
	// 						Id: &rpc.OrderId{
	// 							OrderId: []byte("ID"),
	// 						},
	// 						Type:   0,
	// 						Parity: 0,
	// 						Expiry: 0,
	// 					},
	// 				},
	// 			}
	// 			confirmedBlock := rpc.SyncBlock{
	// 				Signature: []byte{},
	// 				Timestamp: 0,
	// 				EpochHash: []byte{},
	// 				OrderBlock: &rpc.SyncBlock_Confirmed{
	// 					Confirmed: &rpc.Order{
	// 						Id: &rpc.OrderId{
	// 							OrderId: []byte("ID"),
	// 						},
	// 						Type:   0,
	// 						Parity: 0,
	// 						Expiry: 0,
	// 					},
	// 				},
	// 			}

	// 			err := StoreEntryInOrderbook(&openBlock, &book)
	// 			Ω(err).ShouldNot(HaveOccurred())
	// 			err = StoreEntryInOrderbook(&confirmedBlock, &book)
	// 			Ω(err).ShouldNot(HaveOccurred())

	// 			// Check to see if orderbook is as expected
	// 			blocks := book.Blocks()
	// 			Ω(len(blocks)).Should(Equal(1))
	// 			Ω(blocks[0].Status).Should(Equal(order.Confirmed))
	// 		})

	// 		It("should not update entries with a lesser status", func() {
	// 			book := orderbook.NewOrderbook(100)
	// 			confirmedBlock := rpc.SyncBlock{
	// 				Signature: []byte{},
	// 				Timestamp: 0,
	// 				EpochHash: []byte{},
	// 				OrderBlock: &rpc.SyncBlock_Confirmed{
	// 					Confirmed: &rpc.Order{
	// 						Id: &rpc.OrderId{
	// 							OrderId: []byte("ID"),
	// 						},
	// 						Type:   0,
	// 						Parity: 0,
	// 						Expiry: 0,
	// 					},
	// 				},
	// 			}
	// 			openBlock := rpc.SyncBlock{
	// 				Signature: []byte{},
	// 				Timestamp: 0,
	// 				EpochHash: []byte{},
	// 				OrderBlock: &rpc.SyncBlock_Open{
	// 					Open: &rpc.Order{
	// 						Id: &rpc.OrderId{
	// 							OrderId: []byte("ID"),
	// 						},
	// 						Type:   0,
	// 						Parity: 0,
	// 						Expiry: 0,
	// 					},
	// 				},
	// 			}

	// 			err := StoreEntryInOrderbook(&confirmedBlock, [32]byte{}, &book)
	// 			Ω(err).ShouldNot(HaveOccurred())
	// 			err = StoreEntryInOrderbook(&openBlock, [32]byte{}, &book)
	// 			Ω(err).Should(HaveOccurred())

	// 			// Check to see if orderbook is as expected
	// 			blocks := book.Blocks()
	// 			Ω(len(blocks)).Should(Equal(1))
	// 			Ω(blocks[0].Status).Should(Equal(order.Confirmed))
	// 		})
	// 	})

	// 	Context("forwarding orders", func() {
	// 		It("should forward orders read from the connection", func() {
	// 			// Construct channels
	// 			blocks, errs := make(chan *rpc.SyncBlock), make(chan error)
	// 			defer close(blocks)
	// 			defer close(errs)

	// 			connections := int32(1)
	// 			book := orderbook.NewOrderbook(100)
	// 			fstBlock := rpc.SyncBlock{
	// 				Signature: []byte{},
	// 				Timestamp: 0,
	// 				EpochHash: make([]byte, 32),
	// 				OrderBlock: &rpc.SyncBlock_Open{
	// 					Open: &rpc.Order{
	// 						Id: &rpc.OrderId{
	// 							OrderId: []byte("fstID"),
	// 						},
	// 						Type:   0,
	// 						Parity: 0,
	// 						Expiry: 0,
	// 					},
	// 				},
	// 			}
	// 			sndBlock := rpc.SyncBlock{
	// 				Signature: []byte{},
	// 				Timestamp: 0,
	// 				EpochHash: make([]byte, 32),
	// 				OrderBlock: &rpc.SyncBlock_Open{
	// 					Open: &rpc.Order{
	// 						Id: &rpc.OrderId{
	// 							OrderId: []byte("sndID"),
	// 						},
	// 						Type:   0,
	// 						Parity: 0,
	// 						Expiry: 0,
	// 					},
	// 				},
	// 			}

	// 			var wg sync.WaitGroup
	// 			wg.Add(1)
	// 			go func() {
	// 				defer wg.Done()
	// 				blocks <- &fstBlock
	// 				blocks <- &sndBlock
	// 				errs <- errors.New("connection lost")
	// 			}()

	// 			Ω(len(book.Blocks())).Should(Equal(0))
	// 			err := ForwardMessagesToOrderbook(blocks, errs, &connections, &book)
	// 			Ω(err).Should(HaveOccurred())
	// 			Ω(len(book.Blocks())).Should(Equal(2))
	// 			wg.Wait()
	// 		})

	// 		It("should forward orders from multiple connections", func() {
	// 			// Construct channels
	// 			fstBlocks, fstErrs := make(chan *rpc.SyncBlock), make(chan error)
	// 			sndBlocks, sndErrs := make(chan *rpc.SyncBlock), make(chan error)
	// 			defer close(fstBlocks)
	// 			defer close(fstErrs)
	// 			defer close(sndBlocks)
	// 			defer close(sndErrs)

	// 			connections := int32(2)
	// 			book := orderbook.NewOrderbook(100)
	// 			fstBlock := rpc.SyncBlock{
	// 				Signature: []byte{},
	// 				Timestamp: 0,
	// 				EpochHash: make([]byte, 32),
	// 				OrderBlock: &rpc.SyncBlock_Open{
	// 					Open: &rpc.Order{
	// 						Id: &rpc.OrderId{
	// 							OrderId: []byte("fstID"),
	// 						},
	// 						Type:   0,
	// 						Parity: 0,
	// 						Expiry: 0,
	// 					},
	// 				},
	// 			}
	// 			sndBlock := rpc.SyncBlock{
	// 				Signature: []byte{},
	// 				Timestamp: 0,
	// 				EpochHash: make([]byte, 32),
	// 				OrderBlock: &rpc.SyncBlock_Open{
	// 					Open: &rpc.Order{
	// 						Id: &rpc.OrderId{
	// 							OrderId: []byte("sndID"),
	// 						},
	// 						Type:   0,
	// 						Parity: 0,
	// 						Expiry: 0,
	// 					},
	// 				},
	// 			}

	// 			var wg sync.WaitGroup
	// 			wg.Add(1)
	// 			go func() {
	// 				defer wg.Done()
	// 				fstBlocks <- &fstBlock
	// 				fstErrs <- errors.New("connection lost")
	// 			}()

	// 			wg.Add(1)
	// 			go func() {
	// 				defer wg.Done()
	// 				sndBlocks <- &sndBlock
	// 				sndErrs <- errors.New("connection lost")
	// 			}()

	// 			Ω(len(book.Blocks())).Should(Equal(0))
	// 			err := ForwardMessagesToOrderbook(fstBlocks, fstErrs, &connections, &book)
	// 			Ω(err).Should(HaveOccurred())
	// 			err = ForwardMessagesToOrderbook(sndBlocks, sndErrs, &connections, &book)
	// 			Ω(err).Should(HaveOccurred())
	// 			Ω(len(book.Blocks())).Should(Equal(2))
	// 			wg.Wait()
	// 		})

	// 		It("should not forward orders with an invalid epoch hash", func() {
	// 			// Construct channels
	// 			blocks, errs := make(chan *rpc.SyncBlock), make(chan error)
	// 			defer close(blocks)
	// 			defer close(errs)

	// 			connections := int32(1)
	// 			book := orderbook.NewOrderbook(100)
	// 			block := rpc.SyncBlock{
	// 				Signature: []byte{},
	// 				Timestamp: 0,
	// 				EpochHash: []byte{},
	// 				OrderBlock: &rpc.SyncBlock_Open{
	// 					Open: &rpc.Order{
	// 						Id: &rpc.OrderId{
	// 							OrderId: []byte("fstID"),
	// 						},
	// 						Type:   0,
	// 						Parity: 0,
	// 						Expiry: 0,
	// 					},
	// 				},
	// 			}

	// 			var wg sync.WaitGroup
	// 			wg.Add(1)
	// 			go func() {
	// 				defer wg.Done()
	// 				blocks <- &block
	// 			}()

	// 			Ω(len(book.Blocks())).Should(Equal(0))
	// 			err := ForwardMessagesToOrderbook(blocks, errs, &connections, &book)
	// 			Ω(err).Should(HaveOccurred())
	// 			Ω(len(book.Blocks())).Should(Equal(0))
	// 			wg.Wait()
	// 		})
	// 	})

	Context("when sending full orders", func() {

		It("should not return an error", func() {
			fullOrder, err := darknode.CreateOrders(1, true)
			Ω(err).ShouldNot(HaveOccurred())
			err = relayTestNetEnv.Relays[0].SendOrderToDarkOcean(*fullOrder[0])
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("when sending fragmented orders that do not have sufficient fragments", func() {

		It("should return an error", func() {
			fragmentOrder := getFragmentedOrder()
			err := relayTestNetEnv.Relays[0].SendOrderFragmentsToDarkOcean(fragmentOrder)
			
			Ω(err).Should(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("number of fragments do not match pool size"))
		})
	})

	Context("when sending fragmented orders that have sufficient fragments for atleast one dark pool", func() {

		It("should not return an error", func() {
			pools, err := getPools(darknodeTestnetEnv.DarknodeRegistry)
			Ω(err).ShouldNot(HaveOccurred())
			fragmentedOrder, err := generateFragmentedOrderForDarkPool(pools[0])
			Ω(err).ShouldNot(HaveOccurred())

			err = relayTestNetEnv.Relays[0].SendOrderFragmentsToDarkOcean(fragmentedOrder)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	// 	Context("when canceling orders", func() {

	// 		It("should not return an error", func() {
	// 			pools := getPools(darknodeRegistry)
	// 			trader := identity.MultiAddress{}

	// 			keyPair, err := identity.NewKeyPair()
	// 			Expect(err).ShouldNot(HaveOccurred())
	// 			config := Config{
	// 				KeyPair:        keyPair,
	// 				MultiAddress:   trader,
	// 				Token:          "",
	// 				BootstrapNodes: bootstrapNodes,
	// 				BindAddress:    "127.0.0.1:8000",
	// 			}

	// 			relayNode := NewRelay(config, orderbook.NewOrderbook(100), darknodeRegistry)
	// 			orderID := []byte("vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM=")

	// 			err = CancelOrder(orderID, relayNode)
	// 			Ω(err).ShouldNot(HaveOccurred())
	// 		})
	// 	})
})

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

