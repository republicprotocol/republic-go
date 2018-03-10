package node

import (
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

func (worker *DeltaFragmentBroadcastWorker) Run(queues ...chan *compute.Delta) {
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
