package swarm

import (
	"fmt"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/rpc/dht"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Swarm implements the gRPC Swarm service.
type Swarm struct {
	dht *dht.DHT
}

// NewSwarm that will use the dht.DHT to store information about peers.
func NewSwarm(dht *dht.DHT) Swarm {
	return Swarm{
		dht: dht,
	}
}

// Register the gRPC service to a grpc.Server.
func (swarm *Swarm) Register(server *grpc.Server) {
	RegisterSwarmServer(server, swarm)
}

// Ping is an RPC used to notify a Swarm service about the existence of a
// client. In the PingRequest, the client sends a signed identity.MultiAddress
// that the Swarm service will add to its dht.DHT. If successfuly, the Swarm
// service will respond with an empty PingResponse.
func (swarm *Swarm) Ping(ctx context.Context, request *PingRequest) (*PingResponse, error) {
	// TODO: Verify the client signature

	signature := request.GetSignature()
	multiAddress, err := identity.NewMultiAddressFromString(request.GetMultiAddress())
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal multiaddress: %v", err)
	}
	if err := multiAddress.VerifySignature(signature); err != nil {
		return nil, fmt.Errorf("cannot verify multiaddress: %v", err)
	}
	if err := swarm.dht.UpdateMultiAddress(multiAddress); err != nil {
		return nil, fmt.Errorf("cannot update dht: %v", err)
	}
	return &PingResponse{}, nil
}

// Query is an RPC used to find identity.MultiAddresses. In the QueryRequest,
// the client sends an identity.Address and the Swarm service will stream
// identity.MultiAddresses to the client. The Swarm service will stream all
// identity.MultiAddresses that are closer to the queried identity.Address than
// the Swarm service itself.
func (swarm *Swarm) Query(query *QueryRequest, stream Swarm_QueryServer) error {
	// TODO: Verify the client signature

	addr := identity.Address(query.GetAddress())
	multiAddrs := swarm.dht.MultiAddresses()

	for _, multiAddr := range multiAddrs {
		isPeerCloser, err := identity.Closer(multiAddr.Address(), swarm.dht.Address, addr)
		if err != nil {
			return err
		}
		if isPeerCloser {
			// TODO: Send the peer signature for this identity.MultiAddress so
			// that the client can verify it

			queryResponse := &QueryResponse{
				Signature:    []byte{},
				MultiAddress: multiAddr.String(),
			}
			if err := stream.Send(queryResponse); err != nil {
				return err
			}
		}
	}

	return nil
}
