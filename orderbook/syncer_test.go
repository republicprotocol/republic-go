package orderbook_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/orderbook"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/republic-go/testutils"
)

var ()

var _ = Describe("Syncer", func() {
	var (
		NumberOfOrderPairs = 40
		orderbook          Orderbook
		contract           *testutils.MockContractBinder
		storer             *leveldb.Store
		key                crypto.RsaKey
	)

	BeforeEach(func() {
		var err error
		contract = testutils.NewMockContractBinder()
		storer, err = leveldb.NewStore("./tmp/data.out", 72*time.Hour)
		Ω(err).ShouldNot(HaveOccurred())

		key, err = crypto.RandomRsaKey()
		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll("./tmp/data.out")
	})

	Context("when syncing", func() {

		It("should be able to sync new opened orders", func() {
			// logger.SetFilterLevel(logger.LevelDebug)
			done := make(chan struct{})
			defer close(done)

			orders := contract.OpenMatchingOrders(NumberOfOrderPairs)
			_, epoch, err := testutils.RandomEpoch(0)
			Ω(err).ShouldNot(HaveOccurred())

			orderbook = NewOrderbook(key, storer.OrderbookPointerStore(), storer.OrderbookOrderStore(), storer.OrderbookOrderFragmentStore(), contract, time.Millisecond, 100)
			notifications, errs := orderbook.Sync(done)

			orderbook.OnChangeEpoch(epoch)

			go func() {
				for err := range errs {
					fmt.Println(err)
				}
			}()

			countMu := new(sync.Mutex)
			count := 0
			go func() {
				for _ = range notifications {
					countMu.Lock()
					count++
					countMu.Unlock()
				}
			}()

			for _, ord := range orders {
				fragments, err := ord.Split(5, 4)
				Expect(err).ShouldNot(HaveOccurred())
				encFrag, err := fragments[0].Encrypt(key.PublicKey)
				Expect(err).ShouldNot(HaveOccurred())
				err = orderbook.OpenOrder(context.Background(), encFrag)
				Expect(err).ShouldNot(HaveOccurred())
			}
			_, newEpoch, err := testutils.RandomEpoch(1)
			Ω(err).ShouldNot(HaveOccurred())

			orderbook.OnChangeEpoch(newEpoch)
			time.Sleep(time.Second)

			countMu.Lock()
			Expect(count).Should(Equal(2 * NumberOfOrderPairs))
			countMu.Unlock()
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
