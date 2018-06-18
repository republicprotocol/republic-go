package testutils

import (
	"context"
	"sync/atomic"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
)

// Orderbook implements the Orderbook.Orderbook interface
type Orderbook struct {
	numberOfOrderMatches int64
	hasSynced            bool
	rsaKey               crypto.RsaKey
	orderFragments       map[order.ID]order.Fragment
	orders               map[order.ID]order.Order
}

// NewOrderbook return a new mock Orderbook
func NewOrderbook() (*Orderbook, error) {
	// Generate new RSA key
	rsaKey, err := crypto.RandomRsaKey()
	if err != nil {
		return nil, err
	}

	return &Orderbook{
		numberOfOrderMatches: int64(0),
		hasSynced:            false,
		rsaKey:               rsaKey,
		orderFragments:       make(map[order.ID]order.Fragment),
		orders:               make(map[order.ID]order.Order),
	}, nil
}

// OrderFragment returns the order fragment with given order id
func (book *Orderbook) OrderFragment(orderID order.ID) (order.Fragment, error) {
	return book.orderFragments[orderID], nil
}

// OrderFragment returns the order with given order id
func (book *Orderbook) Order(orderID order.ID) (order.Order, error) {
	return book.orders[orderID], nil
}

func (book *Orderbook) ConfirmOrderMatch(buy order.ID, sell order.ID) error {
	atomic.AddInt64(&book.numberOfOrderMatches, 1)
	return nil
}

func (book *Orderbook) OpenOrder(ctx context.Context, orderFragment order.EncryptedFragment) error {
	var err error
	book.orderFragments[orderFragment.OrderID], err = orderFragment.Decrypt(*book.rsaKey.PrivateKey)
	return err
}

func (book *Orderbook) Sync() (orderbook.ChangeSet, error) {
	if !book.hasSynced {
		changes := make(orderbook.ChangeSet, 5)
		i := 0
		for _, orderFragment := range book.orderFragments {
			changes[i] = orderbook.Change{
				OrderID:       orderFragment.OrderID,
				OrderParity:   orderFragment.OrderParity,
				OrderPriority: orderbook.Priority(i),
				OrderStatus:   order.Open,
			}
			i++
		}
		book.hasSynced = true
		return changes, nil
	}
	return orderbook.ChangeSet{}, nil
}

func (book *Orderbook) AddOrder(ord order.Order) {
	if _, ok := book.orders[ord.ID]; !ok {
		book.orders[ord.ID] = ord
	}
}
