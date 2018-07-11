package testutils

import (
	"context"
	"errors"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/registry"
)

// An EmptyOrderbook implements the orderbook.Orderbook interface but it never
// produces notifications.
type EmptyOrderbook struct {
}

// NewEmptyOrderbook returns a new EmptyOrderbook.
func NewEmptyOrderbook() *EmptyOrderbook {
	return &EmptyOrderbook{}
}

// OpenOrder implements the orderbook.Orderbook interface.
func (mock *EmptyOrderbook) OpenOrder(ctx context.Context, encryptedOrderFragment order.EncryptedFragment) error {
	return nil
}

// Sync implements the orderbook.Orderbook interface.
func (mock *EmptyOrderbook) Sync(done <-chan struct{}) (<-chan orderbook.Notification, <-chan error) {
	notifications := make(chan orderbook.Notification)
	errs := make(chan error)

	go func() {
		defer close(notifications)
		defer close(errs)

		<-done
	}()

	return notifications, errs
}

// OnChangeEpoch implements the orderbook.Orderbook interface.
func (mock *EmptyOrderbook) OnChangeEpoch(epoch registry.Epoch) {
}

// A RandOrderbook implements the orderbook.Orderbook interface and produces a
// notification for every opened order.Fragment and randomly produces
// cancelations and confirmations.
type RandOrderbook struct {
	rsaKey              crypto.RsaKey
	orderFragments      []order.Fragment
	orderFragmentsQueue []order.Fragment
	traders             []string
}

// NewRandOrderbook returns a new RandOrderbook that is initially empty and
// uses a crypto.RsaKey to decrypt order.EncryptedFragments.
func NewRandOrderbook(rsaKey crypto.RsaKey) *RandOrderbook {
	return &RandOrderbook{
		rsaKey:              rsaKey,
		orderFragments:      []order.Fragment{},
		orderFragmentsQueue: []order.Fragment{},
		traders:             []string{"trader1", "trader2", "trader3", "trader4"},
	}
}

// OpenOrder implements the orderbook.Orderbook interface.
func (mock *RandOrderbook) OpenOrder(ctx context.Context, encryptedOrderFragment order.EncryptedFragment) error {
	orderFragment, err := encryptedOrderFragment.Decrypt(mock.rsaKey.PrivateKey)
	if err != nil {
		return err
	}
	mock.orderFragmentsQueue = append(mock.orderFragmentsQueue, orderFragment)
	return nil
}

// Sync implements the orderbook.Orderbook interface.
func (mock *RandOrderbook) Sync(done <-chan struct{}) (<-chan orderbook.Notification, <-chan error) {
	notifications := make(chan orderbook.Notification)
	errs := make(chan error)

	go func() {
		defer close(notifications)
		defer close(errs)

		ticker := time.NewTicker(4 * time.Second)
		defer ticker.Stop()

		blockNumber := uint64(0)
		for {
			select {
			case <-done:
			case <-ticker.C:
				// Simulate the block number increasing
				blockNumber++

				// Generate an open notification for all pending order.Fagments
				for i := range mock.orderFragmentsQueue {
					notification := orderbook.NotificationOpenOrder{
						OrderID:       mock.orderFragmentsQueue[i].OrderID,
						OrderFragment: mock.orderFragmentsQueue[i],
						Trader:        mock.traders[rand.Intn(len(mock.traders))],
						BlockNumber:   blockNumber,
					}
					select {
					case <-done:
					case notifications <- notification:
					}
				}
				mock.orderFragments = append(mock.orderFragments, mock.orderFragmentsQueue...)
				mock.orderFragmentsQueue = []order.Fragment{}

				if len(mock.orderFragments) > 0 {
					// Randomly remove an order.Fragment
					n := rand.Intn(len(mock.orderFragments))
					orderFragment := mock.orderFragments[n]
					mock.orderFragments[n] = mock.orderFragments[len(mock.orderFragments)-1]
					mock.orderFragments = mock.orderFragments[:len(mock.orderFragments)-1]

					// Randomly generate a closure notification for it
					var notification orderbook.Notification
					r := rand.Intn(100)
					if r < 50 {
						notification = orderbook.NotificationConfirmOrder{
							OrderID: orderFragment.OrderID,
						}
					} else {
						notification = orderbook.NotificationCancelOrder{
							OrderID: orderFragment.OrderID,
						}
					}
					select {
					case <-done:
					case notifications <- notification:
					}
				}
			}
		}
	}()

	return notifications, errs
}

// OnChangeEpoch implements the orderbook.Orderbook interface.
func (mock *RandOrderbook) OnChangeEpoch(epoch registry.Epoch) {
}

type MockContractBinder struct {
	ordersMu    *sync.Mutex
	orders      []order.ID
	orderStatus map[order.ID]order.Status
	traders     map[order.ID]string
}

// NewMockContractBinder returns a mockContractBinder
func NewMockContractBinder() *MockContractBinder {
	return &MockContractBinder{
		ordersMu:    new(sync.Mutex),
		orders:      []order.ID{},
		orderStatus: map[order.ID]order.Status{},
		traders:     map[order.ID]string{},
	}
}

func (binder *MockContractBinder) Orders(offset, limit int) ([]order.ID, []order.Status, []string, error) {
	statuses := make([]order.Status, 0, len(binder.orders))
	traders := make([]string, 0, len(binder.orders))

	if offset > len(binder.orders) {
		return []order.ID{}, []order.Status{}, []string{}, errors.New("index out of range")
	}

	end := offset + limit
	if end > len(binder.orders) {
		end = len(binder.orders)
	}

	for i := offset; i < end; i++ {
		id := binder.orders[i]
		if status, ok := binder.orderStatus[id]; ok {
			statuses = append(statuses, status)
		}
		if trader, ok := binder.traders[id]; ok {
			traders = append(traders, trader)
		}

	}

	return binder.orders[offset:end], statuses, traders, nil
}

func (binder *MockContractBinder) BlockNumber(orderID order.ID) (*big.Int, error) {
	for i, ord := range binder.orders {
		if ord == orderID {
			return big.NewInt(int64(i)), nil
		}
	}
	return &big.Int{}, orderbook.ErrOrderNotFound
}

func (binder *MockContractBinder) Status(orderID order.ID) (order.Status, error) {
	if status, ok := binder.orderStatus[orderID]; ok {
		return status, nil
	}
	return order.Open, orderbook.ErrOrderNotFound
}

func (binder *MockContractBinder) MinimumEpochInterval() (*big.Int, error) {
	return big.NewInt(2), nil
}

func (binder *MockContractBinder) OpenMatchingOrders(n int) []order.Order {
	binder.ordersMu.Lock()
	defer binder.ordersMu.Unlock()

	orders := []order.Order{}
	for i := 0; i < n; i++ {
		buy, sell := RandomOrderMatch()
		if _, ok := binder.orderStatus[buy.ID]; !ok {
			binder.orders = append(binder.orders, buy.ID)
			binder.orderStatus[buy.ID] = order.Open
			binder.traders[buy.ID] = string(i)
			orders = append(orders, buy)
		}
		if _, ok := binder.orderStatus[sell.ID]; !ok {
			binder.orders = append(binder.orders, sell.ID)
			binder.orderStatus[sell.ID] = order.Open
			binder.traders[sell.ID] = string(i)
			orders = append(orders, sell)
		}
	}
	return orders
}
