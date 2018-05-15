package shamir_test

import (
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/shamir"

	"github.com/republicprotocol/republic-go/stackint"
)

var _ = Describe("Shamir's secret sharing", func() {

	Context("serialization", func() {
		It("should be able to serialize and deserialize shares", func() {
			for i := uint64(0); i < 1000; i++ {
				bytes, err := Share{
					Index: i,
					Value: Prime - i,
				}.MarshalBinary()
				Expect(err).ShouldNot(HaveOccurred())
				share := Share{}
				err = share.UnmarshalBinary(bytes)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(share.Index).Should(Equal(i))
				Ω(share.Value).Should(Equal(Prime - i))
			}
		})
		It("should return an error when deserializing an empty byte slice", func() {
			share := Share{}
			err := share.UnmarshalBinary([]byte{})
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("splitting", func() {
		It("should return the correct number of shares", func() {
			// Shamir parameters.
			n := int64(100)
			k := int64(50)
			secret := uint64(1234)
			// Split the secret.
			shares, err := Split(n, k, secret)
			Ω(err).Should(BeNil())
			Ω(int64(len(shares))).Should(Equal(n))
		})
	})

	Context("joining", func() {
		It("should return the correct secret from K shares", func() {
			// Shamir parameters.
			N := int64(100)
			K := int64(50)
			secret := uint64(1234)
			// Split the secret.
			shares, err := Split(N, K, secret)
			Ω(err).Should(BeNil())
			Ω(int64(len(shares))).Should(Equal(N))
			// For all K greater than, or equal to, 50 attempt to decode the secret.
			for k := int64(50); k < 101; k++ {
				// Pick K unique indices in the range [0, k).
				indices := map[int]struct{}{}
				for i := 0; i < int(k); i++ {
					for {
						index := rand.Intn(int(k))
						if _, ok := indices[index]; !ok {
							indices[index] = struct{}{}
							break
						}
					}
				}
				// Use K shares to reconstruct the secret.
				kShares := make(Shares, k)
				for index := range indices {
					kShares[index] = shares[index]
				}
				decodedSecret := stackint.FromUint(uint(Join(kShares)))
				secretStackInt := stackint.FromUint(uint(secret))
				Ω(decodedSecret.Cmp(&secretStackInt)).Should(Equal(0))
			}
		})

		It("should return an incorrect secret from less than K shares", func() {
			// Shamir parameters.
			N := int64(100)
			K := int64(50)
			secret := uint64(1234)
			// Split the secret.
			shares, err := Split(N, K, secret)
			Ω(err).Should(BeNil())
			Ω(int64(len(shares))).Should(Equal(N))
			// For all K less than 50 attempt to decode the secret.
			for k := int64(1); k < 50; k++ {
				// Pick K unique indices in the range [0, k).
				indices := map[int]struct{}{}
				for i := 0; i < int(k); i++ {
					for {
						index := rand.Intn(int(k))
						if _, ok := indices[index]; !ok {
							indices[index] = struct{}{}
							break
						}
					}
				}
				// Use K shares to reconstruct the secret.
				kShares := make(Shares, k)
				for index := range indices {
					kShares[index] = shares[index]
				}
				decodedSecret := Join(kShares)
				Ω(decodedSecret).Should(Not(Equal(&secret)))
			}
		})
	})
})
