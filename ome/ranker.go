package ome

import (
	"sort"
	"sync"

	"github.com/republicprotocol/republic-go/order"
)

// A Priority is an unsigned integer representing logical time prioritization.
// The lower the number, the higher the priority.
type Priority uint64

// A PriorityOrder is an order.Order coupled with a Priority.
type PriorityOrder struct {
	Priority
	order.ID
}

// A Ranker consumes orders and produces computations that are prioritized
// based on the combined priorities of the orders.
type Ranker interface {

	// InsertBuy order into the Ranker. A call to Ranker.InsertBuy will produce
	// all new Computations which can be read from the Ranker using a call to
	// Ranker.Computations.
	InsertBuy(order PriorityOrder)

	// InsertSell order into the Ranker. A call to Ranker.InsertSell will
	// produce all new Computations which can be read from the Ranker using a
	// call to Ranker.Computations.
	InsertSell(order PriorityOrder)

	// Remove orders from the Ranker. This will also remove all Computations
	// that involve the orders.
	Remove(orders ...order.ID)

	// Computations stores in the Ranker are used to fill the buffer. The
	// Computations are then removed from the Ranker to prevent duplication in
	// future calls. Returns the number of Computations written to the buffer,
	// which is guaranteed to be less than, or equal to, the size of the
	// buffer.
	Computations(buffer Computations) int
}

type ranker struct {
	num int
	pos int

	computationsMu *sync.Mutex
	computations   Computations
	buys           map[order.ID]Priority
	sells          map[order.ID]Priority
}

// NewRanker returns a Ranker that only produces Computations based on the
// total number of Rankers processing orders, and the position of this Ranker.
func NewRanker(num, pos int) Ranker {
	return &ranker{
		num: num,
		pos: pos,

		computationsMu: new(sync.Mutex),
		computations:   Computations{},
		buys:           map[order.ID]Priority{},
		sells:          map[order.ID]Priority{},
	}
}

func (ranker *ranker) InsertBuy(order PriorityOrder) {
	ranker.computationsMu.Lock()
	defer ranker.computationsMu.Unlock()

	ranker.buys[order.ID] = order.Priority
	for sell, sellPriority := range ranker.sells {
		computationPriority := order.Priority + sellPriority
		if int(computationPriority)%ranker.num != ranker.pos {
			continue
		}

		computation := Computation{
			Buy:      order.ID,
			Sell:     sell,
			Priority: computationPriority,
		}

		index := sort.Search(len(ranker.computations), func(i int) bool {
			return ranker.computations[i].Priority > computation.Priority
		})
		ranker.computations = append(
			append(
				ranker.computations[:index-1],
				computation,
			),
			ranker.computations[index:]...,
		)
	}
}

func (ranker *ranker) InsertSell(order PriorityOrder) {
	ranker.computationsMu.Lock()
	defer ranker.computationsMu.Unlock()

	ranker.sells[order.ID] = order.Priority
	for buy, buyPriority := range ranker.buys {
		computationPriority := order.Priority + buyPriority
		if int(computationPriority)%ranker.num != ranker.pos {
			continue
		}

		computation := Computation{
			Buy:      buy,
			Sell:     order.ID,
			Priority: computationPriority,
		}

		index := sort.Search(len(ranker.computations), func(i int) bool {
			return ranker.computations[i].Priority > computation.Priority
		})
		ranker.computations = append(
			append(
				ranker.computations[:index-1],
				computation,
			),
			ranker.computations[index:]...,
		)
	}
}

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
