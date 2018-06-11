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
	address        identity.Address
	computationsMu *sync.Mutex
	computations   []Computation

	rankerMu        *sync.Mutex
	rankers         []*epochRanker
	changeChs       []chan orderbook.Change
	computationsChs []<-chan Computation
	blockNumbers    []uint

	storer Storer
}

// NewRanker returns a Ranker that first filters the Computations it produces
// by checking the Priority. The filter assumes that there are a certain number
// of Rankers, and that each Ranker has a unique position relative to others.
// Priorities that do not match the position of the Ranker, after a modulo of
// the number of Rankers, are filtered. A Storer is used to load existing
// Computations that have not been processed completely, and to store new
// Computations.
func NewRanker(storer Storer, address identity.Address, epoch cal.Epoch) (Ranker, error) {
	ranker := &delegateRanker{
		computationsMu: new(sync.Mutex),
		computations:   []Computation{},

		rankerMu:        new(sync.Mutex),
		rankers:         make([]*epochRanker, 0),
		changeChs:       make([]chan orderbook.Change, 0),
		computationsChs: make([]<-chan Computation, 0),
		blockNumbers:    make([]uint, 0),

		storer: storer,
	}

	pods, pos, err := ranker.getPosFromEpoch(epoch)
	if err != nil {
		return &delegateRanker{}, err
	}
	ranker.changeChs = append(ranker.changeChs, make(chan orderbook.Change))
	ranker.rankers = append(ranker.rankers, newEpochRanker(pods, pos))
	ranker.computationsChs = append(ranker.computationsChs, ranker.rankers[1].run(ranker.changeChs[1]))
	ranker.blockNumbers = append(ranker.blockNumbers, epoch.BlockNumber)

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
		logger.Error(fmt.Sprintf("%v", err))
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
	computations := epochRanker.run(ranker.changeChs[1])
	ranker.rankers = append(ranker.rankers, epochRanker)
	ranker.computationsChs = append(ranker.computationsChs, computations)
	ranker.blockNumbers = append(ranker.blockNumbers, epoch.BlockNumber)
}

func (ranker *delegateRanker) processComputations() {
	go func() {
		for {
			ranker.rankerMu.Lock()
			select {
			case computation, ok := <-ranker.computationsChs[0]:
				if !ok {
					return
				}
				ranker.insertComputation(computation)
			case computation, ok := <-ranker.computationsChs[1]:
				if !ok {
					return
				}
				ranker.insertComputation(computation)
			}
			ranker.rankerMu.Unlock()
		}
	}()
}

func (ranker *delegateRanker) insertStoredComputationsInBackground() {
	go func() {
		// Wait for long enough that the Ome has time to connect to the network
		// for the current epoch before loading computations
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
	buys            map[order.ID]uint64
	sells           map[order.ID]uint64
}

func newEpochRanker(numberOfRankers, pos int) *epochRanker {
	return &epochRanker{
		numberOfRankers: numberOfRankers,
		pos:             pos,
		buys:            map[order.ID]uint64{},
		sells:           map[order.ID]uint64{},
	}
}

func (ranker *epochRanker) run(changes <-chan orderbook.Change) <-chan Computation {
	computationsOut := make(chan Computation)

	go func() {
		defer ranker.cleanup()
		defer close(computationsOut)

		for change := range changes {

			var computations []Computation
			switch change.OrderStatus {
			case order.Open:
				if change.OrderParity == order.ParityBuy {
					computations = ranker.insertBuy(change)
				} else {
					computations = ranker.insertSell(change)
				}
			case order.Canceled, order.Confirmed:
				ranker.remove(change)
			}

			for _, computation := range computations {
				computationsOut <- computation
			}
		}
	}()

	return computationsOut
}

func (ranker *epochRanker) cleanup() {
	ranker.buys = map[order.ID]uint64{}
	ranker.sells = map[order.ID]uint64{}
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
