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

// ErrCloseOnClosedStream is returned when a call to CloseStream.Close happens
// on a closed Stream.
var ErrCloseOnClosedStream = errors.New("close on closed stream")

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

// A CloseStream is a Stream that needs to be explicitly closed.
type CloseStream interface {
	Stream

	// Close the Stream and signal that it is no longer needed.
	Close() error
}

// Client is an interface for connecting to a Server.
type Client interface {

	// Connect to a Server identified by an identity.MultiAddress. Returns a
	// CloseStream for sending and receiving Messages to and from the Server.
	// The CloseStream that must be closed when the CloseStream is no longer
	// needed, otherwise resources will leak.
	Connect(ctx context.Context, multiAddr identity.MultiAddress) (CloseStream, error)
}

// Server is an interface for accepting connections from a Client.
type Server interface {

	// Listen for a connection from a Client identified by an identity.Address.
	// Returns a CloseStream that must be closed when the CloseStream is no
	// longer needed, otherwise resources will leak.
	Listen(ctx context.Context, addr identity.Address) (CloseStream, error)
}

// Streamer abstracts over the Client and Server architecture. By comparing
// identity.Addresses it determines whether opening a Stream should be done by
// listening for a connection as a Server, or connecting to a Server as a
// Client.
type Streamer interface {

	// Open a Stream to an identity.MultiAddress by listening for a Client
	// connection, or actively connecting to a Server. Calls to Streamer.Open
	// are blocking and must be accompanied by a call to Streamer.Close.
	Open(ctx context.Context, multiAddr identity.MultiAddress) (Stream, error)

	// Close a Stream with an identity.Address. Calls to Streamer.Open must be
	// accompanied by a call to Streamer.Close.
	Close(addr identity.Address) error
}

type streamer struct {
	addr   identity.Address
	client Client
	server Server

	streamsMu *sync.Mutex
	streams   map[identity.Address]CloseStream
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

		streamsMu: new(sync.Mutex),
		streams:   map[identity.Address]CloseStream{},
	}
}

// Open implements the Streamer interface.
func (streamer streamer) Open(ctx context.Context, multiAddr identity.MultiAddress) (Stream, error) {
	var stream CloseStream
	var err error

	addr := multiAddr.Address()
	if addr < streamer.addr {
		stream, err = streamer.client.Connect(ctx, multiAddr)
	} else {
		stream, err = streamer.server.Listen(ctx, addr)
	}
	if err != nil {
		return stream, err
	}

	streamer.streamsMu.Lock()
	defer streamer.streamsMu.Unlock()

	streamer.streams[addr] = stream
	return stream, nil
}

// Close implements the Streamer interface.
func (streamer streamer) Close(addr identity.Address) error {
	streamer.streamsMu.Lock()
	defer streamer.streamsMu.Unlock()

	if stream, ok := streamer.streams[addr]; ok {
		err := stream.Close()
		delete(streamer.streams, addr)
		return err
	}
	return nil
}

type streamRecycler struct {
	streamer Streamer

	streamsMu *sync.Mutex
	streamsRc map[identity.Address]int64
	streams   map[identity.Address]Stream
}

// NewStreamRecycler returns a Streamer that wraps another Streamer. It will
// use the Streamer to open and close Streams, but will recycle existing
// Streams when multiple connections to the same identity.Address are needed.
// Streams opened will be safe for concurrent use whenever the inner Streamer
// can open Streams that are safe for concurrent use.
func NewStreamRecycler(streamer Streamer) Streamer {
	return &streamRecycler{
		streamer: streamer,

		streamsMu: new(sync.Mutex),
		streamsRc: map[identity.Address]int64{},
		streams:   map[identity.Address]Stream{},
	}
}

// Open implements the Streamer interface.
func (recycler *streamRecycler) Open(ctx context.Context, multiAddr identity.MultiAddress) (Stream, error) {
	recycler.streamsMu.Lock()
	defer recycler.streamsMu.Unlock()

	addr := multiAddr.Address()
	if _, ok := recycler.streams[addr]; !ok {
		conn, err := recycler.streamer.Open(ctx, multiAddr)
		if err != nil {
			return nil, err
		}
		recycler.streams[addr] = conn
		recycler.streamsRc[addr] = 0
	}
	recycler.streamsRc[addr]++
	return recycler.streams[addr], nil
}

// Close implements the Streamer interface.
func (recycler *streamRecycler) Close(addr identity.Address) error {
	recycler.streamsMu.Lock()
	defer recycler.streamsMu.Unlock()

	if recycler.streamsRc[addr] == 0 {
		return ErrCloseOnClosedStream
	}

	recycler.streamsRc[addr]--
	if recycler.streamsRc[addr] == 0 {
		if err := recycler.streamer.Close(addr); err != nil {
			return err
		}
		delete(recycler.streams, addr)
		delete(recycler.streamsRc, addr)
	}
	return nil
}
