package network

import (
	"log"
	"time"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/network/dht"
	"github.com/republicprotocol/republic-go/network/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// A SwarmDelegate is used as a callback interface to inject behavior into the
// Swarm service.
type SwarmDelegate interface {
	OnPing(from identity.MultiAddress)
	OnQuery(from identity.MultiAddress)
	OnQueryDeep(from identity.MultiAddress)
}

// SwarmService implements the gRPC Swarm service.
type SwarmService struct {
	SwarmDelegate
	Options
	DHT *dht.DHT
}

func NewSwarmService(delegate SwarmDelegate, options Options) *SwarmService {
	return &Swarm{
		Delegate: delegate,
		Options:  options,
		DHT:      dht.NewDHT(options.MultiAddress.Address(), options.MaxBucketLength),
	}
}

// Register the SwarmService as a gRPC service.
func (service *SwarmService) Register(server *grpc.Server) {
	rpc.RegisterSwarmNodeServer(server, service)
}

// Bootstrap the Node into the network. The Node will connect to each bootstrap
// Node and attempt to find itself in the network. This process will ultimately
// connect it to Nodes that are close to it in XOR space.
func (service *SwarmService) Bootstrap() {
	// Add all bootstrap Nodes to the DHT.
	for _, bootstrapMultiAddress := range service.Options.BootstrapMultiAddresses {
		err := service.DHT.UpdateMultiAddress(bootstrapMultiAddress)
		if err != nil && service.Options.Debug >= DebugLow {
			log.Println(err)
		}
	}
	if service.Options.Concurrent {
		// Concurrently search all bootstrap Nodes for itself.
		do.ForAll(service.Options.BootstrapMultiAddresses, func(i int) {
			service.bootstrapUsingMultiAddress(service.Options.BootstrapMultiAddresses[i])
		})
	} else {
		// Sequentially search all bootstrap Nodes for itself.
		for _, bootstrapMultiAddress := range service.Options.BootstrapMultiAddresses {
			service.bootstrapUsingMultiAddress(bootstrapMultiAddress)
		}
	}
	if service.Options.Debug >= DebugMedium {
		log.Printf("%v connected to %v peers after bootstrapping.\n", service.Address(), len(service.DHT.MultiAddresses()))
	}
	if service.Options.Debug >= DebugHigh {
		log.Printf("%v is now connected to:\n", service.Address())
		for _, multiAddress := range service.DHT.MultiAddresses() {
			log.Printf("  %v\n", multiAddress)
		}
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
	if err := rpc.PingTarget(multiAddress, service.MultiAddress(), time.Minute); err != nil {
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
	if service.Options.Debug >= DebugHigh {
		log.Printf("%v was pinged by %v\n", service.Address(), from.Multi)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		nothing, err := service.ping(from)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(nothing)
	})

	select {
	case val := <-wait:
		if multiAddress, ok := val.Ok.(*rpc.MultiAddress); ok {
			return multiAddress, val.Err
		}
		return rpc.SerializeMultiAddress(service.MultiAddress()), val.Err

	case <-ctx.Done():
		return rpc.SerializeMultiAddress(service.MultiAddress()), ctx.Err()
	}
}

func (service *SwarmService) ping(from *rpc.MultiAddress) (*rpc.MultiAddress, error) {
	fromMultiAddress, err := rpc.DeserializeMultiAddress(from)
	if err != nil {
		return &rpc.Nothing{}, err
	}
	service.Delegate.OnPingReceived(fromMultiAddress)
	return rpc.SerializeMultiAddress(service.MultiAddress()), service.updatePeer(from)
}

// Query is used to return MultiAddresses that are closer to the given target
// Address. It will not return MultiAddresses that are further away from the
// target than the node itself, and it will only return MultiAddresses that are
// immediately connected to the service. The MultiAddresses returned are not
// guaranteed to be healthy connections and should be pinged.
func (service *SwarmService) Query(query *rpc.Address, stream rpc.Swarm_QueryServer) error {
	if service.Options.Debug >= DebugHigh {
		log.Printf("%v was queried by %v\n", service.Address(), query.From.Multi)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		return do.Err(service.query(query, stream))
	})

	for {
		select {
		case val := <-wait:
			return val.Err

		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}

func (service *SwarmService) query(query *rpc.Address, stream rpc.Swarm_QueryServer) error {
	deserializedQuery := rpc.DeserializeAddress(query)
	neighbors, err := service.DHT.FindMultiAddressNeighbors(deserializedQuery, service.Options.Alpha)
	if err != nil {
		return err
	}

	// Filter away peers that are further from the target than this service.
	for _, neighbor := range neighbors {
		closer, err := identity.Closer(neighbor.Address(), service.Address(), deserializedQuery)
		if err != nil {
			return err
		}
		if closer {
			if err := stream.Send(rpc.SerializeMultiAddress(neighbor)); err != nil {
				return err
			}
		}
	}

	service.Delegate.OnQuery(fromMultiAddress)
	return nil
}

// QueryDeep is used to return the closest MultiAddresses that can be reached
// from this node, relative to a target Address. It will not return
// MultiAddresses that are further away from the target than the node itself.
// The MultiAddresses returned are not guaranteed to be healthy connections
// and should be pinged.
func (service *SwarmService) QueryDeep(query *rpc.Address, stream rpc.Swarm_QueryDeepServer) error {
	if service.Options.Debug >= DebugHigh {
		log.Printf("%v was frontier queried by %v\n", service.Address(), query.From.Multi)
	}
	if err := stream.Context().Err(); err != nil {
		return err
	}

	wait := do.Process(func() do.Option {
		return do.Err(service.queryDeep(query, stream))
	})

	for {
		select {
		case val := <-wait:
			return val.Err

		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}

func (service *SwarmService) queryDeep(query *rpc.Address, stream rpc.Swarm_QueryDeepServer) error {

	// Get the target identity.Address for which this Node is searching for
	// peers.
	target := identity.Address(query.Query.Address)
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
			if err := stream.Send(rpc.SerializeMultiAddress(peer)); err != nil {
				return err
			}
			frontier = append(frontier, peer)
		}
	}

	// Immediately close the Node that sends the query and Node is running
	// the query and mark all peers in the frontier as seen.
	from, err := identity.NewMultiAddressFromString(query.From.Multi)
	if err != nil {
		return err
	}
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
		candidates, err := rpc.Query(peer, service.MultiAddress(), target, time.Second)
		if err != nil {
			if service.Options.Debug >= DebugLow {
				log.Println(err)
			}
			continue
		}

		// Filter any candidate that is already in the closure.
		for _, candidate := range candidates {
			if _, ok := visited[candidate.Address()]; ok {
				continue
			}
			// Expand the frontier by candidates that have not already been
			// explored, and store them in a persistent list of close peers.
			if err := stream.Send(rpc.SerializeMultiAddress(candidate)); err != nil {
				return err
			}
			frontier = append(frontier, candidate)
			visited[candidate.Address()] = struct{}{}
		}
	}

	fromMultiAddress, err := rpc.DeserializeMultiAddress(query.From)
	if err != nil {
		return err
	}
	service.Delegate.OnQueryCloserPeersOnFrontierReceived(fromMultiAddress)
	return service.updatePeer(query.From)
}

func (service *SwarmService) bootstrapUsingMultiAddress(bootstrapMultiAddress identity.MultiAddress) error {
	var err error
	var peers identity.MultiAddresses

	// The Node attempts to find itself in the network with three attempts
	// backing off by 10 seconds per attempt.
	for attempt := 0; attempt < service.Options.TimeoutRetries; attempt++ {
		// Query the bootstrap service.
		peers, err = rpc.QueryCloserPeersOnFrontierFromTarget(
			bootstrapMultiAddress,
			service.MultiAddress(),
			service.Address(),
			service.Options.Timeout+time.Duration(attempt)*service.Options.TimeoutStep,
		)
		// Errors are not returned because it is reasonable that a bootstrap
		// Node might be unavailable at this time.
		if err == nil {
			break
		}
		if service.Options.Debug >= DebugLow {
			log.Println(err)
		}
		if attempt == service.Options.TimeoutRetries-1 {
			return err
		}
	}

	// Peers returned by the query will be added to the DHT.
	if service.Options.Debug >= DebugMedium {
		log.Printf("%v received %v peers from %v.\n", service.Address(), len(peers), bootstrapMultiAddress.Address())
	}
	for _, peer := range peers {
		if peer.Address() == service.Address() {
			continue
		}
		if err := service.DHT.UpdateMultiAddress(peer); err != nil {
			if service.Options.Debug >= DebugLow {
				log.Println(err)
			}
		}
	}
	return nil
}

func (service *SwarmService) updatePeer(peer *rpc.MultiAddress) error {
	multiAddress, err := rpc.DeserializeMultiAddress(peer)
	if err != nil {
		return err
	}
	if multiAddress.Address() == service.Address() {
		return nil
	}
	if err := service.DHT.UpdateMultiAddress(multiAddress); err != nil {
		if err == dht.ErrFullBucket {
			pruned, err := service.Prune(multiAddress.Address())
			if err != nil {
				return err
			}
			if pruned {
				return service.DHT.UpdateMultiAddress(multiAddress)
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
	peers := service.DHT.MultiAddresses()

	// Create the frontier and a closure map.
	frontier := make(identity.MultiAddresses, 0, len(peers))
	visited := make(map[identity.Address]struct{})

	// Check if we know the target and filter away peers that are further
	// from the target than this Node.
	for _, peer := range peers {
		if peer.Address() == target {
			return &peer, nil
		}
		closer, err := identity.Closer(peer.Address(), service.Address(), target)
		if err != nil {
			return nil, err
		}
		if closer {
			frontier = append(frontier, peer)
		}
	}

	// Immediately close the Node that sends the query and Node is running
	// the query and mark all peers in the frontier as seen.
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

		candidates, err := rpc.QueryCloserPeersFromTarget(peer, service.MultiAddress(), target, time.Second)
		if err != nil {
			if service.Options.Debug >= DebugLow {
				log.Println(err)
			}
			continue
		}

		// Filter any candidate that is already in the closure.
		for _, candidate := range candidates {
			if _, ok := visited[candidate.Address()]; ok {
				continue
			}
			if candidate.Address() == target {
				return &candidate, nil
			}
			frontier = append(frontier, candidate)
			visited[candidate.Address()] = struct{}{}
		}
	}

	return nil, nil
}
