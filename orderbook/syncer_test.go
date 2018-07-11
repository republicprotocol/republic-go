package orderbook_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/republicprotocol/republic-go/logger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/orderbook"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/testutils"
)

var (
	NumberOfOrderPairs = 40
	Limit              = 10
)

var _ = Describe("Syncer", func() {
	var (
		orderbook   Orderbook
		contract    *testutils.MockContractBinder
		storer      *leveldb.Store
		buys, sells []order.Order
		key         crypto.RsaKey
	)

	BeforeEach(func() {
		var err error
		contract = testutils.NewMockContractBinder()
		storer, err = leveldb.NewStore("./tmp/data.out", 72*time.Hour)
		Ω(err).ShouldNot(HaveOccurred())
		buys, sells = generateOrderPairs(NumberOfOrderPairs)

		key, err = crypto.RandomRsaKey()
		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll("./tmp/data.out")
	})

	Context("when syncing", func() {

		FIt("should be able to sync new opened orders", func() {
			logger.SetFilterLevel(logger.LevelDebug)
			done := make(chan struct{})
			defer close(done)

			orders := contract.OpenMatchingOrders(NumberOfOrderPairs)
			addr, err := testutils.RandomAddress()
			Ω(err).ShouldNot(HaveOccurred())
			epoch := newEpoch(0, addr)

			orderbook = NewOrderbook(key, storer.OrderbookPointerStore(), storer.OrderbookOrderStore(), storer.OrderbookOrderFragmentStore(), contract, time.Millisecond, 100)
			notifications, errs := orderbook.Sync(done)
			orderbook.OnChangeEpoch(epoch)

			fmt.Printf("we have %d orders\n", len(orders))
			for i, ord := range orders {
				fmt.Printf("splitting fragments\n")
				fragments, err := ord.Split(5, 4)
				Expect(err).ShouldNot(HaveOccurred())
				// err = storer.OrderbookOrderFragmentStore().PutOrderFragment(epoch, fragments[0])
				fmt.Printf("encrypting public key\n")
				encFrag, err := fragments[0].Encrypt(key.PublicKey)
				Expect(err).ShouldNot(HaveOccurred())
				fmt.Printf("calling orderbook.open\n")
				err = orderbook.OpenOrder(context.Background(), encFrag)
				fmt.Printf("orderbook.open order returned\n")
				Expect(err).ShouldNot(HaveOccurred())
				fmt.Printf("adding order %d fragments to storer\n", i)
			}
			orderbook.OnChangeEpoch(newEpoch(1, addr))
			time.Sleep(time.Second)

			go func() {
				for err := range errs {
					fmt.Println(err)
				}
			}()

			count := 0
			for notification := range notifications {
				fmt.Println("notification: ")
				fmt.Println(notification)
				count++
			}
			Expect(count).Should(Equal(NumberOfOrderPairs))

		})

		/*
			It("should be able to sync confirming order events", func() {
				// Open orders
				openOrders(contract, buys, sells)

				// Confirm orders
				// for i := 0; i < NumberOfOrderPairs; i++ {
				// err := contract.ConfirmOrder(buys[i].ID, sells[i].ID)
				// Ω(err).ShouldNot(HaveOccurred())
				// }
				orderbook.OnChangeEpoch(registry.Epoch{})

				var count = 0

				go func() {
					for {
						select {
						case <-done:
							return
						case _, ok := <-notifications:
							if !ok {
								return
							}
							count++
							log.Println(count)
						case <-errs:
							return
						}
					}
				}()
			})

			It("should be able to sync canceling order events", func() {
				// Open orders
				openOrders(contract, buys, sells)

				// Cancel orders
				// for i := 0; i < NumberOfOrderPairs; i++ {
				// 	err := contract.CancelOrder([65]byte{}, buys[i].ID)
				// 	Ω(err).ShouldNot(HaveOccurred())
				// 	err = contract.CancelOrder([65]byte{}, sells[i].ID)
				// 	Ω(err).ShouldNot(HaveOccurred())
				// }
				orderbook.OnChangeEpoch(registry.Epoch{})

				var count = 0

				go func() {
					for {
						select {
						case <-done:
							return
						case _, ok := <-notifications:
							if !ok {
								return
							}
							count++
							log.Println(count)
						case <-errs:
							return
						}
					}
				}()
			})
		*/
	})
})

func generateOrderPairs(n int) ([]order.Order, []order.Order) {
	buyOrders := make([]order.Order, n)
	sellOrders := make([]order.Order, n)

	for i := 0; i < n; i++ {
		buyOrders[i] = testutils.RandomBuyOrder()
		sellOrders[i] = testutils.RandomSellOrder()
	}

	return buyOrders, sellOrders
}
