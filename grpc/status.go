package grpc

import (
	"context"

	"github.com/republicprotocol/republic-go/dht"
)

// FIXME: Replace StatusService with the HTTP status endpoint.
type StatusService struct {
	dht *dht.DHT
}

func NewStatusService(dht *dht.DHT) StatusService {
	return StatusService{
		dht: dht,
	}
}

func (service *StatusService) Status(ctx context.Context, request *StatusRequest) (*StatusResponse, error) {
	return &StatusResponse{
		Address:      string(service.dht.Address),
		Bootstrapped: true, // FIXME: We probably do not need this status
		Peers:        int64(len(service.dht.MultiAddresses())),
	}, nil
}

// Register the gRPC service to a Server.
func (service *StatusService) Register(server *Server) {
	RegisterStatusServiceServer(server.Server, service)
}
