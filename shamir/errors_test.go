package shamir_test

import (
	"fmt"

	"github.com/republicprotocol/republic-go/stackint"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/shamir"
)

var _ = Describe("Errors", func() {

	Context("caused by NK faults", func() {
		It("should be returned during encoding", func() {
			// Shamir parameters.
			n := int64(50)
			k := int64(100)
			prime := stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")
			secret := stackint.FromUint64(1234)
			// Split the secret.
			_, err := Split(n, k, &prime, &secret)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(Equal(NewNKError(n, k)))
			Ω(err.Error()).Should(Equal(fmt.Sprintf("expected n = %v to be greater than or equal to k = %v", n, k)))
		})
	})

	Context("caused by finite field faults", func() {
		It("should be returned during encoding", func() {
			// Shamir parameters.
			n := int64(100)
			k := int64(50)
			prime := stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")
			one := stackint.One()
			secret := prime.Add(&one)
			// Split the secret.
			_, err := Split(n, k, &prime, &secret)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(Equal(NewFiniteFieldError(&secret)))
			Ω(err.Error()).Should(Equal(fmt.Sprintf("expected secret = %v to be within the finite field", secret.String())))
		})
	})
})
