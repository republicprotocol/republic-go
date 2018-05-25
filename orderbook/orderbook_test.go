package orderbook_test

import (
	"context"
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/crypto"
	. "github.com/republicprotocol/republic-go/orderbook"

	"github.com/republicprotocol/republic-go/order"
)

var _ = Describe("Orderbook", func() {

	Context("when opening new orders", func() {

		It("should not return an error and must add fragment to storer", func() {
			var err error
			numberOfOrders := 10

			// Generate new RSA key
			rsaKey, err := crypto.RandomRsaKey()
			Expect(err).ShouldNot(HaveOccurred())

			// Create mock syncer and storer
			syncer := newMockSyncer(numberOfOrders)
			storer := newMockStorer(numberOfOrders)

			// Create orderbook
			orderbook := NewOrderbook(rsaKey, &syncer, &storer)

			// Create encryptedOrderFragments
			encryptedOrderFragments := make([]order.EncryptedFragment, numberOfOrders)
			for i := 0; i < numberOfOrders; i++ {
				ord := newOrder(true)
				fragments, err := ord.Split(5, 4)
				encryptedOrderFragments[i], err = fragments[0].Encrypt(rsaKey.PublicKey)
				Expect(err).ShouldNot(HaveOccurred())
			}

			ctx, cancelCtx := context.WithCancel(context.Background())
			defer cancelCtx()

			// Open all encrypted order fragments
			for i := 0; i < numberOfOrders; i++ {
				err = orderbook.OpenOrder(ctx, encryptedOrderFragments[i])
				Expect(err).ShouldNot(HaveOccurred())
			}
			Expect(len(storer.orderFragments)).Should(Equal(numberOfOrders))
		})
	})
})

type mockSyncer struct {
	hasSynced       bool
	numberOfMatches int
	orders          []order.Order
}

func (syncer *mockSyncer) Sync() (ChangeSet, error) {
	if !syncer.hasSynced {
		changes := make(ChangeSet, 5)
		i := 0
		for _, ord := range syncer.orders {
			changes[i] = Change{
				OrderID:       ord.ID,
				OrderParity:   ord.Parity,
				OrderPriority: uint64(i),
				OrderStatus:   order.Open,
			}
			i++
		}
		syncer.hasSynced = true
		return changes, nil
	}

	return ChangeSet{}, nil
}

func (syncer *mockSyncer) ConfirmOrderMatch(order.ID, order.ID) error {
	syncer.numberOfMatches++
	return nil
}

func newMockSyncer(numberOfOrders int) mockSyncer {
	return mockSyncer{
		hasSynced:       false,
		numberOfMatches: 0,
		orders:          make([]order.Order, numberOfOrders),
	}
}

type mockStorer struct {
	orderFragments map[order.ID]order.Fragment
	orders         map[order.ID]order.Order
}

func (storer *mockStorer) InsertOrderFragment(orderFragment order.Fragment) error {
	if _, ok := storer.orderFragments[orderFragment.OrderID]; !ok {
		storer.orderFragments[orderFragment.OrderID] = orderFragment
	}
	return nil
}

func (storer *mockStorer) InsertOrder(order order.Order) error {
	if _, ok := storer.orders[order.ID]; !ok {
		storer.orders[order.ID] = order
	}
	return nil
}

func (storer *mockStorer) OrderFragment(id order.ID) (order.Fragment, error) {
	return storer.orderFragments[id], nil
}

func (storer *mockStorer) Order(id order.ID) (order.Order, error) {
	return storer.orders[id], nil
}

func (storer *mockStorer) RemoveOrderFragment(id order.ID) error {
	delete(storer.orderFragments, id)
	return nil
}

func (storer *mockStorer) RemoveOrder(id order.ID) error {
	delete(storer.orders, id)
	return nil
}

func newMockStorer(numberOfOrders int) mockStorer {
	return mockStorer{
		orderFragments: make(map[order.ID]order.Fragment, numberOfOrders),
		orders:         make(map[order.ID]order.Order, numberOfOrders),
	}
}

func newOrder(isBuy bool) order.Order {
	price := uint64(rand.Intn(2000))
	volume := uint64(rand.Intn(2000))
	nonce := int64(rand.Intn(1000000000))
	parity := order.ParityBuy
	if !isBuy {
		parity = order.ParitySell
	}
	return order.NewOrder(order.TypeLimit, parity, time.Now().Add(time.Hour), order.TokensETHREN, order.NewCoExp(price, 26), order.NewCoExp(volume, 26), order.NewCoExp(volume, 26), nonce)
}
