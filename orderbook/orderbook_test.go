package orderbook_test

import (
	"context"
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/testutils"

	"github.com/republicprotocol/republic-go/crypto"
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
			storer := testutils.NewStorer()

			// Create orderbook
			orderbook := NewOrderbook(rsaKey, &syncer, storer)

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
			Expect(storer.NumOrderFragments()).Should(Equal(numberOfOrders))
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
				OrderPriority: Priority(i),
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
