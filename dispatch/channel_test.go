package dispatch_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/dispatch"
)

var _ = Describe("Dispatch Package", func() {

	FContext("Splitter", func() {

		It("Is able to split the channel", func() {
			inCh := make(chan int)
			outChs := make([]chan int, 100)
			for i := 0; i < 100; i++ {
				outChs[i] = make(chan int)
			}

			go Split(inCh, outChs)

			go func() {
				defer close(inCh)
				inCh <- 1729
			}()

			for _, ch := range outChs {
				i := <-ch
				Ω(i).Should(Equal(1729))
			}
		})

	})
})
