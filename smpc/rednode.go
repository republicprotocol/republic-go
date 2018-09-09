package smpc

import (
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"time"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/shamir"
)

// RedNodeBehaviour indicates the malicious behaviours the
// red node will exhibit.
type RedNodeBehaviour int

// Values for a RedNodeBehaviour
const (
	InvalidMessageRequests RedNodeBehaviour = iota
	DropMessages
)

// String returns a human-readable representation of RedNodeTypes.
func (behaviours RedNodeBehaviour) String() string {
	switch behaviours {
	case InvalidMessageRequests:
		return "invalid smpc message requests"
	case DropMessages:
		return "drop smpc messages"
	default:
		return "unexpected behaviour"
	}
}

// RedNodeStreamerTypes contains an array of all possible malicious streaming
// behaviours.
var RedNodeStreamerTypes = []RedNodeBehaviour{
	InvalidMessageRequests,
	DropMessages,
}

func getTamperedMessage(message Message) Message {
	rand.Seed(time.Now().UnixNano())

	redNodeType := RedNodeStreamerTypes[rand.Intn(len(RedNodeStreamerTypes))]

	log.Printf("Red-node streamer will exhibit behaviour: %v; with original message: %v\n", redNodeType, message)
	switch redNodeType {
	case InvalidMessageRequests:
		message = tamperMessage(message)
	case DropMessages:
		message = Message{}
	default:
	}
	log.Printf("Red-node with behaviour %v tampered the message to look like %v", redNodeType, message)

	return message
}

func tamperMessage(message Message) Message {
	r := rand.Intn(100)

	switch message.MessageType {
	case MessageTypeJoin:
		if r < 50 {
			message.MessageType = MessageTypeJoinResponse
			message.MessageJoinResponse.NetworkID = tamperNetworkID(message.MessageJoin.NetworkID)
			message.MessageJoinResponse.Join = tamperMessageJoin(message.MessageJoin.Join)
			return message
		}
		message.MessageJoin.NetworkID = tamperNetworkID(message.MessageJoin.NetworkID)
		message.MessageJoin.Join = tamperMessageJoin(message.MessageJoin.Join)
	case MessageTypeJoinResponse:
		if r < 50 {
			message.MessageType = MessageTypeJoin
			message.MessageJoin.NetworkID = tamperNetworkID(message.MessageJoinResponse.NetworkID)
			message.MessageJoin.Join = tamperMessageJoin(message.MessageJoinResponse.Join)
			return message
		}
		message.MessageJoinResponse.NetworkID = tamperNetworkID(message.MessageJoinResponse.NetworkID)
		message.MessageJoinResponse.Join = tamperMessageJoin(message.MessageJoinResponse.Join)
	default:
		message.MessageType = MessageType(15)
	}
	if r < 80 && r >= 50 {
		message.MessageType = MessageType(0)
	}
	return message
}

func tamperMessageJoin(join Join) Join {
	r := rand.Intn(100)
	// Return an empty Join.
	if r < 10 {
		return Join{}
	}
	// Return an updated Join.
	if r < 90 {
		join.ID = tamperJoinID(join.ID)
		join.Index = tamperJoinIndex(join.Index)
		join.Shares = tamperShares(join.Shares)
		join.Blindings = tamperBlindings(join.Blindings)
	}
	return join
}

func tamperJoinID(joinID JoinID) JoinID {
	r := rand.Intn(100)
	// Return an empty spmc.JoinID.
	if r < 10 {
		return JoinID{}
	}
	// Return a random [33]byte array as JoinID.
	if r < 50 {
		return JoinID(Random33Bytes())
	}
	// Modify the joinID slightly.
	if r < 90 {
		index := rand.Intn(33)
		joinID[index] = byte(index)
		return joinID
	}
	return joinID
}

func tamperJoinIndex(joinIndex JoinIndex) JoinIndex {
	r := rand.Intn(100)
	// Return an 0.
	if r < 10 {
		return 0
	}
	// Return a random uint64 as JoinIndex.
	if r < 50 {
		return JoinIndex(rand.Intn(200))
	}
	// Modify the joinIndex slightly.
	if r < 70 {
		return joinIndex + 1
	}
	if r < 90 {
		return joinIndex - 1
	}
	return joinIndex
}

func tamperShares(shares shamir.Shares) shamir.Shares {
	r := rand.Intn(100)
	// Return an empty list of shamir.Shares.
	if r < 10 {
		return shamir.Shares{}
	}
	// Return a random set of shamir.Shares.
	if r < 50 {
		secret := ((uint64(rand.Int63()) % shamir.Prime) / 2) + (shamir.Prime / 2)
		shares, _ = shamir.Split(72, 48, secret)
		return shares
	}
	// Modify the shares slightly.
	if r < 90 {
		index := rand.Intn(len(shares))
		shares[index] = shamir.Share{Index: uint64(index), Value: uint64(index)}
	}
	return shares
}

func tamperBlindings(blindings shamir.Blindings) shamir.Blindings {
	r := rand.Intn(100)
	// Return an empty list of shamir.Blindings.
	if r < 10 {
		return shamir.Blindings{}
	}
	// Return a completely different random set of shamir.Blindings.
	if r < 50 {
		for i := range blindings {
			blindings[i] = shamir.Blinding{Int: big.NewInt(int64(rand.Intn(100)))}
		}
		return blindings
	}
	// Modify the blindings slightly.
	if r < 90 {
		index := rand.Intn(len(blindings))
		blindings[index] = shamir.Blinding{Int: big.NewInt(int64(index))}
	}
	return blindings
}

func tamperNetworkID(networkID NetworkID) NetworkID {
	r := rand.Intn(100)
	// Return a randomly generated [32]byte array.
	if r < 50 {
		return Random32Bytes()
	}
	// Return a slightly modified networkID.
	if r < 70 {
		index := rand.Intn(32)
		networkID[index] = byte(index)
		return networkID
	}
	// Return a nil array.
	if r < 90 {
		return [32]byte{}
	}
	return networkID
}

// Random32Bytes creates a random [32]byte.
func Random32Bytes() [32]byte {
	var res [32]byte
	i := fmt.Sprintf("%d", rand.Int())
	hash := crypto.Keccak256([]byte(i))
	copy(res[:], hash)
	return res
}

// Random33Bytes creates a random [33]byte.
func Random33Bytes() [33]byte {
	var res [33]byte
	i := fmt.Sprintf("%d", rand.Int())
	hash := crypto.Keccak256([]byte(i))
	copy(res[:], hash)
	return res
}

// Random64Bytes creates a random [64]]byte.
func Random64Bytes() [64]byte {
	var res [64]byte
	i := fmt.Sprintf("%d", rand.Int())
	hash := crypto.Keccak256([]byte(i))
	copy(res[:], hash)
	return res
}
