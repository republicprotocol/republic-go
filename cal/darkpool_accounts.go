package cal

import "github.com/republicprotocol/republic-go/order"

type DarkpoolAccounts interface {

	// Settle the order pair which gets confirmed by the RenLedger
	Settle(buy order.Order, sell order.Order) error

	// Balance returns how many token the trader
	Balance(trader string, token order.Token) (float64, error)
}
