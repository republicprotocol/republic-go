package leveldb_test

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/republic-go/testutils"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/syndtr/goleveldb/leveldb"
)

var orders = make([]order.Order, 100)
var orderFragments = make([]order.Fragment, 100)

const dbFolder = "./tmp/"
const dbFile = dbFolder + "db"
const orderStatus = order.Open

var _ = Describe("Orderbook storage", func() {
	BeforeEach(func() {
		for i := 0; i < 100; i++ {
			ord := order.NewOrder(order.ParityBuy, order.TypeMidpoint, time.Now(), order.SettlementRenEx, testutils.TokensETHREN, uint64(i), uint64(i), uint64(i), uint64(i))
			ordFragments, err := ord.Split(3, 2)
			Expect(err).ShouldNot(HaveOccurred())
			orders[i] = ord
			orderFragments[i] = ordFragments[0]

			// _, err = io.ReadFull(rand.Reader, epoch.Hash[:])
			// Expect(err).ShouldNot(HaveOccurred())

		}
	})

	AfterEach(func() {
		os.RemoveAll(dbFolder)
	})

	Context("when deleting data", func() {
		It("should not retrieve deleted data", func() {
			db := newDB(dbFile)
			orderbookOrderFragmentTable := NewOrderbookOrderFragmentTable(db)
			putAndExpectOrderFragments(orderbookOrderFragmentTable)

			// Attempt to delete and read each of the order fragments
			for i := 0; i < len(orderFragments); i++ {
				err := orderbookOrderFragmentTable.DeleteOrderFragment(orderFragments[i].OrderID)
				Expect(err).ShouldNot(HaveOccurred())
				// Try to read the deleted order fragment
				orderFrag, err := orderbookOrderFragmentTable.OrderFragment(orderFragments[i].OrderID)
				// We should not be able to get the same result back
				Expect(orderFrag.Equal(&orderFragments[i])).Should(BeFalse())
				// We expect a not found error to have occurred
				Expect(err).Should(HaveOccurred())
			}
		})

	})

	Context("when iterating through out of range data", func() {
		It("should trigger an out of range error", func() {
			db := newDB(dbFile)
			orderbookOrderTable := NewOrderbookOrderTable(db)
			orderbookOrderFragmentTable := NewOrderbookOrderFragmentTable(db)

			putAndExpectOrders(orderbookOrderTable)
			putAndExpectOrderFragments(orderbookOrderFragmentTable)

			ordersIter, err := orderbookOrderTable.Orders()
			Expect(err).ShouldNot(HaveOccurred())
			orderFragmentsIter, err := orderbookOrderFragmentTable.OrderFragments()
			Expect(err).ShouldNot(HaveOccurred())
			defer ordersIter.Release()
			defer orderFragmentsIter.Release()

			for ordersIter.Next() {
				_, _, _, _, err = ordersIter.Cursor()
				Expect(err).ShouldNot(HaveOccurred())
			}

			for orderFragmentsIter.Next() {
				_, err = orderFragmentsIter.Cursor()
				Expect(err).ShouldNot(HaveOccurred())
			}

			// These are out of range so we should expect errors
			_, _, _, _, err = ordersIter.Cursor()
			Expect(err).Should(Equal(orderbook.ErrCursorOutOfRange))
			_, err = orderFragmentsIter.Cursor()
			Expect(err).Should(Equal(orderbook.ErrCursorOutOfRange))
		})
	})

	Context("when storing data", func() {

		It("should load data the same data that was stored", func() {
			db := newDB(dbFile)
			orderbookOrderFragmentTable := NewOrderbookOrderFragmentTable(db)
			putAndExpectOrderFragments(orderbookOrderFragmentTable)
		})

		Context("and iterating through", func() {
			It("should load the same amount of data that was stored", func() {
				db := newDB(dbFile)
				orderbookOrderFragmentTable := NewOrderbookOrderFragmentTable(db)
				putAndExpectOrderFragments(orderbookOrderFragmentTable)

				orderFragIter, err := orderbookOrderFragmentTable.OrderFragments()
				Expect(err).ShouldNot(HaveOccurred())
				orderFrags, err := orderFragIter.Collect()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(orderFrags)).Should(Equal(len(orderFragments)))
			})
		})

		Context("when rebooting", func() {
			It("should persist data after reboot", func() {
				db := newDB(dbFile)
				orderbookOrderFragmentTable := NewOrderbookOrderFragmentTable(db)
				putAndExpectOrderFragments(orderbookOrderFragmentTable)

				// Simulate a reboot by closing the database
				err := db.Close()
				Expect(err).ShouldNot(HaveOccurred())

				// Reopen the database and try to read from it
				newDB := newDB(dbFile)
				newOrderbookOrderFragmentTable := NewOrderbookOrderFragmentTable(newDB)
				expectOrderFragments(newOrderbookOrderFragmentTable)
			})
		})
	})

})

func newDB(path string) *leveldb.DB {
	db, err := leveldb.OpenFile(path, nil)
	Expect(err).ShouldNot(HaveOccurred())
	return db
}

func putAndExpectOrders(table *OrderbookOrderTable) {
	for i := 0; i < len(orderFragments); i++ {
		err := table.PutOrder(orders[i].ID, orderStatus, "", uint(i))
		Expect(err).ShouldNot(HaveOccurred())
		Expect(err).ShouldNot(HaveOccurred())
	}
	expectOrders(table)
}

func expectOrders(table *OrderbookOrderTable) {
	for i := 0; i < 100; i++ {
		status, _, _, err := table.Order(orders[i].ID)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(status).Should(Equal(orderStatus))
	}
}

func expectMissingOrders(table *OrderbookOrderTable) {
	for i := 0; i < 100; i++ {
		_, _, _, err := table.Order(orders[i].ID)
		Expect(err).Should(Equal(orderbook.ErrOrderNotFound))
	}
}

func putAndExpectOrderFragments(table *OrderbookOrderFragmentTable) {
	for i := 0; i < len(orderFragments); i++ {
		err := table.PutOrderFragment(orderFragments[i])
		Expect(err).ShouldNot(HaveOccurred())
	}
	expectOrderFragments(table)
}

func expectOrderFragments(table *OrderbookOrderFragmentTable) {
	for i := 0; i < len(orderFragments); i++ {
		orderFrag, err := table.OrderFragment(orderFragments[i].OrderID)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(orderFrag.Equal(&orderFragments[i])).Should(BeTrue())
	}
}

func expectMissingOrderFragments(table *OrderbookOrderFragmentTable) {
	for i := 0; i < len(orderFragments); i++ {
		_, err := table.OrderFragment(orderFragments[i].OrderID)
		Expect(err).Should(Equal(orderbook.ErrOrderFragmentNotFound))
	}
}
