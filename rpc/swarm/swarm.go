package swarm

import (
	"errors"
	"fmt"

	"github.com/republicprotocol/republic-go/identity"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// ErrNotFound is returned when an identity.MultiAddress cannot be found for an
// identity.Address.
var ErrNotFound = errors.New("multiaddress not found")

// Swarm implements the gRPC Swarm service.
type Swarm struct {
	client *Client
}

// NewSwarm that will use the a Client to call RPCs on other Swarm services and
// to identify itself.
func NewSwarm(client *Client) Swarm {
	return Swarm{
		client: client,
	}
}

// Register the gRPC service to a grpc.Server.
func (swarmer *Swarm) Register(server *grpc.Server) {
	RegisterSwarmServer(server, swarmer)
}

// Ping is an RPC used to notify a Swarm service about the existence of a
// client. In the PingRequest, the client sends a signed identity.MultiAddress
// that the Swarm service will add to its dht.DHT. If successfuly, the Swarm
// service will respond with an empty PingResponse.
func (swarmer *Swarm) Ping(ctx context.Context, request *PingRequest) (*PingResponse, error) {
	// FIXME: Verify the client signature
	signature := request.GetSignature()
	multiAddress, err := identity.NewMultiAddressFromString(request.GetMultiAddress())
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal multiaddress: %v", err)
	}
	if err := multiAddress.VerifySignature(signature); err != nil {
		return nil, fmt.Errorf("cannot verify multiaddress: %v", err)
	}
	if err := swarmer.client.UpdateDHT(multiAddress); err != nil {
		return nil, fmt.Errorf("cannot update dht: %v", err)
	}
	return &PingResponse{}, nil
}

// Query is an RPC used to find identity.MultiAddresses. In the QueryRequest,
// the client sends an identity.Address and the Swarm service will stream
// identity.MultiAddresses to the client. The Swarm service will stream all
// identity.MultiAddresses that are closer to the queried identity.Address than
// the Swarm service itself.
func (swarmer *Swarm) Query(request *QueryRequest, stream Swarm_QueryServer) error {
	// FIXME: Verify the client signature

	query := identity.Address(request.GetAddress())
	multiAddrs := swarmer.client.DHT().MultiAddresses()

	for _, multiAddr := range multiAddrs {
		isPeerCloser, err := identity.Closer(multiAddr.Address(), swarmer.client.Address(), query)
		if err != nil {
			return err
		}
		if isPeerCloser {
			// FIXME: Send the peer signature for this identity.MultiAddress so
			// that the client can verify it

			response := &QueryResponse{
				Signature:    []byte{},
				MultiAddress: multiAddr.String(),
			}
			if err := stream.Send(response); err != nil {
				return err
			}
		}
	}

	return nil
}
