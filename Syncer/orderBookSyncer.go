package Syncer

import (
	"sync"

	"github.com/republicprotocol/republic-go/order"
)

type OrderBookSyncer interface {
	Open(ord *order.Order)
	Match(ord *order.Order)
	Confirm(ord *order.Order)
	Release(ord *order.Order)
	Settle(ord *order.Order)
}

type OrderBookCache struct {
	mu *sync.Mutex

	orders map[string]*order.Order
	status map[string]order.Status
}

func NewOrderBookCache() *OrderBookCache {
	return &OrderBookCache{
		mu:     new(sync.Mutex),
		orders: map[string]*order.Order{},
		status: map[string]order.Status{},
	}
}

func (orderBookCache *OrderBookCache) Open(ord *order.Order) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	orderBookCache.orders[string(ord.ID)] = ord
	if _, ok := orderBookCache.status[string(ord.ID)]; !ok {
		orderBookCache.status[string(ord.ID)] = order.Open
	}
	// todo : do we need to something here ?
}

func (orderBookCache *OrderBookCache) Match(ord *order.Order) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	orderBookCache.orders[string(ord.ID)] = ord
	if status := orderBookCache.status[string(ord.ID)]; status == order.Open {
		orderBookCache.status[string(ord.ID)] = order.Unconfirmed
	}
}

func (orderBookCache *OrderBookCache) Confirm(ord *order.Order) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	orderBookCache.orders[string(ord.ID)] = ord
	if status := orderBookCache.status[string(ord.ID)]; status == order.Unconfirmed {
		orderBookCache.status[string(ord.ID)] = order.Confirmed
	}
}

func (orderBookCache *OrderBookCache) Release(ord *order.Order) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	orderBookCache.orders[string(ord.ID)] = ord
	if status := orderBookCache.status[string(ord.ID)]; status == order.Unconfirmed {
		orderBookCache.status[string(ord.ID)] = order.Open
	}
}

func (orderBookCache *OrderBookCache) Settle(ord *order.Order) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	orderBookCache.orders[string(ord.ID)] = ord
	if status := orderBookCache.status[string(ord.ID)]; status == order.Confirmed {
		orderBookCache.status[string(ord.ID)] = order.Settled
	}
}

// todo : finish this
type OrderBookDB struct {
	mu *sync.Mutex
}

func NewOrderBookDB() *OrderBookDB {
	return &OrderBookDB{
		mu: new(sync.Mutex),
	}
}

func (orderBookDB *OrderBookDB) Open(ord *order.Order){
	panic("unimplemented")
}

func (orderBookDB *OrderBookDB) Match(ord *order.Order){
	panic("unimplemented")
}

func (orderBookDB *OrderBookDB) Confirm(ord *order.Order){
	panic("unimplemented")
}

func (orderBookDB *OrderBookDB) Release(ord *order.Order){
	panic("unimplemented")
}

func (orderBookDB *OrderBookDB) Settle(ord *order.Order){
	panic("unimplemented")
}

type OrderStatusEvent struct {
	ID     order.ID
	Status order.Status
}

func NewOrderStatusEvent(id order.ID, status order.Status) OrderStatusEvent {
	return OrderStatusEvent{
		ID:     id,
		Status: status,
	}
}

type OrderBookStreamer struct {
	mu *sync.Mutex

	maxConnections int
	subscribers    map[string]chan OrderStatusEvent
}

func NewOrderBookStreamer(maxConnection int) *OrderBookStreamer {
	return &OrderBookStreamer{
		mu:             new(sync.Mutex),
		maxConnections: maxConnection,
		subscribers:    map[string]chan OrderStatusEvent{},
	}
}

func (orderBookStreamer *OrderBookStreamer) Subscribe(id string, listener chan OrderStatusEvent) {
	orderBookStreamer.mu.Lock()
	defer orderBookStreamer.mu.Unlock()

	if len(orderBookStreamer.subscribers) >= orderBookStreamer.maxConnections {
		// todo : return an error ?
		return
	}
	orderBookStreamer.subscribers[id] = listener
}

func (orderBookStreamer *OrderBookStreamer) Unsubscribe(id string) {
	orderBookStreamer.mu.Lock()
	defer orderBookStreamer.mu.Unlock()

	delete(orderBookStreamer.subscribers, id)
}

func (orderBookStreamer *OrderBookStreamer) Open(ord *order.Order) {
	orderBookStreamer.mu.Lock()
	defer orderBookStreamer.mu.Unlock()

	for _, subscriber := range orderBookStreamer.subscribers {
		// todo : what if block here
		subscriber <- NewOrderStatusEvent(ord.ID, order.Open)
	}
}

func (orderBookStreamer *OrderBookStreamer) Match(ord *order.Order) {
	orderBookStreamer.mu.Lock()
	defer orderBookStreamer.mu.Unlock()

	for _, subscriber := range orderBookStreamer.subscribers {
		subscriber <- NewOrderStatusEvent(ord.ID, order.Unconfirmed)
	}
}

func (orderBookStreamer *OrderBookStreamer) Confirm(ord *order.Order) {
	orderBookStreamer.mu.Lock()
	defer orderBookStreamer.mu.Unlock()

	for _, subscriber := range orderBookStreamer.subscribers {
		subscriber <- NewOrderStatusEvent(ord.ID, order.Confirmed)
	}
}

func (orderBookStreamer *OrderBookStreamer) Release(ord *order.Order) {
	orderBookStreamer.mu.Lock()
	defer orderBookStreamer.mu.Unlock()

	for _, subscriber := range orderBookStreamer.subscribers {
		subscriber <- NewOrderStatusEvent(ord.ID, order.Open)
	}
}

func (orderBookStreamer *OrderBookStreamer) Settle(ord *order.Order) {
	orderBookStreamer.mu.Lock()
	defer orderBookStreamer.mu.Unlock()

	for _, subscriber := range orderBookStreamer.subscribers {
		subscriber <- NewOrderStatusEvent(ord.ID, order.Settled)
	}
}