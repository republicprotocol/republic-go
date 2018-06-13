package ome

import (
	"bytes"
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
	// InsertChange into the Ranker. The orderbook.Change will be forwarded to
	// be handled by the respective internal handler based on the block number
	// of the orderbook.Change. This ensures that Computations can be filtered
	// by their epoch.
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
// epoch Hash and collects computations back from all the epochRankers.
type delegateRanker struct {
	done    <-chan struct{}
	address identity.Address
	storer  Storer

	computationsMu *sync.Mutex
	computations   Computations

	rankerMu        *sync.Mutex
	rankerCurrEpoch *epochRanker
	rankerPrevEpoch *epochRanker
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
		rankerCurrEpoch: nil,
		rankerPrevEpoch: nil,
	}
	ranker.insertStoredComputationsInBackground()

	numberOfRankers, pos, err := ranker.posFromEpoch(epoch)
	if err != nil {
		return &delegateRanker{}, fmt.Errorf("cannot get ranker position from epoch: %v", err)
	}
	ranker.rankerCurrEpoch = newEpochRanker(numberOfRankers, pos, epoch)

	return ranker, nil
}

// InsertChange implements the Ranker interface.
func (ranker *delegateRanker) InsertChange(change orderbook.Change) {
	ranker.rankerMu.Lock()
	defer ranker.rankerMu.Unlock()

	coms := Computations{}
	if ranker.rankerCurrEpoch != nil && change.BlockNumber >= ranker.rankerCurrEpoch.epoch.BlockNumber {
		coms = ranker.rankerCurrEpoch.insertChange(change)
	} else if ranker.rankerPrevEpoch != nil && change.BlockNumber >= ranker.rankerPrevEpoch.epoch.BlockNumber {
		coms = ranker.rankerPrevEpoch.insertChange(change)
	}

	ranker.insertComputations(coms)
}

// Computations implements the Ranker interface.
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
	ranker.computations = ranker.computations[min:]

	return min
}

// OnChangeEpoch implements the Ranker interface.
func (ranker *delegateRanker) OnChangeEpoch(epoch cal.Epoch) {
	ranker.rankerMu.Lock()
	defer ranker.rankerMu.Unlock()

	if ranker.rankerCurrEpoch != nil && bytes.Equal(ranker.rankerCurrEpoch.epoch.Hash[:], epoch.Hash[:]) {
		return
	}
	ranker.rankerPrevEpoch = ranker.rankerCurrEpoch

	numberOfRankers, pos, err := ranker.posFromEpoch(epoch)
	if err != nil {
		logger.Error(fmt.Sprintf("cannot get ranker position from epoch: %v", err))
		return
	}
	ranker.rankerCurrEpoch = newEpochRanker(numberOfRankers, pos, epoch)
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

func (ranker *delegateRanker) insertComputations(coms Computations) {
	ranker.computationsMu.Lock()
	defer ranker.computationsMu.Unlock()

	for _, com := range coms {
		index := sort.Search(len(ranker.computations), func(i int) bool {
			return ranker.computations[i].Priority > com.Priority
		})
		ranker.computations = append(
			ranker.computations[:index],
			append([]Computation{com}, ranker.computations[index:]...)...)
	}
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

func (ranker *delegateRanker) posFromEpoch(epoch cal.Epoch) (int, int, error) {
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
	epoch           cal.Epoch
	numberOfRankers int
	pos             int
	buys            map[order.ID]orderbook.Priority
	sells           map[order.ID]orderbook.Priority
}

func newEpochRanker(numberOfRankers, pos int, epoch cal.Epoch) *epochRanker {
	return &epochRanker{
		epoch:           epoch,
		numberOfRankers: numberOfRankers,
		pos:             pos,
		buys:            map[order.ID]orderbook.Priority{},
		sells:           map[order.ID]orderbook.Priority{},
	}
}

func (ranker *epochRanker) insertChange(change orderbook.Change) Computations {
	if change.OrderParity == order.ParityBuy {
		return ranker.insertBuyChange(change)
	}
	if change.OrderParity == order.ParitySell {
		return ranker.insertSellChange(change)
	}
	return Computations{}
}

func (ranker *epochRanker) insertBuyChange(change orderbook.Change) Computations {
	if change.OrderStatus != order.Open {
		delete(ranker.buys, change.OrderID)
		return Computations{}
	}

	computations := make([]Computation, 0, len(ranker.sells)/2)
	ranker.buys[change.OrderID] = change.OrderPriority
	for sell, sellPriority := range ranker.sells {
		priority := change.OrderPriority + sellPriority
		if int(priority)%ranker.numberOfRankers != ranker.pos {
			continue
		}
		priorityCom := NewComputation(change.OrderID, sell, ranker.epoch.Hash)
		priorityCom.Priority = priority
		priorityCom.Timestamp = time.Now()
		computations = append(computations, priorityCom)
	}
	return computations
}

func (ranker *epochRanker) insertSellChange(change orderbook.Change) Computations {
	if change.OrderStatus != order.Open {
		delete(ranker.sells, change.OrderID)
		return Computations{}
	}

	computations := make([]Computation, 0, len(ranker.buys)/2)
	ranker.sells[change.OrderID] = change.OrderPriority
	for buy, buyPriority := range ranker.buys {
		priority := change.OrderPriority + buyPriority
		if int(priority)%ranker.numberOfRankers != ranker.pos {
			continue
		}
		priorityCom := NewComputation(buy, change.OrderID, ranker.epoch.Hash)
		priorityCom.Priority = priority
		priorityCom.Timestamp = time.Now()
		computations = append(computations, priorityCom)
	}
	return computations
}
