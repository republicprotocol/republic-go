package main

import (
	"fmt"
	"time"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
)

func main() {
	orderID := [32]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	tokens := shamir.Share{
		Index: 1,
	}

	price := order.CoExpShare{
		Co: shamir.Share{
			Index: 1234567834452542542,
			Value: 5409098098584098524,
		},
		Exp: shamir.Share{
			Index: 4235515145415414342,
			Value: 768724151521512563,
		},
	}
	minVolume := order.CoExpShare{
		Co: shamir.Share{
			Index: 2453343252545422462,
			Value: 2499435246542609086,
		},
		Exp: shamir.Share{
			Index: 4354295485094850492,
			Value: 2453243251111122462,
		},
	}
	maxVolume := order.CoExpShare{
		Co: shamir.Share{
			Index: 3439349840928343209,
			Value: 3434345902840329841,
		},
		Exp: shamir.Share{
			Index: 1243542,
			Value: 54524642,
		},
	}

	fragment := order.Fragment{
		OrderID:       orderID,
		OrderType:     order.TypeLimit,
		OrderParity:   order.ParityBuy,
		Tokens:        tokens,
		Price:         price,
		Volume:        maxVolume,
		MinimumVolume: minVolume,
		OrderExpiry:   time.Now(),
	}
	fragment.ID = order.FragmentID(fragment.Hash())

	fmt.Println(fragment.Bytes())
}
