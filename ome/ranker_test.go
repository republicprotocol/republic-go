package ome_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/order"
)

var _ = Describe("OME Ranker", func() {

	Context("when processing buy and sell orders", func() {

		It("should append orders to ranker computations", func() {
			numberOfOrderPairs := 10
			ranker, computations := sendCorrectOrdersToRanker(numberOfOrderPairs)
			Expect(ranker.Computations(computations)).Should(Equal(numberOfOrderPairs * 2))
		})

		It("should remove orders", func() {
			numberOfOrderPairs := 10
			ranker, computations := sendCorrectOrdersToRanker(numberOfOrderPairs)

			// Remove half of the orders
			removeOrderIDs := make([]order.ID, numberOfOrderPairs)
			for i := 0; i < numberOfOrderPairs; i++ {
				removeOrderIDs[i] = [32]byte{byte(i)}
			}
			ranker.Remove(removeOrderIDs...)

			Expect(ranker.Computations(computations)).Should(Equal(numberOfOrderPairs))
		})

		It("should deny a buy-sell order pair that is not meant for the ranker", func() {
			ranker := NewRanker(5, 4)

			// (5+5) => 10 mod 5 != 4
			insertSellOrder([32]byte{byte(1)}, Priority(5), ranker)
			insertBuyOrder([32]byte{byte(2)}, Priority(5), ranker)

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

func sendCorrectOrdersToRanker(numberOfOrderPairs int) (Ranker, Computations) {
	ranker := NewRanker(5, 4)

	orderID := 0
	for i := 0; i < numberOfOrderPairs; i++ {
		insertSellOrder([32]byte{byte(orderID)}, Priority(i), ranker)
		orderID++
		insertBuyOrder([32]byte{byte(orderID)}, Priority(9-i), ranker)
		orderID++
	}

	computations := make([]Computation, numberOfOrderPairs*numberOfOrderPairs)
	return ranker, computations
}
