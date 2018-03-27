package smpc

import (
	"fmt"
	"sync/atomic"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
)

// A WorkerTask is recevied by workers to handle. It represents a task, or
// multiple tasks, that need to be completed.
type WorkerTask struct {
	OrderFragment  *rpc.OrderFragment
	TauMessage     *rpc.TauMessage
	DeltaFragments DeltaFragments
	Deltas         Deltas
}

// Workers is a slice of Worker components.
type Workers []Worker

// A Worker receives messages from a Dispatcher until the Dispatcher is
// shutdown. It is primarily responsible for decoding the message and
// delegating work to the appropriate component.
type Worker struct {
	logger  *logger.Logger
	running int32

	multiplexer   *dispatch.Multiplexer
	messageQueues dispatch.MessageQueues

	deltaFragmentMatrix *DeltaFragmentMatrix
	deltaBuilder        *DeltaBuilder
	deltaHandler        DeltaHandler
}

// NewWorker returns a new Worker that can handle all types of WorkerTasks. It
// will process WorkTasks serially, and if a new WorkerTask is required it will
// send it back through the Multiplexer for scheduling to another worker. This
// prevents new WorkerTasks from jumping the queue, providing a sense of
// fairness in prioritization.
func NewWorker(logger *logger.Logger, peers identity.MultiAddresses, multiplexer *dispatch.Multiplexer, messageQueues dispatch.MessageQueues, deltaFragmentMatrix *DeltaFragmentMatrix, deltaBuilder *DeltaBuilder, deltaHandler DeltaHandler) Worker {
	return Worker{
		logger:  logger,
		running: 1,

		multiplexer:   multiplexer,
		messageQueues: messageQueues,

		deltaFragmentMatrix: deltaFragmentMatrix,
		deltaBuilder:        deltaBuilder,
		deltaHandler:        deltaHandler,
	}
}

// Run until the Multiplexer is shutdown. The worker will read a message from
// the Multiplexer, delegate work to the appropriate component, wait for the
// component to complete, and then read the next message from the Multiplexer.
// This function blocks until the Multiplexer is shutdown.
func (worker *Worker) Run() {
	for atomic.LoadInt32(&worker.running) > 0 {
		message, ok := worker.multiplexer.Recv()
		if !ok {
			break
		}
		switch message := message.(type) {
		case WorkerTask:
			if message.OrderFragment != nil {
				worker.processOrderFragment(message.OrderFragment)
			}
			if message.TauMessage != nil {
				worker.processTauMessage(message.TauMessage)
			}
			if message.DeltaFragments != nil {
				worker.processDeltaFragments(message.DeltaFragments)
			}
			if message.Deltas != nil {
				worker.processDeltas(message.Deltas)
			}
		default:
			// Ignore message that we do not recognize
			break
		}
	}
}

// Shutdown the Worker gracefully.
func (worker *Worker) Shutdown() {
	atomic.StoreInt32(&worker.running, 0)
}

func (worker *Worker) processOrderFragment(orderFragment *rpc.OrderFragment) {
	fragment, err := rpc.DeserializeOrderFragment(orderFragment)
	if err != nil {
		worker.logger.Compute(logger.Error, fmt.Sprintf("cannot deserialize order: %s", err.Error()))
		return
	}

	// Compute all new DeltaFragments
	deltaFragments := DeltaFragments{}
	if fragment.OrderParity == order.ParityBuy {
		deltaFragments = worker.deltaFragmentMatrix.ComputeBuyOrder(fragment)
	} else {
		deltaFragments = worker.deltaFragmentMatrix.ComputeSellOrder(fragment)
	}

	// Send a new WorkerTask directly to the Multiplexer so that the new
	// DeltaFragments can be processed
	if deltaFragments != nil {
		worker.multiplexer.Send(WorkerTask{DeltaFragments: deltaFragments})
	}
}

func (worker *Worker) processTauMessage(message *rpc.TauMessage) {
	if message.GenerateRandomShares != nil {
		worker.processGenerateRandomShares(message.GenerateRandomShares)
	}
	if message.GenerateXiShares != nil {
		worker.processGenerateXiShares(message.GenerateXiShares)
	}
	if message.GenerateXiFragments != nil {
		worker.processGenerateXiFragment(message.GenerateXiFragments)
	}
	if message.RhoSigmaFragments != nil {
		worker.processBroadcastRhoSigmaFragment(message.RhoSigmaFragments)
	}
	if message.DeltaFragments != nil {
		deltaFragments := DeltaFragments{}
		if err := deltaFragments.Unmarshal(message.DeltaFragments); err != nil {

		}
		worker.processDeltaFragments(deltaFragments)
	}
}

func (worker *Worker) processDeltaFragments(deltaFragments DeltaFragments) {
	// Build new Deltas from the DeltaFragments
	newDeltas, newDeltaFragments := worker.deltaBuilder.ComputeDelta(deltaFragments)
	newDeltaFragmentsSerialized := make([]*rpc.DeltaFragment, len(newDeltaFragments))
	for i := range newDeltaFragmentsSerialized {
		newDeltaFragmentsSerialized[i] = newDeltaFragments[i].Marshal()
	}

	// Send a new WorkerTask directly to the Multiplexer so that the new
	// Deltas can be processed
	worker.multiplexer.Send(WorkerTask{
		Deltas: newDeltas,
	})

	// Send a new WorkerTask to a random subset of MessageQueues in the
	// Multiplexer
	for _, messageQueue := range worker.messageQueues {
		messageQueue.Send(WorkerTask{
			TauMessage: &rpc.TauMessage{
				DeltaFragments: &rpc.DeltaFragments{
					DeltaFragments: newDeltaFragmentsSerialized,
				},
			},
		})
	}
}

func (worker *Worker) processDeltas(deltas Deltas) {
	worker.deltaHandler(deltas)
}

func (worker *Worker) processGenerateRandomShares(request *rpc.GenerateRandomShares) {
	panic("unimplemented")
}

func (worker *Worker) processGenerateXiShares(request *rpc.GenerateXiShares) {
	panic("unimplemented")
}

func (worker *Worker) processGenerateXiFragment(request *rpc.GenerateXiFragments) {
	panic("unimplemented")
}

func (worker *Worker) processBroadcastRhoSigmaFragment(request *rpc.RhoSigmaFragments) {
	panic("unimplemented")
}
