package store

import (
	"sync"

	"github.com/republicprotocol/republic-go/order"
)

type OrderBook struct {
	bookMu *sync.RWMutex
	Book   map[string]*order.Order
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		bookMu: new(sync.RWMutex),
		Book:   map[string]*order.Order{},
	}
}

func (orderBook OrderBook) Put(ord *order.Order) order.ID {
	orderBook.bookMu.Lock()
	defer orderBook.bookMu.Unlock()

	orderBook.Book[string(ord.ID)] = ord
	return ord.ID
}

func (orderBook OrderBook) Get(id order.ID) *order.Order {
	orderBook.bookMu.RLock()
	defer orderBook.bookMu.RUnlock()

	if ord, ok := orderBook.Book[string(id)]; ok {
		return ord
	}
	return nil
}

func (orderBook OrderBook) Delete(id order.ID) {
	orderBook.bookMu.Lock()
	defer orderBook.bookMu.Unlock()

	delete(orderBook.Book, string(id))
}
