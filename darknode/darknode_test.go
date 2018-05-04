package darknode_test

import (
	"context"
	"io"
	"log"
	mathRnd "math/rand"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/darknode"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/rpc/client"
	"github.com/republicprotocol/republic-go/rpc/relayer"
)

// TODO: Regression testing for deadlocks when synchronizing the orderbook.

var _ = Describe("Darknode", func() {

	Context("when opening orders", func() {

		// Serialize tests to prevent bleeding the orderbook.
		testMu := new(sync.Mutex)
		BeforeEach(func() {
			testMu.Lock()
			env.ClearOrderbooks()
		})
		AfterEach(func() {
			testMu.Unlock()
		})

		It("should update the orderbook with an open order", func(done Done) {
			stream, conn, cancel, err := createTestRelayClient()
			Expect(err).ShouldNot(HaveOccurred())

			go func() {
				defer GinkgoRecover()
				defer close(done)
				defer cancel()
				defer conn.Close()

				opened := map[string]struct{}{}
				unconfirmed := map[string]struct{}{}
				confirmed := map[string]struct{}{}
				settled := map[string]struct{}{}
				canceled := map[string]struct{}{}
				for len(opened) < 30 {
					syncResp, err := stream.Recv()
					if err != nil {
						if err == io.EOF {
							break
						}
						Expect(err).ShouldNot(HaveOccurred())
					}
					switch syncResp.Entry.OrderStatus {
					case relayer.OrderStatus_Open:
						opened[string(syncResp.Entry.Order.OrderId)] = struct{}{}
					case relayer.OrderStatus_Unconfirmed:
						unconfirmed[string(syncResp.Entry.Order.OrderId)] = struct{}{}
					case relayer.OrderStatus_Confirmed:
						confirmed[string(syncResp.Entry.Order.OrderId)] = struct{}{}
					case relayer.OrderStatus_Settled:
						settled[string(syncResp.Entry.Order.OrderId)] = struct{}{}
					case relayer.OrderStatus_Canceled:
						canceled[string(syncResp.Entry.Order.OrderId)] = struct{}{}
					default:
						break
					}
				}
				Expect(len(opened)).Should(Equal(30))
				Expect(len(unconfirmed)).Should(Equal(0))
				Expect(len(confirmed)).Should(Equal(0))
				Expect(len(settled)).Should(Equal(0))
				Expect(len(canceled)).Should(Equal(0))
			}()

			err = env.SendMatchingOrderPairs(30)
			Expect(err).ShouldNot(HaveOccurred())
		}, 60)

		It("should not match orders that are incompatible", func(done Done) {

			stream, conn, cancel, err := createTestRelayClient()
			Expect(err).ShouldNot(HaveOccurred())

			go func() {
				defer GinkgoRecover()
				defer close(done)
				defer cancel()
				defer conn.Close()

				opened := map[string]struct{}{}
				unconfirmed := map[string]struct{}{}
				confirmed := map[string]struct{}{}
				for len(opened) < 30 {
					syncResp, err := stream.Recv()
					if err != nil {
						if err == io.EOF {
							break
						}
						log.Printf("opened = %d; unconfirmed = %d; confirmed = %d", len(opened), len(unconfirmed), len(confirmed))
						Expect(err).ShouldNot(HaveOccurred())
					}
					switch syncResp.Entry.OrderStatus {
					case relayer.OrderStatus_Open:
						opened[string(syncResp.Entry.Order.OrderId)] = struct{}{}
					case relayer.OrderStatus_Unconfirmed:
						unconfirmed[string(syncResp.Entry.Order.OrderId)] = struct{}{}
					case relayer.OrderStatus_Confirmed:
						confirmed[string(syncResp.Entry.Order.OrderId)] = struct{}{}
					default:
					}
				}
				Expect(len(opened)).Should(Equal(30))
				Expect(len(unconfirmed)).Should(Equal(0))
				Expect(len(confirmed)).Should(Equal(0))
			}()

			orders, err := CreateOrders(30, true)
			Expect(err).ShouldNot(HaveOccurred())
			env.SendOrders(orders)
		}, 30)

		FIt("should confirm orders that are compatible", func(done Done) {
			stream, conn, cancel, err := createTestRelayClient()
			Expect(err).ShouldNot(HaveOccurred())

			go func() {
				defer GinkgoRecover()
				defer close(done)
				defer cancel()
				defer conn.Close()

				opened := map[string]struct{}{}
				unconfirmed := map[string]struct{}{}
				confirmed := map[string]struct{}{}
				for len(confirmed) < 200 {
					syncResp, err := stream.Recv()
					if err != nil {
						if err == io.EOF {
							break
						}
						log.Printf("opened = %d; unconfirmed = %d; confirmed = %d", len(opened), len(unconfirmed), len(confirmed))
						Expect(err).ShouldNot(HaveOccurred())
					}
					switch syncResp.Entry.OrderStatus {
					case relayer.OrderStatus_Open:
						opened[string(syncResp.Entry.Order.OrderId)] = struct{}{}
					case relayer.OrderStatus_Unconfirmed:
						unconfirmed[string(syncResp.Entry.Order.OrderId)] = struct{}{}
					case relayer.OrderStatus_Confirmed:
						confirmed[string(syncResp.Entry.Order.OrderId)] = struct{}{}
					default:
					}
				}
				Expect(len(opened)).Should(Equal(200))
				Expect(len(unconfirmed)).Should(Equal(200))
				Expect(len(confirmed)).Should(Equal(200))
			}()

			err = env.SendMatchingOrderPairs(100) // Send 30 buys and 30 matching sells
			Expect(err).ShouldNot(HaveOccurred())

		}, 120 /* 2 minute timeout */)

		It("should release orders that conflict with a confirmed order", func(done Done) {
			stream, conn, cancel, err := createTestRelayClient()
			Expect(err).ShouldNot(HaveOccurred())

			go func() {
				defer GinkgoRecover()
				defer close(done)
				defer cancel()
				defer conn.Close()

				opened := map[string]struct{}{}
				confirmed := map[string]struct{}{}
				for len(confirmed) < 10 {
					syncResp, err := stream.Recv()
					if err != nil {
						if err == io.EOF {
							break
						}
						log.Printf("opened = %d; confirmed = %d", len(opened), len(confirmed))
						Expect(err).ShouldNot(HaveOccurred())
					}
					switch syncResp.Entry.OrderStatus {
					case relayer.OrderStatus_Open:
						opened[string(syncResp.Entry.Order.OrderId)] = struct{}{}
					case relayer.OrderStatus_Confirmed:
						confirmed[string(syncResp.Entry.Order.OrderId)] = struct{}{}
					default:
					}
				}

				// Wait some time, in case the Darknodes try to find more
				// confirmations
				time.Sleep(30 * time.Second)

				Expect(len(opened)).Should(Equal(30))
				Expect(len(confirmed)).Should(Equal(10))
			}()

			// Send 10 buy orders
			buyOrders, err := CreateOrders(10, true)
			Expect(err).ShouldNot(HaveOccurred())
			env.SendOrders(buyOrders)

			// Send 20 sell orders
			sellOrders, err := CreateOrders(20, false)
			Expect(err).ShouldNot(HaveOccurred())
			env.SendOrders(sellOrders)

		}, 120 /* 1 minute timeout */)

		It("should reject orders from unregistered addresses", func() {
			// FIXME: Implement
		})
	})

	Context("when synchronizing the orderbook", func() {

		// Serialize tests to prevent bleeding the orderbook.
		testMu := new(sync.Mutex)
		BeforeEach(func() {
			testMu.Lock()
			env.ClearOrderbooks()
		})
		AfterEach(func() {
			testMu.Unlock()
		})

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

			err = env.SendMatchingOrderPairs(20)
			Expect(err).ShouldNot(HaveOccurred())

			time.Sleep(10 * time.Second)
			Expect(isDeadlocked).To(Equal(false))

		})

		It("should not deadlock when the sync starts during the updates", func() {
			err := env.SendMatchingOrderPairs(5)
			Expect(err).ShouldNot(HaveOccurred())
			stream1, conn1, cancel1, err := createTestRelayClient()
			defer conn1.Close()
			defer cancel1()
			Expect(err).ShouldNot(HaveOccurred())
			err = env.SendMatchingOrderPairs(10)
			Expect(err).ShouldNot(HaveOccurred())
			stream2, conn2, cancel2, err := createTestRelayClient()
			defer conn2.Close()
			defer cancel2()
			Expect(err).ShouldNot(HaveOccurred())
			err = env.SendMatchingOrderPairs(5)
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
			err := env.SendMatchingOrderPairs(20)
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
})

func createTestRelayClient() (relayer.Relay_SyncClient, *client.Conn, context.CancelFunc, error) {
	// Create a relayer client to sync with the orderbook
	crypter := crypto.NewWeakCrypter()

	n := mathRnd.Intn(len(env.Darknodes))
	conn, err := client.Dial(context.Background(), env.Darknodes[n].MultiAddress())
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
		return nil, conn, cancel, err
	}
	return stream, conn, cancel, nil
}
