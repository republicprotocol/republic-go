package relay_test

import (
	"context"
	"crypto/rand"
	mathrand "math/rand"
	"time"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/stackint"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/order"
	. "github.com/republicprotocol/republic-go/relay"
	"github.com/republicprotocol/republic-go/smpc"
)

type mockDarkpool struct {
}

func (darkpool *mockDarkpool) Pods() ([]cal.Pod, error) {
	pods := make([]cal.Pod, 1)
	darknodeAddrs := make([]identity.Address, 3)
	var podHash [32]byte
	copy(podHash[:], "Td2YBy0MRYPYqqBduRmDsIhTySQUlMhPBM+wnNPWKqq=")
	pods[0] = cal.Pod{
		Hash:      podHash,
		Darknodes: darknodeAddrs,
	}
	return pods, nil
}

func (renLedger *mockRenLedger) OpenOrder(signature [65]byte, orderID order.ID) error {
	return nil
}

type mockRenLedger struct {
}

func (renLedger *mockRenLedger) CancelOrder(signature [65]byte, orderID order.ID) error {
	return nil
}

type mockSwarmer struct {
}

func (swarmer *mockSwarmer) Bootstrap(ctx context.Context, multiAddrs identity.MultiAddresses, depth int) <-chan error {
	errs := make(chan error, 1)
	return errs
}

func (swarmer *mockSwarmer) Query(ctx context.Context, query identity.Address, depth int) (identity.MultiAddress, error) {
	return identity.MultiAddress{}, nil
}

type mockSmpcer struct {
}

func (smpcer *mockSmpcer) OpenOrder(ctx context.Context, to identity.MultiAddress, orderFragment order.Fragment) error {
	return nil
}

var _ = Describe("Relay", func() {
	var relay Relay
	var darkpool mockDarkpool

	BeforeEach(func() {
		darkpool = mockDarkpool{}
		renLedger := mockRenLedger{}
		swarmer := mockSwarmer{}
		smpcer := mockSmpcer{}
		relay = NewRelay(&darkpool, &renLedger, &swarmer, &smpcer)
	})

	Context("when opening orders", func() {

		It("should open orders with a sufficient number of order fragments", func() {
			ord, err := createOrder()
			Expect(err).ShouldNot(HaveOccurred())
			fragments, err := ord.Split(int64(3), int64(8/3), &smpc.Prime)
			Expect(err).ShouldNot(HaveOccurred())

			orderFragmentMappingIn := make(map[[32]byte][]order.Fragment)
			orderFragments := make([]order.Fragment, len(fragments))
			for i, orderFragment := range fragments {
				orderFragments[i] = *orderFragment
			}
			pod, err := darkpool.Pods()
			Expect(err).ShouldNot(HaveOccurred())
			orderFragmentMappingIn[pod[0].Hash] = orderFragments

			signature := [65]byte{}
			_, err = rand.Read(signature[:])
			Expect(err).ShouldNot(HaveOccurred())

			err = relay.OpenOrder(signature, ord.ID, orderFragmentMappingIn)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should not open orders with an insufficient number of order fragments", func() {
			ord, err := createOrder()
			Expect(err).ShouldNot(HaveOccurred())
			fragments, err := ord.Split(int64(1), int64(4/3), &smpc.Prime)
			Expect(err).ShouldNot(HaveOccurred())

			orderFragmentMappingIn := make(map[[32]byte][]order.Fragment)
			orderFragments := make([]order.Fragment, len(fragments))
			for i, orderFragment := range fragments {
				orderFragments[i] = *orderFragment
			}
			pod, err := darkpool.Pods()
			Expect(err).ShouldNot(HaveOccurred())
			orderFragmentMappingIn[pod[0].Hash] = orderFragments

			signature := [65]byte{}
			_, err = rand.Read(signature[:])
			Expect(err).ShouldNot(HaveOccurred())

			err = relay.OpenOrder(signature, ord.ID, orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())
		})

		//FIXME: This test should pass once order fragments are validated in the relay
		It("should not open orders with malformed order fragments", func() {
			ord, err := createOrder()
			Expect(err).ShouldNot(HaveOccurred())

			orderFragmentMappingIn := make(map[[32]byte][]order.Fragment)
			orderFragments := make([]order.Fragment, 3)
			pod, err := darkpool.Pods()
			Expect(err).ShouldNot(HaveOccurred())
			orderFragmentMappingIn[pod[0].Hash] = orderFragments

			signature := [65]byte{}
			_, err = rand.Read(signature[:])
			Expect(err).ShouldNot(HaveOccurred())

			err = relay.OpenOrder(signature, ord.ID, orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())
		})

		It("should not open orders that have not been signed correctly", func() {
			Expect("implemented").To(Equal("unimplemented"))
		})

	})

	Context("when canceling orders", func() {

		It("should cancel orders that are open", func() {
			Expect("implemented").To(Equal("unimplemented"))
		})

		It("should not cancel orders that are not open", func() {
			Expect("implemented").To(Equal("unimplemented"))
		})

		It("should not cancel orders that have not been signed correctly", func() {
			Expect("implemented").To(Equal("unimplemented"))
		})

	})

	Context("when getting orders", func() {

		It("should return an order for verified traders", func() {
			Expect("implemented").To(Equal("unimplemented"))
		})

		It("should not return an order for unverified traders", func() {
			Expect("implemented").To(Equal("unimplemented"))
		})

		It("should return orders for verified traders", func() {
			Expect("implemented").To(Equal("unimplemented"))
		})

		It("should not return orders for unverified traders", func() {
			Expect("implemented").To(Equal("unimplemented"))
		})
	})
})

func createOrder() (order.Order, error) {
	price := mathrand.Intn(1000000000)
	amount := mathrand.Intn(1000000000)

	nonce, err := stackint.Random(rand.Reader, &smpc.Prime)
	if err != nil {
		return order.Order{}, err
	}

	parity := order.ParityBuy

	return *order.NewOrder(order.TypeLimit, parity, time.Now().Add(time.Hour), order.CurrencyCodeETH, order.CurrencyCodeBTC, stackint.FromUint(uint(price)), stackint.FromUint(uint(amount)), stackint.FromUint(uint(amount)), nonce), nil
}
