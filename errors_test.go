package sss_test

import (
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/go-sss"
)

var _ = Describe("Errors", func() {

	Context("caused by NK faults", func() {
		It("should be returned during encoding", func() {
			// Shamir parameters.
			n := int64(50)
			k := int64(100)
			prime, ok := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
			Ω(ok).Should(Equal(true))
			secret := big.NewInt(1234)
			// Encode the secret.
			_, err := NewShamir(n, k, prime).Encode(secret)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(Equal(NewNKError(n, k)))
		})
	})

	Context("caused by finite field faults", func() {
		It("should be returned during encoding", func() {
			// Shamir parameters.
			n := int64(100)
			k := int64(50)
			prime, ok := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
			Ω(ok).Should(Equal(true))
			secret := big.NewInt(0).Add(prime, big.NewInt(1))
			// Encode the secret.
			_, err := NewShamir(n, k, prime).Encode(secret)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(Equal(NewFiniteFieldError(secret)))
		})
	})
})
