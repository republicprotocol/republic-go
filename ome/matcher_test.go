package ome_test

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/leveldb"
	. "github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Matcher", func() {

	Context("when using an smpc that matches all values", func() {
		It("should trigger the callback with matched results", func() {
			store, err := leveldb.NewStore("./data.out")
			Expect(err).ShouldNot(HaveOccurred())
			defer func() {
				os.RemoveAll("./data.out")
			}()
			smpcer := testutils.NewAlwaysMatchSmpc()
			matcher := NewMatcher(store, smpcer)

			numTrials := 100
			numMatches := 0
			for i := 0; i < numTrials; i++ {
				buy, sell := testutils.RandomOrderMatch()
				com := NewComputation(buy.ID, sell.ID, [32]byte{byte(i)})
				com.Timestamp = time.Now()
				com.State = ComputationStateNil
				buyFragments, err := buy.Split(6, 4)
				Expect(err).ShouldNot(HaveOccurred())
				sellFragments, err := sell.Split(6, 4)
				Expect(err).ShouldNot(HaveOccurred())
				matcher.Resolve(com, buyFragments[0], sellFragments[0], func(com Computation) {
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
			store, err := leveldb.NewStore("./data.out")
			Expect(err).ShouldNot(HaveOccurred())
			defer func() {
				os.RemoveAll("./data.out")
			}()
			smpcer := testutils.NewAlwaysMismatchSmpc()
			matcher := NewMatcher(store, smpcer)

			numTrials := 100
			numMatches := 0
			for i := 0; i < numTrials; i++ {
				buy, sell := testutils.RandomOrderMatch()
				com := NewComputation(buy.ID, sell.ID, [32]byte{byte(i)})
				com.Timestamp = time.Now()
				com.State = ComputationStateNil
				buyFragments, err := buy.Split(6, 4)
				Expect(err).ShouldNot(HaveOccurred())
				sellFragments, err := sell.Split(6, 4)
				Expect(err).ShouldNot(HaveOccurred())
				matcher.Resolve(com, buyFragments[0], sellFragments[0], func(com Computation) {
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
			store, err := leveldb.NewStore("./data.out")
			Expect(err).ShouldNot(HaveOccurred())
			defer func() {
				os.RemoveAll("./data.out")
			}()
			smpcer := testutils.NewSmpc()
			matcher := NewMatcher(store, smpcer)

			numTrials := 256
			numMatches := 0
			for i := 0; i < numTrials; i++ {
				buy, sell := testutils.RandomOrderMatch()
				com := NewComputation(buy.ID, sell.ID, [32]byte{byte(i)})
				com.Timestamp = time.Now()
				com.State = ComputationStateNil
				buyFragments, err := buy.Split(6, 4)
				Expect(err).ShouldNot(HaveOccurred())
				sellFragments, err := sell.Split(6, 4)
				Expect(err).ShouldNot(HaveOccurred())
				matcher.Resolve(com, buyFragments[0], sellFragments[0], func(com Computation) {
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
