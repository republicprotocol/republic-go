package hyper_test

import (
	"context"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/hyperdrive"
)

var _ = Describe("Channel set", func() {

	Context("when closing", func() {

		It("should panic when closed more than once", func() {
			chSet := hyper.NewChannelSet(0)
			chSet.Close()
			Ω(func() { chSet.Close }).Should(Panic())
		})

		It("should shutdown gracefully", func() {
			chSet := hyper.NewChannelSet(0)
			defer chSet.Close()
		})
	
	})

	Context("when splitting", func() {

		It("should split all messages to all outputs", func() {
			inputN := 100

			chSet := hyper.NewChannelSet(0)
			chSetsOut := []ChannelSet{
				hyper.NewChannelSet(0),
				hyper.NewChannelSet(0),
				hyper.NewChannelSet(0),
			}

			var writeWg sync.WaitGroup
			writeToChannelSet(&chSet, inputN, &writeWg)

			var splitWg sync.WaitGroup
			splitWg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer splitWg.Done()
				chSet.Split(chSetsOut)
			}()			

			var wg sync.WaitGroup
			for chSetOut := range chSetsOut {
				wg.Add(1)
				go func(chSetOut ChannelSet) {
					defer GinkgoRecover()
					defer wg.Done()
					n := 0
					for range chSetOut.Proposals {
						n++
					}
					Ω(n).Should(Equal(inputN))
				}(chSetOut)

				wg.Add(1)
				go func(chSetOut ChannelSet) {
					defer GinkgoRecover()
					defer wg.Done()
					n := 0
					for range chSetOut.Prepares {
						n++
					}
					Ω(n).Should(Equal(inputN))
				}(chSetOut)

				wg.Add(1)
				go func(chSetOut ChannelSet) {
					defer GinkgoRecover()
					defer wg.Done()
					n := 0
					for range chSetOut.Commits {
						n++
					}
					Ω(n).Should(Equal(inputN))
				}(chSetOut)

				wg.Add(1)
				go func(chSetOut ChannelSet) {
					defer GinkgoRecover()
					defer wg.Done()
					n := 0
					for range chSetOut.Blocks {
						n++
					}
					Ω(n).Should(Equal(inputN))
				}(chSetOut)

				wg.Add(1)
				go func(chSetOut ChannelSet) {
					defer GinkgoRecover()
					defer wg.Done()
					n := 0
					for range chSetOut.Faults {
						n++
					}
					Ω(n).Should(Equal(inputN))
				}(chSetOut)
			}

			writeWg.Wait()
			chSet.Close()

			splitWg.Wait()
			for chSetOut := range chSetsOut {
				chSetOut.Close()
			}

			wg.Wait()
		})

		It("should shutdown gracefully", func() {

			chSet := hyper.NewChannelSet(0)
			chSetsOut := []ChannelSet{
				hyper.NewChannelSet(0),
				hyper.NewChannelSet(0),
				hyper.NewChannelSet(0),
			}

			var writeWg sync.WaitGroup
			writeToChannelSet(&chSet, 100, &writeWg)

			var splitWg sync.WaitGroup
			splitWg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer splitWg.Done()
				chSet.Split(chSetsOut)
			}()

			writeWg.Wait()
			chSet.Close()

			splitWg.Wait()
			for chSetOut := range chSetsOut {
				chSetOut.Close()
			}
		})
	})

	Context("when piping", func() {

		It("should pipe all messages to the output", func() {
			inputN := 100

			chSet := hyper.NewChannelSet(0)
			chSetOut := hyper.NewChannelSet(0)

			var writeWg sync.WaitGroup
			writeToChannelSet(&chSet, inputN, &writeWg)
			
			var pipeWg sync.WaitGroup
			pipeWg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer pipeWg.Done()
				chSet.Pipe(chSetOut)
			}()

			var wg sync.WaitGroup
			for chSetOut := range chSetsOut {
				wg.Add(1)
				go func(chSetOut ChannelSet) {
					defer GinkgoRecover()
					defer wg.Done()
					n := 0
					for range chSetOut.Proposals {
						n++
					}
					Ω(n).Should(Equal(inputN))
				}(chSetOut)

				wg.Add(1)
				go func(chSetOut ChannelSet) {
					defer GinkgoRecover()
					defer wg.Done()
					n := 0
					for range chSetOut.Prepares {
						n++
					}
					Ω(n).Should(Equal(inputN))
				}(chSetOut)

				wg.Add(1)
				go func(chSetOut ChannelSet) {
					defer GinkgoRecover()
					defer wg.Done()
					n := 0
					for range chSetOut.Commits {
						n++
					}
					Ω(n).Should(Equal(inputN))
				}(chSetOut)

				wg.Add(1)
				go func(chSetOut ChannelSet) {
					defer GinkgoRecover()
					defer wg.Done()
					n := 0
					for range chSetOut.Blocks {
						n++
					}
					Ω(n).Should(Equal(inputN))
				}(chSetOut)

				wg.Add(1)
				go func(chSetOut ChannelSet) {
					defer GinkgoRecover()
					defer wg.Done()
					n := 0
					for range chSetOut.Faults {
						n++
					}
					Ω(n).Should(Equal(inputN))
				}(chSetOut)
			}

			writeWg.Wait()
			chSet.Close()

			pipeWg.Wait()
			chSetOut.Close()

			wg.Wait()
		})

		It("should shutdown gracefully", func() {
			chSet := hyper.NewChannelSet(0)
			chSetOut := hyper.NewChannelSet(0)

			var writeWg sync.WaitGroup
			writeToChannelSet(&chSet, 100, &writeWg)

			var pipeWg sync.WaitGroup
			pipeWg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer pipeWg.Done()
				chSet.Pipe(chSetOut)
			}()

			writeWg.Wait()
			chSet.Close()

			pipeWg.Wait()
			chSetOut.Close()
		})

	})

})

// writeToChannelSet writes a number of messages to all channels in the
// ChannelSet in background goroutines. The goroutine are added to the
// sync.WaitGroup but no waiting is done.
func writeToChannelSet(chSet *ChannelSet, n int, wg *sync.WaitGroup) {
	wg.Add(5)
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		for i := 0; i < n; i++ {
			chSet.Proposals <- hyper.Proposal{}
		}
	}()
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		for i := 0; i < n; i++ {
			chSet.Prepares <- hyper.Prepare{}
		}
	}()
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		for i := 0; i < n; i++ {
			chSet.Commits <- hyper.Commit{}
		}
	}()
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		for i := 0; i < n; i++ {
			chSet.Blocks <- hyper.Blocks{}
		}
	}()
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		for i := 0; i < n; i++ {
			chSet.Faults <- hyper.Fault{}
		}
	}()
}

		It("should be able to create, write, read and close a channel set from channels", func() {
			ctx, cancel := context.WithCancel(context.Background())
			chanSet := EmptyChannelSet(ctx, 100)
			chanSet.Proposal <- Proposal{
				Height: 1,
			}
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				select {
				case proposal, ok := <-chanSet.Proposal:
					if !ok {
						return
					}
					Ω(uint64(1)).Should(Equal(proposal.Height))
					wg.Done()
				}
			}()
			time.Sleep(10 * time.Microsecond)
			cancel()
			wg.Wait()
		})

		It("should be able to broadcast a channel set to multiple channel sets", func() {
			ctx, cancel := context.WithCancel(context.Background())
			chanSet := EmptyChannelSet(ctx, 100)
			chanSet.Proposal <- Proposal{
				Height: 1,
			}
			var wg sync.WaitGroup
			chanSets := make([]ChannelSet, 100)

			for i := 0; i < 100; i++ {
				chanSets[i] = EmptyChannelSet(ctx, 100)
			}

			go chanSet.Split(chanSets)

			wg.Add(100)
			for i := 0; i < 100; i++ {
				go func(i int) {
					for {
						select {
						case proposal, ok := <-chanSets[i].Proposal:
							if !ok {
								return
							}
							Ω(uint64(1)).Should(Equal(proposal.Height))
							wg.Done()
							return
						}
					}
				}(i)
			}

			wg.Wait()
			cancel()
		})

		
	})
})
