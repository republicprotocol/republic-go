package dispatch_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/dispatch"
)

var _ = Describe("Coforall loops", func() {

	Context("when iterating over integers", func() {
		It("should apply the function to all elements", func(done Done) {
			defer close(done)

			xs := make(chan int, 10)
			CoForAll(10, func(i int) {
				xs <- (i + 1)
			})
			close(xs)

			for i := range xs {
				Expect(i).Should(BeNumerically(">=", 1))
				Expect(i).Should(BeNumerically("<=", 10))
			}
		})
	})

	Context("when iterating over arrays", func() {
		It("should apply the function to all elements", func() {
			xs := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
			CoForAll(xs, func(i int) {
				xs[i] *= 2
			})
			for i := range xs {
				Expect(xs[i]).Should(Equal(i * 2))
			}
		})
	})

	Context("when iterating over slices", func() {
		It("should apply the function to all elements", func() {
			xs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
			CoForAll(xs[:], func(i int) {
				xs[i] *= 2
			})
			for i := range xs {
				Expect(xs[i]).Should(Equal(i * 2))
			}
		})
	})

})
