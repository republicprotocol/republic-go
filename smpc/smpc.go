package smpc

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/stream"
	"github.com/republicprotocol/republic-go/swarm"
	"golang.org/x/net/context"
)

// ErrSmpcerIsAlreadyRunning is returned when a call to Smpcer.Start happens
// on an Smpcer that has already been started.
var ErrSmpcerIsAlreadyRunning = errors.New("smpcer is already running")

// ErrSmpcerIsNotRunning is returned when a call to Smpcer.Shutdown happens on
// an Smpcer that has not been started yet.
var ErrSmpcerIsNotRunning = errors.New("smpcer is not running")

// ErrUnmarshalNilBytes is returned when a call to UnmarshalBinary happens on
// an empty list of bytes.
var ErrUnmarshalNilBytes = errors.New("unmarshall nil bytes")

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
	swarmer  swarm.Swarmer
	streamer stream.Streamer

	buffer       int
	instructions chan Inst
	results      chan Result

	shutdownMu        *sync.Mutex
	shutdown          chan struct{}
	shutdownDone      chan struct{}
	shutdownInitiated bool

	routesMu *sync.RWMutex
	routes   map[string]map[identity.Address]stream.Stream
}

func NewSmpc(swarmer swarm.Swarmer, streamer stream.Streamer, buffer int) Smpcer {
	return &smpc{
		swarmer:  swarmer,
		streamer: streamer,

		buffer:       buffer,
		instructions: make(chan Inst, buffer),
		results:      make(chan Result, buffer),

		shutdownMu:        new(sync.Mutex),
		shutdown:          nil,
		shutdownDone:      nil,
		shutdownInitiated: true,

		routesMu: new(sync.RWMutex),
		routes:   map[string]map[identity.Address]stream.Stream{},
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

	go smpc.run()
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

	close(smpc.instructions)
	close(smpc.results)
	smpc.instructions = make(chan Inst, smpc.buffer)
	smpc.results = make(chan Result, smpc.buffer)

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
	defer close(smpc.shutdownDone)

	for {
		select {
		case <-smpc.shutdown:
			return

		case inst, ok := <-smpc.instructions:
			if !ok {
				return
			}
			if inst.InstConnect != nil {
				smpc.instConnect(inst.NetworkID, *inst.InstConnect)
			}
			if inst.InstDisconnect != nil {
				smpc.instDisconnect(inst.NetworkID, *inst.InstConnect)
			}
			if inst.InstJ != nil {
				smpc.instJ(inst.ID, inst.NetworkID, *inst.InstJoin)
			}
		}
	}
}

func (smpc *smpc) instConnect(networkID []byte, inst InstConnect) {
	go dispatch.CoForAll(inst.Nodes, func(i int) {
		queryCtx, queryCancel := context.WithTimeout(context.Background(), time.Minute)
		defer queryCancel()

		addr := inst.Nodes[i]
		multiAddr, err := smpc.swarmer.Query(queryCtx, addr, -1)
		if err != nil {
			log.Printf("cannot query smpcer node %v: %v", addr, err)
			return
		}

		openCtx, openCancel := context.WithCancel(context.Background())
		defer openCancel()

		// TODO: Store openCancel so that it can be called by the
		// smpc.instDisconnect method

		s, err := smpc.streamer.Open(openCtx, multiAddr)
		if err != nil {
			log.Printf("cannot connect to smpcer node %v: %v", addr, err)
		}

		smpc.routesMu.Lock()
		smpc.routes[string(networkID)][addr] = s
		smpc.routesMu.Unlock()

		for {
			msg := message{}
			if err := s.Recv(&msg); err != nil {
				if err == stream.ErrRecvOnClosedStream {
					return
				}
				log.Printf("cannot recv message from %v: %v", addr, err)
				continue
			}

		}
	})
}

func (smpc *smpc) instDisconnect(networkID []byte, inst InstDisconnect) {
	panic("unimplemented")
}

func (smpc *smpc) instJ(instID []byte, networkID []byte, inst InstJ) {
}
