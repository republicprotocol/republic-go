package ingress_test

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ingress"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/smpc"
	"github.com/republicprotocol/republic-go/stackint"
)

var _ = Describe("Ingress", func() {
	var ingress Ingresser
	var darkpool mockDarkpool

	BeforeEach(func() {
		darkpool = newMockDarkpool()
		renLedger := newMockRenLedger()
		swarmer := mockSwarmer{}
		smpcer := mockSmpcer{}
		ingress = NewIngress(&darkpool, &renLedger, &swarmer, &smpcer)
		err := ingress.SyncDarkpool()
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("when opening orders", func() {

		It("should open orders with a sufficient number of order fragments", func() {
			ord, err := createOrder()
			Expect(err).ShouldNot(HaveOccurred())
			fragments, err := ord.Split(5, 4, &smpc.Prime)
			Expect(err).ShouldNot(HaveOccurred())

			orderFragmentMappingIn := make(map[[32]byte][]order.Fragment)
			orderFragments := make([]order.Fragment, len(fragments))
			for i, orderFragment := range fragments {
				orderFragments[i] = *orderFragment
			}
			pods, err := darkpool.Pods()
			Expect(err).ShouldNot(HaveOccurred())
			orderFragmentMappingIn[pods[0].Hash] = orderFragments

			signature := [65]byte{}
			_, err = rand.Read(signature[:])
			Expect(err).ShouldNot(HaveOccurred())

			err = ingress.OpenOrder(signature, ord.ID, orderFragmentMappingIn)
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

			err = ingress.OpenOrder(signature, ord.ID, orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())
		})

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

			err = ingress.OpenOrder(signature, ord.ID, orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())
		})

	})

	Context("when canceling orders", func() {

		It("should cancel orders that are open", func() {
			ord, err := createOrder()
			Expect(err).ShouldNot(HaveOccurred())
			fragments, err := ord.Split(5, 4, &smpc.Prime)
			Expect(err).ShouldNot(HaveOccurred())

			orderFragmentMappingIn := make(map[[32]byte][]order.Fragment)
			orderFragments := make([]order.Fragment, len(fragments))
			for i, orderFragment := range fragments {
				orderFragments[i] = *orderFragment
			}
			pods, err := darkpool.Pods()
			Expect(err).ShouldNot(HaveOccurred())
			orderFragmentMappingIn[pods[0].Hash] = orderFragments

			signature := [65]byte{}
			_, err = rand.Read(signature[:])
			Expect(err).ShouldNot(HaveOccurred())

			err = ingress.OpenOrder(signature, ord.ID, orderFragmentMappingIn)
			Expect(err).ShouldNot(HaveOccurred())

			err = ingress.CancelOrder(signature, ord.ID)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should not cancel orders that are not open", func() {
			ord, err := createOrder()
			Expect(err).ShouldNot(HaveOccurred())

			signature := [65]byte{}
			_, err = rand.Read(signature[:])
			Expect(err).ShouldNot(HaveOccurred())

			err = ingress.CancelOrder(signature, ord.ID)
			Expect(err).Should(HaveOccurred())
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

type mockDarkpool struct {
	pods []cal.Pod
}

func newMockDarkpool() mockDarkpool {
	pod := cal.Pod{
		Hash:      [32]byte{},
		Darknodes: []identity.Address{},
	}
	rand.Read(pod.Hash[:])
	for i := 0; i < 5; i++ {
		ecdsaKey, err := crypto.RandomEcdsaKey()
		if err != nil {
			panic(fmt.Sprintf("cannot create mock darkpool %v", err))
		}
		pod.Darknodes = append(pod.Darknodes, identity.Address(ecdsaKey.Address()))
	}
	return mockDarkpool{
		pods: []cal.Pod{pod},
	}
}

func (darkpool *mockDarkpool) Pods() ([]cal.Pod, error) {
	return darkpool.pods, nil
}

type mockRenLedger struct {
	orderStates map[string]struct{}
}

func newMockRenLedger() mockRenLedger {
	return mockRenLedger{
		orderStates: map[string]struct{}{},
	}
}

func (renLedger *mockRenLedger) OpenOrder(signature [65]byte, orderID order.ID) error {
	if _, ok := renLedger.orderStates[string(orderID)]; !ok {
		renLedger.orderStates[string(orderID)] = struct{}{}
		return nil
	}
	return errors.New("cannot open order that is already open")
}

func (renLedger *mockRenLedger) CancelOrder(signature [65]byte, orderID order.ID) error {
	if _, ok := renLedger.orderStates[string(orderID)]; ok {
		delete(renLedger.orderStates, string(orderID))
		return nil
	}
	return errors.New("cannot cancel order that is not open")
}

func (renLedger *mockRenLedger) Fee() (*big.Int, error) {
	return big.NewInt(0), nil
}

type mockSwarmer struct {
}

func (swarmer *mockSwarmer) Bootstrap(ctx context.Context, multiAddrs identity.MultiAddresses, depth int) <-chan error {
	errs := make(chan error)
	defer close(errs)
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
