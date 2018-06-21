package ome_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/leveldb"
	. "github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/testutils"
)

const NumberOfNodes = 24

var _ = Describe("Settler", func() {
	var (
		storers   [NumberOfNodes]*leveldb.Store
		smpcers   [NumberOfNodes]*testutils.Smpc
		contracts [NumberOfNodes]*omeBinder
		settles   [NumberOfNodes]Settler
	)

	BeforeEach(func() {
		for i := 0; i < NumberOfNodes; i++ {
			storer, err := leveldb.NewStore(fmt.Sprintf("./data-%v.out", i))
			Ω(err).ShouldNot(HaveOccurred())
			storers[i] = storer
			smpcers[i] = testutils.NewAlwaysMatchSmpc()
			contracts[i] = newOmeBinder()
			settles[i] = NewSettler(storers[i], smpcers[i], contracts[i])
		}
	})

	AfterEach(func() {
		for i := 0; i < NumberOfNodes; i++ {
			os.RemoveAll(fmt.Sprintf("./data-%v.out", i))
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
				Ω(storers[i].PutOrderFragment(buyShares[i])).ShouldNot(HaveOccurred())
				Ω(storers[i].PutOrderFragment(sellShares[i])).ShouldNot(HaveOccurred())
			}

			for i := 0; i < NumberOfNodes; i++ {
				Ω(settles[i].Settle(comp)).ShouldNot(HaveOccurred())
			}

			for i := 0; i < NumberOfNodes; i++ {
				Ω(contracts[i].SettleCounts()).Should(Equal(1))
			}
		})
	})
})
