package smpcer

import (
	"sync"

	"google.golang.org/grpc"
)

// RendezvousConn is created and returned by a Rendezvous.
type RendezvousConn struct {
	Sender   chan *ComputeMessage
	Receiver chan *ComputeMessage
}

func (conn *RendezvousConn) stream(stream grpc.Stream) error {
	// The done channel will signal to the sender goroutine that it should
	// exit
	done := make(chan struct{})
	defer close(done)

	errs := make(chan error, 1)
	go func() {
		defer close(errs)

		for {
			select {
			case <-done:
				// When the receiver exits the done channel will be closed and
				// this goroutine will eventually exit
				return
			case <-stream.Context().Done():
				errs <- stream.Context().Err()
				return
			case message, ok := <-conn.Sender:
				if !ok {
					return
				}
				if err := stream.SendMsg(message); err != nil {
					errs <- err
					return
				}
			}
		}
	}()

	// Receive messages from the client until the context is done, or an error
	// is received
	for {
		message := new(ComputeMessage)
		if err := stream.RecvMsg(message); err != nil {
			return err
		}

		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case err, ok := <-errs:
			if !ok {
				// When the sender exits this error channel will be closed and
				// this goroutine will eventually exit
				return nil
			}
			return err
		case conn.Receiver <- message:
		}
	}
}

// Rendezvous maintains bi-directional streams between gRPC clients and gRPC
// services.
type Rendezvous struct {
	connsMu *sync.Mutex
	connsRc map[string]int
	conns   map[string]RendezvousConn
}

func newRendezvous() Rendezvous {
	return Rendezvous{
		connsMu: new(sync.Mutex),
		connsRc: map[string]int{},
		conns:   map[string]RendezvousConn{},
	}
}

func (r *Rendezvous) acquireConn(addr string) RendezvousConn {
	r.connsMu.Lock()
	defer r.connsMu.Unlock()

	if r.connsRc[addr] == 0 {
		r.conns[addr] = RendezvousConn{
			Sender:   make(chan *ComputeMessage),
			Receiver: make(chan *ComputeMessage),
		}
	}
	r.connsRc[addr]++
	return r.conns[addr]
}

func (r *Rendezvous) releaseConn(addr string) {
	r.connsMu.Lock()
	defer r.connsMu.Unlock()

	r.connsRc[addr]--
	if r.connsRc[addr] == 0 {
		close(r.conns[addr].Sender)
		close(r.conns[addr].Receiver)
		delete(r.conns, addr)
	}
}
