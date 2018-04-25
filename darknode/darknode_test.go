package darknode_test

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/darknode"
	. "github.com/republicprotocol/republic-go/darknodetest"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dispatch"
	ethereum "github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
	"github.com/republicprotocol/republic-go/ethereum/ganache"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/rpc/client"
	"github.com/republicprotocol/republic-go/rpc/smpcer"
	"github.com/republicprotocol/republic-go/stackint"
)

const (
	GanacheRPC                 = "http://localhost:8545"
	NumberOfDarkNodes          = 12
	NumberOfBootstrapDarkNodes = 3
	NumberOfOrders             = 1000
)

var _ = Describe("Darknode", func() {

	var cmd *exec.Cmd
	var conn ethereum.Connection
	var darknodeRegistry contracts.DarkNodeRegistry
	var darknodes Darknodes

	BeforeEach(func() {
		var err error

		cmd, conn, err = ganache.StartAndConnect()
		Expect(err).ShouldNot(HaveOccurred())

		// Connect to Ganache
		darknodeRegistry, err = contracts.NewDarkNodeRegistry(context.Background(), conn, ganache.GenesisTransactor(), &bind.CallOpts{})
		Expect(err).ShouldNot(HaveOccurred())
		darknodeRegistry.SetGasLimit(1000000)

		// Create DarkNodes and contexts/cancels for running them
		darknodes, err = NewDarknodes(NumberOfDarkNodes, NumberOfBootstrapDarkNodes)
		Expect(err).ShouldNot(HaveOccurred())

		// Register the Darknodes and trigger an epoch to accept their
		// registrations
		err = RegisterDarknodes(darknodes, conn, darknodeRegistry)
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		var err error

		// Deregister the DarkNodes
		err = DeregisterDarknodes(darknodes, conn, darknodeRegistry)
		Expect(err).ShouldNot(HaveOccurred())

		// Refund the DarkNodes
		err = RefundDarknodes(darknodes, conn, darknodeRegistry)
		Expect(err).ShouldNot(HaveOccurred())

		cmd.Process.Kill()
		cmd.Wait()
	})

	Context("when watching the ocean", func() {

		// It("should update local views of the ocean", func() {
		// 	numberOfEpochs := 2
		// 	oceans := make(darkocean.DarkOceans, numberOfDarknodes)

		// 	for j := 0; j < numberOfEpochs; j++ {
		// 		// Store all DarkOceans before the turn of the epoch
		// 		for i := range darknodes {
		// 			oceans[i] = darknodes[i].DarkOcean()
		// 		}

		// 		// Turn the epoch
		// 		_, err := darkNodeRegistry.Epoch()
		// 		Ω(err).ShouldNot(HaveOccurred())

		// 		// Wait for Darknodes to receive a notification and reconfigure
		// 		// themselves
		// 		time.Sleep(time.Second)

		// 		// Verify that all DarkOceans have changed
		// 		for i := range darknodes {
		// 			Ω(oceans[i].Equal(darknodes[i].DarkOcean())).Should(BeFalse())
		// 		}
		// 	}

		// 	// Cancel all Darknodes
		// 	for i := range darknodes {
		// 		cancels[i]()
		// 	}
		// })

		// It("should converge on a global view of the ocean", func() {

		// 	// Turn the epoch
		// 	_, err := darkNodeRegistry.Epoch()
		// 	Ω(err).ShouldNot(HaveOccurred())

		// 	// Wait for Darknodes to receive a notification and reconfigure
		// 	// themselves
		// 	time.Sleep(time.Second)

		// 	// Verify that all Darknodes have converged on the DarkOcean
		// 	ocean := darkocean.NewDarkOcean(darkNodeRegistry)
		// 	for i := range darknodes {
		// 		Ω(ocean.Equal(darknodes[i].DarkOcean())).Should(BeTrue())
		// 	}

		// 	// Cancel all Darknodes
		// 	for i := range darknodes {
		// 		cancels[i]()
		// 	}
		// })

		It("should persist computations from recent epochs", func() {

		})

		It("should not persist computations from distant epochs", func() {

		})
	})

	FContext("when computing order matches", func() {

		FIt("should process the distribute order table in parallel with other pools", func(d Done) {
			defer close(d)

			By("booting darknodes...")
			done := make(chan struct{})
			defer close(done)
			go dispatch.CoForAll(darknodes, func(i int) {
				defer GinkgoRecover()
				for err := range darknodes[i].Serve(done) {
					Expect(err).ShouldNot(HaveOccurred())
				}
			})
			time.Sleep(10 * time.Second)
			for _, node := range darknodes {
				log.Printf("%v has %v peers", node.Address(), len(node.RPC().SwarmerClient().DHT().MultiAddresses()))
			}

			By("sending orders...")
			err := sendOrders(darknodes, NumberOfOrders)
			Ω(err).ShouldNot(HaveOccurred())

			By("waiting for results...")
			time.Sleep(time.Minute)
		}, 24*60*60) // Timeout is set to 24 hours

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

func sendOrders(nodes Darknodes, numberOfOrders int) error {

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
	traderAddr, _, err := identity.NewAddress()
	Expect(err).ShouldNot(HaveOccurred())
	trader, _ := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/80/republic/%v", traderAddr))
	prime, _ := stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")

	crypter := crypto.NewWeakCrypter()
	connPool := client.NewConnPool(256)
	smpcerClient := smpcer.NewClient(&crypter, trader, &connPool)

	for i := 0; i < numberOfOrders; i++ {
		buyOrder, sellOrder := buyOrders[i], sellOrders[i]
		log.Printf("Sending matched order. [BUY] %s <---> [SELL] %s", buyOrder.ID, sellOrder.ID)
		buyShares, err := buyOrder.Split(int64(totalNodes), int64((totalNodes+1)*2/3), &prime)
		if err != nil {
			return err
		}
		sellShares, err := sellOrder.Split(int64(totalNodes), int64((totalNodes+1)*2/3), &prime)
		if err != nil {
			return err
		}

		for _, shares := range [][]*order.Fragment{buyShares, sellShares} {
			do.CoForAll(shares, func(j int) {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if err := smpcerClient.OpenOrder(ctx, nodes[j].MultiAddress(), *shares[j]); err != nil {
					log.Printf("cannot send order fragment to %s: %v", nodes[j].Address(), err)
				}
			})
		}
	}

	return nil
}
