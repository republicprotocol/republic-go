package rpc

import (
	"context"
	"encoding"
	"errors"
	"sync"

	"github.com/republicprotocol/republic-go/identity"
)

var ErrDisconnectedClosedStream = errors.New("disconnected closed stream")

// StreamMessage is an interface for messages sent over a bidirectional
// connection between a client and a server.
type StreamMessage interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	// Types that can be used in an Stream must implement this pass through
	// method. It only exists to restrict Stream to types that have been
	// explicitly marked as compatible to avoid programmer error.
	IsMessage()
}

// Stream is an interface for bidirectional streaming of StreamMessages between
// nodes. It abstracts over the client/server architecture.
type Stream interface {
	Send(StreamMessage) error
	Recv(StreamMessage) error
}

// StreamClient is an interface for connecting to other computational node.
type StreamClient interface {
	Connect(ctx context.Context, multiAddr identity.MultiAddress) (Stream, error)
	Disconnect(multiAddr identity.MultiAddress) error
}

// StreamServer is an interface for accepting connections from other
// computational nodes.
type StreamServer interface {
	ConnectFrom(ctx context.Context, stream Stream) error
}

type RecycleStreamClient struct {
	client StreamClient

	connsMu *sync.Mutex
	connsRc map[string]int64
	conns   map[string]Stream
}

func (client *RecycleStreamClient) Connect(ctx context.Context, multiAddr identity.MultiAddress) (Stream, error) {
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

func (client *RecycleStreamClient) Disconnect(multiAddr identity.MultiAddress) error {
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
