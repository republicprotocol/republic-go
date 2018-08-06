package leveldb_test

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/leveldb"

	"github.com/republicprotocol/republic-go/ome"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/registry"
	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Somer storage", func() {

	computations := make([]ome.Computation, 100)
	buyFragments := make([]order.Fragment, 100)
	sellFragments := make([]order.Fragment, 100)
	dbFolder := "./tmp/"
	dbFile := dbFolder + "db"
	epoch := registry.Epoch{}

	BeforeEach(func() {
		for i := 0; i < 100; i++ {
			buyOrd := order.NewOrder(order.TypeMidpoint, order.ParityBuy, order.SettlementRenEx, time.Now(), order.TokensETHREN, order.NewCoExp(200, 26), order.NewCoExp(200, 26), order.NewCoExp(200, 26), uint64(i))
			buyOrdFragments, err := buyOrd.Split(3, 2)
			sellOrd := order.NewOrder(order.TypeMidpoint, order.ParitySell, order.SettlementRenEx, time.Now(), order.TokensETHREN, order.NewCoExp(200, 26), order.NewCoExp(200, 26), order.NewCoExp(200, 26), uint64(i))
			sellOrdFragments, err := sellOrd.Split(3, 2)
			Expect(err).ShouldNot(HaveOccurred())
			computations[i] = ome.NewComputation([32]byte{byte(i)}, buyOrdFragments[0], sellOrdFragments[0], ome.ComputationStateMatched, true)
			buyFragments[i] = buyFragments[0]
			sellFragments[i] = sellFragments[0]
			_, epoch, err = testutils.RandomEpoch(1)
		}
	})

	AfterEach(func() {
		os.RemoveAll(dbFolder)
	})

	Context("when pruning data", func() {
		It("should not retrieve expired data", func() {
			db := newDB(dbFile)
			somerComputationTable := NewSomerComputationTable(db)
			somerOrderFragmentTable := NewSomerOrderFragmentTable(db, time.Second)

			// Put the computations into the table and attempt to retrieve
			for i := 0; i < len(computations); i++ {
				err := somerComputationTable.PutComputation(computations[i])
				Expect(err).ShouldNot(HaveOccurred())
				somerOrderFragmentTable.PutBuyOrderFragment(epoch, buyFragments[i], "trader1", uint64(i))
				Expect(err).ShouldNot(HaveOccurred())
				somerOrderFragmentTable.PutSellOrderFragment(epoch, sellFragments[i], "trader2", uint64(i))
				Expect(err).ShouldNot(HaveOccurred())
			}
			for i := 0; i < len(computations); i++ {
				com, err := somerComputationTable.Computation(computations[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(com.Equal(&computations[i])).Should(BeTrue())
				buyFragment, trader, _, err := somerOrderFragmentTable.BuyOrderFragment(epoch, buyFragments[i].OrderID)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(buyFragment.Equal(&buyFragments[i])).Should(BeTrue())
				Expect(trader).To(Equal("trader1"))
				sellFragment, trader, _, err := somerOrderFragmentTable.SellOrderFragment(epoch, sellFragments[i].OrderID)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(sellFragment.Equal(&sellFragments[i])).Should(BeTrue())
				Expect(trader).To(Equal("trader2"))
			}

			// Sleep and then prune to expire the data
			time.Sleep(2 * time.Second)
			somerComputationTable.Prune()
			somerOrderFragmentTable.Prune()

			// All data should have expired so we should not get any data back
			comsIter, err := somerComputationTable.Computations()
			Expect(err).ShouldNot(HaveOccurred())
			defer comsIter.Release()
			coms, err := comsIter.Collect()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(coms).Should(HaveLen(0))

			buysIter, err := somerOrderFragmentTable.BuyOrderFragments(epoch)
			Expect(err).ShouldNot(HaveOccurred())
			defer buysIter.Release()
			buys, _, _, err := buysIter.Collect()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(buys).Should(HaveLen(0))

			sellsIter, err := somerOrderFragmentTable.SellOrderFragments(epoch)
			Expect(err).ShouldNot(HaveOccurred())
			defer sellsIter.Release()
			sells, _, _, err := sellsIter.Collect()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(sells).Should(HaveLen(0))
		})
	})

	Context("when iterating through out of range data", func() {
		It("should trigger an out of range error", func() {
			db := newDB(dbFile)
			somerComputationTable := NewSomerComputationTable(db)
			somerOrderFragmentTable := NewSomerOrderFragmentTable(db, time.Second)

			// Put the computations into the table and attempt to retrieve
			for i := 0; i < len(computations); i++ {
				err := somerComputationTable.PutComputation(computations[i])
				Expect(err).ShouldNot(HaveOccurred())
				somerOrderFragmentTable.PutBuyOrderFragment(epoch, buyFragments[i], "trader1", uint64(i))
				Expect(err).ShouldNot(HaveOccurred())
				somerOrderFragmentTable.PutSellOrderFragment(epoch, sellFragments[i], "trader2", uint64(i))
				Expect(err).ShouldNot(HaveOccurred())
			}
			for i := 0; i < len(computations); i++ {
				com, err := somerComputationTable.Computation(computations[i].ID)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(com.Equal(&computations[i])).Should(BeTrue())
				buyFragment, trader, _, err := somerOrderFragmentTable.BuyOrderFragment(epoch, buyFragments[i].OrderID)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(buyFragment.Equal(&buyFragments[i])).Should(BeTrue())
				Expect(trader).To(Equal("trader1"))
				sellFragment, trader, _, err := somerOrderFragmentTable.SellOrderFragment(epoch, sellFragments[i].OrderID)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(sellFragment.Equal(&sellFragments[i])).Should(BeTrue())
				Expect(trader).To(Equal("trader2"))
			}

			comsIter, err := somerComputationTable.Computations()
			defer comsIter.Release()
			for comsIter.Next() {
				_, err := comsIter.Cursor()
				Expect(err).ShouldNot(HaveOccurred())
			}

			// This is out of range so we should expect an error
			_, err = comsIter.Cursor()
			Expect(err).Should(Equal(ome.ErrCursorOutOfRange))

			buysIter, err := somerOrderFragmentTable.BuyOrderFragments(epoch)
			Expect(err).ShouldNot(HaveOccurred())
			defer buysIter.Release()
			for buysIter.Next() {
				_, _, _, err := buysIter.Cursor()
				Expect(err).ShouldNot(HaveOccurred())
			}

			// This is out of range so we should expect an error
			_, _, _, err = buysIter.Cursor()
			Expect(err).Should(Equal(ome.ErrCursorOutOfRange))

			sellsIter, err := somerOrderFragmentTable.SellOrderFragments(epoch)
			Expect(err).ShouldNot(HaveOccurred())
			defer sellsIter.Release()
			for sellsIter.Next() {
				_, _, _, err := sellsIter.Cursor()
				Expect(err).ShouldNot(HaveOccurred())
			}

			// This is out of range so we should expect an error
			_, _, _, err = sellsIter.Cursor()
			Expect(err).Should(Equal(ome.ErrCursorOutOfRange))

		})
	})
})
