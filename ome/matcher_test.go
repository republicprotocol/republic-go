package ome_test

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Matcher", func() {

	var store ComputationStorer
	var buyFragment, sellFragment order.Fragment

	BeforeEach(func() {
		storer, err := leveldb.NewStore("./data.out", 72*time.Hour)
		Expect(err).ShouldNot(HaveOccurred())
		store = storer.SomerComputationStore()

		buyFragments, err := testutils.RandomBuyOrderFragments(6, 4)
		Expect(err).ShouldNot(HaveOccurred())
		buyFragment = buyFragments[0]
		sellFragments, err := testutils.RandomSellOrderFragments(6, 4)
		Expect(err).ShouldNot(HaveOccurred())
		sellFragment = sellFragments[0]
	})

	AfterEach(func() {
		os.RemoveAll("./data.out")
	})

	Context("when using an smpc that matches all values", func() {
		It("should trigger the callback with matched results", func() {
			smpcer := testutils.NewAlwaysMatchSmpc()
			matcher := NewMatcher(store, smpcer)

			numTrials := 100
			numMatches := 0
			for i := 0; i < numTrials; i++ {
				com := NewComputation([32]byte{byte(i)}, buyFragment, sellFragment, ComputationStateNil, true)
				matcher.Resolve(com, func(com Computation) {
					if com.Match {
						numMatches++
					}
				})
			}
			Expect(numMatches).Should(Equal(numTrials))
		})
	})

	Context("when using an smpc that mismatches all values", func() {
		It("should never trigger the callback with matched results", func() {
			smpcer := testutils.NewAlwaysMismatchSmpc()
			matcher := NewMatcher(store, smpcer)

			numTrials := 100
			numMatches := 0
			for i := 0; i < numTrials; i++ {
				com := NewComputation([32]byte{byte(i)}, buyFragment, sellFragment, ComputationStateNil, true)
				matcher.Resolve(com, func(com Computation) {
					if com.Match {
						numMatches++
					}
				})
			}
			Expect(numMatches).Should(Equal(0))
		})
	})

	Context("when using an smpc that randomly matches values", func() {
		It("should randomly trigger the callback with matched results", func() {
			smpcer := testutils.NewSmpc()
			matcher := NewMatcher(store, smpcer)

			numTrials := 1024
			numMatches := 0
			for i := 0; i < numTrials; i++ {
				com := NewComputation([32]byte{byte(i)}, buyFragment, sellFragment, ComputationStateNil, true)
				matcher.Resolve(com, func(com Computation) {
					if com.Match {
						numMatches++
					}
				})
			}
			Expect(numMatches).Should(BeNumerically(">", 0))
			Expect(numMatches).Should(BeNumerically("<", numTrials))
		})
	})
})
