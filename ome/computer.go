package ome

import (
	"bytes"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
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
	Compute(networkID [32]byte, done <-chan struct{}, computations <-chan Computation) (<-chan Computation, <-chan error)
}

type computer struct {
	storer    orderbook.Storer
	smpcer    smpc.Smpcer
	confirmer confirmer

	cmpMu         *sync.Mutex
	cmpPriceExp   map[[32]byte]ComputationState
	cmpPriceCo    map[[32]byte]ComputationState
	cmpBuyVolExp  map[[32]byte]ComputationState
	cmpBuyVolCo   map[[32]byte]ComputationState
	cmpSellVolExp map[[32]byte]ComputationState
	cmpSellVolCo  map[[32]byte]ComputationState
	cmpTokens     map[[32]byte]ComputationState

	joinMu              *sync.Mutex
	orders              map[[32]byte]order.Order
	priceExpPointer     map[[32]byte]*uint64
	volumeExpPointer    map[[32]byte]*uint64
	minVolumeExpPointer map[[32]byte]*uint64
	joinBuyPriceExp     map[[32]byte]ComputationState
	joinBuyPriceCo      map[[32]byte]ComputationState
	joinBuyVolExp       map[[32]byte]ComputationState
	joinBuyVolCo        map[[32]byte]ComputationState
	joinBuyMinVolExp    map[[32]byte]ComputationState
	joinBuyMinVolCo     map[[32]byte]ComputationState
	joinBuyTokens       map[[32]byte]ComputationState
	joinSellPriceExp    map[[32]byte]ComputationState
	joinSellPriceCo     map[[32]byte]ComputationState
	joinSellVolExp      map[[32]byte]ComputationState
	joinSellVolCo       map[[32]byte]ComputationState
	joinSellMinVolExp   map[[32]byte]ComputationState
	joinSellMinVolCo    map[[32]byte]ComputationState
	joinSellTokens      map[[32]byte]ComputationState
}

func NewComputer(storer orderbook.Storer, smpcer smpc.Smpcer, confirmer confirmer) Computer {
	return &computer{
		storer:    storer,
		smpcer:    smpcer,
		confirmer: confirmer,

		cmpMu:         new(sync.Mutex),
		cmpPriceExp:   map[[32]byte]ComputationState{},
		cmpPriceCo:    map[[32]byte]ComputationState{},
		cmpBuyVolExp:  map[[32]byte]ComputationState{},
		cmpBuyVolCo:   map[[32]byte]ComputationState{},
		cmpSellVolExp: map[[32]byte]ComputationState{},
		cmpSellVolCo:  map[[32]byte]ComputationState{},
		cmpTokens:     map[[32]byte]ComputationState{},

		joinMu:              new(sync.Mutex),
		orders:              map[[32]byte]order.Order{},
		priceExpPointer:     map[[32]byte]*uint64{},
		volumeExpPointer:    map[[32]byte]*uint64{},
		minVolumeExpPointer: map[[32]byte]*uint64{},

		joinBuyPriceExp:   map[[32]byte]ComputationState{},
		joinBuyPriceCo:    map[[32]byte]ComputationState{},
		joinBuyVolExp:     map[[32]byte]ComputationState{},
		joinBuyVolCo:      map[[32]byte]ComputationState{},
		joinBuyMinVolExp:  map[[32]byte]ComputationState{},
		joinBuyMinVolCo:   map[[32]byte]ComputationState{},
		joinBuyTokens:     map[[32]byte]ComputationState{},
		joinSellPriceExp:  map[[32]byte]ComputationState{},
		joinSellPriceCo:   map[[32]byte]ComputationState{},
		joinSellVolExp:    map[[32]byte]ComputationState{},
		joinSellVolCo:     map[[32]byte]ComputationState{},
		joinSellMinVolExp: map[[32]byte]ComputationState{},
		joinSellMinVolCo:  map[[32]byte]ComputationState{},
		joinSellTokens:    map[[32]byte]ComputationState{},
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
							}
						}
					}
				}
			},
			func() {
				computations, confirmErrs := computer.confirmer.ConfirmOrderMatches(done, orderMatches)
				for {
					select {
					case computation, ok := <-computations:
						if !ok {
							return
						}
						//Todo: reconstruct the order from the computation
						buyFragment, err := computer.storer.OrderFragment(computation.Buy)
						if err != nil {
							errs <- err
							continue
						}
						sellFragment, err := computer.storer.OrderFragment(computation.Sell)
						if err != nil {
							errs <- err
							continue
						}
						buyOrder := order.Order{
							Type:   buyFragment.OrderType,
							Parity: buyFragment.OrderParity,
							Expiry: buyFragment.OrderExpiry,
						}
						computer.orders[computation.Buy] = buyOrder

						sellOrder := order.Order{
							Type:   sellFragment.OrderType,
							Parity: sellFragment.OrderParity,
							Expiry: sellFragment.OrderExpiry,
						}
						computer.orders[computation.Sell] = sellOrder

						id := computeID(computation)
						id[31] = StageJoinBuyPriceExp
						computer.joinBuyPriceExp[id] = ComputationState{
							Computation: computation,
							State:       StatePending,
						}
					case err, ok := <-confirmErrs:
						if !ok {
							return
						}
						errs <- err
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
					for id, computation := range computer.joinBuyPriceExp {
						if computation.State == StatePending {
							buy, err := computer.storer.OrderFragment(computation.Buy)
							if err != nil {
								log.Printf("cannot get buy order fragment from orderbook: %v", err)
								continue
							}
							computer.joinMu.Lock()
							computation.State = StateComputing
							computer.joinBuyPriceExp[id] = computation
							computer.joinMu.Unlock()
							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: buy.Price.Exp,
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
					for id, computation := range computer.joinBuyPriceCo {
						if computation.State == StatePending {
							buy, err := computer.storer.OrderFragment(computation.Buy)
							if err != nil {
								log.Printf("cannot get buy order fragment from orderbook: %v", err)
								continue
							}
							computer.joinMu.Lock()
							computation.State = StateComputing
							computer.joinBuyPriceCo[id] = computation
							computer.joinMu.Unlock()
							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: buy.Price.Co,
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
					for id, computation := range computer.joinBuyVolExp {
						if computation.State == StatePending {
							buy, err := computer.storer.OrderFragment(computation.Buy)
							if err != nil {
								log.Printf("cannot get buy order fragment from orderbook: %v", err)
								continue
							}
							computer.joinMu.Lock()
							computation.State = StateComputing
							computer.joinBuyVolExp[id] = computation
							computer.joinMu.Unlock()
							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: buy.Volume.Exp,
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
					for id, computation := range computer.joinBuyVolCo {
						if computation.State == StatePending {
							buy, err := computer.storer.OrderFragment(computation.Buy)
							if err != nil {
								log.Printf("cannot get buy order fragment from orderbook: %v", err)
								continue
							}
							computer.joinMu.Lock()
							computation.State = StateComputing
							computer.joinBuyVolCo[id] = computation
							computer.joinMu.Unlock()
							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: buy.Volume.Co,
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
					for id, computation := range computer.joinBuyMinVolExp {
						if computation.State == StatePending {
							buy, err := computer.storer.OrderFragment(computation.Buy)
							if err != nil {
								log.Printf("cannot get buy order fragment from orderbook: %v", err)
								continue
							}
							computer.joinMu.Lock()
							computation.State = StateComputing
							computer.joinBuyMinVolExp[id] = computation
							computer.joinMu.Unlock()
							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: buy.MinimumVolume.Exp,
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
					for id, computation := range computer.joinBuyMinVolCo {
						if computation.State == StatePending {
							buy, err := computer.storer.OrderFragment(computation.Buy)
							if err != nil {
								log.Printf("cannot get buy order fragment from orderbook: %v", err)
								continue
							}
							computer.joinMu.Lock()
							computation.State = StateComputing
							computer.joinBuyMinVolCo[id] = computation
							computer.joinMu.Unlock()
							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: buy.MinimumVolume.Co,
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
					for id, computation := range computer.joinBuyTokens {
						if computation.State == StatePending {
							buy, err := computer.storer.OrderFragment(computation.Buy)
							if err != nil {
								log.Printf("cannot get buy order fragment from orderbook: %v", err)
								continue
							}
							computer.joinMu.Lock()
							computation.State = StateComputing
							computer.joinBuyTokens[id] = computation
							computer.joinMu.Unlock()
							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: buy.Tokens,
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
					for id, computation := range computer.joinSellPriceExp {
						if computation.State == StatePending {
							sell, err := computer.storer.OrderFragment(computation.Sell)
							if err != nil {
								log.Printf("cannot get sell order fragment from orderbook: %v", err)
								continue
							}
							computer.joinMu.Lock()
							computation.State = StateComputing
							computer.joinSellPriceExp[id] = computation
							computer.joinMu.Unlock()
							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: sell.Price.Exp,
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
					for id, computation := range computer.joinSellPriceCo {
						if computation.State == StatePending {
							sell, err := computer.storer.OrderFragment(computation.Sell)
							if err != nil {
								log.Printf("cannot get sell order fragment from orderbook: %v", err)
								continue
							}
							computer.joinMu.Lock()
							computation.State = StateComputing
							computer.joinSellPriceCo[id] = computation
							computer.joinMu.Unlock()
							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: sell.Price.Co,
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
					for id, computation := range computer.joinSellVolExp {
						if computation.State == StatePending {
							sell, err := computer.storer.OrderFragment(computation.Sell)
							if err != nil {
								log.Printf("cannot get sell order fragment from orderbook: %v", err)
								continue
							}
							computer.joinMu.Lock()
							computation.State = StateComputing
							computer.joinSellVolExp[id] = computation
							computer.joinMu.Unlock()
							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: sell.Volume.Exp,
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
					for id, computation := range computer.joinSellVolCo {
						if computation.State == StatePending {
							sell, err := computer.storer.OrderFragment(computation.Sell)
							if err != nil {
								log.Printf("cannot get sell order fragment from orderbook: %v", err)
								continue
							}
							computer.joinMu.Lock()
							computation.State = StateComputing
							computer.joinSellVolCo[id] = computation
							computer.joinMu.Unlock()
							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: sell.Volume.Co,
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
					for id, computation := range computer.joinSellMinVolExp {
						if computation.State == StatePending {
							sell, err := computer.storer.OrderFragment(computation.Sell)
							if err != nil {
								log.Printf("cannot get sell order fragment from orderbook: %v", err)
								continue
							}
							computer.joinMu.Lock()
							computation.State = StateComputing
							computer.joinSellMinVolExp[id] = computation
							computer.joinMu.Unlock()
							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: sell.MinimumVolume.Exp,
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
					for id, computation := range computer.joinSellMinVolCo {
						if computation.State == StatePending {
							sell, err := computer.storer.OrderFragment(computation.Sell)
							if err != nil {
								log.Printf("cannot get sell order fragment from orderbook: %v", err)
								continue
							}
							computer.joinMu.Lock()
							computation.State = StateComputing
							computer.joinSellMinVolCo[id] = computation
							computer.joinMu.Unlock()
							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: sell.MinimumVolume.Co,
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
					for id, computation := range computer.joinSellTokens {
						if computation.State == StatePending {
							sell, err := computer.storer.OrderFragment(computation.Sell)
							if err != nil {
								log.Printf("cannot get sell order fragment from orderbook: %v", err)
								continue
							}
							computer.joinMu.Lock()
							computation.State = StateComputing
							computer.joinSellTokens[id] = computation
							computer.joinMu.Unlock()
							instructions <- smpc.Inst{
								InstID:    id,
								NetworkID: networkID,
								InstJ: &smpc.InstJ{
									Share: sell.Tokens,
								},
							}
						}
					}
					time.Sleep(time.Millisecond)
				}
			},
		)
	}()

	return orderMatches, errs
}

func (computer *computer) processResultJ(instID, networkID [32]byte, resultJ smpc.ResultJ) (*Computation, error) {
	half := shamir.Prime / 2

	computer.cmpMu.Lock()
	computer.joinMu.Lock()
	defer computer.cmpMu.Unlock()
	defer computer.joinMu.Unlock()

	switch instID[31] {

	case StageCmpPriceExp:
		computation, ok := computer.cmpPriceExp[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.cmpPriceExp, instID)
		if resultJ.Value <= half {
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

	case StageJoinBuyPriceExp:
		computation, ok := computer.joinBuyPriceExp[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.joinBuyPriceExp, instID)
		computer.priceExpPointer[computation.Buy] = &resultJ.Value
		computation.State = StatePending
		instID[31] = StageJoinBuyPriceCo
		computer.joinBuyPriceCo[instID] = computation

	case StageJoinBuyPriceCo:
		computation, ok := computer.joinBuyPriceCo[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.joinBuyPriceCo, instID)
		computer.orders[computation.Buy].Price.Co = resultJ.Value
		computation.State = StatePending
		instID[31] = StageJoinBuyVolExp
		computer.joinBuyVolExp[instID] = computation

	case StageJoinBuyVolExp:
		computation, ok := computer.joinBuyVolExp[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.joinBuyVolExp, instID)
		computer.volumeExpPointer[computation.Buy] = &resultJ.Value
		computation.State = StatePending
		instID[31] = StageJoinBuyVolCo
		computer.joinBuyVolCo[instID] = computation

	case StageJoinBuyVolCo:
		computation, ok := computer.joinBuyVolCo[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.joinBuyVolCo, instID)
		computer.orders[computation.Buy].Volume.Co = resultJ.Value
		computation.State = StatePending
		instID[31] = StageJoinBuyTokens
		computer.joinBuyMinVolExp[instID] = computation

	case StageJoinBuyMinVolExp:
		computation, ok := computer.joinBuyMinVolExp[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.joinBuyMinVolExp, instID)
		computer.minVolumeExpPointer[computation.Buy] = &resultJ.Value
		computation.State = StatePending
		instID[31] = StageJoinBuyMinVolCo
		computer.joinBuyMinVolCo[instID] = computation

	case StageJoinBuyMinVolCo:
		computation, ok := computer.joinBuyMinVolCo[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.joinBuyMinVolCo, instID)
		computer.orders[computation.Buy].MinimumVolume.Co = resultJ.Value
		computation.State = StatePending
		instID[31] = StageJoinBuyTokens
		computer.joinBuyTokens[instID] = computation

	case StageJoinBuyTokens:
		computation, ok := computer.joinBuyTokens[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.joinBuyTokens, instID)
		ord := computer.orders[computation.Buy]
		ord.Tokens = order.Tokens(resultJ.Value)
		computer.orders[computation.Buy] = ord
		computation.State = StatePending
		instID[31] = StageJoinSellPriceExp
		computer.joinSellPriceExp[instID] = computation

	case StageJoinSellPriceExp:
		computation, ok := computer.joinSellPriceExp[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.joinSellPriceExp, instID)
		computer.priceExpPointer[computation.Sell] = &resultJ.Value
		computation.State = StatePending
		instID[31] = StageJoinSellPriceCo
		computer.joinSellPriceCo[instID] = computation

	case StageJoinSellPriceCo:
		computation, ok := computer.joinSellPriceCo[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.joinSellPriceCo, instID)
		computer.orders[computation.Sell].Price.Co = resultJ.Value
		computation.State = StatePending
		instID[31] = StageJoinSellVolExp
		computer.joinSellVolExp[instID] = computation

	case StageJoinSellVolExp:
		computation, ok := computer.joinSellVolExp[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.joinSellVolExp, instID)
		computer.volumeExpPointer[computation.Sell] = &resultJ.Value
		computation.State = StatePending
		instID[31] = StageJoinSellVolCo
		computer.joinSellVolCo[instID] = computation

	case StageJoinSellVolCo:
		computation, ok := computer.joinSellVolCo[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.joinSellVolCo, instID)
		computer.orders[computation.Sell].Volume.Co = resultJ.Value
		computation.State = StatePending
		instID[31] = StageJoinSellTokens
		computer.joinSellMinVolExp[instID] = computation

	case StageJoinSellMinVolExp:
		computation, ok := computer.joinSellMinVolExp[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.joinSellMinVolExp, instID)
		computer.minVolumeExpPointer[computation.Sell] = &resultJ.Value
		computation.State = StatePending
		instID[31] = StageJoinSellMinVolCo
		computer.joinSellMinVolCo[instID] = computation

	case StageJoinSellMinVolCo:
		computation, ok := computer.joinSellMinVolCo[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.joinSellMinVolCo, instID)
		computer.orders[computation.Sell].MinimumVolume.Co = resultJ.Value
		computation.State = StatePending
		instID[31] = StageJoinSellTokens
		computer.joinBuyTokens[instID] = computation

	case StageJoinSellTokens:
		computation, ok := computer.joinSellTokens[instID]
		if !ok {
			return nil, nil
		}
		delete(computer.joinSellTokens, instID)
		ord := computer.orders[computation.Sell]
		ord.Tokens = order.Tokens(resultJ.Value)

		buyOrder, sellOrder := computer.orders[computation.Buy], computer.orders[computation.Sell]
		if computer.notNil(buyOrder) {
			computer.storer.InsertOrder(buyOrder)
		} else {
			return nil, errors.New("buy order reconstruction fails, some fields are nil")
		}
		if computer.notNil(sellOrder) {
			computer.storer.InsertOrder(sellOrder)
		} else {
			return nil, errors.New(" sell order reconstruction fails, some fields are nil")
		}
	}

	return nil, nil
}

func computeID(computation Computation) [32]byte {
	id := [32]byte{}
	copy(id[:], crypto.Keccak256(computation.Buy[:], computation.Sell[:]))
	return id
}

func (computer *computer) notNil(ord order.Order) bool {
	if ord.Equal(new(order.Order)) {
		return false
	}
	if bytes.Compare(ord.ID[:], []byte{}) == 0 {
		return false
	}
	// fiXME : orderType , orderParity and expiryTime
	if ord.Tokens == 0 {
		return false
	}
	if ord.Price.Co == 0 || ord.Volume.Co == 0 || ord.MinimumVolume.Co == 0 {
		return false
	}
	if computer.priceExpPointer[ord.ID] == nil || computer.volumeExpPointer[ord.ID] == nil || computer.minVolumeExpPointer[ord.ID] == nil {
		return false
	}

	return true
}
