package ome

import (
	"encoding/base64"
	"fmt"

	"github.com/republicprotocol/republic-go/cal"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/smpc"
)

// A SettleCallback is called when two orders have been settled. This is
// usually when orders have been submitted to the Ethereum blockchain and
// finality has been reached.
type SettleCallback func(buy, sell order.Order)

// A Settler settles a buy order.Order and a sell order.Order.
type Settler interface {

	// Settle a buy order.Order and a sell order.Order. The ξ hash is used to
	// define the ξ in which these orders where successfully matched.
	Settle(ξ [32]byte, buy, sell order.ID) error
}

type settler struct {
	storer   Storer
	smpcer   smpc.Smpcer
	accounts cal.DarkpoolAccounts
}

// NewSettler returns a Settler that settles orders by first using an
// smpc.Smpcer to join all of the composing order.Fragments, and then submits
// them to an Ethereum contract.
func NewSettler(storer Storer, smpcer smpc.Smpcer, accounts cal.DarkpoolAccounts) Settler {
	return &settler{
		storer:   storer,
		smpcer:   smpcer,
		accounts: accounts,
	}
}

// Settle implements the Settler interface.
func (settler *settler) Settle(ξ [32]byte, buy, sell order.ID) error {
	buyFragment, err := matcher.storer.OrderFragment(com.Buy)
	if err != nil {
		return err
	}
	sellFragment, err := matcher.storer.OrderFragment(com.Sell)
	if err != nil {
		return err
	}

	networkID := smpc.NetworkID(ξ)
	settler.joinOrderMatch(networkID, buyFragment, sellFragment)
	return nil
}

func (settler *settler) joinOrderMatch(networkID smpc.NetworkID, buyFragment, sellFragment order.Fragment) {
	join := smpc.Join{
		Index: buyFragment.Tokens.Index,
		Shares: shamir.Shares{
			buyFragment.Tokens,
			buyFragment.Price.Exp, buyFragment.Price.Co,
			buyFragment.Volume.Exp, buyFragment.Volume.Co,
			buyFragment.MinimumVolume.Exp, buyFragment.MinimumVolume.Co,
			sellFragment.Tokens,
			sellFragment.Price.Exp, sellFragment.Price.Co,
			sellFragment.Volume.Exp, sellFragment.Volume.Co,
			sellFragment.MinimumVolume.Exp, sellFragment.MinimumVolume.Co,
		},
	}
	copy(join.ID[:], crypto.Keccak256(buyFragment.OrderID, sellFragment.OrderID))

	err := settler.smpcer.Join(networkID, join, func(joinID join.ID, values []int64) {
		if len(values) != 14 {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot join buy = %v, sell = %v: unexpected number of values: %v", base64.StdEncoding.EncodeToString(buyFragment.OrderID[:8]), base64.StdEncoding.EncodeToString(sellFragment.OrderID[:8]), len(values)))
			return
		}
		buy := order.NewOrder(buyFragment.OrderType, buyFragment.OrderParity, buyFragment.OrderExpiry, Tokens(values[0]), order.NewCoExp(values[1], values[2]), order.NewCoExp(values[3], values[4]), order.NewCoExp(values[5], values[6]), 0)
		buy.ID = buyFragment.OrderID
		sell := order.NewOrder(sellFragment.OrderType, sellFragment.OrderParity, sellFragment.OrderExpiry, Tokens(values[7]), order.NewCoExp(values[8], values[9]), order.NewCoExp(values[10], values[11]), order.NewCoExp(values[12], values[13]), 0)
		sell.ID = sellFragment.OrderID
	})
	if err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot join buy = %v, sell = %v: %v", base64.StdEncoding.EncodeToString(buyFragment.OrderID[:8]), base64.StdEncoding.EncodeToString(sellFragment.OrderID[:8]), err))
	}
}
