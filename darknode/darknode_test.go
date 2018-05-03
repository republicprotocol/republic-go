package darknode_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/darknode"

	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/rpc/client"
	"github.com/republicprotocol/republic-go/rpc/relayer"
	"github.com/republicprotocol/republic-go/rpc/smpcer"
	"github.com/republicprotocol/republic-go/smpc"
	"github.com/republicprotocol/republic-go/stackint"
)

// TODO: Regression testing for deadlocks when synchronizing the orderbook.

var _ = Describe("Darknode", func() {

	Context("when opening orders", func() {

		It("should update the orderbook with an open order", func() {

			// Create a relayer client to sync with the orderbook
			crypter := crypto.NewWeakCrypter()

			conn, err := client.Dial(context.Background(), env.Darknodes[0].MultiAddress())
			Expect(err).ShouldNot(HaveOccurred())

			defer conn.Close()

			traderKeystore, err := crypto.RandomKeystore()
			Expect(err).ShouldNot(HaveOccurred())
			traderAddr := identity.Address(traderKeystore.Address())

			relayClient := relayer.NewRelayClient(conn.ClientConn)
			requestSignature, err := crypter.Sign(traderAddr)
			Expect(err).ShouldNot(HaveOccurred())

			request := &relayer.SyncRequest{
				Signature: requestSignature,
				Address:   traderAddr.String(),
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			stream, err := relayClient.Sync(ctx, request)
			Expect(err).ShouldNot(HaveOccurred())

			// Create order fragment to send
			n := int64(17)
			k := int64(12)
			primeVal, _ := stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")
			prime := &primeVal
			price := stackint.FromUint(10)
			minVolume := stackint.FromUint(100)
			maxVolume := stackint.FromUint(1000)
			nonce := stackint.One()

			fragments, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, price, maxVolume, minVolume, nonce).Split(n, k, prime)
			Expect(err).ShouldNot(HaveOccurred())

			env.Darknodes[0].OnOpenOrder(env.Darknodes[0].MultiAddress(), fragments[0])

			stream.Recv()
		})

		It("should not match orders that are incompatible", func() {

		})

		It("should match orders that are compatible", func() {

		})

		It("should confirm orders that are compatible", func() {

		})

		It("should release orders that conflict with a confirmed order", func() {

		})

		It("should reject orders from unregistered addresses", func() {
			// FIXME: Implement
		})
	})

	Context("when computing order matches", func() {

		It("should process the distribute order table in parallel with other pools", func() {
			for _, node := range env.Darknodes {
				log.Printf("%v has %v peers", node.Address(), len(node.RPC().SwarmerClient().DHT().MultiAddresses()))
			}

			By("sending orders...")
			err := sendOrders(env.Darknodes, NumberOfOrdersPerSecond)
			Expect(err).ShouldNot(HaveOccurred())

			By("verifying that nodes found matches...")

			crypter := crypto.NewWeakCrypter()
			conn, err := client.Dial(context.Background(), env.Darknodes[0].MultiAddress())
			Expect(err).ShouldNot(HaveOccurred())
			defer conn.Close()

			traderKeystore, err := crypto.RandomKeystore()
			Expect(err).ShouldNot(HaveOccurred())
			traderAddr := identity.Address(traderKeystore.Address())

			relayClient := relayer.NewRelayClient(conn.ClientConn)
			requestSignature, err := crypter.Sign(traderAddr)
			Expect(err).ShouldNot(HaveOccurred())

			request := &relayer.SyncRequest{
				Signature: requestSignature,
				Address:   traderAddr.String(),
			}
			stream, err := relayClient.Sync(context.Background(), request)
			Expect(err).ShouldNot(HaveOccurred())

			confirmed := map[string]struct{}{}
			for len(confirmed) < NumberOfOrdersPerSecond {
				syncResp, err := stream.Recv()
				Expect(err).ShouldNot(HaveOccurred())
				log.Printf("synchronizing entry %v => %v", base58.Encode(syncResp.Entry.Order.OrderId), syncResp.Entry.OrderStatus)
				if syncResp.Entry.OrderStatus == relayer.OrderStatus_Confirmed {
					confirmed[string(syncResp.Entry.Order.OrderId)] = struct{}{}
				}
			}

			log.Println("PASSED!")
		})

		It("should update the order book after computing an order match", func() {

		})

	})
})

func sendOrders(nodes Darknodes, numberOfOrders int) error {

	// Generate buy-sell order pairs
	buyOrders, sellOrders := make([]*order.Order, numberOfOrders), make([]*order.Order, numberOfOrders)
	for i := 0; i < numberOfOrders; i++ {
		price := i * 1000000000000
		amount := i * 1000000000000

		nonce, err := stackint.Random(rand.Reader, &smpc.Prime)
		if err != nil {
			return err
		}

		sellOrder := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour),
			order.CurrencyCodeETH, order.CurrencyCodeBTC, stackint.FromUint(uint(price)), stackint.FromUint(uint(amount)),
			stackint.FromUint(uint(amount)), nonce)
		sellOrders[i] = sellOrder

		buyOrder := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour),
			order.CurrencyCodeETH, order.CurrencyCodeBTC, stackint.FromUint(uint(price)), stackint.FromUint(uint(amount)),
			stackint.FromUint(uint(amount)), nonce)
		buyOrders[i] = buyOrder
	}

	// Send order fragment to the nodes
	totalNodes := len(nodes)
	traderKeystore, err := crypto.RandomKeystore()
	Expect(err).ShouldNot(HaveOccurred())
	traderAddr := traderKeystore.Address()
	trader, _ := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/80/republic/%s", traderAddr))
	prime, _ := stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")

	crypter := crypto.NewWeakCrypter()
	connPool := client.NewConnPool(256)
	defer connPool.Close()
	smpcerClient := smpcer.NewClient(&crypter, trader, &connPool)

	for i := 0; i < numberOfOrders; i++ {
		buyOrder, sellOrder := buyOrders[i], sellOrders[i]
		log.Printf("sending buy/sell pair (%s, %s)", buyOrder.ID, sellOrder.ID)
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
