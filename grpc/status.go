package grpc

import (
	"context"

	"github.com/republicprotocol/republic-go/swarm"
)

// FIXME: Replace StatusService with the HTTP status endpoint.
type StatusService struct {
	swarmer swarm.Swarmer
}

func NewStatusService(swarmer swarm.Swarmer) StatusService {
	return StatusService{
		swarmer: swarmer,
	}
}

func (service *StatusService) Status(ctx context.Context, request *StatusRequest) (*StatusResponse, error) {
	multiAddrs, err := service.swarmer.GetConnectedPeers()
	if err != nil {
		return nil, err
	}

	return &StatusResponse{
		Address:      string(service.swarmer.MultiAddress().Address()),
		Bootstrapped: true, // FIXME: We probably do not need this status
		Peers:        int64(len(multiAddrs)),
	}, nil
}

// Register the gRPC service to a Server.
func (service *StatusService) Register(server *Server) {
	RegisterStatusServiceServer(server.Server, service)
}
