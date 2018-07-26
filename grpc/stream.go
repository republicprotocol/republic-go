package grpc

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

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

func (sender *Sender) inject(secret []byte, stream grpc.Stream) {
	sender.streamMu.Lock()
	defer sender.streamMu.Unlock()

	sender.cipher = crypto.NewAESCipher(secret)
	sender.stream = stream
}

func (sender *Sender) release() {
	sender.streamMu.Lock()
	defer sender.streamMu.Unlock()

	if sender.stream == nil {
		return
	}
	if stream, ok := sender.stream.(grpc.ClientStream); ok {
		if err := stream.CloseSend(); err != nil {
			log.Printf("[error] cannot release stream: %v", err)
		}
	}
	sender.stream = nil
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

	// This function is used to read a message from the sender defined above
	addr := to.Address()
	recv := func() error {
		// Block until a message is received or an error occurs
		rawMessage, err := stream.Recv()
		if err != nil {
			log.Printf("[error] cannot receive message from %v on network %v: %v", addr, networkID, err)
			return err
		}
		// Decrypt the message
		data, err := sender.cipher.Decrypt(rawMessage.Data)
		if err != nil {
			log.Printf("[error] received malformed encryption from %v on network %v: %v", addr, networkID, err)
			return err
		}
		// Unmarshal the message
		message := smpc.Message{}
		if err := message.UnmarshalBinary(data); err != nil {
			log.Printf("[error] received malformed message from %v on network %v: %v", addr, networkID, err)
			return err
		}
		// Notify the receiver of the message
		receiver.Receive(addr, message)
		return nil
	}

	go dispatch.CoBegin(
		func() {
			for {
				// Backoff receiving / reconnecting using the stream
				backoffErr := BackoffMax(ctx, func() error {
					// Check the context for termination conditions
					select {
					case <-ctx.Done():
						return nil
					default:
					}

					if err := recv(); err != nil {
						// Reconnect when an error occurs
						secret, stream, err = connector.connect(ctx, networkID, to)
						if err != nil {
							return err
						}
						sender.inject(secret, stream)
						time.Sleep(time.Second)

						// The reconnection is not considered successful until
						// we have successfully received a message
						return recv()
					}
					return nil
				}, 30000)

				// Check the context for termination conditions
				select {
				case <-ctx.Done():
					return
				default:
				}

				// Backoff error indicates that the stream is dead and there is
				// no hope of reconnecting
				if backoffErr != nil {
					log.Printf("[error] cannot reconnect to %v on network %v: %v", to.Address(), networkID, backoffErr)
					return
				}
			}
		},
		func() {
			defer sender.release()
			<-ctx.Done()
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

func (lis *Listener) Listen(ctx context.Context, networkID smpc.NetworkID, to identity.Address, receiver smpc.Receiver) (smpc.Sender, error) {
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

	lis.contexts[networkID][to] = ctx
	lis.senders[networkID][to] = NewSender(nil, nil)
	lis.receivers[networkID][to] = receiver

	return lis.senders[networkID][to], nil
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
	ctx, receiver, sender := func() (context.Context, smpc.Receiver, *Sender) {
		service.lis.mu.Lock()
		defer service.lis.mu.Unlock()
		if _, ok := service.lis.contexts[networkID]; !ok {
			return nil, nil, nil
		}
		if _, ok := service.lis.receivers[networkID]; !ok {
			return nil, nil, nil
		}
		if _, ok := service.lis.senders[networkID]; !ok {
			return nil, nil, nil
		}
		if _, ok := service.lis.contexts[networkID][addr]; !ok {
			return nil, nil, nil
		}
		if _, ok := service.lis.receivers[networkID][addr]; !ok {
			return nil, nil, nil
		}
		if _, ok := service.lis.senders[networkID][addr]; !ok {
			return nil, nil, nil
		}
		return service.lis.contexts[networkID][addr], service.lis.receivers[networkID][addr], service.lis.senders[networkID][addr]
	}()
	if ctx == nil || receiver == nil || sender == nil {
		// TODO: Return a more appropriate error
		return fmt.Errorf("not ready to accept connection")
	}
	sender.inject(secret[:], stream)
	log.Printf("[debug] (stream) accepted connection from %v", addr)

	done := func() chan struct{} {
		service.donesMu.Lock()
		defer service.donesMu.Unlock()
		if _, ok := service.dones[networkID]; !ok {
			service.dones[networkID] = map[identity.Address](chan struct{}){}
		}
		if _, ok := service.dones[networkID][addr]; ok {
			close(service.dones[networkID][addr])
		}
		service.dones[networkID][addr] = make(chan struct{})
		return service.dones[networkID][addr]
	}()

	go func() {
		for {
			var rawMessage *StreamMessage
			var recvErr error

			backoffErr := BackoffMax(stream.Context(), func() error {
				// Receive a message
				rawMessage, recvErr = stream.Recv()
				if recvErr != nil {
					// Check for termination conditions
					select {
					case <-done:
						return nil
					case <-ctx.Done():
						return nil
					case <-stream.Context().Done():
						return nil
					default:
						if recvErr == io.EOF {
							return nil
						}
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

			// Check for termination conditions
			select {
			case <-done:
				return
			case <-ctx.Done():
				return
			case <-stream.Context().Done():
				return
			default:
			}

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
		// TODO: Return better error.
		log.Printf("[debug] (stream) client reconnected to an accepted connection")
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
