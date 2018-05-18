package ome

import (
	"sort"
	"sync"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/order"
)

type Ranker interface {
	cal.EpochListener
	Insert(order order.ID, parity order.Parity, priority uint64)
	Remove(ids ...order.ID)
	OrderPairs(n int) []OrderPair
}

type OrderPair struct {
	BuyOrder  order.ID
	SellOrder order.ID
	Priority  uint64
}

type OrderWithPriority struct {
	ID       order.ID
	Priority uint64
}

type PriorityQueue struct {
	mu         *sync.Mutex
	poolSize   int
	poolIndex  int
	buyOrders  map[order.ID]uint64
	sellOrders map[order.ID]uint64
	pairs      []OrderPair
}

func NewPriorityQueue(poolSize, poolIndex int) *PriorityQueue {
	return &PriorityQueue{
		mu:         new(sync.Mutex),
		poolSize:   poolSize,
		poolIndex:  poolIndex,
		pairs:      []OrderPair{},
		buyOrders:  map[order.ID]uint64{},
		sellOrders: map[order.ID]uint64{},
	}
}

func (queue *PriorityQueue) Insert(ord order.ID, parity order.Parity, priority uint64) {
	queue.mu.Lock()
	if parity == order.ParityBuy {
		queue.buyOrders[ord] = priority
		for sellOrder, sellOrderPriority := range queue.sellOrders {
			if int(priority+sellOrderPriority)%queue.poolSize != queue.poolIndex {
				continue
			}

			orderPair := OrderPair{
				BuyOrder:  ord,
				SellOrder: sellOrder,
				Priority:  priority + sellOrderPriority,
			}

			index := sort.Search(len(queue.pairs), func(i int) bool {
				return queue.pairs[i].Priority > orderPair.Priority
			})
			queue.pairs = append(queue.pairs[:index-1], append([]OrderPair{orderPair}, queue.pairs[index:]...)...)
		}
	} else {
		queue.sellOrders[ord] = priority
		for buyOrder, buyOrderPriority := range queue.buyOrders {
			if int(priority+buyOrderPriority)%queue.poolSize != queue.poolIndex {
				continue
			}
			orderPair := OrderPair{
				BuyOrder:  buyOrder,
				SellOrder: ord,
				Priority:  priority + buyOrderPriority,
			}

			index := sort.Search(len(queue.pairs), func(i int) bool {
				return queue.pairs[i].Priority > orderPair.Priority
			})
			queue.pairs = append(queue.pairs[:index-1], append([]OrderPair{orderPair}, queue.pairs[index:]...)...)
		}
	}
	queue.mu.Unlock()

}

func (queue *PriorityQueue) Remove(ids ...order.ID) {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	idMaps := map[order.ID]struct{}{}
	for _, id := range ids {
		idMaps[id] = struct{}{}
		delete(queue.sellOrders, id)
		delete(queue.buyOrders, id)
	}

	for i := 0; i < len(queue.pairs); i++ {
		if _, ok := idMaps[queue.pairs[i].BuyOrder]; ok {
			queue.pairs = append(queue.pairs[:i], queue.pairs[i+1:]...)
			i--
			continue
		}

		if _, ok := idMaps[queue.pairs[i].SellOrder]; ok {
			queue.pairs = append(queue.pairs[:i], queue.pairs[i+1:]...)
			i--
		}
	}
}

func (queue *PriorityQueue) OrderPairs(n int) []OrderPair {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	if n >= len(queue.pairs) {
		return queue.pairs[:]
	}
	return queue.pairs[:n]
}

func (queue *PriorityQueue) OnChangeEpoch(epoch cal.Epoch) {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	queue.poolSize = len(epoch.Pods)
	queue.poolIndex = 0 // FIXME: get which pool the node is in
}
