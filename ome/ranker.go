package ome

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/logger"

	"github.com/republicprotocol/republic-go/order"
)

// A PriorityOrder is Priority coupled with an order.Order.
type PriorityOrder struct {
	Priority Priority
	Order    order.ID
}

// A Ranker consumes orders and produces Computations that are prioritized
// based on the combined priorities of the involved orders.
type Ranker interface {

	// InsertBuy order.Order into the Ranker. A call to Ranker.InsertBuy will
	// combine the newly inserted order.Order with others to create pending
	// Computations. These can be read from the Ranker using a call to
	// Ranker.Computations.
	InsertBuy(PriorityOrder)

	// InsertSell order.Order into the Ranker. A call to Ranker.InsertBuy will
	// combine the newly inserted order.Order with others to create pending
	// Computations. These can be read from the Ranker using a call to
	// Ranker.Computations.
	InsertSell(PriorityOrder)

	// Remove orders from the Ranker. This will also remove all Computations
	// that involve these orders.
	Remove(...order.ID)

	// Computations stored in the Ranker are written to the input buffer. The
	// written Computations are removed from the Ranker to prevent duplicate
	// Computations. Returns the number of Computations written to the buffer,
	// which is guaranteed to be less than, or equal to, the size of the
	// buffer.
	Computations(Computations) int
}

type ranker struct {
	numberOfRankers int
	pos             int

	computationsMu *sync.Mutex
	computations   []Computation
	buys           map[order.ID]Priority
	sells          map[order.ID]Priority

	storer Storer
}

// NewRanker returns a Ranker that first filters the Computations it produces
// by checking the Priority. The filter assumes that there are a certain number
// of Rankers, and that each Ranker has a unique position relative to others.
// Priorities that do not match the position of the Ranker, after a modulo of
// the number of Rankers, are filtered. A Storer is used to load existing
// Computations that have not been processed completely, and to store new
// Computations.
func NewRanker(numberOfRankers, pos int, storer Storer) Ranker {
	ranker := &ranker{
		numberOfRankers: numberOfRankers,
		pos:             pos,

		computationsMu: new(sync.Mutex),
		computations:   []Computation{},
		buys:           map[order.ID]Priority{},
		sells:          map[order.ID]Priority{},

		storer: storer,
	}
	ranker.insertStoredComputationsInBackground()
	return ranker
}

// InsertBuy implements the Ranker interface.
func (ranker *ranker) InsertBuy(priorityOrder PriorityOrder) {
	ranker.computationsMu.Lock()
	defer ranker.computationsMu.Unlock()

	ranker.buys[priorityOrder.Order] = priorityOrder.Priority
	for sell, sellPriority := range ranker.sells {

		priority := priorityOrder.Priority + sellPriority
		if int(priority)%ranker.numberOfRankers != ranker.pos {
			continue
		}

		priorityCom := NewComputation(priorityOrder.Order, sell)
		priorityCom.Priority = priority

		ranker.insertComputation(priorityCom)
		if err := ranker.storer.InsertComputation(priorityCom); err != nil {
			logger.Error(fmt.Sprintf("cannot store computation buy = %v, sell = %v", priorityCom.Buy, priorityCom.Sell))
		}
	}
}

// InsertSell implements the Ranker interface.
func (ranker *ranker) InsertSell(priorityOrder PriorityOrder) {
	ranker.computationsMu.Lock()
	defer ranker.computationsMu.Unlock()

	ranker.sells[priorityOrder.Order] = priorityOrder.Priority
	for buy, buyPriority := range ranker.buys {

		priority := priorityOrder.Priority + buyPriority
		if int(priority)%ranker.numberOfRankers != ranker.pos {
			continue
		}

		priorityCom := NewComputation(buy, priorityOrder.Order)
		priorityCom.Priority = priority

		ranker.insertComputation(priorityCom)
		if err := ranker.storer.InsertComputation(priorityCom); err != nil {
			logger.Error(fmt.Sprintf("cannot store computation buy = %v, sell = %v", priorityCom.Buy, priorityCom.Sell))
		}
	}
}

// Remove implements the Ranker interface.
func (ranker *ranker) Remove(orders ...order.ID) {
	ranker.computationsMu.Lock()
	defer ranker.computationsMu.Unlock()

	mapping := map[order.ID]struct{}{}
	for _, order := range orders {
		mapping[order] = struct{}{}
		delete(ranker.buys, order)
		delete(ranker.sells, order)
	}

	for i := 0; i < len(ranker.computations); i++ {
		if _, ok := mapping[ranker.computations[i].Buy]; ok {
			ranker.computations = append(ranker.computations[:i], ranker.computations[i+1:]...)
			i--
			continue
		}
		if _, ok := mapping[ranker.computations[i].Sell]; ok {
			ranker.computations = append(ranker.computations[:i], ranker.computations[i+1:]...)
			i--
			continue
		}
	}
}

// Computations implements the Ranker interface.
func (ranker *ranker) Computations(buffer Computations) int {
	ranker.computationsMu.Lock()
	defer ranker.computationsMu.Unlock()

	n := 0
	for i := 0; i < len(buffer) && i < len(ranker.computations); i++ {
		buffer[i] = ranker.computations[i]
		n++
	}

	if n >= len(ranker.computations) {
		ranker.computations = ranker.computations[0:0]
	} else {
		ranker.computations = ranker.computations[n:]
	}
	return n
}

func (ranker *ranker) insertStoredComputationsInBackground() {
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

func (ranker *ranker) insertComputation(com Computation) {
	index := sort.Search(len(ranker.computations), func(i int) bool {
		return ranker.computations[i].Priority > com.Priority
	})
	ranker.computations = append(
		ranker.computations[:index],
		append([]Computation{com}, ranker.computations[index:]...)...)
}
