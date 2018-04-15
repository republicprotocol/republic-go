package main_test

import (
	"context"
	"log"
	"sync"
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/cmd/darknode"
	. "github.com/republicprotocol/republic-go/darknode_test"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/republic-go/darknode"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
	"github.com/republicprotocol/republic-go/ethereum/ganache"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/rpc"
	"github.com/republicprotocol/republic-go/stackint"
)

const (
	GanacheRPC                 = "http://localhost:8545"
	NumberOfDarkNodes          = 5
	NumberOfBootstrapDarkNodes = 5
	NumberOfOrders             = 1
)

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
			darkNodeRegistry, err = contracts.NewDarkNodeRegistry(context.Background(), connection, ganache.GenesisTransactor(), &bind.CallOpts{})
			Ω(err).ShouldNot(HaveOccurred())
			darkNodeRegistry.SetGasLimit(1000000)

			// Create DarkNodes and contexts/cancels for running them
			darkNodes, ctxs, cancels, err = NewDarkNodes(NumberOfDarkNodes, NumberOfBootstrapDarkNodes)
			Ω(err).ShouldNot(HaveOccurred())

			// Register the DarkNodes
			err := RegisterDarkNodes(darkNodes, darkNodeRegistry)
			Ω(err).ShouldNot(HaveOccurred())

			shutdown = make(chan struct{})
			go func() {
				defer close(shutdown)

				do.CoForAll(darkNodes, func(i int) {
					darkNodes[i].Run(ctxs[i])
				})
			}()

			// Wait for the Darknodes to boot
			time.Sleep(time.Second)
		})

		AfterEach(func() {
			// Wait for the DarkNodes to shutdown
			<-shutdown

			// Deregister the DarkNodes
			err := DeregisterDarkNodes(darkNodes, darkNodeRegistry)
			Ω(err).ShouldNot(HaveOccurred())

			// Refund the DarkNodes
			err := RefundDarkNodes(darkNodes, darkNodeRegistry)
			Ω(err).ShouldNot(HaveOccurred())
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

	FContext("when computing order matches", func() {
		var darkNodeRegistry contracts.DarkNodeRegistry
		var darkNodes darknode.DarkNodes
		var ctxs []context.Context
		var cancels []context.CancelFunc
		var shutdown chan struct{}

		BeforeEach(func() {
			// Bind to the DarkNodeRegistry contract
			connection, err := ganache.Connect(GanacheRPC)
			Ω(err).ShouldNot(HaveOccurred())
			darkNodeRegistry, err = contracts.NewDarkNodeRegistry(context.Background(), connection, ganache.GenesisTransactor(), &bind.CallOpts{})
			Ω(err).ShouldNot(HaveOccurred())
			darkNodeRegistry.SetGasLimit(1000000)

			// Create DarkNodes and contexts/cancels for running them
			darkNodes, ctxs, cancels, err = NewLocalDarkNodes(NumberOfDarkNodes, NumberOfBootstrapDarkNodes)
			Ω(err).ShouldNot(HaveOccurred())

			shutdown = make(chan struct{}, 0)
			go func() {
				defer close(shutdown)

				do.CoForAll(darkNodes, func(i int) {
					darkNodes[i].Run(ctxs[i])
				})
			}()

			// Wait for the DarkNodes to boot
			time.Sleep(time.Second)
		})

		AfterEach(func() {
			// Wait for the DarkNodes to shutdown
			<-shutdown
			for _, node := range darkNodes {
				node.Stop()
			}
		})

		It("should process the distribute order table in parallel with other pools", func() {
			By("start sending orders ")
			err := sendOrders(darkNodes, NumberOfOrders)
			Ω(err).ShouldNot(HaveOccurred())
			By("finish sending orders ")

			time.Sleep(time.Minute)
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

func sendOrders(nodes darknode.DarkNodes, numberOfOrders int) error {
	// Generate buy-sell order pairs
	buyOrders, sellOrders := make([]*order.Order, numberOfOrders), make([]*order.Order, numberOfOrders)
	for i := 0; i < numberOfOrders; i++ {
		price := i * 1000000000000
		amount := i * 1000000000000
		sellOrder := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour),
			order.CurrencyCodeETH, order.CurrencyCodeBTC, stackint.FromUint(uint(price)), stackint.FromUint(uint(amount)),
			stackint.FromUint(uint(amount)), stackint.FromUint(1))
		sellOrders[i] = sellOrder

		buyOrder := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour),
			order.CurrencyCodeETH, order.CurrencyCodeBTC, stackint.FromUint(uint(price)), stackint.FromUint(uint(amount)),
			stackint.FromUint(uint(amount)), stackint.FromUint(1))
		buyOrders[i] = buyOrder
	}

	// Send order fragment to the nodes
	totalNodes := len(nodes)
	trader, _ := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f")
	prime, _ := stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")
	pool := rpc.NewClientPool(trader).WithTimeout(5 * time.Second).WithTimeoutBackoff(3 * time.Second)

	for i := 0; i < numberOfOrders; i++ {
		buyOrder, sellOrder := buyOrders[i], sellOrders[i]
		log.Printf("Sending matched order. [BUY] %s <---> [SELL] %s", buyOrder.ID, sellOrder.ID)
		buyShares, err := buyOrder.Split(int64(totalNodes), int64(totalNodes*2/3+1), &prime)
		if err != nil {
			return err
		}
		sellShares, err := sellOrder.Split(int64(totalNodes), int64(totalNodes*2/3+1), &prime)
		if err != nil {
			return err
		}

		for _, shares := range [][]*order.Fragment{buyShares, sellShares} {
			do.CoForAll(shares, func(j int) {
				orderRequest := &rpc.OpenOrderRequest{
					From: &rpc.MultiAddress{
						Signature:    []byte{},
						MultiAddress: trader.String(),
					},
					OrderFragment: rpc.MarshalOrderFragment(shares[j]),
				}
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				err := pool.OpenOrder(ctx, nodes[j].NetworkOption.MultiAddress, orderRequest)
				if err != nil {
					log.Printf("Coudln't send order fragment to %s\n", nodes[j].NetworkOption.MultiAddress.ID())
					log.Fatal(err)
				}
			})
		}
	}

	return nil
}
