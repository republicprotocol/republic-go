package ome

import (
	"encoding/base64"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/smpc"
)

// Computations is an alias type.
type Computations []Computation

// A Computation involving a buy order, and a sell order, with a combined
// Priority.
type Computation struct {
	Buy      order.ID
	Sell     order.ID
	Priority Priority
}

// TODO: Stage bytes are a really ugly way of tracking our computations. We
// need a proper SMPC VM.
type Stage byte

const (
	StageCmpPriceExp   = 1
	StageCmpPriceCo    = 2
	StageCmpBuyVolExp  = 3
	StageCmpBuyVolCo   = 4
	StageCmpSellVolExp = 5
	StageCmpSellVolCo  = 6
	StageCmpTokens     = 7

	StageJoinBuyPriceExp   = 8
	StageJoinBuyPriceCo    = 9
	StageJoinBuyVolExp     = 10
	StageJoinBuyVolCo      = 11
	StageJoinBuyMinVolExp  = 12
	StageJoinBuyMinVolCo   = 13
	StageJoinBuyTokens     = 14
	StageJoinSellPriceExp  = 15
	StageJoinSellPriceCo   = 16
	StageJoinSellVolExp    = 17
	StageJoinSellVolCo     = 18
	StageJoinSellMinVolExp = 19
	StageJoinSellMinVolCo  = 20
	StageJoinSellTokens    = 21
)

type ComputationEpoch struct {
	Computation

	ID    [32]byte // Used for SMPC InstID
	Epoch [32]byte // Used for SMPC NetworkID
}

// Results is an alias type.
type Results []Result

// A Result is the output of performing Computations.
type Result struct {
	Buy     order.ID
	Sell    order.ID
	IsMatch bool
}

// A Computer consumes Computations that need to be processed and outputs
// Computations that result in an order match.
type Computer interface {

	// Compute order comparisons, component by component. If all components
	// match, then an order match has been found. Computations are consumed
	// from an input channel, and Computations that resulted in an order match
	// are produced to the output channel.
	// TODO: NetworkID should be constructed as part of the computation, and
	// there will need to be done by the orderbook.Server (which knows the
	// epoch at which an order fragment was received).
	Compute(done <-chan struct{}, computations <-chan ComputationEpoch) <-chan error
}

type computer struct {
	storer    orderbook.Storer
	smpcer    smpc.Smpcer
	confirmer Confirmer
	ledger    cal.RenLedger
	accounts  cal.DarkpoolAccounts

	computationsMu    *sync.Mutex
	computationsState map[[32]byte]ComputationEpoch
	computations      chan ComputationEpoch

	matchingComputationsMu    *sync.Mutex
	matchingComputationsState map[[32]byte]ComputationEpoch
	matchingComputations      chan Computation

	priceExpPointer  map[[32]byte]*uint64
	priceCoPointer   map[[32]byte]*uint64
	volExpPointer    map[[32]byte]*uint64
	volCoPointer     map[[32]byte]*uint64
	minVolExpPointer map[[32]byte]*uint64
	minVolCoPointer  map[[32]byte]*uint64
	tokensPointer    map[[32]byte]*order.Tokens
}

func NewComputer(storer orderbook.Storer, smpcer smpc.Smpcer, confirmer Confirmer, ledger cal.RenLedger, accounts cal.DarkpoolAccounts) Computer {
	return &computer{
		storer:    storer,
		smpcer:    smpcer,
		confirmer: confirmer,
		ledger:    ledger,
		accounts:  accounts,

		computationsMu:    new(sync.Mutex),
		computationsState: map[[32]byte]ComputationEpoch{},
		computations:      make(chan ComputationEpoch),

		matchingComputationsMu:    new(sync.Mutex),
		matchingComputationsState: map[[32]byte]ComputationEpoch{},
		matchingComputations:      make(chan Computation),

		priceExpPointer:  map[[32]byte]*uint64{},
		priceCoPointer:   map[[32]byte]*uint64{},
		volExpPointer:    map[[32]byte]*uint64{},
		volCoPointer:     map[[32]byte]*uint64{},
		minVolExpPointer: map[[32]byte]*uint64{},
		minVolCoPointer:  map[[32]byte]*uint64{},
		tokensPointer:    map[[32]byte]*order.Tokens{},
	}
}

func (computer *computer) Compute(done <-chan struct{}, computations <-chan ComputationEpoch) <-chan error {
	errs := make(chan error, 1)

	go func() {
		defer close(errs)

		instructions := computer.smpcer.Instructions()
		results := computer.smpcer.Results()

		computer.processComputations(done, instructions)
		computer.processResults(done, results)

		confirmedMatchingComputations, confirmedErrs := computer.confirmer.ConfirmOrderMatches(done, computer.matchingComputations)

		for {
			select {
			case <-done:
				return

			case computation, ok := <-computations:
				if !ok {
					return
				}
				computation.ID = computeID(computation.Computation)
				computation.ID[31] = StageCmpPriceExp

				computer.computationsMu.Lock()
				computer.computationsState[computation.ID] = computation
				computer.computationsMu.Unlock()

				select {
				case <-done:
				case computer.computations <- computation:
				}
			case confirmedMatchingComputation, ok := <-confirmedMatchingComputations:
				if !ok {
					return
				}
				confirmedMatchingComputationID := computeID(confirmedMatchingComputation)
				computer.matchingComputationsMu.Lock()
				computation := computer.matchingComputationsState[confirmedMatchingComputationID]
				delete(computer.matchingComputationsState, confirmedMatchingComputationID)
				computer.matchingComputationsMu.Unlock()
				computation.ID[31] = StageJoinBuyPriceExp

				computer.computationsMu.Lock()
				computer.computationsState[computation.ID] = computation
				computer.computationsMu.Unlock()

				select {
				case <-done:
				case computer.computations <- computation:
				}

			case err, ok := <-confirmedErrs:
				if !ok {
					return
				}
				select {
				case <-done:
				case errs <- err:
				}
			}
		}
	}()

	return errs
}

func (computer *computer) processComputations(done <-chan struct{}, insts chan<- smpc.Inst) {
	go func() {

		pendingComputations := map[[32]byte]ComputationEpoch{}
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case computation, ok := <-computer.computations:
				if !ok {
					return
				}
				computer.processComputation(computation, pendingComputations, done, insts)
			case <-ticker.C:
				if len(pendingComputations) == 0 {
					continue
				}
				for _, computation := range pendingComputations {
					computer.processComputation(computation, pendingComputations, done, insts)
				}
			}
		}
	}()
}

func (computer *computer) processComputation(computation ComputationEpoch, pendingComputations map[[32]byte]ComputationEpoch, done <-chan struct{}, insts chan<- smpc.Inst) {

	buy, err := computer.storer.OrderFragment(computation.Buy)
	if err != nil {
		pendingComputations[computation.ID] = computation
		return
	}
	sell, err := computer.storer.OrderFragment(computation.Sell)
	if err != nil {
		pendingComputations[computation.ID] = computation
		return
	}

	buyID := base64.StdEncoding.EncodeToString(computation.Buy[:])
	if strings.HasPrefix(buyID, "dcF") {
		log.Println(buy)
	}
	sellID := base64.StdEncoding.EncodeToString(computation.Sell[:])
	if strings.HasPrefix(sellID, "dcF") {
		log.Println(sell)
	}

	delete(pendingComputations, computation.ID)
	var share shamir.Share
	switch computation.ID[31] {
	case StageCmpPriceExp:
		share = buy.Price.Exp.Sub(&sell.Price.Exp)
	case StageCmpPriceCo:
		share = buy.Price.Co.Sub(&sell.Price.Co)
	case StageCmpBuyVolExp:
		share = buy.Volume.Exp.Sub(&sell.MinimumVolume.Exp)
	case StageCmpBuyVolCo:
		share = buy.Volume.Co.Sub(&sell.MinimumVolume.Co)
	case StageCmpSellVolExp:
		share = sell.Volume.Exp.Sub(&buy.MinimumVolume.Exp)
	case StageCmpSellVolCo:
		share = sell.Volume.Co.Sub(&buy.MinimumVolume.Co)
	case StageCmpTokens:
		share = buy.Tokens.Sub(&sell.Tokens)

	case StageJoinBuyPriceExp:
		share = buy.Price.Exp
	case StageJoinBuyPriceCo:
		share = buy.Price.Co
	case StageJoinBuyVolExp:
		share = buy.Volume.Exp
	case StageJoinBuyVolCo:
		share = buy.Volume.Co
	case StageJoinBuyMinVolExp:
		share = buy.MinimumVolume.Exp
	case StageJoinBuyMinVolCo:
		share = buy.MinimumVolume.Co
	case StageJoinBuyTokens:
		share = buy.Tokens
	case StageJoinSellPriceExp:
		share = sell.Price.Exp
	case StageJoinSellPriceCo:
		share = sell.Price.Co
	case StageJoinSellVolExp:
		share = sell.Volume.Exp
	case StageJoinSellVolCo:
		share = sell.Volume.Co
	case StageJoinSellMinVolExp:
		share = sell.MinimumVolume.Exp
	case StageJoinSellMinVolCo:
		share = sell.MinimumVolume.Co
	case StageJoinSellTokens:
		share = sell.Tokens
	}

	inst := smpc.Inst{
		InstID:    computation.ID,
		NetworkID: computation.Epoch,
		InstJ: &smpc.InstJ{
			Share: share,
		},
	}

	log.Printf("[stage => %v] processing computation: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))

	// Write instruction
	select {
	case insts <- inst:
	default:
	}

	// Invert flow if necessary
	select {
	case <-done:
	case insts <- inst:
	case computation, ok := <-computer.computations:
		log.Printf("compute inversion: buy = %v; sell = %v", base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		if !ok {
			break
		}
		pendingComputations[computation.ID] = computation
	}
}

func (computer *computer) processResults(done <-chan struct{}, results <-chan smpc.Result) {
	go func() {

		for {
			select {
			case <-done:
				return
			case result, ok := <-results:
				if !ok {
					return
				}
				if result.ResultJ != nil {
					computer.processResultJ(result.InstID, result.NetworkID, *result.ResultJ, done)
				}
			}
		}

	}()
}

func (computer *computer) processResultJ(instID, networkID [32]byte, resultJ smpc.ResultJ, done <-chan struct{}) {
	half := shamir.Prime / 2
	computer.computationsMu.Lock()
	computation, ok := computer.computationsState[instID]
	delete(computer.computationsState, instID)
	computer.computationsMu.Unlock()
	if !ok {
		return
	}

	log.Printf("[stage => %v] received result: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:]), base64.StdEncoding.EncodeToString(computation.Sell[:]))

	switch instID[31] {

	case StageCmpPriceExp:
		if resultJ.Value > half {
			log.Printf("[stage => %v] halt: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
			return
		}
		log.Printf("[stage => %v] ok: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		computation.ID[31] = StageCmpPriceCo

	case StageCmpPriceCo:
		if resultJ.Value > half {
			log.Printf("[stage => %v] halt: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
			return
		}
		log.Printf("[stage => %v] ok: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		computation.ID[31] = StageCmpBuyVolExp

	case StageCmpBuyVolExp:
		if resultJ.Value > half {
			log.Printf("[stage => %v] halt: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
			return
		}
		log.Printf("[stage => %v] ok: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		computation.ID[31] = StageCmpBuyVolCo

	case StageCmpBuyVolCo:
		if resultJ.Value > half {
			log.Printf("[stage => %v] halt: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
			return
		}
		log.Printf("[stage => %v] ok: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		computation.ID[31] = StageCmpSellVolExp

	case StageCmpSellVolExp:
		if resultJ.Value > half {
			log.Printf("[stage => %v] halt: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
			return
		}
		log.Printf("[stage => %v] ok: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		computation.ID[31] = StageCmpSellVolCo

	case StageCmpSellVolCo:
		if resultJ.Value > half {
			log.Printf("[stage => %v] halt: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
			return
		}
		log.Printf("[stage => %v] ok: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		computation.ID[31] = StageCmpTokens

	case StageCmpTokens:
		if resultJ.Value == 0 {
			computation.ID = computeID(computation.Computation)
			computer.matchingComputationsMu.Lock()
			computer.matchingComputationsState[computation.ID] = computation
			computer.matchingComputationsMu.Unlock()
			select {
			case <-done:
			case computer.matchingComputations <- computation.Computation:
				log.Printf("✔ [stage => %v] matched: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
			}
		} else {
			log.Printf("[stage => %v] halt: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		}
		return

	case StageJoinBuyPriceExp:
		computer.priceExpPointer[computation.Buy] = &resultJ.Value
		computation.ID[31] = StageJoinBuyPriceCo
	case StageJoinBuyPriceCo:
		computer.priceCoPointer[computation.Buy] = &resultJ.Value
		computation.ID[31] = StageJoinBuyVolExp
	case StageJoinBuyVolExp:
		computer.volExpPointer[computation.Buy] = &resultJ.Value
		computation.ID[31] = StageJoinBuyVolCo
	case StageJoinBuyVolCo:
		computer.volCoPointer[computation.Buy] = &resultJ.Value
		computation.ID[31] = StageJoinBuyMinVolExp
	case StageJoinBuyMinVolExp:
		computer.minVolExpPointer[computation.Buy] = &resultJ.Value
		computation.ID[31] = StageJoinBuyMinVolCo
	case StageJoinBuyMinVolCo:
		computer.minVolCoPointer[computation.Buy] = &resultJ.Value
		computation.ID[31] = StageJoinBuyTokens
	case StageJoinBuyTokens:
		tokens := order.Tokens(resultJ.Value)
		computer.tokensPointer[computation.Buy] = &tokens
		computation.ID[31] = StageJoinSellPriceExp

	case StageJoinSellPriceExp:
		computer.priceExpPointer[computation.Sell] = &resultJ.Value
		computation.ID[31] = StageJoinSellPriceCo
	case StageJoinSellPriceCo:
		computer.priceCoPointer[computation.Sell] = &resultJ.Value
		computation.ID[31] = StageJoinSellVolExp
	case StageJoinSellVolExp:
		computer.volExpPointer[computation.Sell] = &resultJ.Value
		computation.ID[31] = StageJoinSellVolCo
	case StageJoinSellVolCo:
		computer.volCoPointer[computation.Sell] = &resultJ.Value
		computation.ID[31] = StageJoinSellMinVolExp
	case StageJoinSellMinVolExp:
		computer.minVolExpPointer[computation.Sell] = &resultJ.Value
		computation.ID[31] = StageJoinSellMinVolCo
	case StageJoinSellMinVolCo:
		computer.minVolCoPointer[computation.Sell] = &resultJ.Value
		computation.ID[31] = StageJoinSellTokens
	case StageJoinSellTokens:
		tokens := order.Tokens(resultJ.Value)
		computer.tokensPointer[computation.Sell] = &tokens

		// TODO:
		// 1. Settle buy order and sell order
		// 2. Delete orders from pointer maps
		buy, err := computer.reconstructOrder(computation.Buy)
		if err != nil {
			log.Printf("cannot reconstruct buy order %v : %v", base64.StdEncoding.EncodeToString(computation.Buy[:]), err)
			return
		}
		sell, err := computer.reconstructOrder(computation.Sell)
		if err != nil {
			log.Printf("cannot reconstruct sell order %v : %v", base64.StdEncoding.EncodeToString(computation.Sell[:]), err)
			return
		}
		err = computer.accounts.Settle(buy, sell)
		if err != nil {
			log.Printf("cannot settle orders: %v", err)
			return
		}
		log.Printf("-------Order Settled (Buy:%v , Sell: %v)-------", base64.StdEncoding.EncodeToString(computation.Buy[:]), base64.StdEncoding.EncodeToString(computation.Sell[:]))
		return
	}

	computer.computationsMu.Lock()
	computer.computationsState[computation.ID] = computation
	computer.computationsMu.Unlock()
	select {
	case <-done:
	case computer.computations <- computation:
		// FIXME: Removing the buffered channel from the mockSMPC causes a
		// deadlock. This is because there is a cyclic hold-and-wait between
		// reading a result, and writing an instruction. This can be solved
		// using very careful prioritization of channels (but first, we will
		// need to create more channels and merge them cleverly). It will also
		// involve fine-tuned buffering of the channels.
	}
}

func computeID(computation Computation) [32]byte {
	id := [32]byte{}
	copy(id[:], crypto.Keccak256(computation.Buy[:], computation.Sell[:]))
	return id
}

func (computer *computer) reconstructOrder(id order.ID) (order.Order, error) {
	fragment, err := computer.storer.OrderFragment(id)
	if err != nil {
		return order.Order{}, err
	}
	price := order.CoExp{
		Co:  *computer.priceCoPointer[id],
		Exp: *computer.priceExpPointer[id],
	}
	volume := order.CoExp{
		Co:  *computer.volCoPointer[id],
		Exp: *computer.volExpPointer[id],
	}
	minVolume := order.CoExp{
		Co:  *computer.minVolCoPointer[id],
		Exp: *computer.minVolExpPointer[id],
	}

	return order.NewOrder(order.TypeLimit, order.ParityBuy, fragment.OrderExpiry, *computer.tokensPointer[fragment.OrderID], price, volume, minVolume, 1), nil
}
