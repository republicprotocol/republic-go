package smpc_test

import (
	"fmt"
	"log"
	"math/big"
	"runtime"
	"time"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/smpc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/dispatch"
	. "github.com/republicprotocol/republic-go/smpc"
)

var _ = Describe("Smpc workers", func() {

	Context("when receiving order fragment tasks", func() {

		It("it should produce all matching deltas correctly", func() {
			runtime.GOMAXPROCS(runtime.NumCPU())

			queueLimit := 1
			numOrders := 1
			n := 72
			k := 48
			prime, ok := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
			Ω(ok).Should(BeTrue())

			By("configuring a secure multiparty computation")

			deltaQueue := NewDeltaQueue(queueLimit)
			messageQueues := make(dispatch.MessageQueues, n)
			for i := 0; i < n; i++ {
				messageQueues[i] = dispatch.NewChannelQueue(queueLimit, false)
			}
			multiplexers := make([]dispatch.Multiplexer, n)
			for i := 0; i < n; i++ {
				multiplexers[i] = dispatch.NewMultiplexer(queueLimit)
			}
			workers := make(Workers, n)
			broadcasters := make(Broadcasters, n)
			for i := 0; i < n; i++ {
				go func(i int) {
					// defer GinkgoRecover()

					deltaFragmentMatrix := smpc.NewDeltaFragmentMatrix(prime)
					deltaBuilder := smpc.NewDeltaBuilder(int64(k), prime)

					// Run the MessageQueue that will feed data into the Multiplexer
					err := multiplexers[i].RunMessageQueue(fmt.Sprintf("%d", i), messageQueues[i])
					Ω(err).ShouldNot(HaveOccurred())

					// Create a Worker that is connected to all other parties
					workerPeerQueues := make(dispatch.MessageQueues, 0, n-1)
					for j := 0; j < n; j++ {
						if i == j {
							continue
						}
						workerPeerQueues = append(workerPeerQueues, messageQueues[j])
					}

					broadcasters[i] = NewBroadcaster(nil, workerPeerQueues, deltaFragmentMatrix, deltaBuilder, &deltaQueue)
					workers[i] = NewWorker(nil, workerPeerQueues, &multiplexers[i], deltaFragmentMatrix, deltaBuilder, &deltaQueue)

					go broadcasters[i].Run()
					workers[i].Run()
				}(i)
			}

			By("sending order fragments")
			go func() {
				for i := 0; i < numOrders; i++ {
					go func(i int) {
						buyOrderFragments, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(int64(i))).Split(int64(n), int64(k), prime)
						Ω(err).ShouldNot(HaveOccurred())
						for j := range multiplexers {
							multiplexers[j].Send(Message{
								OrderFragment: buyOrderFragments[j],
							})
						}
					}(i)
					go func(i int) {
						sellOrderFragments, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(int64(i))).Split(int64(n), int64(k), prime)
						Ω(err).ShouldNot(HaveOccurred())
						for j := range multiplexers {
							multiplexers[j].Send(Message{
								OrderFragment: sellOrderFragments[j],
							})
						}
					}(i)
				}
			}()

			By("waiting for deltas")
			for i := 0; i < n*numOrders*numOrders; i++ {
				delta, ok := deltaQueue.Recv()
				Ω(ok).Should(BeTrue())

				if (i+1)%(n*numOrders) == 0 {
					log.Println("found", i+1, "using", runtime.NumGoroutine())
				}

				switch delta := delta.(type) {
				case Delta:
					Ω(delta.IsMatch(prime)).Should(BeTrue())
				default:
					Fail(fmt.Sprintf("unrecognized type %T: expected smpc.Delta", delta))
				}
			}

			By("shutting down")
			for i := range multiplexers {
				multiplexers[i].Shutdown()
			}
			for i := range workers {
				workers[i].Shutdown()
			}
		})
	})

})
