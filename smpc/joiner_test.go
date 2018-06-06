package smpc_test

import (
	"bytes"
	"sync/atomic"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/smpc"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/testutils"
)

var k = int64(24)

var _ = Describe("Joiner", func() {
	var joiner *Joiner

	Context("Insert join and set callback", func() {
		var joinId = JoinID{1}

		BeforeEach(func() {
			joiner = NewJoiner(k)
		})

		It("should call the callback when have enough shares ", func() {
			ord, joins := generateJoins(joinId)
			var getsCalled = int64(0)
			callback := generateCallback(&getsCalled, ord)

			for i := int64(0); i < k; i++ {
				if i == 0 {
					Ω(joiner.InsertJoinAndSetCallback(joins[i], callback)).ShouldNot(HaveOccurred())
				} else {
					Ω(joiner.InsertJoin(joins[i])).ShouldNot(HaveOccurred())
				}

				if i == k-1 {
					Ω(atomic.LoadInt64(&getsCalled)).Should(Equal(int64(1)))
				} else {
					Ω(atomic.LoadInt64(&getsCalled)).Should(Equal(int64(0)))

				}
			}
		})

		It("the insertion order should not matter", func() {
			ord, joins := generateJoins(joinId)
			var getsCalled = int64(0)
			callback := generateCallback(&getsCalled, ord)

			for i := int64(0); i < k; i++ {
				if i == k-1 {
					Ω(joiner.InsertJoinAndSetCallback(joins[i], callback)).ShouldNot(HaveOccurred())
					Ω(atomic.LoadInt64(&getsCalled)).Should(Equal(int64(1)))
				} else {
					Ω(joiner.InsertJoin(joins[i])).ShouldNot(HaveOccurred())
					Ω(atomic.LoadInt64(&getsCalled)).Should(Equal(int64(0)))
				}
			}
		})
	})

	Context("marshal and unmarshal join", func() {
		var joinId = JoinID{2}

		It("should get the same join after marshal and unmarshal", func() {
			_, joins := generateJoins(joinId)
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
})

func generateJoins(id JoinID) (order.Order, []Join) {
	ord := testutils.RandomOrder()
	fragments, err := ord.Split(k, k)
	Ω(err).ShouldNot(HaveOccurred())
	joins := make([]Join, k)
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
			ID:     id,
			Index:  JoinIndex(i),
			Shares: shares,
		}
	}

	return ord, joins
}

func generateCallback(getsCalled *int64, ord order.Order) func(id JoinID, values []uint64) {
	return func(id JoinID, values []uint64) {
		atomic.AddInt64(getsCalled, 1)
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
