package smpcer

import (
	"errors"
	"sync"

	"github.com/republicprotocol/republic-go/dispatch"
	"google.golang.org/grpc"
)

// MaxConnections that can be serviced by an Smpc service.
const MaxConnections = 256

var ErrRendezvousIsClosed = errors.New("rendezvous is closed")

// A Rendezvous is a collection of channels for sending, and receiving,
// ComputeMessages, and errors, between a client and an Smpc service. It is
// used by clients to wait for Smpc services to accept their connection
// requests, and by Smpc service to wait for clients to request a connection.
type Rendezvous struct {
	Sender    chan *ComputeMessage
	Receiver  chan *ComputeMessage
	senders   chan (<-chan *ComputeMessage)
	receivers chan (chan<- *ComputeMessage)

	qMu *sync.RWMutex
	q   bool
}

func NewRendezvous(done <-chan struct{}) *Rendezvous {
	rendezvous := &Rendezvous{
		Sender:    make(chan *ComputeMessage),
		Receiver:  make(chan *ComputeMessage),
		senders:   make(chan (<-chan *ComputeMessage)),
		receivers: make(chan (chan<- *ComputeMessage)),

		qMu: new(sync.RWMutex),
		q:   false,
	}
	go rendezvous.mergeSenders(done)
	go rendezvous.splitReceivers(done)
	return rendezvous
}

func (rendezvous *Rendezvous) Connect(conn *RendezvousConn) error {
	rendezvous.qMu.RLock()
	defer rendezvous.qMu.RUnlock()
	if rendezvous.q {
		return ErrRendezvousIsClosed
	}
	rendezvous.senders <- conn.Sender
	rendezvous.receivers <- conn.Receiver
	return nil
}

func (rendezvous *Rendezvous) mergeSenders(done <-chan struct{}) {
	defer close(rendezvous.Sender)
	defer close(rendezvous.senders)
	defer rendezvous.quit()

	var wg sync.WaitGroup
	quit := false
	for !quit {
		select {
		case <-done:
			quit = true
		case sender, ok := <-rendezvous.senders:
			if !ok {
				quit = true
				break
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case <-done:
						return
					case send, ok := <-sender:
						if !ok {
							return
						}
						select {
						case <-done:
							return
						case rendezvous.Sender <- send:
						}
					}

				}
			}()
		}
	}
	wg.Wait()
}

func (rendezvous *Rendezvous) splitReceivers(done <-chan struct{}) {
	defer close(rendezvous.Receiver)
	defer close(rendezvous.receivers)
	defer rendezvous.quit()

	receivers := []chan<- *ComputeMessage{}
	defer func() {
		for _, receiver := range receivers {
			close(receiver)
		}
	}()

	for {
		select {
		case <-done:
			return
		case receiver, ok := <-rendezvous.receivers:
			if !ok {
				return
			}
			receivers = append(receivers, receiver)
		case recv, ok := <-rendezvous.Receiver:
			if !ok {
				return
			}
			for _, receiver := range receivers {
				select {
				case <-done:
					return
				case receiver <- recv:
				}
			}
		}
	}
}

func (rendezvous *Rendezvous) quit() {
	rendezvous.qMu.Lock()
	defer rendezvous.qMu.Unlock()
	rendezvous.q = true
}

type RendezvousConn struct {
	Sender   chan *ComputeMessage
	Receiver chan *ComputeMessage
}

func NewRendezvousConn() *RendezvousConn {
	return &RendezvousConn{
		Sender:   make(chan *ComputeMessage),
		Receiver: make(chan *ComputeMessage),
	}
}

func (conn *RendezvousConn) Close() {
	close(conn.Sender)
}

// A RendezvousRouter stores Rendezvous between clients and Smpc services. It
// is used to acquire and release Rendezvous (see RendezvousRouter.Acquire and
// RendezvousRouter.Release). Calls to RendezvousRouter.Acquire reuse
// Rendezvous when one is already open.
type RendezvousRouter struct {
	rendezvousMu       *sync.Mutex
	rendezvousSplitter map[string]*dispatch.Splitter
}

// NewRendezvousRouter returns an empty RendezvousRouter.
func NewRendezvousRouter() RendezvousRouter {
	return RendezvousRouter{
		rendezvousMu:       new(sync.Mutex),
		rendezvousSplitter: map[string]*dispatch.Splitter{},
	}
}

// waitForClient to have its Compute RPC connection accepted by the address.
// Once the connection is accepted, all messages written to the sender channel
// will be forwarded to the address, and all messages received from the address
// will be written to the returned receiver channel. Calls to
// Rendezvous.waitForClient must only be made by a Client.
func (router *RendezvousRouter) waitForClient(addr string, sender <-chan *ComputeMessage) (<-chan *ComputeMessage, <-chan error) {
	router.rendezvousMu.Lock()
	defer router.rendezvousMu.Unlock()

	panic("unimplemented")
}

// waitForService to accept a Compute RPC connection from the address. Once the
// connection is accepted, all messages written to the sender channel will be
// forwarded to the address, and all messages received from the address will be
// written to the returned receiver channel. Calls to Rendezvous.waitForService
// must only be made by a Client.
func (router *RendezvousRouter) waitForService(addr string, sender <-chan *ComputeMessage) (<-chan *ComputeMessage, <-chan error) {
	router.rendezvousMu.Lock()
	defer router.rendezvousMu.Unlock()

	panic("unimplemented")
}

// connect an accepted Compute RPC connection from an address. Messages written
// to the sender channels passed in from calls to Rendezvous.waitForClient and
// Rendezvous.waitForService are written to the stream. Messages read from the
// stream are split and written to all receiver channels created in calls to
// Rendezvous.waitForClient and Rendezvous.waitForService. Calls to
// Rendezvous.connect must only be made by an Smpc service.
func (router *RendezvousRouter) connect(addr string, stream Smpc_ComputeServer) error {
	router.rendezvousMu.Lock()
	defer router.rendezvousMu.Unlock()

	panic("unimplemented")
}

func handleStream(stream grpc.Stream, in <-chan *ComputeMessage, out chan<- *ComputeMessage) error {
	defer func() {
		if clientStream, ok := stream.(Smpc_ComputeClient); ok {
			clientStream.CloseSend()
		}
	}()

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
			case computeMessage, ok := <-in:
				if !ok {
					return
				}
				if err := stream.SendMsg(computeMessage); err != nil {
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
		case out <- message:
		}
	}
}
