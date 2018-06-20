package ingress_test

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	mathRand "math/rand"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/registry"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ingress"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
)

var _ = Describe("Ingress", func() {

	var rsaKey crypto.RsaKey
	var contract *ingressBinder
	var ingress Ingress
	var done chan struct{}
	var errChSync <-chan error
	var errChProcess <-chan error

	BeforeEach(func() {
		var err error
		done = make(chan struct{})

		rsaKey, err = crypto.RandomRsaKey()
		Expect(err).ShouldNot(HaveOccurred())

		contract = newIngressBinder()

		swarmer := mockSwarmer{}
		orderbookClient := mockOrderbookClient{}
		ingress = NewIngress(contract, &swarmer, &orderbookClient)
		errChSync = ingress.Sync(done)
		errChProcess = ingress.ProcessRequests(done)

		// Consume errors in the background to allow progress when an event occurs
		go captureErrorsFromErrorChannel(errChSync)
		go captureErrorsFromErrorChannel(errChProcess)

		time.Sleep(time.Millisecond)
	})

	AfterEach(func() {
		close(done)

		// Wait for all errors to close
		captureErrorsFromErrorChannel(errChSync)
		captureErrorsFromErrorChannel(errChProcess)

		time.Sleep(time.Second)
	})

	Context("when opening orders", func() {

		It("should open orders with a sufficient number of order fragments", func() {
			ord, err := createOrder()
			Expect(err).ShouldNot(HaveOccurred())
			fragments, err := ord.Split(6, 4)
			Expect(err).ShouldNot(HaveOccurred())

			orderFragmentMappingIn := OrderFragmentMapping{}
			pods, err := contract.Pods()
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
			pods, err := contract.Pods()
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
			pods, err := contract.Pods()
			Expect(err).ShouldNot(HaveOccurred())
			orderFragmentMappingIn[pods[0].Hash] = make([]OrderFragment, 3)

			signature := [65]byte{}
			_, err = rand.Read(signature[:])
			Expect(err).ShouldNot(HaveOccurred())

			err = ingress.OpenOrder(signature, ord.ID, orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())
		})

		It("should not open orders with empty orderFragmentMappings", func() {
			ord, err := createOrder()
			Expect(err).ShouldNot(HaveOccurred())

			orderFragmentMappingIn := OrderFragmentMapping{}

			signature := [65]byte{}
			_, err = rand.Read(signature[:])
			Expect(err).ShouldNot(HaveOccurred())

			err = ingress.OpenOrder(signature, ord.ID, orderFragmentMappingIn)
			Expect(err).Should(HaveOccurred())
		})

		It("should not open orders with unknown pod hashes", func() {
			ord, err := createOrder()
			Expect(err).ShouldNot(HaveOccurred())

			orderFragmentMappingIn := OrderFragmentMapping{}
			orderFragmentMappingIn[[32]byte{byte(1)}] = make([]OrderFragment, 3)

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
			pods, err := contract.Pods()
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

			time.Sleep(time.Second)
			err = ingress.CancelOrder(signature, ord.ID)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should cancel orders that are not open", func() {
			ord, err := createOrder()
			Expect(err).ShouldNot(HaveOccurred())

			signature := [65]byte{}
			_, err = rand.Read(signature[:])
			Expect(err).ShouldNot(HaveOccurred())

			err = ingress.CancelOrder(signature, ord.ID)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})

// ErrOpenOpenedOrder is returned when trying to open an opened order.
var ErrOpenOpenedOrder = errors.New("cannot open order that is already open")

// ingressBinder is a mock implementation of ingress.ContractBinder.
type ingressBinder struct {
	buyOrdersMu *sync.Mutex
	buyOrders   []order.ID

	sellOrdersMu *sync.Mutex
	sellOrders   []order.ID

	ordersMu    *sync.Mutex
	orders      map[order.ID]int
	orderStatus map[order.ID]order.Status

	numberOfDarknodes int
	pods              []registry.Pod
}

// NewingressBinder returns a mock RenLedger.
func newIngressBinder() *ingressBinder {
	pod := registry.Pod{
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

	return &ingressBinder{
		buyOrdersMu: new(sync.Mutex),
		buyOrders:   []order.ID{},

		sellOrdersMu: new(sync.Mutex),
		sellOrders:   []order.ID{},

		ordersMu:    new(sync.Mutex),
		orders:      map[order.ID]int{},
		orderStatus: map[order.ID]order.Status{},

		numberOfDarknodes: 6,
		pods:              []registry.Pod{pod},
	}
}

// OpenBuyOrder in the ledger.
func (binder *ingressBinder) OpenBuyOrder(signature [65]byte, orderID order.ID) error {
	binder.ordersMu.Lock()
	binder.buyOrdersMu.Lock()
	defer binder.ordersMu.Unlock()
	defer binder.buyOrdersMu.Unlock()

	if _, ok := binder.orders[orderID]; !ok {
		binder.orders[orderID] = len(binder.buyOrders)
		binder.buyOrders = append(binder.buyOrders, orderID)
		binder.orderStatus[orderID] = order.Open
		return nil
	}

	return errors.New("cannot open order that is already open")
}

// OpenSellOrder in the ledger.
func (binder *ingressBinder) OpenSellOrder(signature [65]byte, orderID order.ID) error {
	binder.ordersMu.Lock()
	binder.sellOrdersMu.Lock()
	defer binder.ordersMu.Unlock()
	defer binder.sellOrdersMu.Unlock()

	if _, ok := binder.orders[orderID]; !ok {
		binder.orders[orderID] = len(binder.sellOrders)
		binder.sellOrders = append(binder.sellOrders, orderID)
		binder.orderStatus[orderID] = order.Open
		return nil
	}

	return ErrOpenOpenedOrder
}

// CancelOrder in the ledger.
func (binder *ingressBinder) CancelOrder(signature [65]byte, orderID order.ID) error {
	return binder.setOrderStatus(orderID, order.Canceled)
}

func (binder *ingressBinder) Darknodes() (identity.Addresses, error) {
	darknodes := identity.Addresses{}
	for _, pod := range binder.pods {
		darknodes = append(darknodes, pod.Darknodes...)
	}
	return darknodes, nil
}

func (binder *ingressBinder) NextEpoch() (registry.Epoch, error) {
	return binder.Epoch()
}

func (binder *ingressBinder) Epoch() (registry.Epoch, error) {
	darknodes, err := binder.Darknodes()
	if err != nil {
		return registry.Epoch{}, err
	}
	return registry.Epoch{
		Hash:      [32]byte{1},
		Pods:      binder.pods,
		Darknodes: darknodes,
	}, nil
}

func (binder *ingressBinder) MinimumEpochInterval() (*big.Int, error) {
	return big.NewInt(1), nil
}

func (binder *ingressBinder) Pods() ([]registry.Pod, error) {
	return binder.pods, nil
}

func (binder *ingressBinder) setOrderStatus(orderID order.ID, status order.Status) error {
	binder.ordersMu.Lock()
	defer binder.ordersMu.Unlock()

	switch status {
	case order.Open:
		binder.orderStatus[orderID] = order.Open
	case order.Confirmed:
		if binder.orderStatus[orderID] != order.Open {
			return errors.New("order not open")
		}
		binder.orderStatus[orderID] = order.Confirmed
	case order.Canceled:
		if binder.orderStatus[orderID] != order.Open {
			return errors.New("order not open")
		}
		binder.orderStatus[orderID] = order.Canceled
	}
	return nil
}

func createOrder() (order.Order, error) {
	parity := order.ParityBuy
	price := uint64(mathRand.Intn(2000))
	volume := uint64(mathRand.Intn(2000))
	nonce := uint64(mathRand.Intn(1000000000))
	return order.NewOrder(order.TypeLimit, parity, order.SettlementRenEx, time.Now().Add(time.Hour), order.TokensETHREN, order.NewCoExp(price, 26), order.NewCoExp(volume, 26), order.NewCoExp(volume, 26), nonce), nil
}

type mockSwarmer struct {
}

func (swarmer *mockSwarmer) Bootstrap(ctx context.Context, multiAddrs identity.MultiAddresses) error {
	return nil
}

func (swarmer *mockSwarmer) Query(ctx context.Context, query identity.Address, depth int) (identity.MultiAddress, error) {
	return identity.MultiAddress{}, nil
}

func (swarmer *mockSwarmer) MultiAddress() identity.MultiAddress {
	return identity.MultiAddress{}
}

type mockOrderbookClient struct {
}

func (client *mockOrderbookClient) OpenOrder(ctx context.Context, to identity.MultiAddress, orderFragment order.EncryptedFragment) error {
	return nil
}

func captureErrorsFromErrorChannel(errs <-chan error) {
	for range errs {
	}
}
