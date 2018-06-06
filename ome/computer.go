package ome

import (
	"encoding/base64"
	"log"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
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
	StageNil = iota
	StageCmpPriceExp
	StageCmpPriceCo
	StageCmpBuyVolExp
	StageCmpBuyVolCo
	StageCmpSellVolExp
	StageCmpSellVolCo
	StageCmpTokens

	StageJoin

	StageJoinBuyPriceExp
	StageJoinBuyPriceCo
	StageJoinBuyVolExp
	StageJoinBuyVolCo
	StageJoinBuyMinVolExp
	StageJoinBuyMinVolCo
	StageJoinBuyTokens
	StageJoinSellPriceExp
	StageJoinSellPriceCo
	StageJoinSellVolExp
	StageJoinSellVolCo
	StageJoinSellMinVolExp
	StageJoinSellMinVolCo
	StageJoinSellTokens
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

		computer.processComputations(done)

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

				// storedComputation, err := computer.storer.Computation(computation.ID)
				// if err == nil {
				// 	switch storedComputation.Result {
				// 	case ComputationResultMatched:
				// 		continue // FIXME: Skip matching and send to the confirmer
				// 	case ComputationResultMismatched:
				// 		continue
				// 	case ComputationResultConfirmAccepted:
				// 		computation.ID[31] = StageJoin
				// 	case ComputationResultConfirmRejected:
				// 		continue
				// 	case ComputationResultSettled:
				// 		continue
				// 	}
				// } else {
				// 	computation.ID[31] = StageCmpPriceExp
				// }
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
				computation, ok := computer.matchingComputationsState[confirmedMatchingComputationID]
				if !ok {
					log.Println("COMPUTATION CONFIRMED doesn't exist in the state map")
					continue
				}
				log.Println("COMPUTATION CONFIRMED:", computation)
				delete(computer.matchingComputationsState, confirmedMatchingComputationID)
				computer.matchingComputationsMu.Unlock()

				computer.computationsMu.Lock()
				for i := StageJoinBuyPriceExp; i <= StageJoinSellTokens; i++ {
					computation.ID[31] = byte(i)
					computer.computationsState[computation.ID] = computation
				}
				computer.computationsMu.Unlock()

				computation.ID[31] = StageJoin
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

func (computer *computer) processComputations(done <-chan struct{}) {
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
				computer.processComputation(computation, pendingComputations)
			case <-ticker.C:
				log.Printf("there are %d computations in the pending list ", len(pendingComputations))
				if len(pendingComputations) == 0 {
					continue
				}
				for _, computation := range pendingComputations {
					computer.processComputation(computation, pendingComputations)
				}
			}
		}
	}()
}

func (computer *computer) processComputation(computation ComputationEpoch, pendingComputations map[[32]byte]ComputationEpoch) {
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

	var components smpc.Components
	switch computation.ID[31] {
	case StageCmpPriceExp:
		components = smpc.Components{smpc.Component{ComponentID: smpc.ComponentID(computation.ID), Share: buy.Price.Exp.Sub(&sell.Price.Exp)}}
	case StageCmpPriceCo:
		components = smpc.Components{smpc.Component{ComponentID: smpc.ComponentID(computation.ID), Share: buy.Price.Co.Sub(&sell.Price.Co)}}
	case StageCmpBuyVolExp:
		components = smpc.Components{smpc.Component{ComponentID: smpc.ComponentID(computation.ID), Share: buy.Volume.Exp.Sub(&sell.MinimumVolume.Exp)}}
	case StageCmpBuyVolCo:
		components = smpc.Components{smpc.Component{ComponentID: smpc.ComponentID(computation.ID), Share: buy.Volume.Co.Sub(&sell.MinimumVolume.Co)}}
	case StageCmpSellVolExp:
		components = smpc.Components{smpc.Component{ComponentID: smpc.ComponentID(computation.ID), Share: sell.Volume.Exp.Sub(&buy.MinimumVolume.Exp)}}
	case StageCmpSellVolCo:
		components = smpc.Components{smpc.Component{ComponentID: smpc.ComponentID(computation.ID), Share: sell.Volume.Co.Sub(&buy.MinimumVolume.Co)}}
	case StageCmpTokens:
		components = smpc.Components{smpc.Component{ComponentID: smpc.ComponentID(computation.ID), Share: buy.Tokens.Sub(&sell.Tokens)}}
	case StageJoin:
		log.Println("JOIN COMPUTATION")
		components = smpc.Components{
			smpc.Component{ComponentID: computation.ID, Share: buy.Price.Exp},
			smpc.Component{ComponentID: computation.ID, Share: buy.Price.Co},
			smpc.Component{ComponentID: computation.ID, Share: buy.Volume.Exp},
			smpc.Component{ComponentID: computation.ID, Share: buy.Volume.Co},
			smpc.Component{ComponentID: computation.ID, Share: buy.MinimumVolume.Exp},
			smpc.Component{ComponentID: computation.ID, Share: buy.MinimumVolume.Co},
			smpc.Component{ComponentID: computation.ID, Share: buy.Tokens},

			smpc.Component{ComponentID: computation.ID, Share: sell.Price.Exp},
			smpc.Component{ComponentID: computation.ID, Share: sell.Price.Co},
			smpc.Component{ComponentID: computation.ID, Share: sell.Volume.Exp},
			smpc.Component{ComponentID: computation.ID, Share: sell.Volume.Co},
			smpc.Component{ComponentID: computation.ID, Share: sell.MinimumVolume.Exp},
			smpc.Component{ComponentID: computation.ID, Share: sell.MinimumVolume.Co},
			smpc.Component{ComponentID: computation.ID, Share: sell.Tokens},
		}
		components[0].ComponentID[31] = StageJoinBuyPriceExp
		components[1].ComponentID[31] = StageJoinBuyPriceCo
		components[2].ComponentID[31] = StageJoinBuyVolExp
		components[3].ComponentID[31] = StageJoinBuyVolCo
		components[4].ComponentID[31] = StageJoinBuyMinVolExp
		components[5].ComponentID[31] = StageJoinBuyMinVolCo
		components[6].ComponentID[31] = StageJoinBuyTokens

		components[7].ComponentID[31] = StageJoinSellPriceExp
		components[8].ComponentID[31] = StageJoinSellPriceCo
		components[9].ComponentID[31] = StageJoinSellVolExp
		components[10].ComponentID[31] = StageJoinSellVolCo
		components[11].ComponentID[31] = StageJoinSellMinVolExp
		components[12].ComponentID[31] = StageJoinSellMinVolCo
		components[13].ComponentID[31] = StageJoinSellTokens
	}

	computer.smpcer.JoinComponents(computation.Epoch, components, computer)
}

func (computer *computer) OnNotifyBuild(componentID smpc.ComponentID, networkID smpc.NetworkID, value uint64) {
	half := shamir.Prime / 2
	computer.computationsMu.Lock()
	computation, ok := computer.computationsState[componentID]
	delete(computer.computationsState, componentID)
	computer.computationsMu.Unlock()
	if !ok {
		return
	}

	log.Printf("[stage => %v] received result: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:]), base64.StdEncoding.EncodeToString(computation.Sell[:]))

	switch componentID[31] {

	case StageCmpPriceExp:
		if value > half {
			//log.Printf("[stage => %v] halt: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
			err := storeComputationResult(computer.storer, computation, ComputationResultMismatched)
			if err != nil {
				log.Printf("fail to store the computaion result: %v", err)
			}
			return
		}
		//log.Printf("[stage => %v] ok: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		computation.ID[31] = StageCmpPriceCo

	case StageCmpPriceCo:
		if value > half {
			//log.Printf("[stage => %v] halt: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
			err := storeComputationResult(computer.storer, computation, ComputationResultMismatched)
			if err != nil {
				log.Printf("fail to store the computaion result: %v", err)
			}
			return
		}
		//log.Printf("[stage => %v] ok: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		computation.ID[31] = StageCmpBuyVolExp

	case StageCmpBuyVolExp:
		if value > half {
			//log.Printf("[stage => %v] halt: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
			err := storeComputationResult(computer.storer, computation, ComputationResultMismatched)
			if err != nil {
				log.Printf("fail to store the computaion result: %v", err)
			}
			return
		}
		//log.Printf("[stage => %v] ok: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		computation.ID[31] = StageCmpBuyVolCo

	case StageCmpBuyVolCo:
		if value > half {
			//log.Printf("[stage => %v] halt: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
			err := storeComputationResult(computer.storer, computation, ComputationResultMismatched)
			if err != nil {
				log.Printf("fail to store the computaion result: %v", err)
			}
			return
		}
		//log.Printf("[stage => %v] ok: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		computation.ID[31] = StageCmpSellVolExp

	case StageCmpSellVolExp:
		if value > half {
			//log.Printf("[stage => %v] halt: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
			err := storeComputationResult(computer.storer, computation, ComputationResultMismatched)
			if err != nil {
				log.Printf("fail to store the computaion result: %v", err)
			}
			return
		}
		//log.Printf("[stage => %v] ok: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		computation.ID[31] = StageCmpSellVolCo

	case StageCmpSellVolCo:
		if value > half {
			//log.Printf("[stage => %v] halt: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
			err := storeComputationResult(computer.storer, computation, ComputationResultMismatched)
			if err != nil {
				log.Printf("fail to store the computaion result: %v", err)
			}
			return
		}
		//log.Printf("[stage => %v] ok: buy = %v; sell = %v", computation.ID[31], base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		computation.ID[31] = StageCmpTokens

	case StageCmpTokens:
		if value == 0 {
			computation.ID = computeID(computation.Computation)
			computer.matchingComputationsMu.Lock()
			computer.matchingComputationsState[computation.ID] = computation
			computer.matchingComputationsMu.Unlock()

			err := storeComputationResult(computer.storer, computation, ComputationResultMatched)
			if err != nil {
				log.Printf("fail to store the computaion result: %v", err)
				return
			}

			// FIXME: This cannot be escaped when the done channel is closed
			computer.matchingComputations <- computation.Computation
			log.Printf("✔ [stage => matched] buy = %v; sell = %v", base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
			return
		}
		err := storeComputationResult(computer.storer, computation, ComputationResultMismatched)
		if err != nil {
			log.Printf("fail to store the computaion result: %v", err)
		}
		log.Printf("[stage => %v] halt: buy = %v; sell = %v", StageCmpTokens, base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
		return

	case StageJoinBuyPriceExp:
		log.Printf("[stage => %v] join buy price exp : result = %v", computation.ID[31], value)
		computer.priceExpPointer[computation.Buy] = &value
		computation.ID[31] = StageJoinBuyPriceCo
	case StageJoinBuyPriceCo:
		log.Printf("[stage => %v] join buy price co : result = %v", computation.ID[31], value)
		computer.priceCoPointer[computation.Buy] = &value
		computation.ID[31] = StageJoinBuyVolExp
	case StageJoinBuyVolExp:
		log.Printf("[stage => %v] join buy vol exp: result = %v", computation.ID[31], value)
		computer.volExpPointer[computation.Buy] = &value
		computation.ID[31] = StageJoinBuyVolCo
	case StageJoinBuyVolCo:
		log.Printf("[stage => %v] join buy vol co: result = %v", computation.ID[31], value)
		computer.volCoPointer[computation.Buy] = &value
		computation.ID[31] = StageJoinBuyMinVolExp
	case StageJoinBuyMinVolExp:
		log.Printf("[stage => %v] join buy min vol exp: result = %v", computation.ID[31], value)
		computer.minVolExpPointer[computation.Buy] = &value
		computation.ID[31] = StageJoinBuyMinVolCo
	case StageJoinBuyMinVolCo:
		log.Printf("[stage => %v] join buy min vol co: result = %v", computation.ID[31], value)
		computer.minVolCoPointer[computation.Buy] = &value
		computation.ID[31] = StageJoinBuyTokens
	case StageJoinBuyTokens:
		log.Printf("[stage => %v] join buy token: result = %v", computation.ID[31], value)
		tokens := order.Tokens(value)
		computer.tokensPointer[computation.Buy] = &tokens
		computation.ID[31] = StageJoinSellPriceExp
	case StageJoinSellPriceExp:
		log.Printf("[stage => %v] join sell price exp: result = %v", computation.ID[31], value)
		computer.priceExpPointer[computation.Sell] = &value
		computation.ID[31] = StageJoinSellPriceCo
	case StageJoinSellPriceCo:
		log.Printf("[stage => %v] join sell price co: result = %v", computation.ID[31], value)
		computer.priceCoPointer[computation.Sell] = &value
		computation.ID[31] = StageJoinSellVolExp
	case StageJoinSellVolExp:
		log.Printf("[stage => %v] join sell vol exp: result = %v", computation.ID[31], value)
		computer.volExpPointer[computation.Sell] = &value
		computation.ID[31] = StageJoinSellVolCo
	case StageJoinSellVolCo:
		log.Printf("[stage => %v] join sell vol co: result = %v", computation.ID[31], value)
		computer.volCoPointer[computation.Sell] = &value
		computation.ID[31] = StageJoinSellMinVolExp
	case StageJoinSellMinVolExp:
		log.Printf("[stage => %v] join sell min vol exp: result = %v", computation.ID[31], value)
		computer.minVolExpPointer[computation.Sell] = &value
		computation.ID[31] = StageJoinSellMinVolCo
	case StageJoinSellMinVolCo:
		log.Printf("[stage => %v] join sell min vol co: result = %v", computation.ID[31], value)
		computer.minVolCoPointer[computation.Sell] = &value
		computation.ID[31] = StageJoinSellTokens
	case StageJoinSellTokens:
		log.Printf("[stage => %v] join sell token: result = %v", computation.ID[31], value)
		tokens := order.Tokens(value)
		computer.tokensPointer[computation.Sell] = &tokens
		log.Printf(" start settl orders (Buy:%v , Sell: %v)-------", common.ToHex(computation.Buy[:]), common.ToHex(computation.Sell[:]))

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
		log.Printf("<sell order> price<Co: %v, Exp: %v>, volumn<Co: %v, Exp: %v>,", buy.Price.Co, buy.Price.Exp, buy.Volume.Co, buy.Volume.Exp)

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

	// FIXME: This cannot be escaped when the done channel is closed and will
	// deadlock if the notification does not come through on a different
	// goroutine from the goroutine that emits computations.
	log.Println("writing computation")
	computer.computations <- computation
	log.Println("computation written")
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
