package stream

import (
	"context"
	"encoding"
	"errors"
	"sync"

	"github.com/republicprotocol/republic-go/identity"
)

// ErrSendOnClosedStream is returned when a call to Stream.Send happens on a
// closed Stream.
var ErrSendOnClosedStream = errors.New("send on closed stream")

// ErrRecvOnClosedStream is returned when a call to Stream.Recv happens on a
// closed Stream.
var ErrRecvOnClosedStream = errors.New("receive on closed stream")

// Message is an interface for data that can be sent over a bidirectional
// stream between nodes.
type Message interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	// Types that can be used in an Stream must implement this pass through
	// method. It only exists to restrict Stream to types that have been
	// explicitly marked as compatible to avoid programmer error.
	IsMessage()
}

// Stream is an interface for sending and receiving Messages over a
// bidirectional stream. It abstracts over the client and server architecture.
type Stream interface {
	// Send a Message on the Stream. Calls to Stream.Send might be blocking
	// depending on the underlying implementation.
	Send(Message) error

	// Recv a Message from the Stream. Calls to Stream.Recv will block until
	// a Message is received.
	Recv(Message) error
}

// Client is an interface for connecting to a Server.
type Client interface {

	// Connect to a Server identified by an identity.MultiAddress. Returns a
	// Stream for sending, and receiving, Messages to, and from, the Server.
	// The context.Context can be used to close the Stream.
	Connect(ctx context.Context, multiAddr identity.MultiAddress) (Stream, error)
}

// Server is an interface for accepting connections from a Client.
type Server interface {

	// Listen for a connection from a Client identified by an identity.Address.
	// Returns a Stream for sending, and receiving, Messages to, and from, the
	// Client. The context.Context can be used to close the Stream.
	Listen(ctx context.Context, addr identity.Address) (Stream, error)
}

// Streamer abstracts over the Client and Server model. By comparing
// identity.Addresses it determines whether opening a Stream should be done by
// listening for a connection as a Server, or connecting to a Server as a
// Client.
type Streamer interface {

	// Open a Stream to an identity.MultiAddress by listening for a Client
	// connection, or connecting to a Server. Calls to Streamer.Open are
	// blocking. The context.Context can be used to close the Stream.
	Open(ctx context.Context, multiAddr identity.MultiAddress) (Stream, error)
}

type streamer struct {
	addr   identity.Address
	client Client
	server Server
}

// NewStreamer returns a Streamer that uses an identity.Address to identify
// itself. It will use the Client to connect streams when opening streams to an
// identity.Address greater than its own, and it will use the Server to listen
// for connections when opening streams to an identity.Address lower than its
// own.
func NewStreamer(addr identity.Address, client Client, server Server) Streamer {
	return &streamer{
		addr:   addr,
		client: client,
		server: server,
	}
}

// Open implements the Streamer interface.
func (streamer streamer) Open(ctx context.Context, multiAddr identity.MultiAddress) (Stream, error) {
	addr := multiAddr.Address()
	if streamer.addr < addr {
		return streamer.client.Connect(ctx, multiAddr)
	}
	return streamer.server.Listen(ctx, addr)
}

type streamRecycler struct {
	streamer Streamer

	mu            *sync.Mutex
	streams       map[identity.Address]Stream
	streamsCancel map[identity.Address]context.CancelFunc
	streamsRc     map[identity.Address]int
	streamsMu     map[identity.Address]*sync.Mutex
}

// NewStreamRecycler returns a Streamer that wraps another Streamer. It will
// use the wrapper Streamer to open and close Streams, but will recycle
// existing Streams when multiple connections to the same identity.Address are
// needed. Streams opened will be safe for concurrent use whenever the wrapped
// Streamer can open Streams that are safe for concurrent use.
func NewStreamRecycler(streamer Streamer) Streamer {
	return &streamRecycler{
		streamer: streamer,

		mu:            new(sync.Mutex),
		streams:       map[identity.Address]Stream{},
		streamsCancel: map[identity.Address]context.CancelFunc{},
		streamsRc:     map[identity.Address]int{},
		streamsMu:     map[identity.Address]*sync.Mutex{},
	}
}

// Open implements the Streamer interface.
func (recycler *streamRecycler) Open(ctx context.Context, multiAddr identity.MultiAddress) (Stream, error) {
	addr := multiAddr.Address()

	mu := func() *sync.Mutex {
		recycler.mu.Lock()
		defer recycler.mu.Unlock()

		if recycler.streamsMu[addr] == nil {
			recycler.streamsMu[addr] = new(sync.Mutex)
		}
		return recycler.streamsMu[addr]
	}()
	mu.Lock()
	defer mu.Unlock()

	recycler.mu.Lock()
	defer recycler.mu.Unlock()

	if recycler.streamsRc[addr] == 0 {
		recycler.mu.Unlock()
		ctx, cancel := context.WithCancel(context.Background())
		stream, err := recycler.streamer.Open(ctx, multiAddr)
		recycler.mu.Lock()

		if err != nil {
			cancel()
			return stream, err
		}

		recycler.streams[addr] = stream
		recycler.streamsCancel[addr] = cancel
	}

	go func() {
		<-ctx.Done()

		recycler.mu.Lock()
		defer recycler.mu.Unlock()

		recycler.streamsRc[addr]--
		if recycler.streamsRc[addr] == 0 {
			recycler.streamsCancel[addr]()
			delete(recycler.streams, addr)
			delete(recycler.streamsCancel, addr)
			delete(recycler.streamsRc, addr)
			delete(recycler.streamsMu, addr)
		}
	}()

	recycler.streamsRc[addr]++
	return recycler.streams[addr], nil
}
