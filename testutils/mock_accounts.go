package testutils

import (
	"github.com/pkg/errors"
	"github.com/republicprotocol/republic-go/order"
)

// ErrInsufficientBalance is returned when the buyer or seller doesn't have
// enough balance in the contract to settle the match.
var ErrInsufficientBalance = errors.New("insufficient balance")

// ErrTraderNotFound is returned when the trader is found in the accounts.
var ErrTraderNotFound = errors.New("trader not found")

// initialBalance returns a map representing the balance of a trader.
var initialBalance = map[order.Token]float64{
	order.TokenBTC: 1000.0,
	order.TokenETH: 1000.0,
	order.TokenDGX: 1000.0,
	order.TokenREN: 1000.0,
}

// DarkpoolAccounts is a mock implementation of the cal.DarkpoolAccounts which
// stores all the settled orders as references and a counter to show how many
// times the Settle function has been called
type DarkpoolAccounts struct {
	settledBuy  []order.Order
	settledSell []order.Order
	counter     int
}

// NewDarkpoolAccounts creates a new mock cal.DarkpoolAccounts.
func NewDarkpoolAccounts() *DarkpoolAccounts {
	return &DarkpoolAccounts{
		settledBuy:  make([]order.Order, 0),
		settledSell: make([]order.Order, 0),

		counter: 0,
	}
}

// Settle the order pair which gets confirmed.
func (accounts *DarkpoolAccounts) Settle(buy order.Order, sell order.Order) error {
	accounts.settledBuy = append(accounts.settledBuy, buy)
	accounts.settledSell = append(accounts.settledSell, sell)
	accounts.counter++

	return nil
}

// Balance returns the balance of a trader for a particular order.Token.
func (accounts *DarkpoolAccounts) Balance(trader string, token order.Token) (float64, error) {
	panic("unimplemented")
}

// Count returns how many time the settle function being called
func (accounts *DarkpoolAccounts) Count() int {
	return accounts.counter
}

// IsSettle returns whether the given order id has been settled.
func (accounts *DarkpoolAccounts) IsSettle(id order.ID) bool {
	for _, ord := range accounts.settledBuy {
		if ord.ID.Equal(id) {
			return true
		}
	}

	for _, ord := range accounts.settledSell {
		if ord.ID.Equal(id) {
			return true
		}
	}

	return false
}
