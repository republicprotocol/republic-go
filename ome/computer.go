package ome

import (
	"encoding/base64"
	"log"
	"sync"
	"time"

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

	StageJoinBuyPriceExp  = 8
	StageJoinBuyPriceCo   = 9
	StageJoinBuyVolExp    = 10
	StageJoinBuyVolCo     = 11
	StageJoinBuyMinVolExp = 12
	StageJoinBuyMinVolCo  = 13
	StageJoinBuyTokens    = 14

	StageJoinSellPriceExp  = 15
	StageJoinSellPriceCo   = 16
	StageJoinSellVolExp    = 17
	StageJoinSellVolCo     = 18
	StageJoinSellMinVolExp = 19
	StageJoinSellMinVolCo  = 20
	StageJoinSellTokens    = 21
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

// A Computer consumes Computations that need to be procesed and outputs
// Computations that result in an order match.
type Computer interface {

	// Compute order comparisons, component by component. If all components
	// match, then an order match has been found. Computations are consumed
	// from an input channel, and Computations that resulted in an order match
	// are produced to the output channel. A call to Computer.Compute should
	// use a background goroutine to avoid blocking the caller.
	// TODO: NetworkID should be constructed as part of the computation, and
	// there will need to be done by the orderbook.Server (which knows the
	// epoch at which an order fragment was received).
	Compute(networkID [32]byte, done <-chan struct{}, computations <-chan Computation) (<-chan Computation, <-chan error)
}

type computer struct {
	storer orderbook.Storer
	smpcer smpc.Smpcer

	cmpMu         *sync.Mutex
	cmpPriceExp   map[[32]byte]ComputationState
	cmpPriceCo    map[[32]byte]ComputationState
	cmpBuyVolExp  map[[32]byte]ComputationState
	cmpBuyVolCo   map[[32]byte]ComputationState
	cmpSellVolExp map[[32]byte]ComputationState
	cmpSellVolCo  map[[32]byte]ComputationState
	cmpTokens     map[[32]byte]ComputationState
}

func NewComputer(storer orderbook.Storer, smpcer smpc.Smpcer) Computer {
	return &computer{
		storer: storer,
		smpcer: smpcer,

		cmpMu:         new(sync.Mutex),
		cmpPriceExp:   map[[32]byte]ComputationState{},
		cmpPriceCo:    map[[32]byte]ComputationState{},
		cmpBuyVolExp:  map[[32]byte]ComputationState{},
		cmpBuyVolCo:   map[[32]byte]ComputationState{},
		cmpSellVolExp: map[[32]byte]ComputationState{},
		cmpSellVolCo:  map[[32]byte]ComputationState{},
		cmpTokens:     map[[32]byte]ComputationState{},
	}
}

func (computer *computer) Compute(networkID [32]byte, done <-chan struct{}, computations <-chan Computation) (<-chan Computation, <-chan error) {
	instructions := computer.smpcer.Instructions()
	results := computer.smpcer.Results()

	orderMatches := make(chan Computation)
	errs := make(chan error, 1)

	go func() {
		defer close(orderMatches)
		defer close(errs)

		dispatch.CoBegin(
			func() {
				for {
					select {
					case <-done:
						return
					case computation, ok := <-computations:
						if !ok {
							return
						}
						id := computeID(computation)
						id[31] = StageCmpPriceExp
						computer.cmpPriceExp[id] = ComputationState{
							Computation: computation,
							State:       StatePending,
						}
					}
				}
			},
			func() {
				for {
					select {
					case <-done:
						return
					default:
					}
					for id, computation := range computer.cmpPriceExp {
						if computation.State == StatePending {
							buy, err := computer.storer.OrderFragment(computation.Buy)
							if err != nil {
								continue
							}
							sell, err := computer.storer.OrderFragment(computation.Sell)
							if err != nil {
								continue
							}
							computer.cmpMu.Lock()
							computation.State = StateComputing
							computer.cmpPriceExp[id] = computation
							computer.cmpMu.Unlock()

							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: buy.Price.Exp.Sub(&sell.Price.Exp),
								},
							}
						}
					}
					time.Sleep(time.Millisecond)
				}
			},
			func() {
				for {
					select {
					case <-done:
						return
					default:
					}
					for id, computation := range computer.cmpPriceCo {
						if computation.State == StatePending {
							log.Printf("DEBUG => computing price co: %v, %v", base64.StdEncoding.EncodeToString(computation.Buy[:]), base64.StdEncoding.EncodeToString(computation.Sell[:]))

							buy, err := computer.storer.OrderFragment(computation.Buy)
							if err != nil {
								log.Printf("cannot get buy order fragment from orderbook: %v", err)
								continue
							}
							sell, err := computer.storer.OrderFragment(computation.Sell)
							if err != nil {
								log.Printf("cannot get sell order fragment from orderbook: %v", err)
								continue
							}
							computer.cmpMu.Lock()
							computation.State = StateComputing
							computer.cmpPriceCo[id] = computation
							computer.cmpMu.Unlock()
							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: buy.Price.Co.Sub(&sell.Price.Co),
								},
							}
						}
					}
					time.Sleep(time.Millisecond)
				}
			},
			func() {
				for {
					select {
					case <-done:
						return
					default:
					}
					for id, computation := range computer.cmpBuyVolExp {
						if computation.State == StatePending {
							log.Printf("DEBUG => computing vol exp: %v, %v", base64.StdEncoding.EncodeToString(computation.Buy[:]), base64.StdEncoding.EncodeToString(computation.Sell[:]))

							buy, err := computer.storer.OrderFragment(computation.Buy)
							if err != nil {
								log.Printf("cannot get buy order fragment from orderbook: %v", err)
								continue
							}
							sell, err := computer.storer.OrderFragment(computation.Sell)
							if err != nil {
								log.Printf("cannot get sell order fragment from orderbook: %v", err)
								continue
							}
							computer.cmpMu.Lock()
							computation.State = StateComputing
							computer.cmpBuyVolExp[id] = computation
							computer.cmpMu.Unlock()
							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: buy.Volume.Exp.Sub(&sell.MinimumVolume.Exp),
								},
							}
						}
					}
					time.Sleep(time.Millisecond)
				}
			},
			func() {
				for {
					select {
					case <-done:
						return
					default:
					}
					for id, computation := range computer.cmpBuyVolCo {
						if computation.State == StatePending {
							buy, err := computer.storer.OrderFragment(computation.Buy)
							if err != nil {
								log.Printf("cannot get buy order fragment from orderbook: %v", err)
								continue
							}
							sell, err := computer.storer.OrderFragment(computation.Sell)
							if err != nil {
								log.Printf("cannot get sell order fragment from orderbook: %v", err)
								continue
							}
							computer.cmpMu.Lock()
							computation.State = StateComputing
							computer.cmpBuyVolCo[id] = computation
							computer.cmpMu.Unlock()
							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: buy.Volume.Co.Sub(&sell.MinimumVolume.Co),
								},
							}
						}
					}
					time.Sleep(time.Millisecond)
				}
			},
			func() {
				for {
					select {
					case <-done:
						return
					default:
					}
					for id, computation := range computer.cmpSellVolExp {
						if computation.State == StatePending {
							buy, err := computer.storer.OrderFragment(computation.Buy)
							if err != nil {
								log.Printf("cannot get buy order fragment from orderbook: %v", err)
								continue
							}
							sell, err := computer.storer.OrderFragment(computation.Sell)
							if err != nil {
								log.Printf("cannot get sell order fragment from orderbook: %v", err)
								continue
							}
							computer.cmpMu.Lock()
							computation.State = StateComputing
							computer.cmpSellVolExp[id] = computation
							computer.cmpMu.Unlock()
							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: sell.Volume.Exp.Sub(&buy.MinimumVolume.Exp),
								},
							}
						}
					}
					time.Sleep(time.Millisecond)
				}
			},
			func() {
				for {
					select {
					case <-done:
						return
					default:
					}
					for id, computation := range computer.cmpSellVolCo {
						if computation.State == StatePending {
							buy, err := computer.storer.OrderFragment(computation.Buy)
							if err != nil {
								log.Printf("cannot get buy order fragment from orderbook: %v", err)
								continue
							}
							sell, err := computer.storer.OrderFragment(computation.Sell)
							if err != nil {
								log.Printf("cannot get sell order fragment from orderbook: %v", err)
								continue
							}
							computer.cmpMu.Lock()
							computation.State = StateComputing
							computer.cmpSellVolCo[id] = computation
							computer.cmpMu.Unlock()
							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: sell.Volume.Co.Sub(&buy.MinimumVolume.Co),
								},
							}
						}
					}
					time.Sleep(time.Millisecond)
				}
			},
			func() {
				for {
					select {
					case <-done:
						return
					default:
					}
					for id, computation := range computer.cmpTokens {
						if computation.State == StatePending {
							buy, err := computer.storer.OrderFragment(computation.Buy)
							if err != nil {
								log.Printf("cannot get buy order fragment from orderbook: %v", err)
								continue
							}
							sell, err := computer.storer.OrderFragment(computation.Sell)
							if err != nil {
								log.Printf("cannot get sell order fragment from orderbook: %v", err)
								continue
							}
							computer.cmpMu.Lock()
							computation.State = StateComputing
							computer.cmpTokens[id] = computation
							computer.cmpMu.Unlock()
							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: sell.Tokens.Sub(&buy.Tokens),
								},
							}
						}
					}
					time.Sleep(time.Millisecond)
				}
			},
			func() {
				for {
					select {
					case <-done:
						return
					case result, ok := <-results:
						if !ok {
							return
						}
						if result.ResultJ != nil {
							orderMatch, err := computer.processResultJ(result.InstID, result.NetworkID, *result.ResultJ)
							if err != nil {
								select {
								case <-done:
									return
								case errs <- err:
									continue
								}
							}
							if orderMatch == nil {
								continue
							}
							select {
							case <-done:
								return
							case orderMatches <- *orderMatch:
								// FIXME: We should write this to a confirmer
							}
						}
					}
				}
			},
		)
	}()

	return orderMatches, errs
}

func (computer *computer) processResultJ(instID, networkID [32]byte, resultJ smpc.ResultJ) (*Computation, error) {
	half := shamir.Prime / 2

	computer.cmpMu.Lock()
	defer computer.cmpMu.Unlock()

	switch instID[31] {

	case StageCmpPriceExp:
		computation, ok := computer.cmpPriceExp[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.cmpPriceExp, instID)
		if resultJ.Value <= half {
			log.Println("price exp -> price co ")
			computation.State = StatePending
			instID[31] = StageCmpPriceCo
			computer.cmpPriceCo[instID] = computation
		}

	case StageCmpPriceCo:
		computation, ok := computer.cmpPriceCo[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.cmpPriceCo, instID)
		if resultJ.Value <= half {
			log.Println("price co -> buy volume exp ")
			computation.State = StatePending
			instID[31] = StageCmpBuyVolExp
			computer.cmpBuyVolExp[instID] = computation
		}

	case StageCmpBuyVolExp:
		computation, ok := computer.cmpBuyVolExp[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.cmpBuyVolExp, instID)
		if resultJ.Value <= half {
			log.Println("buy volume exp -> buy volume co")
			computation.State = StatePending
			instID[31] = StageCmpBuyVolCo
			computer.cmpBuyVolCo[instID] = computation
		}

	case StageCmpBuyVolCo:
		computation, ok := computer.cmpBuyVolCo[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.cmpBuyVolCo, instID)
		if resultJ.Value <= half {
			log.Println("buy volume co -> sell volumn exp")
			computation.State = StatePending
			instID[31] = StageCmpSellVolExp
			computer.cmpSellVolExp[instID] = computation
		}

	case StageCmpSellVolExp:
		computation, ok := computer.cmpSellVolExp[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.cmpSellVolExp, instID)
		if resultJ.Value <= half {
			log.Println("sell volume exp -> sell volumn co ")
			computation.State = StatePending
			instID[31] = StageCmpSellVolCo
			computer.cmpSellVolCo[instID] = computation
		}

	case StageCmpSellVolCo:
		computation, ok := computer.cmpSellVolCo[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.cmpSellVolCo, instID)
		if resultJ.Value <= half {
			log.Println("sell volume co -> tokens")
			computation.State = StatePending
			instID[31] = StageCmpTokens
			computer.cmpTokens[instID] = computation
		}

	case StageCmpTokens:
		computation, ok := computer.cmpTokens[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.cmpTokens, instID)
		if resultJ.Value == 0 {
			return &computation.Computation, nil
		}
	}
	return nil, nil
}

func computeID(computation Computation) [32]byte {
	id := [32]byte{}
	copy(id[:], crypto.Keccak256(computation.Buy[:], computation.Sell[:]))
	return id
}
