package network

import (
	"fmt"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/network/dht"
	"github.com/republicprotocol/republic-go/network/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// A SwarmDelegate is used as a callback interface to inject behavior into the
// Swarm service.
type SwarmDelegate interface {
	// OnPing(from identity.MultiAddress)
	// OnQuery(from identity.MultiAddress)
	// OnQueryDeep(from identity.MultiAddress)
}

// SwarmService implements the gRPC Swarm service.
type SwarmService struct {
	SwarmDelegate
	Options
	Logger     *logger.Logger
	ClientPool *rpc.ClientPool
	DHT        *dht.DHT
}

func NewSwarmService(delegate SwarmDelegate, options Options, logger *logger.Logger, clientPool *rpc.ClientPool, dht *dht.DHT) *SwarmService {
	return &SwarmService{
		SwarmDelegate: delegate,
		Options:       options,
		Logger:        logger,
		ClientPool:    clientPool,
		DHT:           dht,
	}
}

// Register the gRPC service.
func (service *SwarmService) Register(server *grpc.Server) {
	rpc.RegisterSwarmServer(server, service)
}

// Bootstrap the Node into the network. The Node will connect to each bootstrap
// Node and attempt to find itself in the network. This process will ultimately
// connect it to Nodes that are close to it in XOR space.
func (service *SwarmService) Bootstrap() {
	// Add all bootstrap Nodes to the DHT.
	for _, bootstrapMultiAddress := range service.Options.BootstrapMultiAddresses {
		err := service.DHT.UpdateMultiAddress(bootstrapMultiAddress)
		if err != nil && service.Options.Debug >= DebugLow {
			service.Logger.Error(logger.TagNetwork, err.Error())
		}
	}

	if service.Options.Concurrent {
		// Concurrently search all bootstrap Nodes for itself.
		do.ForAll(service.Options.BootstrapMultiAddresses, func(i int) {
			bootstrapMultiAddress := service.Options.BootstrapMultiAddresses[i]
			service.bootstrapUsingMultiAddress(bootstrapMultiAddress)
		})
	} else {
		// Sequentially search all bootstrap Nodes for itself.
		for _, bootstrapMultiAddress := range service.Options.BootstrapMultiAddresses {
			service.bootstrapUsingMultiAddress(bootstrapMultiAddress)
		}
	}
	if service.Options.Debug >= DebugMedium {
		service.Logger.Info(logger.TagNetwork, fmt.Sprintf("%v connected to %v peers", service.Address(), len(service.DHT.MultiAddresses())))
	}
}

// Prune an identity.Address from the dht.DHT. Returns a boolean indicating
// whether or not an identity.Address was pruned.
func (service *SwarmService) Prune(target identity.Address) (bool, error) {
	bucket, err := service.DHT.FindBucket(target)
	if err != nil {
		return false, err
	}
	if bucket == nil || bucket.Length() == 0 {
		return false, nil
	}
	multiAddress := bucket.MultiAddresses[0]

	client, err := service.ClientPool.FindOrCreateClient(multiAddress)
	if err != nil {
		return true, service.DHT.RemoveMultiAddress(multiAddress)
	}
	if err := client.Ping(); err != nil {
		return true, service.DHT.RemoveMultiAddress(multiAddress)
	}
	return false, service.DHT.UpdateMultiAddress(multiAddress)
}

// Address returns the identity.Address of the Node.
func (service *SwarmService) Address() identity.Address {
	return service.Options.MultiAddress.Address()
}

// MultiAddress returns the identity.MultiAddress of the Node.
func (service *SwarmService) MultiAddress() identity.MultiAddress {
	return service.Options.MultiAddress
}

// Ping is used to test the connection to the Node and exchange
// identity.MultiAddresses. If the Node does not respond, or it responds with
// an error, then the connection should be considered unhealthy.
func (service *SwarmService) Ping(ctx context.Context, from *rpc.MultiAddress) (*rpc.MultiAddress, error) {
	wait := do.Process(func() do.Option {
		return do.Err(service.ping(from))
	})

	select {
	case val := <-wait:
		return rpc.SerializeMultiAddress(service.MultiAddress()), val.Err

	case <-ctx.Done():
		return rpc.SerializeMultiAddress(service.MultiAddress()), ctx.Err()
	}
}

func (service *SwarmService) ping(from *rpc.MultiAddress) error {
	return service.updatePeer(from)
}

// QueryPeers is used to return MultiAddresses that are closer to the given
// target Address. It will not return MultiAddresses that are further away from
// the target than the node itself, and it will only return MultiAddresses that
// are immediately connected to the service. The MultiAddresses returned are
// not guaranteed to be healthy connections and should be pinged.
func (service *SwarmService) QueryPeers(query *rpc.Query, stream rpc.Swarm_QueryPeersServer) error {
	wait := do.Process(func() do.Option {
		return do.Err(service.queryPeers(query, stream))
	})

	select {
	case val := <-wait:
		return val.Err

	case <-stream.Context().Done():
		return stream.Context().Err()
	}
}

func (service *SwarmService) queryPeers(query *rpc.Query, stream rpc.Swarm_QueryPeersServer) error {
	target := rpc.DeserializeAddress(query.Target)
	peers, err := service.DHT.FindMultiAddressNeighbors(target, service.Options.Alpha)
	if err != nil {
		return err
	}

	// Filter away peers that are further from the target than this service.
	for _, peer := range peers {
		closer, err := identity.Closer(peer.Address(), service.Address(), target)
		if err != nil {
			return err
		}
		if closer {
			if err := stream.Send(rpc.SerializeMultiAddress(peer)); err != nil {
				return err
			}
		}
	}
	return service.updatePeer(query.From)
}

// QueryPeersDeep is used to return the closest MultiAddresses that can be
// reached from this node, relative to a target Address. It will not return
// MultiAddresses that are further away from the target than the node itself.
// The MultiAddresses returned are not guaranteed to be healthy connections
// and should be pinged.
func (service *SwarmService) QueryPeersDeep(query *rpc.Query, stream rpc.Swarm_QueryPeersDeepServer) error {
	wait := do.Process(func() do.Option {
		return do.Err(service.queryPeersDeep(query, stream))
	})

	select {
	case val := <-wait:
		return val.Err

	case <-stream.Context().Done():
		return stream.Context().Err()
	}
}

func (service *SwarmService) queryPeersDeep(query *rpc.Query, stream rpc.Swarm_QueryPeersDeepServer) error {
	from, err := rpc.DeserializeMultiAddress(query.From)
	if err != nil {
		return err
	}
	target := rpc.DeserializeAddress(query.Target)
	peers, err := service.DHT.FindMultiAddressNeighbors(target, service.Options.Alpha)
	if err != nil {
		return err
	}

	// Create the frontier and a closure map.
	frontier := make(identity.MultiAddresses, 0, len(peers))
	visited := make(map[identity.Address]struct{})

	// Filter away peers that are further from the target than this Node.
	for _, peer := range peers {
		closer, err := identity.Closer(peer.Address(), service.Address(), target)
		if err != nil {
			return err
		}
		if closer {
			if err := stream.Send(rpc.SerializeMultiAddress(peer)); err != nil {
				service.Logger.Error(logger.TagNetwork, fmt.Sprintf("cannot send deep query result: %s", err.Error()))
				return err
			}
			frontier = append(frontier, peer)
		}
	}

	// Immediately close the Node that sends the query and Node is running
	// the query and mark all peers in the frontier as seen.
	visited[from.Address()] = struct{}{}
	visited[service.Address()] = struct{}{}
	for _, peer := range frontier {
		visited[peer.Address()] = struct{}{}
	}

	// While there are still Nodes to be explored in the frontier.
	for len(frontier) > 0 {
		// Pop the first peer off the frontier.
		peer := frontier[0]
		frontier = frontier[1:]

		// Close the peer and use it to find peers that are even closer to the
		// target.
		visited[peer.Address()] = struct{}{}
		if peer.Address() == target {
			continue
		}

		candidates, err := service.ClientPool.QueryPeers(peer, query.Target)
		if err != nil {
			service.Logger.Error(logger.TagNetwork, fmt.Sprintf("cannot deepen query: %s", err.Error()))
			continue
		}

		for serializedCandidate := range candidates {

			candidate, err := rpc.DeserializeMultiAddress(serializedCandidate)
			if err != nil {
				service.Logger.Error(logger.TagNetwork, fmt.Sprintf("cannot deserialize multiaddress: %s", err.Error()))
				return err
			}
			if _, ok := visited[candidate.Address()]; ok {
				continue
			}
			// Expand the frontier by candidates that have not already been
			// explored, and store them in a persistent list of close peers.
			if err := stream.Send(serializedCandidate); err != nil {
				service.Logger.Error(logger.TagNetwork, fmt.Sprintf("cannot send query result: %s", err.Error()))
				return err
			}
			frontier = append(frontier, candidate)
			visited[candidate.Address()] = struct{}{}
		}
	}
	return service.updatePeer(query.From)
}

func (service *SwarmService) bootstrapUsingMultiAddress(bootstrapMultiAddress identity.MultiAddress) error {
	var err error
	var peers chan *rpc.MultiAddress

	// Query the bootstrap service.
	peers, err = service.ClientPool.QueryPeersDeep(bootstrapMultiAddress, rpc.SerializeAddress(service.Address()))
	if err != nil {
		if service.Options.Debug >= DebugLow {
			service.Logger.Error(logger.TagNetwork, err.Error())
		}
		return err
	}
	// Peers returned by the query will be added to the DHT.
	for serializedPeer := range peers {
		peer, err := rpc.DeserializeMultiAddress(serializedPeer)
		if err != nil {
			service.Logger.Error(logger.TagNetwork, fmt.Sprintf("cannot deserialize multiaddress: %s", err.Error()))
			continue
		}
		if peer.Address() == service.Address() {
			continue
		}
		if err := service.DHT.UpdateMultiAddress(peer); err != nil {
			if service.Options.Debug >= DebugLow {
				service.Logger.Error(logger.TagNetwork, fmt.Sprintf("cannot update DHT: %s", err.Error()))
			}
		}
	}
	return nil
}

func (service *SwarmService) updatePeer(peer *rpc.MultiAddress) error {
	peerMultiAddress, err := rpc.DeserializeMultiAddress(peer)
	if err != nil {
		return err
	}
	if service.Address() == peerMultiAddress.Address() {
		return nil
	}
	if err := service.DHT.UpdateMultiAddress(peerMultiAddress); err != nil {
		if err == dht.ErrFullBucket {
			pruned, err := service.Prune(peerMultiAddress.Address())
			if err != nil {
				return err
			}
			if pruned {
				return service.DHT.UpdateMultiAddress(peerMultiAddress)
			}
			return nil
		}
		return err
	}
	return nil
}

// FindNode will try to find the node multiAddress by its republic ID.
func (service *SwarmService) FindNode(targetID identity.ID) (*identity.MultiAddress, error) {
	target := targetID.Address()
	targetMultiAddress, err := service.DHT.FindMultiAddress(target)
	if err != nil {
		return nil, err
	}
	if targetMultiAddress != nil {
		return targetMultiAddress, nil
	}

	peers, err := service.DHT.FindMultiAddressNeighbors(target, service.Options.Alpha)
	if err != nil {
		return nil, err
	}

	serializedTarget := rpc.SerializeAddress(target)
	for _, peer := range peers {
		candidates, err := service.ClientPool.QueryPeersDeep(peer, serializedTarget)
		if err != nil {
			service.Logger.Error(logger.TagNetwork, fmt.Sprintf("cannot find node: %s", err.Error()))
			continue
		}
		for candidate := range candidates {
			deserializedCandidate, err := rpc.DeserializeMultiAddress(candidate)
			if err != nil {
				service.Logger.Error(logger.TagNetwork, fmt.Sprintf("cannot deserialize multiaddress: %s", err.Error()))
				continue
			}
			if target == deserializedCandidate.Address() {
				close(candidates)
				return &deserializedCandidate, nil
			}
		}
	}

	return nil, nil
}
