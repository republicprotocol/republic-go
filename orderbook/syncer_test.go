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

		It("should be able to sync opening, confirming and canceling order events", func() {
			done := make(chan struct{})
			defer close(done)

			// Open matching order pairs
			orders := contract.OpenMatchingOrders(NumberOfOrderPairs)

			// Create and start orderbook
			orderbook = NewOrderbook(key, storer.OrderbookPointerStore(), storer.OrderbookOrderStore(), storer.OrderbookOrderFragmentStore(), contract, time.Millisecond, 30)
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
				},
			)

			// Change to first epoch
			_, epoch, err := testutils.RandomEpoch(0)
			Ω(err).ShouldNot(HaveOccurred())
			orderbook.OnChangeEpoch(epoch)

			// Send encrypted order fragments to the orderbook
			for _, ord := range orders {
				fragments, err := ord.Split(5, 4)
				Expect(err).ShouldNot(HaveOccurred())
				encFrag, err := fragments[0].Encrypt(key.PublicKey)
				Expect(err).ShouldNot(HaveOccurred())
				err = orderbook.OpenOrder(context.Background(), encFrag)
				Expect(err).ShouldNot(HaveOccurred())
			}
			time.Sleep(time.Second)

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
			time.Sleep(time.Second)

			// Notifications for all the confirmations must be returned
			// on the notifications channel
			countMu.Lock()
			Expect(countConfirms).Should(Equal(numConfirms))
			countConfirms = 0
			countMu.Unlock()

			// Cancel random orders in the contract
			numCancels := contract.UpdateStatusRandomly(order.Canceled)
			time.Sleep(time.Second)

			// Notifications for all the canceled orders must be returned
			// on the notifications channel
			countMu.Lock()
			Expect(countCancels).Should(Equal(numCancels))
			countMu.Unlock()
		})
	})
})
