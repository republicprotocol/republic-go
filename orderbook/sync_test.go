package orderbook_test

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/orderbook"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/testutils"
)

var (
	NumberOfOrderPairs = 40
	RenLimit           = 10
)

var _ = Describe("Syncer", func() {
	var (
		renLedger   cal.RenLedger
		storer      SyncStorer
		syncer      Syncer
		buys, sells []order.Order
	)

	Context("when syncing with the ledger", func() {

		BeforeEach(func() {
			renLedger = testutils.NewRenLedger()
			storer = testutils.NewStorer()
			buys, sells = generateOrderPairs(NumberOfOrderPairs)

			syncer = NewSyncer(storer, renLedger, RenLimit)
			changeSet, err := syncer.Sync()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(changeSet)).Should(Equal(0))
		})

		It("should be able to sync new opened orders", func() {
			priority := uint64(1)
			for i := 0; i < NumberOfOrderPairs; i++ {
				err := renLedger.OpenBuyOrder([65]byte{}, buys[i].ID)
				Ω(err).ShouldNot(HaveOccurred())
				changeSet, err := syncer.Sync()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(len(changeSet)).Should(Equal(1))
				Ω(bytes.Compare(changeSet[0].OrderID[:], buys[i].ID[:])).Should(Equal(0))
				Ω(changeSet[0].OrderParity).Should(Equal(order.ParityBuy))
				Ω(changeSet[0].OrderPriority).Should(Equal(priority))
				Ω(changeSet[0].OrderStatus).Should(Equal(order.Open))

				err = renLedger.OpenSellOrder([65]byte{}, sells[i].ID)
				Ω(err).ShouldNot(HaveOccurred())
				changeSet, err = syncer.Sync()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(len(changeSet)).Should(Equal(1))
				Ω(bytes.Compare(changeSet[0].OrderID[:], sells[i].ID[:])).Should(Equal(0))
				Ω(changeSet[0].OrderParity).Should(Equal(order.ParitySell))
				Ω(changeSet[0].OrderPriority).Should(Equal(priority))
				Ω(changeSet[0].OrderStatus).Should(Equal(order.Open))
				priority++
			}
		})

		It("should be able to sync confirming order events", func() {
			// Open orders
			openOrders(renLedger, syncer, buys, sells)

			// Confirm orders
			for i := 0; i < NumberOfOrderPairs; i++ {
				err := renLedger.ConfirmOrder(buys[i].ID, sells[i].ID)
				Ω(err).ShouldNot(HaveOccurred())
			}
			changeSet, err := syncer.Sync()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(changeSet)).Should(Equal(NumberOfOrderPairs * 2))
			for i := range changeSet {
				Ω(changeSet[i].OrderStatus).Should(Equal(order.Confirmed))
			}
		})

		It("should be able to sync canceling order events", func() {
			// Open orders
			openOrders(renLedger, syncer, buys, sells)

			// Cancel orders
			for i := 0; i < NumberOfOrderPairs; i++ {
				err := renLedger.CancelOrder([65]byte{}, buys[i].ID)
				Ω(err).ShouldNot(HaveOccurred())
				err = renLedger.CancelOrder([65]byte{}, sells[i].ID)
				Ω(err).ShouldNot(HaveOccurred())
			}
			changeSet, err := syncer.Sync()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(changeSet)).Should(Equal(NumberOfOrderPairs * 2))
			for i := range changeSet {
				Ω(changeSet[i].OrderStatus).Should(Equal(order.Canceled))
			}
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

func openOrders(renLedger cal.RenLedger, syncer Syncer, buys, sells []order.Order) {
	for i := 0; i < NumberOfOrderPairs; i++ {
		err := renLedger.OpenBuyOrder([65]byte{}, buys[i].ID)
		Ω(err).ShouldNot(HaveOccurred())
		err = renLedger.OpenSellOrder([65]byte{}, sells[i].ID)
		Ω(err).ShouldNot(HaveOccurred())
	}
	// Test the renLimit
	for i := 0; i < NumberOfOrderPairs/RenLimit; i++ {
		changeSet, err := syncer.Sync()
		Ω(err).ShouldNot(HaveOccurred())
		Ω(len(changeSet)).Should(Equal(RenLimit * 2))
		Ω(changeSet[i].OrderStatus).Should(Equal(order.Open))
	}
}
