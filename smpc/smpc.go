package smpc

import (
	"encoding/base64"
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

// ErrJoinOnDisconnectedNetwork is returned when an Smpcer attempts to access a
// Joiner for a NetworkID that has not been connected to.
var ErrJoinOnDisconnectedNetwork = errors.New("join on disconnected network")

// NetworkID for a network of Smpcer nodes. Using a NetworkID allows nodes to
// be involved in multiple distinct computation networks in parallel.
type NetworkID [32]byte

// String returns a human-readable representation of a NetworkID.
func (id NetworkID) String() string {
	return base64.StdEncoding.EncodeToString(id[:8])
}

// Smpcer is an interface for a secure multi-party computer. It asynchronously
// consumes computation instructions and produces computation results.
type Smpcer interface {

	// Connect to a network of nodes and assign this network to a NetworkID.
	Connect(networkID NetworkID, nodes identity.Addresses)

	// Disconnect from a network of nodes.
	Disconnect(networkID NetworkID)

	// Join a set of shamir.Shares for distinct values. This involves broadcast
	// communication with the nodes in the network. On a success, the Callback
	// is called.
	Join(networkID NetworkID, join Join, callback Callback) error
}

type smpcer struct {
	swarmer  swarm.Swarmer
	streamer stream.Streamer

	networkMu      *sync.RWMutex
	network        map[NetworkID]identity.Addresses
	networkCancels map[NetworkID][]context.CancelFunc

	lookupMu *sync.RWMutex
	lookup   map[identity.Address]identity.MultiAddress

	joinersMu *sync.RWMutex
	joiners   map[NetworkID]*Joiner

	selfJoinsMu *sync.RWMutex
	selfJoins   map[JoinID]Join
}

// NewSmpcer returns an Smpcer node that is not connected to a network.
func NewSmpcer(swarmer swarm.Swarmer, streamer stream.Streamer) Smpcer {
	return &smpcer{
		swarmer:  swarmer,
		streamer: streamer,

		networkMu:      new(sync.RWMutex),
		network:        map[NetworkID]identity.Addresses{},
		networkCancels: map[NetworkID][]context.CancelFunc{},

		lookupMu: new(sync.RWMutex),
		lookup:   map[identity.Address]identity.MultiAddress{},

		joinersMu: new(sync.RWMutex),
		joiners:   map[NetworkID]*Joiner{},

		selfJoinsMu: new(sync.RWMutex),
		selfJoins:   map[JoinID]Join{},
	}
}

// Connect implements the Smpcer interface.
func (smpc *smpcer) Connect(networkID NetworkID, nodes identity.Addresses) {
	k := int64(2 * (len(nodes) + 1) / 3)

	smpc.networkMu.Lock()
	smpc.network[networkID] = nodes
	smpc.networkMu.Unlock()

	smpc.joinersMu.Lock()
	smpc.joiners[networkID] = NewJoiner(k)
	smpc.joinersMu.Unlock()

	logger.Network(logger.LevelInfo, fmt.Sprintf("connecting to network = %v, thresold = (%v, %v)", networkID, len(nodes), k))

	go dispatch.CoForAll(nodes, func(i int) {
		addr := nodes[i]
		if addr == smpc.swarmer.MultiAddress().Address() {
			// Skip trying to connect to ourself
			return
		}
		multiAddr, err := smpc.query(addr)
		if err != nil {
			logger.Network(logger.LevelError, fmt.Sprintf("cannot connect to smpc node %v: %v", addr, err))
			return
		}

		// Store the identity.Identity to identity.MultiAddress mapping
		smpc.lookupMu.Lock()
		smpc.lookup[addr] = multiAddr
		smpc.lookupMu.Unlock()

		logger.Network(logger.LevelInfo, fmt.Sprintf("connecting to %v", addr))

		// Open a stream to the node and store the context.CancelFunc so that
		// we can call it when we need to disconnect
		ctx, cancel := context.WithCancel(context.Background())
		stream, err := smpc.streamer.Open(ctx, multiAddr)
		if err != nil {
			log.Println(fmt.Errorf("cannot open stream to smpc node %v: %v", addr, err))
			return
		}

		smpc.networkMu.Lock()
		if _, ok := smpc.networkCancels[networkID]; !ok {
			smpc.networkCancels[networkID] = make([]context.CancelFunc, 0, len(nodes))
		}
		smpc.networkCancels[networkID] = append(smpc.networkCancels[networkID], cancel)
		smpc.networkMu.Unlock()

		// A background goroutine will handle the stream
		logger.Network(logger.LevelDebug, fmt.Sprintf("connected to %v in network %v", addr, networkID))
		go smpc.handleStream(ctx, addr, stream)
	})
}

// Disconnect implements the Smpcer interface.
func (smpc *smpcer) Disconnect(networkID NetworkID) {
	logger.Network(logger.LevelInfo, fmt.Sprintf("disconnecting from network = %v", networkID))

	smpc.networkMu.Lock()
	if _, ok := smpc.networkCancels[networkID]; ok {
		for _, cancel := range smpc.networkCancels[networkID] {
			cancel()
		}
	}
	delete(smpc.network, networkID)
	delete(smpc.networkCancels, networkID)
	smpc.networkMu.Unlock()

	smpc.joinersMu.Lock()
	delete(smpc.joiners, networkID)
	smpc.joinersMu.Unlock()
}

// Join implements the Smpcer interface.
func (smpc *smpcer) Join(networkID NetworkID, join Join, callback Callback) error {
	smpc.selfJoinsMu.Lock()
	smpc.selfJoins[join.ID] = join
	smpc.selfJoinsMu.Unlock()

	smpc.joinersMu.RLock()
	joiner, joinerOk := smpc.joiners[networkID]
	smpc.joinersMu.RUnlock()
	if !joinerOk {
		return ErrJoinOnDisconnectedNetwork
	}
	if err := joiner.InsertJoinAndSetCallback(join, callback); err != nil {
		return err
	}

	message := Message{
		MessageType: MessageTypeJoin,
		MessageJoin: &MessageJoin{
			NetworkID: networkID,
			Join:      join,
		},
	}

	func() {
		smpc.networkMu.RLock()
		defer smpc.networkMu.RUnlock()

		for _, addr := range smpc.network[networkID] {
			go smpc.sendMessage(addr, &message)
		}
	}()

	return nil
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

func (smpc *smpcer) handleStream(ctx context.Context, remoteAddr identity.Address, remoteStream stream.Stream) {
	timeout := 1000
	for {
		// Check whether or not the stream has closed
		select {
		case <-ctx.Done():
			return
		default:
		}

		message := Message{}
		if err := remoteStream.Recv(&message); err != nil {
			log.Printf("[error] (network) cannot receive message from %v: %v", remoteAddr, err)
			time.Sleep(time.Duration(timeout) * time.Millisecond)
			timeout = int(float64(timeout) * 1.6)
			continue
		}
		timeout = 1000

		switch message.MessageType {
		case MessageTypeJoin:
			if err := smpc.handleMessageJoin(remoteAddr, message.MessageJoin); err != nil {
				logger.Network(logger.LevelError, fmt.Sprintf("error handling join message from smpc node %v: %v", remoteAddr, err))
			}
		case MessageTypeJoinResponse:
			if err := smpc.handleMessageJoinResponse(message.MessageJoinResponse); err != nil {
				logger.Network(logger.LevelError, fmt.Sprintf("error handling join message from smpc node %v: %v", remoteAddr, err))
			}
		default:
			logger.Network(logger.LevelError, fmt.Sprintf("error receiving message from smpc node %v: %v", remoteAddr, ErrUnexpectedMessageType))
		}
	}
}

func (smpc *smpcer) handleMessageJoin(remoteAddr identity.Address, message *MessageJoin) error {
	var err error
	smpc.joinersMu.RLock()
	if joiner, ok := smpc.joiners[message.NetworkID]; ok {
		err = joiner.InsertJoin(message.Join)
	}
	smpc.joinersMu.RUnlock()
	if err != nil {
		return err
	}

	go func() {
		smpc.selfJoinsMu.RLock()
		join, ok := smpc.selfJoins[message.Join.ID]
		smpc.selfJoinsMu.RUnlock()
		if !ok {
			return
		}

		response := Message{
			MessageType: MessageTypeJoinResponse,
			MessageJoinResponse: &MessageJoinResponse{
				NetworkID: message.NetworkID,
				Join:      join,
			},
		}
		smpc.sendMessage(remoteAddr, &response)
	}()

	return nil
}

func (smpc *smpcer) handleMessageJoinResponse(message *MessageJoinResponse) error {
	var err error
	smpc.joinersMu.RLock()
	if joiner, ok := smpc.joiners[message.NetworkID]; ok {
		err = joiner.InsertJoin(message.Join)
	}
	smpc.joinersMu.RUnlock()
	if err != nil {
		return err
	}
	return nil
}

func (smpc *smpcer) sendMessage(addr identity.Address, msg *Message) {
	smpc.lookupMu.RLock()
	defer smpc.lookupMu.RUnlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if multiAddr, ok := smpc.lookup[addr]; ok {
		stream, err := smpc.streamer.Open(ctx, multiAddr)
		if err != nil {
			logger.Network(logger.LevelWarn, fmt.Sprintf("cannot open messaging stream to smpc node %v: %v", addr, err))
			return
		}
		if err := stream.Send(msg); err != nil {
			logger.Network(logger.LevelWarn, fmt.Sprintf("cannot send message to smpc node %v: %v", addr, err))
		}
	}
}
