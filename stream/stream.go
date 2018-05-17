package stream

import (
	"context"
	"encoding"
	"errors"
	"sync"

	"github.com/republicprotocol/republic-go/identity"
)

// ErrDisconnectedClosedStream is returned when a call to
// StreamClient.Disconnect happens for an identity.MultiAddress that is not
// connected.
var ErrDisconnectedClosedStream = errors.New("disconnected closed stream")

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
	Close() error
	Send(Message) error
	Recv(Message) error
}

type Streamer interface {
	Open(ctx context.Context, multiAddr identity.MultiAddress) (Stream, error)
	Close(multiAddr identity.MultiAddress) error
}

type streamer struct {
	addr   identity.Address
	client Client
	server Server

	streamsMu *sync.Mutex
	streams   map[identity.Address]Stream
}

func NewStreamer(addr identity.Address, client Client, server Server) Streamer {
	return &streamer{
		addr:   addr,
		client: client,
		server: server,
	}
}

func (streamer streamer) Open(ctx context.Context, multiAddr identity.MultiAddress) (Stream, error) {
	if streamer.addr < multiAddr.Address() {
		return streamer.client.Connect(ctx, multiAddr)
	}

	stream, err := streamer.server.Listen(ctx, multiAddr.Address())
	if err != nil {
		return nil, err
	}

	streamer.streamsMu.Lock()
	defer streamer.streamsMu.Unlock()
	streamer.streams[multiAddr.Address()] = stream

	return stream, nil
}

func (streamer streamer) Close(multiAddr identity.MultiAddress) error {
	if streamer.addr < multiAddr.Address() {
		return streamer.client.Disconnect(multiAddr)
	}

	streamer.streamsMu.Lock()
	defer streamer.streamsMu.Unlock()

	if stream, ok := streamer.streams[multiAddr.Address()]; ok {
		err := stream.Close()
		delete(streamer.streams, multiAddr.Address())
		return err
	}

	return nil
}

// Client is an interface for connecting to a Server.
type Client interface {

	// Connect to a Server identified by an identity.MultiAddress. Returns a
	// Stream for sending and receiving Messages to and from the Server. To be
	// used with a ClientRecycler, the Stream must be safe for concurrent use.
	// A call to Client.Connect for an identity.MultiAddress must eventually be
	// followed by a call to Client.Disconnect for the same
	// identity.MultiAddress.
	Connect(ctx context.Context, multiAddr identity.MultiAddress) (Stream, error)

	// Disconnect from a Server identified by an identity.MultiAddress. A call
	// to Client.Connect for an identity.MultiAddress must eventually be
	// followed by a call to Client.Disconnect for the same
	// identity.MultiAddress. Returns an ErrDisconnectedClosedStream when there
	// is no previous call to Client.Connect.
	Disconnect(multiAddr identity.MultiAddress) error
}

// Server is an interface for accepting connections from a Client.
type Server interface {

	// Listen for a connection from a Client identified by an identity.Address.
	Listen(ctx context.Context, addr identity.Address) (Stream, error)
}

// clientRecycler encapsulates a Client and reuses a Stream that has been
// connected to Server when multiple connections to the Server are needd. It
// does not protect the Stream from concurrent use.
type clientRecycler struct {
	client Client

	connsMu *sync.Mutex
	connsRc map[identity.Address]int64
	conns   map[identity.Address]Stream
}

// NewClientRecycler returns a Client that wraps another Client. It will
// use the Client to create Stream, but will recycle connected Streams when
// multiple connections to the same Server are needed. The wrapped Client must
// ensure that the Stream is safe for concurrent use.
func NewClientRecycler(client Client) Client {
	return &clientRecycler{
		client:  client,
		connsMu: new(sync.Mutex),
		connsRc: map[identity.Address]int64{},
		conns:   map[identity.Address]Stream{},
	}
}

// Connect implements the Client interface.
func (client *clientRecycler) Connect(ctx context.Context, multiAddr identity.MultiAddress) (Stream, error) {
	client.connsMu.Lock()
	defer client.connsMu.Unlock()

	addr := multiAddr.Address()
	if _, ok := client.conns[addr]; !ok {
		conn, err := client.client.Connect(ctx, multiAddr)
		if err != nil {
			return nil, err
		}
		client.conns[addr] = conn
		client.connsRc[addr] = 0
	}
	client.connsRc[addr]++
	return client.conns[addr], nil
}

// Disconnect implements the Client interface.
func (client *clientRecycler) Disconnect(multiAddr identity.MultiAddress) error {
	client.connsMu.Lock()
	defer client.connsMu.Unlock()

	addr := multiAddr.Address()
	if client.connsRc[addr] == 0 {
		return ErrDisconnectedClosedStream
	}

	client.connsRc[addr]--
	if client.connsRc[addr] == 0 {
		if err := client.client.Disconnect(multiAddr); err != nil {
			return err
		}
		delete(client.conns, addr)
		delete(client.connsRc, addr)
	}
	return nil
}
