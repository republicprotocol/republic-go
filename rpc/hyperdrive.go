package rpc

import (
	"context"
	"errors"
	"sync"

	"github.com/republicprotocol/republic-go/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HyperdriveDelegate interface {
	OnTxReceived(tx *Tx)
}

type HyperdriveService struct {
	HyperdriveDelegate

	senderSignalsMu *sync.Mutex
	senderSignals   map[string]chan (<-chan *DriveMessage)

	receiverSignalsMu *sync.Mutex
	receiverSignals   map[string]chan (<-chan *DriveMessage)

	errSignalsMu *sync.Mutex
	errSignals   map[string]chan (<-chan error)
}

func NewHyperdriveService(delegate HyperdriveDelegate) HyperdriveService {
	return HyperdriveService{
		HyperdriveDelegate: delegate,

		senderSignalsMu: new(sync.Mutex),
		senderSignals:   map[string]chan (<-chan *DriveMessage){},

		receiverSignalsMu: new(sync.Mutex),
		receiverSignals:   map[string]chan (<-chan *DriveMessage){},

		errSignalsMu: new(sync.Mutex),
		errSignals:   map[string]chan (<-chan error){},
	}
}

func (service *HyperdriveService) Register(server *grpc.Server) {
	RegisterHyperdriveServer(server, service)
}

func (service *HyperdriveService) SendTx(ctx context.Context, tx *Tx) (*Nothing, error) {
	service.OnTxReceived(tx)
	return nil, nil
}

func (service *HyperdriveService) WaitForDrive(multiAddress identity.MultiAddress, driveMessageIn <-chan *DriveMessage) (<-chan *DriveMessage, <-chan error) {
	multiAddressAsStr := multiAddress.String()

	senderSignal := service.senderSignal(multiAddressAsStr)
	defer service.closeSenderSignal(multiAddressAsStr)
	senderSignal <- driveMessageIn

	errSignal := service.errSignal(multiAddressAsStr)
	errCh := <-errSignal

	receiverSignal := service.receiverSignal(multiAddressAsStr)
	receiverCh := <-receiverSignal

	return receiverCh, errCh
}

func (service *HyperdriveService) SyncBlock(nothing *Nothing, stream Hyperdrive_SyncBlockServer) error {
	return nil
}

func (service *HyperdriveService) Drive(stream Hyperdrive_DriveServer) error {

	// Accept an initial authentication message to verify the identity of the
	// client
	authentication, err := stream.Recv()
	if err != nil {
		return err
	}
	if authentication.MultiAddress == nil {
		return errors.New("cannot connect to stream: no authentication message")
	}
	multiAddress, _, err := UnmarshalMultiAddress(authentication.MultiAddress)
	// FIXME: Validate the multiaddress signature.
	multiAddressAsStr := multiAddress.String()

	errCh := make(chan error)
	defer close(errCh)

	errSignal := service.errSignal(multiAddressAsStr)
	defer service.closeErrSignal(multiAddressAsStr)
	errSignal <- errCh

	senderErrCh := make(chan error, 1)
	go func() {
		defer close(senderErrCh)
		senderSignal := service.senderSignal(multiAddressAsStr)
		senderCh := <-senderSignal

		for message := range senderCh {
			if err := stream.Send(message); err != nil {
				s, _ := status.FromError(err)
				if s.Code() != codes.Canceled && s.Code() != codes.DeadlineExceeded {
					errCh <- err
				}
				senderErrCh <- err
				return
			}
		}
	}()

	receiverCh := make(chan *DriveMessage)
	defer close(receiverCh)

	receiverSignal := service.receiverSignal(multiAddressAsStr)
	defer service.closeReceiverSignal(multiAddressAsStr)
	receiverSignal <- receiverCh

	for {
		message, err := stream.Recv()
		if err != nil {
			errCh <- err
			return err
		}

		select {
		case err, ok := <-senderErrCh:
			if !ok {
				return nil
			}
			errCh <- err
			return err
		case <-stream.Context().Done():
			errCh <- stream.Context().Err()
			return stream.Context().Err()
		case receiverCh <- message:
		}
	}
}

func (service *HyperdriveService) senderSignal(key string) chan (<-chan *DriveMessage) {
	service.senderSignalsMu.Lock()
	defer service.senderSignalsMu.Unlock()

	if _, ok := service.senderSignals[key]; !ok {
		service.senderSignals[key] = make(chan (<-chan *DriveMessage), 1)
	}
	return service.senderSignals[key]
}

func (service *HyperdriveService) closeSenderSignal(key string) {
	service.senderSignalsMu.Lock()
	defer service.senderSignalsMu.Unlock()

	if _, ok := service.senderSignals[key]; ok {
		close(service.senderSignals[key])
		delete(service.senderSignals, key)
	}
}

func (service *HyperdriveService) receiverSignal(key string) chan (<-chan *DriveMessage) {
	service.receiverSignalsMu.Lock()
	defer service.receiverSignalsMu.Unlock()

	if _, ok := service.receiverSignals[key]; !ok {
		service.receiverSignals[key] = make(chan (<-chan *DriveMessage), 1)
	}
	return service.receiverSignals[key]
}

func (service *HyperdriveService) closeReceiverSignal(key string) {
	service.receiverSignalsMu.Lock()
	defer service.receiverSignalsMu.Unlock()

	if _, ok := service.receiverSignals[key]; ok {
		close(service.receiverSignals[key])
		delete(service.receiverSignals, key)
	}
}

func (service *HyperdriveService) errSignal(key string) chan (<-chan error) {
	service.errSignalsMu.Lock()
	defer service.errSignalsMu.Unlock()

	if _, ok := service.errSignals[key]; !ok {
		service.errSignals[key] = make(chan (<-chan error), 1)
	}
	return service.errSignals[key]
}

func (service *HyperdriveService) closeErrSignal(key string) {
	service.errSignalsMu.Lock()
	defer service.errSignalsMu.Unlock()

	if _, ok := service.errSignals[key]; ok {
		close(service.errSignals[key])
		delete(service.errSignals, key)
	}
}
