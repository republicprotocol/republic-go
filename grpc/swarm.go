package grpc

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/swarm"
)

type SwarmService struct {
	server swarm.Server
}

func NewSwarmService(server swarm.Server) SwarmService {
	return SwarmService{
		server: server,
	}
}

// Register the gRPC service to a Server.
func (service *SwarmService) Register(server *Server) {
	RegisterSwarmServiceServer(server.Server, service)
}

// Ping is an RPC used to notify a Swarm service about the existence of a
// client. In the PingRequest, the client sends a signed identity.MultiAddress
// that the Swarm service will add to its dht.DHT. If successfuly, the Swarm
// service will respond with an empty PingResponse.
func (service *SwarmService) Ping(ctx context.Context, request *PingRequest) (*PingResponse, error) {
	from, err := identity.NewMultiAddressFromString(request.GetMultiAddress())
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal multiaddress: %v", err)
	}
	from.Signature = request.GetSignature()
	multiAddr, err := service.server.Ping(ctx, from)
	if err != nil {
		return &PingResponse{}, fmt.Errorf("cannot update dht: %v", err)
	}
	return &PingResponse{
		Signature:    multiAddr.Signature,
		MultiAddress: multiAddr.String(),
	}, nil
}

// Query is an RPC used to find identity.MultiAddresses. In the QueryRequest,
// the client sends an identity.Address and the Swarm service will stream
// identity.MultiAddresses to the client. The Swarm service will stream all
// identity.MultiAddresses that are closer to the queried identity.Address than
// the Swarm service itself.
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

type swarmClient struct {
	multiAddr identity.MultiAddress
	connPool  *ConnPool
}

// NewSwarmClient returns an implementation of the swarm.Client interface that
// uses gRPC and a recycled connection pool.
func NewSwarmClient(multiAddr identity.MultiAddress, connPool *ConnPool) swarm.Client {
	return &swarmClient{
		multiAddr: multiAddr,
		connPool:  connPool,
	}
}

// Ping implements the swarm.Client interface.
func (client *swarmClient) Ping(ctx context.Context, to identity.MultiAddress) (identity.MultiAddress, error) {
	conn, err := client.connPool.Dial(ctx, to)
	if err != nil {
		return identity.MultiAddress{}, fmt.Errorf("cannot dial %v: %v", to, err)
	}
	defer conn.Close()

	request := &PingRequest{
		Signature:    client.multiAddr.Signature,
		MultiAddress: client.multiAddr.String(),
	}
	response := &PingResponse{}
	if err := Backoff(ctx, func() error {
		response, err = NewSwarmServiceClient(conn.ClientConn).Ping(ctx, request)
		return err
	}); err != nil {
		return identity.MultiAddress{}, err
	}

	multiAddr, err := identity.NewMultiAddressFromString(response.GetMultiAddress())
	if err != nil {
		return identity.MultiAddress{}, fmt.Errorf("cannot parse %v: %v", response.GetMultiAddress(), err)
	}
	multiAddr.Signature = response.GetSignature()
	return multiAddr, nil
}

// Query implements the swarm.Client interface.
func (client *swarmClient) Query(ctx context.Context, to identity.MultiAddress, query identity.Address, querySignature [65]byte) (identity.MultiAddresses, error) {
	if _, err := client.Ping(ctx, to); err != nil {
		return identity.MultiAddresses{}, fmt.Errorf("cannot ping before query: %v", err)
	}

	conn, err := client.connPool.Dial(ctx, to)
	if err != nil {
		return identity.MultiAddresses{}, fmt.Errorf("cannot dial %v: %v", to, err)
	}
	defer conn.Close()

	request := &QueryRequest{
		Signature: querySignature[:],
		Address:   query.String(),
	}

	var stream SwarmService_QueryClient
	if err := Backoff(ctx, func() error {
		stream, err = NewSwarmServiceClient(conn.ClientConn).Query(ctx, request)
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
			log.Printf("cannot parse %v: %v", message.GetMultiAddress(), err)
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
