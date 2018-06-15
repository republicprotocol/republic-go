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

// A Settler settles Computations that have been resolved to matches and have
// been confirmed.
type Settler interface {

	// Settle a Computation that has been resolved to match and has been
	// confirmed. Computations are usually settled after they have been through
	// the Matcher and Confirmer interfaces.
	Settle(com Computation) error
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
func (settler *settler) Settle(com Computation) error {
	buyFragment, err := settler.storer.OrderFragment(com.Buy)
	if err != nil {
		return err
	}
	sellFragment, err := settler.storer.OrderFragment(com.Sell)
	if err != nil {
		return err
	}

	networkID := smpc.NetworkID(com.EpochHash)
	settler.joinOrderMatch(networkID, com, buyFragment, sellFragment)
	return nil
}

func (settler *settler) joinOrderMatch(networkID smpc.NetworkID, com Computation, buyFragment, sellFragment order.Fragment) {
	join := smpc.Join{
		Index: smpc.JoinIndex(buyFragment.Tokens.Index),
		Shares: shamir.Shares{
			buyFragment.Tokens,
			buyFragment.Price.Co, buyFragment.Price.Exp,
			buyFragment.Volume.Co, buyFragment.Volume.Exp,
			buyFragment.MinimumVolume.Co, buyFragment.MinimumVolume.Exp,
			sellFragment.Tokens,
			sellFragment.Price.Co, sellFragment.Price.Exp,
			sellFragment.Volume.Co, sellFragment.Volume.Exp,
			sellFragment.MinimumVolume.Co, sellFragment.MinimumVolume.Exp,
		},
	}
	copy(join.ID[:], crypto.Keccak256(buyFragment.OrderID[:], sellFragment.OrderID[:]))

	err := settler.smpcer.Join(networkID, join, func(joinID smpc.JoinID, values []uint64) {
		if len(values) != 14 {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot join buy = %v, sell = %v: unexpected number of values: %v", base64.StdEncoding.EncodeToString(buyFragment.OrderID[:8]), base64.StdEncoding.EncodeToString(sellFragment.OrderID[:8]), len(values)))
			return
		}
		buy := order.NewOrder(buyFragment.OrderType, buyFragment.OrderParity, buyFragment.OrderExpiry, order.Tokens(values[0]), order.NewCoExp(values[1], values[2]), order.NewCoExp(values[3], values[4]), order.NewCoExp(values[5], values[6]), 0)
		buy.ID = buyFragment.OrderID

		sell := order.NewOrder(sellFragment.OrderType, sellFragment.OrderParity, sellFragment.OrderExpiry, order.Tokens(values[7]), order.NewCoExp(values[8], values[9]), order.NewCoExp(values[10], values[11]), order.NewCoExp(values[12], values[13]), 0)
		sell.ID = sellFragment.OrderID

		settler.settleOrderMatch(com, buy, sell)

	})
	if err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot join buy = %v, sell = %v: %v", base64.StdEncoding.EncodeToString(buyFragment.OrderID[:8]), base64.StdEncoding.EncodeToString(sellFragment.OrderID[:8]), err))
	}
}

func (settler *settler) settleOrderMatch(com Computation, buy, sell order.Order) {
	if err := settler.accounts.Settle(buy, sell); err != nil {
		logger.Error(fmt.Sprintf("cannot settle buy = %v, sell = %v: %v", base64.StdEncoding.EncodeToString(buy.ID[:8]), base64.StdEncoding.EncodeToString(sell.ID[:8]), err))
		return
	}

	com.State = ComputationStateSettled
	if err := settler.storer.InsertComputation(com); err != nil {
		logger.Error(fmt.Sprintf("cannot settle buy = %v, sell = %v: %v", base64.StdEncoding.EncodeToString(buy.ID[:8]), base64.StdEncoding.EncodeToString(sell.ID[:8]), err))
		return
	}

	logger.Info(fmt.Sprintf("$ settled buy = %v, sell = %v", buy.ID, sell.ID))
}
