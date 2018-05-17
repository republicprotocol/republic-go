package smpc

import (
	"errors"
	"sync"
)

// ErrSmpcerIsAlreadyRunning is returned when a call to Smpcer.Start happens
// on an Smpcer that has already been started.
var ErrSmpcerIsAlreadyRunning = errors.New("smpcer is already running")

// ErrSmpcerIsNotRunning is returned when a call to Smpcer.Shutdown happens on
// an Smpcer that has not been started yet.
var ErrSmpcerIsNotRunning = errors.New("smpcer is not running")

// Smpcer is an interface for a secure multi-party computer. It asynchronously
// consumes computation instructions and produces computation results.
type Smpcer interface {

	// Start the Smpcer. Until a call to Smpcer.Start, no computation
	// instruction will be processed.
	Start() error

	// Shutdown the Smpcer. After a call to Smpcer.Shutdown, no computation
	// instruction will be processed.
	Shutdown() error

	// Instructions channel for sending computation instructions to the Smpcer.
	Instructions() chan<- Inst

	// Results channel for receiving computation results from the Smpcer.
	Results() <-chan Result
}

type smpc struct {
	instructions chan Inst
	results      chan Result

	shutdownMu        *sync.Mutex
	shutdown          chan struct{}
	shutdownDone      chan struct{}
	shutdownInitiated bool
}

func NewSmpc(buffer int) Smpcer {
	return &smpc{
		instructions: make(chan Inst, buffer),
		results:      make(chan Result, buffer),

		shutdownMu:        new(sync.Mutex),
		shutdown:          nil,
		shutdownDone:      nil,
		shutdownInitiated: true,
	}
}

// Start implements the Smpcer interface.
func (smpc *smpc) Start() error {
	smpc.shutdownMu.Lock()
	defer smpc.shutdownMu.Unlock()

	if smpc.shutdown != nil {
		return ErrSmpcerIsAlreadyRunning
	}
	smpc.shutdown = make(chan struct{})
	smpc.shutdownDone = make(chan struct{})
	smpc.shutdownInitiated = false

	go smpcer.run()
	return nil
}

// Shutdown implements the Smpcer interface.
func (smpc *smpc) Shutdown() error {
	smpc.shutdownMu.Lock()
	defer smpc.shutdownMu.Unlock()

	if smpc.shutdownInitiated {
		return ErrSmpcerIsNotRunning
	}
	smpc.shutdownInitiated = true

	close(smpc.shutdown)
	<-smpc.shutdownDone

	smpc.shutdown = nil
	smpc.shutdownDone = nil

	return nil
}

// Instructions implements the Smpcer interface.
func (smpc *smpc) Instructions() chan<- Inst {
	return smpc.instructions
}

// Results implements the Smpcer interface.
func (smpc *smpc) Results() <-chan Result {
	return smpc.results
}

func (smpc *smpc) run() {
	for {
		select {
		case <-smpc.shutdown:
			close(smpc.shutdownDone)
			return
		}
	}
}
