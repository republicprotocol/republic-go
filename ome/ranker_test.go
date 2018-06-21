package ome_test

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/leveldb"
	. "github.com/republicprotocol/republic-go/ome"
	"github.com/republicprotocol/republic-go/registry"

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
	var storer *leveldb.Store
	var ranker Ranker
	var done chan struct{}
	var epoch registry.Epoch
	var addr identity.Address
	var err error

	Context("when processing inserted changes", func() {
		BeforeEach(func() {
			addr, err = testutils.RandomAddress()
			Ω(err).ShouldNot(HaveOccurred())
			epoch = newEpoch(0, addr)
			done = make(chan struct{})
			storer, err = leveldb.NewStore("./data.out")
			Ω(err).ShouldNot(HaveOccurred())
			ranker, err = NewRanker(done, addr, storer, storer, epoch)
			Ω(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			close(done)
			os.RemoveAll("./data.out")
		})

		It("should create computations with new opened orders", func() {
			for i := 0; i < NumberOfOrderPairs; i++ {
				buy, sell := testutils.RandomBuyOrder(), testutils.RandomSellOrder()
				buyChange := newChange(buy, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				buyChange.Trader = "buyer"
				Ω(storer.PutChange(buyChange)).ShouldNot(HaveOccurred())
				sellChange := newChange(sell, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				sellChange.Trader = "seller"
				Ω(storer.PutChange(sellChange)).ShouldNot(HaveOccurred())
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
				buyChange.Trader = "buyer"
				Ω(storer.PutChange(buyChange)).ShouldNot(HaveOccurred())
				ranker.InsertChange(buyChange)
				cancelBuy := newChange(buy, order.Canceled, orderbook.Priority(i), EpochBlockNumbers[0])
				cancelBuy.Trader = "buyer"
				Ω(storer.PutChange(cancelBuy)).ShouldNot(HaveOccurred())
				ranker.InsertChange(cancelBuy)

				sellChange := newChange(sell, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				sellChange.Trader = "seller"
				Ω(storer.PutChange(sellChange)).ShouldNot(HaveOccurred())
				ranker.InsertChange(sellChange)
				cancelSell := newChange(sell, order.Canceled, orderbook.Priority(i), EpochBlockNumbers[0])
				cancelSell.Trader = "seller"
				Ω(storer.PutChange(cancelBuy)).ShouldNot(HaveOccurred())
				ranker.InsertChange(cancelSell)
			}
			computations := make([]Computation, 128)
			i := ranker.Computations(computations)
			Ω(i).Should(BeZero())
		})

		It("should return number of computations no more the size of the given variable", func() {
			for i := 0; i < NumberOfOrderPairs; i++ {
				buy, sell := testutils.RandomBuyOrder(), testutils.RandomSellOrder()
				buyChange := newChange(buy, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				buyChange.Trader = "buyer"
				Ω(storer.PutChange(buyChange)).ShouldNot(HaveOccurred())
				sellChange := newChange(sell, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				sellChange.Trader = "seller"
				Ω(storer.PutChange(sellChange)).ShouldNot(HaveOccurred())
				ranker.InsertChange(buyChange)
				ranker.InsertChange(sellChange)
			}
			computations := make([]Computation, 10)
			i := ranker.Computations(computations)
			Ω(i).Should(Equal(len(computations)))
		})

		It("should not create computations from orders from same trader", func() {
			for i := 0; i < NumberOfOrderPairs; i++ {
				buy, sell := testutils.RandomBuyOrder(), testutils.RandomSellOrder()
				buyChange := newChange(buy, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				buyChange.Trader = "trader"
				Ω(storer.PutChange(buyChange)).ShouldNot(HaveOccurred())
				sellChange := newChange(sell, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				sellChange.Trader = "trader"
				Ω(storer.PutChange(sellChange)).ShouldNot(HaveOccurred())
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
			storer, err = leveldb.NewStore("./data.out")
			Ω(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			close(done)
			os.RemoveAll("./data.out")
		})

		It("should insert computations from the storer when started", func() {
			for i := 0; i < NumberOfOrderPairs; i++ {
				comp := testutils.RandomComputation()
				Ω(storer.PutComputation(comp)).ShouldNot(HaveOccurred())
			}

			ranker, err = NewRanker(done, addr, storer, storer, epoch)
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
			storer, err = leveldb.NewStore("./data.out")
			Ω(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			close(done)
			os.RemoveAll("./data.out")
		})

		It("should only generate computations meant for its position", func() {
			// add one more pod in the epoch
			another, err := testutils.RandomAddress()
			Ω(err).ShouldNot(HaveOccurred())
			anotherPod := registry.Pod{
				Position:  0,
				Hash:      testutils.Random32Bytes(),
				Darknodes: []identity.Address{another},
			}
			epoch.Pods = append(epoch.Pods, anotherPod)
			epoch.Darknodes = append(epoch.Darknodes, another)
			ranker, err = NewRanker(done, addr, storer, storer, epoch)
			Ω(err).ShouldNot(HaveOccurred())

			for i := 0; i < NumberOfOrderPairs; i++ {
				buy, sell := testutils.RandomBuyOrder(), testutils.RandomSellOrder()
				buyChange := newChange(buy, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				buyChange.Trader = "buyer"
				Ω(storer.PutChange(buyChange)).ShouldNot(HaveOccurred())
				sellChange := newChange(sell, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				sellChange.Trader = "seller"
				Ω(storer.PutChange(sellChange)).ShouldNot(HaveOccurred())
				ranker.InsertChange(buyChange)
				ranker.InsertChange(sellChange)
			}

			computations := make([]Computation, 128)
			i := ranker.Computations(computations)
			Ω(i).Should(BeNumerically(">=", NumberOfOrderPairs*NumberOfOrderPairs/2))
		})
	})

	Context("when epoch happens", func() {
		BeforeEach(func() {
			addr, err = testutils.RandomAddress()
			Ω(err).ShouldNot(HaveOccurred())
			epoch = newEpoch(0, addr)
			done = make(chan struct{})
			storer, err = leveldb.NewStore("./data.out")
			Ω(err).ShouldNot(HaveOccurred())
			ranker, err = NewRanker(done, addr, storer, storer, epoch)
			Ω(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			close(done)
			os.RemoveAll("./data.out")
		})

		It("should ignore same epoch change event", func() {
			ranker.OnChangeEpoch(epoch)
			for i := 0; i < NumberOfOrderPairs/2; i++ {
				buy, sell := testutils.RandomBuyOrder(), testutils.RandomSellOrder()
				buyChange := newChange(buy, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				buyChange.Trader = "buyer"
				Ω(storer.PutChange(buyChange)).ShouldNot(HaveOccurred())
				sellChange := newChange(sell, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				sellChange.Trader = "seller"
				Ω(storer.PutChange(sellChange)).ShouldNot(HaveOccurred())
				ranker.InsertChange(buyChange)
				ranker.InsertChange(sellChange)
			}
			ranker.OnChangeEpoch(epoch)
			for i := NumberOfOrderPairs / 2; i < NumberOfOrderPairs; i++ {
				buy, sell := testutils.RandomBuyOrder(), testutils.RandomSellOrder()
				buyChange := newChange(buy, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				buyChange.Trader = "buyer"
				Ω(storer.PutChange(buyChange)).ShouldNot(HaveOccurred())
				sellChange := newChange(sell, order.Open, orderbook.Priority(i), EpochBlockNumbers[0])
				sellChange.Trader = "seller"
				Ω(storer.PutChange(sellChange)).ShouldNot(HaveOccurred())
				ranker.InsertChange(buyChange)
				ranker.InsertChange(sellChange)
			}

			computations := make([]Computation, 128)
			n := ranker.Computations(computations)
			Ω(n).Should(Equal(NumberOfOrderPairs * NumberOfOrderPairs))
		})

		It("should not crash when epoch change", func() {
			for i := 0; i < len(EpochBlockNumbers)-1; i++ {
				for j := 0; j < NumberOfOrderPairs; j++ {
					buy, sell := testutils.RandomBuyOrder(), testutils.RandomSellOrder()
					buyChange := newChange(buy, order.Open, orderbook.Priority(j), EpochBlockNumbers[i])
					buyChange.Trader = "buyer"
					Ω(storer.PutChange(buyChange)).ShouldNot(HaveOccurred())
					sellChange := newChange(sell, order.Open, orderbook.Priority(j), EpochBlockNumbers[i])
					sellChange.Trader = "seller"
					Ω(storer.PutChange(sellChange)).ShouldNot(HaveOccurred())
					ranker.InsertChange(buyChange)
					ranker.InsertChange(sellChange)
				}

				epoch := newEpoch(i+1, addr)
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
					buyChange.Trader = "buyer"
					Ω(storer.PutChange(buyChange)).ShouldNot(HaveOccurred())
					sellChange := newChange(sell, order.Open, orderbook.Priority(j), EpochBlockNumbers[i-1])
					sellChange.Trader = "seller"
					Ω(storer.PutChange(sellChange)).ShouldNot(HaveOccurred())
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
					buyChange.Trader = "buyer"
					Ω(storer.PutChange(buyChange)).ShouldNot(HaveOccurred())
					sellChange := newChange(sell, order.Open, orderbook.Priority(j), EpochBlockNumbers[i-2])
					sellChange.Trader = "buyer"
					Ω(storer.PutChange(sellChange)).ShouldNot(HaveOccurred())
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
			storer, err = leveldb.NewStore("./data.out")
			Ω(err).ShouldNot(HaveOccurred())
			wrongAddr, err = testutils.RandomAddress()
			Ω(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			close(done)
			os.RemoveAll("./data.out")
		})

		It("should error with a wrong epoch", func() {
			ranker, err = NewRanker(done, wrongAddr, storer, storer, epoch)
			Ω(err).Should(HaveOccurred())
		})

		It("should ignore the epoch when the ranker is not in and print out an error message ", func() {
			ranker, err = NewRanker(done, addr, storer, storer, epoch)
			Ω(err).ShouldNot(HaveOccurred())
			epoch := newEpoch(1, wrongAddr)
			ranker.OnChangeEpoch(epoch)
		})
	})
})

// newEpoch returns a new epoch with only one pod and one darknode.
func newEpoch(i int, node identity.Address) registry.Epoch {
	return registry.Epoch{
		Hash: testutils.Random32Bytes(),
		Pods: []registry.Pod{
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
