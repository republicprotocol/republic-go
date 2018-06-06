package smpc_test

import (
	"math/rand"
	"sync/atomic"

	"github.com/republicprotocol/republic-go/smpc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/smpc"

	"github.com/republicprotocol/republic-go/shamir"
)

var _ = Describe("Share builder", func() {

	Context("when inserting shares", func() {

		It("should return the value after at least k shares have joined", func() {
			n := int64(24)
			k := int64(16)
			shareBuilderObserver := mockShareBuilderObserver{}
			shareBuilder := NewShareBuilder(k)

			for i := 0; i < 100; i++ {
				shareBuilderObserver.value = 0
				shareBuilder.Observe([32]byte{byte(i)}, [32]byte{byte(i)}, &shareBuilderObserver)

				secret := uint64(rand.Intn(100))
				shares, err := shamir.Split(n, k, secret)
				Expect(err).ShouldNot(HaveOccurred())
				for j := int64(0); j < n; j++ {
					err := shareBuilder.Insert([32]byte{byte(i)}, shares[j])
					if j < k-1 {
						Expect(err).Should(HaveOccurred())
						Expect(err).To(Equal(ErrInsufficientSharesToJoin))
						Expect(shareBuilderObserver.value).To(Equal(uint64(0)))
					} else {
						Expect(err).ShouldNot(HaveOccurred())
						Expect(shareBuilderObserver.value).To(Equal(secret))
					}
				}

				shareBuilder.Observe([32]byte{byte(i)}, [32]byte{byte(i)}, nil)
			}
		})
	})
})

type mockShareBuilderObserver struct {
	value            uint64
	numNotifications uint64
}

func (observer *mockShareBuilderObserver) OnNotifyBuild(id smpc.ComponentID, networkID smpc.NetworkID, value uint64) {
	atomic.StoreUint64(&observer.value, value)
	atomic.AddUint64(&observer.numNotifications, 1)
}
