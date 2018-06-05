package adapter_test

import (
	"crypto/rand"
	"encoding/base64"
	mathRand "math/rand"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/crypto"
	. "github.com/republicprotocol/republic-go/http/adapter"

	"github.com/republicprotocol/republic-go/ingress"
	"github.com/republicprotocol/republic-go/order"
)

var _ = Describe("Ingress Adapter", func() {

	Context("when marshaling and unmarshaling order fragment mappings", func() {

		var ord order.Order
		var orderFragmentMappingIn OrderFragmentMapping
		var podHashBytes [32]byte
		var podHash string

		BeforeEach(func() {
			rsaKey, err := crypto.RandomRsaKey()
			Expect(err).ShouldNot(HaveOccurred())
			ord, err = createOrder()
			Expect(err).ShouldNot(HaveOccurred())
			fragments, err := ord.Split(24, 16)
			Expect(err).ShouldNot(HaveOccurred())

			signatureBytes := [65]byte{}
			_, err = rand.Read(signatureBytes[:])
			Expect(err).ShouldNot(HaveOccurred())

			podHashBytes = [32]byte{}
			_, err = rand.Read(podHashBytes[:])
			Expect(err).ShouldNot(HaveOccurred())
			podHash = base64.StdEncoding.EncodeToString(podHashBytes[:])

			orderFragmentMappingIn = OrderFragmentMapping{}
			orderFragmentMappingIn[podHash] = []OrderFragment{}
			for i, fragment := range fragments {
				orderFragment := ingress.OrderFragment{
					Index: int64(i),
				}
				orderFragment.EncryptedFragment, err = fragment.Encrypt(rsaKey.PublicKey)
				Expect(err).ShouldNot(HaveOccurred())
				orderFragmentMappingIn[podHash] = append(
					orderFragmentMappingIn[podHash],
					MarshalOrderFragment(orderFragment))
			}
		})

		It("should return the same data after marshaling and unmarshaling well formed data", func() {
			ordID, orderFragmentMapping, err := UnmarshalOrderFragmentMapping(orderFragmentMappingIn)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ordID).Should(Equal(ord.ID))
			Expect(orderFragmentMapping).Should(HaveLen(1))

			for i, fragment := range orderFragmentMapping[podHashBytes] {
				orderFragmentIn := MarshalOrderFragment(fragment)
				Expect(orderFragmentIn).Should(Equal(orderFragmentMappingIn[podHash][i]))
			}
		})

		It("should return an error for malformed order fragment IDs", func() {
			for i := range orderFragmentMappingIn[podHash] {
				orderFragmentMappingIn[podHash][i].ID = orderFragmentMappingIn[podHash][i].ID[1:]
			}
			_, _, err := UnmarshalOrderFragmentMapping(orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())

			for i := range orderFragmentMappingIn[podHash] {
				orderFragmentMappingIn[podHash][i].ID = "this is invalid"
			}
			_, _, err = UnmarshalOrderFragmentMapping(orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())
		})

		It("should return an error for malformed order fragment IDs", func() {
			for i := range orderFragmentMappingIn[podHash] {
				orderFragmentMappingIn[podHash][i].OrderID = orderFragmentMappingIn[podHash][i].OrderID[1:]
			}
			_, _, err := UnmarshalOrderFragmentMapping(orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())

			for i := range orderFragmentMappingIn[podHash] {
				orderFragmentMappingIn[podHash][i].OrderID = "this is invalid"
			}
			_, _, err = UnmarshalOrderFragmentMapping(orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())
		})

		It("should return an error for malformed pod hashes", func() {
			orderFragmentMappingIn[podHash[16:]] = orderFragmentMappingIn[podHash]
			_, _, err := UnmarshalOrderFragmentMapping(orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())

			delete(orderFragmentMappingIn, podHash[16:])
			orderFragmentMappingIn["this is invalid"] = orderFragmentMappingIn[podHash]
			_, _, err = UnmarshalOrderFragmentMapping(orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())
		})

		It("should return an error for malformed tokens", func() {
			for i := range orderFragmentMappingIn[podHash] {
				orderFragmentMappingIn[podHash][i].Tokens = "this is invalid"
			}
			_, _, err := UnmarshalOrderFragmentMapping(orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())
		})

		It("should return an error for malformed prices", func() {
			for i := range orderFragmentMappingIn[podHash] {
				orderFragmentMappingIn[podHash][i].Price = []string{}
			}
			_, _, err := UnmarshalOrderFragmentMapping(orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())

			for i := range orderFragmentMappingIn[podHash] {
				orderFragmentMappingIn[podHash][i].Price = []string{"this is invalid", "this is also invalid"}
			}
			_, _, err = UnmarshalOrderFragmentMapping(orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())
		})

		It("should return an error for malformed volumes", func() {
			for i := range orderFragmentMappingIn[podHash] {
				orderFragmentMappingIn[podHash][i].Volume = []string{}
			}
			_, _, err := UnmarshalOrderFragmentMapping(orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())

			for i := range orderFragmentMappingIn[podHash] {
				orderFragmentMappingIn[podHash][i].Volume = []string{"this is invalid", "this is also invalid"}
			}
			_, _, err = UnmarshalOrderFragmentMapping(orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())
		})

		It("should return an error for malformed minimum volumes", func() {
			for i := range orderFragmentMappingIn[podHash] {
				orderFragmentMappingIn[podHash][i].MinimumVolume = []string{}
			}
			_, _, err := UnmarshalOrderFragmentMapping(orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())

			for i := range orderFragmentMappingIn[podHash] {
				orderFragmentMappingIn[podHash][i].MinimumVolume = []string{"this is invalid", "this is also invalid"}
			}
			_, _, err = UnmarshalOrderFragmentMapping(orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("when opening orders", func() {

		It("should forward data to the ingress if the signature and mapping are well formed", func() {
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

func (ingress *mockIngress) OpenOrderFragmentsProcess(done chan struct{}) <-chan error {
	return nil
}

func (ingress *mockIngress) OpenOrderProcess(done chan struct{}) <-chan error {
	return nil
}

func (ingress *mockIngress) CancelOrder(signature [65]byte, orderID order.ID) error {
	atomic.AddInt64(&ingress.numCanceled, 1)
	return nil
}

func (ingress *mockIngress) Sync() error {
	return nil
}

func createOrder() (order.Order, error) {
	parity := order.ParityBuy
	price := uint64(mathRand.Intn(2000))
	volume := uint64(mathRand.Intn(2000))
	nonce := int64(mathRand.Intn(1000000000))
	return order.NewOrder(order.TypeLimit, parity, time.Now().Add(time.Hour), order.TokensETHREN, order.NewCoExp(price, 26), order.NewCoExp(volume, 26), order.NewCoExp(volume, 26), nonce), nil
}
