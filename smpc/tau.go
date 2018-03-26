package smpc

import (
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/network/rpc"
)

// Tau implements the Tau gRPC service. Tau creates MessageQueues for each gRPC
// stream, and runs them on a Dispatcher. The closure of a gRPC stream, by the
// client or by the server, will prompt Tau to shutdown the MessageQueue.
type Tau struct {
	dispatcher *Dispatcher
}

// NewTau returns a new Tau service that will run MessageQueues on the given
// Dispatcher.
func NewTau(dispatcher *Dispatcher) Tau {
	return Tau{
		dispatcher: dispatcher,
	}
}

// Connect to the Tau service and begin streaming requests and responses for
// the Tau sMPC protocol.
func (τ *Tau) Connect(stream rpc.TauService_ConnectServer) error {

	// Use a background MessageQueue to handle the connection until an error
	// is returned by the MessageQueue
	ch := make(chan error, 1)
	quit := make(chan struct{}, 1)
	go func() { ch <- τ.connect(stream, quit) }()
	defer close(quit)

	// Select between the context finishing and the background worker
	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	case err := <-ch:
		return err
	}
}

func (τ *Tau) connect(stream rpc.TauService_ConnectServer, quit chan struct{}) error {

	// Create a MessageQueue that owns this gRPC stream and run it on the
	// Dispatcher
	messageQueue := NewMessageQueue(stream)
	multiAddress, err := identity.NewMultiAddressFromString("unimplemented")
	if err != nil {
		return err
	}

	// Shutdown the MessageQueue
	go func() {
		<-quit
		τ.dispatcher.ShutdownMessageQueue(multiAddress)
	}()

	return τ.dispatcher.RunMessageQueue(multiAddress, &messageQueue)
}
