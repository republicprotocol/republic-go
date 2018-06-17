package ome_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/testutils"
)

// WARNING: The expected number of computations is only
// (NumberOfOrderPairs^2) / numberOfRankers when the numberOfRankers divides
// the NumberOfOrderPairs^2 perfectly.
var (
	EpochBlockNumbers  = []uint{5, 10, 15, 20, 25}
	NumberOfOrderPairs = 10
)

var _ = Describe("OME Ranker", func() {
	var storer Storer
	var ranker Ranker
	var done chan struct{}
	var epoch cal.Epoch
	var addr identity.Address
	var err error

	Context("when processing inserted changes", func() {
		BeforeEach(func() {
			addr, err = testutils.RandomAddress()
			Ω(err).ShouldNot(HaveOccurred())
			epoch = newEpoch(0, addr)
			done = make(chan struct{})
			storer = testutils.NewStorer()
			ranker, err = NewRanker(done, addr, storer, epoch)
			Ω(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			close(done)
		})

		It("should create computations with new opened orders ", func() {
			for i := 0; i < NumberOfOrderPairs; i++ {
				buy, sell := testutils.RandomBuyOrder(), testutils.RandomSellOrder()
				buyChange := newChange(buy, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				sellChange := newChange(sell, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				ranker.InsertChange(buyChange)
				ranker.InsertChange(sellChange)
			}
			computations := make([]Computation, 128)
			i := ranker.Computations(computations)
			Ω(i).Should(Equal(NumberOfOrderPairs * NumberOfOrderPairs))
		})

		It("should not create computations with removed orders", func() {
			for i := 0; i < NumberOfOrderPairs; i++ {
				buy, sell := testutils.RandomBuyOrder(), testutils.RandomSellOrder()
				buyChange := newChange(buy, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				ranker.InsertChange(buyChange)
				cancelBuy := newChange(buy, order.Canceled, orderbook.Priority(i), EpochBlockNumbers[0])
				ranker.InsertChange(cancelBuy)

				sellChange := newChange(sell, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				ranker.InsertChange(sellChange)
				cancelSell := newChange(sell, order.Canceled, orderbook.Priority(i), EpochBlockNumbers[0])
				ranker.InsertChange(cancelSell)
			}
			computations := make([]Computation, 128)
			i := ranker.Computations(computations)
			Ω(i).Should(BeZero())
		})

		It("should return number of computations no more the size of the given variable ", func() {
			for i := 0; i < NumberOfOrderPairs; i++ {
				buy, sell := testutils.RandomBuyOrder(), testutils.RandomSellOrder()
				buyChange := newChange(buy, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				sellChange := newChange(sell, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				ranker.InsertChange(buyChange)
				ranker.InsertChange(sellChange)
			}
			computations := make([]Computation, 10)
			i := ranker.Computations(computations)
			Ω(i).Should(Equal(len(computations)))
		})

		It("should create computations from orders from same trader ", func() {
			for i := 0; i < NumberOfOrderPairs; i++ {
				buy, sell := testutils.RandomBuyOrder(), testutils.RandomSellOrder()
				buyChange := newChange(buy, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				buyChange.Trader = "trader"
				sellChange := newChange(sell, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				sellChange.Trader = "trader"
				ranker.InsertChange(buyChange)
				ranker.InsertChange(sellChange)
			}
			computations := make([]Computation, 10)
			i := ranker.Computations(computations)
			Ω(i).Should(Equal(0))
		})
	})

	Context("when loading stored computations from the storer", func() {
		BeforeEach(func() {
			addr, err = testutils.RandomAddress()
			Ω(err).ShouldNot(HaveOccurred())
			epoch = newEpoch(0, addr)
			done = make(chan struct{})
			storer = testutils.NewStorer()
		})

		AfterEach(func() {
			close(done)
		})

		It("should insert computations from the storer when started", func() {
			for i := 0; i < NumberOfOrderPairs; i++ {
				comp := testutils.RandomComputation()
				Ω(storer.InsertComputation(comp)).ShouldNot(HaveOccurred())
			}

			ranker, err = NewRanker(done, addr, storer, epoch)
			Ω(err).ShouldNot(HaveOccurred())
			time.Sleep(15 * time.Second)

			computations := make([]Computation, 20)
			i := ranker.Computations(computations)
			Ω(i).Should(Equal(NumberOfOrderPairs))
		})
	})

	Context(" when there are multiple rankers", func() {
		BeforeEach(func() {
			addr, err = testutils.RandomAddress()
			Ω(err).ShouldNot(HaveOccurred())
			epoch = newEpoch(0, addr)
			done = make(chan struct{})
			storer = testutils.NewStorer()
		})

		It("should only generate computations meant for its position", func() {
			// add one more pod in the epoch
			another, err := testutils.RandomAddress()
			Ω(err).ShouldNot(HaveOccurred())
			anotherPod := cal.Pod{
				Position:  0,
				Hash:      testutils.Random32Bytes(),
				Darknodes: []identity.Address{another},
			}
			epoch.Pods = append(epoch.Pods, anotherPod)
			epoch.Darknodes = append(epoch.Darknodes, another)
			ranker, err = NewRanker(done, addr, storer, epoch)
			Ω(err).ShouldNot(HaveOccurred())

			for i := 0; i < NumberOfOrderPairs; i++ {
				buy, sell := testutils.RandomBuyOrder(), testutils.RandomSellOrder()
				buyChange := newChange(buy, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				sellChange := newChange(sell, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				ranker.InsertChange(buyChange)
				ranker.InsertChange(sellChange)
			}

			computations := make([]Computation, 128)
			i := ranker.Computations(computations)
			Ω(i).Should(BeNumerically(">=", NumberOfOrderPairs * NumberOfOrderPairs / 2))
		})
	})

	Context("when epoch happens", func() {
		BeforeEach(func() {
			addr, err = testutils.RandomAddress()
			Ω(err).ShouldNot(HaveOccurred())
			epoch = newEpoch(0, addr)
			done = make(chan struct{})
			storer = testutils.NewStorer()
			ranker, err = NewRanker(done, addr, storer, epoch)
			Ω(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			close(done)
		})

		It("should ignore same epoch change event", func() {
			ranker.OnChangeEpoch(epoch)
			for i := 0; i < NumberOfOrderPairs/2; i++ {
				buy, sell := testutils.RandomBuyOrder(), testutils.RandomSellOrder()
				buyChange := newChange(buy, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				sellChange := newChange(sell, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				ranker.InsertChange(buyChange)
				ranker.InsertChange(sellChange)
			}
			ranker.OnChangeEpoch(epoch)
			for i := NumberOfOrderPairs / 2; i < NumberOfOrderPairs; i++ {
				buy, sell := testutils.RandomBuyOrder(), testutils.RandomSellOrder()
				buyChange := newChange(buy, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				sellChange := newChange(sell, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				ranker.InsertChange(buyChange)
				ranker.InsertChange(sellChange)
			}

			computations := make([]Computation, 128)
			n := ranker.Computations(computations)
			Ω(n).Should(Equal(NumberOfOrderPairs * NumberOfOrderPairs))
		})

		It("should not crash when epoch change", func() {
			for i := 0; i < len(EpochBlockNumbers); i++ {
				for j := 0; j < NumberOfOrderPairs; j++ {
					buy, sell := testutils.RandomBuyOrder(), testutils.RandomSellOrder()
					buyChange := newChange(buy, order.Open, orderbook.Priority(j), EpochBlockNumbers[i])
					sellChange := newChange(sell, order.Open, orderbook.Priority(j), EpochBlockNumbers[i])
					ranker.InsertChange(buyChange)
					ranker.InsertChange(sellChange)
				}

				epoch := newEpoch(i, addr)
				ranker.OnChangeEpoch(epoch)
				computations := make([]Computation, 128)
				n := ranker.Computations(computations)
				Ω(n).Should(Equal(NumberOfOrderPairs * NumberOfOrderPairs))
			}
		})

		It("should be handle orders from previous epoch", func() {
			for i := 1; i < len(EpochBlockNumbers); i++ {
				epoch := newEpoch(i, addr)
				ranker.OnChangeEpoch(epoch)

				for j := 0; j < NumberOfOrderPairs; j++ {
					buy, sell := testutils.RandomBuyOrder(), testutils.RandomSellOrder()
					buyChange := newChange(buy, order.Open, orderbook.Priority(j), EpochBlockNumbers[i-1])
					sellChange := newChange(sell, order.Open, orderbook.Priority(j), EpochBlockNumbers[i-1])
					ranker.InsertChange(buyChange)
					ranker.InsertChange(sellChange)
				}

				computations := make([]Computation, 128)
				n := ranker.Computations(computations)
				Ω(n).Should(Equal(NumberOfOrderPairs * NumberOfOrderPairs))
			}
		})

		It("should be ignore ordres from 2 epochs ago", func() {
			epoch := newEpoch(1, addr)
			ranker.OnChangeEpoch(epoch)

			for i := 2; i < len(EpochBlockNumbers); i++ {
				epoch := newEpoch(i, addr)
				ranker.OnChangeEpoch(epoch)

				for j := 0; j < NumberOfOrderPairs; j++ {
					buy, sell := testutils.RandomBuyOrder(), testutils.RandomSellOrder()
					buyChange := newChange(buy, order.Open, orderbook.Priority(j), EpochBlockNumbers[i-2])
					sellChange := newChange(sell, order.Open, orderbook.Priority(j), EpochBlockNumbers[i-2])
					ranker.InsertChange(buyChange)
					ranker.InsertChange(sellChange)
				}

				computations := make([]Computation, 128)
				n := ranker.Computations(computations)
				Ω(n).Should(Equal(0))
			}
		})

	})

	Context("negative tests", func() {
		var wrongAddr identity.Address

		BeforeEach(func() {
			addr, err = testutils.RandomAddress()
			Ω(err).ShouldNot(HaveOccurred())
			epoch = newEpoch(0, addr)
			done = make(chan struct{})
			storer = testutils.NewStorer()
			wrongAddr, err = testutils.RandomAddress()
			Ω(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			close(done)
		})

		It("should error with a wrong epoch", func() {
			ranker, err = NewRanker(done, wrongAddr, storer, epoch)
			Ω(err).Should(HaveOccurred())
		})

		It("should ignore the epoch when the ranker is not in and print out an error message ", func() {
			ranker, err = NewRanker(done, addr, storer, epoch)
			Ω(err).ShouldNot(HaveOccurred())
			epoch := newEpoch(1, wrongAddr)
			ranker.OnChangeEpoch(epoch)
		})
	})
})

// newEpoch returns a new epoch with only one pod and one darknode.
func newEpoch(i int, node identity.Address) cal.Epoch {
	return cal.Epoch{
		Hash: testutils.Random32Bytes(),
		Pods: []cal.Pod{
			{
				Position:  0,
				Hash:      testutils.Random32Bytes(),
				Darknodes: []identity.Address{node},
			},
		},
		Darknodes:   []identity.Address{node},
		BlockNumber: EpochBlockNumbers[i],
	}
}

// newChange returns a new change which can be passed to the ranker.
func newChange(ord order.Order, status order.Status, priority orderbook.Priority, blockNumber uint) orderbook.Change {
	return orderbook.Change{
		OrderID:       ord.ID,
		OrderParity:   ord.Parity,
		OrderStatus:   status,
		OrderPriority: priority,
		BlockNumber:   blockNumber,
	}
}
