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
	rcs       map[identity.Address]int
}

func NewRendezvous() Rendezvous {
	return Rendezvous{
		mu:        new(sync.Mutex),
		senders:   map[identity.Address]chan *ComputeMessage{},
		receivers: map[identity.Address]*dispatch.Splitter{},
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

		dispatch.Pipe(done, sender, rendezvous.senders[addr])
	}()
	return receiver, errs
}

func (rendezvous *Rendezvous) acquireConn(addr identity.Address) {
	rendezvous.mu.Lock()
	defer rendezvous.mu.Unlock()

	if rendezvous.rcs[addr] == 0 {
		rendezvous.senders[addr] = make(chan *ComputeMessage)
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
		delete(rendezvous.rcs, addr)
	}
}

type listener struct {
	ch   chan *ComputeMessage
	done <-chan struct{}
}

type dispatcher struct {
	ch        chan *ComputeMessage
	listeners chan listener
}

func newDispatcher(done <-chan struct{}) *dispatcher {
	d := &dispatcher{
		ch:        make(chan *ComputeMessage),
		listeners: make(chan listener),
	}

	go func() {
		listeners := map[chan *ComputeMessage](<-chan struct{}){}
		defer func() {
			for listener, listenerDone := range listeners {
				<-listenerDone
				close(listener)
			}
		}()

		for {
			select {
			case <-done:
				return
			case msg, ok := <-d.ch:
				if !ok {
					return
				}
				deleteListener := func(listener chan *ComputeMessage) {
					close(listener)
					delete(listeners, listener)
				}
				for listener, listenerDone := range listeners {
					select {
					case <-done:
						deleteListener(listener)
					case <-listenerDone:
						deleteListener(listener)
					case listener <- msg:
					}
				}
			case listener, ok := <-d.listeners:
				if !ok {
					return
				}
				listeners[listener.ch] = listener.done
			}
		}
	}()

	return d
}

func (d *dispatcher) broadcast(done <-chan struct{}, ch <-chan *ComputeMessage) {
	for {
		select {
		case <-done:
			return
		case msg, ok := <-ch:
			if !ok {
				<-done
				return
			}
			select {
			case <-done:
				return
			case d.ch <- msg:
			}
		}
	}
}

func (d *dispatcher) listen(done <-chan struct{}) <-chan *ComputeMessage {
	listener := listener{
		ch:   make(chan *ComputeMessage),
		done: done,
	}
	select {
	case <-done:
		close(listener.ch)
	case d.listeners <- listener:
	}
	return listener.ch
}
