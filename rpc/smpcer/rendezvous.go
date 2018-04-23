package smpcer

import (
	"sync"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
)

type Rendezvous struct {
	mu        *sync.Mutex
	senders   map[identity.Address]chan *ComputeMessage
	receivers map[identity.Address]*dispatch.Splitter
	errs      map[identity.Address]*dispatch.Splitter
	rcs       map[identity.Address]int
}

func NewRendezvous() Rendezvous {
	return Rendezvous{
		senders:   map[identity.Address]chan *ComputeMessage{},
		receivers: map[identity.Address]*dispatch.Splitter{},
		errs:      map[identity.Address]*dispatch.Splitter{},
		rcs:       map[identity.Address]int{},
	}
}

// connect an accepted Compute RPC connection from an address. Messages written
// to the sender channels passed in from calls to Rendezvous.waitForClient and
// Rendezvous.waitForService are written to the stream. Messages read from the
// stream are split and written to all receiver channels created in calls to
// Rendezvous.waitForClient and Rendezvous.waitForService. Calls to
// Rendezvous.connect must only be made by an Smpc service.
func (rendezvous *Rendezvous) connect(addr identity.Address, done <-chan struct{}, receiver <-chan *ComputeMessage) <-chan *ComputeMessage {
	sender := make(chan *ComputeMessage)
	go func() {
		defer close(sender)

		rendezvous.acquireConn(addr)
		defer rendezvous.releaseConn(addr)

		dispatch.CoBegin(func() {
			rendezvous.receivers[addr].Split(receiver)
		}, func() {
			dispatch.Pipe(done, rendezvous.senders[addr], sender)
		})
	}()
	return sender
}

// wait for the Rendezvous to be connected to from an address. The Rendezvous
// read messages from the sender channel and forward them to the address. The
// user can read messages sent by the address from the returned channel. Calls
// to Rendezvous.wait must only be made by a Client.
func (rendezvous *Rendezvous) wait(addr identity.Address, done <-chan struct{}, sender <-chan *ComputeMessage) (<-chan *ComputeMessage, <-chan error) {
	receiver := make(chan *ComputeMessage)
	errs := make(chan error, 1)

	go func() {
		defer close(receiver)
		defer close(errs)

		rendezvous.acquireConn(addr)
		defer rendezvous.releaseConn(addr)

		if err := rendezvous.receivers[addr].Subscribe(receiver); err != nil {
			errs <- err
			return
		}
		defer rendezvous.receivers[addr].Unsubscribe(receiver)
	}()
	return receiver, errs
}

func (rendezvous *Rendezvous) acquireConn(addr identity.Address) {
	rendezvous.mu.Lock()
	defer rendezvous.mu.Unlock()

	if rendezvous.rcs[addr] == 0 {
		sender := make(chan *ComputeMessage)
		rendezvous.senders[addr] = sender
		rendezvous.receivers[addr] = dispatch.NewSplitter(MaxConnections)
	}
	rendezvous.rcs[addr]++
}

func (rendezvous *Rendezvous) releaseConn(addr identity.Address) {
	rendezvous.mu.Lock()
	defer rendezvous.mu.Unlock()

	rendezvous.rcs[addr]--
	if rendezvous.rcs[addr] == 0 {
		close(rendezvous.senders[addr])
		delete(rendezvous.senders, addr)
		delete(rendezvous.receivers, addr)
		delete(rendezvous.errs, addr)
	}
}

// func (connector *Connector) compute(ctx context.Context, multiAddress identity.MultiAddress) (Smpc_ComputeClient, error) {

// 	// Dial the client.ConnPool for a client.Conn
// 	conn, err := connector.connPool.Dial(ctx, multiAddress)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot dial %v: %v", multiAddress, err)
// 	}
// 	defer conn.Close()

// 	// Create an SmpcClient and call the Compute RPC
// 	smpcClient := NewSmpcClient(conn.ClientConn)
// 	stream, err := smpcClient.Compute(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot open stream to %v: %v", multiAddress, err)
// 	}

// 	// Send an authentication message to the Smpc service
// 	auth := &ComputeMessage{
// 		Signature: []byte{}, // FIXME: Provide verifiable signature
// 		Value: &ComputeMessage_Address{
// 			Address: connector.multiAddress.Address().String(),
// 		},
// 	}
// 	if err := stream.Send(auth); err != nil {
// 		return nil, fmt.Errorf("cannot authenticate with %v: %v", multiAddress, err)
// 	}

// 	return stream, nil
// }
