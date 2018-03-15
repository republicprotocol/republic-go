package stackint_test

import (
	"crypto/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/stackint"
)

var _ = Describe("Int1024 crypto utilities", func() {
	Context("creating random Int1024s", func() {
		It("should work", func() {
			Ω(FromBytes([]byte{01})).Should(Equal(FromUint64(1)))
			Ω(FromBytes([]byte{0xFF, 0xFF})).Should(Equal(FromUint64(65535)))

			for i := 0; i < 10000; i++ {
				max := FromUint64(2)
				r, err := Random(rand.Reader, &max)
				Ω(err).Should(BeNil())
				Ω(r.LessThan(&max) && r.GreaterThanOrEqual(&ZERO)).Should(BeTrue())
			}

			for i := 0; i < 10000; i++ {
				max := FromUint64(10)
				r, err := Random(rand.Reader, &max)
				Ω(err).Should(BeNil())
				Ω(r.LessThan(&max) && r.GreaterThanOrEqual(&ZERO)).Should(BeTrue())
			}
		})
	})
})
