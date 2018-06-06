package smpc

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
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
	Connect(networkID NetworkID, nodes identity.Addresses, k int64)
	Disconnect(networkID NetworkID)
	JoinComponents(networkID NetworkID, components Components, observer ComponentBuilderObserver)
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

	shareBuildersMu       *sync.RWMutex
	shareBuilders         map[[32]byte]*ShareBuilder
	shareBuildersJoinable map[[32]byte]Component

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

		shareBuildersMu:       new(sync.RWMutex),
		shareBuilders:         map[[32]byte]*ShareBuilder{},
		shareBuildersJoinable: map[[32]byte]Component{},

		ctxCancelsMu: new(sync.Mutex),
		ctxCancels:   map[[32]byte]map[identity.Address]context.CancelFunc{},
	}
}

func (smpc *smpcer) Connect(networkID NetworkID, nodes identity.Addresses, k int64) {
	smpc.instConnect(networkID, nodes, k)
}

func (smpc *smpcer) Disconnect(networkID NetworkID) {
	smpc.instDisconnect(networkID)
}

func (smpc *smpcer) JoinComponents(networkID NetworkID, components Components, observer ComponentBuilderObserver) {
	message := Message{
		MessageType: MessageTypeJoinComponents,
		MessageJoinComponents: &MessageJoinComponents{
			NetworkID:  networkID,
			Components: components,
		},
	}

	func() {
		smpc.shareBuildersMu.RLock()
		defer smpc.shareBuildersMu.RUnlock()
		if shareBuilder, ok := smpc.shareBuilders[networkID]; ok {
			for _, component := range components {
				shareBuilder.Observe(component.ComponentID, networkID, observer)
			}
		}
	}()
	smpc.processMessageJoinComponents(nil, *message.MessageJoinComponents)

	smpc.networkMu.RLock()
	defer smpc.networkMu.RUnlock()
	for _, addr := range smpc.network[networkID] {
		log.Println("sending message to", addr)
		go smpc.sendMessage(addr, &message)
	}
}

func (smpc *smpcer) instConnect(networkID NetworkID, nodes identity.Addresses, k int64) {
	smpc.networkMu.Lock()
	smpc.shareBuildersMu.Lock()
	defer smpc.networkMu.Unlock()
	defer smpc.shareBuildersMu.Unlock()

	smpc.network[networkID] = nodes
	smpc.shareBuilders[networkID] = NewShareBuilder(k)

	go dispatch.CoForAll(nodes, func(i int) {

		addr := nodes[i]
		multiAddr, err := smpc.query(addr)
		if err != nil {
			logger.Network(logger.LevelError, fmt.Sprintf("%v", err))
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

func (smpc *smpcer) instDisconnect(networkID NetworkID) {
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

func (smpc *smpcer) sendMessage(addr identity.Address, msg *Message) {
	smpc.lookupMu.RLock()
	defer smpc.lookupMu.RUnlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if multiAddr, ok := smpc.lookup[addr]; ok {
		stream, err := smpc.streamer.Open(ctx, multiAddr)
		if err != nil {
			logger.Network(logger.LevelWarn, fmt.Sprintf("cannot open stream for messageTypeJ to %v: %v", addr, err))
			return
		}
		if err := stream.Send(msg); err != nil {
			logger.Network(logger.LevelWarn, fmt.Sprintf("cannot send messageTypeJ to %v: %v", addr, err))
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
			logger.Network(logger.LevelDebug, fmt.Sprintf("closing stream with %v: %v", remoteAddr, err))
			return
		}

		switch msg.MessageType {
		case MessageTypeJoinComponents:
			smpc.processMessageJoinComponents(&remoteAddr, *msg.MessageJoinComponents)
		case MessageTypeJoinComponentsResponse:
			smpc.processMessageJoinComponentsResponse(*msg.MessageJoinComponentsResponse)
		default:
			logger.Network(logger.LevelDebug, fmt.Sprintf("cannot recv message from %v: %v", remoteAddr, ErrUnexpectedMessageType))
		}
	}
}

func (smpc *smpcer) processMessageJoinComponents(remoteAddr *identity.Address, message MessageJoinComponents) {
	smpc.shareBuildersMu.Lock()
	defer smpc.shareBuildersMu.Unlock()

	shareBuilder, shareBuilderOk := smpc.shareBuilders[message.NetworkID]
	for _, component := range message.Components {
		if remoteAddr == nil {
			// we sent this message to ourselves so we store the component for
			// forwarding to senders
			smpc.shareBuildersJoinable[component.ComponentID] = component
		}
		if shareBuilderOk {
			if err := shareBuilder.Insert(component.ComponentID, component.Share); err != nil {
				if err != ErrInsufficientSharesToJoin {
					log.Printf("could not insert share: %v", err)
					return
				}
			}
		}
	}

	if remoteAddr != nil {
		response := Message{
			MessageType: MessageTypeJoinComponentsResponse,
			MessageJoinComponentsResponse: &MessageJoinComponentsResponse{
				NetworkID:  message.NetworkID,
				Components: make(Components, 0, len(message.Components)/2),
			},
		}

		for _, component := range message.Components {
			if responseComponent, ok := smpc.shareBuildersJoinable[component.ComponentID]; ok {
				response.MessageJoinComponentsResponse.Components = append(response.MessageJoinComponentsResponse.Components, responseComponent)
			}
		}

		go smpc.sendMessage(*remoteAddr, &response)
	}
}

func (smpc *smpcer) processMessageJoinComponentsResponse(message MessageJoinComponentsResponse) {
	smpc.shareBuildersMu.Lock()
	defer smpc.shareBuildersMu.Unlock()

	shareBuilder, shareBuilderOk := smpc.shareBuilders[message.NetworkID]
	for _, component := range message.Components {
		if shareBuilderOk {
			if err := shareBuilder.Insert(component.ComponentID, component.Share); err != nil {
				if err != ErrInsufficientSharesToJoin {
					log.Printf("could not insert share: %v", err)
					return
				}
			}
		}
	}
}
