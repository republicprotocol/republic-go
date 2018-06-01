package cal

import "github.com/republicprotocol/republic-go/order"

type DarkpoolAccounts interface {
	Settle(buy order.Order, sell order.Order) error
	Balance(trader string, token order.Token) (float64, error)
}
