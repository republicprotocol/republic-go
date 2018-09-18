package grpc_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/grpc"

	"golang.org/x/time/rate"
)

var _ = Describe("Rate-limiter", func() {
	var limits []float64
	var bursts []int
	var addrs []string

	BeforeEach(func() {
		limits = []float64{9, 10, 2, 4, 30}
		bursts = []int{15, 20, 4, 6, 40}
		addrs = []string{
			"8MGs29U416jt8WqaUyRZsXyGqhXxTZ",
			"8MGs29U416jt8WqaUyRZsXyGqhXxTe",
			"8MGs29U416jt8WqaUyRZsXyGqhXxT2",
			"8MGs29U416jt8WqaUyRZsXyGqhXxTq",
			"8MGs29U416jt8WqaUyRZsXyGqhXxTc"}
	})

	Context("when an address sends too many requests and limit has changed", func() {

		It("should block requests after the number of requests has exceeded the new limit", func() {
			rateLimiter := NewRateLimiter(rate.NewLimiter(40, 100), 8, 20)

			for i := 0; i < len(limits); i++ {
				rateLimiter.SetLimit(limits[i])
				rateLimiter.SetBurst(bursts[i])

				for j := 0; j < bursts[i]; j++ {
					Expect(rateLimiter.Allow(addrs[i])).To(BeTrue())
				}

				Expect(rateLimiter.Allow(addrs[i])).To(BeFalse())

				time.Sleep(1 * time.Second)

				for j := 0; j < int(limits[i]); j++ {
					Expect(rateLimiter.Allow(addrs[i])).To(BeTrue())
				}

				Expect(rateLimiter.Allow(addrs[i])).To(BeFalse())

				r := rateLimiter.Reserve(addrs[i])
				Expect(r.OK()).To(BeTrue())
			}
		})
	})
})
