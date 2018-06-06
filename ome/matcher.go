package ome

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/smpc"
)

// A MatchCallback is called when a Computation is finished. The Computation
// can then be inspected to determine if the result is a match.
type MatchCallback func(Computation)

// A Matcher resolves Computations into a matched, or mismatched, result.
type Matcher interface {

	// Resolve a Computation to determine whether or not the orders involved
	// are a match. The ξ hash is used to define the ξ in which this
	// Computation exists, and the MatchCallback is called when a result has
	// be determined.
	Resolve(ξ [32]byte, com Computation, callback MatchCallback) error
}

type matcher struct {
	storer Storer
	smpcer smpc.Smpcer
}

// NewMatcher returns a Matcher that will resolve Computations by resolving
// each component in a pipeline. If a mismatch is encounterd at any stage of
// the pipeline, the Computation is short circuited and the MatchCallback will
// be called immediately.
func NewMatcher(storer Storer, smpcer smpc.Smpcer) Matcher {
	return &matcher{
		storer: storer,
		smpcer: smpcer,
	}
}

// Resolve implements the Matcher interface.
func (matcher *matcher) Resolve(ξ [32]byte, com Computation, callback MatchCallback) error {
	buyFragment, err := matcher.storer.OrderFragment(com.Buy)
	if err != nil {
		return err
	}
	sellFragment, err := matcher.storer.OrderFragment(com.Sell)
	if err != nil {
		return err
	}

	matcher.resolvePriceExp(smpc.NetworkID(ξ), buyFragment, sellFragment, com, callback)
	return nil
}

func (matcher *matcher) resolvePriceExp(networkID smpc.NetworkID, buyFragment, sellFragment order.Fragment, com Computation, callback MatchCallback) {
	priceExpShare := buyFragment.Price.Exp.Sub(&sellFragment.Price.Exp)
	priceExpJoin := smpc.Join{
		ID:     smpc.JoinID(com.ID),
		Index:  smpc.JoinIndex(priceExpShare.Index),
		Shares: shamir.Shares{priceExpShare},
	}

	err := matcher.smpcer.Join(networkID, priceExpJoin, func(joinID smpc.JoinID, values []uint64) {
		if len(values) != 1 {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot com priceExp: unexpected number of values: %v", len(values)))
			return
		}
		if isGreaterThanOrEqualToZero(values[0], com, "priceExp") {
			matcher.resolvePriceCo(networkID, buyFragment, sellFragment, com, callback)
			return
		}
		com.IsMatch = false
		callback(com)
	})
	if err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot join priceExp: %v", err))
	}
}

func (matcher *matcher) resolvePriceCo(networkID smpc.NetworkID, buyFragment, sellFragment order.Fragment, com Computation, callback MatchCallback) {
	priceCoShare := buyFragment.Price.Co.Sub(&sellFragment.Price.Co)
	priceCoJoin := smpc.Join{
		ID:     smpc.JoinID(com.ID),
		Index:  smpc.JoinIndex(priceCoShare.Index),
		Shares: shamir.Shares{priceCoShare},
	}

	err := matcher.smpcer.Join(networkID, priceCoJoin, func(joinID smpc.JoinID, values []uint64) {
		if len(values) != 1 {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot com priceCo: unexpected number of values: %v", len(values)))
			return
		}
		if isGreaterThanOrEqualToZero(values[0], com, "priceCo") {
			matcher.resolveBuyVolumeExp(networkID, buyFragment, sellFragment, com, callback)
			return
		}
		com.IsMatch = false
		callback(com)
	})
	if err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot join priceCo: %v", err))
	}
}

func (matcher *matcher) resolveBuyVolumeExp(networkID smpc.NetworkID, buyFragment, sellFragment order.Fragment, com Computation, callback MatchCallback) {
	buyVolumeExpShare := buyFragment.Volume.Co.Sub(&sellFragment.MinimumVolume.Exp)
	buyVolumeExpJoin := smpc.Join{
		ID:     smpc.JoinID(com.ID),
		Index:  smpc.JoinIndex(buyVolumeExpShare.Index),
		Shares: shamir.Shares{buyVolumeExpShare},
	}

	err := matcher.smpcer.Join(networkID, buyVolumeExpJoin, func(joinID smpc.JoinID, values []uint64) {
		if len(values) != 1 {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot com buyVolumeExp: unexpected number of values: %v", len(values)))
			return
		}
		if isGreaterThanOrEqualToZero(values[0], com, "buyVolumeExp") {
			matcher.resolveBuyVolumeCo(networkID, buyFragment, sellFragment, com, callback)
			return
		}
		com.IsMatch = false
		callback(com)
	})
	if err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot join buyVolumeExp: %v", err))
	}
}

func (matcher *matcher) resolveBuyVolumeCo(networkID smpc.NetworkID, buyFragment, sellFragment order.Fragment, com Computation, callback MatchCallback) {
	buyVolumeCoShare := buyFragment.Volume.Co.Sub(&sellFragment.MinimumVolume.Co)
	buyVolumeCoJoin := smpc.Join{
		ID:     smpc.JoinID(com.ID),
		Index:  smpc.JoinIndex(buyVolumeCoShare.Index),
		Shares: shamir.Shares{buyVolumeCoShare},
	}

	err := matcher.smpcer.Join(networkID, buyVolumeCoJoin, func(joinID smpc.JoinID, values []uint64) {
		if len(values) != 1 {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot com buyVolumeCo: unexpected number of values: %v", len(values)))
			return
		}
		if isGreaterThanOrEqualToZero(values[0], com, "buyVolumeCo") {
			matcher.resolveSellVolumeExp(networkID, buyFragment, sellFragment, com, callback)
			return
		}
		com.IsMatch = false
		callback(com)
	})
	if err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot join buyVolumeCo: %v", err))
	}
}

func (matcher *matcher) resolveSellVolumeExp(networkID smpc.NetworkID, buyFragment, sellFragment order.Fragment, com Computation, callback MatchCallback) {
	sellVolumeExpShare := sellFragment.Volume.Exp.Sub(&buyFragment.MinimumVolume.Exp)
	sellVolumeExpJoin := smpc.Join{
		ID:     smpc.JoinID(com.ID),
		Index:  smpc.JoinIndex(sellVolumeExpShare.Index),
		Shares: shamir.Shares{sellVolumeExpShare},
	}

	err := matcher.smpcer.Join(networkID, sellVolumeExpJoin, func(joinID smpc.JoinID, values []uint64) {
		if len(values) != 1 {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot com sellVolumeExp: unexpected number of values: %v", len(values)))
			return
		}
		if isGreaterThanOrEqualToZero(values[0], com, "sellVolumeExp") {
			matcher.resolveSellVolumeCo(networkID, sellFragment, sellFragment, com, callback)
			return
		}
		com.IsMatch = false
		callback(com)
	})
	if err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot join sellVolumeExp: %v", err))
	}
}

func (matcher *matcher) resolveSellVolumeCo(networkID smpc.NetworkID, buyFragment, sellFragment order.Fragment, com Computation, callback MatchCallback) {
	sellVolumeCoShare := sellFragment.Volume.Co.Sub(&buyFragment.MinimumVolume.Co)
	sellVolumeCoJoin := smpc.Join{
		ID:     smpc.JoinID(com.ID),
		Index:  smpc.JoinIndex(sellVolumeCoShare.Index),
		Shares: shamir.Shares{sellVolumeCoShare},
	}

	err := matcher.smpcer.Join(networkID, sellVolumeCoJoin, func(joinID smpc.JoinID, values []uint64) {
		if len(values) != 1 {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot com sellVolumeCo: unexpected number of values: %v", len(values)))
			return
		}
		if isGreaterThanOrEqualToZero(values[0], com, "sellVolumeCo") {
			matcher.resolveTokens(networkID, sellFragment, sellFragment, com, callback)
			return
		}
		com.IsMatch = false
		callback(com)
	})
	if err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot join sellVolumeCo: %v", err))
	}
}

func (matcher *matcher) resolveTokens(networkID smpc.NetworkID, buyFragment, sellFragment order.Fragment, com Computation, callback MatchCallback) {
	tokensShare := buyFragment.Tokens.Sub(&sellFragment.Tokens)
	tokensJoin := smpc.Join{
		ID:     smpc.JoinID(com.ID),
		Index:  smpc.JoinIndex(tokensShare.Index),
		Shares: shamir.Shares{tokensShare},
	}

	err := matcher.smpcer.Join(networkID, tokensJoin, func(joinID smpc.JoinID, values []uint64) {
		if len(values) != 1 {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot com tokens: unexpected number of values: %v", len(values)))
			return
		}
		if isEqualToZero(values[0], com, "tokens") {
			com.IsMatch = true
			callback(com)
			return
		}
		com.IsMatch = false
		callback(com)
	})
	if err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot join sellVolumeCo: %v", err))
	}
}

func isGreaterThanOrEqualToZero(value uint64, com Computation, stages ...string) bool {
	stage := ""
	if stages != nil && len(stages) > 0 {
		stage = "[" + strings.Join(stages, " => ") + "]"
	}
	if value > shamir.Prime/2 {
		logger.Compute(logger.LevelDebugHigh, fmt.Sprintf("✗ %v: mismatch: %v, buy = %v, sell = %v", stage, base64.StdEncoding.EncodeToString(com.Buy[:8]), base64.StdEncoding.EncodeToString(com.Sell[:8])))
		return false
	}
	logger.Compute(logger.LevelDebug, fmt.Sprintf("✔ %v: buy = %v, sell = %v", stage, base64.StdEncoding.EncodeToString(com.Buy[:8]), base64.StdEncoding.EncodeToString(com.Sell[:8])))
	return true
}

func isEqualToZero(value uint64, com Computation, stages ...string) bool {
	stage := ""
	if stages != nil && len(stages) > 0 {
		stage = "[" + strings.Join(stages, " => ") + "]"
	}
	if value == 0 || value == shamir.Prime {
		logger.Compute(logger.LevelDebugHigh, fmt.Sprintf("✗ %v: mismatch: %v, buy = %v, sell = %v", stage, base64.StdEncoding.EncodeToString(com.Buy[:8]), base64.StdEncoding.EncodeToString(com.Sell[:8])))
		return false
	}
	logger.Compute(logger.LevelDebug, fmt.Sprintf("✔ %v: buy = %v, sell = %v", stage, base64.StdEncoding.EncodeToString(com.Buy[:8]), base64.StdEncoding.EncodeToString(com.Sell[:8])))
	return true
}
