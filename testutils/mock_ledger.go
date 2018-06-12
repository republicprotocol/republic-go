package testutils

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/republicprotocol/republic-go/order"
)

// ErrOpenOpenedOrder is returned when trying to open an opened order.
var ErrOpenOpenedOrder = errors.New("cannot open order that is already open")

// Order is the internal struct used in the renLedger.
type Order struct {
	status   order.Status
	parity   order.Parity
	priority uint64
}

// RenLedger is a mock implementation of the cal.Ledger.
type RenLedger struct {
	buyOrdersMu *sync.Mutex
	buyOrders   []order.ID

	sellOrdersMu *sync.Mutex
	sellOrders   []order.ID

	ordersMu *sync.Mutex
	orders   map[order.ID]Order
}

// NewRenLedger returns a mock RenLedger.
func NewRenLedger() *RenLedger {
	return &RenLedger{
		buyOrdersMu: new(sync.Mutex),
		buyOrders:   []order.ID{},

		sellOrdersMu: new(sync.Mutex),
		sellOrders:   []order.ID{},

		ordersMu: new(sync.Mutex),
		orders:   map[order.ID]Order{},
	}
}

// OpenBuyOrder in the ledger.
func (renLedger *RenLedger) OpenBuyOrder(signature [65]byte, orderID order.ID) error {
	renLedger.ordersMu.Lock()
	renLedger.buyOrdersMu.Lock()
	defer renLedger.ordersMu.Unlock()
	defer renLedger.buyOrdersMu.Unlock()

	if _, ok := renLedger.orders[orderID]; !ok {
		renLedger.orders[orderID] = Order{
			status:   order.Open,
			parity:   order.ParityBuy,
			priority: binary.LittleEndian.Uint64(orderID[:]),
		}
		renLedger.buyOrders = append(renLedger.buyOrders, orderID)
		return nil
	}

	return errors.New("cannot open order that is already open")
}

// OpenSellOrder in the ledger.
func (renLedger *RenLedger) OpenSellOrder(signature [65]byte, orderID order.ID) error {
	renLedger.ordersMu.Lock()
	renLedger.sellOrdersMu.Lock()
	defer renLedger.ordersMu.Unlock()
	defer renLedger.sellOrdersMu.Unlock()

	if _, ok := renLedger.orders[orderID]; !ok {
		renLedger.orders[orderID] = Order{
			status:   order.Open,
			parity:   order.ParitySell,
			priority: binary.LittleEndian.Uint64(orderID[:]),
		}
		renLedger.sellOrders = append(renLedger.sellOrders, orderID)
		return nil
	}

	return ErrOpenOpenedOrder
}

// CancelOrder in the ledger.
func (renLedger *RenLedger) CancelOrder(signature [65]byte, orderID order.ID) error {
	return renLedger.setOrderStatus(orderID, order.Canceled)
}

// ConfirmOrder confirm a order pair is a match.
func (renLedger *RenLedger) ConfirmOrder(id order.ID, match order.ID) error {
	if err := renLedger.setOrderStatus(id, order.Confirmed); err != nil {
		return fmt.Errorf("cannot confirm order that is not open: %v", err)
	}
	renLedger.setOrderStatus(match, order.Confirmed)
	return nil
}

// Fee returns how much it costs to open an order in the ledger.
func (renLedger *RenLedger) Fee() (*big.Int, error) {
	return big.NewInt(0), nil
}

// Status returns the status of the order by the order ID.
func (renLedger *RenLedger) Status(orderID order.ID) (order.Status, error) {
	renLedger.ordersMu.Lock()
	defer renLedger.ordersMu.Unlock()

	if ord, ok := renLedger.orders[orderID]; ok {
		return ord.status, nil
	}
	return order.Nil, ErrOrderNotFound
}

// Priority returns the priority of the order by the order ID.
func (renLedger *RenLedger) Priority(orderID order.ID) (uint64, error) {
	renLedger.ordersMu.Lock()
	defer renLedger.ordersMu.Unlock()

	if ord, ok := renLedger.orders[orderID]; ok {
		return ord.priority, nil
	}
	return uint64(0), ErrOrderNotFound
}

// Trader returns the trader of the order by the order ID.
func (renLedger *RenLedger) Trader(orderID order.ID) (string, error) {
	panic("unimplemented")
}

// Trader returns the matched order of the order by the order ID.
func (renLedger *RenLedger) OrderMatch(orderID order.ID) (order.ID, error) {
	panic("unimplemented")
}

// Depth returns the block depth since the order been confirmed.
func (renLedger *RenLedger) Depth(orderID order.ID) (uint, error) {
	panic("unimplemented")
}

// BlockNumber returns the block number when the order being last modified.
func (renLedger *RenLedger) BlockNumber(orderID order.ID) (uint, error) {
	panic("unimplemented")
}

// BuyOrders returns a limit of buy orders starting from the offset.
func (renLedger *RenLedger) BuyOrders(offset, limit int) ([]order.ID, error) {
	renLedger.ordersMu.Lock()
	renLedger.buyOrdersMu.Lock()
	defer renLedger.ordersMu.Unlock()
	defer renLedger.buyOrdersMu.Unlock()

	orders := []order.ID{}
	end := offset + limit
	if end > len(renLedger.buyOrders) {
		end = len(renLedger.buyOrders)
	}
	for i := offset; i < end; i++ {
		orderID := renLedger.buyOrders[i]
		if buyOrder, ok := renLedger.orders[orderID]; ok {
			if buyOrder.parity == order.ParityBuy && buyOrder.status == order.Open {
				orders = append(orders, orderID)
			}
		}
	}

	return orders, nil
}

// SellOrders returns a limit sell orders starting from the offset.
func (renLedger *RenLedger) SellOrders(offset, limit int) ([]order.ID, error) {
	renLedger.ordersMu.Lock()
	renLedger.buyOrdersMu.Lock()
	defer renLedger.ordersMu.Unlock()
	defer renLedger.buyOrdersMu.Unlock()

	orders := []order.ID{}
	end := offset + limit
	if end > len(renLedger.sellOrders) {
		end = len(renLedger.sellOrders)
	}
	for i := offset; i < end; i++ {
		orderID := renLedger.sellOrders[i]
		if sellOrder, ok := renLedger.orders[orderID]; ok {
			if sellOrder.parity == order.ParitySell && sellOrder.status == order.Open {
				orders = append(orders, orderID)
			}
		}
	}
	return orders, nil
}

func (renLedger *RenLedger) setOrderStatus(orderID order.ID, status order.Status) error {
	renLedger.ordersMu.Lock()
	defer renLedger.ordersMu.Unlock()

	if _, ok := renLedger.orders[orderID]; ok {
		ord := renLedger.orders[orderID]
		ord.status = status
		renLedger.orders[orderID] = ord
		return nil
	}
	return ErrOrderNotFound
}

func (renLedger *RenLedger) openBuyAndSellOrders(n int) error {
	for i := 0; i < 2*n; i += 2 {
		if err := renLedger.OpenBuyOrder([65]byte{}, [32]byte{byte(i)}); err != nil {
			return err
		}
		if err := renLedger.OpenSellOrder([65]byte{}, [32]byte{byte(i + 1)}); err != nil {
			return err
		}
	}
	return nil
}
