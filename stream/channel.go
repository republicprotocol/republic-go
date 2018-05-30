package stream

import (
	"context"
	"errors"
	"sync"

	"github.com/republicprotocol/republic-go/identity"
)

// ErrAlreadyRegistered is returned when a client tries to register
// again with the same server
var ErrAlreadyRegistered = errors.New("client has already been registered with the server")

// channelStream implements a Stream interface using channels. It stores one
// channel for sending Messages, and another channel for receiving Messages. A
// channelStream must not be used unless it was returned from a call to
// channelStreamer.Open.
type channelStream struct {
	send chan []byte
	recv chan []byte
}

// Send implements the Stream interface by marshaling the Message to binary and
// writing it to the sending channel.
func (stream channelStream) Send(message Message) error {
	data, err := message.MarshalBinary()
	if err != nil {
		return err
	}
	stream.send <- data
	return nil
}

// Recv implements the Stream interface by reading from the receiving channel
// and unmarshaling the data into a Message.
func (stream channelStream) Recv(message Message) error {
	data, ok := <-stream.recv
	if !ok {
		return ErrRecvOnClosedStream
	}
	return message.UnmarshalBinary(data)
}

// A ChannelHub will store a map of all active channelStreams between
// identity.Addresses and ensures that the mapping is symmetrical.
type ChannelHub struct {
	connsMu *sync.Mutex
	conns   map[identity.Address]map[identity.Address]channelStream
}

// NewChannelHub returns a ChannelHub with no connections in the map.
func NewChannelHub() ChannelHub {
	return ChannelHub{
		connsMu: new(sync.Mutex),
		conns:   map[identity.Address]map[identity.Address]channelStream{},
	}
}

func (hub *ChannelHub) register(clientAddr, serverAddr identity.Address) Stream {
	hub.connsMu.Lock()
	defer hub.connsMu.Unlock()

	// Ensure that both mappings are initialized before using them
	if _, ok := hub.conns[clientAddr]; !ok {
		hub.conns[clientAddr] = map[identity.Address]channelStream{}
	}
	if _, ok := hub.conns[serverAddr]; !ok {
		hub.conns[serverAddr] = map[identity.Address]channelStream{}
	}

	// An assymetric connection should be unreachable so we explicitly check
	// and panic for clarity and easier debugging
	_, clientOk := hub.conns[clientAddr][serverAddr]
	_, serverOk := hub.conns[serverAddr][clientAddr]
	if clientOk && !serverOk {
		panic("assymetric connection from client to server")
	}
	if serverOk && !clientOk {
		panic("assymetric connection from server to client")
	}

	// A symmetric connection has already been established
	if clientOk && serverOk {
		return hub.conns[clientAddr][serverAddr]
	}

	// Build a symmetric connection between the client and the server
	hub.conns[clientAddr][serverAddr] = channelStream{
		send: make(chan []byte),
		recv: make(chan []byte),
	}
	hub.conns[serverAddr][clientAddr] = channelStream{
		send: hub.conns[clientAddr][serverAddr].recv,
		recv: hub.conns[clientAddr][serverAddr].send,
	}
	return hub.conns[clientAddr][serverAddr]
}

type channelStreamer struct {
	addr identity.Address
	hub  *ChannelHub
}

// NewChannelStreamer returns a Streamer that uses channel to implement the
// Stream interface. Streams are recycled whenever multiple connections between
// two identity.Addresses is needed.
func NewChannelStreamer(addr identity.Address, hub *ChannelHub) Streamer {
	streamer := channelStreamer{
		addr: addr,
		hub:  hub,
	}
	return NewStreamRecycler(&streamer)
}

// Open implements the Streamer interface by using the ChannelHub to register
// connections between two identity.Addresses.
func (streamer *channelStreamer) Open(ctx context.Context, multiAddr identity.MultiAddress) (Stream, error) {
	return streamer.hub.register(streamer.addr, multiAddr.Address()), nil
}
