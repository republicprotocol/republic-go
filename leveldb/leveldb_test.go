package leveldb_test

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/leveldb"

	"github.com/republicprotocol/republic-go/ome"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
)

var _ = Describe("LevelDB storage", func() {

	orders := make([]order.Order, 100)
	orderFragments := make([]order.Fragment, 100)
	computations := make([]ome.Computation, 100)

	BeforeEach(func() {
		for i := 0; i < 100; i++ {
			ord := order.NewOrder(order.TypeMidpoint, order.ParityBuy, time.Now(), order.TokensETHREN, order.NewCoExp(200, 26), order.NewCoExp(200, 26), order.NewCoExp(200, 26), int64(i))
			ordFragments, err := ord.Split(3, 2)
			Expect(err).ShouldNot(HaveOccurred())
			orders[i] = ord
			orderFragments[i] = ordFragments[0]
			computations[i] = ome.NewComputation(ord.ID, ord.ID, [32]byte{})
		}
	})

	AfterEach(func() {
		os.RemoveAll("./tmp/")
	})

	Context("when storing, loading, and removing data", func() {

		It("should return a default value when loading pointers before storing them", func() {
			db, err := NewStore("./tmp")
			Expect(err).ShouldNot(HaveOccurred())
			buyPointer, err := db.BuyPointer()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buyPointer).Should(Equal(0))
			sellPointer, err := db.SellPointer()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sellPointer).Should(Equal(0))
			err = db.Close()
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return an error when loading data before storing it", func() {
			db, err := NewStore("./tmp")
			Expect(err).ShouldNot(HaveOccurred())
			for i := 0; i < 100; i++ {
				_, err = db.OrderFragment(orders[i].ID)
				Expect(err).Should(Equal(orderbook.ErrOrderFragmentNotFound))
				_, err = db.Order(orders[i].ID)
				Expect(err).Should(Equal(orderbook.ErrOrderNotFound))
				_, err = db.Computation(computations[i].ID)
				Expect(err).Should(Equal(ome.ErrComputationNotFound))
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
				err = db.RemoveComputation(computations[i].ID)
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
				err = db.InsertComputation(computations[i])
				Expect(err).ShouldNot(HaveOccurred())
			}
			for i := 0; i < 100; i++ {
				orderFragment, err := db.OrderFragment(orders[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
				order, err := db.Order(orders[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
				com, err := db.Computation(computations[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(orderFragment.Equal(&orderFragments[i])).Should(BeTrue())
				Expect(order.Equal(&orders[i])).Should(BeTrue())
				Expect(com.Equal(&computations[i])).Should(BeTrue())
			}
			err = db.InsertBuyPointer(42)
			Expect(err).ShouldNot(HaveOccurred())
			buyPointer, err := db.BuyPointer()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buyPointer).Should(Equal(42))
			err = db.InsertSellPointer(420)
			Expect(err).ShouldNot(HaveOccurred())
			sellPointer, err := db.SellPointer()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sellPointer).Should(Equal(420))
			err = db.Close()
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should load all data that was stored", func() {
			db, err := NewStore("./tmp")
			Expect(err).ShouldNot(HaveOccurred())
			for i := 0; i < 100; i++ {
				err = db.InsertOrder(orders[i])
				Expect(err).ShouldNot(HaveOccurred())
				err = db.InsertComputation(computations[i])
				Expect(err).ShouldNot(HaveOccurred())
			}

			ords, err := db.Orders()
			Expect(err).ShouldNot(HaveOccurred())
			coms, err := db.Computations()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ords).Should(HaveLen(len(orders)))
			Expect(coms).Should(HaveLen(len(computations)))
			for i := 0; i < 100; i++ {
				foundOrd := false
				foundCom := false
				for j := 0; j < 100; j++ {
					if foundOrd || ords[i].Equal(&orders[j]) {
						foundOrd = true
					}
					if foundCom || coms[i].Equal(&computations[j]) {
						foundCom = true
					}
					if foundOrd && foundCom {
						break
					}
				}
				Expect(foundOrd).Should(BeTrue())
				Expect(foundCom).Should(BeTrue())
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
				err = db.InsertComputation(computations[i])
				Expect(err).ShouldNot(HaveOccurred())
			}
			for i := 0; i < 100; i++ {
				err = db.RemoveOrderFragment(orders[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
				err = db.RemoveOrder(orders[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
				err = db.RemoveComputation(computations[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
			}
			for i := 0; i < 100; i++ {
				_, err = db.OrderFragment(orders[i].ID)
				Expect(err).Should(Equal(orderbook.ErrOrderFragmentNotFound))
				_, err = db.Order(orders[i].ID)
				Expect(err).Should(Equal(orderbook.ErrOrderNotFound))
				_, err = db.Computation(computations[i].ID)
				Expect(err).Should(Equal(ome.ErrComputationNotFound))
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
				err = db.InsertComputation(computations[i])
				Expect(err).ShouldNot(HaveOccurred())
			}
			err = db.InsertBuyPointer(42)
			Expect(err).ShouldNot(HaveOccurred())
			err = db.InsertSellPointer(420)
			Expect(err).ShouldNot(HaveOccurred())
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
					com, err := nextDb.Computation(computations[i].ID)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(orderFragment.Equal(&orderFragments[i])).Should(BeTrue())
					Expect(order.Equal(&orders[i])).Should(BeTrue())
					Expect(com.Equal(&computations[i])).Should(BeTrue())
				}
				buyPointer, err := nextDb.BuyPointer()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(buyPointer).Should(Equal(42))
				sellPointer, err := nextDb.SellPointer()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(sellPointer).Should(Equal(420))
				err = nextDb.Close()
				Expect(err).ShouldNot(HaveOccurred())
			}
		})

	})

})
