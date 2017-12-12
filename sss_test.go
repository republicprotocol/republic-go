package sss_test

import (
	"math/big"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/go-sss"
)

var _ = Describe("Shamir's secret sharing", func() {

	Context("configuration", func() {
		It("should be able to represent primes larger than 1024 bits", func() {
			// The maximum 1024 bit integer.
			max, ok := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137215", 10)
			Ω(ok).Should(Equal(true))
			max.Sub(max, big.NewInt(1))
			// The first prime above 1024 bits.
			prime, ok := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
			Ω(ok).Should(Equal(true))
			Ω(prime.Cmp(max) > 0).Should(Equal(true))
		})
	})

	Context("encoding", func() {
		It("should return the correct number of shares", func() {
			// Shamir parameters.
			n := int64(100)
			k := int64(50)
			prime, ok := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
			Ω(ok).Should(Equal(true))
			secret := big.NewInt(1234)
			// Encode the secret.
			shares, err := NewShamir(n, k, prime).Encode(secret)
			Ω(err).Should(BeNil())
			Ω(int64(len(shares))).Should(Equal(n))
		})
	})

	Context("decoding", func() {
		It("should return the correct secret from K shares", func() {
			// Shamir parameters.
			N := int64(100)
			K := int64(50)
			prime, ok := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
			Ω(ok).Should(Equal(true))
			secret := big.NewInt(1234)
			// Encode the secret.
			shares, err := NewShamir(N, K, prime).Encode(secret)
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
				decodedSecret, err := NewShamir(N, K, prime).Decode(kShares)
				Ω(err).Should(BeNil())
				Ω(decodedSecret.Cmp(secret)).Should(Equal(0))
			}
		})

		It("should return an incorrect secret from less than K shares", func() {
			// Shamir parameters.
			N := int64(100)
			K := int64(50)
			prime, ok := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
			Ω(ok).Should(Equal(true))
			secret := big.NewInt(1234)
			// Encode the secret.
			shares, err := NewShamir(N, K, prime).Encode(secret)
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
				decodedSecret, err := NewShamir(N, K, prime).Decode(kShares)
				Ω(err).Should(BeNil())
				Ω(decodedSecret.Cmp(secret)).ShouldNot(Equal(0))
			}
		})
	})
})
