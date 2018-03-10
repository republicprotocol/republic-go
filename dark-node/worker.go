package node

import (
	"fmt"
	"time"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/compute"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
)

// An OrderFragmentWorker consumes order fragments and computes all
// combinations of delta fragments.
type OrderFragmentWorker struct {
	queue               chan *order.Fragment
	deltaFragmentMatrix *compute.DeltaFragmentMatrix
}

// NewOrderFragmentWorker returns an OrderFragmentWorker that reads work from
// a queue and uses a DeltaFragmentMatrix to do computations.
func NewOrderFragmentWorker(queue chan *order.Fragment, deltaFragmentMatrix *compute.DeltaFragmentMatrix) *OrderFragmentWorker {
	return &OrderFragmentWorker{
		queue:               queue,
		deltaFragmentMatrix: deltaFragmentMatrix,
	}
}

// Run the OrderFragmentWorker and write all delta fragments to an output
// queue.
func (worker *OrderFragmentWorker) Run(queues ...chan *compute.DeltaFragment) error {
	for orderFragment := range worker.queue {
		deltaFragments, err := worker.deltaFragmentMatrix.InsertOrderFragment(orderFragment)
		if err != nil {
			return err
		}
		if deltaFragments != nil {
			// Write to channels that might be closed
			func() {
				defer func() { recover() }()
				for _, deltaFragment := range deltaFragments {
					for _, queue := range queues {
						queue <- deltaFragment
					}
				}
			}()
		}
	}
	return nil
}

type DeltaFragmentBroadcastWorker struct {
	queue      chan *compute.DeltaFragment
	clientPool *rpc.ClientPool
	darkPool   identity.MultiAddresses
	logger     *logger.Logger
}

func NewDeltaFragmentBroadcastWorker(logger *logger.Logger, queue chan *compute.DeltaFragment, clientPool *rpc.ClientPool, darkPool identity.MultiAddresses) *DeltaFragmentBroadcastWorker {
	return &DeltaFragmentBroadcastWorker{
		logger:     logger,
		queue:      queue,
		clientPool: clientPool,
		darkPool:   darkPool,
	}
}

func (worker *DeltaFragmentBroadcastWorker) Run() {
	for deltaFragment := range worker.queue {
		serializedDeltaFragment := rpc.SerializeDeltaFragment(deltaFragment)
		do.CoForAll(worker.darkPool, func(i int) {
			_, err := worker.clientPool.BroadcastDeltaFragment(worker.darkPool[i], serializedDeltaFragment)
			if err != nil {
				worker.logger.Error(logger.TagNetwork, err.Error())
			}
		})
	}
}

// An DeltaFragmentWorker consumes delta fragments and reconstructs deltas.
type DeltaFragmentWorker struct {
	queue        chan *compute.DeltaFragment
	deltaBuilder *compute.DeltaBuilder
}

// NewDeltaFragmentWorker returns an DeltaFragmentWorker that reads work from
// a queue and uses a DeltaBuilder to do reconstructions.
func NewDeltaFragmentWorker(queue chan *compute.DeltaFragment, deltaBuilder *compute.DeltaBuilder) *DeltaFragmentWorker {
	return &DeltaFragmentWorker{
		queue:        queue,
		deltaBuilder: deltaBuilder,
	}
}

// Run the DeltaFragmentWorker and write all deltas to  an output queue.
func (worker *DeltaFragmentWorker) Run(queues ...chan *compute.Delta) {
	for deltaFragment := range worker.queue {
		delta := worker.deltaBuilder.InsertDeltaFragment(deltaFragment)
		if delta != nil {
			// Write to channels that might be closed
			func() {
				defer func() { recover() }()
				for _, queue := range queues {
					queue <- delta
				}
			}()
		}
	}
}

type GossipWorker struct {
	expiryTime map[string]time.Time
	bestMatch  map[string]*compute.Delta
	queue      chan *compute.Delta
	clientPool *rpc.ClientPool
	gossipers  identity.MultiAddresses
}

func NewGossipWorker(clientPool *rpc.ClientPool, gossipers identity.MultiAddresses, queue chan *compute.Delta) *GossipWorker {
	return &GossipWorker{
		expiryTime: make(map[string]time.Time),
		bestMatch:  make(map[string]*compute.Delta),
		queue:      queue,
		clientPool: clientPool,
		gossipers:  gossipers,
	}
}

// Starts timers for each new id
func (worker *GossipWorker) Run(queues ...chan *compute.Delta) {
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()
	for {
		select {
		case newDelta := <-worker.queue:
			// Set up timers
			if worker.expiryTime[string(newDelta.BuyOrderID)] == (time.Time{}) {
				worker.expiryTime[string(newDelta.BuyOrderID)] = time.Now().Add(10 * time.Second)
			}
			if worker.expiryTime[string(newDelta.SellOrderID)] == (time.Time{}) {
				worker.expiryTime[string(newDelta.SellOrderID)] = time.Now().Add(10 * time.Second)
			}

			previousBuyBest := worker.bestMatch[string(newDelta.BuyOrderID)]
			previousSellBest := worker.bestMatch[string(newDelta.SellOrderID)]
			newBuyBest := bestFitDelta(newDelta, previousBuyBest)
			newSellBest := bestFitDelta(newDelta, previousSellBest)
			worker.bestMatch[string(newDelta.BuyOrderID)] = newBuyBest
			worker.bestMatch[string(newDelta.SellOrderID)] = newSellBest

			if string(newBuyBest.ID) != string(previousBuyBest.ID) ||
				string(newSellBest.ID) != string(previousSellBest.ID) {

				// Gossip to others.
				for _, multi := range worker.gossipers {
					worker.clientPool.Gossip(multi, &rpc.Rumor{
						BuyOrderId:  newDelta.BuyOrderID,
						SellOrderId: newDelta.SellOrderID,
					})
				}
			}
		case now := <-timer.C:
			for k, v := range worker.expiryTime {
				if v.Unix() < now.Unix() { // TODO: How to better compare time
					// Safe to do in loop
					// https://golang.org/doc/effective_go.html#for

					for _, queue := range queues {
						queue <- worker.bestMatch[k]
					}

					// TODO: If we receive a new delta for k, it's timer will start again
					delete(worker.expiryTime, k)
					delete(worker.bestMatch, k)
				}
			}
		}
	}
}

func bestFitDelta(left, right *compute.Delta) *compute.Delta {
	if left.ID.LessThan(right.ID) {
		return left
	}
	return right
}

type Data struct {
	Delta *compute.Delta
	Vote  int
	Sent  bool
}

type FinalizeWorker struct {
	queue    chan *compute.Delta
	deltas   map[string]*Data
	poolSize int
}

func NewFinalizeWorker(queue chan *compute.Delta, poolSize int) *FinalizeWorker {
	return &FinalizeWorker{
		queue:    queue,
		deltas:   map[string]*Data{},
		poolSize: poolSize,
	}
}

func (worker *FinalizeWorker) Run(queues ...chan *compute.Delta) {
	for delta := range worker.queue {
		d, ok := worker.deltas[string(delta.ID)]
		if !ok {
			worker.deltas[string(delta.ID)] = &Data{
				Delta: delta,
				Vote:  1,
				Sent:  false,
			}
		} else {
			d.Vote += 1
		}

		if d.Vote > worker.poolSize/2 && !d.Sent {
			for _, queue := range queues {
				queue <- worker.deltas[string(delta.ID)].Delta
			}
			d.Sent = true
		}
	}
}

type ConsensusWorker struct {
	logger              *logger.Logger
	queue               chan *compute.Delta
	deltaFragmentMatrix *compute.DeltaFragmentMatrix
}

func NewConsensusWorker(logger *logger.Logger, queue chan *compute.Delta, deltaFragmentMatrix *compute.DeltaFragmentMatrix) *ConsensusWorker {
	return &ConsensusWorker{
		logger:              logger,
		queue:               queue,
		deltaFragmentMatrix: deltaFragmentMatrix,
	}
}

func (worker *ConsensusWorker) Run() {
	for delta := range worker.queue {
		worker.logger.Info(logger.TagConsensus, fmt.Sprintf("(%s, %s)", delta.BuyOrderID.String(), delta.SellOrderID.String()))
	}
}
