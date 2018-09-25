package ome

import (
	"fmt"
	"log"
	"math/big"

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
	computationStore    ComputationStorer
	smpcer              smpc.Smpcer
	contract            ContractBinder
	minimumSettleVolume uint64 // In units of 1e-12 ETH
}

// NewSettler returns a Settler that settles orders by first using an
// smpc.Smpcer to join all of the composing order.Fragments, and then submits
// them to an Ethereum contract.
func NewSettler(computationStore ComputationStorer, smpcer smpc.Smpcer, contract ContractBinder, minimumSettleVolume uint64) Settler {
	return &settler{
		computationStore:    computationStore,
		smpcer:              smpcer,
		contract:            contract,
		minimumSettleVolume: minimumSettleVolume,
	}
}

// Settle implements the Settler interface.
func (settler *settler) Settle(com Computation) error {
	networkID := smpc.NetworkID(com.Epoch)
	settler.joinOrderMatch(networkID, com)
	return nil
}

func (settler *settler) joinOrderMatch(networkID smpc.NetworkID, com Computation) {

	// Create the blinding to verify the computation
	blinding := shamir.Blinding{
		Int: big.NewInt(0).Add(com.Buy.Blinding.Int, big.NewInt(0).Sub(shamir.CommitP, com.Sell.Blinding.Int)),
	}
	blinding.Mod(blinding.Int, shamir.CommitP)

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
		Blindings: shamir.Blindings{
			com.Buy.Blinding,
			com.Buy.Blinding, com.Buy.Blinding,
			com.Buy.Blinding, com.Buy.Blinding,
			com.Buy.Blinding, com.Buy.Blinding,
			com.Buy.Blinding,

			com.Sell.Blinding,
			com.Sell.Blinding, com.Sell.Blinding,
			com.Sell.Blinding, com.Sell.Blinding,
			com.Sell.Blinding, com.Sell.Blinding,
			com.Sell.Blinding,
		},
	}
	copy(join.ID[:], com.ID[:])
	join.ID[32] = byte(ResolveStageSettlement)

	err := settler.smpcer.Join(networkID, join, func(joinID smpc.JoinID, values []uint64) {
		if len(values) != 16 {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot join buy = %v, sell = %v: unexpected number of values: %v", com.Buy.OrderID, com.Sell.OrderID, len(values)))
			return
		}
		buy := order.NewOrder(com.Buy.OrderParity, com.Buy.OrderType, com.Buy.OrderExpiry, com.Buy.OrderSettlement, order.Tokens(values[0]), order.PriceFromCoExp(values[1], values[2]), order.VolumeFromCoExp(values[3], values[4]), order.VolumeFromCoExp(values[5], values[6]), values[7])
		buy.MinimumVolume = order.VolumeFromCoExp(values[5], values[6])

		sell := order.NewOrder(com.Sell.OrderParity, com.Sell.OrderType, com.Sell.OrderExpiry, com.Sell.OrderSettlement, order.Tokens(values[8]), order.PriceFromCoExp(values[9], values[10]), order.VolumeFromCoExp(values[11], values[12]), order.VolumeFromCoExp(values[13], values[14]), values[15])
		sell.MinimumVolume = order.VolumeFromCoExp(values[13], values[14])

		settler.settleOrderMatch(com, buy, sell)
	}, true /* delay message sending to ensure the round-robin */)
	if err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot join buy = %v, sell = %v: %v", com.Buy.OrderID, com.Sell.OrderID, err))
	}
}

func (settler *settler) settleOrderMatch(com Computation, buy, sell order.Order) {
	// Submit a challenge if the orders do not match.
	if buy.Tokens != sell.Tokens ||
		buy.Volume < sell.MinimumVolume ||
		sell.Volume < buy.MinimumVolume ||
		buy.Price < sell.Price {
		if err := settler.contract.SubmitChallengeOrder(buy); err != nil {
			log.Printf("[error] (settle) cannot submit challenge for buy order = %v: %v", buy.ID, err)
		}
		if err := settler.contract.SubmitChallengeOrder(sell); err != nil {
			log.Printf("[error] (settle) cannot submit challenge for sell order = %v: %v", sell.ID, err)
		}
		if err := settler.contract.SubmitChallenge(buy.ID, sell.ID); err != nil {
			log.Printf("[error] (settle) cannot submit challenge buy = %v, sell = %v: %v", buy.ID, sell.ID, err)
		}
		log.Printf("[info] (slash) found mismatched order confirmation")
		return
	}

	// Leave the orders if volume is too low and there is no profit for
	// submitting such orders. Note: minimum volume is set to 1 ETH.
	settleVolume := volumeInEth(buy, sell)
	if settleVolume < settler.minimumSettleVolume {
		log.Printf("[info] (settle) cannot execute settlement buy = %v, sell = %v: volume = %v ETH too low", buy.ID, sell.ID, settleVolume)
		return
	}

	// Try settling the orders for at most 3 times.
	err := settler.contract.Settle(buy, sell)
	if err != nil {
		log.Printf("[error] (settle) cannot store settlement buy = %v, sell = %v: %v", buy.ID, sell.ID, err)
		return
	}

	com.State = ComputationStateSettled
	if err := settler.computationStore.PutComputation(com); err != nil {
		log.Printf("[error] (settle) cannot store settlement buy = %v, sell = %v: %v", buy.ID, sell.ID, err)
		return
	}
}

func volumeInEth(buy, sell order.Order) uint64 {
	if buy.Tokens.PriorityToken() == order.TokenETH {
		// BTC-ETH
		if buy.Volume >= sell.Volume {
			return sell.Volume
		} else {
			return buy.Volume
		}
	} else {
		// ETH-ERC20
		var erc20Volume uint64
		if buy.Volume >= sell.Volume {
			erc20Volume = sell.Volume
		} else {
			erc20Volume = buy.Volume
		}

		x := big.NewInt(0)
		y := big.NewInt(0)

		x.SetUint64(buy.Price)
		y.SetUint64(sell.Price)
		x.Add(x, y)

		y.SetUint64(2)
		x.Div(x, y)

		y.SetUint64(erc20Volume)
		x.Mul(x, y)

		y.SetUint64(1e12)
		x.Div(x, y)

		return x.Uint64()
	}
}
