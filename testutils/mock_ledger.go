package testutils

import (
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/republicprotocol/republic-go/order"
)

// ErrOpenOpenedOrder is returned when trying to open an opened order.
var ErrOpenOpenedOrder = errors.New("cannot open order that is already open")

// ErrOrderNotFound is returnd when the ledger attempts to access an order that
// cannot be found.
var ErrOrderNotFound = errors.New("order not found")

// Constant value of trader address for testing
const (
	GenesisBuyer  = "0x90e6572eF66a11690b09dd594a18f36Cf76055C8"
	GenesisSeller = "0x8DF05f77e8aa74D3D8b5342e6007319A470a64ce"
)

// RenLedger is a mock implementation of the cal.Ledger.
type RenLedger struct {
	buyOrdersMu *sync.Mutex
	buyOrders   []order.ID

	sellOrdersMu *sync.Mutex
	sellOrders   []order.ID

	ordersMu    *sync.Mutex
	orders      map[order.ID]int
	orderStatus map[order.ID]order.Status
}

// NewRenLedger returns a mock RenLedger.
func NewRenLedger() *RenLedger {
	return &RenLedger{
		buyOrdersMu: new(sync.Mutex),
		buyOrders:   []order.ID{},

		sellOrdersMu: new(sync.Mutex),
		sellOrders:   []order.ID{},

		ordersMu:    new(sync.Mutex),
		orders:      map[order.ID]int{},
		orderStatus: map[order.ID]order.Status{},
	}
}

// OpenOrders implements the ledger.
func (renLedger *RenLedger) OpenOrders(signatures [][65]byte, orderIDs []order.ID, orderParities []order.Parity) (int, error) {
	if len(signatures) != len(orderIDs) || len(signatures) != len(orderParities) {
		return 0, errors.New("mismatched order lengths")
	}
	for i := range signatures {
		if orderParities[i] == order.ParityBuy {
			if err := renLedger.OpenBuyOrder(signatures[i], orderIDs[i]); err != nil {
				return i, err
			}
		} else {
			if err := renLedger.OpenSellOrder(signatures[i], orderIDs[i]); err != nil {
				return i, err
			}
		}
		renLedger.orderStatus[orderIDs[i]] = order.Open
	}
	return len(signatures), nil
}

// OpenBuyOrder in the ledger.
func (renLedger *RenLedger) OpenBuyOrder(signature [65]byte, orderID order.ID) error {
	renLedger.ordersMu.Lock()
	renLedger.buyOrdersMu.Lock()
	defer renLedger.ordersMu.Unlock()
	defer renLedger.buyOrdersMu.Unlock()

	if _, ok := renLedger.orders[orderID]; !ok {
		renLedger.orders[orderID] = len(renLedger.buyOrders)
		renLedger.buyOrders = append(renLedger.buyOrders, orderID)
		renLedger.orderStatus[orderID] = order.Open
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
		renLedger.orders[orderID] = len(renLedger.sellOrders)
		renLedger.sellOrders = append(renLedger.sellOrders, orderID)
		renLedger.orderStatus[orderID] = order.Open
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
	if err := renLedger.setOrderStatus(match, order.Confirmed); err != nil {
		return fmt.Errorf("cannot confirm order that is not open: %v", err)
	}
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

	if status, ok := renLedger.orderStatus[orderID]; ok {
		return status, nil
	}
	return order.Nil, ErrOrderNotFound
}

// Priority returns the priority of the order by the order ID.
func (renLedger *RenLedger) Priority(orderID order.ID) (uint64, error) {
	renLedger.ordersMu.Lock()
	defer renLedger.ordersMu.Unlock()

	for _, id := range renLedger.buyOrders {
		if orderID.Equal(id) {
			return uint64(order.ParityBuy), nil
		}
	}
	for _, id := range renLedger.sellOrders {
		if orderID.Equal(id) {
			return uint64(order.ParitySell), nil
		}
	}

	return uint64(0), ErrOrderNotFound
}

// Trader returns the trader of the order by the order ID.
func (renLedger *RenLedger) Trader(orderID order.ID) (string, error) {
	return GenesisBuyer, nil
}

// Trader returns the matched order of the order by the order ID.
func (renLedger *RenLedger) OrderMatch(orderID order.ID) (order.ID, error) {
	renLedger.ordersMu.Lock()
	defer renLedger.ordersMu.Unlock()

	if renLedger.orderStatus[orderID] != order.Confirmed {
		return order.ID{}, errors.New("order is not open ")
	}
	for i, id := range renLedger.buyOrders {
		if orderID.Equal(id) {
			return renLedger.sellOrders[i], nil
		}
	}
	for i, id := range renLedger.sellOrders {
		if orderID.Equal(id) {
			return renLedger.buyOrders[i], nil
		}
	}
	return order.ID{}, ErrOrderNotFound
}

// Depth returns the block depth since the order been confirmed.
func (renLedger *RenLedger) Depth(orderID order.ID) (uint, error) {
	return 100, nil
}

// BlockNumber returns the block number when the order being last modified.
func (renLedger *RenLedger) BlockNumber(orderID order.ID) (uint, error) {
	return 100, nil
}

// BuyOrders returns a limit of buy orders starting from the offset.
func (renLedger *RenLedger) BuyOrders(offset, limit int) ([]order.ID, error) {
	renLedger.ordersMu.Lock()
	renLedger.buyOrdersMu.Lock()
	defer renLedger.ordersMu.Unlock()
	defer renLedger.buyOrdersMu.Unlock()

	if offset > len(renLedger.buyOrders) {
		return []order.ID{}, errors.New("index out of range")
	}
	end := offset + limit
	if end > len(renLedger.buyOrders) {
		end = len(renLedger.buyOrders)
	}
	return renLedger.buyOrders[offset:end], nil
}

// SellOrders returns a limit sell orders starting from the offset.
func (renLedger *RenLedger) SellOrders(offset, limit int) ([]order.ID, error) {
	renLedger.ordersMu.Lock()
	renLedger.buyOrdersMu.Lock()
	defer renLedger.ordersMu.Unlock()
	defer renLedger.buyOrdersMu.Unlock()

	if offset > len(renLedger.sellOrders) {
		return []order.ID{}, errors.New("index out of range")
	}
	end := offset + limit
	if end > len(renLedger.sellOrders) {
		end = len(renLedger.sellOrders)
	}
	return renLedger.sellOrders[offset:end], nil
}

func (renLedger *RenLedger) setOrderStatus(orderID order.ID, status order.Status) error {
	renLedger.ordersMu.Lock()
	defer renLedger.ordersMu.Unlock()

	switch status {
	case order.Open:
		renLedger.orderStatus[orderID] = order.Open
	case order.Confirmed:
		if renLedger.orderStatus[orderID] != order.Open {
			return errors.New("order not open")
		}
		renLedger.orderStatus[orderID] = order.Confirmed
	case order.Canceled:
		if renLedger.orderStatus[orderID] != order.Open {
			return errors.New("order not open")
		}
		renLedger.orderStatus[orderID] = order.Canceled
	}

	return nil
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
