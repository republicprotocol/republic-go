package smpc

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
)

// Workers is a slice of Worker components.
type Workers []Worker

// A Worker receives messages from a Dispatcher until the Dispatcher is
// shutdown. It is primarily responsible for decoding the message and
// delegating work to the appropriate component.
type Worker struct {
	running int32
	logger  *logger.Logger

	messageQueuesMu *sync.RWMutex
	messageQueues   dispatch.MessageQueues

	multiplexer         *dispatch.Multiplexer
	deltaFragmentMatrix *DeltaFragmentMatrix
	deltaBuilder        *DeltaBuilder
	deltaQueue          *DeltaQueue
}

// NewWorker returns a new Worker that can handle all types of WorkerTasks. It
// will process WorkTasks serially, and if a new Message is required it will
// send it back through the Multiplexer for scheduling to another worker. This
// prevents new WorkerTasks from jumping the queue, providing a sense of
// fairness in prioritization.
func NewWorker(logger *logger.Logger, messageQueues dispatch.MessageQueues, multiplexer *dispatch.Multiplexer, deltaFragmentMatrix *DeltaFragmentMatrix, deltaBuilder *DeltaBuilder, deltaQueue *DeltaQueue) Worker {
	return Worker{
		running: 1,
		logger:  logger,

		messageQueuesMu: new(sync.RWMutex),
		messageQueues:   messageQueues,

		multiplexer:         multiplexer,
		deltaFragmentMatrix: deltaFragmentMatrix,
		deltaBuilder:        deltaBuilder,
		deltaQueue:          deltaQueue,
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
		case Message:
			if message.OrderFragment != nil {
				worker.processOrderFragment(message.OrderFragment)
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

func (worker *Worker) processOrderFragment(orderFragment *order.Fragment) {

	// Compute all new DeltaFragments
	deltaFragments := DeltaFragments{}
	if orderFragment.OrderParity == order.ParityBuy {
		deltaFragments = worker.deltaFragmentMatrix.ComputeBuyOrder(orderFragment)
	} else {
		deltaFragments = worker.deltaFragmentMatrix.ComputeSellOrder(orderFragment)
	}

	// Send a new Message directly to the Multiplexer so that the new
	// DeltaFragments can be processed
	if deltaFragments != nil && len(deltaFragments) > 0 {
		worker.multiplexer.Send(Message{DeltaFragments: deltaFragments})
	}
}

func (worker *Worker) processDeltaFragments(deltaFragments DeltaFragments) {
	if deltaFragments == nil || len(deltaFragments) == 0 {
		return
	}

	// Build new Deltas from the DeltaFragments
	newDeltas, newDeltaFragments := worker.deltaBuilder.ComputeDelta(deltaFragments)

	// Send a new Message directly to the Multiplexer so that the new
	// Deltas can be processed
	if newDeltas != nil && len(newDeltas) > 0 {
		worker.multiplexer.Send(Message{
			Deltas: newDeltas,
		})
	}

	if newDeltaFragments != nil && len(newDeltaFragments) > 0 {
		// Send a new Message to all MessageQueues available to this Worker
		worker.messageQueuesMu.RLock()
		defer worker.messageQueuesMu.RUnlock()

		for _, queue := range worker.messageQueues {
			queue.Send(Message{
				DeltaFragments: newDeltaFragments,
			})
		}
	}
}

func (worker *Worker) processDeltas(deltas Deltas) {
	for _, delta := range deltas {
		if err := worker.deltaQueue.Send(delta); err != nil {
			worker.logger.Compute(logger.Error, fmt.Sprintf("cannot send delta notification: %s", err.Error()))
		}
	}
}
