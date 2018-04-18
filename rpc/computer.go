package rpc

import (
	"errors"
	"sync"

	"github.com/republicprotocol/republic-go/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// A ComputerService implements the Computer gRPC service. It exposes methods
// for accepting gRPC client connections, and blocking until a gRPC client
// connection is received.
type ComputerService struct {
	senderSignalsMu *sync.Mutex
	senderSignals   map[string]chan (<-chan *Computation)

	receiverSignalsMu *sync.Mutex
	receiverSignals   map[string]chan (<-chan *Computation)

	errSignalsMu *sync.Mutex
	errSignals   map[string]chan (<-chan error)
}

// NewComputerService returns a new ComputerService.
func NewComputerService() *ComputerService {
	return &ComputerService{
		senderSignalsMu: new(sync.Mutex),
		senderSignals:   map[string]chan (<-chan *Computation){},

		receiverSignalsMu: new(sync.Mutex),
		receiverSignals:   map[string]chan (<-chan *Computation){},

		errSignalsMu: new(sync.Mutex),
		errSignals:   map[string]chan (<-chan error){},
	}
}

// Register the SmpcService with a gRPC server.
func (service *ComputerService) Register(server *grpc.Server) {
	RegisterComputerServer(server, service)
}

// WaitForCompute stream to be opened by a gRPC client. The input channel will
// be used to write Computation messages to the client, and the output channel
// can be used to read Computation messages from the client. The error channel
// will carry errors that occur when writing, or reading, Computation messages
// to, or from, the client.
func (service *ComputerService) WaitForCompute(multiAddress identity.MultiAddress, computationChIn <-chan *Computation) (<-chan *Computation, <-chan error) {
	multiAddressAsStr := multiAddress.String()

	senderSignal := service.senderSignal(multiAddressAsStr)
	defer service.closeSenderSignal(multiAddressAsStr)
	senderSignal <- computationChIn

	errSignal := service.errSignal(multiAddressAsStr)
	errCh := <-errSignal

	receiverSignal := service.receiverSignal(multiAddressAsStr)
	receiverCh := <-receiverSignal

	return receiverCh, errCh
}

// Compute opens a gRPC stream for streaming Computation messages to, and from,
// a client.
func (service *ComputerService) Compute(stream Computer_ComputeServer) error {

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

	receiverCh := make(chan *Computation)
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

func (service *ComputerService) senderSignal(key string) chan (<-chan *Computation) {
	service.senderSignalsMu.Lock()
	defer service.senderSignalsMu.Unlock()

	if _, ok := service.senderSignals[key]; !ok {
		service.senderSignals[key] = make(chan (<-chan *Computation), 1)
	}
	return service.senderSignals[key]
}

func (service *ComputerService) closeSenderSignal(key string) {
	service.senderSignalsMu.Lock()
	defer service.senderSignalsMu.Unlock()

	if _, ok := service.senderSignals[key]; ok {
		close(service.senderSignals[key])
		delete(service.senderSignals, key)
	}
}

func (service *ComputerService) receiverSignal(key string) chan (<-chan *Computation) {
	service.receiverSignalsMu.Lock()
	defer service.receiverSignalsMu.Unlock()

	if _, ok := service.receiverSignals[key]; !ok {
		service.receiverSignals[key] = make(chan (<-chan *Computation), 1)
	}
	return service.receiverSignals[key]
}

func (service *ComputerService) closeReceiverSignal(key string) {
	service.receiverSignalsMu.Lock()
	defer service.receiverSignalsMu.Unlock()

	if _, ok := service.receiverSignals[key]; ok {
		close(service.receiverSignals[key])
		delete(service.receiverSignals, key)
	}
}

func (service *ComputerService) errSignal(key string) chan (<-chan error) {
	service.errSignalsMu.Lock()
	defer service.errSignalsMu.Unlock()

	if _, ok := service.errSignals[key]; !ok {
		service.errSignals[key] = make(chan (<-chan error), 1)
	}
	return service.errSignals[key]
}

func (service *ComputerService) closeErrSignal(key string) {
	service.errSignalsMu.Lock()
	defer service.errSignalsMu.Unlock()

	if _, ok := service.errSignals[key]; ok {
		close(service.errSignals[key])
		delete(service.errSignals, key)
	}
}

type ComputationWithID struct {
	ID          [32]byte
	Computation *Computation
}

// func ZipComputationsWithID(done <-chan struct{}, computations <-chan *Computation, id [32]byte) <-chan ComputationWithID {
// 	computationsWithID := make([]<-chan ComputationWithID, len(computationChs))

// 	go func() {
// 		defer close(computationsWithID)
// 		for {
// 			select {
// 			case <-done:
// 				return
// 			case computation, ok := <-computations:
// 				if !ok {
// 					return
// 				}
// 				computationWithID := ComputationWithID{
// 					ID:          id,
// 					Computation: computation,
// 				}
// 				select {
// 				case <-done:
// 					return
// 				case computationsWithID <- computationWithID:
// 				}
// 			}
// 		}
// 	}()

// 	return computationsWithID
// }

// type ComputeMap struct {
// 	clientPool ClientPool

// 	inputs map[[32]byte]chan ComputationWithID
// 	outputs map[[32]byte]chan ComputationWithID
// }

// func (computeMap *ComputeMap) Setup(multiAddresses []identity.MultiAddress) {



// }

// func (computeMap *ComputeMap) Compute(done chan struct{}, multiAddresses []identity.MultiAddress) <-chan ComputationWithID {
// 	computationsWithID := make(chan ComputationWithID)

// 	go func() {
// 		defer close(computationsWithID)

// 		arrayOfComputationsWithID := []<-chan ComputationWithID{}
// 		for i := range multiAddresses {

// 			id := [32]byte{}
// 			copy(id[:], multiAddresses[i].ID())
	
// 			computations, errs := clientPool.Compute(...)
// 			arrayOfComputationsWithID = append(arrayOfComputationsWithID, ZipComputationsWithID(done, computations, id))
// 		}

// 		dispatch.Merge(computationsWithID, arrayOfComputationsWithID)
// 	}()

// 	return computationsWithID
// }

// func (computeMap) 
