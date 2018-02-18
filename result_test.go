package compute_test

import (
	"math/big"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	compute "github.com/republicprotocol/go-order-compute"
)

var _ = Describe("Results and result fragments", func() {

	n := int64(8)
	k := int64(6)
	prime, _ := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)

	Context("when serializing IDs to strings", func() {

		It("should return the string representation of the ID", func() {
			result, err := computeRandomResult(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(result.ID.String()).Should(Equal(string(result.ID)))
		})
	})
})

func computeRandomResult(n, k int64, prime *big.Int) (*compute.Result, error) {
	lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
	if err != nil {
		return nil, err
	}
	rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
	if err != nil {
		return nil, err
	}
	return computeResultFromOrderFragments(lhs, rhs, n, prime)
}
