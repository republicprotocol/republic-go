package rpc_test

import (
	"context"
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc/codes"

	"github.com/republicprotocol/republic-go/identity"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var _ = Describe("Secure multi-party computer service", func() {

	Context("when streaming computations", func() {

		It("should setup a rendezvous between clients and servers", func(done Done) {
			defer close(done)

			numberOfMessages := 1

			serverKeyPair, err := identity.NewKeyPair()
			Ω(err).ShouldNot(HaveOccurred())
			serverMulti, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/0.0.0.0/tcp/3000/republic/%s", serverKeyPair.Address().String()))
			Ω(err).ShouldNot(HaveOccurred())
			server := grpc.NewServer()

			service := NewComputerService()
			service.Register(server)

			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				listener, err := net.Listen("tcp", "0.0.0.0:3000")
				Ω(err).ShouldNot(HaveOccurred())

				err = server.Serve(listener)
				Ω(err).ShouldNot(HaveOccurred())
			}()

			clientKeyPair, err := identity.NewKeyPair()
			Ω(err).ShouldNot(HaveOccurred())
			clientMulti, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/0.0.0.0/tcp/3001/republic/%s", clientKeyPair.Address().String()))
			Ω(err).ShouldNot(HaveOccurred())
			client, err := NewClient(context.Background(), serverMulti, clientMulti)
			Ω(err).ShouldNot(HaveOccurred())

			clientCtx, clientCtxCancel := context.WithCancel(context.Background())
			clientChIn := make(chan *Computation)
			clientChOut, clientErrCh := client.Compute(clientCtx, clientChIn)

			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				for err := range clientErrCh {
					if s, _ := status.FromError(err); s != nil {
						Ω(s.Code()).Should(Equal(codes.Canceled))
					} else {
						Ω(err).Should(Equal(context.Canceled))
					}
				}
			}()

			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				serviceChIn := make(chan *Computation)
				serviceChOut, serviceErrCh := service.WaitForCompute(clientMulti, serviceChIn)

				var serviceWg sync.WaitGroup
				serviceWg.Add(1)
				go func() {
					defer GinkgoRecover()
					defer wg.Done()

					for err := range serviceErrCh {
						if s, _ := status.FromError(err); s != nil {
							Ω(s.Code()).Should(Equal(codes.Canceled))
						} else {
							Ω(err).Should(Equal(context.Canceled))
						}
					}
				}()

				serviceWg.Add(1)
				go func() {
					defer GinkgoRecover()
					defer serviceWg.Done()

					for i := 0; i < numberOfMessages; i++ {
						serviceChIn <- &Computation{MultiAddress: MarshalMultiAddress(&serverMulti)}
					}
					close(serviceChIn)
				}()
				serviceWg.Add(1)
				go func() {
					defer GinkgoRecover()
					defer serviceWg.Done()

					for i := 0; i < numberOfMessages; i++ {
						<-serviceChOut
					}
				}()

				serviceWg.Wait()
			}()

			clientChIn <- &Computation{MultiAddress: MarshalMultiAddress(&clientMulti)}

			wg.Add(2)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				for i := 0; i < numberOfMessages; i++ {
					clientChIn <- &Computation{}
				}

				clientCtxCancel()
				close(clientChIn)
			}()
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				for i := 0; i < numberOfMessages; i++ {
					<-clientChOut
				}
			}()

			wg.Wait()
		})

	})

})

//import (
//	"context"
//	"fmt"
//	"log"
//	"math/big"
//	"net"
//	"runtime"
//	"time"
//
//	. "github.com/onsi/ginkgo"
//	. "github.com/onsi/gomega"
//	"github.com/republicprotocol/republic-go/dispatch"
//	"github.com/republicprotocol/republic-go/identity"
//	"github.com/republicprotocol/republic-go/order"
//	. "github.com/republicprotocol/republic-go/rpc"
//	"google.golang.org/grpc"
//)
//
//var _ = Describe("Smpc service", func() {
//
//	Context("when streaming messages", func() {
//
//		It("it should produce all matching deltas correctly", func() {
//			runtime.GOMAXPROCS(runtime.NumCPU())
//
//			queueLimit := 400
//			numOrders := 1
//			numWorkers := 1
//			n := 18
//			k := 12
//			prime, ok := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
//			Ω(ok).Should(BeTrue())
//
//			By("creating multiplexers for each service")
//
//			multiplexers := make([]dispatch.Multiplexer, n)
//			for i := 0; i < n; i++ {
//				multiplexers[i] = dispatch.NewMultiplexer(queueLimit)
//			}
//
//			By("creating identities for each service")
//
//			keyPairs := make([]identity.KeyPair, n)
//			multiAddresses := make([]identity.MultiAddress, n)
//			for i := 0; i < n; i++ {
//				keyPair, err := identity.NewKeyPair()
//				Ω(err).ShouldNot(HaveOccurred())
//				multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/republic/%s", 3000+i, keyPair.Address().String()))
//				Ω(err).ShouldNot(HaveOccurred())
//				keyPairs[i] = keyPair
//				multiAddresses[i] = multiAddress
//			}
//
//			By("running each service")
//
//			services := make([]SmpcService, n)
//			servers := make([]*grpc.Server, n)
//			for i := 0; i < n; i++ {
//				services[i] = NewSmpcService(&multiAddresses[i], &multiplexers[i], queueLimit)
//				servers[i] = grpc.NewServer()
//				services[i].Register(servers[i])
//				go func(i int) {
//					// defer GinkgoRecover()
//					listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 3000+i))
//					Ω(err).ShouldNot(HaveOccurred())
//					Ω(servers[i].Serve(listener)).ShouldNot(HaveOccurred())
//				}(i)
//			}
//			time.Sleep(time.Second)
//
//			By("configuring a secure multiparty computation")
//
//			deltaQueue := smpc.NewDeltaQueue(queueLimit)
//
//			workers := make(smpc.Workers, n*numWorkers)
//			for i := 0; i < n; i++ {
//				go func(i int) {
//					// defer GinkgoRecover()
//
//					// Create a Worker that is connected to all other parties
//					deltaFragmentMatrix := smpc.NewDeltaFragmentMatrix(prime)
//					deltaBuilder := smpc.NewDeltaBuilder(int64(k), prime)
//					workerPeerQueues := make(dispatch.MessageQueues, 0, n-1)
//					for j := 0; j < n; j++ {
//						if i == j {
//							continue
//						}
//
//						conn, err := grpc.DialContext(context.Background(), fmt.Sprintf("0.0.0.0:%d", 3000+j), grpc.WithInsecure())
//						Ω(err).ShouldNot(HaveOccurred())
//						client := NewSmpcClient(conn)
//						stream, err := client.Compute(context.Background(), grpc.FailFast(true))
//						if err == nil {
//							stream.Send(&SmpcMessage{
//								MultiAddress: MarshalMultiAddress(&multiAddresses[i]),
//							})
//							queue := NewSmpcClientStreamQueue(stream, queueLimit)
//							go func() {
//								log.Println(queue.Run())
//							}()
//							workerPeerQueues = append(workerPeerQueues, queue)
//						} else {
//							log.Println(err)
//						}
//					}
//					for j := 0; j < numWorkers; j++ {
//						go func(i, j int) {
//							workers[i*numWorkers+j] = smpc.NewWorker(nil, workerPeerQueues, &multiplexers[i], &deltaFragmentMatrix, &deltaBuilder, &deltaQueue)
//							workers[i*numWorkers+j].Run()
//						}(i, j)
//					}
//				}(i)
//			}
//			time.Sleep(time.Second)
//
//			By("sending order fragments")
//			go func() {
//				for i := 0; i < numOrders; i++ {
//					// go func(i int) {
//					buyOrderFragments, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(int64(i))).Split(int64(n), int64(k), prime)
//					Ω(err).ShouldNot(HaveOccurred())
//					for j := range multiplexers {
//						multiplexers[j].Send(smpc.Message{
//							OrderFragment: buyOrderFragments[j],
//						})
//					}
//					log.Println("SENT BUY ORDER", i)
//				}
//				// }(i)
//				// go func(i int) {
//				for i := 0; i < numOrders; i++ {
//					sellOrderFragments, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(int64(i))).Split(int64(n), int64(k), prime)
//					Ω(err).ShouldNot(HaveOccurred())
//					for j := range multiplexers {
//						multiplexers[j].Send(smpc.Message{
//							OrderFragment: sellOrderFragments[j],
//						})
//					}
//					log.Println("SENT SELL ORDER", i)
//					// }(i)
//				}
//			}()
//			time.Sleep(time.Second)
//
//			By("waiting for deltas")
//			for i := 0; i < n*numOrders*numOrders; i++ {
//				delta, ok := deltaQueue.Recv()
//				Ω(ok).Should(BeTrue())
//
//				// if (i+1)%(n*numOrders) == 0 {
//				log.Println("found", i+1, "[", runtime.NumGoroutine(), "]")
//				// }
//
//				switch delta := delta.(type) {
//				case smpc.Delta:
//					Ω(delta.IsMatch(prime)).Should(BeTrue())
//				default:
//					Fail(fmt.Sprintf("unrecognized type %T: expected smpc.Delta", delta))
//				}
//			}
//
//			By("shutting down")
//			for i := range multiplexers {
//				multiplexers[i].Shutdown()
//			}
//			for i := range workers {
//				workers[i].Shutdown()
//			}
//		})
//	})
//
//})
