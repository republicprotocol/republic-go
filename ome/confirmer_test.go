package ome_test

import (
	"log"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/testutils"
)

var numberOfComputationsToTest = 100

var _ = Describe("Confirmer", func() {
	var confirmer Confirmer
	var renLedger cal.RenLedger
	var storer Storer

	BeforeEach(func() {
		depth, pollInterval := uint(0), time.Second
		renLedger = testutils.NewRenLedger()
		storer = testutils.NewStorer()
		confirmer = NewConfirmer(storer, renLedger, pollInterval, depth)
	})

	FIt("should be able to confirm order on the ren ledger", func(d Done) {
		defer close(d)

		done := make(chan struct{})
		orderMatches := make(chan Computation)
		orderIDs := map[[32]byte]struct{}{}
		computations := make([]Computation, numberOfComputationsToTest)
		for i := 0; i < numberOfComputationsToTest; i++ {
			computations[i] = testutils.RandomComputation()
			orderIDs[computations[i].Buy] = struct{}{}
			orderIDs[computations[i].Sell] = struct{}{}
		}

		// Open all the orders
		for i := 0; i < numberOfComputationsToTest; i++ {
			err := renLedger.OpenBuyOrder([65]byte{}, computations[i].Buy)
			Expect(err).ShouldNot(HaveOccurred())
			err = renLedger.OpenSellOrder([65]byte{}, computations[i].Sell)
			Expect(err).ShouldNot(HaveOccurred())
		}

		go func() {
			defer GinkgoRecover()
			defer close(done)

			for i := 0; i < numberOfComputationsToTest; i++ {
				log.Println(i)
				orderMatches <- computations[i]
			}
		}()

		confirmedMatches, errs := confirmer.Confirm(done, orderMatches)
		go func() {
			defer GinkgoRecover()

			for err := range errs {
				Ω(err).ShouldNot(HaveOccurred())
			}
		}()

		for match := range confirmedMatches {
			_, ok := orderIDs[match.Buy]
			Ω(ok).Should(BeTrue())
			delete(orderIDs, match.Buy)

			_, ok = orderIDs[match.Sell]
			Ω(ok).Should(BeTrue())
			delete(orderIDs, match.Sell)
		}

		Ω(len(orderIDs)).Should(Equal(0))
	}, 100)
})
