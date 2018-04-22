package rpc

import (
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/orderbook"
	"google.golang.org/grpc"
)

// SyncerOptions defines the options specifically for syncer service
type SyncerOptions struct {
	MaxConnections int `json:"maxConnections"`
}

// SyncerService implements the syncer gRPC service. It creates
// MessageQueues for each gRPC stream.
type SyncerService struct {
	Options

	Logger    *logger.Logger
	Orderbook *orderbook.Orderbook
}

// NewSyncerService creates a SyncerService with the given options.
func NewSyncerService(options Options, logger *logger.Logger, orderbook *orderbook.Orderbook) SyncerService {
	return SyncerService{
		Options:   options,
		Logger:    logger,
		Orderbook: orderbook,
	}
}

// Register the SyncerService with a gRPC server.
func (service *SyncerService) Register(server *grpc.Server) {
	RegisterSyncerServer(server, service)
}

func (service *SyncerService) Sync(req *SyncRequest, stream Syncer_SyncServer) error {
	// Use a background MessageQueue to handle the connection until an error
	// is returned by the MessageQueue
	ch := make(chan error, 1)
	quit := make(chan struct{}, 1)
	go func() {
		ch <- service.sync(req, stream, quit)
	}()
	defer close(quit)

	// Select between the context finishing and the background worker
	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	case err := <-ch:
		return err
	}
}

func (service *SyncerService) sync(req *SyncRequest, stream Syncer_SyncServer, quit chan struct{}) error {

	ch := make(chan orderbook.Entry)

	// Shutdown the MessageQueue when the quit signal is received
	go func() {
		defer close(ch)
		defer service.Orderbook.Unsubscribe(ch)
		for {
			select {
			case entry := <-ch:
				if err := stream.SendMsg(entry); err != nil {
					return
				}
			case <-quit:
				return
			}
		}
	}()

	// Subscribe to the orderbook to received updates
	return service.Orderbook.Subscribe(ch)
}
