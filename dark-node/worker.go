package node

import (
	"fmt"
	"log"
	"time"

	"github.com/republicprotocol/republic-go/compute"
	"github.com/republicprotocol/republic-go/dark"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
)

// An OrderFragmentWorker consumes order fragments and computes all
// combinations of delta fragments.
type OrderFragmentWorker struct {
	logger              *logger.Logger
	deltaFragmentMatrix *compute.DeltaFragmentMatrix
	queue               chan *order.Fragment
}

// NewOrderFragmentWorker returns an OrderFragmentWorker that reads work from
// a queue and uses a DeltaFragmentMatrix to do computations.
func NewOrderFragmentWorker(logger *logger.Logger, deltaFragmentMatrix *compute.DeltaFragmentMatrix, queue chan *order.Fragment) *OrderFragmentWorker {
	return &OrderFragmentWorker{
		logger:              logger,
		deltaFragmentMatrix: deltaFragmentMatrix,
		queue:               queue,
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
	logger     *logger.Logger
	clientPool *rpc.ClientPool
	darkPool   *dark.Pool
	queue      chan *compute.DeltaFragment
}

func NewDeltaFragmentBroadcastWorker(logger *logger.Logger, clientPool *rpc.ClientPool, darkPool *dark.Pool, queue chan *compute.DeltaFragment) *DeltaFragmentBroadcastWorker {
	return &DeltaFragmentBroadcastWorker{
		logger:     logger,
		clientPool: clientPool,
		darkPool:   darkPool,
		queue:      queue,
	}
}

func (worker *DeltaFragmentBroadcastWorker) Run() {
	for deltaFragment := range worker.queue {
		serializedDeltaFragment := rpc.SerializeDeltaFragment(deltaFragment)
		log.Println(worker.darkPool)
		worker.darkPool.CoForAll(func(node *dark.Node) {
			multiAddress := node.MultiAddress()
			if multiAddress == nil {
				return
			}
			_, err := worker.clientPool.BroadcastDeltaFragment(*multiAddress, serializedDeltaFragment)
			if err != nil {
				worker.logger.Error(logger.TagNetwork, err.Error())
			}
		})
	}
}

// An DeltaFragmentWorker consumes delta fragments and reconstructs deltas.
type DeltaFragmentWorker struct {
	logger       *logger.Logger
	deltaBuilder *compute.DeltaBuilder
	queue        chan *compute.DeltaFragment
}

// NewDeltaFragmentWorker returns an DeltaFragmentWorker that reads work from
// a queue and uses a DeltaBuilder to do reconstructions.
func NewDeltaFragmentWorker(logger *logger.Logger, deltaBuilder *compute.DeltaBuilder, queue chan *compute.DeltaFragment) *DeltaFragmentWorker {
	return &DeltaFragmentWorker{
		logger:       logger,
		deltaBuilder: deltaBuilder,
		queue:        queue,
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
				// FIXME: Stop doing this shit.
				if delta.IsMatch(Prime) {
					worker.logger.Info(logger.TagCompute, fmt.Sprintf("(%s, %s) ✓", delta.BuyOrderID.String(), delta.SellOrderID.String()))
				} else {
					worker.logger.Info(logger.TagCompute, fmt.Sprintf("(%s, %s)", delta.BuyOrderID.String(), delta.SellOrderID.String()))
				}
				for _, queue := range queues {
					queue <- delta
				}
			}()
		}
	}
}

type GossipWorker struct {
	logger     *logger.Logger
	clientPool *rpc.ClientPool
	gossipers  identity.MultiAddresses
	queue      chan *compute.Delta

	expiryTime map[string]time.Time
	bestMatch  map[string]*compute.Delta
}

func NewGossipWorker(logger *logger.Logger, clientPool *rpc.ClientPool, gossipers identity.MultiAddresses, queue chan *compute.Delta) *GossipWorker {
	return &GossipWorker{
		logger:     logger,
		clientPool: clientPool,
		gossipers:  gossipers,
		queue:      queue,

		expiryTime: make(map[string]time.Time),
		bestMatch:  make(map[string]*compute.Delta),
	}
}

// Starts timers for each new id
func (worker *GossipWorker) Run(queues ...chan *compute.Delta) {
	// timer := time.NewTimer(5 * time.Second)
	// defer timer.Stop()
	// for {
	// 	select {
	// 	case newDelta := <-worker.queue:
	// 		// Set up timers
	// 		if worker.expiryTime[string(newDelta.BuyOrderID)] == (time.Time{}) {
	// 			worker.expiryTime[string(newDelta.BuyOrderID)] = time.Now().Add(10 * time.Second)
	// 		}
	// 		if worker.expiryTime[string(newDelta.SellOrderID)] == (time.Time{}) {
	// 			worker.expiryTime[string(newDelta.SellOrderID)] = time.Now().Add(10 * time.Second)
	// 		}

	// 		previousBuyBest := worker.bestMatch[string(newDelta.BuyOrderID)]
	// 		previousSellBest := worker.bestMatch[string(newDelta.SellOrderID)]
	// 		newBuyBest := bestFitDelta(newDelta, previousBuyBest)
	// 		newSellBest := bestFitDelta(newDelta, previousSellBest)
	// 		worker.bestMatch[string(newDelta.BuyOrderID)] = newBuyBest
	// 		worker.bestMatch[string(newDelta.SellOrderID)] = newSellBest

	// 		if string(newBuyBest.ID) != string(previousBuyBest.ID) ||
	// 			string(newSellBest.ID) != string(previousSellBest.ID) {

	// 			// Gossip to others.
	// 			for _, multi := range worker.gossipers {
	// 				worker.clientPool.Gossip(multi, &rpc.Rumor{
	// 					BuyOrderId:  newDelta.BuyOrderID,
	// 					SellOrderId: newDelta.SellOrderID,
	// 				})
	// 			}
	// 		}
	// 	case now := <-timer.C:
	// 		for k, v := range worker.expiryTime {
	// 			if v.Unix() < now.Unix() { // TODO: How to better compare time
	// 				// Safe to do in loop
	// 				// https://golang.org/doc/effective_go.html#for

	// 				for _, queue := range queues {
	// 					queue <- worker.bestMatch[k]
	// 				}

	// 				// TODO: If we receive a new delta for k, it's timer will start again
	// 				delete(worker.expiryTime, k)
	// 				delete(worker.bestMatch, k)
	// 			}
	// 		}
	// 	}
	// }
}

func bestFitDelta(left, right *compute.Delta) *compute.Delta {
	if left.ID.LessThan(right.ID) {
		return left
	}
	return right
}

type Data struct {
	Delta *compute.Delta
	Vote  int64
	Sent  bool
}

type FinalizeWorker struct {
	logger   *logger.Logger
	queue    chan *compute.Delta
	poolSize int64

	deltas map[string]*Data
}

func NewFinalizeWorker(logger *logger.Logger, poolSize int64, queue chan *compute.Delta) *FinalizeWorker {
	return &FinalizeWorker{
		logger:   logger,
		poolSize: poolSize,
		queue:    queue,

		deltas: map[string]*Data{},
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
	deltaFragmentMatrix *compute.DeltaFragmentMatrix
	queue               chan *compute.Delta
}

func NewConsensusWorker(logger *logger.Logger, deltaFragmentMatrix *compute.DeltaFragmentMatrix, queue chan *compute.Delta) *ConsensusWorker {
	return &ConsensusWorker{
		logger:              logger,
		deltaFragmentMatrix: deltaFragmentMatrix,
		queue:               queue,
	}
}

func (worker *ConsensusWorker) Run() {
	for delta := range worker.queue {
		worker.logger.Info(logger.TagConsensus, fmt.Sprintf("(%s, %s)", delta.BuyOrderID.String(), delta.SellOrderID.String()))
	}
}
