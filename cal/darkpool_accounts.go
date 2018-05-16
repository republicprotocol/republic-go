package cal

import "github.com/republicprotocol/republic-go/order"

type DarkpoolAccounts interface {
	Settle(id order.Order, matches []order.Order) error
}
