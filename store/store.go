package store

import "github.com/republicprotocol/republic-go/order"

type Store interface {
	Get(id order.ID) order.Order
	Put(ord order.Order) order.ID
	Delete(id order.ID)
}
