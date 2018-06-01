package smpc

import (
	"errors"
	"fmt"
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

// ErrInsufficientSharesToJoin is returned when a share is inserted into a
// shareBuilder and there are fewer than k total shares.
var ErrInsufficientSharesToJoin = errors.New("insufficient shares to join")

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

type smpcer struct {
	swarmer  swarm.Swarmer
	streamer stream.Streamer

	buffer       int
	instructions chan Inst
	results      chan Result

	shutdownMu        *sync.Mutex
	shutdown          chan struct{}
	shutdownDone      chan struct{}
	shutdownInitiated bool

	networkMu *sync.RWMutex
	network   map[[32]byte]identity.Addresses

	lookupMu *sync.RWMutex
	lookup   map[identity.Address]identity.MultiAddress

	shareBuildersMu *sync.RWMutex
	shareBuilders   map[[32]byte]*ShareBuilder

	ctxCancelsMu *sync.Mutex
	ctxCancels   map[[32]byte]map[identity.Address]context.CancelFunc
}

func NewSmpcer(swarmer swarm.Swarmer, streamer stream.Streamer, buffer int) Smpcer {
	return &smpcer{
		swarmer:  swarmer,
		streamer: streamer,

		buffer:       buffer,
		instructions: make(chan Inst, buffer),
		results:      make(chan Result, buffer),

		shutdownMu:        new(sync.Mutex),
		shutdown:          nil,
		shutdownDone:      nil,
		shutdownInitiated: true,

		networkMu: new(sync.RWMutex),
		network:   map[[32]byte]identity.Addresses{},

		lookupMu: new(sync.RWMutex),
		lookup:   map[identity.Address]identity.MultiAddress{},

		shareBuildersMu: new(sync.RWMutex),
		shareBuilders:   map[[32]byte]*ShareBuilder{},

		ctxCancelsMu: new(sync.Mutex),
		ctxCancels:   map[[32]byte]map[identity.Address]context.CancelFunc{},
	}
}

// Start implements the Smpcer interface.
func (smpc *smpcer) Start() error {
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
func (smpc *smpcer) Shutdown() error {
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
func (smpc *smpcer) Instructions() chan<- Inst {
	return smpc.instructions
}

// Results implements the Smpcer interface.
func (smpc *smpcer) Results() <-chan Result {
	return smpc.results
}

// OnNotifyBuild implements the ShareBuilderObserver interface. It is used to
// notify the Smpcer that a value of interest has been reconstructed by the
// ShareBuilder.
func (smpc *smpcer) OnNotifyBuild(id, networkID [32]byte, value uint64) {
	log.Println("notified successfully!")
	result := Result{
		InstID:    id,
		NetworkID: networkID,
		ResultJ: &ResultJ{
			Value: value,
		},
	}
	select {
	case <-smpc.shutdown:
	case smpc.results <- result:
	}
}

func (smpc *smpcer) run() {
	defer close(smpc.shutdownDone)

	for {
		select {
		case <-smpc.shutdown:
			// FIXME: Close all open streams so that resources do not leak.
			return

		case inst, ok := <-smpc.instructions:
			if !ok {
				return
			}
			if inst.InstConnect != nil {
				smpc.instConnect(inst.NetworkID, *inst.InstConnect)
			}
			if inst.InstDisconnect != nil {
				smpc.instDisconnect(inst.NetworkID, *inst.InstDisconnect)
			}
			if inst.InstJ != nil {
				smpc.instJ(inst.InstID, inst.NetworkID, *inst.InstJ)
			}
		}
	}
}

func (smpc *smpcer) instConnect(networkID [32]byte, inst InstConnect) {
	smpc.networkMu.Lock()
	smpc.shareBuildersMu.Lock()
	defer smpc.networkMu.Unlock()
	defer smpc.shareBuildersMu.Unlock()

	smpc.network[networkID] = inst.Nodes
	smpc.shareBuilders[networkID] = NewShareBuilder(inst.K)

	go dispatch.CoForAll(inst.Nodes, func(i int) {

		addr := inst.Nodes[i]
		multiAddr, err := smpc.query(addr)
		if err != nil {
			log.Println(err)
			return
		}

		smpc.lookupMu.Lock()
		smpc.lookup[addr] = multiAddr
		smpc.lookupMu.Unlock()

		ctx, cancel := context.WithCancel(context.Background())
		stream, err := smpc.streamer.Open(ctx, multiAddr)
		if err != nil {
			log.Println(fmt.Errorf("cannot connect to smpcer node %v: %v", addr, err))
			return
		}
		go smpc.stream(addr, stream)

		smpc.ctxCancelsMu.Lock()
		if _, ok := smpc.ctxCancels[networkID]; !ok {
			smpc.ctxCancels[networkID] = map[identity.Address]context.CancelFunc{}
		}
		smpc.ctxCancels[networkID][addr] = cancel
		smpc.ctxCancelsMu.Unlock()
	})
}

func (smpc *smpcer) instDisconnect(networkID [32]byte, inst InstDisconnect) {
	smpc.networkMu.Lock()
	smpc.shareBuildersMu.Lock()
	smpc.ctxCancelsMu.Lock()
	defer smpc.networkMu.Unlock()
	defer smpc.shareBuildersMu.Unlock()
	defer smpc.ctxCancelsMu.Unlock()

	if _, ok := smpc.ctxCancels[networkID]; ok {
		for _, addr := range smpc.network[networkID] {
			smpc.ctxCancels[networkID][addr]()
		}
	}

	delete(smpc.network, networkID)
	delete(smpc.shareBuilders, networkID)
	delete(smpc.ctxCancels, networkID)
}

func (smpc *smpcer) instJ(instID, networkID [32]byte, inst InstJ) {
	msg := Message{
		MessageType: MessageTypeJ,
		MessageJ: &MessageJ{
			InstID:    instID,
			NetworkID: networkID,
			Share:     inst.Share,
		},
	}

	smpc.networkMu.RLock()
	smpc.lookupMu.RLock()
	smpc.shareBuildersMu.RLock()
	defer smpc.networkMu.RUnlock()
	defer smpc.lookupMu.RUnlock()
	defer smpc.shareBuildersMu.RUnlock()

	if shareBuilder, ok := smpc.shareBuilders[networkID]; ok {
		shareBuilder.Observe(instID, networkID, smpc)
	}
	smpc.processMessageJ(*msg.MessageJ)

	for _, addr := range smpc.network[networkID] {
		go smpc.sendMessage(addr, &msg)
	}
}

func (smpc *smpcer) sendMessage(addr identity.Address, msg *Message) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if multiAddr, ok := smpc.lookup[addr]; ok {
		stream, err := smpc.streamer.Open(ctx, multiAddr)
		if err != nil {
			log.Printf("cannot open stream for messageTypeJ to %v: %v", addr, err)
			return
		}
		if err := stream.Send(msg); err != nil {
			log.Printf("cannot send messageTypeJ to %v: %v", addr, err)
		}
	}
}

func (smpc *smpcer) query(addr identity.Address) (identity.MultiAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	multiAddr, err := smpc.swarmer.Query(ctx, addr, -1)
	if err != nil {
		return multiAddr, fmt.Errorf("cannot query smpcer node %v: %v", addr, err)
	}
	return multiAddr, nil
}

func (smpc *smpcer) stream(remoteAddr identity.Address, remoteStream stream.Stream) {
	for {
		msg := Message{}
		if err := remoteStream.Recv(&msg); err != nil {
			log.Printf("closing stream with %v: %v", remoteAddr, err)
			return
		}

		switch msg.MessageType {
		case MessageTypeJ:
			smpc.processMessageJ(*msg.MessageJ)
		default:
			log.Printf("cannot recv message from %v: %v", remoteAddr, ErrUnexpectedMessageType)
		}
	}
}

func (smpc *smpcer) processMessageJ(message MessageJ) {
	smpc.shareBuildersMu.RLock()
	defer smpc.shareBuildersMu.RUnlock()

	if shareBuilder, ok := smpc.shareBuilders[message.NetworkID]; ok {
		log.Println("inserting value into share builder")
		if err := shareBuilder.Insert(message.InstID, message.Share); err != nil {
			if err == ErrInsufficientSharesToJoin {
				return
			}
			log.Printf("could not insert share: %v", err)
			return
		}
	}
}
