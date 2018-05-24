package ingress_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"math/big"
	mathRand "math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ingress"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
)

var _ = Describe("Ingress", func() {

	var rsaKey crypto.RsaKey
	var darkpool mockDarkpool
	var ingress Ingress

	BeforeEach(func() {
		var err error
		rsaKey, err = crypto.RandomRsaKey()
		Expect(err).ShouldNot(HaveOccurred())
		darkpool = newMockDarkpool()
		renLedger := newMockRenLedger()
		swarmer := mockSwarmer{}
		orderbookClient := mockOrderbookClient{}
		ingress = NewIngress(&darkpool, &renLedger, &swarmer, &orderbookClient)
		err = ingress.Sync()
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("when opening orders", func() {

		It("should open orders with a sufficient number of order fragments", func() {
			ord, err := createOrder()
			Expect(err).ShouldNot(HaveOccurred())
			fragments, err := ord.Split(5, 4)
			Expect(err).ShouldNot(HaveOccurred())

			orderFragmentMappingIn := OrderFragmentMapping{}
			pods, err := darkpool.Pods()
			Expect(err).ShouldNot(HaveOccurred())
			orderFragmentMappingIn[pods[0].Hash] = []OrderFragment{}
			for i, fragment := range fragments {
				orderFragment := OrderFragment{
					Index: int64(i),
				}
				orderFragment.EncryptedFragment, err = fragment.Encrypt(rsaKey.PublicKey)
				Expect(err).ShouldNot(HaveOccurred())
				orderFragmentMappingIn[pods[0].Hash] = append(orderFragmentMappingIn[pods[0].Hash], orderFragment)
			}

			signature := [65]byte{}
			_, err = rand.Read(signature[:])
			Expect(err).ShouldNot(HaveOccurred())

			err = ingress.OpenOrder(signature, ord.ID, orderFragmentMappingIn)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should not open orders with an insufficient number of order fragments", func() {
			ord, err := createOrder()
			Expect(err).ShouldNot(HaveOccurred())
			fragments, err := ord.Split(int64(1), int64(4/3))
			Expect(err).ShouldNot(HaveOccurred())

			orderFragmentMappingIn := OrderFragmentMapping{}
			pods, err := darkpool.Pods()
			Expect(err).ShouldNot(HaveOccurred())
			orderFragmentMappingIn[pods[0].Hash] = []OrderFragment{}
			for i, fragment := range fragments {
				orderFragment := OrderFragment{
					Index: int64(i),
				}
				orderFragment.EncryptedFragment, err = fragment.Encrypt(rsaKey.PublicKey)
				Expect(err).ShouldNot(HaveOccurred())
				orderFragmentMappingIn[pods[0].Hash] = append(orderFragmentMappingIn[pods[0].Hash], orderFragment)
			}

			signature := [65]byte{}
			_, err = rand.Read(signature[:])
			Expect(err).ShouldNot(HaveOccurred())

			err = ingress.OpenOrder(signature, ord.ID, orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())
		})

		It("should not open orders with malformed order fragments", func() {
			ord, err := createOrder()
			Expect(err).ShouldNot(HaveOccurred())

			orderFragmentMappingIn := OrderFragmentMapping{}
			pods, err := darkpool.Pods()
			Expect(err).ShouldNot(HaveOccurred())
			orderFragmentMappingIn[pods[0].Hash] = make([]OrderFragment, 3)

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
			fragments, err := ord.Split(5, 4)
			Expect(err).ShouldNot(HaveOccurred())

			orderFragmentMappingIn := OrderFragmentMapping{}
			pods, err := darkpool.Pods()
			Expect(err).ShouldNot(HaveOccurred())
			orderFragmentMappingIn[pods[0].Hash] = []OrderFragment{}
			for i, fragment := range fragments {
				orderFragment := OrderFragment{
					Index: int64(i),
				}
				orderFragment.EncryptedFragment, err = fragment.Encrypt(rsaKey.PublicKey)
				Expect(err).ShouldNot(HaveOccurred())
				orderFragmentMappingIn[pods[0].Hash] = append(orderFragmentMappingIn[pods[0].Hash], orderFragment)
			}

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
	parity := order.ParityBuy
	price := uint64(mathRand.Intn(2000))
	volume := uint64(mathRand.Intn(2000))
	nonce := int64(mathRand.Intn(1000000000))
	return order.NewOrder(order.TypeLimit, parity, time.Now().Add(time.Hour), order.TokensETHREN, order.NewCoExp(price, 26), order.NewCoExp(volume, 26), order.NewCoExp(volume, 26), nonce), nil
}

type mockDarkpool struct {
	numberOfDarknodes int
	pods              []cal.Pod
}

func newMockDarkpool() mockDarkpool {
	pod := cal.Pod{
		Hash:      [32]byte{},
		Darknodes: []identity.Address{},
	}
	rand.Read(pod.Hash[:])
	for i := 0; i < 6; i++ {
		ecdsaKey, err := crypto.RandomEcdsaKey()
		if err != nil {
			panic(fmt.Sprintf("cannot create mock darkpool %v", err))
		}
		pod.Darknodes = append(pod.Darknodes, identity.Address(ecdsaKey.Address()))
	}
	return mockDarkpool{
		numberOfDarknodes: 6,
		pods:              []cal.Pod{pod},
	}
}

func (darkpool *mockDarkpool) Darknodes() (identity.Addresses, error) {
	darknodes := identity.Addresses{}
	for _, pod := range darkpool.pods {
		darknodes = append(darknodes, pod.Darknodes...)
	}
	return darknodes, nil
}

func (darkpool *mockDarkpool) Epoch() (cal.Epoch, error) {
	darknodes, err := darkpool.Darknodes()
	if err != nil {
		return cal.Epoch{}, err
	}
	return cal.Epoch{
		Hash:      [32]byte{},
		Pods:      darkpool.pods,
		Darknodes: darknodes,
	}, nil
}

func (darkpool *mockDarkpool) Pods() ([]cal.Pod, error) {
	return darkpool.pods, nil
}

func (darkpool *mockDarkpool) Pod(addr identity.Address) (cal.Pod, error) {
	panic("unimplemented")
}

func (darkpool *mockDarkpool) PublicKey(addr identity.Address) (rsa.PublicKey, error) {
	panic("unimplemented")
}

func (darkpool *mockDarkpool) IsRegistered(addr identity.Address) (bool, error) {
	panic("unimplemented")
}

type mockRenLedger struct {
	orderStates map[string]struct{}
}

func newMockRenLedger() mockRenLedger {
	return mockRenLedger{
		orderStates: map[string]struct{}{},
	}
}

func (renLedger *mockRenLedger) OpenBuyOrder(signature [65]byte, orderID order.ID) error {
	if _, ok := renLedger.orderStates[string(orderID[:])]; !ok {
		renLedger.orderStates[string(orderID[:])] = struct{}{}
		return nil
	}
	return errors.New("cannot open order that is already open")
}

func (renLedger *mockRenLedger) OpenSellOrder(signature [65]byte, orderID order.ID) error {
	if _, ok := renLedger.orderStates[string(orderID[:])]; !ok {
		renLedger.orderStates[string(orderID[:])] = struct{}{}
		return nil
	}
	return errors.New("cannot open order that is already open")
}

func (renLedger *mockRenLedger) CancelOrder(signature [65]byte, orderID order.ID) error {
	if _, ok := renLedger.orderStates[string(orderID[:])]; ok {
		delete(renLedger.orderStates, string(orderID[:]))
		return nil
	}
	return errors.New("cannot cancel order that is not open")
}

func (renLedger *mockRenLedger) ConfirmOrder(id order.ID, matches []order.ID) error {
	if _, ok := renLedger.orderStates[string(id[:])]; ok {
		delete(renLedger.orderStates, string(id[:]))
		for _, matchID := range matches {
			if _, ok := renLedger.orderStates[string(matchID[:])]; ok {
				delete(renLedger.orderStates, string(matchID[:]))
				continue
			}
			return errors.New("cannot confirm order that is not open")
		}
		return nil
	}
	return errors.New("cannot confirm order that is not open")
}

func (renLedger *mockRenLedger) Fee() (*big.Int, error) {
	return big.NewInt(0), nil
}

func (renLedger *mockRenLedger) Status(orderID order.ID) (order.Status, error) {
	panic("unimplemented")
}

func (renLedger *mockRenLedger) Priority(orderID order.ID) (uint64, error) {
	panic("unimplemented")
}

func (renLedger *mockRenLedger) Depth(orderID order.ID) (uint, error) {
	panic("unimplemented")
}

func (renLedger *mockRenLedger) BuyOrders(offset, limit int) ([]order.ID, error) {
	panic("unimplemented")
}

func (renLedger *mockRenLedger) SellOrders(offset, limit int) ([]order.ID, error) {
	panic("unimplemented")
}

type mockSwarmer struct {
}

func (swarmer *mockSwarmer) Bootstrap(ctx context.Context, multiAddrs identity.MultiAddresses) error {
	return nil
}

func (swarmer *mockSwarmer) Query(ctx context.Context, query identity.Address, depth int) (identity.MultiAddress, error) {
	return identity.MultiAddress{}, nil
}

type mockOrderbookClient struct {
}

func (client *mockOrderbookClient) OpenOrder(ctx context.Context, to identity.MultiAddress, orderFragment order.EncryptedFragment) error {
	return nil
}
