package smpc

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/republicprotocol/republic-go/shamir"

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
	// The ordering of the nodes is used to determine which nodes are expected
	// to hold which fragments.
	Connect(networkID NetworkID, nodes identity.Addresses)

	// Disconnect from a network of nodes.
	Disconnect(networkID NetworkID)

	// Join a set of shamir.Shares for distinct values. This involves broadcast
	// communication with the nodes in the network. On a success, the Callback
	// is called.
	Join(networkID NetworkID, join Join, callback Callback) error

	// InsertCommitments for the shamir.Shares inside a Join. These commitments
	// are used to blind shamir.Shares while being able to verify that the
	// computations performed have been done correctly.
	InsertCommitments(networkID NetworkID, joinID JoinID, joinCommitments JoinCommitments)
}

type smpcer struct {
	network Network

	joinersMu *sync.RWMutex
	joiners   map[NetworkID]*Joiner

	selfJoinsMu *sync.RWMutex
	selfJoins   map[JoinID]Join

	commitmentsMu *sync.RWMutex
	commitments   map[NetworkID](map[JoinID]JoinCommitments)
}

// NewSmpcer returns an Smpcer node that is not connected to a network.
func NewSmpcer(conn ConnectorListener, swarmer swarm.Swarmer) Smpcer {
	smpc := &smpcer{
		joinersMu: new(sync.RWMutex),
		joiners:   map[NetworkID]*Joiner{},

		selfJoinsMu: new(sync.RWMutex),
		selfJoins:   map[JoinID]Join{},

		commitmentsMu: new(sync.RWMutex),
		commitments:   map[NetworkID](map[JoinID]JoinCommitments){},
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

	smpc.commitmentsMu.Lock()
	smpc.commitments[networkID] = map[JoinID]JoinCommitments{}
	smpc.commitmentsMu.Unlock()

	smpc.network.Connect(networkID, addrs)
}

// Disconnect implements the Smpcer interface.
func (smpc *smpcer) Disconnect(networkID NetworkID) {
	smpc.network.Disconnect(networkID)

	smpc.joinersMu.Lock()
	delete(smpc.joiners, networkID)
	smpc.joinersMu.Unlock()

	smpc.commitmentsMu.Lock()
	delete(smpc.commitments, networkID)
	smpc.commitmentsMu.Unlock()
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

func (smpc *smpcer) InsertCommitments(networkID NetworkID, joinID JoinID, joinCommitements JoinCommitments) {
	smpc.commitmentsMu.Lock()
	defer smpc.commitmentsMu.Unlock()

	if _, ok := smpc.commitments[networkID]; !ok {
		return
	}
	smpc.commitments[networkID][joinID] = joinCommitements
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
			logger.Network(logger.LevelError, fmt.Sprintf("error handling joinResponse message from smpc node %v: %v", from, err))
		}
	default:
		logger.Network(logger.LevelError, fmt.Sprintf("error receiving message from smpc node %v: %v", from, ErrUnexpectedMessageType))
	}
}

func (smpc *smpcer) handleMessageJoin(from identity.Address, message *MessageJoin) error {
	if !smpc.verifyJoin(message.NetworkID, message.Join) {
		return ErrUnverifiedJoin
	}

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
	if !smpc.verifyJoin(message.NetworkID, message.Join) {
		return ErrUnverifiedJoin
	}

	var err error
	smpc.joinersMu.RLock()
	if joiner, ok := smpc.joiners[message.NetworkID]; ok {
		err = joiner.InsertJoin(message.Join)
	}
	smpc.joinersMu.RUnlock()

	return err
}

func (smpc *smpcer) verifyJoin(networkID NetworkID, join Join) bool {
	// Always require that each share has a blinding
	if len(join.Shares) != len(join.Blindings) {
		logger.Debug(fmt.Sprintf("share and blindings have different length, Blindings= %v, shares= %v", len(join.Shares), len(join.Blindings)))
		return false
	}

	// Get the JoinCommitments stored for the network and join ID being
	// verified
	joinCommitments, ok := func() (JoinCommitments, bool) {
		smpc.commitmentsMu.RLock()
		defer smpc.commitmentsMu.RUnlock()

		if _, ok := smpc.commitments[networkID]; !ok {
			return JoinCommitments{}, false
		}
		if _, ok := smpc.commitments[networkID][join.ID]; !ok {
			return JoinCommitments{}, false
		}

		return smpc.commitments[networkID][join.ID], true
	}()
	// We accept the join if we do not have the JoinCommitments to check for
	// incorrectness
	if !ok {
		return true
	}

	for i := range join.Shares {

		// Get the relevant commitments for the LHS and RHS of the computation
		// in the join
		lhs, ok := joinCommitments.LHS[join.Shares[i].Index]
		if !ok {
			// We accept the join if we do not have the JoinCommitments to
			// check for incorrectness
			continue
		}
		rhs, ok := joinCommitments.RHS[join.Shares[i].Index]
		if !ok {
			// We accept the join if we do not have the JoinCommitments to
			// check for incorrectness
			continue
		}
		rhs.Int.ModInverse(rhs.Int, shamir.CommitP)

		// Check the expected commitment against the commitment we actually
		// received
		expected := big.NewInt(0).Mul(lhs.Int, rhs.Int)
		expected.Mod(expected, shamir.CommitP)
		got := shamir.NewCommitment(join.Shares[i], join.Blindings[i])

		if expected.Cmp(got.Int) != 0 {
			// Reject the join
			log.Printf("reject the join due to %vth share, expected=%v , got=%v", i, expected.Int64(), got.Int64())
			return false
		}
	}

	// Accept the join
	return true
}
