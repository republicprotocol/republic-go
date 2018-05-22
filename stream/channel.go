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

// ChannelStream implements a Stream Interface for local send and recv channels
type ChannelStream struct {
	send chan []byte
	recv chan []byte
}

// Send will send Message to the send channel
func (stream *ChannelStream) Send(message Message) error {
	data, err := message.MarshalBinary()
	if err != nil {
		return err
	}
	stream.send <- data
	return nil
}

// Recv will receive Message from the recv channel
func (stream *ChannelStream) Recv(message Message) error {
	data, ok := <-stream.recv
	if !ok {
		return ErrRecvOnClosedStream
	}
	return message.UnmarshalBinary(data)
}

// Close wil close both send and recv channels
func (stream *ChannelStream) Close() error {
	close(stream.send)
	close(stream.recv)
	return nil
}

// ChannelHub will store a map of all active streams
// between clients and servers. The map will always be
// symmetric over a client-server pair
type ChannelHub struct {
	connsMu *sync.Mutex
	conns   map[identity.Address]map[identity.Address]ChannelStream
}

// NewChannelHub will initialize and return a ChannelHub
func NewChannelHub() ChannelHub {
	return ChannelHub{
		connsMu: new(sync.Mutex),
		conns:   map[identity.Address]map[identity.Address]ChannelStream{},
	}
}

func (hub *ChannelHub) RegisterClient(clientAddr, serverAddr identity.Address) CloseStream {
	hub.connsMu.Lock()
	defer hub.connsMu.Unlock()

	hub.register(clientAddr, serverAddr)
	closeStream := hub.conns[clientAddr][serverAddr]
	return &(closeStream)
}

func (hub *ChannelHub) RegisterServer(serverAddr, clientAddr identity.Address) CloseStream {
	hub.connsMu.Lock()
	defer hub.connsMu.Unlock()

	hub.register(clientAddr, serverAddr)
	closeStream := hub.conns[serverAddr][clientAddr]
	return &closeStream
}

func (hub *ChannelHub) register(clientAddr, serverAddr identity.Address) {

	// Ensure that both mappings are initialized before using them
	if _, ok := hub.conns[clientAddr]; !ok {
		hub.conns[clientAddr] = map[identity.Address]ChannelStream{}
	}
	if _, ok := hub.conns[serverAddr]; !ok {
		hub.conns[serverAddr] = map[identity.Address]ChannelStream{}
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
		return
	}

	// Build a symmetric connection between the client and the server
	hub.conns[clientAddr][serverAddr] = ChannelStream{
		send: make(chan []byte),
		recv: make(chan []byte),
	}
	hub.conns[serverAddr][clientAddr] = ChannelStream{
		send: hub.conns[clientAddr][serverAddr].recv,
		recv: hub.conns[clientAddr][serverAddr].send,
	}
}

type channelClient struct {
	addr identity.Address
	hub  *ChannelHub
}

func NewChannelClient(addr identity.Address, hub *ChannelHub) Client {
	return &channelClient{
		addr: addr,
		hub:  hub,
	}
}

func (client *channelClient) Connect(ctx context.Context, multiAddr identity.MultiAddress) (CloseStream, error) {
	return client.hub.RegisterClient(client.addr, multiAddr.Address()), nil
}

type channelServer struct {
	addr identity.Address
	hub  *ChannelHub
}

func NewChannelServer(addr identity.Address, hub *ChannelHub) Server {
	return &channelServer{
		addr: addr,
		hub:  hub,
	}
}

func (server *channelServer) Listen(ctx context.Context, addr identity.Address) (CloseStream, error) {
	return server.hub.RegisterServer(server.addr, addr), nil
}
