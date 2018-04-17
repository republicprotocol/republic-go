package hyperdrive_test

import (
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/hyperdrive"
)

var _ = Describe("Filters", func() {

	Context("when filtering duplicates", func() {

		It("should shutdown gracefully", func() {
			chSet := NewChannelSet(0)
			chSetOut := FilterDuplicates(chSet, 0)

			var writeWg sync.WaitGroup
			writeToChannelSet(chSet, 100, &writeWg)

			var n int64
			var readWg sync.WaitGroup
			readFromChannelSet(chSetOut, &readWg, &n)

			writeWg.Wait()
			chSet.Close()

			readWg.Wait()
		})

		It("should never produce a duplicate", func() {
			chSet := NewChannelSet(0)
			chSetOut := FilterDuplicates(chSet, 0)

			var writeWg sync.WaitGroup
			writeToChannelSet(chSet, 100, &writeWg)

			var n int64
			var readWg sync.WaitGroup
			readFromChannelSet(chSetOut, &readWg, &n)

			writeWg.Wait()
			chSet.Close()

			readWg.Wait()
			Ω(n).Should(Equal(int64(5)))
		})

	})

	Context("when filtering heights", func() {

		FIt("should shutdown gracefully", func() {
			capacity := 0
			height := make(chan int, capacity)
			chSet := NewChannelSet(capacity)
			chSetOut := FilterHeight(chSet, height, capacity)

			var heightWg sync.WaitGroup
			heightWg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer heightWg.Done()

				for i := 0; i < 10; i++ {
					height <- i
					time.Sleep(2 * time.Second)
				}
			}()

			var writeWg sync.WaitGroup
			writeToChannelSet(chSet, 100, &writeWg)

			var n int64
			var readWg sync.WaitGroup
			readFromChannelSet(chSetOut, &readWg, &n)

			writeWg.Wait()
			chSet.Close()

			readWg.Wait()
			heightWg.Wait()
		})

		It("should only produce messages for the current height", func() {
			numberOfMessages := 100
			chSet := NewChannelSet(0)
			chSetOut := FilterDuplicates(chSet, 0)

			var writeWg sync.WaitGroup
			for h := 0; h < 5; h++ {
				writeToChannelSet(chSet, numberOfMessages, &writeWg)
			}

			var n int64
			var readWg sync.WaitGroup
			readFromChannelSet(chSetOut, &readWg, &n)

			writeWg.Wait()
			chSet.Close()

			readWg.Wait()
			Ω(n).Should(Equal(numberOfMessages))
		})

		It("should produce buffered messages when the height changes", func() {
			numberOfMessages, capacity := 100, 0
			chSet := NewChannelSet(0)
			chSetOut := FilterDuplicates(chSet, 0)
			height := make(chan int, capacity)

			var heightWg sync.WaitGroup
			heightWg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer heightWg.Done()

				for i := 0; i < 10; i++ {
					height <- i
					time.Sleep(2 * time.Second)
				}
			}()

			var writeWg sync.WaitGroup
			for height := 5; height > 0; height-- {
				writeToChannelSetWithHeight(chSet, numberOfMessages, height, &writeWg)
			}

			var n int64
			var readWg sync.WaitGroup
			readFromChannelSet(chSetOut, &readWg, &n)

			writeWg.Wait()
			chSet.Close()

			readWg.Wait()

			heightWg.Wait()

			Ω(n).Should(Equal(int64(500)))
		})
	})
})
