package rpc

import (
	"context"
	"fmt"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/dht"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"google.golang.org/grpc"
)

// DebugLevel defines the debug level.
type DebugLevel int

// Constants for different debug options.
const (
	DebugOff = DebugLevel(iota)
	DebugLow
	DebugMedium
	DebugHigh
)

// SwarmService implements the gRPC Swarm service.
type SwarmService struct {
	Options

	ClientPool *ClientPool
	DHT        *dht.DHT
	Logger     *logger.Logger
}

// NewSwarmService returns a SwarmService.
func NewSwarmService(options Options, clientPool *ClientPool, dht *dht.DHT, logger *logger.Logger) SwarmService {
	return SwarmService{
		Options:    options,
		ClientPool: clientPool,
		DHT:        dht,
		Logger:     logger,
	}
}

// Register the gRPC service.
func (service *SwarmService) Register(server *grpc.Server) {
	RegisterSwarmServer(server, service)
}

// Bootstrap the Node into the network. The Node will connect to each bootstrap
// Node and attempt to find itself in the network. This process will ultimately
// connect it to Nodes that are close to it in XOR space.
func (service *SwarmService) Bootstrap() {
	// Add all bootstrap Nodes to the DHT.
	for _, bootstrapMultiAddress := range service.Options.BootstrapMultiAddresses {
		if err := service.DHT.UpdateMultiAddress(bootstrapMultiAddress); err != nil {
			service.Logger.Error(err.Error())
		}
	}
	if service.Options.Concurrent {
		// Concurrently search all bootstrap Nodes for itself.
		do.ForAll(service.Options.BootstrapMultiAddresses, func(i int) {
			bootstrapMultiAddress := service.Options.BootstrapMultiAddresses[i]
			if err := service.bootstrapUsingMultiAddress(bootstrapMultiAddress); err != nil {
				service.Logger.Error(fmt.Sprintf("error bootstrapping with %s: %s", bootstrapMultiAddress.Address(), err.Error()))
			}
		})
	} else {
		// Sequentially search all bootstrap Nodes for itself.
		for _, bootstrapMultiAddress := range service.Options.BootstrapMultiAddresses {
			if err := service.bootstrapUsingMultiAddress(bootstrapMultiAddress); err != nil {
				service.Logger.Error(fmt.Sprintf("error bootstrapping with %s: %s", bootstrapMultiAddress.Address(), err.Error()))
			}
		}
	}
	service.Logger.Info(fmt.Sprintf("boostrap connected to %v peers", len(service.DHT.MultiAddresses())))
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

	ctx, cancel := context.WithTimeout(context.Background(), service.Options.Timeout)
	defer cancel()

	if err := client.Ping(ctx); err != nil {
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
func (service *SwarmService) Ping(ctx context.Context, from *MultiAddress) (*MultiAddress, error) {
	wait := do.Process(func() do.Option {
		return do.Err(service.ping(from))
	})

	select {
	case val := <-wait:
		multi := service.MultiAddress()
		return MarshalMultiAddress(&multi), val.Err

	case <-ctx.Done():
		multi := service.MultiAddress()
		return MarshalMultiAddress(&multi), ctx.Err()
	}
}

func (service *SwarmService) ping(from *MultiAddress) error {
	return service.updatePeer(from)
}

// QueryPeers is used to return MultiAddresses that are closer to the given
// target Address. It will not return MultiAddresses that are further away from
// the target than the node itself, and it will only return MultiAddresses that
// are immediately connected to the service. The MultiAddresses returned are
// not guaranteed to be healthy connections and should be pinged.
func (service *SwarmService) QueryPeers(query *Query, stream Swarm_QueryPeersServer) error {
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

func (service *SwarmService) queryPeers(query *Query, stream Swarm_QueryPeersServer) error {
	target := UnmarshalAddress(query.Target)
	peers := service.DHT.MultiAddresses()

	// Filter away peers that are further from the target than this service.
	for _, peer := range peers {
		closer, err := identity.Closer(peer.Address(), service.Address(), target)
		if err != nil {
			return err
		}
		if closer {
			if err := stream.Send(MarshalMultiAddress(&peer)); err != nil {
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
func (service *SwarmService) QueryPeersDeep(query *Query, stream Swarm_QueryPeersDeepServer) error {
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

func (service *SwarmService) queryPeersDeep(query *Query, stream Swarm_QueryPeersDeepServer) error {
	from, _, err := UnmarshalMultiAddress(query.From)
	if err != nil {
		return err
	}
	target := UnmarshalAddress(query.Target)
	peers := service.DHT.MultiAddresses()

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
			if err := stream.Send(MarshalMultiAddress(&peer)); err != nil {
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

		ctx, cancel := context.WithTimeout(context.Background(), service.Timeout)
		defer cancel()

		candidates, errs := service.ClientPool.QueryPeers(ctx, peer, query.Target)
		continuing := true
		for continuing {
			select {
			case err := <-errs:
				if err != nil {
					service.Logger.Error(fmt.Sprintf("cannot deepen query: %s", err.Error()))
				}
				continuing = false
			case marshaledCandidate, ok := <-candidates:
				if !ok {
					continuing = false
					break
				}

				candidate, _, err := UnmarshalMultiAddress(marshaledCandidate)
				if err != nil {
					return err
				}
				if _, ok := visited[candidate.Address()]; ok {
					continue
				}
				// Expand the frontier by candidates that have not already been
				// explored, and store them in a persistent list of close peers.
				if err := stream.Send(marshaledCandidate); err != nil {
					return err
				}
				frontier = append(frontier, candidate)
				visited[candidate.Address()] = struct{}{}
			}

		}
	}
	return service.updatePeer(query.From)
}

func (service *SwarmService) bootstrapUsingMultiAddress(bootstrapMultiAddress identity.MultiAddress) error {
	ctx, cancel := context.WithTimeout(context.Background(), service.Timeout)
	defer cancel()

	// Query the bootstrap service.
	peers, errs := service.ClientPool.QueryPeersDeep(ctx, bootstrapMultiAddress, MarshalAddress(service.Address()))

	continuing := true
	numberOfPeers := 0
	for continuing {
		select {
		case err := <-errs:
			if err != nil {
				service.Logger.Error(fmt.Sprintf("cannot deepen query: %s", err.Error()))
			}
			continuing = false
		case marshaledPeer, ok := <-peers:
			if !ok {
				continuing = false
				break
			}

			peer, _, err := UnmarshalMultiAddress(marshaledPeer)
			if err != nil {
				service.Logger.Error(fmt.Sprintf("cannot deserialize multiaddress: %s", err.Error()))
				continue
			}
			if peer.Address() == service.Address() {
				continue
			}
			numberOfPeers++
			if err := service.DHT.UpdateMultiAddress(peer); err != nil {
				service.Logger.Error(fmt.Sprintf("cannot update DHT: %s", err.Error()))
			}
		}
	}

	service.Logger.Info(fmt.Sprintf("bootstrapping with %s returned %d peers", bootstrapMultiAddress.Address(), numberOfPeers))
	return nil
}

func (service *SwarmService) updatePeer(peer *MultiAddress) error {
	peerMultiAddress, sig, err := UnmarshalMultiAddress(peer)
	if err != nil {
		return err
	}
	err = peerMultiAddress.VerifySignature(sig)
	if err != nil {
		return nil
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

	marshaledTarget := MarshalAddress(target)
	for _, peer := range peers {
		ctx, cancel := context.WithTimeout(context.Background(), service.Timeout)
		defer cancel()

		candidates, errs := service.ClientPool.QueryPeersDeep(ctx, peer, marshaledTarget)
		for err := range errs {
			service.Logger.Error(fmt.Sprintf("error finding node: %s", err.Error()))
			continue
		}
		for candidate := range candidates {
			unmarshalCandidate, _, err := UnmarshalMultiAddress(candidate)
			if err != nil {
				service.Logger.Error(fmt.Sprintf("cannot deserialize multiaddress: %s", err.Error()))
				continue
			}
			if target == unmarshalCandidate.Address() {
				return &unmarshalCandidate, nil
			}
		}
	}
	return nil, nil
}
