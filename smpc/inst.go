package smpc

import (
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/shamir"
)

// Inst stores instructions that will be consumed by the Smpcer and executed
// asynchronously. All instructions involve a network of nodes that the Smpcer
// will communicate with, identified by the network ID.
type Inst struct {
	InstID    [32]byte
	NetworkID [32]byte

	*InstConnect
	*InstDisconnect
	*InstJ
}

// InstConnect instructs the Smpcer to connect to another network of nodes,
// identified by their identity.Address. It is the responsibility of the Smpcer
// to know how to use these identity.Addresses to connect with other Smpcer
// nodes.
type InstConnect struct {
	Nodes identity.Addresses
	N     int64
	K     int64
}

// InstDisconnect instructs the Smpcer to disconnect from a network.
type InstDisconnect struct {
}

// InstJ instructs the Smpcer to join shares into a value. The Smpcer will
// will eventually output a respective ResultValue.
type InstJ struct {
	shamir.Share
}

// Result stores the results of an Inst after it has been executed to
// completion by the Smpcer. An ID associates it with the respective Inst and a
// NetworkID associates it with the network of Smpcer nodes that executed the
// Inst.
type Result struct {
	InstID    [32]byte
	NetworkID [32]byte

	*ResultJ
}

// ResultJ is a result from the execution of an InstJ instruction.
type ResultJ struct {
	Value uint64
	Err   error
}
