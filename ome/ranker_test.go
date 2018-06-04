package ome_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/order"
)

// WARNING: The expected number of computations is only
// (numberOfOrderPairs^2) / numberOfRankers when the numberOfRankers divides
// the numberOfOrderPairs^2 perfectly.
const numberOfRankers = 5
const numberOfOrderPairs = 10

var _ = Describe("OME Ranker", func() {

	Context("when processing buy and sell orders", func() {

		It("should append orders to ranker computations", func() {
			ranker, computations := sendCorrectOrdersToRanker(numberOfRankers, numberOfOrderPairs)
			Expect(ranker.Computations(computations)).Should(Equal((numberOfOrderPairs * numberOfOrderPairs) / numberOfRankers))
		})

		It("should remove all computations if all buy orders are removed", func() {
			ranker, computations := sendCorrectOrdersToRanker(numberOfRankers, numberOfOrderPairs)

			// Remove half of the orders
			removeOrderIDs := make([]order.ID, 0, numberOfOrderPairs)
			for i := 0; i < numberOfOrderPairs; i++ {
				removeOrderIDs = append(removeOrderIDs, [32]byte{byte(i)})
			}
			ranker.Remove(removeOrderIDs...)

			Expect(ranker.Computations(computations)).Should(Equal(0))
		})

		It("should remove all computations if all sell orders are removed", func() {
			ranker, computations := sendCorrectOrdersToRanker(numberOfRankers, numberOfOrderPairs)

			// Remove half of the orders
			removeOrderIDs := make([]order.ID, 0, numberOfOrderPairs)
			for i := numberOfOrderPairs; i < 2*numberOfOrderPairs; i++ {
				removeOrderIDs = append(removeOrderIDs, [32]byte{byte(i)})
			}
			ranker.Remove(removeOrderIDs...)

			Expect(ranker.Computations(computations)).Should(Equal(0))
		})

		It("should quarter the number of computations when halving the number of orders", func() {
			ranker, computations := sendCorrectOrdersToRanker(numberOfRankers, numberOfOrderPairs)

			// Remove half of the orders
			removeOrderIDs := make([]order.ID, 0, numberOfOrderPairs)
			for i := numberOfOrderPairs; i < 2*numberOfOrderPairs; i += 2 {
				removeOrderIDs = append(removeOrderIDs, [32]byte{byte(i)})
			}
			ranker.Remove(removeOrderIDs...)

			Expect(ranker.Computations(computations)).Should(Equal(((numberOfOrderPairs / 2) * (numberOfOrderPairs / 2)) / numberOfRankers))
		})

		It("should not return computations that are not meant for the ranker", func() {
			ranker := NewRanker(numberOfRankers, 0)

			// (1+1) => 2 mod 5 != 0
			insertSellOrder([32]byte{byte(1)}, Priority(numberOfRankers+1), ranker)
			insertBuyOrder([32]byte{byte(2)}, Priority(numberOfRankers+1), ranker)

			// Check that the order-pair was not added
			computations := make([]Computation, 5)
			Expect(ranker.Computations(computations)).Should(Equal(0))
		})
	})
})

func insertSellOrder(orderID order.ID, priority Priority, ranker Ranker) {
	sellOrder := PriorityOrder{
		ID:       orderID,
		Priority: priority,
	}

	// Insert a sell order into the ranker
	ranker.InsertSell(sellOrder)
}

func insertBuyOrder(orderID order.ID, priority Priority, ranker Ranker) {
	buyOrder := PriorityOrder{
		ID:       orderID,
		Priority: priority,
	}

	// Insert a buy order into the ranker
	ranker.InsertBuy(buyOrder)
}

func sendCorrectOrdersToRanker(numberOfRankers, numberOfOrderPairs int) (Ranker, Computations) {
	ranker := NewRanker(numberOfRankers, 0)

	for i := 0; i < numberOfOrderPairs; i++ {
		insertBuyOrder([32]byte{byte(i)}, Priority(i), ranker)
	}
	for i := 0; i < numberOfOrderPairs; i++ {
		insertSellOrder([32]byte{byte(numberOfOrderPairs + i)}, Priority(i), ranker)
	}

	computations := make([]Computation, numberOfOrderPairs*numberOfOrderPairs)
	return ranker, computations
}
