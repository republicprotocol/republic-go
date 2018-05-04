package smpcer_test

import (
	"context"
	"net"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/rpc/smpcer"

	"github.com/republicprotocol/republic-go/rpc"
	"google.golang.org/grpc"
)

var _ = Describe("Smpcer", func() {

	Context("NewSmpcer method", func() {

		It("should return a Smpcer object", func() {
			client, _, err := createNewClient("127.0.0.1", "3000")
			Expect(err).ShouldNot(HaveOccurred())
			rpc := new(rpc.RPC)
			Expect(NewSmpcer(&client, rpc)).ShouldNot(BeNil())
		})
	})

	Context("Smpc method", func() {

		PIt("should create two smpc clients and attempt to connect", func() {
			client1, multiaddr1, err := createNewClient("127.0.0.1", "3001")
			Expect(err).ShouldNot(HaveOccurred())
			rpc := new(rpc.RPC)
			smpc1 := NewSmpcer(&client1, rpc)

			client2, multiaddr2, err := createNewClient("127.0.0.1", "3002")
			Expect(err).ShouldNot(HaveOccurred())
			smpc2 := NewSmpcer(&client2, rpc)

			// Create sender channels
			ch1 := make(chan interface{})
			defer close(ch1)
			ch2 := make(chan interface{})
			defer close(ch2)

			server1 := grpc.NewServer()
			server2 := grpc.NewServer()
			smpc1.Register(server1)
			smpc2.Register(server2)

			listener1, err := net.Listen("tcp", "127.0.0.1:3001")
			Expect(err).ShouldNot(HaveOccurred())
			go func() {
				defer GinkgoRecover()
				err := server1.Serve(listener1)
				Expect(err).ShouldNot(HaveOccurred())
			}()

			listener2, err := net.Listen("tcp", "127.0.0.1:3002")
			Expect(err).ShouldNot(HaveOccurred())
			go func() {
				defer GinkgoRecover()
				err := server2.Serve(listener2)
				Expect(err).ShouldNot(HaveOccurred())
			}()

			messageCount, errCount := 0, 0
			var wg sync.WaitGroup
			wg.Add(3)
			for i := 0; i < 3; i++ {

				msgs1, errs1 := client1.Compute(context.Background(), multiaddr2, ch1)
				msgs2, errs2 := client2.Compute(context.Background(), multiaddr1, ch2)

				go func() {
					defer wg.Done()
					select {
					case _, ok := <-msgs1:
						if !ok {
							return
						}
						messageCount++
					case _, ok := <-msgs2:
						if !ok {
							return
						}
						messageCount++
					case val, ok := <-errs1:
						if !ok {
							return
						}
						_, ok = val.(error)
						Expect(ok).To(BeTrue())
						errCount++
					case val, ok := <-errs2:
						if !ok {
							return
						}
						_, ok = val.(error)
						Expect(ok).To(BeTrue())
						errCount++
					}
				}()
			}
			wg.Wait()

			Expect(messageCount).Should(Equal(0))
			Expect(errCount).Should(Equal(0))
		})
	})
})
