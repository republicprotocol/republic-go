package swarm

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/republicprotocol/go-dht"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-swarm/rpc"
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
	open, err := node.DHT.FindBucket(target).MultiAddresses()
	closed := make(map[identity.MultiAddresses]bool, 0, len(open))

	for len(open) > 0 {
		var wg sync.WaitGroup
		wg.Add(α)

		for i := 0; i < α; i++ {
			multi := open[0].MultiAddress
			open = open[1:]
			closed[multi] = true
			go func() {
				defer wg.Done()

				// Create a client connection to the peer.
				client, conn, err := NewNodeClient(multi)
				if err != nil {
					return
				}
				defer conn.Close()

				// Get peers of the peer.
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				peers, err = client.Peers(ctx, &rpc.Nothing{})
				if err != nil {
					return
				}

				// Traverse all peers.
				openNext := make(identity.MultiAddresses, 0, len(peers))
				for peer := range peers {
					multi, err := identity.NewMultiAddress(peer.Multi)
					if err != nil {
						continue
					}
					address, err := multi.Address()
					if err != nil {
						continue
					}
					if string(target) == string(address) {
						targetMultiMu.Lock()
						if targetMulti == nil {
							targetMulti = multi
						}
						targetMultiMu.Unlock()
						break
					}
					openNext = append(openNext, multi)
				}

				targetMultiMu.Lock()
				if targetMulti == nil {
					targetMultiMu.Unlock()
					openMu.Lock()
					for next := range openNext {
						if c, ok := closed[next]; !c || !ok {
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
			break
		}
		open.Sort()
	}

	if targetMulti == nil {
		return fmt.Errorf("cannot find target")
	}

	// Create a client connection to the peer.
	client, conn, err := NewNodeClient(targetMulti)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Send the order fragment on to the peer.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_, err = client.SendOrderFragment(ctx, orderFragment)
	if err != nil {
		return err
	}
	return false, nil
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

	// Create a client connection to the peer.
	client, conn, err := NewNodeClient(*multi)
	if err != nil {
		// If the connection could not be made, prune the peer.
		if err == context.DeadlineExceeded || err == context.Canceled {
			return true, node.DHT.Remove(*multi)
		}
		return false, err
	}
	defer conn.Close()

	// Ping the peer.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_, err = client.Ping(ctx, &rpc.MultiAddress{Multi: node.MultiAddress.String()})
	if err != nil {
		// If the ping could not be made, prune the peer.
		if err == context.DeadlineExceeded || err == context.Canceled {
			return true, node.DHT.Remove(*multi)
		}
		return false, err
	}
	return false, nil
}
