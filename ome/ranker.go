package ome

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
)

// A Ranker consumes orders and produces Computations that are prioritized
// based on the combined priorities of the involved orders.
type Ranker interface {
	// InsertChange into the Ranker. Ranker will generate new computations or
	// remove order from the ranker depending on the new status of the order.
	// New computations can be read by using a call to Ranker.Computations.
	InsertChange(change orderbook.Change)

	// Computations stored in the Ranker are written to the input buffer. The
	// written Computations are removed from the Ranker to prevent duplicate
	// Computations. Returns the number of Computations written to the buffer,
	// which is guaranteed to be less than, or equal to, the size of the
	// buffer.
	Computations(Computations) int

	// OnChangeEpoch should be called whenever a new cal.Epoch is observed.
	OnChangeEpoch(epoch cal.Epoch)
}

// delegateRanker delegates orders to specific epochRanker according to the
// epoch Hash and collects computations back from all the epochRanker.
type delegateRanker struct {
	done    <-chan struct{}
	address identity.Address
	storer  Storer

	computationsMu *sync.Mutex
	computations   Computations

	rankerMu        *sync.Mutex
	rankers         []*epochRanker
	changeChs       []chan orderbook.Change
	computationsChs []<-chan Computations
	blockNumbers    []uint
}

// NewRanker returns a Ranker that first filters the Computations it produces
// by checking the Priority. The filter assumes that there are a certain number
// of Rankers, and that each Ranker has a unique position relative to others.
// Priorities that do not match the position of the Ranker, after a modulo of
// the number of Rankers, are filtered. A Storer is used to load existing
// Computations that have not been processed completely, and to store new
// Computations. The Ranker will run background processes until the done
// channel is closed, after which the Ranker will no longer consume
// orderbook.Changeset or produce Computation.
func NewRanker(done <-chan struct{}, address identity.Address, storer Storer, epoch cal.Epoch) (Ranker, error) {
	ranker := &delegateRanker{
		done:    done,
		address: address,
		storer:  storer,

		computationsMu: new(sync.Mutex),
		computations:   Computations{},

		rankerMu:        new(sync.Mutex),
		rankers:         make([]*epochRanker, 0),
		changeChs:       make([]chan orderbook.Change, 0),
		computationsChs: make([]<-chan Computations, 0),
		blockNumbers:    make([]uint, 0),
	}

	pods, pos, err := ranker.getPosFromEpoch(epoch)
	if err != nil {
		return &delegateRanker{}, err
	}
	ranker.changeChs = append(ranker.changeChs, make(chan orderbook.Change))
	ranker.rankers = append(ranker.rankers, newEpochRanker(pods, pos))
	ranker.computationsChs = append(ranker.computationsChs, ranker.rankers[1].run(ranker.done, ranker.changeChs[1]))
	ranker.blockNumbers = append(ranker.blockNumbers, epoch.BlockNumber)

	ranker.processComputations(done)
	ranker.insertStoredComputationsInBackground()

	return ranker, nil
}

// InsertChange
func (ranker *delegateRanker) InsertChange(change orderbook.Change) {
	ranker.rankerMu.Lock()
	defer ranker.rankerMu.Unlock()

	if change.BlockNumber >= ranker.blockNumbers[0] && change.BlockNumber < ranker.blockNumbers[1] {
		ranker.changeChs[0] <- change
	} else if change.BlockNumber >= ranker.blockNumbers[2] {
		ranker.changeChs[1] <- change
	}
}

func (ranker *delegateRanker) Computations(buffer Computations) int {
	ranker.computationsMu.Lock()
	defer ranker.computationsMu.Unlock()

	var min int
	if len(buffer) < len(ranker.computations) {
		min = len(buffer)
	} else {
		min = len(ranker.computations)
	}
	for i := 0; i < min; i++ {
		buffer[i] = ranker.computations[i]
	}
	ranker.computations = ranker.computations[:min]

	return min
}

func (ranker *delegateRanker) OnChangeEpoch(epoch cal.Epoch) {
	ranker.rankerMu.Lock()
	defer ranker.rankerMu.Unlock()

	pods, pos, err := ranker.getPosFromEpoch(epoch)
	if err != nil {
		logger.Error(fmt.Sprintf("cannot get ranker position from epoch: %v", err))
		return
	}

	epochRanker := newEpochRanker(pods, pos)
	if len(ranker.rankers) > 1 {
		close(ranker.changeChs[0])
		ranker.changeChs = ranker.changeChs[1:]
		ranker.rankers = ranker.rankers[1:]
		ranker.computationsChs = ranker.computationsChs[1:]
		ranker.blockNumbers = ranker.blockNumbers[1:]
	}

	ranker.changeChs = append(ranker.changeChs, make(chan orderbook.Change))
	computations := epochRanker.run(ranker.done, ranker.changeChs[1])
	ranker.rankers = append(ranker.rankers, epochRanker)
	ranker.computationsChs = append(ranker.computationsChs, computations)
	ranker.blockNumbers = append(ranker.blockNumbers, epoch.BlockNumber)
}

func (ranker *delegateRanker) processComputations(done <-chan struct{}) {
	go func() {
		for {
			ranker.rankerMu.Lock()
			currEpochRankerCh := ranker.computationsChs[0]
			prevEpochRankerCh := ranker.computationsChs[1]
			ranker.rankerMu.Unlock()

			select {
			case <-done:
				return
			case coms, ok := <-currEpochRankerCh:
				if !ok {
					return
				}
				for _, com := range coms {
					ranker.insertComputation(com)
				}
			case coms, ok := <-prevEpochRankerCh:
				if !ok {
					return
				}
				for _, com := range coms {
					ranker.insertComputation(com)
				}
			}
		}
	}()
}

func (ranker *delegateRanker) insertStoredComputationsInBackground() {
	go func() {
		// Wait for long enough that the Ome has time to connect to the network
		// for the current epoch before loading computations (approximately one
		// block)
		timer := time.NewTimer(14 * time.Second)

		coms, err := ranker.storer.Computations()
		if err != nil {
			logger.Error(fmt.Sprintf("cannot load existing computations into ranker: %v", err))
		}

		<-timer.C
		for _, com := range coms {
			if com.State != ComputationStateMismatched && com.State != ComputationStateRejected && com.State != ComputationStateSettled {
				ranker.insertComputation(com)
			}
		}
	}()
}

func (ranker *delegateRanker) insertComputation(com Computation) {
	ranker.computationsMu.Lock()
	defer ranker.computationsMu.Unlock()

	index := sort.Search(len(ranker.computations), func(i int) bool {
		return ranker.computations[i].Priority > com.Priority
	})
	ranker.computations = append(
		ranker.computations[:index],
		append([]Computation{com}, ranker.computations[index:]...)...)
}

func (ranker *delegateRanker) getPosFromEpoch(epoch cal.Epoch) (int, int, error) {
	pod, err := epoch.Pod(ranker.address)
	if err != nil {
		return 0, 0, err
	}
	return len(epoch.Pods), pod.Position, nil
}

// epochRanker forms new computation and rank them depending on the priority.
// It only cares about orders from one dedicated epoch, so that we won't
// cross match orders from different epoch.
type epochRanker struct {
	numberOfRankers int
	pos             int
	buys            map[order.ID]orderbook.Priority
	sells           map[order.ID]orderbook.Priority
}

func newEpochRanker(numberOfRankers, pos int) *epochRanker {
	return &epochRanker{
		numberOfRankers: numberOfRankers,
		pos:             pos,
		buys:            map[order.ID]orderbook.Priority{},
		sells:           map[order.ID]orderbook.Priority{},
	}
}

func (ranker *epochRanker) run(done <-chan struct{}, changes <-chan orderbook.Change) <-chan Computations {
	computations := make(chan Computations)

	go func() {
		defer ranker.cleanup()
		defer close(computations)

		for change := range changes {
			switch change.OrderStatus {
			case order.Open:
				if change.OrderParity == order.ParityBuy {
					select {
					case <-done:
						return
					case computations <- ranker.insertBuy(change):
					}
				} else {
					select {
					case <-done:
						return
					case computations <- ranker.insertSell(change):
					}
				}
			case order.Canceled, order.Confirmed:
				ranker.remove(change)
			}
		}
	}()

	return computations
}

func (ranker *epochRanker) cleanup() {
	ranker.buys = map[order.ID]orderbook.Priority{}
	ranker.sells = map[order.ID]orderbook.Priority{}
}

func (ranker *epochRanker) insertBuy(change orderbook.Change) []Computation {
	computations := make([]Computation, 0)
	ranker.buys[change.OrderID] = change.OrderPriority
	for sell, sellPriority := range ranker.sells {
		priority := change.OrderPriority + sellPriority
		if int(priority)%ranker.numberOfRankers != ranker.pos {
			continue
		}

		priorityCom := NewComputation(change.OrderID, sell)
		priorityCom.Priority = priority
		priorityCom.Timestamp = time.Now()

		computations = append(computations, priorityCom)
	}

	return computations
}

func (ranker *epochRanker) insertSell(change orderbook.Change) []Computation {
	computations := make([]Computation, 0)
	ranker.buys[change.OrderID] = change.OrderPriority
	for buy, buyPriority := range ranker.buys {
		priority := change.OrderPriority + buyPriority
		if int(priority)%ranker.numberOfRankers != ranker.pos {
			continue
		}

		priorityCom := NewComputation(buy, change.OrderID)
		priorityCom.Priority = priority
		priorityCom.Timestamp = time.Now()

		computations = append(computations, priorityCom)
	}

	return computations
}

func (ranker *epochRanker) remove(change orderbook.Change) {
	if change.OrderParity == order.ParityBuy {
		delete(ranker.buys, change.OrderID)
	} else {
		delete(ranker.sells, change.OrderID)
	}
}
