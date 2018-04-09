package relay

import (
	"fmt"

	"github.com/republicprotocol/republic-go/order"
)

type Relay struct {
}

func SendOrderToDarkOcean (order order.Order) error {
	fmt.Println(string(order.Type))
	return nil
}
