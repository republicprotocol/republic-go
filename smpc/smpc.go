package smpc

import (
	"errors"
	"fmt"
	"sync"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/swarm"
)

// ErrJoinOnDisconnectedNetwork is returned when an Smpcer attempts to access a
// Joiner for a NetworkID that has not been connected to.
var ErrJoinOnDisconnectedNetwork = errors.New("join on disconnected network")

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
	network Network

	joinersMu *sync.RWMutex
	joiners   map[NetworkID]*Joiner

	selfJoinsMu *sync.RWMutex
	selfJoins   map[JoinID]Join
}

// NewSmpcer returns an Smpcer node that is not connected to a network.
func NewSmpcer(conn ConnectorListener, swarmer swarm.Swarmer) Smpcer {
	smpc := &smpcer{
		joinersMu: new(sync.RWMutex),
		joiners:   map[NetworkID]*Joiner{},

		selfJoinsMu: new(sync.RWMutex),
		selfJoins:   map[JoinID]Join{},
	}
	smpc.network = NewNetwork(conn, smpc, swarmer)
	return smpc
}

// Connect implements the Smpcer interface.
func (smpc *smpcer) Connect(networkID NetworkID, addrs identity.Addresses) {
	k := int64(2 * (len(addrs) + 1) / 3)

	smpc.joinersMu.Lock()
	smpc.joiners[networkID] = NewJoiner(k)
	smpc.joinersMu.Unlock()

	smpc.network.Connect(networkID, addrs)
}

// Disconnect implements the Smpcer interface.
func (smpc *smpcer) Disconnect(networkID NetworkID) {
	smpc.network.Disconnect(networkID)

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
	smpc.network.Send(networkID, message)

	return nil
}

// Receive implements the Receiver interface.
func (smpc *smpcer) Receive(from identity.Address, message Message) {
	switch message.MessageType {
	case MessageTypeJoin:
		if err := smpc.handleMessageJoin(from, message.MessageJoin); err != nil {
			logger.Network(logger.LevelError, fmt.Sprintf("error handling join message from smpc node %v: %v", from, err))
		}
	case MessageTypeJoinResponse:
		if err := smpc.handleMessageJoinResponse(message.MessageJoinResponse); err != nil {
			logger.Network(logger.LevelError, fmt.Sprintf("error handling join message from smpc node %v: %v", from, err))
		}
	default:
		logger.Network(logger.LevelError, fmt.Sprintf("error receiving message from smpc node %v: %v", from, ErrUnexpectedMessageType))
	}
}

func (smpc *smpcer) handleMessageJoin(from identity.Address, message *MessageJoin) error {
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
		smpc.network.SendTo(message.NetworkID, from, response)
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
