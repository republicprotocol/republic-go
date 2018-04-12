package main_test

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/cmd/darknode"
	"github.com/republicprotocol/republic-go/darknode"
	"github.com/republicprotocol/republic-go/darkocean"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
	"github.com/republicprotocol/republic-go/ethereum/ganache"
)

const (
	GanacheRPC                 = "http://localhost:8545"
	NumberOfDarkNodes          = 48
	NumberOfBootstrapDarkNodes = 5
	PoolSize                   = 5
)

var _ = Describe("DarkNode", func() {

	Context("when watching the ocean", func() {

		var darkNodeRegistry contracts.DarkNodeRegistry
		var DarkNodes darknode.DarkNodes
		var ctxs []context.Context
		var cancels []context.CancelFunc
		var shutdown chan struct{}

		BeforeEach(func() {

			// Bind to the DarkNodeRegistry contract
			connection, err := ganache.Connect(GanacheRPC)
			Ω(err).ShouldNot(HaveOccurred())
			darkNodeRegistry, err = contracts.NewDarkNodeRegistry(context.Background(), connection, ganache.GenesisTransactor(), &bind.CallOpts{})
			Ω(err).ShouldNot(HaveOccurred())

			// Create DarkNodes and contexts/cancels for running them
			DarkNodes, ctxs, cancels = NewLocalDarkNodes(NumberOfDarkNodes, NumberOfBootstrapDarkNodes)

			shutdown = make(chan struct{})

			var wg sync.WaitGroup
			wg.Add(len(DarkNodes))
			for i := range DarkNodes {
				go func(i int) {
					defer wg.Done()

					DarkNodes[i].Run(ctxs[i])
				}(i)
			}

			go func() {
				defer close(shutdown)

				wg.Wait()
			}()

			// Wait for the DarkNodes to boot
			time.Sleep(time.Second)
		})

		AfterEach(func() {

			// Wait for the DarkNodes to shutdown
			<-shutdown
		})

		It("should update local views of the ocean", func() {
			numberOfEpochs := 2
			oceans := make(darkocean.Oceans, NumberOfDarkNodes)

			By("start calling epoch")
			for j := 0; j < numberOfEpochs; j++ {
				// Store all DarkOceans before the turn of the epoch
				for i := range DarkNodes {
					oceans[i] = DarkNodes[i].Ocean()
				}

				// Turn the epoch
				_, err := darkNodeRegistry.Epoch()
				Ω(err).ShouldNot(HaveOccurred())

				// Wait for DarkNodes to receive a notification and reconfigure
				// themselves
				time.Sleep(time.Second)

				// Verify that all DarkOceans have changed
				for i := range DarkNodes {
					Ω(oceans[i].Equal(DarkNodes[i].Ocean())).Should(BeFalse())
				}
			}

			By("stop all the nodes")
			// Cancel all DarkNodes
			for i := range DarkNodes {
				cancels[i]()
			}
		})

		It("should converge on a global view of the ocean", func() {

			// Turn the epoch
			_, err := darkNodeRegistry.Epoch()
			Ω(err).ShouldNot(HaveOccurred())

			// Wait for DarkNodes to receive a notification and reconfigure
			// themselves
			time.Sleep(time.Second)

			// Verify that all DarkNodes have converged on the DarkOcean
			ocean, err := darkocean.NewOcean(PoolSize, darkNodeRegistry)
			Ω(err).ShouldNot(HaveOccurred())
			for i := range DarkNodes {
				Ω(ocean.Equal(DarkNodes[i].Ocean())).Should(BeTrue())
			}

			// Cancel all DarkNodes
			for i := range DarkNodes {
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
