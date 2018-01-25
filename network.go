package network

import (
	"fmt"
	"log"
	"net"
	"sort"
	"sync"

	"github.com/republicprotocol/go-dht"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-network/rpc"
	"github.com/republicprotocol/go-order-compute"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// The Delegate is used as a callback interface to inject logic into the
// different RPCs.
type Delegate interface {
	OnPingReceived(from identity.MultiAddress)
	OnOrderFragmentReceived(from identity.MultiAddress, orderFragment *compute.OrderFragment)
	OnResultFragmentReceived(from identity.MultiAddress, resultFragment *compute.ResultFragment)
}

// Node implements the gRPC Node service.
type Node struct {
	Delegate
	Server  *grpc.Server
	DHT     *dht.DHT
	Options Options
}

// NewNode returns a Node with the given its own identity.MultiAddress, a list
// of boostrap node identity.MultiAddresses, and a delegate that defines
// callbacks for each RPC.
func NewNode(delegate Delegate, options Options) *Node {
	return &Node{
		Delegate: delegate,
		Server:   grpc.NewServer(),
		DHT:      dht.NewDHT(options.MultiAddress.Address(), options.MaxBucketLength),
		Options:  options,
	}
}

// Serve starts the gRPC server.
func (node *Node) Serve() error {
	rpc.RegisterNodeServer(node.Server, node)
	host, err := node.MultiAddress().ValueForProtocol(identity.IP4Code)
	if err != nil {
		return err
	}
	port, err := node.MultiAddress().ValueForProtocol(identity.TCPCode)
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return err
	}
	if node.Options.Debug >= DebugLow {
		log.Printf("Listening at %s:%s\n", host, port)
	}
	return node.Server.Serve(listener)
}

// Stop the gRPC server.
func (node *Node) Stop() {
	if node.Options.Debug >= DebugLow {
		log.Printf("Stopping\n")
	}
	node.Server.Stop()
}

// Bootstrap the Node into the network. The Node will connect to each bootstrap
// Node and attempt to find itself in the network. This process will ultimately
// connect it to Nodes that are close to it in XOR space.
func (node *Node) Bootstrap() {
	do.CoForAll(node.Options.BootstrapMultiAddresses, func(i int) {
		// The Node attempts to find itself in the network.
		bootstrapMultiAddress := node.Options.BootstrapMultiAddresses[i]
		peers, err := QueryCloserPeersFromTarget(
			SerializeMultiAddress(bootstrapMultiAddress),
			SerializeMultiAddress(node.MultiAddress()),
			SerializeAddress(node.Address()),
			true,
		)
		if err != nil {
			if node.Options.Debug >= DebugLow {
				log.Println(err)
			}
			return
		}
		// All of the peers that it gets back will be added to the DHT.
		for _, peer := range peers.Multis {
			multiAddress, err := DeserializeMultiAddress(peer)
			if err != nil {
				if node.Options.Debug >= DebugLow {
					log.Println(err)
				}
				continue
			}
			node.DHT.UpdateMultiAddress(multiAddress)
		}
	})
}

// Prune an identity.Address from the dht.DHT. Returns a boolean indicating
// whether or not an identity.Address was pruned.
func (node *Node) Prune(target identity.Address) (bool, error) {
	bucket, err := node.DHT.FindBucket(target)
	if err != nil {
		return false, err
	}
	if bucket == nil || bucket.Length() == 0 {
		return false, nil
	}
	multiAddress := bucket.MultiAddresses[0]
	if _, err := PingTarget(SerializeMultiAddress(multiAddress), SerializeMultiAddress(node.MultiAddress())); err != nil {
		return true, node.DHT.RemoveMultiAddress(multiAddress)
	}
	return false, node.DHT.UpdateMultiAddress(multiAddress)
}

// Address returns the identity.Address of the Node.
func (node *Node) Address() identity.Address {
	return node.Options.MultiAddress.Address()
}

// MultiAddress returns the identity.MultiAddress of the Node.
func (node *Node) MultiAddress() identity.MultiAddress {
	return node.Options.MultiAddress
}

// Ping is used to test the connection to the Node and exchange
// identity.MultiAddresses. If the Node does not respond, or it responds with
// an error, then the connection should be considered unhealthy.
func (node *Node) Ping(ctx context.Context, from *rpc.MultiAddress) (*rpc.Nothing, error) {
	if node.Options.Debug >= DebugMedium {
		log.Printf("Ping received\n")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		nothing, err := node.ping(from)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(nothing)
	})

	select {
	case val := <-wait:
		if nothing, ok := val.Ok.(*rpc.Nothing); ok {
			return nothing, val.Err
		}
		return nil, val.Err

	case <-ctx.Done():
		return &rpc.Nothing{}, ctx.Err()
	}
}

// Peers is used to return the rpc.MultiAddresses to which a Node is connected.
// The rpc.MultiAddresses returned are not guaranteed to provide healthy
// connections and should be pinged.
func (node *Node) Peers(ctx context.Context, from *rpc.MultiAddress) (*rpc.MultiAddresses, error) {
	if node.Options.Debug >= DebugMedium {
		log.Printf("Peers received\n")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		peers, err := node.peers(from)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(peers)
	})

	select {
	case val := <-wait:
		if multiAddresses, ok := val.Ok.(*rpc.MultiAddresses); ok {
			return multiAddresses, val.Err
		}
		return nil, val.Err

	case <-ctx.Done():
		return &rpc.MultiAddresses{Multis: []*rpc.MultiAddress{}}, ctx.Err()
	}
}

// QueryCloserPeers is used to return the closest rpc.MultiAddresses to a peer
// with the given target rpc.Address. It will not return rpc.MultiAddresses
// that are further away from the target than the Node itself. The
// rpc.MultiAddresses returned are not guaranteed to provide healthy
// connections and should be pinged.
func (node *Node) QueryCloserPeers(ctx context.Context, query *rpc.Query) (*rpc.MultiAddresses, error) {
	if node.Options.Debug >= DebugMedium {
		log.Printf("QueryCloserPeers received\n")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		peers, err := node.queryCloserPeers(query)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(peers)
	})

	select {
	case val := <-wait:
		if multiAddresses, ok := val.Ok.(*rpc.MultiAddresses); ok {
			return multiAddresses, val.Err
		}
		return nil, val.Err

	case <-ctx.Done():
		return &rpc.MultiAddresses{Multis: []*rpc.MultiAddress{}}, ctx.Err()
	}
}

// SendOrderFragment to the Node. If the rpc.OrderFragment is not destined for
// this Node then it will be forwarded on to the correct destination.
func (node *Node) SendOrderFragment(ctx context.Context, orderFragment *rpc.OrderFragment) (*rpc.Nothing, error) {
	if node.Options.Debug >= DebugMedium {
		log.Printf("SendOrderFragment received\n")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		nothing, err := node.sendOrderFragment(orderFragment)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(nothing)
	})

	select {
	case val := <-wait:
		if nothing, ok := val.Ok.(*rpc.Nothing); ok {
			return nothing, val.Err
		}
		return nil, val.Err

	case <-ctx.Done():
		return &rpc.Nothing{}, ctx.Err()
	}
}

// SendResultFragment to the Node. If the rpc.ResultFragment is not destined
// for this Node then it will be forwarded on to the correct destination.
func (node *Node) SendResultFragment(ctx context.Context, resultFragment *rpc.ResultFragment) (*rpc.Nothing, error) {
	if node.Options.Debug >= DebugMedium {
		log.Printf("SendResultFragment received\n")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		nothing, err := node.sendResultFragment(resultFragment)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(nothing)
	})

	select {
	case val := <-wait:
		if nothing, ok := val.Ok.(*rpc.Nothing); ok {
			return nothing, val.Err
		}
		return nil, val.Err

	case <-ctx.Done():
		return &rpc.Nothing{}, ctx.Err()
	}
}

func (node *Node) ping(from *rpc.MultiAddress) (*rpc.Nothing, error) {
	// Update the DHT.
	fromMultiAddress, err := DeserializeMultiAddress(from)
	if err != nil {
		return &rpc.Nothing{}, err
	}
	if err := node.updatePeer(fromMultiAddress); err != nil {
		return &rpc.Nothing{}, err
	}

	// Notify the delegate of the ping.
	node.Delegate.OnPingReceived(fromMultiAddress)
	return &rpc.Nothing{}, nil
}

func (node *Node) peers(from *rpc.MultiAddress) (*rpc.MultiAddresses, error) {
	// Update the DHT.
	fromMultiAddress, err := DeserializeMultiAddress(from)
	if err != nil {
		return &rpc.MultiAddresses{Multis: []*rpc.MultiAddress{}}, nil
	}
	if err := node.updatePeer(fromMultiAddress); err != nil {
		return &rpc.MultiAddresses{Multis: []*rpc.MultiAddress{}}, nil
	}

	// Return all peers in the DHT.
	peers := node.DHT.MultiAddresses()
	return SerializeMultiAddresses(peers), nil
}

func (node *Node) queryCloserPeers(query *rpc.Query) (*rpc.MultiAddresses, error) {
	// Update the DHT.
	if query.From != nil {
		fromMultiAddress, err := DeserializeMultiAddress(query.From)
		if err != nil {
			return &rpc.MultiAddresses{Multis: []*rpc.MultiAddress{}}, err
		}
		if err := node.updatePeer(fromMultiAddress); err != nil {
			return &rpc.MultiAddresses{Multis: []*rpc.MultiAddress{}}, err
		}
	}

	// Get the target identity.Address for which this Node is searching for
	// peers.
	target := identity.Address(query.Query.Address)
	targetPeers := &rpc.MultiAddresses{Multis: make([]*rpc.MultiAddress, 0, node.Options.Alpha)}
	peers, err := node.DHT.FindMultiAddressNeighbors(target, node.Options.Alpha)
	if err != nil {
		return targetPeers, err
	}

	// Filter away peers that are further from the target than this Node.
	for _, peer := range peers {
		peerAddress := peer.Address()
		closer, err := identity.Closer(peerAddress, node.Address(), target)
		if err != nil {
			return targetPeers, err
		}
		if closer {
			targetPeers.Multis = append(targetPeers.Multis, SerializeMultiAddress(peer))
		}
	}

	// If this is not a deep query, stop here.
	if !query.Deep {
		return targetPeers, nil
	}

	mu := new(sync.Mutex)
	open := true
	openList := make([]*rpc.MultiAddress, len(targetPeers.Multis))
	closeMap := map[string]bool{}
	do.ForAll(targetPeers.Multis, func(i int) {
		openList[i] = targetPeers.Multis[i]
	})
	for open {
		open = false
		openNext := make([]*rpc.MultiAddress, 0, len(openList))
		do.ForAll(openList, func(i int) {
			peers, err := QueryCloserPeersFromTarget(openList[i], SerializeMultiAddress(node.MultiAddress()), query.Query, false)
			if err != nil {
				if node.Options.Debug >= DebugLow {
					log.Println(err)
					return
				}
			}

			mu.Lock()
			defer mu.Unlock()

			closeMap[openList[i].Multi] = true
			for _, nextPeer := range peers.Multis {
				if closeMap[nextPeer.Multi] {
					continue
				}
				open = true
				openNext = append(openNext, nextPeer)
			}
		})
		targetPeers.Multis = append(targetPeers.Multis, openList...)
		openList = openNext
	}

	sort.Slice(openList, func(i, j int) bool {
		leftMultiAddress, _ := DeserializeMultiAddress(openList[i])
		left := leftMultiAddress.Address()
		rightMultiAddress, _ := DeserializeMultiAddress(openList[j])
		right := rightMultiAddress.Address()
		closer, _ := identity.Closer(left, right, target)
		return closer
	})

	minLength := len(openList)
	if minLength > node.Options.Alpha {
		minLength = node.Options.Alpha
	}

	return &rpc.MultiAddresses{Multis: targetPeers.Multis[:minLength]}, nil
}

func (node *Node) sendOrderFragment(orderFragment *rpc.OrderFragment) (*rpc.Nothing, error) {
	// Update the DHT.
	var err error
	var fromMultiAddress identity.MultiAddress
	if orderFragment.From != nil {
		fromMultiAddress, err := DeserializeMultiAddress(orderFragment.From)
		if err != nil {
			return &rpc.Nothing{}, err
		}
		if err := node.updatePeer(fromMultiAddress); err != nil {
			return &rpc.Nothing{}, err
		}
	}

	// Check if the rpc.OrderFragment has reached its destination.
	if orderFragment.To.Address == node.Address().String() {
		deserializedOrderFragment, err := DeserializeOrderFragment(orderFragment)
		if err != nil {
			return &rpc.Nothing{}, err
		}
		node.OnOrderFragmentReceived(fromMultiAddress, deserializedOrderFragment)
		return &rpc.Nothing{}, nil
	}

	// Forward the rpc.OrderFragment to the closest peers.
	peers, err := node.queryCloserPeers(&rpc.Query{From: nil, Query: orderFragment.To})
	if err != nil {
		return &rpc.Nothing{}, err
	}
	do.CoForAll(peers.Multis, func(i int) {
		_, err := SendOrderFragmentToTarget(peers.Multis[i], orderFragment)
		if err != nil && node.Options.Debug >= DebugLow {
			log.Println(err)
		}
	})
	return &rpc.Nothing{}, nil
}

func (node *Node) sendResultFragment(resultFragment *rpc.ResultFragment) (*rpc.Nothing, error) {
	// Update the DHT.
	var err error
	var fromMultiAddress identity.MultiAddress
	if resultFragment.From != nil {
		if resultFragment.From.Multi == node.MultiAddress().String() {
			return &rpc.Nothing{}, nil
		}
		fromMultiAddress, err = DeserializeMultiAddress(resultFragment.From)
		if err != nil {
			return &rpc.Nothing{}, err
		}
		if err := node.updatePeer(fromMultiAddress); err != nil {
			return &rpc.Nothing{}, err
		}
	}

	// Check if the rpc.OrderFragment has reached its destination.
	if resultFragment.To.String() == node.Address().String() {
		deserializedResultFragment, err := DeserializeResultFragment(resultFragment)
		if err != nil {
			return &rpc.Nothing{}, err
		}
		node.OnResultFragmentReceived(fromMultiAddress, deserializedResultFragment)
		return &rpc.Nothing{}, nil
	}

	// Forward the rpc.OrderFragment to the closest peers.
	peers, err := node.queryCloserPeers(&rpc.Query{From: nil, Query: resultFragment.To, Deep: true})
	if err != nil {
		return &rpc.Nothing{}, err
	}
	for _, peer := range peers.Multis {
		go SendResultFragmentToTarget(peer, resultFragment)
	}
	return &rpc.Nothing{}, nil
}

func (node *Node) updatePeer(multiAddress identity.MultiAddress) error {
	if err := node.DHT.UpdateMultiAddress(multiAddress); err != nil {
		if err == dht.ErrFullBucket {
			pruned, err := node.Prune(multiAddress.Address())
			if err != nil {
				return err
			}
			if pruned {
				return node.DHT.UpdateMultiAddress(multiAddress)
			}
			return nil
		}
		log.Println(err)
		return err
	}
	return nil
}
