package ome_test

import (
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/republic-go/testutils"
)

const NumberOfNodes = 24

var _ = Describe("Settler", func() {
	var (
		storers   [NumberOfNodes]ComputationStorer
		smpcers   [NumberOfNodes]*testutils.Smpc
		contracts [NumberOfNodes]*omeBinder
		settles   [NumberOfNodes]Settler
	)

	BeforeEach(func() {
		for i := 0; i < NumberOfNodes; i++ {
			storer, err := leveldb.NewStore(fmt.Sprintf("./data-%v.out", i), 72*time.Hour)
			Ω(err).ShouldNot(HaveOccurred())
			storers[i] = storer.SomerComputationStore()
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
			buyFragments, err := testutils.RandomBuyOrderFragments(int64(NumberOfNodes), int64(2*(NumberOfNodes+1)/3))
			Expect(err).ShouldNot(HaveOccurred())
			sellFragments, err := testutils.RandomSellOrderFragments(int64(NumberOfNodes), int64(2*(NumberOfNodes+1)/3))
			Expect(err).ShouldNot(HaveOccurred())
			comp := NewComputation([32]byte{}, buyFragments[0], sellFragments[0], ComputationStateNil, true)
			for i := 0; i < NumberOfNodes; i++ {
				Ω(storers[i].PutComputation(comp)).ShouldNot(HaveOccurred())
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
