package adapter_test

import (
	"crypto/rand"
	"encoding/base64"
	"sync/atomic"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/http/adapter"
	"github.com/republicprotocol/republic-go/relay"

	"github.com/republicprotocol/republic-go/order"
)

type mockRelayer struct {
	numOpened   int64
	numCanceled int64
}

func (relayer *mockRelayer) OpenOrder(signature [65]byte, orderID order.ID, orderFragmentMapping relay.OrderFragmentMapping) error {
	atomic.AddInt64(&relayer.numOpened, 1)
	return nil
}

func (relayer *mockRelayer) CancelOrder(signature [65]byte, orderID order.ID) error {
	atomic.AddInt64(&relayer.numCanceled, 1)
	return nil
}

func (relayer *mockRelayer) SyncDarkpool() error {
	return nil
}

var _ = Describe("Relay Adapter", func() {

	Context("when opening orders", func() {

		It("should call relayer.OpenOrder if signature and pool is valid", func() {
			relayer := &mockRelayer{}
			relayAdapter := NewRelayAdapter(relayer)
			signatureBytes := [65]byte{}
			_, err := rand.Read(signatureBytes[:])
			Expect(err).ShouldNot(HaveOccurred())
			signature := base64.StdEncoding.EncodeToString(signatureBytes[:])
			orderFragmentMappingIn := make(map[string][]order.Fragment, 1)
			orderFragmentMappingIn["Td2YBy0MRYPYqqBduRmDsIhTySQUlMhPBM+wnNPWKqq="] = []order.Fragment{}
			err = relayAdapter.OpenOrder(signature, orderFragmentMappingIn)
			Expect(err).To(BeNil())
			Expect(atomic.LoadInt64(&relayer.numOpened)).To(Equal(int64(1)))
		})

		It("should not call relayer.OpenOrder if signature is invalid", func() {
			relayer := &mockRelayer{}
			relayAdapter := NewRelayAdapter(relayer)
			signatureBytes := []byte{}
			copy(signatureBytes[:], "incorrect signature")
			orderFragmentMappingIn := make(map[string][]order.Fragment, 1)
			orderFragmentMappingIn["Td2YBy0MRYPYqqBduRmDsIhTySQUlMhPBM+wnNPWKqq="] = []order.Fragment{}
			err := relayAdapter.OpenOrder(string(signatureBytes), orderFragmentMappingIn)
			Expect(err.Error()).To(ContainSubstring("invalid signature length"))
			Expect(atomic.LoadInt64(&relayer.numOpened)).To(Equal(int64(0)))
		})

		It("should not call relayer.OpenOrder if pool hash is invalid", func() {
			relayer := &mockRelayer{}
			relayAdapter := NewRelayAdapter(relayer)
			signatureBytes := [65]byte{}
			_, err := rand.Read(signatureBytes[:])
			Expect(err).ShouldNot(HaveOccurred())
			signature := base64.StdEncoding.EncodeToString(signatureBytes[:])
			orderFragmentMappingIn := make(map[string][]order.Fragment, 1)
			orderFragmentMappingIn["some invalid hash"] = []order.Fragment{}
			err = relayAdapter.OpenOrder(signature, orderFragmentMappingIn)
			Expect(err.Error()).To(ContainSubstring("cannot decode pool hash"))
			Expect(atomic.LoadInt64(&relayer.numOpened)).To(Equal(int64(0)))
		})
	})

	Context("when canceling orders", func() {

		It("should call relayer.CancelOrder if signature and orderID is valid", func() {
			relayer := &mockRelayer{}
			relayAdapter := NewRelayAdapter(relayer)
			signatureBytes := [65]byte{}
			_, err := rand.Read(signatureBytes[:])
			Expect(err).ShouldNot(HaveOccurred())
			signature := base64.StdEncoding.EncodeToString(signatureBytes[:])
			err = relayAdapter.CancelOrder(signature, "Td2YBy0MRYPYqqBduRmDsIhTySQUlMhPBM+wnNPWKqq=")

			Expect(err).To(BeNil())
			Expect(atomic.LoadInt64(&relayer.numCanceled)).To(Equal(int64(1)))
		})

		It("should not call relayer.CancelOrder if orderID is invalid", func() {
			relayer := &mockRelayer{}
			relayAdapter := NewRelayAdapter(relayer)
			signatureBytes := [65]byte{}
			_, err := rand.Read(signatureBytes[:])
			Expect(err).ShouldNot(HaveOccurred())
			signature := base64.StdEncoding.EncodeToString(signatureBytes[:])
			err = relayAdapter.CancelOrder(signature, "")
			Expect(err.Error()).To(ContainSubstring("invalid order id length"))
			Expect(atomic.LoadInt64(&relayer.numCanceled)).To(Equal(int64(0)))
		})
	})
})
