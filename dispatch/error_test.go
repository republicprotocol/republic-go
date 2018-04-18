package dispatch_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/dispatch"
)

var _ = Describe("Error channels", func() {

	Context("Merge errors", func() {

		It("should merge multiple error channels", func() {
			errCh1 := make(chan error)
			errCh2 := make(chan error)
			errCh3 := make(chan error)

			err1 := errors.New("1")
			err2 := errors.New("2")
			err3 := errors.New("3")

			errCh1 <- err1
			errCh2 <- err2
			errCh3 <- err3

			errCh := MergeErrors(errCh1, errCh2, errCh3)
			Close(errCh1, errCh2, errCh3)

			Ω(len(errCh)).Should(Equal(3))
		})

		It("should be able to read the errors originated from all the error channels", func() {
			errCh1 := make(chan error)
			errCh2 := make(chan error)
			errCh3 := make(chan error)

			err1 := errors.New("1")
			err2 := errors.New("2")
			err3 := errors.New("3")

			errCh1 <- err1
			errCh2 <- err2
			errCh3 <- err3

			errCh := MergeErrors(errCh1, errCh2, errCh3)

			Ω(len(errCh)).Should(Equal(3))

			Close(errCh1, errCh2, errCh3)

			for err := range errCh {
				Ω(err == err1 || err == err2 || err == err3).Should(BeTrue())
			}

		})

	})

	Context("Filter errors", func() {
		It("", func() {

		})
	})

	Context("Consume errors", func() {
		It("", func() {

		})
	})
})
