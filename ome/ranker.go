package ome

import (
	"encoding/base64"
	"fmt"
	"log"
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

	rankerMu           *sync.Mutex
	rankerCurrEpoch    *epochRanker
	rankerCurrEpochIn  chan orderbook.Change
	rankerCurrEpochOut <-chan Computations
	rankerCurrBlockNum uint
	rankerPrevEpoch    *epochRanker
	rankerPrevEpochIn  chan orderbook.Change
	rankerPrevEpochOut <-chan Computations
	rankerPrevBlockNum uint
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

		rankerMu:           new(sync.Mutex),
		rankerCurrEpoch:    nil,
		rankerCurrEpochIn:  nil,
		rankerCurrEpochOut: nil,
		rankerCurrBlockNum: 0,
		rankerPrevEpoch:    nil,
		rankerPrevEpochIn:  nil,
		rankerPrevEpochOut: nil,
		rankerPrevBlockNum: 0,
	}

	numberOfRankers, pos, err := ranker.getPosFromEpoch(epoch)
	if err != nil {
		return &delegateRanker{}, err
	}
	ranker.rankerCurrEpoch = newEpochRanker(numberOfRankers, pos)
	ranker.rankerCurrEpochIn = make(chan orderbook.Change)
	ranker.rankerCurrEpochOut = ranker.rankerCurrEpoch.run(ranker.done, ranker.rankerCurrEpochIn)
	ranker.rankerCurrBlockNum = epoch.BlockNumber
	log.Printf("rankers: %d, pos : %d, blockNumber: %d", numberOfRankers, pos, epoch.BlockNumber)

	ranker.run(done)
	ranker.insertStoredComputationsInBackground()

	return ranker, nil
}

// InsertChange implements the Ranker interface.
func (ranker *delegateRanker) InsertChange(change orderbook.Change) {
	ranker.rankerMu.Lock()
	defer ranker.rankerMu.Unlock()

	log.Printf("[change detected] order %v status change to %v at block %d", base64.StdEncoding.EncodeToString(change.OrderID[:]), change.OrderStatus, change.BlockNumber)
	// FIXME : Change blockNumber can be different from the epoch blockNumber
	if change.BlockNumber >= ranker.rankerCurrBlockNum {
		select {
		case <-ranker.done:
		case ranker.rankerCurrEpochIn <- change:
			log.Println("pumping change to current epochRanker")
		}
		return
	}
	if change.BlockNumber >= ranker.rankerPrevBlockNum {
		select {
		case <-ranker.done:
		case ranker.rankerPrevEpochIn <- change:
		}
		return
	}
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

	if epoch.BlockNumber == ranker.rankerCurrBlockNum {
		return
	}

	if ranker.rankerPrevEpoch != nil {
		close(ranker.rankerPrevEpochIn)
	}
	ranker.rankerPrevEpoch = ranker.rankerCurrEpoch
	ranker.rankerPrevEpochIn = ranker.rankerCurrEpochIn
	ranker.rankerPrevEpochOut = ranker.rankerCurrEpochOut
	ranker.rankerPrevBlockNum = ranker.rankerCurrBlockNum

	numberOfRankers, pos, err := ranker.getPosFromEpoch(epoch)
	if err != nil {
		logger.Error(fmt.Sprintf("cannot get ranker position from epoch: %v", err))
		return
	}
	ranker.rankerCurrEpoch = newEpochRanker(numberOfRankers, pos)
	ranker.rankerCurrEpochIn = make(chan orderbook.Change)
	ranker.rankerCurrEpochOut = ranker.rankerCurrEpoch.run(ranker.done, ranker.rankerCurrEpochIn)
	ranker.rankerCurrBlockNum = epoch.BlockNumber
	log.Printf("[epoch change] rankers: %d, pos : %d, blockNumber: %d", numberOfRankers, pos, epoch.BlockNumber)
}

func (ranker *delegateRanker) run(done <-chan struct{}) {
	go func() {
		for {
			ranker.rankerMu.Lock()
			currEpochRankerCh := ranker.rankerCurrEpochOut
			prevEpochRankerCh := ranker.rankerPrevEpochOut
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

		log.Printf("load %d computations from local storage.", len(coms))
		for i := range coms {
			log.Printf("computation %v with state %v", base64.StdEncoding.EncodeToString(coms[i].ID[:]), coms[i].State)
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
		defer close(computations)

		for change := range changes {
			switch change.OrderStatus {
			case order.Open:
				if change.OrderParity == order.ParityBuy {
					select {
					case <-done:
						return
					case computations <- ranker.insertBuy(change):
						log.Println("inserting buy into the epochRanker")
					}
				} else {
					select {
					case <-done:
						return
					case computations <- ranker.insertSell(change):
						log.Println("inserting sell into the epochRanker")
					}
				}
			case order.Canceled, order.Confirmed:
				ranker.remove(change)
			}
		}
	}()

	return computations
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
	log.Printf("return %d computations after inserting the buy order %v", len(computations), base64.StdEncoding.EncodeToString(change.OrderID[:]))

	return computations
}

func (ranker *epochRanker) insertSell(change orderbook.Change) []Computation {
	computations := make([]Computation, 0)
	ranker.sells[change.OrderID] = change.OrderPriority
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
	log.Printf("return %d computations after inserting the sell order %v", len(computations), base64.StdEncoding.EncodeToString(change.OrderID[:]))
	return computations
}

func (ranker *epochRanker) remove(change orderbook.Change) {
	if change.OrderParity == order.ParityBuy {
		delete(ranker.buys, change.OrderID)
	} else {
		delete(ranker.sells, change.OrderID)
	}
}
