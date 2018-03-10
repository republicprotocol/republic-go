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
	firstTimers map[string]*time.Timer
	bestMatch   map[string]order.ID
	queue       chan *compute.Delta
}

func NewGossipWorker(queue chan *compute.Delta) *GossipWorker {
	return &GossipWorker{
		firstTimers: make(map[string]*time.Timer),
		bestMatch:   make(map[string]order.ID),
		queue:       queue,
	}
}

// Starts timers for each new id
func (worker *GossipWorker) Run(queues ...chan *compute.Delta) {
	select {
	case new := <-worker.queue:
		// Set up timers
		if worker.firstTimers[string(new.BuyOrderID)] == nil {
			worker.firstTimers[string(new.BuyOrderID)] = time.NewTimer(1 * time.Second)
		}
		if worker.firstTimers[string(new.SellOrderID)] == nil {
			worker.firstTimers[string(new.SellOrderID)] = time.NewTimer(1 * time.Second)
		}

		// TODO: Gossip to others
		// rpc...
		// rpc.NewGossipClient()

		worker.bestMatch = /* calculate new best match */ worker.bestMatch
	}
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
