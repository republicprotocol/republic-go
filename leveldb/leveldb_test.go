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

	changes := make([]orderbook.Change, 100)
	orders := make([]order.Order, 100)
	orderFragments := make([]order.Fragment, 100)
	computations := make([]ome.Computation, 100)

	BeforeEach(func() {
		for i := 0; i < 100; i++ {
			ord := order.NewOrder(order.TypeMidpoint, order.ParityBuy, order.SettlementRenEx, time.Now(), order.TokensETHREN, order.NewCoExp(200, 26), order.NewCoExp(200, 26), order.NewCoExp(200, 26), uint64(i))
			change := orderbook.NewChange(ord.ID, ord.Parity, order.Open, orderbook.Priority(i), "", uint(i))
			ordFragments, err := ord.Split(3, 2)
			Expect(err).ShouldNot(HaveOccurred())
			changes[i] = change
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
			Expect(buyPointer).Should(Equal(orderbook.Pointer(0)))
			sellPointer, err := db.SellPointer()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sellPointer).Should(Equal(orderbook.Pointer(0)))
			err = db.Close()
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return an error when loading data before storing it", func() {
			db, err := NewStore("./tmp")
			Expect(err).ShouldNot(HaveOccurred())
			for i := 0; i < 100; i++ {
				_, err = db.Change(changes[i].OrderID)
				Expect(err).Should(Equal(orderbook.ErrChangeNotFound))
				_, err = db.OrderFragment(orderFragments[i].OrderID)
				Expect(err).Should(Equal(orderbook.ErrOrderFragmentNotFound))
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
				err = db.DeleteChange(changes[i].OrderID)
				Expect(err).ShouldNot(HaveOccurred())
				err = db.DeleteOrderFragment(orderFragments[i].OrderID)
				Expect(err).ShouldNot(HaveOccurred())
				err = db.DeleteComputation(computations[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
			}
			err = db.Close()
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should load data equal to when it was stored", func() {
			db, err := NewStore("./tmp")
			Expect(err).ShouldNot(HaveOccurred())
			for i := 0; i < 100; i++ {
				err = db.PutChange(changes[i])
				Expect(err).ShouldNot(HaveOccurred())
				err = db.PutOrderFragment(orderFragments[i])
				Expect(err).ShouldNot(HaveOccurred())
				err = db.PutComputation(computations[i])
				Expect(err).ShouldNot(HaveOccurred())
			}
			for i := 0; i < 100; i++ {
				change, err := db.Change(changes[i].OrderID)
				Expect(err).ShouldNot(HaveOccurred())
				orderFragment, err := db.OrderFragment(orders[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
				com, err := db.Computation(computations[i].ID)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(change.Equal(&changes[i])).Should(BeTrue())
				Expect(orderFragment.Equal(&orderFragments[i])).Should(BeTrue())
				Expect(com.Equal(&computations[i])).Should(BeTrue())
			}
			err = db.PutBuyPointer(42)
			Expect(err).ShouldNot(HaveOccurred())
			buyPointer, err := db.BuyPointer()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buyPointer).Should(Equal(orderbook.Pointer(42)))

			err = db.PutSellPointer(420)
			Expect(err).ShouldNot(HaveOccurred())
			sellPointer, err := db.SellPointer()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sellPointer).Should(Equal(orderbook.Pointer(420)))

			err = db.Close()
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should load all data that was stored", func() {
			db, err := NewStore("./tmp")
			Expect(err).ShouldNot(HaveOccurred())
			for i := 0; i < 100; i++ {
				err = db.PutChange(changes[i])
				Expect(err).ShouldNot(HaveOccurred())
				err = db.PutOrderFragment(orderFragments[i])
				Expect(err).ShouldNot(HaveOccurred())
				err = db.PutComputation(computations[i])
				Expect(err).ShouldNot(HaveOccurred())
			}

			changesIter, err := db.Changes()
			Expect(err).ShouldNot(HaveOccurred())
			defer changesIter.Release()
			chngs, err := changesIter.Collect()
			Expect(err).ShouldNot(HaveOccurred())

			fragmentsIter, err := db.OrderFragments()
			Expect(err).ShouldNot(HaveOccurred())
			defer fragmentsIter.Release()
			frgmnts, err := fragmentsIter.Collect()
			Expect(err).ShouldNot(HaveOccurred())

			comsIter, err := db.Computations()
			Expect(err).ShouldNot(HaveOccurred())
			defer comsIter.Release()
			coms, err := comsIter.Collect()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(chngs).Should(HaveLen(len(changes)))
			Expect(frgmnts).Should(HaveLen(len(orderFragments)))
			Expect(coms).Should(HaveLen(len(computations)))

			changesIter, err = db.Changes()
			Expect(err).ShouldNot(HaveOccurred())
			changesIter.Next()
			changesCur, err := changesIter.Cursor()
			Expect(err).ShouldNot(HaveOccurred())

			fragmentsIter, err = db.OrderFragments()
			Expect(err).ShouldNot(HaveOccurred())
			fragmentsIter.Next()
			fragmentsCur, err := fragmentsIter.Cursor()
			Expect(err).ShouldNot(HaveOccurred())

			comsIter, err = db.Computations()
			Expect(err).ShouldNot(HaveOccurred())
			comsIter.Next()
			comsCur, err := comsIter.Cursor()
			Expect(err).ShouldNot(HaveOccurred())

			for i := 0; i < 100; i++ {
				foundChange := false
				foundFragment := false
				foundCom := false
				for j := 0; j < 100; j++ {
					if foundChange || changesCur.Equal(&changes[j]) {
						foundChange = true
					}
					if foundFragment || fragmentsCur.Equal(&orderFragments[j]) {
						foundFragment = true
					}
					if foundCom || comsCur.Equal(&computations[j]) {
						foundCom = true
					}
					if foundChange && foundCom {
						break
					}
					changesIter.Next()
					fragmentsIter.Next()
					comsIter.Next()
				}
				Expect(foundChange).Should(BeTrue())
				Expect(foundFragment).Should(BeTrue())
				Expect(foundCom).Should(BeTrue())
			}
			err = db.Close()
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return an error when loading data after removing it", func() {
			db, err := NewStore("./tmp")
			Expect(err).ShouldNot(HaveOccurred())
			for i := 0; i < 100; i++ {
				err = db.PutChange(changes[i])
				Expect(err).ShouldNot(HaveOccurred())
				err = db.PutOrderFragment(orderFragments[i])
				Expect(err).ShouldNot(HaveOccurred())
				err = db.PutComputation(computations[i])
				Expect(err).ShouldNot(HaveOccurred())
			}
			for i := 0; i < 100; i++ {
				err = db.DeleteChange(changes[i].OrderID)
				Expect(err).ShouldNot(HaveOccurred())
				err = db.DeleteOrderFragment(orders[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
				err = db.DeleteComputation(computations[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
			}
			for i := 0; i < 100; i++ {
				_, err = db.Change(changes[i].OrderID)
				Expect(err).Should(Equal(orderbook.ErrChangeNotFound))
				_, err = db.OrderFragment(orders[i].ID)
				Expect(err).Should(Equal(orderbook.ErrOrderFragmentNotFound))
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
				err = db.PutChange(changes[i])
				Expect(err).ShouldNot(HaveOccurred())
				err = db.PutOrderFragment(orderFragments[i])
				Expect(err).ShouldNot(HaveOccurred())
				err = db.PutComputation(computations[i])
				Expect(err).ShouldNot(HaveOccurred())
			}
			err = db.PutBuyPointer(42)
			Expect(err).ShouldNot(HaveOccurred())
			err = db.PutSellPointer(420)
			Expect(err).ShouldNot(HaveOccurred())
			err = db.Close()
			Expect(err).ShouldNot(HaveOccurred())

			for i := 0; i < 10; i++ {
				nextDb, err := NewStore("./tmp")
				Expect(err).ShouldNot(HaveOccurred())
				for j := 0; j < 100; j++ {
					change, err := nextDb.Change(changes[i].OrderID)
					Expect(err).ShouldNot(HaveOccurred())
					orderFragment, err := nextDb.OrderFragment(orders[i].ID)
					Expect(err).ShouldNot(HaveOccurred())
					com, err := nextDb.Computation(computations[i].ID)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(change.Equal(&changes[i])).Should(BeTrue())
					Expect(orderFragment.Equal(&orderFragments[i])).Should(BeTrue())
					Expect(com.Equal(&computations[i])).Should(BeTrue())
				}
				buyPointer, err := nextDb.BuyPointer()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(buyPointer).Should(Equal(orderbook.Pointer(42)))
				sellPointer, err := nextDb.SellPointer()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(sellPointer).Should(Equal(orderbook.Pointer(420)))
				err = nextDb.Close()
				Expect(err).ShouldNot(HaveOccurred())
			}
		})

	})

})
