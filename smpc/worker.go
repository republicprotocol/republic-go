package smpc

import (
	"log"
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

	peerQueues dispatch.MessageQueues

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
func NewWorker(logger *logger.Logger, peerQueues dispatch.MessageQueues, multiplexer *dispatch.Multiplexer, deltaFragmentMatrix *DeltaFragmentMatrix, deltaBuilder *DeltaBuilder, deltaQueue *DeltaQueue) Worker {
	return Worker{
		running: 1,
		logger:  logger,

		peerQueues: peerQueues,

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
	log.Println("Run")
	for atomic.LoadInt32(&worker.running) > 0 {
		message, ok := worker.multiplexer.Recv()
		if !ok {
			break
		}
		switch message := message.(type) {
		case Message:
			if message.Error != nil {
				log.Println(message.Error)
			}
			if message.OrderFragment != nil {
				worker.processOrderFragment(message.OrderFragment)
			}
			if message.DeltaFragments != nil {
				worker.processDeltaFragments(message.DeltaFragments)
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
	if orderFragment.OrderParity == order.ParityBuy {
		worker.deltaFragmentMatrix.ComputeBuyOrder(orderFragment)
	} else {
		worker.deltaFragmentMatrix.ComputeSellOrder(orderFragment)
	}
}

func (worker *Worker) processDeltaFragments(deltaFragments DeltaFragments) {
	if deltaFragments == nil || len(deltaFragments) == 0 {
		return
	}

	// Build new Deltas from the DeltaFragments
	worker.deltaBuilder.ComputeDelta(deltaFragments)
}

type Broadcasters []Broadcaster

type Broadcaster struct {
	running int32
	logger  *logger.Logger

	peerQueues dispatch.MessageQueues

	deltaFragmentMatrix *DeltaFragmentMatrix
	deltaBuilder        *DeltaBuilder
	deltaQueue          *DeltaQueue
}

func NewBroadcaster(logger *logger.Logger, peerQueues dispatch.MessageQueues, deltaFragmentMatrix *DeltaFragmentMatrix, deltaBuilder *DeltaBuilder, deltaQueue *DeltaQueue) Broadcaster {
	return Broadcaster{
		running: 1,
		logger:  logger,

		peerQueues: peerQueues,

		deltaFragmentMatrix: deltaFragmentMatrix,
		deltaBuilder:        deltaBuilder,
		deltaQueue:          deltaQueue,
	}
}

func (broadcaster *Broadcaster) Run() {
	go func() {
		for atomic.LoadInt32(&broadcaster.running) != 0 {
			deltaFragments := [128]DeltaFragment{}
			n := broadcaster.deltaFragmentMatrix.WaitForDeltaFragments(deltaFragments[:])
			for _, queue := range broadcaster.peerQueues {
				queue.Send(Message{
					DeltaFragments: deltaFragments[:n],
				})
			}
		}
	}()

	deltas := [128]Delta{}
	for atomic.LoadInt32(&broadcaster.running) != 0 {
		n := broadcaster.deltaBuilder.WaitForDeltas(deltas[:])
		for i := 0; i < n; i++ {
			broadcaster.deltaQueue.Send(deltas[i])
		}
	}
}

func (broadcaster *Broadcaster) Shutdown() {
	atomic.StoreInt32(&broadcaster.running, 0)
}
