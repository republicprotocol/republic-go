package smpc_test

import (
	"math/rand"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/order"
	. "github.com/republicprotocol/republic-go/smpc"

	"github.com/republicprotocol/republic-go/shamir"
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
})

func generateJoins(id JoinID) (order.Order, []Join) {
	ord := newRandomOrder()
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

func newRandomOrder() order.Order {
	parity := []order.Parity{order.ParityBuy, order.ParitySell}[rand.Intn(2)]
	tokens := []order.Tokens{order.TokensBTCETH,
		order.TokensBTCDGX,
		order.TokensBTCREN,
		order.TokensETHDGX,
		order.TokensETHREN,
		order.TokensDGXREN,
	}[rand.Intn(6)]

	ord := order.NewOrder(order.TypeLimit, parity, time.Now().Add(1*time.Hour), tokens, randomCoExp(), randomCoExp(), randomCoExp(), rand.Int63())
	return ord
}

func randomCoExp() order.CoExp {
	co := uint64(rand.Intn(1999) + 1)
	exp := uint64(rand.Intn(25))
	return order.CoExp{
		Co:  co,
		Exp: exp,
	}
}
