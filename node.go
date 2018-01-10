package x

import (
	"fmt"
	"log"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/republicprotocol/go-dht"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-x/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// α determines the maximum number of concurrent client connections that the
// Node is expected to use when running a distributed Dijkstra search.
const α = 3

// Node implements the gRPC Node service.
type Node struct {
	*grpc.Server
	KeyPair      identity.KeyPair
	MultiAddress identity.MultiAddress
	DHT          *dht.DHT
}

// NewNode returns a Node with the given Config, a new DHT, and a new set of grpc.Connections.
func NewNode(config *Config) (*Node, error) {
	dht := dht.NewDHT(config.KeyPair.Address())
	for _, peer := range config.Peers {
		if err := dht.Update(peer); err != nil {
			return nil, err
		}
	}
	return &Node{
		Server:       grpc.NewServer(),
		KeyPair:      config.KeyPair,
		MultiAddress: config.MultiAddress,
		DHT:          dht,
	}, nil
}

// Serve starts the gRPC server.
func (node *Node) Serve() error {
	host, err := node.MultiAddress.ValueForProtocol(identity.IP4Code)
	if err != nil {
		return err
	}
	port, err := node.MultiAddress.ValueForProtocol(identity.TCPCode)
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return err
	}
	rpc.RegisterNodeServer(node.Server, node)
	reflection.Register(node.Server)
	return node.Server.Serve(listener)
}

// Ping is used to test the connection to the Node and exchange MultiAddresses.
// If the Node does not respond, or it responds with an error, then the
// connection is considered unhealthy.
func (node *Node) Ping(ctx context.Context, peer *rpc.MultiAddress) (*rpc.MultiAddress, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Spawn a goroutine to evaluate the return value.
	wait := make(chan error)
	go func() {
		defer close(wait)
		wait <- node.ping(peer)
	}()

	select {
	case ret := <-wait:
		return &rpc.MultiAddress{Multi: node.MultiAddress.String()}, ret

	// Select the timeout from the context.
	case <-ctx.Done():
		return &rpc.MultiAddress{Multi: node.MultiAddress.String()}, ctx.Err()
	}
}

// Peers is used to return the rpc.MultiAddresses to which a Node is connected.
// The rpc.MultiAddresses returned are not guaranteed to provide healthy
// connections and should be pinged.
func (node *Node) Peers(ctx context.Context, sender *rpc.Nothing) (*rpc.MultiAddresses, error) {
	if err := ctx.Err(); err != nil {
		return &rpc.MultiAddresses{}, err
	}

	// Spawn a goroutine to evaluate the return value.
	wait := make(chan *rpc.MultiAddresses)
	go func() {
		defer close(wait)
		wait <- node.peers()
	}()

	select {
	case ret := <-wait:
		return ret, nil

	// Select the timeout from the context.
	case <-ctx.Done():
		return &rpc.MultiAddresses{}, ctx.Err()
	}
}

// SendOrderFragment is used to forward an rpc.OrderFragment through the X
// Network to its destination Node. This forwarding is done using a distributed
// Dijkstra search, using the XOR distance between identity.Addresses as the
// distance heuristic.
func (node *Node) SendOrderFragment(ctx context.Context, orderFragment *rpc.OrderFragment) (*rpc.Nothing, error) {
	// Check for errors in the context.
	if err := ctx.Err(); err != nil {
		return &rpc.Nothing{}, err
	}

	// Spawn a goroutine to evaluate the return value.
	wait := make(chan error)
	go func() {
		defer close(wait)
		wait <- node.sendOrderFragment(orderFragment)
	}()

	select {
	// Select the timeout from the context.
	case <-ctx.Done():
		return &rpc.Nothing{}, ctx.Err()

	// Select the value passed by the goroutine.
	case ret := <-wait:
		if ret != nil {
			return &rpc.Nothing{}, ret
		}
		return &rpc.Nothing{}, nil
	}
}

func (node *Node) ping(peer *rpc.MultiAddress) error {
	multi, err := identity.NewMultiAddress(peer.Multi)
	if err != nil {
		return err
	}
	// Attempt to update the DHT.
	err = node.DHT.Update(multi)
	if err == dht.ErrFullBucket {
		// If the DHT is full then try and prune disconnected peers.
		address, err := multi.Address()
		if err != nil {
			return err
		}
		pruned, err := node.pruneMostRecentPeer(address)
		if err != nil {
			return err
		}
		// If a peer was pruned, then update the DHT again.
		if pruned {
			return node.DHT.Update(multi)
		}
		return nil
	}
	return err
}

func (node *Node) peers() *rpc.MultiAddresses {
	// Get all identity.MultiAddresses in the DHT.
	multis := node.DHT.MultiAddresses()
	ret := &rpc.MultiAddresses{
		Multis: make([]*rpc.MultiAddress, len(multis)),
	}
	// Transform them into rpc.MultiAddresses.
	for i, multi := range multis {
		ret.Multis[i] = &rpc.MultiAddress{Multi: multi.String()}
	}
	return ret
}

func (node *Node) sendOrderFragment(orderFragment *rpc.OrderFragment) error {
	target := identity.Address(orderFragment.To)
	if string(target) == string(node.DHT.Address) {
		// TODO: This Node is the intended target! Do something with the
		//       rpc.OrderFragment.
		log.Println("rpc.OrderFragment", orderFragment.OrderFragmentID, "received!")
		return nil
	}

	targetMultiMu := new(sync.Mutex)
	var targetMulti *identity.MultiAddress

	openMu := new(sync.Mutex)
	bucket, err := node.DHT.FindBucket(target)
	if err != nil {
		return err
	}
	open := bucket.Multiaddresses()
	closed := make(map[identity.MultiAddress]bool)

	// TODO: We are only using one dht.Bucket to search the network. If this
	//       dht.Bucket is not sufficient, we should also search the
	//       neighborhood of the dht.Bucket. The neighborhood should expand by
	//       a distance of one, until the entire DHT has been searched, or the
	//       targetMulti has been found.
	for len(open) > 0 {
		var wg sync.WaitGroup
		wg.Add(α)
		badNodes := make(chan identity.MultiAddress, α)

		// Take the first α multi-addresses from the open list and expand them
		// concurrently. This moves them from the open list to the closed list,
		// preventing the same multi-address from being expanded more than
		// once.
		for i := 0; i < α; i++ {
			multi := open[0]
			open = open[1:]
			closed[multi] = true

			go func() {
				defer wg.Done()

				// Get all peers of this multi-address. This is the expansion
				// step of the search.
				peers, err := Peers(multi,&rpc.MultiAddress{Multi:node.MultiAddress.String()})
				if err != nil {
					badNodes <- multi
					return
				}

				// Traverse all peers and collect them into the openNext, a
				// list of peers that we want to add the open list.
				openNext := make(identity.MultiAddresses, 0, len(peers))
				for _, peer := range peers {
					address, err := peer.Address()
					if err != nil {
						badNodes <- peer
						return
					}
					if string(target) == string(address) {
						// If we have found the target, set the targetMulti and
						// exit the loop. There is no point acquiring more
						// peers for the open list.
						targetMultiMu.Lock()
						if targetMulti == nil {
							targetMulti = &peer
						}
						targetMultiMu.Unlock()
						break
					}
					// Otherwise, store this peer's multi-address in the
					// openNext list. It will be added to the open list if it
					// has not already been closed.
					openNext = append(openNext, peer)
				}

				targetMultiMu.Lock()
				if targetMulti == nil {
					// We have not found the targetMulti yet.
					targetMultiMu.Unlock()
					openMu.Lock()
					// Add new peers to the open list if they have not been
					// closed.
					for _, next := range openNext {
						if _, ok := closed[next];!ok {
							open = append(open, next)
						}
					}
					openMu.Unlock()
				} else {
					targetMultiMu.Unlock()
				}
			}()
		}
		wg.Wait()
		if targetMulti != nil {
			// If we have found the targetMulti then we can exit the search
			// loop.
			break
		}

		// Remove bad nodes which do not respond
		for n := range badNodes{
			node.DHT.Remove(n)
		}

		// Otherwise, sort the open list by distance to the target.
		sort.SliceStable(open, func(i, j int) bool {
			left := open[i]
			right := open[j]
			closer, _ := identity.Closer(left,right,target)
			return closer
		})
	}

	if targetMulti == nil {
		return fmt.Errorf("cannot find target")
	}
	err = SendOrderFragment(*targetMulti, orderFragment)
	if err!= nil {
		return err
	}
	return nil
}

func (node *Node) pruneMostRecentPeer(target identity.Address) (bool, error) {
	bucket, err := node.DHT.FindBucket(target)
	if err != nil {
		return false, err
	}
	multi := bucket.OldestMultiAddress()
	if multi == nil {
		return false, nil
	}


	err = Ping(*multi,&rpc.MultiAddress{Multi:node.MultiAddress.String()})
	if err != nil {
		// If the connection could not be made, prune the peer.
		return true, node.DHT.Remove(*multi)
	}
	return false,nil
}
