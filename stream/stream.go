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
	Send(Message) error
	Recv(Message) error
}

// Client is an interface for connecting to a Server.
type Client interface {
	Connect(ctx context.Context, multiAddr identity.MultiAddress) (Stream, error)
	Disconnect(multiAddr identity.MultiAddress) error
}

// Server is an interface for accepting connections from a Client.
type Server interface {
	Listen(ctx context.Context) (Stream, error)
}

// clientRecycler encapsulates a Client and reuses a Stream that has been
// connected to Server when multiple connections are needed for the same
// identity.MultiAddress. It does not protect the Stream from concurrent use.
type clientRecycler struct {
	client Client

	connsMu *sync.Mutex
	connsRc map[string]int64
	conns   map[string]Stream
}

func NewClientRecycler(client Client) Client {
	return &clientRecycler{
		client:  client,
		connsMu: new(sync.Mutex),
		connsRc: map[string]int64{},
		conns:   map[string]Stream{},
	}
}

func (client *clientRecycler) Connect(ctx context.Context, multiAddr identity.MultiAddress) (Stream, error) {
	client.connsMu.Lock()
	defer client.connsMu.Unlock()

	if _, ok := client.conns[multiAddr.String()]; !ok {
		conn, err := client.client.Connect(ctx, multiAddr)
		if err != nil {
			return nil, err
		}
		client.conns[multiAddr.String()] = conn
		client.connsRc[multiAddr.String()] = 0
	}
	client.connsRc[multiAddr.String()]++
	return client.conns[multiAddr.String()], nil
}

func (client *clientRecycler) Disconnect(multiAddr identity.MultiAddress) error {
	client.connsMu.Lock()
	defer client.connsMu.Unlock()

	if client.connsRc[multiAddr.String()] == 0 {
		return ErrDisconnectedClosedStream
	}

	client.connsRc[multiAddr.String()]--
	if client.connsRc[multiAddr.String()] == 0 {
		if err := client.client.Disconnect(multiAddr); err != nil {
			return err
		}
		delete(client.conns, multiAddr.String())
		delete(client.connsRc, multiAddr.String())
	}
	return nil
}
