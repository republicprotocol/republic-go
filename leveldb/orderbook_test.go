package leveldb_test

import (
	"crypto/rand"
	"io"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/registry"
	"github.com/syndtr/goleveldb/leveldb"
)

var _ = Describe("LevelDB storage", func() {
	orderFragments := make([]order.Fragment, 100)
	epoch := registry.Epoch{}
	dbFolder := "./tmp/"
	dbFile := "db"

	BeforeEach(func() {
		for i := 0; i < 100; i++ {
			ord := order.NewOrder(order.TypeMidpoint, order.ParityBuy, order.SettlementRenEx, time.Now(), order.TokensETHREN, order.NewCoExp(200, 26), order.NewCoExp(200, 26), order.NewCoExp(200, 26), uint64(i))
			ordFragments, err := ord.Split(3, 2)
			Expect(err).ShouldNot(HaveOccurred())
			orderFragments[i] = ordFragments[0]

			_, err = io.ReadFull(rand.Reader, epoch.Hash[:])
			Expect(err).ShouldNot(HaveOccurred())

		}
	})

	AfterEach(func() {
		os.RemoveAll(dbFolder)
	})

	Context("when deleting data", func() {
		It("should not retrieve deleted data", func() {
			db, err := leveldb.OpenFile(filepath.Join(dbFolder, dbFile), nil)
			Expect(err).ShouldNot(HaveOccurred())

			orderbookOrderFragmentTable := NewOrderbookOrderFragmentTable(db)

			// Put the order fragments into the table and attempt to retrieve
			for i := 0; i < len(orderFragments); i++ {
				err := orderbookOrderFragmentTable.PutOrderFragment(epoch, orderFragments[i])
				Expect(err).ShouldNot(HaveOccurred())
				orderFrag, err := orderbookOrderFragmentTable.OrderFragment(epoch, orderFragments[i].OrderID)
				Expect(err).ShouldNot(HaveOccurred())
				// We should be able to get the same result back
				Expect(orderFrag.Equal(&orderFragments[i])).Should(BeTrue())
				err = orderbookOrderFragmentTable.DeleteOrderFragment(epoch, orderFrag.OrderID)
				Expect(err).ShouldNot(HaveOccurred())
				// Try to read the deleted order fragment
				orderFrag, err = orderbookOrderFragmentTable.OrderFragment(epoch, orderFragments[i].OrderID)
				// We should not be able to get the same result back
				Expect(orderFrag.Equal(&orderFragments[i])).Should(BeFalse())
				// We expect a not found error to have occurred
				Expect(err).Should(HaveOccurred())
			}
		})

	})

	Context("when storing data", func() {

		It("should load data the same data that was stored", func() {
			db, err := leveldb.OpenFile(filepath.Join(dbFolder, dbFile), nil)
			Expect(err).ShouldNot(HaveOccurred())

			orderbookOrderFragmentTable := NewOrderbookOrderFragmentTable(db)

			// Put the order fragments into the table and attempt to retrieve
			for i := 0; i < len(orderFragments); i++ {
				err = orderbookOrderFragmentTable.PutOrderFragment(epoch, orderFragments[i])
				Expect(err).ShouldNot(HaveOccurred())
				orderFrag, err := orderbookOrderFragmentTable.OrderFragment(epoch, orderFragments[i].OrderID)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(orderFrag.Equal(&orderFragments[i])).Should(BeTrue())
			}
		})

		Context("and iterating through", func() {
			It("should load the same amount of data that was stored", func() {
				db, err := leveldb.OpenFile(filepath.Join(dbFolder, dbFile), nil)
				Expect(err).ShouldNot(HaveOccurred())

				orderbookOrderFragmentTable := NewOrderbookOrderFragmentTable(db)

				// Put the order fragments into the table and attempt to retrieve
				for i := 0; i < len(orderFragments); i++ {
					err := orderbookOrderFragmentTable.PutOrderFragment(epoch, orderFragments[i])
					Expect(err).ShouldNot(HaveOccurred())
				}

				orderFragIter, err := orderbookOrderFragmentTable.OrderFragments(epoch)
				Expect(err).ShouldNot(HaveOccurred())
				orderFrags, err := orderFragIter.Collect()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(orderFrags)).Should(Equal(len(orderFragments)))
			})
		})

		Context("when rebooting", func() {
			It("should persist data after reboot", func() {
				db, err := leveldb.OpenFile(filepath.Join(dbFolder, dbFile), nil)
				Expect(err).ShouldNot(HaveOccurred())

				orderbookOrderFragmentTable := NewOrderbookOrderFragmentTable(db)

				// Put the order fragments into the table and attempt to retrieve
				for i := 0; i < len(orderFragments); i++ {
					err = orderbookOrderFragmentTable.PutOrderFragment(epoch, orderFragments[i])
					Expect(err).ShouldNot(HaveOccurred())
					orderFrag, err := orderbookOrderFragmentTable.OrderFragment(epoch, orderFragments[i].OrderID)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(orderFrag.Equal(&orderFragments[i])).Should(BeTrue())
				}

				// Simulate a reboot by closing the database
				err = db.Close()
				Expect(err).ShouldNot(HaveOccurred())

				// Reopen the database and try to read from it
				newDB, err := leveldb.OpenFile(filepath.Join(dbFolder, dbFile), nil)
				Expect(err).ShouldNot(HaveOccurred())

				newOrderbookOrderFragmentTable := NewOrderbookOrderFragmentTable(newDB)
				for i := 0; i < len(orderFragments); i++ {
					orderFrag, err := newOrderbookOrderFragmentTable.OrderFragment(epoch, orderFragments[i].OrderID)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(orderFrag.Equal(&orderFragments[i])).Should(BeTrue())
				}

			})
		})

	})

})
