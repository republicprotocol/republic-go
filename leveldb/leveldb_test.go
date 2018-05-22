package leveldb_test

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/leveldb"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
)

var _ = Describe("LevelDB storage", func() {

	orders := make([]order.Order, 100)
	orderFragments := make([]order.Fragment, 100)

	BeforeEach(func() {
		for i := 0; i < 100; i++ {
			ord := order.NewOrder(order.TypeMidpoint, order.ParityBuy, time.Now(), order.TokensETHREN, order.NewCoExp(200, 26), order.NewCoExp(200, 26), order.NewCoExp(200, 26), int64(i))
			ordFragments, err := ord.Split(3, 2)
			Expect(err).ShouldNot(HaveOccurred())
			orders[i] = ord
			orderFragments[i] = ordFragments[0]
		}
	})

	AfterEach(func() {
		os.RemoveAll("./tmp/")
	})

	Context("when storing, loading, and removing data", func() {

		It("should return an error when loading data before storing it", func() {
			db, err := NewStore("./tmp")
			Expect(err).ShouldNot(HaveOccurred())
			for i := 0; i < 100; i++ {
				_, err = db.OrderFragment(orders[i].ID)
				Expect(err).Should(Equal(orderbook.ErrOrderFragmentNotFound))
				_, err = db.Order(orders[i].ID)
				Expect(err).Should(Equal(orderbook.ErrOrderNotFound))
			}
			err = db.Close()
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should not return an error when removing data before storing it", func() {
			db, err := NewStore("./tmp")
			Expect(err).ShouldNot(HaveOccurred())
			for i := 0; i < 100; i++ {
				err = db.RemoveOrderFragment(orders[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
				err = db.RemoveOrder(orders[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
			}
			err = db.Close()
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should load data equal to when it was stored", func() {
			db, err := NewStore("./tmp")
			Expect(err).ShouldNot(HaveOccurred())
			for i := 0; i < 100; i++ {
				err = db.InsertOrderFragment(orderFragments[i])
				Expect(err).ShouldNot(HaveOccurred())
				err = db.InsertOrder(orders[i])
				Expect(err).ShouldNot(HaveOccurred())
			}
			for i := 0; i < 100; i++ {
				orderFragment, err := db.OrderFragment(orders[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
				order, err := db.Order(orders[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(orderFragment.Equal(&orderFragments[i])).Should(BeTrue())
				Expect(order.Equal(&orders[i])).Should(BeTrue())
			}
			err = db.Close()
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return an error when loading data after removing it", func() {
			db, err := NewStore("./tmp")
			Expect(err).ShouldNot(HaveOccurred())
			for i := 0; i < 100; i++ {
				err = db.InsertOrderFragment(orderFragments[i])
				Expect(err).ShouldNot(HaveOccurred())
				err = db.InsertOrder(orders[i])
				Expect(err).ShouldNot(HaveOccurred())
			}
			for i := 0; i < 100; i++ {
				err = db.RemoveOrderFragment(orders[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
				err = db.RemoveOrder(orders[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
			}
			for i := 0; i < 100; i++ {
				_, err = db.OrderFragment(orders[i].ID)
				Expect(err).Should(Equal(orderbook.ErrOrderFragmentNotFound))
				_, err = db.Order(orders[i].ID)
				Expect(err).Should(Equal(orderbook.ErrOrderNotFound))
			}
			err = db.Close()
			Expect(err).ShouldNot(HaveOccurred())
		})

	})

	Context("when rebooting", func() {

		It("should load data that were stored before rebooting", func() {
			db, err := NewStore("./tmp")
			Expect(err).ShouldNot(HaveOccurred())
			for i := 0; i < 100; i++ {
				err = db.InsertOrderFragment(orderFragments[i])
				Expect(err).ShouldNot(HaveOccurred())
				err = db.InsertOrder(orders[i])
				Expect(err).ShouldNot(HaveOccurred())
			}
			err = db.Close()
			Expect(err).ShouldNot(HaveOccurred())

			for i := 0; i < 10; i++ {
				nextDb, err := NewStore("./tmp")
				Expect(err).ShouldNot(HaveOccurred())
				for j := 0; j < 100; j++ {
					orderFragment, err := nextDb.OrderFragment(orders[i].ID)
					Expect(err).ShouldNot(HaveOccurred())
					order, err := nextDb.Order(orders[i].ID)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(orderFragment.Equal(&orderFragments[i])).Should(BeTrue())
					Expect(order.Equal(&orders[i])).Should(BeTrue())
				}
				err = nextDb.Close()
				Expect(err).ShouldNot(HaveOccurred())
			}
		})

	})

})
