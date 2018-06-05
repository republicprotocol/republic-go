package ome

import (
	"encoding/base64"
	"log"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/smpc"
)

type ComputationResult byte

const (
	ComputationResultMatched         = 1
	ComputationResultMismatched      = 2
	ComputationResultConfirmAccepted = 3
	ComputationResultConfirmRejected = 4
	ComputationResultSettled         = 5
)

// Computations is an alias type.
type Computations []Computation

// A Computation involving a buy order, and a sell order, with a combined
// Priority.
type Computation struct {
	Buy      order.ID          `json:"buy"`
	Sell     order.ID          `json:"sell"`
	Priority Priority          `json:"priority"`
	Result   ComputationResult `json:"result"`
}

func (computation Computation) ID() [32]byte {
	id := [32]byte{}
	copy(id[:], crypto.Keccak256(computation.Buy[:], computation.Sell[:]))
	return id
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
	Computation `json:"computation"`

	ID    [32]byte `json:"id"`    // Used for SMPC InstID
	Epoch [32]byte `json:"epoch"` // Used for SMPC NetworkID
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
	storer    Storer
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

func NewComputer(storer Storer, smpcer smpc.Smpcer, confirmer Confirmer, ledger cal.RenLedger, accounts cal.DarkpoolAccounts) Computer {
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

				storedComputation, err := computer.storer.Computation(computation.ID)
				if err == nil {
					switch storedComputation.Result {
					case ComputationResultMatched:
						continue // FIXME: Skip matching and send to the confirmer
					case ComputationResultMismatched:
						continue
					case ComputationResultConfirmAccepted:
						computation.ID[31] = StageJoinBuyPriceExp
					case ComputationResultConfirmRejected:
						continue
					case ComputationResultSettled:
						continue
					}
				} else {
					computation.ID[31] = StageCmpPriceExp
				}

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
				computation, ok := computer.matchingComputationsState[confirmedMatchingComputationID]
				if !ok {
					log.Println("COMPUTATION CONFIRMED doesn't exist in the state map")
					continue
				}
				log.Println("COMPUTATION CONFIRMED:", computation)
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
			err := storeComputationResult(computer.storer, computation, ComputationResultMismatched)
			if err != nil {
				log.Printf("fail to store the computaion result: %v", err)
			}
			return
		}
		log.Printf("[stage => %v] ok: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		computation.ID[31] = StageCmpPriceCo

	case StageCmpPriceCo:
		if resultJ.Value > half {
			log.Printf("[stage => %v] halt: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
			err := storeComputationResult(computer.storer, computation, ComputationResultMismatched)
			if err != nil {
				log.Printf("fail to store the computaion result: %v", err)
			}
			return
		}
		log.Printf("[stage => %v] ok: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		computation.ID[31] = StageCmpBuyVolExp

	case StageCmpBuyVolExp:
		if resultJ.Value > half {
			log.Printf("[stage => %v] halt: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
			err := storeComputationResult(computer.storer, computation, ComputationResultMismatched)
			if err != nil {
				log.Printf("fail to store the computaion result: %v", err)
			}
			return
		}
		log.Printf("[stage => %v] ok: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		computation.ID[31] = StageCmpBuyVolCo

	case StageCmpBuyVolCo:
		if resultJ.Value > half {
			log.Printf("[stage => %v] halt: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
			err := storeComputationResult(computer.storer, computation, ComputationResultMismatched)
			if err != nil {
				log.Printf("fail to store the computaion result: %v", err)
			}
			return
		}
		log.Printf("[stage => %v] ok: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		computation.ID[31] = StageCmpSellVolExp

	case StageCmpSellVolExp:
		if resultJ.Value > half {
			log.Printf("[stage => %v] halt: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
			err := storeComputationResult(computer.storer, computation, ComputationResultMismatched)
			if err != nil {
				log.Printf("fail to store the computaion result: %v", err)
			}
			return
		}
		log.Printf("[stage => %v] ok: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		computation.ID[31] = StageCmpSellVolCo

	case StageCmpSellVolCo:
		if resultJ.Value > half {
			log.Printf("[stage => %v] halt: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
			err := storeComputationResult(computer.storer, computation, ComputationResultMismatched)
			if err != nil {
				log.Printf("fail to store the computaion result: %v", err)
			}
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

			err := storeComputationResult(computer.storer, computation, ComputationResultMatched)
			if err != nil {
				log.Printf("fail to store the computaion result: %v", err)
				return
			}
			select {
			case <-done:
			case computer.matchingComputations <- computation.Computation:
				log.Printf("✔ [stage => matched] buy = %v; sell = %v", base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
			}
			return
		}
		err := storeComputationResult(computer.storer, computation, ComputationResultMismatched)
		if err != nil {
			log.Printf("fail to store the computaion result: %v", err)
		}
		log.Printf("[stage => %v] halt: buy = %v; sell = %v", StageCmpTokens, base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		return

	case StageJoinBuyPriceExp:
		log.Printf("[stage => %v] join buy price exp : result = %v", computation.ID[31], resultJ.Value)
		computer.priceExpPointer[computation.Buy] = &resultJ.Value
		computation.ID[31] = StageJoinBuyPriceCo
	case StageJoinBuyPriceCo:
		log.Printf("[stage => %v] join buy price co : result = %v", computation.ID[31], resultJ.Value)
		computer.priceCoPointer[computation.Buy] = &resultJ.Value
		computation.ID[31] = StageJoinBuyVolExp
	case StageJoinBuyVolExp:
		log.Printf("[stage => %v] join buy vol exp: result = %v", computation.ID[31], resultJ.Value)
		computer.volExpPointer[computation.Buy] = &resultJ.Value
		computation.ID[31] = StageJoinBuyVolCo
	case StageJoinBuyVolCo:
		log.Printf("[stage => %v] join buy vol co: result = %v", computation.ID[31], resultJ.Value)
		computer.volCoPointer[computation.Buy] = &resultJ.Value
		computation.ID[31] = StageJoinBuyMinVolExp
	case StageJoinBuyMinVolExp:
		log.Printf("[stage => %v] join buy min vol exp: result = %v", computation.ID[31], resultJ.Value)
		computer.minVolExpPointer[computation.Buy] = &resultJ.Value
		computation.ID[31] = StageJoinBuyMinVolCo
	case StageJoinBuyMinVolCo:
		log.Printf("[stage => %v] join buy min vol co: result = %v", computation.ID[31], resultJ.Value)
		computer.minVolCoPointer[computation.Buy] = &resultJ.Value
		computation.ID[31] = StageJoinBuyTokens
	case StageJoinBuyTokens:
		log.Printf("[stage => %v] join buy token: result = %v", computation.ID[31], resultJ.Value)
		tokens := order.Tokens(resultJ.Value)
		computer.tokensPointer[computation.Buy] = &tokens
		computation.ID[31] = StageJoinSellPriceExp
	case StageJoinSellPriceExp:
		log.Printf("[stage => %v] join sell price exp: result = %v", computation.ID[31], resultJ.Value)
		computer.priceExpPointer[computation.Sell] = &resultJ.Value
		computation.ID[31] = StageJoinSellPriceCo
	case StageJoinSellPriceCo:
		log.Printf("[stage => %v] join sell price co: result = %v", computation.ID[31], resultJ.Value)
		computer.priceCoPointer[computation.Sell] = &resultJ.Value
		computation.ID[31] = StageJoinSellVolExp
	case StageJoinSellVolExp:
		log.Printf("[stage => %v] join sell vol exp: result = %v", computation.ID[31], resultJ.Value)
		computer.volExpPointer[computation.Sell] = &resultJ.Value
		computation.ID[31] = StageJoinSellVolCo
	case StageJoinSellVolCo:
		log.Printf("[stage => %v] join sell vol co: result = %v", computation.ID[31], resultJ.Value)
		computer.volCoPointer[computation.Sell] = &resultJ.Value
		computation.ID[31] = StageJoinSellMinVolExp
	case StageJoinSellMinVolExp:
		log.Printf("[stage => %v] join sell min vol exp: result = %v", computation.ID[31], resultJ.Value)
		computer.minVolExpPointer[computation.Sell] = &resultJ.Value
		computation.ID[31] = StageJoinSellMinVolCo
	case StageJoinSellMinVolCo:
		log.Printf("[stage => %v] join sell min vol co: result = %v", computation.ID[31], resultJ.Value)
		computer.minVolCoPointer[computation.Sell] = &resultJ.Value
		computation.ID[31] = StageJoinSellTokens
	case StageJoinSellTokens:
		log.Printf("[stage => %v] join sell token: result = %v", computation.ID[31], resultJ.Value)
		tokens := order.Tokens(resultJ.Value)
		computer.tokensPointer[computation.Sell] = &tokens
		log.Printf(" start settl orders (Buy:%v , Sell: %v)-------", base64.StdEncoding.EncodeToString(computation.Buy[:]), base64.StdEncoding.EncodeToString(computation.Sell[:]))

		// TODO:
		// 1. Settle buy order and sell order
		// 2. Delete orders from pointer maps
		buy, err := computer.reconstructOrder(computation.Buy)
		log.Printf("<buy order> price<Co: %v, Exp: %v>, volumn<Co: %v, Exp: %v>,", buy.Price.Co, buy.Price.Exp, buy.Volume.Co, buy.Volume.Exp)
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
	if computer.priceCoPointer[id] == nil {
		log.Println("price co pointer nil ")
	}
	if computer.priceExpPointer[id] == nil {
		log.Println("price exp pointer nil ")
	}
	if computer.volCoPointer[id] == nil {
		log.Println("vol co pointer nil ")
	}
	if computer.volExpPointer[id] == nil {
		log.Println("vol exp pointer nil ")
	}
	if computer.minVolCoPointer[id] == nil {
		log.Println("min vol cp pointer nil ")
	}
	if computer.minVolExpPointer[id] == nil {
		log.Println("min vol exp pointer nil ")
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

	ord := order.NewOrder(order.TypeLimit, order.ParityBuy, fragment.OrderExpiry, *computer.tokensPointer[fragment.OrderID], price, volume, minVolume, 1)
	ord.ID = id
	return ord, nil
}

func storeComputationResult(store Storer, computation ComputationEpoch, result ComputationResult) error {
	computationResult := Computation{
		Buy:      computation.Buy,
		Sell:     computation.Sell,
		Priority: computation.Priority,
		Result:   result,
	}

	return store.InsertComputation(computationResult)
}
