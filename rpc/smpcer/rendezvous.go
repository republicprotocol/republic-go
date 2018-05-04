package smpcer

import (
	"sync"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
)

type Rendezvous struct {
	mu        *sync.Mutex
	rcs       map[identity.Address]int
	senders   map[identity.Address]chan *ComputeMessage
	receivers map[identity.Address]*dispatch.Broadcaster
}

func NewRendezvous() Rendezvous {
	return Rendezvous{
		mu:        new(sync.Mutex),
		rcs:       map[identity.Address]int{},
		senders:   map[identity.Address]chan *ComputeMessage{},
		receivers: map[identity.Address]*dispatch.Broadcaster{},
	}
}

// connect an accepted Compute RPC connection from an address. Messages written
// to the sender channels passed in from calls to Rendezvous.waitForClient and
// Rendezvous.waitForService are written to the stream. Messages read from the
// stream are split and written to all receiver channels created in calls to
// Rendezvous.waitForClient and Rendezvous.waitForService. Calls to
// Rendezvous.connect must only be made by an Smpc service.
func (rendezvous *Rendezvous) connect(addr identity.Address, done <-chan struct{}, receiver <-chan interface{}) <-chan *ComputeMessage {
	sender := make(chan *ComputeMessage)
	go func() {
		defer close(sender)
		defer rendezvous.releaseConn(addr)

		rendezvous.acquireConn(addr)
		dispatch.CoBegin(func() {
			rendezvous.receivers[addr].Broadcast(done, receiver)
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
func (rendezvous *Rendezvous) wait(addr identity.Address, done <-chan struct{}, sender <-chan interface{}) <-chan interface{} {
	rendezvous.acquireConn(addr)
	receiver := rendezvous.receivers[addr].Listen(done)
	go func() {
		defer rendezvous.releaseConn(addr)
		dispatch.Pipe(done, sender, rendezvous.senders[addr])
	}()
	return receiver
}

func (rendezvous *Rendezvous) acquireConn(addr identity.Address) {
	rendezvous.mu.Lock()
	defer rendezvous.mu.Unlock()

	rendezvous.rcs[addr]++
	if rendezvous.rcs[addr] == 1 {
		rendezvous.senders[addr] = make(chan *ComputeMessage)
		rendezvous.receivers[addr] = dispatch.NewBroadcaster()
	}
}

func (rendezvous *Rendezvous) releaseConn(addr identity.Address) {
	rendezvous.mu.Lock()
	defer rendezvous.mu.Unlock()

	rendezvous.rcs[addr]--
	if rendezvous.rcs[addr] == 0 {
		close(rendezvous.senders[addr])
		delete(rendezvous.rcs, addr)
		delete(rendezvous.senders, addr)
		delete(rendezvous.receivers, addr)
	}
}
