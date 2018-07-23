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
	"github.com/republicprotocol/republic-go/registry"
)

var expiry = 72 * time.Hour
var _ = Describe("LevelDB storage", func() {

	orders := make([]order.Order, 100)
	orderFragments := make([]order.Fragment, 100)
	computations := make([]ome.Computation, 50)

	BeforeEach(func() {
		j := 0
		for i := 0; i < 100; i++ {
			ord := order.NewOrder(order.TypeMidpoint, order.ParityBuy, order.SettlementRenEx, time.Now(), order.TokensETHREN, order.NewCoExp(200, 26), order.NewCoExp(200, 26), order.NewCoExp(200, 26), uint64(i))
			ordFragments, err := ord.Split(3, 2)
			Expect(err).ShouldNot(HaveOccurred())
			orders[i] = ord
			orderFragments[i] = ordFragments[0]
			if (i+1)%2 == 0 {
				computations[j] = ome.NewComputation([32]byte{byte(j)}, orderFragments[i], orderFragments[i-1], ome.ComputationStateMatched, true)
				j++
			}
		}
	})

	AfterEach(func() {
		os.RemoveAll("./tmp/")
	})

	Context("when pruning data", func() {

		It("should not load any expired data", func() {
			db, err := NewStore("./tmp", 2*time.Second)
			Expect(err).ShouldNot(HaveOccurred())

			for i := 0; i < 100; i++ {
				err = db.OrderbookOrderStore().PutOrder(orders[i].ID, order.Open, "")
				Expect(err).ShouldNot(HaveOccurred())
				err = db.OrderbookOrderFragmentStore().PutOrderFragment(registry.Epoch{}, orderFragments[i])
				Expect(err).ShouldNot(HaveOccurred())
				if i < 50 {
					err = db.SomerComputationStore().PutComputation(computations[i])
					Expect(err).ShouldNot(HaveOccurred())
				}
			}

			// Prune the data
			time.Sleep(2 * time.Second)
			db.Prune()

			for i := 0; i < 100; i++ {
				_, _, err = db.OrderbookOrderStore().Order(orders[i].ID)
				Expect(err).Should(Equal(orderbook.ErrOrderNotFound))
				_, err = db.OrderbookOrderFragmentStore().OrderFragment(registry.Epoch{}, orders[i].ID)
				Expect(err).Should(Equal(orderbook.ErrOrderFragmentNotFound))
				if i < 50 {
					_, err = db.SomerComputationStore().Computation(computations[i].ID)
					Expect(err).Should(Equal(ome.ErrComputationNotFound))
				}
			}

			err = db.Release()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("when storing, loading, and removing data", func() {

		It("should return a default value when loading pointers before storing them", func() {
			db, err := NewStore("./tmp", expiry)
			Expect(err).ShouldNot(HaveOccurred())
			pointer, err := db.OrderbookPointerStore().Pointer()
			Expect(pointer).Should(Equal(orderbook.Pointer(0)))
			err = db.Release()
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return an error when loading data before storing it", func() {
			db, err := NewStore("./tmp", expiry)
			Expect(err).ShouldNot(HaveOccurred())

			for i := 0; i < 100; i++ {
				_, _, err = db.OrderbookOrderStore().Order(orders[i].ID)
				Expect(err).Should(Equal(orderbook.ErrOrderNotFound))
				_, err = db.OrderbookOrderFragmentStore().OrderFragment(registry.Epoch{}, orderFragments[i].OrderID)
				Expect(err).Should(Equal(orderbook.ErrOrderFragmentNotFound))
				if i < 50 {
					_, err = db.SomerComputationStore().Computation(computations[i].ID)
					Expect(err).Should(Equal(ome.ErrComputationNotFound))
				}
			}

			err = db.Release()
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should not return an error when removing data before storing it", func() {
			db, err := NewStore("./tmp", expiry)
			Expect(err).ShouldNot(HaveOccurred())

			for i := 0; i < 100; i++ {
				err = db.OrderbookOrderStore().DeleteOrder(orders[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
				err = db.OrderbookOrderFragmentStore().DeleteOrderFragment(registry.Epoch{}, orderFragments[i].OrderID)
				Expect(err).ShouldNot(HaveOccurred())
				if i < 50 {
					err = db.SomerComputationStore().DeleteComputation(computations[i].ID)
					Expect(err).ShouldNot(HaveOccurred())
				}
			}

			err = db.Release()
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should load data equal to when it was stored", func() {
			db, err := NewStore("./tmp", expiry)
			Expect(err).ShouldNot(HaveOccurred())

			for i := 0; i < 100; i++ {
				err = db.OrderbookOrderStore().PutOrder(orders[i].ID, order.Open, "")
				Expect(err).ShouldNot(HaveOccurred())
				err = db.OrderbookOrderFragmentStore().PutOrderFragment(registry.Epoch{}, orderFragments[i])
				Expect(err).ShouldNot(HaveOccurred())
				if i < 50 {
					err = db.SomerComputationStore().PutComputation(computations[i])
					Expect(err).ShouldNot(HaveOccurred())
				}
			}

			for i := 0; i < 100; i++ {
				status, _, err := db.OrderbookOrderStore().Order(orders[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
				orderFragment, err := db.OrderbookOrderFragmentStore().OrderFragment(registry.Epoch{}, orders[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
				if i < 50 {
					com, err := db.SomerComputationStore().Computation(computations[i].ID)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(com.Equal(&computations[i])).Should(BeTrue())
				}

				Expect(status).Should(Equal(order.Open))
				Expect(orderFragment.Equal(&orderFragments[i])).Should(BeTrue())

			}
			err = db.OrderbookPointerStore().PutPointer(42)
			Expect(err).ShouldNot(HaveOccurred())
			pointer, err := db.OrderbookPointerStore().Pointer()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(pointer).Should(Equal(orderbook.Pointer(42)))

			err = db.Release()
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should load all data that was stored", func() {
			db, err := NewStore("./tmp", expiry)
			Expect(err).ShouldNot(HaveOccurred())

			for i := 0; i < 100; i++ {
				err = db.OrderbookOrderStore().PutOrder(orders[i].ID, order.Open, "")
				Expect(err).ShouldNot(HaveOccurred())
				err = db.OrderbookOrderFragmentStore().PutOrderFragment(registry.Epoch{}, orderFragments[i])
				Expect(err).ShouldNot(HaveOccurred())
				if i < 50 {
					err = db.SomerComputationStore().PutComputation(computations[i])
					Expect(err).ShouldNot(HaveOccurred())
				}
			}

			ordersIter, err := db.OrderbookOrderStore().Orders()
			Expect(err).ShouldNot(HaveOccurred())
			defer ordersIter.Release()
			orderIDs, _, _, err := ordersIter.Collect()
			Expect(err).ShouldNot(HaveOccurred())

			fragmentsIter, err := db.OrderbookOrderFragmentStore().OrderFragments(registry.Epoch{})
			Expect(err).ShouldNot(HaveOccurred())
			defer fragmentsIter.Release()
			frgmnts, err := fragmentsIter.Collect()
			Expect(err).ShouldNot(HaveOccurred())

			comsIter, err := db.SomerComputationStore().Computations()
			Expect(err).ShouldNot(HaveOccurred())
			defer comsIter.Release()
			coms, err := comsIter.Collect()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(orderIDs).Should(HaveLen(len(orders)))
			Expect(frgmnts).Should(HaveLen(len(orderFragments)))
			Expect(coms).Should(HaveLen(len(computations)))

			ordersIter, err = db.OrderbookOrderStore().Orders()
			Expect(err).ShouldNot(HaveOccurred())
			ordersIter.Next()
			ordersID, _, _, err := ordersIter.Cursor()
			Expect(err).ShouldNot(HaveOccurred())

			fragmentsIter, err = db.OrderbookOrderFragmentStore().OrderFragments(registry.Epoch{})
			Expect(err).ShouldNot(HaveOccurred())
			fragmentsIter.Next()
			fragmentsCur, err := fragmentsIter.Cursor()
			Expect(err).ShouldNot(HaveOccurred())

			comsIter, err = db.SomerComputationStore().Computations()
			Expect(err).ShouldNot(HaveOccurred())
			comsIter.Next()
			comsCur, err := comsIter.Cursor()
			Expect(err).ShouldNot(HaveOccurred())

			for i := 0; i < 100; i++ {
				foundChange := false
				foundFragment := false
				foundCom := false
				for j := 0; j < 100; j++ {
					if foundChange || ordersID.Equal(orders[j].ID) {
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
					ordersIter.Next()
					fragmentsIter.Next()
					comsIter.Next()
				}
				Expect(foundChange).Should(BeTrue())
				Expect(foundFragment).Should(BeTrue())
				Expect(foundCom).Should(BeTrue())
			}

			err = db.Release()
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return an error when loading data after removing it", func() {
			db, err := NewStore("./tmp", expiry)
			Expect(err).ShouldNot(HaveOccurred())

			for i := 0; i < 100; i++ {
				err = db.OrderbookOrderStore().PutOrder(orders[i].ID, order.Open, "")
				Expect(err).ShouldNot(HaveOccurred())
				err = db.OrderbookOrderFragmentStore().PutOrderFragment(registry.Epoch{}, orderFragments[i])
				Expect(err).ShouldNot(HaveOccurred())
				if i < 50 {
					err = db.SomerComputationStore().PutComputation(computations[i])
					Expect(err).ShouldNot(HaveOccurred())
				}
			}

			for i := 0; i < 100; i++ {
				err = db.OrderbookOrderStore().DeleteOrder(orders[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
				err = db.OrderbookOrderFragmentStore().DeleteOrderFragment(registry.Epoch{}, orders[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
				if i < 50 {
					err = db.SomerComputationStore().DeleteComputation(computations[i].ID)
					Expect(err).ShouldNot(HaveOccurred())
				}
			}

			for i := 0; i < 100; i++ {
				_, _, err = db.OrderbookOrderStore().Order(orders[i].ID)
				Expect(err).Should(Equal(orderbook.ErrOrderNotFound))
				_, err = db.OrderbookOrderFragmentStore().OrderFragment(registry.Epoch{}, orders[i].ID)
				Expect(err).Should(Equal(orderbook.ErrOrderFragmentNotFound))
				if i < 50 {
					_, err = db.SomerComputationStore().Computation(computations[i].ID)
					Expect(err).Should(Equal(ome.ErrComputationNotFound))
				}
			}

			err = db.Release()
			Expect(err).ShouldNot(HaveOccurred())
		})

	})

	Context("when rebooting", func() {

		It("should load data that were stored before rebooting", func() {
			db, err := NewStore("./tmp", expiry)
			Expect(err).ShouldNot(HaveOccurred())
			for i := 0; i < 100; i++ {
				err = db.OrderbookOrderStore().PutOrder(orders[i].ID, order.Open, "")
				Expect(err).ShouldNot(HaveOccurred())
				err = db.OrderbookOrderFragmentStore().PutOrderFragment(registry.Epoch{}, orderFragments[i])
				Expect(err).ShouldNot(HaveOccurred())
				if i < 50 {
					err = db.SomerComputationStore().PutComputation(computations[i])
					Expect(err).ShouldNot(HaveOccurred())
				}
			}
			err = db.OrderbookPointerStore().PutPointer(42)
			Expect(err).ShouldNot(HaveOccurred())
			err = db.Release()
			Expect(err).ShouldNot(HaveOccurred())

			for i := 0; i < 10; i++ {
				nextDb, err := NewStore("./tmp", expiry)
				Expect(err).ShouldNot(HaveOccurred())
				for j := 0; j < 100; j++ {
					status, _, err := nextDb.OrderbookOrderStore().Order(orders[i].ID)
					Expect(err).ShouldNot(HaveOccurred())
					orderFragment, err := nextDb.OrderbookOrderFragmentStore().OrderFragment(registry.Epoch{}, orders[i].ID)
					Expect(err).ShouldNot(HaveOccurred())
					if i < 50 {
						com, err := nextDb.SomerComputationStore().Computation(computations[i].ID)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(com.Equal(&computations[i])).Should(BeTrue())
					}
					Expect(status).Should(Equal(order.Open))
					Expect(orderFragment.Equal(&orderFragments[i])).Should(BeTrue())

				}
				pointer, err := nextDb.OrderbookPointerStore().Pointer()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(pointer).Should(Equal(orderbook.Pointer(42)))
				err = nextDb.Release()
				Expect(err).ShouldNot(HaveOccurred())
			}
		})

	})

})
