package node

import (
	"fmt"
	"time"

	"github.com/republicprotocol/republic-go/compute"
	"github.com/republicprotocol/republic-go/logger"
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

type FinalizeWorker struct {
}

func NewFinalizeWorker(queue chan *compute.Delta) *FinalizeWorker {
	return &FinalizeWorker{}
}

func (worker *FinalizeWorker) Run(queues ...chan *compute.Delta) {
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
