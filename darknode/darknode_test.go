package darknode_test

import (
	"context"
	"io"
	"log"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/rpc/client"
	"github.com/republicprotocol/republic-go/rpc/relayer"
	"github.com/republicprotocol/republic-go/stackint"
)

// TODO: Regression testing for deadlocks when synchronizing the orderbook.

var _ = Describe("Darknode", func() {

	Context("when opening orders", func() {

		It("should update the orderbook with an open order", func(done Done) {
			stream, conn, cancel, err := createTestRelayClient()
			Expect(err).ShouldNot(HaveOccurred())

			go func() {
				defer close(done)
				defer GinkgoRecover()
				defer cancel()
				defer conn.Close()

				// time.Sleep(1*time.Second)
				matched := map[string]struct{}{}
				for len(matched) < 30 {
					syncResp, err := stream.Recv()
					if err != nil {
						if err == io.EOF {
							log.Println("EOF")
							break
						}
						Expect(err).ShouldNot(HaveOccurred())
					}
					log.Printf("received message: %v", syncResp.Entry.OrderStatus)
					if syncResp.Entry.OrderStatus == relayer.OrderStatus_Open {
						matched[string(syncResp.Entry.Order.OrderId)] = struct{}{}
					}
				}
				Expect(len(matched)).Should(Equal(30))
			}()

			err = env.SendOrders(30)
			Expect(err).ShouldNot(HaveOccurred())
		}, 60)

		It("should not match orders that are incompatible", func(done Done) {
			stream, conn, cancel, err := createTestRelayClient()
			Expect(err).ShouldNot(HaveOccurred())

			go func() {
				defer close(done)
				defer GinkgoRecover()
				defer cancel()
				defer conn.Close()

				// time.Sleep(1*time.Second)
				matched := map[string]struct{}{}
				for len(matched) < 30 {
					syncResp, err := stream.Recv()
					if err != nil {
						if err == io.EOF {
							log.Println("EOF")
							break
						}
						Expect(err).ShouldNot(HaveOccurred())
					}
					log.Printf("received message: %v", syncResp.Entry.OrderStatus)
					if syncResp.Entry.OrderStatus == relayer.OrderStatus_Unconfirmed {
						matched[string(syncResp.Entry.Order.OrderId)] = struct{}{}
					}
				}
				Expect(len(matched)).Should(Equal(30))
			}()

			for i := 0; i < 30; i++ {
				fragments, err := createNewOrderFragments(true)
				Expect(err).ShouldNot(HaveOccurred())
				env.Darknodes[0].OnOpenOrder(env.Darknodes[0].MultiAddress(), fragments[0])
			}
		}, 60)

		It("should match orders that are compatible", func(done Done) {
			stream, conn, cancel, err := createTestRelayClient()
			Expect(err).ShouldNot(HaveOccurred())

			go func() {
				defer close(done)
				defer GinkgoRecover()
				defer cancel()
				defer conn.Close()

				// time.Sleep(1*time.Second)
				matched := map[string]struct{}{}
				for len(matched) < 30 {
					syncResp, err := stream.Recv()
					if err != nil {
						if err == io.EOF {
							log.Println("EOF")
							break
						}
						Expect(err).ShouldNot(HaveOccurred())
					}
					log.Printf("received message: %v", syncResp.Entry.OrderStatus)
					if syncResp.Entry.OrderStatus == relayer.OrderStatus_Unconfirmed {
						matched[string(syncResp.Entry.Order.OrderId)] = struct{}{}
					}
				}
				Expect(len(matched)).Should(Equal(30))
			}()

			err = env.SendOrders(30)
			Expect(err).ShouldNot(HaveOccurred())
			time.Sleep(30 * time.Second)

		}, 60 /* 1 minute timeout */)

		It("should confirm orders that are compatible", func(done Done) {
			stream, conn, cancel, err := createTestRelayClient()
			Expect(err).ShouldNot(HaveOccurred())

			go func() {
				defer close(done)
				defer GinkgoRecover()
				defer cancel()
				defer conn.Close()

				// time.Sleep(1*time.Second)
				matched := map[string]struct{}{}
				for len(matched) < 30 {
					syncResp, err := stream.Recv()
					if err != nil {
						if err == io.EOF {
							log.Println("EOF")
							break
						}
						Expect(err).ShouldNot(HaveOccurred())
					}
					log.Printf("received message: %v", syncResp.Entry.OrderStatus)
					if syncResp.Entry.OrderStatus == relayer.OrderStatus_Confirmed {
						matched[string(syncResp.Entry.Order.OrderId)] = struct{}{}
					}
				}
				Expect(len(matched)).Should(Equal(30))
			}()

			err = env.SendOrders(30)
			Expect(err).ShouldNot(HaveOccurred())
			time.Sleep(30 * time.Second)

		}, 60 /* 1 minute timeout */)

		It("should release orders that conflict with a confirmed order", func(done Done) {
			stream, conn, cancel, err := createTestRelayClient()
			Expect(err).ShouldNot(HaveOccurred())

			go func() {
				defer close(done)
				defer GinkgoRecover()
				defer cancel()
				defer conn.Close()

				// time.Sleep(1*time.Second)
				matched := map[string]struct{}{}
				openOrders := map[string]struct{}{}
				for len(matched) < 30 {
					syncResp, err := stream.Recv()
					if err != nil {
						if err == io.EOF {
							log.Println("EOF")
							break
						}
						Expect(err).ShouldNot(HaveOccurred())
					}
					log.Printf("received message: %v", syncResp.Entry.OrderStatus)
					if syncResp.Entry.OrderStatus == relayer.OrderStatus_Open {
						openOrders[string(syncResp.Entry.Order.OrderId)] = struct{}{}
					}
					if syncResp.Entry.OrderStatus == relayer.OrderStatus_Confirmed {
						matched[string(syncResp.Entry.Order.OrderId)] = struct{}{}
					}

				}
				Expect(len(matched)).Should(Equal(20))
			}()

			for i := 0; i < 10; i++ {
				fragments, err := createNewOrderFragments(true)
				Expect(err).ShouldNot(HaveOccurred())
				env.Darknodes[0].OnOpenOrder(env.Darknodes[0].MultiAddress(), fragments[0])
			}
			for i := 0; i < 20; i++ {
				fragments, err := createNewOrderFragments(false)
				Expect(err).ShouldNot(HaveOccurred())
				env.Darknodes[0].OnOpenOrder(env.Darknodes[0].MultiAddress(), fragments[0])
			}

		}, 60 /* 1 minute timeout */)

		It("should reject orders from unregistered addresses", func() {
			// FIXME: Implement
		})
	})

	Context("when synchronizing the orderbook", func() {

		It("should not deadlock when the sync starts before the updates", func() {
			stream1, conn1, cancel1, err := createTestRelayClient()
			defer conn1.Close()
			defer cancel1()
			Expect(err).ShouldNot(HaveOccurred())

			stream2, conn2, cancel2, err := createTestRelayClient()
			defer conn2.Close()
			defer cancel2()
			Expect(err).ShouldNot(HaveOccurred())

			isDeadlocked := true
			go func() {
				defer GinkgoRecover()
				status1 := map[string]struct{}{}
				status2 := map[string]struct{}{}
				for len(status1) < 20 && len(status2) < 20 {
					resp1, err := stream1.Recv()
					if err != nil {
						if err == io.EOF {
							break
						}
						Expect(err).ShouldNot(HaveOccurred())
					}
					status1[string(resp1.Entry.Order.OrderId)] = struct{}{}
					resp2, err := stream2.Recv()
					if err != nil {
						if err == io.EOF {
							break
						}
						Expect(err).ShouldNot(HaveOccurred())
					}
					status2[string(resp2.Entry.Order.OrderId)] = struct{}{}
				}
				isDeadlocked = false
			}()

			err = env.SendOrders(20)
			Expect(err).ShouldNot(HaveOccurred())

			time.Sleep(10 * time.Second)
			Expect(isDeadlocked).To(Equal(false))

		})

		It("should not deadlock when the sync starts during the updates", func() {
			err := env.SendOrders(5)
			Expect(err).ShouldNot(HaveOccurred())
			stream1, conn1, cancel1, err := createTestRelayClient()
			defer conn1.Close()
			defer cancel1()
			Expect(err).ShouldNot(HaveOccurred())
			err = env.SendOrders(10)
			Expect(err).ShouldNot(HaveOccurred())
			stream2, conn2, cancel2, err := createTestRelayClient()
			defer conn2.Close()
			defer cancel2()
			Expect(err).ShouldNot(HaveOccurred())
			err = env.SendOrders(5)
			Expect(err).ShouldNot(HaveOccurred())

			isDeadlocked := true
			go func() {
				defer GinkgoRecover()
				status1 := map[string]struct{}{}
				status2 := map[string]struct{}{}
				for len(status1) < 20 && len(status2) < 20 {
					resp1, err := stream1.Recv()
					if err != nil {
						if err == io.EOF {
							break
						}
						Expect(err).ShouldNot(HaveOccurred())
					}
					status1[string(resp1.Entry.Order.OrderId)] = struct{}{}
					resp2, err := stream2.Recv()
					if err != nil {
						if err == io.EOF {
							break
						}
						Expect(err).ShouldNot(HaveOccurred())
					}
					status2[string(resp2.Entry.Order.OrderId)] = struct{}{}
				}
				isDeadlocked = false
			}()

			time.Sleep(10 * time.Second)
			Expect(isDeadlocked).To(Equal(false))
		})

		It("should not deadlock when the sync starts after the updates", func() {
			err := env.SendOrders(20)
			Expect(err).ShouldNot(HaveOccurred())

			stream1, conn1, cancel1, err := createTestRelayClient()
			defer conn1.Close()
			defer cancel1()
			Expect(err).ShouldNot(HaveOccurred())
			stream2, conn2, cancel2, err := createTestRelayClient()
			defer conn2.Close()
			defer cancel2()
			Expect(err).ShouldNot(HaveOccurred())

			isDeadlocked := true
			go func() {
				defer GinkgoRecover()
				status1 := map[string]struct{}{}
				status2 := map[string]struct{}{}
				for len(status1) < 20 && len(status2) < 20 {
					resp1, err := stream1.Recv()
					if err != nil {
						if err == io.EOF {
							break
						}
						Expect(err).ShouldNot(HaveOccurred())
					}
					status1[string(resp1.Entry.Order.OrderId)] = struct{}{}
					resp2, err := stream2.Recv()
					if err != nil {
						if err == io.EOF {
							break
						}
						Expect(err).ShouldNot(HaveOccurred())
					}
					status2[string(resp2.Entry.Order.OrderId)] = struct{}{}
				}
				isDeadlocked = false
			}()

			time.Sleep(10 * time.Second)
			Expect(isDeadlocked).To(Equal(false))
		})

	})

	// Context("when computing order matches", func() {

	// 	It("should process the distribute order table in parallel with other pools", func() {
	// 		for _, node := range env.Darknodes {
	// 			log.Printf("%v has %v peers", node.Address(), len(node.RPC().SwarmerClient().DHT().MultiAddresses()))
	// 		}

	// 		By("sending orders...")
	// 		err := sendOrders(env.Darknodes, NumberOfOrdersPerSecond)
	// 		Expect(err).ShouldNot(HaveOccurred())

	// 		By("verifying that nodes found matches...")

	// 		crypter := crypto.NewWeakCrypter()
	// 		conn, err := client.Dial(context.Background(), env.Darknodes[0].MultiAddress())
	// 		Expect(err).ShouldNot(HaveOccurred())
	// 		defer conn.Close()

	// 		traderKeystore, err := crypto.RandomKeystore()
	// 		Expect(err).ShouldNot(HaveOccurred())
	// 		traderAddr := identity.Address(traderKeystore.Address())

	// 		relayClient := relayer.NewRelayClient(conn.ClientConn)
	// 		requestSignature, err := crypter.Sign(traderAddr)
	// 		Expect(err).ShouldNot(HaveOccurred())

	// 		request := &relayer.SyncRequest{
	// 			Signature: requestSignature,
	// 			Address:   traderAddr.String(),
	// 		}
	// 		stream, err := relayClient.Sync(context.Background(), request)
	// 		Expect(err).ShouldNot(HaveOccurred())

	// 		confirmed := map[string]struct{}{}
	// 		for len(confirmed) < NumberOfOrdersPerSecond {
	// 			syncResp, err := stream.Recv()
	// 			Expect(err).ShouldNot(HaveOccurred())
	// 			log.Printf("synchronizing entry %v => %v", base58.Encode(syncResp.Entry.Order.OrderId), syncResp.Entry.OrderStatus)
	// 			if syncResp.Entry.OrderStatus == relayer.OrderStatus_Confirmed {
	// 				confirmed[string(syncResp.Entry.Order.OrderId)] = struct{}{}
	// 			}
	// 		}

	// 		log.Println("PASSED!")
	// 	})

	// 	It("should update the order book after computing an order match", func() {

	// 	})

	// })
})

func createTestRelayClient() (relayer.Relay_SyncClient, *client.Conn, context.CancelFunc, error) {
	// Create a relayer client to sync with the orderbook
	crypter := crypto.NewWeakCrypter()

	conn, err := client.Dial(context.Background(), env.Darknodes[0].MultiAddress())
	if err != nil {
		return nil, &(client.Conn{}), nil, err
	}

	traderKeystore, err := crypto.RandomKeystore()
	if err != nil {
		return nil, conn, nil, err
	}
	traderAddr := identity.Address(traderKeystore.Address())

	relayClient := relayer.NewRelayClient(conn.ClientConn)
	requestSignature, err := crypter.Sign(traderAddr)
	if err != nil {
		return nil, conn, nil, err
	}

	request := &relayer.SyncRequest{
		Signature: requestSignature,
		Address:   traderAddr.String(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	stream, err := relayClient.Sync(ctx, request)
	if err != nil {
		return nil, conn, nil, err
	}
	return stream, conn, cancel, nil
}

func createNewOrderFragments(buy bool) ([]*order.Fragment, error) {
	// Create order fragment to send
	n := int64(17)
	k := int64(12)
	primeVal, _ := stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")
	prime := &primeVal
	price := stackint.FromUint(10)
	minVolume := stackint.FromUint(100)
	maxVolume := stackint.FromUint(1000)
	nonce := stackint.One()
	parity := order.ParityBuy
	if !buy {
		parity = order.ParitySell
	}

	fragments, err := order.NewOrder(order.TypeLimit, parity, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, price, maxVolume, minVolume, nonce).Split(n, k, prime)
	if err != nil {
		return nil, err
	}

	return fragments, nil
}
