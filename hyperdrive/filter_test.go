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
			numberOfMessages := 100
			chSet := NewChannelSet(0)
			chSetOut := FilterDuplicates(chSet, 0)

			var writeWg sync.WaitGroup
			writeToChannelSet(chSet, numberOfMessages, &writeWg)

			var n int64
			var readWg sync.WaitGroup
			readFromChannelSet(chSetOut, numberOfMessages, &readWg, &n)

			writeWg.Wait()
			chSet.Close()

			readWg.Wait()
		})

		It("should never produce a duplicate", func() {
			numberOfMessages := 100
			chSet := NewChannelSet(0)
			chSetOut := FilterDuplicates(chSet, 0)

			var writeWg sync.WaitGroup
			writeToChannelSet(chSet, numberOfMessages, &writeWg)

			var n int64
			var readWg sync.WaitGroup
			readFromChannelSet(chSetOut, numberOfMessages, &readWg, &n)

			writeWg.Wait()
			chSet.Close()

			readWg.Wait()
			Ω(n).Should(Equal(int64(5)))
		})

	})

	Context("when filtering heights", func() {

		It("should shutdown gracefully", func() {
			numberOfMessages := 100
			capacity := 0
			height := make(chan Height, capacity)
			chSet := NewChannelSet(capacity)
			chSetOut := FilterHeight(chSet, height, capacity)

			var heightWg sync.WaitGroup
			heightWg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer heightWg.Done()

				for i := 0; i < 5; i++ {
					height <- Height(i)
					time.Sleep(time.Second)
				}
			}()

			var writeWg sync.WaitGroup
			writeToChannelSet(chSet, numberOfMessages, &writeWg)

			var n int64
			var readWg sync.WaitGroup
			readFromChannelSet(chSetOut, numberOfMessages, &readWg, &n)

			writeWg.Wait()
			heightWg.Wait()

			chSet.Close()
			close(height)

			readWg.Wait()
		})

		It("should only produce messages for the current height", func() {
			numberOfMessages := 100
			capacity := 0
			height := make(chan Height, capacity)
			chSet := NewChannelSet(capacity)
			chSetOut := FilterHeight(chSet, height, numberOfMessages)

			h := Height(1)
			hErrCh := make(chan error)
			go func() {
				// defer GinkgoRecover()
				for err := range hErrCh {
					Ω(err).ShouldNot(HaveOccurred())
				}
			}()

			var heightWg sync.WaitGroup
			heightWg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer heightWg.Done()

				height <- 1
			}()

			var writeWg sync.WaitGroup
			writeToChannelSetWithHeight(chSet, numberOfMessages, Height(1), &writeWg)
			writeToChannelSetWithHeight(chSet, numberOfMessages, Height(2), &writeWg)
			writeToChannelSetWithHeight(chSet, numberOfMessages, Height(3), &writeWg)

			var n int64
			var readWg sync.WaitGroup
			readFromChannelSetForHeight(chSetOut, numberOfMessages, &readWg, &n, &h, hErrCh)

			writeWg.Wait()
			heightWg.Wait()

			chSet.Close()
			close(height)

			readWg.Wait()
			close(hErrCh)
		})

		It("should continue to produce messages when the height changes", func() {
			numberOfMessages := 100
			capacity := 0
			height := make(chan Height, capacity)
			chSet := NewChannelSet(capacity)
			chSetOut := FilterHeight(chSet, height, numberOfMessages)

			hErrCh := make(chan error)
			go func() {
				// defer GinkgoRecover()
				for err := range hErrCh {
					Ω(err).ShouldNot(HaveOccurred())
				}
			}()

			for i := 0; i < 10; i++ {
				h := Height(i)
				height <- h

				var writeWg sync.WaitGroup
				writeToChannelSetWithHeight(chSet, numberOfMessages, h, &writeWg)

				var n int64
				var readWg sync.WaitGroup
				readFromChannelSetForHeight(chSetOut, numberOfMessages, &readWg, &n, &h, hErrCh)

				writeWg.Wait()
				readWg.Wait()
			}

			chSet.Close()
			close(height)
			close(hErrCh)
		})

		It("should produce buffered messages when the height changes", func() {
			numberOfMessages := 100
			capacity := 0
			height := make(chan Height, capacity)
			chSet := NewChannelSet(capacity)
			chSetOut := FilterHeight(chSet, height, numberOfMessages)

			h1 := Height(1)
			hErrCh := make(chan error)
			go func() {
				// defer GinkgoRecover()
				for err := range hErrCh {
					Ω(err).ShouldNot(HaveOccurred())
				}
			}()

			height <- 1

			var writeWg sync.WaitGroup
			writeToChannelSetWithHeight(chSet, numberOfMessages, Height(1), &writeWg)
			writeToChannelSetWithHeight(chSet, numberOfMessages, Height(2), &writeWg)

			var n int64
			var readWg sync.WaitGroup
			readFromChannelSetForHeight(chSetOut, numberOfMessages, &readWg, &n, &h1, hErrCh)

			writeWg.Wait()
			readWg.Wait()

			height <- 2

			h2 := Height(2)
			readFromChannelSetForHeight(chSetOut, numberOfMessages, &readWg, &n, &h2, hErrCh)

			chSet.Close()
			close(height)

			readWg.Wait()
			close(hErrCh)
		})
	})
})
