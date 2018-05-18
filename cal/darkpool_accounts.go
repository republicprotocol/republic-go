package cal

import "github.com/republicprotocol/republic-go/order"

type DarkpoolAccounts interface {
	Settle(buy order.Order, sell order.Order) error
}
