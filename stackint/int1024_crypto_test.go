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
			Ω(FromBytes([]byte{01})).Should(Equal(FromUint(1)))
			Ω(FromBytes([]byte{0xFF, 0xFF})).Should(Equal(FromUint(65535)))

			// Test at least one case with bitlength % 8 = 0
			for i := 1; i <= (1<<8 + 1); i++ {
				max := FromUint(uint(i))

				r, err := Random(rand.Reader, &max)
				Ω(err).Should(BeNil())
				Ω(r.LessThan(&max) && r.GreaterThanOrEqual(&zero)).Should(BeTrue())
			}
		})

		It("should handle edge cases", func() {
			r, err := Random(rand.Reader, &zero)
			Ω(err).Should(BeNil())
			Ω(r).Should(Equal(zero))
		})
	})
})
