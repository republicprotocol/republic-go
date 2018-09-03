package grpc

import (
	"log"
	"math/rand"
	"time"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/smpc"
	"github.com/republicprotocol/republic-go/testutils"
)

// RedNodeBehaviour indicates the malicious behaviours the
// red node will exhibit.
type RedNodeBehaviour int

// Values for a RedNodeBehaviour
const (
	InvalidRequests RedNodeBehaviour = iota
	InvalidNonce
	InvalidSignature
	InvalidBlindings
	InvalidJoins
	DropMessages
	DropMultiAddresses
	DropSignatures
)

// String returns a human-readable representation of RedNodeTypes.
func (behaviours RedNodeBehaviour) String() string {
	switch behaviours {
	case InvalidRequests:
		return "invalid requests"
	case InvalidNonce:
		return "invalid nonce"
	case InvalidSignature:
		return "invalid multi-address signature"
	case DropMultiAddresses:
		return "drop multi-addresses"
	case DropSignatures:
		return "drop multi-address signatures"
	default:
		return "unexpected behaviour"
	}
}

// RedNodeSwarmerTypes contains an array of all possible malicious swarming
// behaviours.
var RedNodeSwarmerTypes = []RedNodeBehaviour{
	InvalidRequests,
	InvalidNonce,
	InvalidSignature,
	DropMultiAddresses,
	DropSignatures,
}

// RedNodeStreamerTypes contains an array of all possible malicious streaming
// behaviours.
var RedNodeStreamerTypes = []RedNodeBehaviour{
	InvalidRequests,
	InvalidBlindings,
	InvalidJoins,
	DropMessages,
}

func getTamperedMultiAddress(multiAddr identity.MultiAddress) MultiAddress {
	redNodeType := RedNodeSwarmerTypes[rand.Intn(len(RedNodeSwarmerTypes))]

	rand.Seed(time.Now().UnixNano())
	multiAddress := MultiAddress{
		Signature:         multiAddr.Signature,
		MultiAddress:      multiAddr.String(),
		MultiAddressNonce: multiAddr.Nonce,
	}

	switch redNodeType {
	case InvalidRequests:
		multiAddress.Signature = tamperSignature(multiAddr)
		multiAddress.MultiAddressNonce = tamperNonce(multiAddr)
		multiAddress.MultiAddress = tamperMultiAddress(multiAddr)
	case InvalidNonce:
		multiAddress.MultiAddressNonce = tamperNonce(multiAddr)
	case InvalidSignature:
		multiAddress.Signature = tamperSignature(multiAddr)
	case DropMultiAddresses:
		multiAddress.MultiAddress = ""
	case DropSignatures:
		multiAddress.Signature = []byte{}
	default:
	}

	log.Printf("Red-node swarmer will exhibit behaviour: %v\n", redNodeType)
	log.Printf("Red-node tampered multi-address %v to look like %v", multiAddr, multiAddress)
	return multiAddress
}

func tamperMessage(message smpc.Message) smpc.Message {
	// r := rand.Intn(100)
	// if r < 75 {
	// message.MessageJoin.Join
	// }
	return message
}

func tamperSignature(multiAddr identity.MultiAddress) []byte {
	r := rand.Intn(100)
	if r < 50 {
		randBytes := testutils.Random64Bytes()
		return randBytes[:]
	}
	multiAddr.Signature[rand.Intn(64)] = byte(rand.Intn(100))
	return multiAddr.Signature
}

func tamperMultiAddress(multiAddr identity.MultiAddress) string {
	r := rand.Intn(100)
	if r < 75 {
		multiAddr, _ := testutils.RandomMultiAddress()
		return multiAddr.String()
	}
	return multiAddr.String()
}

func tamperNonce(multiAddr identity.MultiAddress) uint64 {
	r := rand.Intn(100)
	if r < 33 {
		return multiAddr.Nonce + uint64(r)
	}
	if r < 66 {
		return multiAddr.Nonce - uint64(r)
	}
	if r < 99 {
		return 0
	}
	return multiAddr.Nonce
}
