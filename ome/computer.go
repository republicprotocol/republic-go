package ome

import (
	"log"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dispatch"
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

type State uint64

const (
	StatePending   = 0
	StateComputing = 1
)

type Stage byte

const (
	StageCmpPriceExp   = 1
	StageCmpPriceCo    = 2
	StageCmpBuyVolExp  = 3
	StageCmpBuyVolCo   = 4
	StageCmpSellVolExp = 5
	StageCmpSellVolCo  = 6
	StageCmpTokens     = 7
)

type ComputationState struct {
	Computation
	State
}

// Results is an alias type.
type Results []Result

// A Result is the output of performing Computations.
type Result struct {
	Buy     order.ID
	Sell    order.ID
	IsMatch bool
}

type Computer interface {

	// Compute drives the Computer. It accepts new Computations into its
	// pipeline.
	// TODO: NetworkID should be constructed as part of the computation, and
	// there will need to be done by the orderbook.Server (which knows the
	// epoch at which an order fragment was received).
	Compute(networkID [32]byte, computations Computations)
}

type computer struct {
	orderbook orderbook.Orderbook
	smpcer    smpc.Smpcer

	cmpPriceExp   map[[32]byte]ComputationState
	cmpPriceCo    map[[32]byte]ComputationState
	cmpBuyVolExp  map[[32]byte]ComputationState
	cmpBuyVolCo   map[[32]byte]ComputationState
	cmpSellVolExp map[[32]byte]ComputationState
	cmpSellVolCo  map[[32]byte]ComputationState
	cmpTokens     map[[32]byte]ComputationState
}

func NewComputer(orderbook orderbook.Orderbook, smpcer smpc.Smpcer) Computer {
	return &computer{
		orderbook: orderbook,
		smpcer:    smpcer,

		cmpPriceExp:   map[[32]byte]ComputationState{},
		cmpPriceCo:    map[[32]byte]ComputationState{},
		cmpBuyVolExp:  map[[32]byte]ComputationState{},
		cmpBuyVolCo:   map[[32]byte]ComputationState{},
		cmpSellVolExp: map[[32]byte]ComputationState{},
		cmpSellVolCo:  map[[32]byte]ComputationState{},
		cmpTokens:     map[[32]byte]ComputationState{},
	}
}

func (computer *computer) Compute(networkID [32]byte, computations Computations) {
	instructions := computer.smpcer.Instructions()
	results := computer.smpcer.Results()

	for _, computation := range computations {
		id := computeID(computation)
		id[31] = StageCmpPriceExp
		computer.cmpPriceExp[id] = ComputationState{
			Computation: computation,
			State:       StatePending,
		}
	}

	dispatch.CoBegin(
		func() {
			// Compute price exponenets
			for id, computation := range computer.cmpPriceExp {
				if computation.State == StatePending {
					buy, err := computer.orderbook.OrderFragment(computation.Buy)
					if err != nil {
						log.Printf("cannot get buy order fragment from orderbook: %v", err)
						continue
					}
					sell, err := computer.orderbook.OrderFragment(computation.Sell)
					if err != nil {
						log.Printf("cannot get buy order fragment from orderbook: %v", err)
						continue
					}
					instructions <- smpc.Inst{
						InstID:    id,
						NetworkID: networkID,
						InstJ: &smpc.InstJ{
							Share: buy.Price.Exp.Sub(&sell.Price.Exp),
						},
					}
					computation.State = StateComputing
					computer.cmpPriceExp[id] = computation
				}
			}
		},
		func() {
			// Compute price cos
			for id, computation := range computer.cmpPriceCo {
				if computation.State == StatePending {
					buy, err := computer.orderbook.OrderFragment(computation.Buy)
					if err != nil {
						log.Printf("cannot get buy order fragment from orderbook: %v", err)
						continue
					}
					sell, err := computer.orderbook.OrderFragment(computation.Sell)
					if err != nil {
						log.Printf("cannot get buy order fragment from orderbook: %v", err)
						continue
					}
					instructions <- smpc.Inst{
						InstID:    id,
						NetworkID: networkID,
						InstJ: &smpc.InstJ{
							Share: buy.Price.Co.Sub(&sell.Price.Co),
						},
					}
					computation.State = StateComputing
					computer.cmpPriceCo[id] = computation
				}
			}
		},
		func() {
			// Compute buy volume exponenets
			for id, computation := range computer.cmpBuyVolExp {
				if computation.State == StatePending {
					buy, err := computer.orderbook.OrderFragment(computation.Buy)
					if err != nil {
						log.Printf("cannot get buy order fragment from orderbook: %v", err)
						continue
					}
					sell, err := computer.orderbook.OrderFragment(computation.Sell)
					if err != nil {
						log.Printf("cannot get buy order fragment from orderbook: %v", err)
						continue
					}
					instructions <- smpc.Inst{
						InstID:    id,
						NetworkID: networkID,
						InstJ: &smpc.InstJ{
							Share: buy.Volume.Exp.Sub(&sell.MinimumVolume.Exp),
						},
					}
					computation.State = StateComputing
					computer.cmpBuyVolExp[id] = computation
				}
			}
		},
		func() {
			// Compute buy volume cos
			for id, computation := range computer.cmpBuyVolCo {
				if computation.State == StatePending {
					buy, err := computer.orderbook.OrderFragment(computation.Buy)
					if err != nil {
						log.Printf("cannot get buy order fragment from orderbook: %v", err)
						continue
					}
					sell, err := computer.orderbook.OrderFragment(computation.Sell)
					if err != nil {
						log.Printf("cannot get buy order fragment from orderbook: %v", err)
						continue
					}
					instructions <- smpc.Inst{
						InstID:    id,
						NetworkID: networkID,
						InstJ: &smpc.InstJ{
							Share: buy.Volume.Co.Sub(&sell.MinimumVolume.Co),
						},
					}
					computation.State = StateComputing
					computer.cmpBuyVolCo[id] = computation
				}
			}
		},
		func() {
			// Compute sell volume exponenets
			for id, computation := range computer.cmpSellVolExp {
				if computation.State == StatePending {
					buy, err := computer.orderbook.OrderFragment(computation.Buy)
					if err != nil {
						log.Printf("cannot get buy order fragment from orderbook: %v", err)
						continue
					}
					sell, err := computer.orderbook.OrderFragment(computation.Sell)
					if err != nil {
						log.Printf("cannot get buy order fragment from orderbook: %v", err)
						continue
					}
					instructions <- smpc.Inst{
						InstID:    id,
						NetworkID: networkID,
						InstJ: &smpc.InstJ{
							Share: sell.Volume.Exp.Sub(&buy.MinimumVolume.Exp),
						},
					}
					computation.State = StateComputing
					computer.cmpSellVolExp[id] = computation
				}
			}
		},
		func() {
			// Compute sell volume cos
			for id, computation := range computer.cmpSellVolCo {
				if computation.State == StatePending {
					buy, err := computer.orderbook.OrderFragment(computation.Buy)
					if err != nil {
						log.Printf("cannot get buy order fragment from orderbook: %v", err)
						continue
					}
					sell, err := computer.orderbook.OrderFragment(computation.Sell)
					if err != nil {
						log.Printf("cannot get buy order fragment from orderbook: %v", err)
						continue
					}
					instructions <- smpc.Inst{
						InstID:    id,
						NetworkID: networkID,
						InstJ: &smpc.InstJ{
							Share: sell.Volume.Co.Sub(&buy.MinimumVolume.Co),
						},
					}
					computation.State = StateComputing
					computer.cmpSellVolCo[id] = computation
				}
			}
		},
		func() {
			// Compute tokens
			for id, computation := range computer.cmpTokens {
				if computation.State == StatePending {
					buy, err := computer.orderbook.OrderFragment(computation.Buy)
					if err != nil {
						log.Printf("cannot get buy order fragment from orderbook: %v", err)
						continue
					}
					sell, err := computer.orderbook.OrderFragment(computation.Sell)
					if err != nil {
						log.Printf("cannot get buy order fragment from orderbook: %v", err)
						continue
					}
					instructions <- smpc.Inst{
						InstID:    id,
						NetworkID: networkID,
						InstJ: &smpc.InstJ{
							Share: sell.Tokens.Sub(&buy.Tokens),
						},
					}
					computation.State = StateComputing
					computer.cmpTokens[id] = computation
				}
			}
		},
	)

	nMax := 128
	for n := 0; n < nMax; n++ {
		select {
		case result, ok := <-results:
			if !ok {
				return
			}
			if result.ResultJ != nil {
				computer.processResultJ(result.InstID, result.NetworkID, *result.ResultJ)
			}
		default:
			return
		}
	}
	return
}

func (computer *computer) processResultJ(instID, networkID [32]byte, resultJ smpc.ResultJ) {
	half := shamir.Prime / 2
	switch instID[31] {

	case StageCmpPriceExp:
		computation := computer.cmpPriceExp[instID]
		delete(computer.cmpPriceExp, instID)
		if resultJ.Value <= half {
			computation.State = StatePending
			instID[31] = StageCmpPriceCo
			computer.cmpPriceCo[instID] = computation
		}

	case StageCmpPriceCo:
		computation := computer.cmpPriceCo[instID]
		delete(computer.cmpPriceCo, instID)
		if resultJ.Value <= half {
			computation.State = StatePending
			instID[31] = StageCmpBuyVolExp
			computer.cmpBuyVolExp[instID] = computation
		}

	case StageCmpBuyVolExp:
		computation := computer.cmpBuyVolExp[instID]
		delete(computer.cmpBuyVolExp, instID)
		if resultJ.Value <= half {
			computation.State = StatePending
			instID[31] = StageCmpBuyVolCo
			computer.cmpBuyVolCo[instID] = computation
		}

	case StageCmpBuyVolCo:
		computation := computer.cmpBuyVolCo[instID]
		delete(computer.cmpBuyVolCo, instID)
		if resultJ.Value <= half {
			computation.State = StatePending
			instID[31] = StageCmpSellVolCo
			computer.cmpSellVolCo[instID] = computation
		}

	case StageCmpSellVolExp:
		computation := computer.cmpSellVolExp[instID]
		delete(computer.cmpSellVolExp, instID)
		if resultJ.Value <= half {
			instID[31] = StageCmpSellVolCo
			computer.cmpSellVolCo[instID] = computation
		}

	case StageCmpSellVolCo:
		computation := computer.cmpSellVolCo[instID]
		delete(computer.cmpSellVolCo, instID)
		if resultJ.Value <= half {
			computation.State = StatePending
			instID[31] = StageCmpSellVolCo
			computer.cmpSellVolCo[instID] = computation
		}

	case StageCmpTokens:
		computation := computer.cmpTokens[instID]
		delete(computer.cmpTokens, instID)
		if resultJ.Value == 0 {
			computer.orderbook.ConfirmOrderMatch(computation.Buy, computation.Sell)
		}
	}
}

func computeID(computation Computation) [32]byte {
	id := [32]byte{}
	copy(id[:], crypto.Keccak256(computation.Buy[:], computation.Sell[:]))
	return id
}
