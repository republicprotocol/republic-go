package grpc

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/smpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// ErrUnverifiedConnection is returned when a StreamClient does not produce a
// verifiable connection signature as its first StreamMessage to a
// StreamServer.
var ErrUnverifiedConnection = errors.New("unverified connection")

// ErrMalformedSignature is returned when a malformed signature, or address, is
// provided when  authenticating a connection.
var ErrMalformedSignature = errors.New("malformed signature")

// ErrMalformedEncryptionSecret is returned when a malformed encryption secret
// is provided when authenticating a connection.
var ErrMalformedEncryptionSecret = errors.New("malformed encryption secret")

// ErrStreamDisconnected is returned when a stream.Stream is disconnected and
// a connection cannot be re-established.
var ErrStreamDisconnected = errors.New("stream disconnected")

// ErrCannotGenerateSecret is returned when a random secret cannot be created
// for encrypting a connection.
var ErrCannotGenerateSecret = errors.New("cannot generate secret")

// ErrCannotEncryptSecret is returned when a secret cannot be encrypted for
// secure transfer.
var ErrCannotEncryptSecret = errors.New("cannot encrypt secret")

type Sender struct {
	cipher crypto.AESCipher

	streamMu *sync.Mutex
	stream   grpc.Stream
}

func NewSender(secret []byte, stream grpc.Stream) *Sender {
	return &Sender{
		cipher: crypto.NewAESCipher(secret),

		streamMu: new(sync.Mutex),
		stream:   stream,
	}
}

func (sender *Sender) Send(message smpc.Message) error {
	sender.streamMu.Lock()
	defer sender.streamMu.Unlock()

	if sender.stream == nil {
		return ErrStreamDisconnected
	}

	data, err := message.MarshalBinary()
	if err != nil {
		return err
	}
	data, err = sender.cipher.Encrypt(data)
	if err != nil {
		return err
	}

	return sender.stream.SendMsg(&StreamMessage{
		Data: data,
	})
}

type Connector struct {
	addr      identity.Address
	signer    crypto.Signer
	encrypter crypto.Encrypter
}

func NewConnector(addr identity.Address, signer crypto.Signer, encrypter crypto.Encrypter) *Connector {
	return &Connector{
		addr:      addr,
		signer:    signer,
		encrypter: encrypter,
	}
}

func (connector *Connector) Connect(ctx context.Context, networkID smpc.NetworkID, to identity.MultiAddress, receiver smpc.Receiver) (smpc.Sender, error) {
	secret, stream, err := connector.connect(ctx, networkID, to)
	if err != nil {
		return nil, err
	}

	sender := NewSender(secret, stream)
	recv := func() error {
		rawMessage, err := stream.Recv()
		if err != nil {
			return err
		}

		data, err := sender.cipher.Decrypt(rawMessage.Data)
		if err != nil {
			log.Printf("[error] received malformed encryption from %v on network %v: %v", to, networkID, err)
			return err
		}
		message := smpc.Message{}
		if err := message.UnmarshalBinary(data); err != nil {
			log.Printf("[error] received malformed message from %v on network %v: %v", to, networkID, err)
			return err
		}

		receiver.Receive(to.Address(), message)
		return nil
	}

	go dispatch.CoBegin(
		func() {
			for {
				// Backoff receive from the stream
				recvErrNum := 0
				recvErr := error(nil)
				backoffErr := BackoffMax(ctx, func() error {
					recvErr = recv()
					if recvErr != nil {
						// Check for valid termination conditions
						select {
						case <-ctx.Done():
							return nil
						default:
						}
						if recvErr == io.EOF {
							return nil
						}
						// Retry 10 times before attempting a reconnect
						recvErrNum++
						if recvErrNum < 10 {
							return recvErr
						}
						// Reconnect
						return BackoffMax(ctx, func() error {
							secret, stream, err = connector.connect(ctx, networkID, to)
							if err != nil {
								return err
							}
							sender.streamMu.Lock()
							sender.cipher = crypto.NewAESCipher(secret[:])
							sender.stream = stream
							sender.streamMu.Unlock()
							return nil
						}, 30000)
					}
					return nil
				}, 30000)

				// Backoff error indicates that the stream is dead and there is
				// no hope of reconnecting
				if backoffErr != nil {
					log.Printf("[error] cannot reconnect to %v on network %v: %v", to.Address(), networkID, err)
					return
				}

				// After reading a message successfully, we still need to check
				// termination conditions
				select {
				case <-ctx.Done():
					return
				default:
					if recvErr == io.EOF {
						return
					}
				}
			}
		},
		func() {
			<-ctx.Done()
			if err := stream.CloseSend(); err != nil {
				log.Printf("[error] cannot close stream to peer %v on network %v: %v", to, networkID, err)
			}
		})

	return sender, nil
}

func (connector *Connector) connect(ctx context.Context, networkID smpc.NetworkID, to identity.MultiAddress) ([]byte, StreamService_ConnectClient, error) {
	// Establish a connection to the identity.MultiAddress and clean the
	// connection once the context.Context is done
	log.Printf("[debug] (stream) dialing...")
	conn, err := Dial(ctx, to)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot dial %v: %v", to, err)
	}
	go func() {
		defer conn.Close()
		<-ctx.Done()
	}()

	// Open a bidirectional stream
	var stream StreamService_ConnectClient
	if err := BackoffMax(ctx, func() error {
		// On an error backoff and retry until the context.Context is done
		stream, err = NewStreamServiceClient(conn).Connect(ctx)
		if err != nil {
			return err
		}
		return nil
	}, 30000 /* 30s max backoff */); err != nil {
		return nil, nil, fmt.Errorf("cannot open stream: %v", err)
	}

	// Generate a secret
	log.Printf("[debug] (stream) authorising...")
	secret := make([]byte, 16)
	if _, err := rand.Read(secret); err != nil {
		return secret, nil, ErrCannotGenerateSecret
	}
	encryptedSecret, err := connector.encrypter.Encrypt(to.Address().String(), secret[:])
	if err != nil {
		return secret, nil, ErrCannotEncryptSecret
	}

	// Sign an authentication message so that the StreamService can verify the
	// identity.Address of the client
	signature := append([]byte(fmt.Sprintf("Republic Protocol: connect: from %v to %v on ", connector.addr, to.Address())), networkID[:]...)
	signature = crypto.Keccak256(signature)
	signature, err = connector.signer.Sign(signature)
	if err != nil {
		return secret, nil, fmt.Errorf("cannot sign stream authentication: %v", err)
	}

	// Send the authentication message
	if err := stream.Send(&StreamMessage{
		Signature: signature,
		Address:   connector.addr.String(),
		Network:   networkID[:],
		Data:      encryptedSecret,
	}); err != nil {
		return secret, nil, fmt.Errorf("cannot send stream address: %v", err)
	}

	return secret, stream, err
}

type Listener struct {
	mu        *sync.Mutex
	contexts  map[smpc.NetworkID]map[identity.Address]context.Context
	receivers map[smpc.NetworkID]map[identity.Address]smpc.Receiver
	senders   map[smpc.NetworkID]map[identity.Address]*Sender
}

func NewListener() *Listener {
	return &Listener{
		mu:        new(sync.Mutex),
		contexts:  map[smpc.NetworkID]map[identity.Address]context.Context{},
		receivers: map[smpc.NetworkID]map[identity.Address]smpc.Receiver{},
		senders:   map[smpc.NetworkID]map[identity.Address]*Sender{},
	}
}

func (lis *Listener) Listen(ctx context.Context, networkID smpc.NetworkID, to identity.MultiAddress, receiver smpc.Receiver) (smpc.Sender, error) {
	lis.mu.Lock()
	defer lis.mu.Unlock()

	if _, ok := lis.contexts[networkID]; !ok {
		lis.contexts[networkID] = map[identity.Address]context.Context{}
	}
	if _, ok := lis.receivers[networkID]; !ok {
		lis.receivers[networkID] = map[identity.Address]smpc.Receiver{}
	}
	if _, ok := lis.senders[networkID]; !ok {
		lis.senders[networkID] = map[identity.Address]*Sender{}
	}

	addr := to.Address()
	lis.contexts[networkID][addr] = ctx
	lis.senders[networkID][addr] = NewSender(nil, nil)
	lis.receivers[networkID][addr] = receiver

	return lis.senders[networkID][addr], nil
}

type ConnectorListener struct {
	*Connector
	*Listener
}

func NewConnectorListener(addr identity.Address, signer crypto.Signer, encrypter crypto.Encrypter) ConnectorListener {
	return ConnectorListener{
		Connector: NewConnector(addr, signer, encrypter),
		Listener:  NewListener(),
	}
}

// StreamerService implements the gRPC StreamService. After being registered to
// a gRPC Server it will listen for requests to the StreamService.Connect RPC
// and pass the connections to a Streamer.
type StreamerService struct {
	addr      identity.Address
	verifier  crypto.Verifier
	decrypter crypto.Decrypter
	lis       *Listener

	donesMu *sync.Mutex
	dones   map[smpc.NetworkID]map[identity.Address](chan struct{})
}

// NewStreamerService returns an implementation of the gRPC StreamService that
// connects stream.Streams from clients to a Streamer.
func NewStreamerService(addr identity.Address, verifier crypto.Verifier, decrypter crypto.Decrypter, lis *Listener) StreamerService {
	return StreamerService{
		addr:      addr,
		verifier:  verifier,
		decrypter: decrypter,
		lis:       lis,

		donesMu: new(sync.Mutex),
		dones:   map[smpc.NetworkID]map[identity.Address](chan struct{}){},
	}
}

// Register the StreamerService to a Server.
func (service *StreamerService) Register(server *Server) {
	RegisterStreamServiceServer(server.Server, service)
}

func (service *StreamerService) Connect(stream StreamService_ConnectServer) error {
	// Verify the address of this connection
	message, err := stream.Recv()
	if err != nil {
		log.Printf("[error] cannot receive authorisation message on network: %v", err)
		return err
	}
	addr, networkID, secret, err := service.verifyAuthentication(message.GetSignature(), message.GetAddress(), message.GetNetwork(), message.GetData())
	if err != nil {
		log.Printf("[error] cannot authorise stream on network: %v", err)
		return err
	}
	log.Printf("[debug] (stream) accepted connection")

	// TODO: Return a more appropriate error
	service.lis.mu.Lock()
	if _, ok := service.lis.contexts[networkID]; !ok {
		return nil
	}
	if _, ok := service.lis.receivers[networkID]; !ok {
		return nil
	}
	if _, ok := service.lis.senders[networkID]; !ok {
		return nil
	}
	if _, ok := service.lis.contexts[networkID][addr]; !ok {
		return nil
	}
	if _, ok := service.lis.receivers[networkID][addr]; !ok {
		return nil
	}
	if _, ok := service.lis.senders[networkID][addr]; !ok {
		return nil
	}
	ctx := service.lis.contexts[networkID][addr]
	receiver := service.lis.receivers[networkID][addr]
	sender := service.lis.senders[networkID][addr]
	service.lis.mu.Unlock()

	sender.streamMu.Lock()
	sender.cipher = crypto.NewAESCipher(secret[:])
	sender.stream = stream
	sender.streamMu.Unlock()

	service.donesMu.Lock()
	if _, ok := service.dones[networkID]; !ok {
		service.dones[networkID] = map[identity.Address](chan struct{}){}
	}
	if _, ok := service.dones[networkID][addr]; ok {
		close(service.dones[networkID][addr])
	}
	done := make(chan struct{})
	service.dones[networkID][addr] = done
	service.donesMu.Unlock()

	go func() {
		for {
			var rawMessage *StreamMessage
			var recvErr error

			backoffErr := BackoffMax(stream.Context(), func() error {
				// Receive a message
				rawMessage, recvErr = stream.Recv()
				if recvErr != nil {
					if recvErr == io.EOF {
						return nil
					}
					log.Printf("[error] cannot receive message from %v on network %v: %v", addr, networkID, recvErr)
					return recvErr
				}
				// Decrypt the message
				data, err := sender.cipher.Decrypt(rawMessage.Data)
				if err != nil {
					log.Printf("[error] received malformed encryption from %v on network %v: %v", addr, networkID, err)
					return err
				}
				message := smpc.Message{}
				if err := message.UnmarshalBinary(data); err != nil {
					log.Printf("[error] received malformed message from %v on network %v: %v", addr, networkID, err)
					return err
				}
				receiver.Receive(addr, message)
				return nil
			}, 30000)
			if backoffErr != nil {
				log.Printf("[error] cannot relisten to %v on network %v: %v", addr, networkID, backoffErr)
				return
			}
		}
	}()

	// Wait until the client closes the connection or the stream itself is
	// closed
	select {
	case <-done:
		log.Printf("[debug] (stream) client reconnected to an accepted connection")
		// TODO: Return better error.
		return nil
	case <-ctx.Done():
		log.Printf("[debug] (stream) server closed accepted connection")
		return nil
	case <-stream.Context().Done():
		log.Printf("[debug] (stream) client closed accepted connection")
		return nil
	}
}

func (service *StreamerService) verifyAuthentication(signature []byte, addr string, networkID []byte, encryptedSecret []byte) (identity.Address, smpc.NetworkID, [16]byte, error) {

	if signature == nil || len(signature) != 65 || networkID == nil || len(networkID) != 32 || addr == "" {
		return identity.Address(""), smpc.NetworkID{}, [16]byte{}, ErrMalformedSignature
	}
	if encryptedSecret == nil {
		return identity.Address(""), smpc.NetworkID{}, [16]byte{}, ErrMalformedEncryptionSecret
	}

	message := append([]byte(fmt.Sprintf("Republic Protocol: connect: from %v to %v on ", addr, service.addr)), networkID...)
	message = crypto.Keccak256(message)
	if err := service.verifier.Verify(message, signature); err != nil {
		return identity.Address(""), smpc.NetworkID{}, [16]byte{}, err
	}

	secret, err := service.decrypter.Decrypt(encryptedSecret)
	if err != nil {
		return identity.Address(""), smpc.NetworkID{}, [16]byte{}, err
	}
	if len(secret) != 16 {
		return identity.Address(""), smpc.NetworkID{}, [16]byte{}, ErrMalformedEncryptionSecret
	}
	secret16 := [16]byte{}
	copy(secret16[:], secret)
	networkID32 := [32]byte{}
	copy(networkID32[:], networkID)

	return identity.Address(addr), smpc.NetworkID(networkID32), secret16, nil
}
