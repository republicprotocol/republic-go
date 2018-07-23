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
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/testutils"
)

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

		It("should be able generate the correct number of notifications", func() {
			done := make(chan struct{})
			defer close(done)

			// Open matching order pairs
			orders := contract.OpenMatchingOrders(NumberOfOrderPairs, order.Open)

			// Create and start orderbook
			orderbook = NewOrderbook(key, storer.OrderbookPointerStore(), storer.OrderbookOrderStore(), storer.OrderbookOrderFragmentStore(), contract, time.Millisecond, 80)
			notifications, errs := orderbook.Sync(done)

			// Start reading notifications and errs
			countMu := new(sync.Mutex)
			countOpens := 0
			countCancels := 0
			countConfirms := 0

			go dispatch.CoBegin(
				func() {
					for err := range errs {
						fmt.Println(err)
					}
				},
				func() {
					for notification := range notifications {
						countMu.Lock()
						switch notification.(type) {
						case NotificationOpenOrder:
							countOpens++
						case NotificationConfirmOrder:
							countConfirms++
						case NotificationCancelOrder:
							countCancels++
						}
						countMu.Unlock()
					}
				})

			// Change to first epoch
			_, epoch, err := testutils.RandomEpoch(0)
			Ω(err).ShouldNot(HaveOccurred())
			orderbook.OnChangeEpoch(epoch)

			err = sendOrdersToOrderbook(orders, key, orderbook, 0)
			Ω(err).ShouldNot(HaveOccurred())
			time.Sleep(15 * time.Millisecond)

			// Notifications channel must have emitted open order notifications
			// for all the opened orders
			countMu.Lock()
			Expect(countOpens).Should(Equal(2 * NumberOfOrderPairs))
			Expect(countConfirms).Should(BeZero())
			Expect(countCancels).Should(BeZero())
			countOpens = 0
			countMu.Unlock()

			// Confirm random orders in the contract
			numConfirms := contract.UpdateStatusRandomly(order.Confirmed)
			time.Sleep(15 * time.Millisecond)

			// Notifications for all the confirmations must be returned
			// on the notifications channel
			countMu.Lock()
			Expect(countConfirms).Should(Equal(numConfirms))
			Expect(countOpens).Should(BeZero())
			Expect(countCancels).Should(BeZero())
			countConfirms = 0
			countMu.Unlock()

			// Cancel random orders in the contract
			numCancels := contract.UpdateStatusRandomly(order.Canceled)
			time.Sleep(15 * time.Millisecond)

			// Notifications for all the canceled orders must be returned
			// on the notifications channel
			countMu.Lock()
			Expect(countCancels).Should(Equal(numCancels))
			Expect(countConfirms).Should(BeZero())
			Expect(countOpens).Should(BeZero())
			countMu.Unlock()
		})
	})

	Context("when syncing over multiple epochs", func() {

		It("should not create any open order notifications for orders in a different epoch depth", func() {
			done := make(chan struct{})
			defer close(done)

			// Open matching order pairs
			orders := contract.OpenMatchingOrders(NumberOfOrderPairs, order.Open)

			// Create and start orderbook
			orderbook = NewOrderbook(key, storer.OrderbookPointerStore(), storer.OrderbookOrderStore(), storer.OrderbookOrderFragmentStore(), contract, time.Millisecond, 80)
			notifications, errs := orderbook.Sync(done)

			// Start reading notifications and errs
			countMu := new(sync.Mutex)
			countOpens := 0
			countCancels := 0
			countConfirms := 0

			go dispatch.CoBegin(
				func() {
					for err := range errs {
						fmt.Println(err)
					}
				},
				func() {
					for notification := range notifications {
						countMu.Lock()
						switch notification.(type) {
						case NotificationOpenOrder:
							countOpens++
						case NotificationConfirmOrder:
							countConfirms++
						case NotificationCancelOrder:
							countCancels++
						}
						countMu.Unlock()
					}
				})

			// Change to first epoch
			_, epoch, err := testutils.RandomEpoch(0)
			Ω(err).ShouldNot(HaveOccurred())
			orderbook.OnChangeEpoch(epoch)

			// Send encrypted order fragments at depth 1 to the orderbook
			err = sendOrdersToOrderbook(orders, key, orderbook, 1)
			Ω(err).ShouldNot(HaveOccurred())
			time.Sleep(15 * time.Millisecond)

			// No open order notifications should be created for depth 1
			// in the first epoch
			countMu.Lock()
			Expect(countOpens).Should(BeZero())
			Expect(countConfirms).Should(BeZero())
			Expect(countCancels).Should(BeZero())
			countMu.Unlock()

			confirmedOrders := contract.OpenMatchingOrders(NumberOfOrderPairs/2, order.Confirmed)
			time.Sleep(100 * time.Millisecond)

			// Notifications for all the confirmations must be returned
			// on the notifications channel
			countMu.Lock()
			Expect(countConfirms).Should(Equal(len(confirmedOrders)))
			Expect(countOpens).Should(BeZero())
			Expect(countCancels).Should(BeZero())
			countConfirms = 0
			countMu.Unlock()

			canceledOrders := contract.OpenMatchingOrders(NumberOfOrderPairs/2, order.Canceled)
			time.Sleep(100 * time.Millisecond)

			// Notifications for all the cancelations must be returned
			// on the notifications channel
			countMu.Lock()
			Expect(countCancels).Should(Equal(len(canceledOrders)))
			Expect(countOpens).Should(BeZero())
			Expect(countConfirms).Should(BeZero())
			countCancels = 0
			countMu.Unlock()

			// Change to next epoch
			_, epoch, err = testutils.RandomEpoch(1)
			Ω(err).ShouldNot(HaveOccurred())
			orderbook.OnChangeEpoch(epoch)

			// Send encrypted order fragments at depth 0 to the orderbook
			err = sendOrdersToOrderbook(orders, key, orderbook, 0)
			Ω(err).ShouldNot(HaveOccurred())
			time.Sleep(100 * time.Millisecond)

			// Notifications channel must have emitted open order notifications
			// for all opened fragments at depth 0 in the second epoch
			countMu.Lock()
			Expect(countOpens).Should(Equal(2 * NumberOfOrderPairs))
			Expect(countConfirms).Should(BeZero())
			Expect(countCancels).Should(BeZero())
			countOpens = 0
			countMu.Unlock()

			// Send encrypted order fragments at depth 1 to the orderbook
			err = sendOrdersToOrderbook(orders, key, orderbook, 1)
			Ω(err).ShouldNot(HaveOccurred())
			time.Sleep(100 * time.Millisecond)

			// Notifications channel must have emitted open order notifications
			// for all the opened fragments at depth 1 in the second epoch
			countMu.Lock()
			Expect(countOpens).Should(Equal(2 * NumberOfOrderPairs))
			Expect(countConfirms).Should(BeZero())
			Expect(countCancels).Should(BeZero())
			countOpens = 0
			countMu.Unlock()
		})
	})
})

// Send encrypted order fragments to the orderbook
func sendOrdersToOrderbook(orders []order.Order, key crypto.RsaKey, orderbook Orderbook, depth order.FragmentEpochDepth) error {

	for _, ord := range orders {
		fragments, err := ord.Split(5, 4)
		if err != nil {
			return err
		}
		fragments[0].EpochDepth = depth
		encFrag, err := fragments[0].Encrypt(key.PublicKey)
		if err != nil {
			return err
		}
		err = orderbook.OpenOrder(context.Background(), encFrag)
		if err != nil {
			return err
		}
	}
	return nil
}
