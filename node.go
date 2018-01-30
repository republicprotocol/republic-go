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
	"github.com/republicprotocol/go-rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// The Delegate is used as a callback interface to inject logic into the
// different RPCs.
type Delegate interface {
	OnPingReceived(from identity.MultiAddress)
	OnQueryCloserPeersReceived(from identity.MultiAddress)
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
		peers, err := rpc.QueryCloserPeersFromTarget(
			bootstrapMultiAddress,
			node.MultiAddress(),
			node.Address(),
			true,
		)
		if err != nil {
			if node.Options.Debug >= DebugLow {
				log.Println(err)
			}
			return
		}
		// All of the peers that it gets back will be added to the DHT.
		for _, peer := range peers {
			node.DHT.UpdateMultiAddress(peer)
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
	if _, err := rpc.PingTarget(multiAddress, node.MultiAddress()); err != nil {
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

func (node *Node) ping(from *rpc.MultiAddress) (*rpc.Nothing, error) {
	// Update the DHT.
	fromMultiAddress, err := rpc.DeserializeMultiAddress(from)
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
	fromMultiAddress, err := rpc.DeserializeMultiAddress(from)
	if err != nil {
		return &rpc.MultiAddresses{Multis: []*rpc.MultiAddress{}}, nil
	}
	if err := node.updatePeer(fromMultiAddress); err != nil {
		return &rpc.MultiAddresses{Multis: []*rpc.MultiAddress{}}, nil
	}

	// Return all peers in the DHT.
	peers := node.DHT.MultiAddresses()
	return rpc.SerializeMultiAddresses(peers), nil
}

func (node *Node) queryCloserPeers(query *rpc.Query) (*rpc.MultiAddresses, error) {
	// Update the DHT.
	if query.From != nil {
		fromMultiAddress, err := rpc.DeserializeMultiAddress(query.From)
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
	targetPeers := make(identity.MultiAddresses, 0, node.Options.Alpha)
	peers, err := node.DHT.FindMultiAddressNeighbors(target, node.Options.Alpha)
	if err != nil {
		return targetPeers, err
	}

	// Filter away peers that are further from the target than this Node.
	for _, peer := range peers {
		peerAddress := peer.Address()
		closer, err := identity.Closer(peerAddress, node.Address(), target)
		if err != nil {
			return rpc.SerializeMultiAddresses(targetPeers), err
		}
		if closer {
			targetPeers = append(targetPeers, peer)
		}
	}

	// If this is not a deep query, stop here.
	if !query.Deep {
		return targetPeers, nil
	}

	mu := new(sync.Mutex)
	open := true
	openList := make(identity.MultiAddresses, len(targetPeers))
	closeMap := map[string]bool{}
	do.ForAll(targetPeers, func(i int) {
		openList[i] = targetPeers[i]
	})
	for open {
		open = false
		openNext := make(identity.MultiAddresses, 0, len(openList))
		do.ForAll(openList, func(i int) {
			peers, err := rpc.QueryCloserPeersFromTarget(openList[i], node.MultiAddress(), query.Query, false)
			if err != nil {
				if node.Options.Debug >= DebugLow {
					log.Println(err)
					return
				}
			}

			mu.Lock()
			defer mu.Unlock()

			closeMap[openList[i]] = true
			for _, nextPeer := range peers {
				if closeMap[nextPeer] {
					continue
				}
				nextPeerAddress := nextPeer.Address()
				if closer, err := identity.Closer(nextPeerAddress, node.Address(), target); closer && err != nil {
					open = true
					openNext = append(openNext, nextPeer)
				}
			}
		})
		targetPeers = append(targetPeers, openList...)
		openList = openNext
	}

	sort.Slice(targetPeers, func(i, j int) bool {
		left := targetPeers[i].Address()
		right := targetPeers[j].Address()
		closer, _ := identity.Closer(left, right, target)
		return closer
	})

	minLength := len(targetPeers)
	if minLength > node.Options.Alpha {
		minLength = node.Options.Alpha
	}

	return rpc.SerializeMultiAddresses(targetPeers[:minLength]), nil
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
