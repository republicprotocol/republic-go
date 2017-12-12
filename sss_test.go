package sss_test

import (
	"log"
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/go-sss"
)

var _ = Describe("Shamir's secret sharing", func() {

	Context("configuration", func() {

		It("should be able to represent primes larger than 1024 bits", func() {
			// The maximum 1024 bit integer.
			max := big.NewInt(2)
			for i := 0; i < 10; i++ {
				max.Mul(max, max)
			}
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
			n := 6
			k := 3
			prime, ok := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
			Ω(ok).Should(Equal(true))
			secret := big.NewInt(1234)
			// Encode the secret.
			shares, err := NewShamir(n, k, prime).Encode(secret)
			Ω(err).Should(BeNil())
			Ω(len(shares)).Should(Equal(n))
		})

		It("should return a FiniteFieldError when the secret is too large", func() {
			// Shamir parameters.
			n := 6
			k := 3
			prime, ok := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
			Ω(ok).Should(Equal(true))
			secret := big.NewInt(0).Add(prime, big.NewInt(1))
			// Encode the secret.
			_, err := NewShamir(n, k, prime).Encode(secret)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(Equal(NewFiniteFieldError(secret)))
		})
	})

	Context("decoding", func() {

		It("should return the correct secret from all shares", func() {
			// Shamir parameters.
			n := 6
			k := 3
			prime, ok := big.NewInt(0).SetString("1237", 10)
			Ω(ok).Should(Equal(true))
			secret := big.NewInt(1234)
			// Encode the secret.
			shares, err := NewShamir(n, k, prime).Encode(secret)
			Ω(err).Should(BeNil())
			Ω(len(shares)).Should(Equal(n))
			// Decode the secret.
			decodedSecret, err := NewShamir(n, k, prime).Decode(shares)
			Ω(err).Should(BeNil())
			Ω(decodedSecret.Cmp(secret)).Should(Equal(0))
		})

		It("should return the correct secret from K shares", func() {
			// Shamir parameters.
			n := 6
			k := 3
			prime, ok := big.NewInt(0).SetString("1237", 10)
			Ω(ok).Should(Equal(true))
			secret := big.NewInt(1234)
			// Encode the secret.
			shares, err := NewShamir(n, k, prime).Encode(secret)
			Ω(err).Should(BeNil())
			Ω(len(shares)).Should(Equal(n))
			log.Println(shares)
			// Decode the secret.
			kShares := Shares{shares[1], shares[3], shares[4]}
			decodedSecret, err := NewShamir(n, k, prime).Decode(kShares)
			Ω(err).Should(BeNil())
			Ω(decodedSecret.Cmp(secret)).Should(Equal(0))
		})
	})
})
