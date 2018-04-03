package smpc

import (
	"log"
	"sync/atomic"
	"time"

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
	debug   bool
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
func NewWorker(debug bool, logger *logger.Logger, peerQueues dispatch.MessageQueues, multiplexer *dispatch.Multiplexer, deltaFragmentMatrix *DeltaFragmentMatrix, deltaBuilder *DeltaBuilder, deltaQueue *DeltaQueue) Worker {
	return Worker{
		debug:   debug,
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
	if worker.debug {
		log.Printf("%p starting", worker)
		defer log.Printf("%p shutting down", worker)
	}

	for atomic.LoadInt32(&worker.running) > 0 {
		if worker.debug {
			log.Printf("%p is recving", worker)
		}
		message, ok := worker.multiplexer.Recv()
		if !ok {
			break
		}
		if worker.debug {
			log.Printf("%p recvd", worker)
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
			log.Fatalf("unrecognized message type %T", message)
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
	debug   bool
	running int32
	logger  *logger.Logger

	peerQueues dispatch.MessageQueues

	deltaFragmentMatrix *DeltaFragmentMatrix
	deltaBuilder        *DeltaBuilder
	deltaQueue          *DeltaQueue
}

func NewBroadcaster(debug bool, logger *logger.Logger, peerQueues dispatch.MessageQueues, deltaFragmentMatrix *DeltaFragmentMatrix, deltaBuilder *DeltaBuilder, deltaQueue *DeltaQueue) Broadcaster {
	return Broadcaster{
		debug:   debug,
		running: 1,
		logger:  logger,

		peerQueues: peerQueues,

		deltaFragmentMatrix: deltaFragmentMatrix,
		deltaBuilder:        deltaBuilder,
		deltaQueue:          deltaQueue,
	}
}

func (broadcaster *Broadcaster) Run() {
	if broadcaster.debug {
		log.Printf("%p starting", broadcaster)
		defer log.Printf("%p shutting down", broadcaster)
	}

	go func() {
		for atomic.LoadInt32(&broadcaster.running) != 0 {
			time.Sleep(2 * time.Second)
			deltaFragments := [128]DeltaFragment{}
			if broadcaster.debug {
				log.Printf("%p is waiting for delta fragments", broadcaster)
			}
			n := broadcaster.deltaFragmentMatrix.WaitForDeltaFragments(deltaFragments[:])
			if n == 0 {
				continue
			}
			if broadcaster.debug {
				log.Printf("%p got delta fragments", broadcaster)
			}

			if broadcaster.debug {
				log.Printf("%p is broadcasting delta fragments", broadcaster)
			}
			for _, queue := range broadcaster.peerQueues {
				if err := queue.Send(Message{
					DeltaFragments: deltaFragments[:n],
				}); err != nil {
					log.Fatal(err)
				}
			}

			if broadcaster.debug {
				log.Printf("%p broadcast delta fragments", broadcaster)
			}
		}
	}()

	deltas := [128]Delta{}
	for atomic.LoadInt32(&broadcaster.running) != 0 {
		time.Sleep(2 * time.Second)
		if broadcaster.debug {
			log.Printf("%p is waiting for deltas", broadcaster)
		}
		n := broadcaster.deltaBuilder.WaitForDeltas(deltas[:])
		if n == 0 {
			continue
		}
		if broadcaster.debug {
			log.Printf("%p got deltas", broadcaster)
		}

		if broadcaster.debug {
			log.Printf("%p is broadcasting deltas", broadcaster)
		}
		for i := 0; i < n; i++ {
			if broadcaster.debug {
				log.Printf("%p broadcast delta %d", broadcaster, i)
			}
			if err := broadcaster.deltaQueue.Send(deltas[i]); err != nil {
				log.Fatal(err)
			}
		}
		if broadcaster.debug {
			log.Printf("%p broadcast deltas", broadcaster)
		}
	}
}

func (broadcaster *Broadcaster) Shutdown() {
	atomic.StoreInt32(&broadcaster.running, 0)
}
