package grpc

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/stream"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// ErrUnverifiedConnection is returned when a StreamClient does not produce a
// verifiable connection signature as its first StreamMessage to a
// StreamServer.
var ErrUnverifiedConnection = errors.New("unverified connection")

// ErrNilAuthentication is returned when no authentication message is provided
// for a connection.
var ErrNilAuthentication = errors.New("nil authentication")

// ErrStreamDisconnected is returned when a stream.Stream is disconnected and
// a connection cannot be re-established.
var ErrStreamDisconnected = errors.New("stream disconnected")

// concurrentStream is a grpc.Stream that is safe for concurrent reading and
// writing and implements the stream.Stream interface.
type concurrentStream struct {
	done   chan struct{}
	closed bool

	grpcSendMu *sync.Mutex
	grpcRecvMu *sync.Mutex
	grpcStream grpc.Stream
}

func newConcurrentStream(grpcStream grpc.Stream) *concurrentStream {
	return &concurrentStream{
		done:   make(chan struct{}),
		closed: false,

		grpcSendMu: new(sync.Mutex),
		grpcRecvMu: new(sync.Mutex),
		grpcStream: grpcStream,
	}
}

// Send implements the stream.Stream interface.
func (concurrentStream *concurrentStream) Send(message stream.Message) error {
	concurrentStream.grpcSendMu.Lock()
	defer concurrentStream.grpcSendMu.Unlock()

	if concurrentStream.closed {
		return stream.ErrSendOnClosedStream
	}

	data, err := message.MarshalBinary()
	if err != nil {
		return err
	}
	return concurrentStream.grpcStream.SendMsg(&StreamMessage{
		Data: data,
	})
}

// Recv implements the stream.Stream interface.
func (concurrentStream *concurrentStream) Recv(message stream.Message) error {
	concurrentStream.grpcRecvMu.Lock()
	defer concurrentStream.grpcRecvMu.Unlock()

	if concurrentStream.closed {
		return stream.ErrRecvOnClosedStream
	}

	data := StreamMessage{}
	if err := concurrentStream.grpcStream.RecvMsg(&data); err != nil {
		return err
	}
	return message.UnmarshalBinary(data.Data)
}

// Close the stream. This will close the done channel, and prevents future
// sending and receiving on this stream.
func (concurrentStream *concurrentStream) Close() error {
	concurrentStream.grpcSendMu.Lock()
	concurrentStream.grpcRecvMu.Lock()
	defer concurrentStream.grpcSendMu.Unlock()
	defer concurrentStream.grpcRecvMu.Unlock()

	if concurrentStream.closed {
		return nil
	}
	concurrentStream.closed = true
	close(concurrentStream.done)

	if grpcClientStream, ok := concurrentStream.grpcStream.(grpc.ClientStream); ok {
		return grpcClientStream.CloseSend()
	}
	return nil
}

// Done returns a read-only channel that can be used to wait for the stream to
// be closed.
func (concurrentStream *concurrentStream) Done() <-chan struct{} {
	return concurrentStream.done
}

// streamClient implements the stream.Client by using gRPC to create
// concurrentConnStreams.
type streamClient struct {
	signer crypto.Signer
	addr   identity.Address
}

func newStreamClient(signer crypto.Signer, addr identity.Address) *streamClient {
	return &streamClient{
		signer: signer,
		addr:   addr,
	}
}

// Connect implements the stream.Client interface.
func (client *streamClient) Connect(ctx context.Context, multiAddr identity.MultiAddress) (stream.Stream, error) {
	// Establish a connection to the identity.MultiAddress and clean the
	// connection once the context.Context is done
	conn, err := Dial(ctx, multiAddr)
	if err != nil {
		return nil, fmt.Errorf("cannot dial %v: %v", multiAddr, err)
	}
	go func() {
		defer conn.Close()
		<-ctx.Done()
	}()

	// Open a bidirectional stream
	var grpcStream StreamService_ConnectClient
	if err := BackoffMax(ctx, func() error {
		// On an error backoff and retry until the context.Context is done
		grpcStream, err = NewStreamServiceClient(conn).Connect(ctx)
		return err
	}, 30000 /* 30s max backoff */); err != nil {
		return nil, fmt.Errorf("cannot open stream: %v", err)
	}

	// Sign an authentication message so that the StreamService can verify the
	// identity.Address of the client
	data := []byte(fmt.Sprintf("Republic Protocol: connect: from %v to %v", client.addr, multiAddr.Address()))
	data = crypto.Keccak256(data)
	dataSignature, err := client.signer.Sign(data)
	if err != nil {
		return nil, fmt.Errorf("cannot sign stream authentication: %v", err)
	}
	// Send the authentication message
	if err := grpcStream.Send(&StreamMessage{
		Authentication: &StreamAuthentication{
			Signature: dataSignature,
			Address:   client.addr.String(),
		},
		Data: []byte{},
	}); err != nil {
		return nil, fmt.Errorf("cannot send stream address: %v", err)
	}

	return newConcurrentStream(grpcStream), err
}

// concurrentStreamConnector connects and reconnects concurrentStreams.
type concurrentStreamConnector struct {
	client *streamClient

	mu            *sync.Mutex
	streamsMu     map[identity.Address]*sync.Mutex
	streams       map[identity.Address]*concurrentStream
	streamsCtx    map[identity.Address]context.Context
	streamsCancel map[identity.Address]context.CancelFunc
	streamsRc     map[identity.Address]int
}

func newConcurrentStreamConnector(client *streamClient) *concurrentStreamConnector {
	return &concurrentStreamConnector{
		client: client,

		mu:            new(sync.Mutex),
		streamsMu:     map[identity.Address]*sync.Mutex{},
		streams:       map[identity.Address]*concurrentStream{},
		streamsCtx:    map[identity.Address]context.Context{},
		streamsCancel: map[identity.Address]context.CancelFunc{},
		streamsRc:     map[identity.Address]int{},
	}
}

func (connector *concurrentStreamConnector) connect(ctx context.Context, multiAddr identity.MultiAddress) (*concurrentStream, error) {
	addr := multiAddr.Address()

	connector.mu.Lock()
	defer connector.mu.Unlock()
	if connector.streamsMu[addr] == nil {
		connector.streamsMu[addr] = new(sync.Mutex)
	}
	connector.streamsMu[addr].Lock()
	defer connector.streamsMu[addr].Unlock()

	if connector.streamsRc[addr] == 0 {

		connector.mu.Unlock()
		ctx, cancel := context.WithCancel(context.Background())
		stream, err := connector.client.Connect(ctx, multiAddr)
		connector.mu.Lock()

		if err != nil {
			cancel()
			return nil, err
		}

		connector.streams[addr] = stream.(*concurrentStream)
		connector.streamsCtx[addr] = ctx
		connector.streamsCancel[addr] = cancel
	}

	go func() {
		<-ctx.Done()

		connector.mu.Lock()
		defer connector.mu.Unlock()

		connector.streamsRc[addr]--
		if connector.streamsRc[addr] == 0 {

			connector.streams[addr].Close()
			connector.streamsCancel[addr]()

			delete(connector.streamsMu, addr)
			delete(connector.streams, addr)
			delete(connector.streamsCtx, addr)
			delete(connector.streamsCancel, addr)
			delete(connector.streamsRc, addr)
		}
	}()

	connector.streamsRc[addr]++
	return connector.streams[addr], nil
}

func (connector *concurrentStreamConnector) listen(ctx context.Context, multiAddr identity.MultiAddress) (*concurrentStream, error) {
	addr := multiAddr.Address()
	stream := func() *concurrentStream {

		connector.mu.Lock()
		defer connector.mu.Unlock()
		if connector.streamsMu[addr] == nil {
			connector.streamsMu[addr] = new(sync.Mutex)
		}
		connector.streamsMu[addr].Lock()
		defer connector.streamsMu[addr].Unlock()

		if connector.streamsRc[addr] == 0 {
			ctx, cancel := context.WithCancel(context.Background())
			connector.streamsCtx[addr] = ctx
			connector.streamsCancel[addr] = cancel
		}

		go func() {
			<-ctx.Done()

			connector.mu.Lock()
			defer connector.mu.Unlock()

			connector.streamsRc[addr]--
			if connector.streamsRc[addr] == 0 {

				if connector.streams[addr] != nil {
					connector.streams[addr].Close()
				}
				connector.streamsCancel[addr]()

				delete(connector.streamsMu, addr)
				delete(connector.streams, addr)
				delete(connector.streamsCtx, addr)
				delete(connector.streamsCancel, addr)
				delete(connector.streamsRc, addr)
			}
		}()

		connector.streamsRc[addr]++
		return connector.streams[addr]
	}()

	if stream == nil {
		err := BackoffMax(ctx, func() error {
			if stream = connector.connection(addr); stream == nil {
				return ErrStreamDisconnected
			}
			return nil
		}, 30000 /* Maximum backoff 30s */)
		return stream, err
	}
	return stream, nil
}

func (connector *concurrentStreamConnector) reconnect(multiAddr identity.MultiAddress) (*concurrentStream, error) {
	addr := multiAddr.Address()

	connector.mu.Lock()
	defer connector.mu.Unlock()
	if connector.streamsMu[addr] == nil {
		connector.streamsMu[addr] = new(sync.Mutex)
	}
	connector.streamsMu[addr].Lock()
	defer connector.streamsMu[addr].Unlock()

	if connector.streams[addr] != nil {
		connector.streams[addr].Close()
		ctx := connector.streamsCtx[addr]
		connector.mu.Unlock()
		stream, err := connector.client.Connect(ctx, multiAddr)
		connector.mu.Lock()
		if err != nil {
			return nil, err
		}
		connector.streams[addr] = stream.(*concurrentStream)
	}

	return connector.streams[addr], nil
}

func (connector *concurrentStreamConnector) accept(addr identity.Address, stream *concurrentStream) {
	connector.mu.Lock()
	defer connector.mu.Unlock()
	if connector.streamsMu[addr] == nil {
		connector.streamsMu[addr] = new(sync.Mutex)
	}
	connector.streamsMu[addr].Lock()
	defer connector.streamsMu[addr].Unlock()

	if connector.streams[addr] != nil {
		connector.streams[addr].Close()
	}
	connector.streams[addr] = stream
}

func (connector *concurrentStreamConnector) connection(addr identity.Address) *concurrentStream {
	connector.mu.Lock()
	defer connector.mu.Unlock()
	if connector.streamsMu[addr] == nil {
		connector.streamsMu[addr] = new(sync.Mutex)
	}
	connector.streamsMu[addr].Lock()
	defer connector.streamsMu[addr].Unlock()

	return connector.streams[addr]
}

// Streamer implements the stream.Streamer interface by using gRPC to create
// stream.Streams. Internally, it uses the concurrentStreamConnector to keep
// stream.Streams alive until the opening context.Context is done. If a
// bidirectional stream closes prematurely, the streamer will attempt to
// reconnect without disruption.
type Streamer struct {
	addr      identity.Address
	connector *concurrentStreamConnector
}

// NewStreamer returns an implementation of stream.Streamer that uses gRPC to
// create bidirectional streams and keeps its stream.Streams alive until the
// opening context.Context is done. The user does not need to explicitly
// attempt to reconnect in the event of a fault.
func NewStreamer(signer crypto.Signer, addr identity.Address) *Streamer {
	return &Streamer{
		addr:      addr,
		connector: newConcurrentStreamConnector(newStreamClient(signer, addr)),
	}
}

// Open implements the stream.Streamer interface. The stream.Stream returned is
// valid when an error is returned and can be used to send and receive
// stream.Messages.
func (streamer *Streamer) Open(ctx context.Context, multiAddr identity.MultiAddress) (stream.Stream, error) {
	var err error
	if streamer.addr < multiAddr.Address() {
		_, err = streamer.connector.connect(ctx, multiAddr)
	} else {
		_, err = streamer.connector.listen(ctx, multiAddr)
	}
	return newCtxStreamer(ctx, multiAddr, streamer), err
}

type ctxStreamer struct {
	*Streamer

	ctx             context.Context
	remoteMultiAddr identity.MultiAddress
}

func newCtxStreamer(ctx context.Context, remoteMultiAddr identity.MultiAddress, streamer *Streamer) *ctxStreamer {
	return &ctxStreamer{
		Streamer: streamer,

		ctx:             ctx,
		remoteMultiAddr: remoteMultiAddr,
	}
}

// Send implements the stream.Stream interface.
func (streamer *ctxStreamer) Send(message stream.Message) error {
	return streamer.message(message, func(stream *concurrentStream, message stream.Message) error { return stream.Send(message) })
}

// Recv implements the stream.Stream interface.
func (streamer *ctxStreamer) Recv(message stream.Message) error {
	return streamer.message(message, func(stream *concurrentStream, message stream.Message) error { return stream.Recv(message) })
}

func (streamer *ctxStreamer) message(message stream.Message, f func(stream *concurrentStream, message stream.Message) error) error {
	var stream *concurrentStream
	var err error

	remoteAddr := streamer.remoteMultiAddr.Address()
	if stream = streamer.connector.connection(remoteAddr); stream == nil {
		if streamer.addr < remoteAddr {
			stream, err = streamer.connector.connect(streamer.ctx, streamer.remoteMultiAddr)
		} else {
			stream, err = streamer.connector.listen(streamer.ctx, streamer.remoteMultiAddr)
		}
		if err != nil {
			return err
		}
		if stream == nil {
			return ErrStreamDisconnected
		}
	}

	if err := f(stream, message); err != nil {
		if streamer.addr < remoteAddr {
			stream, err = streamer.connector.reconnect(streamer.remoteMultiAddr)
		} else {
			// There is no such this as "relistening" so if the connection dies
			// all we can do is hope that the client will eventually attempt to
			// reconnect and until then, we simply let the errors happen
			log.Println("cannot relisten")
		}
		if err != nil {
			return err
		}
		if stream == nil {
			return ErrStreamDisconnected
		}
		return f(stream, message)
	}
	return nil
}

// StreamerService implements the gRPC StreamService. After being registered to
// a gRPC Server it will listen for requests to the StreamService.Connect RPC
// and pass the connections to a Streamer.
type StreamerService struct {
	verifier  crypto.Verifier
	addr      identity.Address
	connector *concurrentStreamConnector
}

// NewStreamerService returns an implementation of the gRPC StreamService that
// connects stream.Streams from clients to a Streamer.
func NewStreamerService(verifier crypto.Verifier, streamer *Streamer) StreamerService {
	return StreamerService{
		verifier:  verifier,
		addr:      streamer.addr,
		connector: streamer.connector,
	}
}

// Register the StreamerService to a Server.
func (service *StreamerService) Register(server *Server) {
	RegisterStreamServiceServer(server.Server, service)
}

func (service *StreamerService) Connect(grpcStream StreamService_ConnectServer) error {
	// Verify the address of this connection
	message, err := grpcStream.Recv()
	if err != nil {
		return err
	}
	addr, err := service.verifyAuthentication(message.GetAuthentication())
	if err != nil {
		return err
	}

	// Establish a connection with the recycler so that the stream can be used
	// outside of this service
	concurrentStream := newConcurrentStream(grpcStream)
	service.connector.accept(addr, concurrentStream)

	// Wait until the client closes the connection or the stream itself is
	// closed
	select {
	case <-grpcStream.Context().Done():
		return grpcStream.Context().Err()
	case <-concurrentStream.Done():
		return nil
	}
}

func (service *StreamerService) verifyAuthentication(auth *StreamAuthentication) (identity.Address, error) {
	if auth == nil || auth.Address == "" || auth.Signature == nil || len(auth.Signature) != 65 {
		return identity.Address(""), ErrNilAuthentication
	}

	data := []byte(fmt.Sprintf("Republic Protocol: connect: from %v to %v", auth.Address, service.addr))
	data = crypto.Keccak256(data)

	return identity.Address(auth.Address), service.verifier.Verify(data, auth.Signature)
}
