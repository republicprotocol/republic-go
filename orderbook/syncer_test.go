package orderbook_test

import (
	"log"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/orderbook"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/registry"
	"github.com/republicprotocol/republic-go/testutils"
)

var (
	NumberOfOrderPairs = 40
	Limit              = 10
)

var _ = Describe("Syncer", func() {
	var (
		notifications <-chan Notification
		errs          <-chan error
		done          chan struct{}
		orderbook     Orderbook
		contract      *orderbookBinder
		storer        *leveldb.Store
		buys, sells   []order.Order
	)

	BeforeEach(func() {
		var err error
		contract = newOrderbookBinder()
		storer, err = leveldb.NewStore("./data.out", 72*time.Hour)
		Ω(err).ShouldNot(HaveOccurred())
		buys, sells = generateOrderPairs(NumberOfOrderPairs)

		key, err := crypto.RandomRsaKey()
		Ω(err).ShouldNot(HaveOccurred())
		orderbook = NewOrderbook(key, storer.OrderbookPointerStore(), storer.OrderbookOrderStore(), storer.OrderbookOrderFragmentStore(), contract, 72*time.Hour, 100)

		done = make(chan struct{})
		notifications, errs = orderbook.Sync(done)
	})

	AfterEach(func() {
		close(done)
		os.RemoveAll("./data.out")
	})

	Context("when syncing", func() {

		It("should be able to sync new opened orders", func() {
			// priority := Priority(1)

			for i := 0; i < NumberOfOrderPairs; i++ {
				err := contract.OpenBuyOrder([65]byte{}, buys[i].ID)
				Ω(err).ShouldNot(HaveOccurred())

				err = contract.OpenSellOrder([65]byte{}, sells[i].ID)
				Ω(err).ShouldNot(HaveOccurred())
			}
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

		It("should be able to sync confirming order events", func() {
			// Open orders
			openOrders(contract, buys, sells)

			// Confirm orders
			for i := 0; i < NumberOfOrderPairs; i++ {
				err := contract.ConfirmOrder(buys[i].ID, sells[i].ID)
				Ω(err).ShouldNot(HaveOccurred())
			}
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

		FIt("should be able to sync canceling order events", func() {
			// Open orders
			openOrders(contract, buys, sells)

			// Cancel orders
			for i := 0; i < NumberOfOrderPairs; i++ {
				err := contract.CancelOrder([65]byte{}, buys[i].ID)
				Ω(err).ShouldNot(HaveOccurred())
				err = contract.CancelOrder([65]byte{}, sells[i].ID)
				Ω(err).ShouldNot(HaveOccurred())
			}
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

func openOrders(contract *orderbookBinder, buys, sells []order.Order) {
	for i := 0; i < NumberOfOrderPairs; i++ {
		err := contract.OpenBuyOrder([65]byte{}, buys[i].ID)
		Ω(err).ShouldNot(HaveOccurred())
		err = contract.OpenSellOrder([65]byte{}, sells[i].ID)
		Ω(err).ShouldNot(HaveOccurred())
	}
}
