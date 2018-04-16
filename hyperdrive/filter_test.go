package hyper_test

import (
	"sync"

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

})
