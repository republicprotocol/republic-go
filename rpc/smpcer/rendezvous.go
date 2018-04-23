package smpcer

import (
	"errors"
	"sync"
)

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
	rendezvousMu     *sync.Mutex
	rendezvousRc     map[string]int
	rendezvousQ      map[string]chan struct{}
	rendezvousRoutes map[string]*Rendezvous
}

// NewRendezvousRouter returns an empty RendezvousRouter.
func NewRendezvousRouter() RendezvousRouter {
	return RendezvousRouter{
		rendezvousMu:     new(sync.Mutex),
		rendezvousRc:     map[string]int{},
		rendezvousQ:      map[string]chan struct{}{},
		rendezvousRoutes: map[string]*Rendezvous{},
	}
}

// Acquire a Rendezvous with an address. If a Rendezvous with that address is
// already open, then it will be reused. Otherwise, a new Rendezvous will be
// created. A call to RendezvousRouter.Acquire must be accompanied by exactly
// one call to RendezvousRouter.Release, using the same address.
func (router *RendezvousRouter) Acquire(addr string) *RendezvousConn {
	router.rendezvousMu.Lock()
	defer router.rendezvousMu.Unlock()

	if router.rendezvousRc[addr] == 0 {
		router.rendezvousQ[addr] = make(chan struct{})
		router.rendezvousRoutes[addr] = NewRendezvous(router.rendezvousQ[addr])
	}
	router.rendezvousRc[addr]++

	conn := NewRendezvousConn()
	router.rendezvousRoutes[addr].Connect(conn)
	return conn
}

// Release a Rendezvous with an address. If no other references to the
// Rendezvous exists, then the Rendezvous will be closed and deleted.
func (router *RendezvousRouter) Release(addr string) {
	router.rendezvousMu.Lock()
	defer router.rendezvousMu.Unlock()

	router.rendezvousRc[addr]--
	if router.rendezvousRc[addr] == 0 {
		close(router.rendezvousQ[addr])
		delete(router.rendezvousRoutes, addr)
	}
}
