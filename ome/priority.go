package ome

import (
	"sort"
	"sync"

	"github.com/republicprotocol/republic-go/order"
)

type Ranker interface {
	Insert(order order.ID, parity order.Parity, priority uint64)
	Remove(id order.ID)
	OrderPairs(done <-chan struct{}) <-chan OrderPair
}

type OrderPair struct {
	BuyOrder  order.ID
	SellOrder order.ID
	Priority  uint64
}

type PriorityQueue struct {
	mu       *sync.Mutex
	pairs    []OrderPair
	newPairs chan OrderPair
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		mu:       new(sync.Mutex),
		pairs:    []OrderPair{},
		newPairs: make(chan OrderPair),
	}
}

func (queue *PriorityQueue) Insert(pair OrderPair) {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	index := sort.Search(len(queue.pairs), func(i int) bool {
		return queue.pairs[i].Priority > pair.Priority
	})
	queue.pairs = append(queue.pairs[:index-1], append([]OrderPair{pair}, queue.pairs[index:]...)...)
}

func (queue *PriorityQueue) Remove(id order.ID) {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	for i := 0; i < len(queue.pairs); i++ {
		remove := false
		if queue.pairs[i].ID == id {
			remove = true
		}
		for _, match := range queue.pairs[i].Matches {
			if match == id {
				remove = true
				break
			}
		}

		if remove == true {
			// a:= []{1,2,3}
			// log.Println(a[3])   # this will panic
			// log.Println(a[3:])  # this will not panic ,amazing!
			queue.pairs = append(queue.pairs[:i], queue.pairs[i+1:]...)
			i--
		}
	}
}

func (queue *PriorityQueue) OrderPairs() <-chan OrderPair {
	orderPairs := make(chan OrderPair)

	queue.mu.Lock()
	defer queue.mu.Unlock()

	defer func() {
		queue.pairs = queue.pairs[number:]
	}()

	return queue.pairs[:number]

	return orderPairs
}
