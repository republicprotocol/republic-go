package dispatch

import (
	"sync/atomic"
)

// MaxListeners defines the maximum number of Listeners that can be subscribed
// to a Broadcaster. At this limit, the Broadcaster will return errors when a
// Listener attempts to subscribe.
// TODO: Make this constant configurable.
const MaxListeners = int32(128)

type Broadcaster struct {
	done         chan struct{}
	ch           chan interface{}
	listeners    chan listener
	numListeners int32
}

func NewBroadcaster() *Broadcaster {
	broadcaster := &Broadcaster{
		done:         make(chan struct{}),
		ch:           make(chan interface{}, MaxListeners),
		listeners:    make(chan listener, MaxListeners),
		numListeners: 0,
	}

	go func() {
		listeners := [MaxListeners]listener{}
		defer func() {
			for i := int32(0); i < broadcaster.numListeners; i++ {
				close(listeners[i].ch)
			}
		}()

		for {
			select {
			case <-broadcaster.done:
				return
			case msg, ok := <-broadcaster.ch:
				if !ok {
					// The broadcasting channel is intentionally left open to
					// prevent writes panicking
					panic("broadcasting channel closed")
				}

				for i := int32(0); i < broadcaster.numListeners; i++ {
					listener := listeners[i]
					select {
					case <-broadcaster.done:
						// All listeners will be cleaned up during the defer
						// phase and the closure of broadcaster signals that we
						// should return
						return

					case <-listener.done:
						// The listener has signaled that they are done, but
						// the broadcaster is not necessarily done, so we clean
						// up the listener immediately
						close(listener.ch)
						atomic.AddInt32(&broadcaster.numListeners, -1)
						listeners[i] = listeners[broadcaster.numListeners]
						i--

					case listener.ch <- msg:
					}
				}
			case listener, ok := <-broadcaster.listeners:
				if !ok {
					// The broadcasting listeners are intentionally left open
					// to prevent writes panicking
					panic("broadcasting listeners closed")
				}
				if atomic.LoadInt32(&broadcaster.numListeners) >= MaxListeners {
					close(listener.ch)
					continue
				}
				listeners[broadcaster.numListeners] = listener
				atomic.AddInt32(&broadcaster.numListeners, 1)
			}
		}
	}()

	return broadcaster
}

func (broadcaster *Broadcaster) Close() {
	close(broadcaster.done)
}

func (broadcaster *Broadcaster) Broadcast(done <-chan struct{}, ch <-chan interface{}) {
	for {
		select {
		case <-done:
			return
		case <-broadcaster.done:
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			select {
			case <-done:
				return
			case <-broadcaster.done:
				return
			case broadcaster.ch <- msg:
			}
		}
	}
}

func (broadcaster *Broadcaster) Listen(done <-chan struct{}) <-chan interface{} {
	lis := listener{
		done: done,
		ch:   make(chan interface{}),
	}
	if atomic.LoadInt32(&broadcaster.numListeners) >= MaxListeners {
		close(lis.ch)
		return lis.ch
	}
	select {
	case <-done:
		close(lis.ch)
	case <-broadcaster.done:
		close(lis.ch)
	case broadcaster.listeners <- lis:
	}
	return lis.ch
}

type listener struct {
	done <-chan struct{}
	ch   chan interface{}
}
