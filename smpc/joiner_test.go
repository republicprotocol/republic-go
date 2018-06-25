package smpc_test

import (
	"bytes"
	"sync/atomic"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/crypto"
	. "github.com/republicprotocol/republic-go/smpc"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Joiner", func() {

	var n = int64(24)
	var k = 2 * (n + 1) / 3
	var joiner *Joiner

	BeforeEach(func() {
		joiner = NewJoiner(n)
	})

	Context("when joining", func() {

		Context("when setting the callback at the first insertion", func() {
			It("should call the callback after inserting n joins", func() {
				ord, joins := generateJoins(n, k)
				called := int64(0)
				callback := generateCallback(&called, ord)

				for i := int64(0); i < n; i++ {
					if i == 0 {
						Ω(joiner.InsertJoinAndSetCallback(joins[i], callback)).ShouldNot(HaveOccurred())
					} else {
						Ω(joiner.InsertJoin(joins[i])).ShouldNot(HaveOccurred())
					}
					if i == n-1 {
						Ω(atomic.LoadInt64(&called)).Should(Equal(int64(1)))
					} else {
						Ω(atomic.LoadInt64(&called)).Should(Equal(int64(0)))
					}
				}
			})
		})

		Context("when setting the callback at the last insertion", func() {
			It("should call the callback after inserting joins", func() {
				ord, joins := generateJoins(n, k)
				called := int64(0)
				callback := generateCallback(&called, ord)

				for i := int64(0); i < n; i++ {
					if i == n-1 {
						Ω(joiner.InsertJoinAndSetCallback(joins[i], callback)).ShouldNot(HaveOccurred())
						Ω(atomic.LoadInt64(&called)).Should(Equal(int64(1)))
					} else {
						Ω(joiner.InsertJoin(joins[i])).ShouldNot(HaveOccurred())
						Ω(atomic.LoadInt64(&called)).Should(Equal(int64(0)))
					}
				}
			})
		})

		Context("when the callback is set multiple times", func() {
			It("should call the latest callback after inserting n joins", func() {
				ord, joins := generateJoins(n, k)
				var called = int64(0)
				callback := generateCallback(&called, ord)
				callbackOverride := func(id JoinID, values []uint64) {
					atomic.AddInt64(&called, 2)
				}

				for i := int64(0); i < n; i++ {
					if i == 0 {
						Ω(joiner.InsertJoinAndSetCallback(joins[i], callback)).ShouldNot(HaveOccurred())
					} else if i == 2 {
						Ω(joiner.InsertJoinAndSetCallback(joins[i], callbackOverride)).ShouldNot(HaveOccurred())
					} else {
						Ω(joiner.InsertJoin(joins[i])).ShouldNot(HaveOccurred())
					}
					if i == n-1 {
						Ω(atomic.LoadInt64(&called)).Should(Equal(int64(2)))
					} else {
						Ω(atomic.LoadInt64(&called)).Should(Equal(int64(0)))
					}
				}
			})
		})

		Context("when inserting computed joins", func() {
			It("should pass the computed values to the callback", func() {
				joins := generateMatchedJoins(n, k)
				called := int64(0)
				callback := func(id JoinID, values []uint64) {
					atomic.AddInt64(&called, 1)
					Ω(len(values)).Should(Equal(7))
					Ω(values[0]).Should(BeNumerically("<=", shamir.Prime/2))
					Ω(values[1]).Should(BeNumerically("<=", shamir.Prime/2))
					Ω(values[2]).Should(BeNumerically("<=", shamir.Prime/2))
					Ω(values[3]).Should(BeNumerically("<=", shamir.Prime/2))
					Ω(values[4]).Should(BeNumerically("<=", shamir.Prime/2))
					Ω(values[5]).Should(BeNumerically("<=", shamir.Prime/2))
					Ω(values[6]).Should(BeZero())
				}

				for i := int64(0); i < n; i++ {
					if i == n-1 {
						Ω(joiner.InsertJoinAndSetCallback(joins[i], callback)).ShouldNot(HaveOccurred())
						Ω(atomic.LoadInt64(&called)).Should(Equal(int64(1)))
					} else {
						Ω(joiner.InsertJoin(joins[i])).ShouldNot(HaveOccurred())
						Ω(atomic.LoadInt64(&called)).Should(Equal(int64(0)))
					}
				}
			})
		})
	})

	Context("when marshaling and unmarshaling joins", func() {
		It("should get the same join after marshal and unmarshal", func() {
			_, joins := generateJoins(n, k)
			for i := range joins {
				data, err := joins[i].MarshalBinary()
				Ω(err).ShouldNot(HaveOccurred())

				newJoin := new(Join)
				err = newJoin.UnmarshalBinary(data)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(bytes.Compare(joins[i].ID[:], newJoin.ID[:])).Should(Equal(0))
				Ω(joins[i].Index).Should(Equal(newJoin.Index))
				Ω(len(joins[i].Shares)).Should(Equal(len(newJoin.Shares)))
				for j := range joins[i].Shares {
					Ω(joins[i].Shares[j].Equal(&newJoin.Shares[j]))
				}
			}
		})
	})

	Context("when inserting joins with shares that exceed the maximum", func() {
		It("should return an error", func() {
			_, joins := generateJoins(n, k)
			for i := range joins {
				shares := make([]shamir.Share, MaxJoinLength+1)
				for j := range shares {
					shares[j] = joins[i].Shares[i%(len(joins[i].Shares))]
				}
				joins[i].Shares = shares
			}
			for i := range joins {
				Ω(joiner.InsertJoin(joins[i])).Should(Equal(ErrJoinLengthExceedsMax))
			}
		})
	})

	Context("when inserting joins with different numbers of shares", func() {
		It("should return an error", func() {
			_, joins := generateJoins(n, k)

			for i := int64(0); i < n; i++ {
				if i > n/2 {
					joins[i].Shares = joins[i].Shares[:3]
					Ω(joiner.InsertJoin(joins[i])).Should(Equal(ErrJoinLengthUnequal))
				} else {
					Ω(joiner.InsertJoin(joins[i])).ShouldNot(HaveOccurred())
				}
			}
		})

		It("should not call the callback more than once", func() {
			ord, joins := generateJoins(n, k)
			called := int64(0)
			callback := generateCallback(&called, ord)

			for i := int64(0); i < n; i++ {
				if i == 0 {
					Ω(joiner.InsertJoinAndSetCallback(joins[i], callback)).ShouldNot(HaveOccurred())
				} else {
					Ω(joiner.InsertJoin(joins[i])).ShouldNot(HaveOccurred())
				}
			}
			Ω(atomic.LoadInt64(&called)).Should(Equal(int64(1)))
		})
	})
})

func generateJoins(n, k int64) (order.Order, []Join) {
	ord := testutils.RandomOrder()
	fragments, err := ord.Split(n, k)
	Ω(err).ShouldNot(HaveOccurred())
	joins := make([]Join, n)
	for i := range joins {
		shares := []shamir.Share{
			fragments[i].Price.Co,
			fragments[i].Price.Exp,
			fragments[i].Volume.Co,
			fragments[i].Volume.Exp,
			fragments[i].MinimumVolume.Co,
			fragments[i].MinimumVolume.Exp,
			fragments[i].Tokens,
		}
		joins[i] = Join{
			Index:  JoinIndex(i + 1),
			Shares: shares,
		}
		copy(joins[i].ID[:], ord.ID[:])
	}

	return ord, joins
}

func generateMatchedJoins(n, k int64) []Join {
	buy, sell := testutils.RandomOrderMatch()
	buyFragments, err := buy.Split(n, k)
	Ω(err).ShouldNot(HaveOccurred())
	sellFragments, err := sell.Split(n, k)
	Ω(err).ShouldNot(HaveOccurred())

	joins := make([]Join, n)
	for i := range joins {
		shares := []shamir.Share{
			buyFragments[i].Price.Co.Sub(&sellFragments[i].Price.Co),
			buyFragments[i].Price.Exp.Sub(&sellFragments[i].Price.Exp),
			buyFragments[i].Volume.Co.Sub(&sellFragments[i].MinimumVolume.Co),
			buyFragments[i].Volume.Exp.Sub(&sellFragments[i].MinimumVolume.Exp),
			sellFragments[i].Volume.Co.Sub(&buyFragments[i].MinimumVolume.Co),
			sellFragments[i].Volume.Exp.Sub(&buyFragments[i].MinimumVolume.Exp),
			buyFragments[i].Tokens.Sub(&sellFragments[i].Tokens),
		}
		joins[i] = Join{
			Index:  JoinIndex(i),
			Shares: shares,
		}
		copy(joins[i].ID[:], crypto.Keccak256(buy.ID[:], sell.ID[:]))
	}

	return joins
}

func generateCallback(called *int64, ord order.Order) func(id JoinID, values []uint64) {
	return func(id JoinID, values []uint64) {
		atomic.AddInt64(called, 1)
		Ω(len(values)).Should(Equal(7))
		Ω(values[0]).Should(Equal(ord.Price.Co))
		Ω(values[1]).Should(Equal(ord.Price.Exp))
		Ω(values[2]).Should(Equal(ord.Volume.Co))
		Ω(values[3]).Should(Equal(ord.Volume.Exp))
		Ω(values[4]).Should(Equal(ord.MinimumVolume.Co))
		Ω(values[5]).Should(Equal(ord.MinimumVolume.Exp))
		Ω(values[6]).Should(Equal(uint64(ord.Tokens)))
	}
}
