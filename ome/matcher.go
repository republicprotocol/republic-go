package ome

import (
	"errors"
	"fmt"
	"time"

	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/smpc"
)

// ErrUnexpectedResolveStage is returned when a ResolveStage is not one of the
// explicitly enumerated values.
var ErrUnexpectedResolveStage = errors.New("unexpected resolve stage")

// ResolveStage defines the various stages that resolving can be in for any
// given Computation.
type ResolveStage byte

// Values for ResolveStage.
const (
	ResolveStageNil ResolveStage = iota
	ResolveStagePriceExp
	ResolveStagePriceCo
	ResolveStageBuyVolumeExp
	ResolveStageBuyVolumeCo
	ResolveStageSellVolumeExp
	ResolveStageSellVolumeCo
	ResolveStageTokens
)

// String returns the human-readable representation of a ResolveStage.
func (stage ResolveStage) String() string {
	switch stage {
	case ResolveStageNil:
		return "nil"
	case ResolveStagePriceExp:
		return "priceExp"
	case ResolveStagePriceCo:
		return "priceCo"
	case ResolveStageBuyVolumeExp:
		return "buyVolumeExp"
	case ResolveStageBuyVolumeCo:
		return "buyVolumeCo"
	case ResolveStageSellVolumeExp:
		return "sellVolumeExp"
	case ResolveStageSellVolumeCo:
		return "sellVolumeCo"
	case ResolveStageTokens:
		return "tokens"
	}
	return ""
}

// A MatchCallback is called when a Computation is finished. The Computation
// can then be inspected to determine if the result is a match.
type MatchCallback func(Computation)

// A Matcher resolves Computations into a matched, or mismatched, result.
type Matcher interface {

	// Resolve a Computation to determine whether or not the orders involved
	// are a match. The epoch hash of the Computation and is used to
	// differentiate between the various networks required for SMPC. The
	// MatchCallback is called when a result has be determined.
	Resolve(com Computation, buyFragment, sellFragment order.Fragment, callback MatchCallback)
}

type matcher struct {
	storer Storer
	smpcer smpc.Smpcer
}

// NewMatcher returns a Matcher that will resolve Computations by resolving
// each component in a pipeline. If a mismatch is encountered at any stage of
// the pipeline, the Computation is short circuited and the MatchCallback will
// be called immediately.
func NewMatcher(storer Storer, smpcer smpc.Smpcer) Matcher {
	return &matcher{
		storer: storer,
		smpcer: smpcer,
	}
}

// Resolve implements the Matcher interface.
func (matcher *matcher) Resolve(com Computation, buyFragment, sellFragment order.Fragment, callback MatchCallback) {
	if buyFragment.OrderSettlement != sellFragment.OrderSettlement {
		// Store the computation as a mismatch
		com.State = ComputationStateMismatched
		com.Match = false
		if err := matcher.storer.InsertComputation(com); err != nil {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot store mismatched computation buy = %v, sell = %v", com.Buy, com.Sell))
		}
		// Trigger the callback with a mismatch
		logger.Compute(logger.LevelDebug, fmt.Sprintf("✗ settlement => buy = %v, sell = %v", com.Buy, com.Sell))
		callback(com)
		return
	}

	matcher.resolve(smpc.NetworkID(com.EpochHash), com, buyFragment, sellFragment, callback, ResolveStagePriceExp)
}

func (matcher *matcher) resolve(networkID smpc.NetworkID, com Computation, buyFragment, sellFragment order.Fragment, callback MatchCallback, stage ResolveStage) {
	if isExpired(com, buyFragment, sellFragment) {
		com.State = ComputationStateRejected
		if err := matcher.storer.InsertComputation(com); err != nil {
			logger.Error(fmt.Sprintf("cannot store expired computation buy = %v, sell = %v: %v", com.Buy, com.Sell, err))
		}
		return
	}

	join, err := buildJoin(com, buyFragment, sellFragment, stage)
	if err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot build %v join: %v", stage, err))
		return
	}
	err = matcher.smpcer.Join(networkID, join, func(joinID smpc.JoinID, values []uint64) {
		matcher.resolveValues(values, networkID, com, buyFragment, sellFragment, callback, stage)
	})
	if err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot resolve %v: cannot join computation = %v: %v", stage, com.ID, err))
	}
}

func (matcher *matcher) resolveValues(values []uint64, networkID smpc.NetworkID, com Computation, buyFragment, sellFragment order.Fragment, callback MatchCallback, stage ResolveStage) {
	if len(values) != 1 {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot resolve %v: unexpected number of values: %v", stage, len(values)))
		return
	}

	switch stage {
	case ResolveStagePriceExp, ResolveStageBuyVolumeExp, ResolveStageSellVolumeExp:
		if isGreaterThanZero(values[0]) {
			logger.Compute(logger.LevelDebug, fmt.Sprintf("✔ %v => buy = %v, sell = %v", stage, com.Buy, com.Sell))
			matcher.resolve(networkID, com, buyFragment, sellFragment, callback, stage+2)
			return
		}
		if isEqualToZero(values[0]) {
			logger.Compute(logger.LevelDebug, fmt.Sprintf("✔ %v => buy = %v, sell = %v", stage, com.Buy, com.Sell))
			matcher.resolve(networkID, com, buyFragment, sellFragment, callback, stage+1)
			return
		}

	case ResolveStagePriceCo, ResolveStageBuyVolumeCo, ResolveStageSellVolumeCo:
		if isGreaterThanOrEqualToZero(values[0]) {
			logger.Compute(logger.LevelDebug, fmt.Sprintf("✔ %v => buy = %v, sell = %v", stage, com.Buy, com.Sell))
			matcher.resolve(networkID, com, buyFragment, sellFragment, callback, stage+1)
			return
		}

	case ResolveStageTokens:
		if isEqualToZero(values[0]) {
			// Store the computation as a match
			com.State = ComputationStateMatched
			com.Match = true
			if err := matcher.storer.InsertComputation(com); err != nil {
				logger.Compute(logger.LevelError, fmt.Sprintf("cannot store matched computation buy = %v, sell = %v", com.Buy, com.Sell))
			}

			// Trigger the callback with a match
			logger.Compute(logger.LevelDebug, fmt.Sprintf("✔ %v => buy = %v, sell = %v", stage, com.Buy, com.Sell))
			callback(com)
			return
		}

	default:
		// If the stage is unknown it is always considered a mismatch
	}

	// Store the computation as a mismatch
	com.State = ComputationStateMismatched
	com.Match = false
	if err := matcher.storer.InsertComputation(com); err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot store mismatched computation buy = %v, sell = %v", com.Buy, com.Sell))
	}

	// Trigger the callback with a mismatch
	logger.Compute(logger.LevelDebug, fmt.Sprintf("✗ %v => buy = %v, sell = %v", stage, com.Buy, com.Sell))
	callback(com)
}

func buildJoin(com Computation, buyFragment, sellFragment order.Fragment, stage ResolveStage) (smpc.Join, error) {
	var share shamir.Share
	switch stage {
	case ResolveStagePriceExp:
		share = buyFragment.Price.Exp.Sub(&sellFragment.Price.Exp)

	case ResolveStagePriceCo:
		share = buyFragment.Price.Co.Sub(&sellFragment.Price.Co)

	case ResolveStageBuyVolumeExp:
		share = buyFragment.Volume.Exp.Sub(&sellFragment.MinimumVolume.Exp)

	case ResolveStageBuyVolumeCo:
		share = buyFragment.Volume.Co.Sub(&sellFragment.MinimumVolume.Co)

	case ResolveStageSellVolumeExp:
		share = sellFragment.Volume.Exp.Sub(&buyFragment.MinimumVolume.Exp)

	case ResolveStageSellVolumeCo:
		share = sellFragment.Volume.Co.Sub(&buyFragment.MinimumVolume.Co)

	case ResolveStageTokens:
		share = buyFragment.Tokens.Sub(&sellFragment.Tokens)
	default:
		return smpc.Join{}, ErrUnexpectedResolveStage
	}
	join := smpc.Join{
		ID:     smpc.JoinID(com.ID),
		Index:  smpc.JoinIndex(share.Index),
		Shares: shamir.Shares{share},
	}
	join.ID[31] = byte(stage)
	return join, nil
}

func isGreaterThanOrEqualToZero(value uint64) bool {
	return value >= 0 && value < shamir.Prime/2
}

func isGreaterThanZero(value uint64) bool {
	return value > 0 && value < shamir.Prime/2
}

func isEqualToZero(value uint64) bool {
	return value == 0 || value == shamir.Prime
}

func isExpired(com Computation, buyFragment, sellFragment order.Fragment) bool {
	if time.Now().After(buyFragment.OrderExpiry) || time.Now().After(sellFragment.OrderExpiry) {
		logger.Compute(logger.LevelDebug, fmt.Sprintf("⧖ expired => buy = %v, sell = %v", com.Buy, com.Sell))
		return true
	}
	return false
}
