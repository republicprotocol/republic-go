package smpc_test

import (
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/smpc"

	"github.com/republicprotocol/republic-go/shamir"
)

var _ = Describe("share builder", func() {

	Context("insert shares", func() {

		It("should return the secret after atleast k shares have joined", func() {
			n := int64(24)
			k := int64(16)
			shareBuilder := NewShareBuilder(k)

			for i := 0; i < 100; i++ {
				secret := uint64(rand.Intn(100))
				shares, err := shamir.Split(n, k, secret)
				Expect(err).ShouldNot(HaveOccurred())
				for j := int64(0); j < n; j++ {
					val, err := shareBuilder.InsertShare([32]byte{byte(i)}, shares[j])
					if j < k-1 {
						Expect(err).Should(HaveOccurred())
						Expect(err).To(Equal(ErrInsufficientSharesToJoin))
						Expect(val).To(Equal(uint64(0)))
					} else {
						Expect(err).ShouldNot(HaveOccurred())
						Expect(val).To(Equal(secret))
					}
				}
			}
		})
	})
})
