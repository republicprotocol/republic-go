package swarm

import (
	"errors"
	"fmt"
	"time"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/rpc/dht"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// ErrNotFound is returned when an identity.MultiAddress cannot be found for an
// identity.Address.
var ErrNotFound = errors.New("multiaddress not found")

// Swarm implements the gRPC Swarm service.
type Swarm struct {
	multiAddress identity.MultiAddress
	dht          *dht.DHT
}

// NewSwarm that will use the dht.DHT to store information about peers.
func NewSwarm(multiAddress identity.MultiAddress, dht *dht.DHT) Swarm {
	return Swarm{
		multiAddress: multiAddress,
		dht:          dht,
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
	// FIXME: Verify the client signature
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
func (swarm *Swarm) Query(request *QueryRequest, stream Swarm_QueryServer) error {
	// TODO: Verify the client signature

	query := identity.Address(request.GetAddress())
	multiAddrs := swarm.dht.MultiAddresses()

	for _, multiAddr := range multiAddrs {
		isPeerCloser, err := identity.Closer(multiAddr.Address(), swarm.dht.Address, query)
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

func (swarm *Swarm) QueryDeep(request *QueryRequest, depth int) (identity.MultiAddress, error) {

	whitelist := identity.MultiAddresses{}
	blacklist := map[identity.Address]struct{}{}

	// Build a list of identity.MultiAddresses that are closer to the query
	// than the Swarm service
	query := identity.Address(request.GetAddress())
	multiAddrs := swarm.dht.MultiAddresses()
	for _, multiAddr := range multiAddrs {
		// Shrort circuit if the Swarm service is directly connected to the
		// query
		if query == multiAddr.Address() {
			return multiAddr, nil
		}

		isPeerCloser, err := identity.Closer(multiAddr.Address(), swarm.dht.Address, query)
		if err != nil {
			return identity.MultiAddress{}, err
		}
		if isPeerCloser {
			whitelist = append(whitelist, multiAddr)
		}
	}

	// Search all peers for identiyt.MultiAddresses that are even closer to the
	// query until the depth limit is reach or there are no more peers left to
	// search
	for i := 0; (i < depth || depth == 0) && len(whitelist) > 0; i++ {

		peer := whitelist[0]
		whitelist = whitelist[1:]
		if _, ok := blacklist[peer.Address()]; ok {
			continue
		}
		blacklist[peer.Address()] = struct{}{}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		conn, err := swarm.connectionPool.Dial(ctx, peer)
		if err != nil {
			return identity.MultiAddress{}, fmt.Errorf("cannot dial %v: %v", peer, err)
		}
		defer conn.Close()

		client := NewSwarmClient(conn)
		// FIXME: Provide verifiable signature
		stream, err := client.Query(ctx, request)
		if err != nil {
			return identity.MultiAddress{}, err
		}
		if err := swarm.updateDHT(peer); err != nil {
			return identity.MultiAddress{}, fmt.Errorf("cannot update dht: %v", err)
		}

		for {
			message, err := stream.Recv()
			if err != nil {
				return identity.MultiAddress{}, err
			}

			// FIXME: Verify the message signature
			signature := message.GetSignature()
			multiAddress, err := identity.NewMultiAddressFromString(message.GetMultiAddress())
			if err != nil {
				continue
			}
			if err := multiAddress.VerifySignature(signature); err != nil {
				continue
			}
			if err := swarm.dht.UpdateMultiAddress(multiAddress); err != nil {
				continue
			}

			if query == multiAddress.Address() {
				return multiAddress, nil
			}
			isPeerCloser, err := identity.Closer(multiAddress.Address(), swarm.dht.Address, query)
			if err != nil {
				return identity.MultiAddress{}, err
			}
			if isPeerCloser {
				whitelist = append(whitelist, multiAddress)
			}
		}
	}

	return identity.MultiAddress{}, ErrNotFound
}

func (swarm *Swarm) updateDHT(multiAddress identity.MultiAddress) error {
	if swarm.multiAddress.Address() == multiAddress.Address() {
		return nil
	}
	if err := swarm.dht.UpdateMultiAddress(multiAddress); err != nil {
		if err == dht.ErrFullBucket {
			prune := swarm.pruneDHT(multiAddress.Address())
			if prune {
				return swarm.dht.UpdateMultiAddress(multiAddress)
			}
			return nil
		}
		return err
	}
	return nil
}

func (swarm *Swarm) pruneDHT(addr identity.Address) bool {
	bucket, err := swarm.dht.FindBucket(addr)
	if err != nil {
		return false
	}
	if bucket == nil || bucket.Length() == 0 {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Ping the first multiaddress in the bucket and see if they are still
	// responsive
	multiAddress := bucket.MultiAddresses[0]
	conn, err := swarm.connectionPool.Dial(ctx, multiAddress)
	if err != nil {
		swarm.dht.RemoveMultiAddress(multiAddress)
		return true
	}
	defer conn.Close()

	// FIXME: Provide verifiable signature of the multiaddress
	client := NewSwarmClient(conn)
	if _, err := client.Ping(ctx, &PingRequest{
		Signature:    []byte{},
		MultiAddress: swarm.multiAddress.String(),
	}); err != nil {
		swarm.dht.RemoveMultiAddress(multiAddress)
		return true
	}

	swarm.dht.UpdateMultiAddress(multiAddress)
	return false
}
