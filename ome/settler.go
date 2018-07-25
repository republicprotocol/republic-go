package ome

import (
	"fmt"
	"log"

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
	computationStore ComputationStorer
	smpcer           smpc.Smpcer
	contract         ContractBinder
}

// NewSettler returns a Settler that settles orders by first using an
// smpc.Smpcer to join all of the composing order.Fragments, and then submits
// them to an Ethereum contract.
func NewSettler(computationStore ComputationStorer, smpcer smpc.Smpcer, contract ContractBinder) Settler {
	return &settler{
		computationStore: computationStore,
		smpcer:           smpcer,
		contract:         contract,
	}
}

// Settle implements the Settler interface.
func (settler *settler) Settle(com Computation) error {
	networkID := smpc.NetworkID(com.Epoch)
	settler.joinOrderMatch(networkID, com)
	return nil
}

func (settler *settler) joinOrderMatch(networkID smpc.NetworkID, com Computation) {
	join := smpc.Join{
		Index: smpc.JoinIndex(com.Buy.Tokens.Index),
		Shares: shamir.Shares{
			com.Buy.Tokens,
			com.Buy.Price.Co, com.Buy.Price.Exp,
			com.Buy.Volume.Co, com.Buy.Volume.Exp,
			com.Buy.MinimumVolume.Co, com.Buy.MinimumVolume.Exp,
			com.Buy.Nonce,
			com.Sell.Tokens,
			com.Sell.Price.Co, com.Sell.Price.Exp,
			com.Sell.Volume.Co, com.Sell.Volume.Exp,
			com.Sell.MinimumVolume.Co, com.Sell.MinimumVolume.Exp,
			com.Sell.Nonce,
		},
	}
	copy(join.ID[:], com.ID[:])
	join.ID[32] = byte(ResolveStageSettlement)

	err := settler.smpcer.Join(networkID, join, func(joinID smpc.JoinID, values []uint64) {
		if len(values) != 16 {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot join buy = %v, sell = %v: unexpected number of values: %v", com.Buy.OrderID, com.Sell.OrderID, len(values)))
			return
		}
		buy := order.NewOrder(com.Buy.OrderType, com.Buy.OrderParity, com.Buy.OrderSettlement, com.Buy.OrderExpiry, order.Tokens(values[0]), order.NewCoExp(values[1], values[2]), order.NewCoExp(values[3], values[4]), order.NewCoExp(values[5], values[6]), values[7])
		sell := order.NewOrder(com.Sell.OrderType, com.Sell.OrderParity, com.Sell.OrderSettlement, com.Sell.OrderExpiry, order.Tokens(values[8]), order.NewCoExp(values[9], values[10]), order.NewCoExp(values[11], values[12]), order.NewCoExp(values[13], values[14]), values[15])

		settler.settleOrderMatch(com, buy, sell)

	})
	if err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot join buy = %v, sell = %v: %v", com.Buy.OrderID, com.Sell.OrderID, err))
	}
}

func (settler *settler) settleOrderMatch(com Computation, buy, sell order.Order) {
	if err := settler.contract.Settle(buy, sell); err != nil {
		log.Printf("[error] (settle) cannot execute settlement buy = %v, sell = %v: %v", buy.ID, sell.ID, err)
		return
	}
	log.Printf("[info] (settle) 💰💰💰 buy = %v, sell = %v 💰💰💰", buy.ID, sell.ID)

	com.State = ComputationStateSettled
	if err := settler.computationStore.PutComputation(com); err != nil {
		log.Printf("[error] (settle) cannot store settlement buy = %v, sell = %v: %v", buy.ID, sell.ID, err)
		return
	}
}
