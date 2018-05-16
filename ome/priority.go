package ome

import (
	"sort"
	"sync"

	"github.com/republicprotocol/republic-go/order"
)

type Ranker interface {
	Insert(pair OrderPair)
	Remove(id order.ID)
	Get(number int) []OrderPair
}

type PriorityQueue struct {
	mu    *sync.Mutex
	pairs []OrderPair
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		mu:    new(sync.Mutex),
		pairs: []OrderPair{},
	}
}

func (queue *PriorityQueue) Insert(pair OrderPair) {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	index := sort.Search(len(queue.pairs), func(i int) bool {
		return queue.pairs[i].priority > pair.priority
	})
	queue.pairs = append(queue.pairs[:index-1], append([]OrderPair{pair}, queue.pairs[index:]...)...)
}

func (queue *PriorityQueue) Remove(id order.ID) {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	for i := 0; i < len(queue.pairs); i++ {
		remove := false
		if queue.pairs[i].orderID == id {
			remove = true
		}
		for _, match := range queue.pairs[i].matches {
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

func (queue *PriorityQueue) Get(number int) []OrderPair {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	defer func() {
		queue.pairs = queue.pairs[number:]
	}()

	return queue.pairs[:number]

}
