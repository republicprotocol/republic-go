package grpc

import (
	"context"
	"fmt"
	"io"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/swarm"
)

type swarmClient struct {
	multiAddr identity.MultiAddress
}

// NewSwarmClient returns an implementation of the swarm.Client interface that
// uses gRPC and a recycled connection pool.
func NewSwarmClient(multiAddr identity.MultiAddress) swarm.Client {
	return &swarmClient{
		multiAddr: multiAddr,
	}
}

// Ping implements the swarm.Client interface.
func (client *swarmClient) Ping(ctx context.Context, to identity.MultiAddress) (identity.MultiAddress, error) {
	conn, err := Dial(ctx, to)
	if err != nil {
		logger.Network(logger.LevelError, fmt.Sprintf("cannot dial %v: %v", to, err))
		return identity.MultiAddress{}, fmt.Errorf("cannot dial %v: %v", to, err)
	}
	defer conn.Close()

	request := &PingRequest{
		Signature:    client.multiAddr.Signature,
		MultiAddress: client.multiAddr.String(),
	}
	response := &PingResponse{}
	if err := Backoff(ctx, func() error {
		response, err = NewSwarmServiceClient(conn).Ping(ctx, request)
		return err
	}); err != nil {
		return identity.MultiAddress{}, err
	}

	multiAddr, err := identity.NewMultiAddressFromString(response.GetMultiAddress())
	if err != nil {
		logger.Network(logger.LevelError, fmt.Sprintf("cannot parse %v: %v", response.GetMultiAddress(), err))
		return identity.MultiAddress{}, fmt.Errorf("cannot parse %v: %v", response.GetMultiAddress(), err)
	}
	multiAddr.Signature = response.GetSignature()
	return multiAddr, nil
}

// Query implements the swarm.Client interface.
func (client *swarmClient) Query(ctx context.Context, to identity.MultiAddress, query identity.Address, querySignature [65]byte) (identity.MultiAddresses, error) {

	if _, err := client.Ping(ctx, to); err != nil {
		logger.Network(logger.LevelError, fmt.Sprintf("cannot ping before query: %v", err))
		return identity.MultiAddresses{}, fmt.Errorf("cannot ping before query: %v", err)
	}

	conn, err := Dial(ctx, to)
	if err != nil {
		logger.Network(logger.LevelError, fmt.Sprintf("cannot dial %v: %v", to, err))
		return identity.MultiAddresses{}, fmt.Errorf("cannot dial %v: %v", to, err)
	}
	defer conn.Close()

	request := &QueryRequest{
		Signature: querySignature[:],
		Address:   query.String(),
	}

	var stream SwarmService_QueryClient
	if err := Backoff(ctx, func() error {
		stream, err = NewSwarmServiceClient(conn).Query(ctx, request)
		return err
	}); err != nil {
		return identity.MultiAddresses{}, err
	}

	multiAddrs := identity.MultiAddresses{}
	for {
		message, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return multiAddrs, nil
			}
			return multiAddrs, err
		}
		multiAddr, err := identity.NewMultiAddressFromString(message.GetMultiAddress())
		if err != nil {
			logger.Network(logger.LevelWarn, fmt.Sprintf("cannot parse %v: %v", message.GetMultiAddress(), err))
			continue
		}
		multiAddr.Signature = message.GetSignature()
		multiAddrs = append(multiAddrs, multiAddr)
	}
}

// MultiAddress implements the swarm.Client interface.
func (client *swarmClient) MultiAddress() identity.MultiAddress {
	return client.multiAddr
}

// SwarmService is a Service that implements the gRPC SwarmService defined in
// protobuf. It delegates responsibility for handling the Ping and Query RPCs
// to a swarm.Server.
type SwarmService struct {
	server swarm.Server
}

// NewSwarmService returns a SwarmService that uses the swarm.Server as a
// delegate.
func NewSwarmService(server swarm.Server) SwarmService {
	return SwarmService{
		server: server,
	}
}

// Register implements the Service interface.
func (service *SwarmService) Register(server *Server) {
	RegisterSwarmServiceServer(server.Server, service)
}

// Ping is an RPC used to notify a SwarmService about the existence of a
// client. In the PingRequest, the client sends a signed identity.MultiAddress
// and the SwarmService delegates the responsibility of handling this signed
// identity.MultiAddress to its swarm.Server. If its swarm.Server accepts the
// signed identity.MultiAddress of the client it will return its own signed
// identity.MultiAddress in a PingResponse.
func (service *SwarmService) Ping(ctx context.Context, request *PingRequest) (*PingResponse, error) {
	from, err := identity.NewMultiAddressFromString(request.GetMultiAddress())
	if err != nil {
		logger.Network(logger.LevelError, fmt.Sprintf("cannot unmarshal multiaddress: %v", err))
		return nil, fmt.Errorf("cannot unmarshal multiaddress: %v", err)
	}
	from.Signature = request.GetSignature()
	multiAddr, err := service.server.Ping(ctx, from)
	if err != nil {
		logger.Network(logger.LevelInfo, fmt.Sprintf("cannot update dht with %v: %v", multiAddr, err))
		return &PingResponse{}, fmt.Errorf("cannot update dht: %v", err)
	}
	return &PingResponse{
		Signature:    multiAddr.Signature,
		MultiAddress: multiAddr.String(),
	}, nil
}

// Query is an RPC used to find identity.MultiAddresses. In the QueryRequest,
// the client sends an identity.Address and the SwarmService will stream
// identity.MultiAddresses to the client. The SwarmService delegates
// responsibility to its swarm.Server to return identity.MultiAddresses that
// are close to the queried identity.Address.
func (service *SwarmService) Query(request *QueryRequest, stream SwarmService_QueryServer) error {
	query := identity.Address(request.GetAddress())
	querySig := [65]byte{}
	copy(querySig[:], request.GetSignature())

	multiAddrs, err := service.server.Query(stream.Context(), query, querySig)
	if err != nil {
		return err
	}

	for _, multiAddr := range multiAddrs {
		response := &QueryResponse{
			Signature:    multiAddr.Signature,
			MultiAddress: multiAddr.String(),
		}
		if err := stream.Send(response); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}
