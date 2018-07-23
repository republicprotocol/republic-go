package ome_test

import (
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/testutils"
)

var numberOfComputationsToTest = 100

var _ = Describe("Confirmer", func() {

	It("should be able to confirm order on the ren ledger", func(d Done) {
		defer close(d)

		var err error
		depth, pollInterval := uint(0), time.Second
		contract := newOmeBinder()
		db, err := leveldb.NewStore("./data1.out", time.Hour)
		defer os.RemoveAll("./data1.out")
		Expect(err).ShouldNot(HaveOccurred())
		storer := db.SomerComputationStore()
		confirmer := NewConfirmer(storer, contract, pollInterval, depth)

		done := make(chan struct{})
		orderMatches := make(chan Computation)
		orderIDs := map[[32]byte]struct{}{}
		computations := make([]Computation, numberOfComputationsToTest)

		for i := 0; i < numberOfComputationsToTest; i++ {
			computations[i], err = testutils.RandomComputation()
			Expect(err).ShouldNot(HaveOccurred())
			orderIDs[computations[i].Buy.OrderID] = struct{}{}
			orderIDs[computations[i].Sell.OrderID] = struct{}{}
			storer.PutComputation(computations[i])

			// Open the orders
			err := contract.OpenBuyOrder([65]byte{}, computations[i].Buy.OrderID)
			Expect(err).ShouldNot(HaveOccurred())
			err = contract.OpenSellOrder([65]byte{}, computations[i].Sell.OrderID)
			Expect(err).ShouldNot(HaveOccurred())
		}

		go func() {
			defer GinkgoRecover()
			defer close(done)

			for i := 0; i < numberOfComputationsToTest; i++ {
				orderMatches <- computations[i]
			}
			time.Sleep(5 * time.Second)
		}()

		confirmedMatches, errs := confirmer.Confirm(done, orderMatches)

		go dispatch.CoBegin(
			func() {
				defer GinkgoRecover()

				for err := range errs {
					fmt.Println(err)
				}
			},
			func() {
				defer GinkgoRecover()

				for match := range confirmedMatches {
					_, ok := orderIDs[match.Buy.OrderID]
					Ω(ok).Should(BeTrue())
					delete(orderIDs, match.Buy.OrderID)

					_, ok = orderIDs[match.Sell.OrderID]
					Ω(ok).Should(BeTrue())
					delete(orderIDs, match.Sell.OrderID)
				}
				Ω(len(orderIDs)).Should(Equal(0))
			})
	}, 100)

	It("should return error for invalid computations", func() {
		var err error
		depth, pollInterval := uint(0), time.Second
		contract := newOmeBinder()
		db, err := leveldb.NewStore("./data2.out", time.Hour)
		defer os.RemoveAll("./data2.out")
		Expect(err).ShouldNot(HaveOccurred())
		storer := db.SomerComputationStore()
		confirmer := NewConfirmer(storer, contract, pollInterval, depth)

		done := make(chan struct{})
		orderMatches := make(chan Computation)
		orderIDs := map[[32]byte]struct{}{}
		computations := make([]Computation, numberOfComputationsToTest)

		for i := 0; i < numberOfComputationsToTest; i++ {
			computations[i], err = testutils.RandomComputation()
			computations[i].EpochDepth = order.FragmentEpochDepth(100)
			Expect(err).ShouldNot(HaveOccurred())
			orderIDs[computations[i].Buy.OrderID] = struct{}{}
			orderIDs[computations[i].Sell.OrderID] = struct{}{}
			storer.PutComputation(computations[i])
		}

		go func() {
			defer GinkgoRecover()
			defer close(done)

			for i := 0; i < numberOfComputationsToTest; i++ {
				orderMatches <- computations[i]
			}
			time.Sleep(5 * time.Second)
		}()

		confirmedMatches, errs := confirmer.Confirm(done, orderMatches)
		go func() {
			defer GinkgoRecover()

			for err := range errs {
				Ω(err).Should(HaveOccurred())
			}
		}()

		Ω(len(confirmedMatches)).Should(BeZero())
	})
})
