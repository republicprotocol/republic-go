package main_test

import (
	"sync"
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/cmd/darknode"
	"github.com/republicprotocol/republic-go/darknode"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
	"github.com/republicprotocol/republic-go/ethereum/ganache"
)

const ganacheRPC = "http://localhost:8545"
const numberOfDarknodes = 48
const numberOfBootstrapDarknodes = 5

var _ = Describe("Darknode", func() {

	Context("when watching the ocean", func() {

		var err error
		var darkNodeRegistry contracts.DarkNodeRegistry
		var darkNodes, ctxs, cancels darknode.Darknodes, []context.Context, []context.CancelFunc
		var shutdown chan struct{}

		BeforeEach(func() {

			// Bind to the DarkNodeRegistry contract
			darkNodeRegistry, err = contracts.NewDarkNodeRegistry(context.Background(), ganache.Connect(ganacheRPC), ganache.GenesisTransactor(), &bind.CallOpts{})
			Ω(err).ShouldNot(HaveOccurred())

			// Create Darknodes and contexts/cancels for running them
			darknodes, ctxs, cancels = NewLocalDarknodes()
			shutdown = make(chan struct{})

			var wg sync.WaitGroup
			wg.Add(len(darknodes))

			for i := range darknodes {
				go func(i int) {
					defer wg.Done()

					darknodes[i].Run(ctx[i])
				}(i)
			}
			go func() {
				defer close(shutdown)
	
				wg.Wait()
			}()

			// Wait for the Darknodes to boot
			time.Sleep(time.Second)
		})

		AfterEach(func() {

			// Wait for the Darknodes to shutdown
			<-shutdown
		})

		It("should update local views of the ocean", func() {
			numberOfEpochs := 2
			oceans := make(darkocean.DarkOceans, numberOfDarknodes)

			for j := 0; j < numberOfEpochs; j++ {
				// Store all DarkOceans before the turn of the epoch
				for i := range darknodes {
					oceans[i] = darknodes[i].DarkOcean()
				}

				// Turn the epoch
				_, err := darkNodeRegistry.Epoch()
				Ω(err).ShouldNot(HaveOccurred())

				// Wait for Darknodes to receive a notification and reconfigure
				// themselves
				time.Sleep(time.Second)

				// Verify that all DarkOceans have changed
				for i := range darknodes {
					Ω(oceans[i].Equal(darknodes[i].DarkOcean())).Should(BeFalse())
				}
			}

			// Cancel all Darknodes
			for i := range darknodes {
				cancels[i]()
			}
		})

		It("should converge on a global view of the ocean", func() {

			// Turn the epoch
			_, err := darkNodeRegistry.Epoch()
			Ω(err).ShouldNot(HaveOccurred())

			// Wait for Darknodes to receive a notification and reconfigure
			// themselves
			time.Sleep(time.Second)

			// Verify that all Darknodes have converged on the DarkOcean
			ocean := darkocean.NewDarkOcean(darkNodeRegistry)
			for i := range darknodes {
				Ω(ocean.Equal(darknodes[i].DarkOcean())).Should(BeTrue())
			}

			// Cancel all Darknodes
			for i := range darknodes {
				cancels[i]()
			}
		})

		It("should persist computations from recent epochs", func() {

		})

		It("should not persist computations from distant epochs", func() {

		})
	})

	Context("when computing order matches", func() {

		It("should process the distribute order table in parallel with other pools", func() {

		})

		It("should update the order book after computing an order match", func() {

		})

	})

	Context("when confirming order matches", func() {

		It("should update the order book after confirming an order match", func() {

		})

		It("should update the order book after releasing an order match", func() {

		})
	})
})
