package ome_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/testutils"
)

const NumberOfNodes = 24

var _ = Describe("Settler", func() {
	var (
		storers  [NumberOfNodes]Storer
		smpcers  [NumberOfNodes]*testutils.Smpc
		accounts [NumberOfNodes]*testutils.DarkpoolAccounts
		settles  [NumberOfNodes]Settler
	)

	BeforeEach(func() {
		for i := 0; i < NumberOfNodes; i++ {
			storers[i] = testutils.NewStorer()
			smpcers[i] = testutils.NewAlwaysMatchSmpc()
			accounts[i] = testutils.NewDarkpoolAccounts()
			settles[i] = NewSettler(storers[i], smpcers[i], accounts[i])
		}
	})

	Context("when a computation has been resolved to a match and been confirmed ", func() {
		It("should be able to reconstruct the order and settle it.", func() {
			buy, sell := testutils.RandomBuyOrder(), testutils.RandomSellOrder()
			comp := NewComputation(buy.ID, sell.ID, [32]byte{})
			buyShares, err := buy.Split(int64(NumberOfNodes), int64(2*(NumberOfNodes+1)/3))
			Ω(err).ShouldNot(HaveOccurred())
			sellShares, err := sell.Split(int64(NumberOfNodes), int64(2*(NumberOfNodes+1)/3))
			Ω(err).ShouldNot(HaveOccurred())
			for i := 0; i < NumberOfNodes; i++ {
				Ω(storers[i].InsertOrderFragment(buyShares[i])).ShouldNot(HaveOccurred())
				Ω(storers[i].InsertOrderFragment(sellShares[i])).ShouldNot(HaveOccurred())
			}

			for i := 0; i < NumberOfNodes; i++ {
				Ω(settles[i].Settle(comp)).ShouldNot(HaveOccurred())
			}

			for i := 0; i < NumberOfNodes; i++ {
				Ω(accounts[i].SettleCounts()).Should(Equal(1))
			}
		})
	})
})
