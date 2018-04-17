package hyperdrive_test

import (
	"fmt"
	"sync"
	"sync/atomic"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/hyperdrive"
)

var _ = Describe("Channel set", func() {

	Context("when closing", func() {
		It("should shutdown gracefully", func(done Done) {
			defer close(done)

			chSet := NewChannelSet(0)
			defer chSet.Close()
		})

		It("should panic when closed more than once", func(done Done) {
			defer close(done)

			chSet := NewChannelSet(0)
			chSet.Close()
			Ω(func() { chSet.Close() }).Should(Panic())
		})
	})

	Context("when splitting", func() {
		It("should shutdown gracefully", func(done Done) {
			defer close(done)

			numberOfMessages := 100
			chSet := NewChannelSet(300)
			chSetsOut := []ChannelSet{
				NewChannelSet(0),
				NewChannelSet(0),
				NewChannelSet(0),
			}

			var writeWg sync.WaitGroup
			writeToChannelSet(chSet, numberOfMessages, &writeWg)

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
				readFromChannelSet(chSetOut, numberOfMessages, &readWg, &n)
			}

			writeWg.Wait()
			chSet.Close()

			splitWg.Wait()
			for _, chSetOut := range chSetsOut {
				chSetOut.Close()
			}

			readWg.Wait()
		})

		It("should split all messages to all outputs", func(done Done) {
			defer close(done)

			numberOfMessages := 100
			chSet := NewChannelSet(0)
			chSetsOut := []ChannelSet{
				NewChannelSet(0),
				NewChannelSet(0),
				NewChannelSet(0),
			}

			var writeWg sync.WaitGroup
			writeToChannelSet(chSet, numberOfMessages, &writeWg)

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
				readFromChannelSet(chSetOut, numberOfMessages, &readWg, &n)
			}

			writeWg.Wait()
			chSet.Close()

			splitWg.Wait()
			for _, chSetOut := range chSetsOut {
				chSetOut.Close()
			}

			readWg.Wait()
			Ω(n).Should(Equal(int64(numberOfMessages * len(chSetsOut) * 5)))
		})

	})

	Context("when piping", func() {

		It("should shutdown gracefully", func(done Done) {
			defer close(done)

			numberOfMessages := 100
			chSet := NewChannelSet(0)
			chSetOut := NewChannelSet(0)

			var writeWg sync.WaitGroup
			writeToChannelSet(chSet, numberOfMessages, &writeWg)

			var pipeWg sync.WaitGroup
			pipeWg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer pipeWg.Done()
				chSet.Pipe(chSetOut)
			}()

			var n int64
			var readWg sync.WaitGroup
			readFromChannelSet(chSetOut, numberOfMessages, &readWg, &n)

			writeWg.Wait()
			chSet.Close()

			pipeWg.Wait()
			chSetOut.Close()

			readWg.Wait()
		})

		It("should pipe all messages to the output", func(done Done) {
			defer close(done)

			numberOfMessages := 100
			chSet := NewChannelSet(0)
			chSetOut := NewChannelSet(0)

			var writeWg sync.WaitGroup
			writeToChannelSet(chSet, numberOfMessages, &writeWg)

			var pipeWg sync.WaitGroup
			pipeWg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer pipeWg.Done()
				chSet.Pipe(chSetOut)
			}()

			var n int64
			var readWg sync.WaitGroup
			readFromChannelSet(chSetOut, numberOfMessages, &readWg, &n)

			writeWg.Wait()
			chSet.Close()

			pipeWg.Wait()
			chSetOut.Close()

			readWg.Wait()
			Ω(n).Should(Equal(int64(numberOfMessages * 5)))
		})
	})
})

// writeToChannelSet writes a number of messages to all channels in the
// ChannelSet in background goroutines. The goroutine are added to the
// sync.WaitGroup but no waiting is done.
func writeToChannelSet(chSet ChannelSet, numberOfMessages int, wg *sync.WaitGroup) {
	writeToChannelSetWithHeight(chSet, numberOfMessages, 0, wg)
}

func writeToChannelSetWithHeight(chSet ChannelSet, numberOfMessages int, height Height, wg *sync.WaitGroup) {
	wg.Add(5)
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		for i := 0; i < numberOfMessages; i++ {
			chSet.Proposals <- Proposal{Block: Block{Height: height}}
		}
	}()
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		for i := 0; i < numberOfMessages; i++ {
			chSet.Prepares <- Prepare{
				Proposal: Proposal{
					Block:     Block{Height: height},
					Signature: Signature{},
				},
			}
		}
	}()
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		for i := 0; i < numberOfMessages; i++ {
			chSet.Commits <- Commit{
				Prepare: Prepare{
					Proposal: Proposal{
						Block:     Block{Height: height},
						Signature: Signature{},
					},
				},
			}
		}
	}()
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		for i := 0; i < numberOfMessages; i++ {
			chSet.Blocks <- Block{Height: height}
		}
	}()
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		for i := 0; i < numberOfMessages; i++ {
			chSet.Faults <- Fault{Height: height}
		}
	}()
}

// readFromChannelSet reads all messages from all channels in the ChannelSet in
// background goroutines. The goroutine are added to the sync.WaitGroup but no
// waiting is done. Writes the total number of messages read to the int64
// pointer using atomic increments.
func readFromChannelSet(chSet ChannelSet, numberOfMessages int, wg *sync.WaitGroup, n *int64) {
	readFromChannelSetForHeight(chSet, numberOfMessages, wg, n, nil, nil)
}

func readFromChannelSetForHeight(chSet ChannelSet, numberOfMessages int, wg *sync.WaitGroup, n *int64, height *Height, errCh chan<- error) {
	wg.Add(1)
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		i := 0
		if i == numberOfMessages {
			return
		}
		for proposal := range chSet.Proposals {
			if height != nil {
				if proposal.Height != *height {
					errCh <- fmt.Errorf("unexpected height: expected %v got %v", *height, proposal.Height)
					continue
				}
			}
			atomic.AddInt64(n, 1)
			if i++; i == numberOfMessages {
				break
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		i := 0
		if i == numberOfMessages {
			return
		}
		for prepare := range chSet.Prepares {
			if height != nil {
				if prepare.Height != *height {
					errCh <- fmt.Errorf("unexpected height: expected %v got %v", *height, prepare.Height)
					continue
				}
			}
			atomic.AddInt64(n, 1)
			if i++; i == numberOfMessages {
				break
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		i := 0
		if i == numberOfMessages {
			return
		}
		for commit := range chSet.Commits {
			if height != nil {
				if commit.Height != *height {
					errCh <- fmt.Errorf("unexpected height: expected %v got %v", *height, commit.Height)
					continue
				}
			}
			atomic.AddInt64(n, 1)
			if i++; i == numberOfMessages {
				break
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		i := 0
		if i == numberOfMessages {
			return
		}
		for block := range chSet.Blocks {
			if height != nil {
				if block.Height != *height {
					errCh <- fmt.Errorf("unexpected height: expected %v got %v", *height, block.Height)
					continue
				}
			}
			atomic.AddInt64(n, 1)
			if i++; i == numberOfMessages {
				break
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer GinkgoRecover()
		defer wg.Done()
		i := 0
		if i == numberOfMessages {
			return
		}
		for fault := range chSet.Faults {
			if height != nil {
				if fault.Height != *height {
					errCh <- fmt.Errorf("unexpected height: expected %v got %v", *height, fault.Height)
					continue
				}
			}
			atomic.AddInt64(n, 1)
			if i++; i == numberOfMessages {
				break
			}
		}
	}()
}
