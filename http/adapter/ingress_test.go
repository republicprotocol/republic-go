package adapter_test

import (
	"crypto/rand"
	"encoding/base64"
	"sync/atomic"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/http/adapter"

	"github.com/republicprotocol/republic-go/ingress"
	"github.com/republicprotocol/republic-go/order"
)

var _ = Describe("Ingress Adapter", func() {

	Context("when opening orders", func() {

		It("should call ingress.OpenOrder if signature and pool is valid", func() {
			ingress := &mockIngress{}
			ingressAdapter := NewIngressAdapter(ingress)
			signatureBytes := [65]byte{}
			_, err := rand.Read(signatureBytes[:])
			Expect(err).ShouldNot(HaveOccurred())
			signature := base64.StdEncoding.EncodeToString(signatureBytes[:])
			orderFragmentMappingIn := OrderFragmentMapping{}
			orderFragmentMappingIn["Td2YBy0MRYPYqqBduRmDsIhTySQUlMhPBM+wnNPWKqq="] = []OrderFragment{}
			err = ingressAdapter.OpenOrder(signature, orderFragmentMappingIn)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(atomic.LoadInt64(&ingress.numOpened)).To(Equal(int64(1)))
		})

		It("should not call ingress.OpenOrder if signature is invalid", func() {
			ingress := &mockIngress{}
			ingressAdapter := NewIngressAdapter(ingress)
			signatureBytes := []byte{}
			copy(signatureBytes[:], "incorrect signature")
			orderFragmentMappingIn := OrderFragmentMapping{}
			orderFragmentMappingIn["Td2YBy0MRYPYqqBduRmDsIhTySQUlMhPBM+wnNPWKqq="] = []OrderFragment{}
			err := ingressAdapter.OpenOrder(string(signatureBytes), orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())
			Expect(atomic.LoadInt64(&ingress.numOpened)).To(Equal(int64(0)))
		})

		It("should not call ingress.OpenOrder if pool hash is invalid", func() {
			ingress := &mockIngress{}
			ingressAdapter := NewIngressAdapter(ingress)
			signatureBytes := [65]byte{}
			_, err := rand.Read(signatureBytes[:])
			Expect(err).ShouldNot(HaveOccurred())
			signature := base64.StdEncoding.EncodeToString(signatureBytes[:])
			orderFragmentMappingIn := OrderFragmentMapping{}
			orderFragmentMappingIn["some invalid hash"] = []OrderFragment{OrderFragment{OrderID: "thisisanorderid"}}
			err = ingressAdapter.OpenOrder(signature, orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())
			Expect(atomic.LoadInt64(&ingress.numOpened)).To(Equal(int64(0)))
		})
	})

	Context("when canceling orders", func() {

		It("should call ingress.CancelOrder if signature and orderID is valid", func() {
			ingress := &mockIngress{}
			ingressAdapter := NewIngressAdapter(ingress)
			signatureBytes := [65]byte{}
			_, err := rand.Read(signatureBytes[:])
			Expect(err).ShouldNot(HaveOccurred())
			signature := base64.StdEncoding.EncodeToString(signatureBytes[:])
			err = ingressAdapter.CancelOrder(signature, "Td2YBy0MRYPYqqBduRmDsIhTySQUlMhPBM+wnNPWKqq=")

			Expect(err).To(BeNil())
			Expect(atomic.LoadInt64(&ingress.numCanceled)).To(Equal(int64(1)))
		})

		It("should not call ingress.CancelOrder if orderID is invalid", func() {
			ingress := &mockIngress{}
			ingressAdapter := NewIngressAdapter(ingress)
			signatureBytes := [65]byte{}
			_, err := rand.Read(signatureBytes[:])
			Expect(err).ShouldNot(HaveOccurred())
			signature := base64.StdEncoding.EncodeToString(signatureBytes[:])
			err = ingressAdapter.CancelOrder(signature, "")
			Expect(err).Should(Equal(ErrInvalidOrderIDLength))
			Expect(atomic.LoadInt64(&ingress.numCanceled)).To(Equal(int64(0)))
		})
	})
})

type mockIngress struct {
	numOpened   int64
	numCanceled int64
}

func (ingress *mockIngress) OpenOrder(signature [65]byte, orderID order.ID, orderFragmentMapping ingress.OrderFragmentMapping) error {
	atomic.AddInt64(&ingress.numOpened, 1)
	return nil
}

func (ingress *mockIngress) CancelOrder(signature [65]byte, orderID order.ID) error {
	atomic.AddInt64(&ingress.numCanceled, 1)
	return nil
}

func (ingress *mockIngress) Sync() error {
	return nil
}
