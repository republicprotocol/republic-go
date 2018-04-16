package hyper_test

import (
	"sync"
	"sync/atomic"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/hyperdrive"
)

var _ = Describe("Channel set", func() {

	Context("when closing", func() {

		It("should panic when closed more than once", func() {
			chSet := NewChannelSet(0)
			chSet.Close()
			Ω(func() { chSet.Close() }).Should(Panic())
		})

		It("should shutdown gracefully", func() {
			chSet := NewChannelSet(0)
			defer chSet.Close()
		})

	})

	Context("when splitting", func() {

		It("should shutdown gracefully", func() {

			chSet := NewChannelSet(300)
			chSetsOut := []ChannelSet{
				NewChannelSet(0),
				NewChannelSet(0),
				NewChannelSet(0),
			}

			var writeWg sync.WaitGroup
			writeToChannelSet(chSet, 100, &writeWg)

			var splitWg sync.WaitGroup
			splitWg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer splitWg.Done()
				chSet.Split(chSetsOut...)
			}()

			var n int64
			var readWg sync.WaitGroup
			for _, chSetOut := range chSetsOut {
				readFromChannelSet(chSetOut, &readWg, &n)
			}

			writeWg.Wait()
			chSet.Close()

			splitWg.Wait()
			for _, chSetOut := range chSetsOut {
				chSetOut.Close()
			}

			readWg.Wait()
		})

		It("should split all messages to all outputs", func() {
			inputN := 100

			chSet := NewChannelSet(0)
			chSetsOut := []ChannelSet{
				NewChannelSet(0),
				NewChannelSet(0),
				NewChannelSet(0),
			}

			var writeWg sync.WaitGroup
			writeToChannelSet(chSet, inputN, &writeWg)

			var splitWg sync.WaitGroup
			splitWg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer splitWg.Done()
				chSet.Split(chSetsOut...)
			}()

			var n int64
			var readWg sync.WaitGroup
			for _, chSetOut := range chSetsOut {
				readFromChannelSet(chSetOut, &readWg, &n)
			}

			writeWg.Wait()
			chSet.Close()

			splitWg.Wait()
			for _, chSetOut := range chSetsOut {
				chSetOut.Close()
			}

			readWg.Wait()
			Ω(n).Should(Equal(int64(inputN * len(chSetsOut) * 5)))
		})

	})

	Context("when piping", func() {

		It("should shutdown gracefully", func() {
			chSet := NewChannelSet(0)
			chSetOut := NewChannelSet(0)

			var writeWg sync.WaitGroup
			writeToChannelSet(chSet, 100, &writeWg)

			var pipeWg sync.WaitGroup
			pipeWg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer pipeWg.Done()
				chSet.Pipe(chSetOut)
			}()

			var n int64
			var readWg sync.WaitGroup
			readFromChannelSet(chSetOut, &readWg, &n)

			writeWg.Wait()
			chSet.Close()

			pipeWg.Wait()
			chSetOut.Close()

			readWg.Wait()
		})

		It("should pipe all messages to the output", func() {
			inputN := 100

			chSet := NewChannelSet(0)
			chSetOut := NewChannelSet(0)

			var writeWg sync.WaitGroup
			writeToChannelSet(chSet, inputN, &writeWg)

			var pipeWg sync.WaitGroup
			pipeWg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer pipeWg.Done()
				chSet.Pipe(chSetOut)
			}()

			var n int64
			var readWg sync.WaitGroup
			readFromChannelSet(chSetOut, &readWg, &n)

			writeWg.Wait()
			chSet.Close()

			pipeWg.Wait()
			chSetOut.Close()

			readWg.Wait()
			Ω(n).Should(Equal(int64(inputN * 5)))
		})

	})

})

// writeToChannelSet writes a number of messages to all channels in the
// ChannelSet in background goroutines. The goroutine are added to the
// sync.WaitGroup but no waiting is done.
func writeToChannelSet(chSet ChannelSet, n int, wg *sync.WaitGroup) {
	writeToChannelSetWithHeight(chSet, n, 0, wg)
}

func writeToChannelSetWithHeight(chSet ChannelSet, n, height int, wg *sync.WaitGroup) {
	wg.Add(5)
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		for i := 0; i < n; i++ {
			chSet.Proposals <- Proposal{Height: height}
		}
	}()
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		for i := 0; i < n; i++ {
			chSet.Prepares <- Prepare{Height: height}
		}
	}()
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		for i := 0; i < n; i++ {
			chSet.Commits <- Commit{Height: height}
		}
	}()
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		for i := 0; i < n; i++ {
			chSet.Blocks <- Block{Height: height}
		}
	}()
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		for i := 0; i < n; i++ {
			chSet.Faults <- Fault{Height: height}
		}
	}()
}

// readFromChannelSet reads all messages from all channels in the ChannelSet in
// background goroutines. The goroutine are added to the sync.WaitGroup but no
// waiting is done. Writes the total number of messages read to the int64
// pointer using atomic increments.
func readFromChannelSet(chSet ChannelSet, wg *sync.WaitGroup, n *int64) {
	wg.Add(1)
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		for range chSet.Proposals {
			atomic.AddInt64(n, 1)
		}
	}()

	wg.Add(1)
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		for range chSet.Prepares {
			atomic.AddInt64(n, 1)
		}
	}()

	wg.Add(1)
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		for range chSet.Commits {
			atomic.AddInt64(n, 1)
		}
	}()

	wg.Add(1)
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		for range chSet.Blocks {
			atomic.AddInt64(n, 1)
		}
	}()

	wg.Add(1)
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		for range chSet.Faults {
			atomic.AddInt64(n, 1)
		}
	}()
}
