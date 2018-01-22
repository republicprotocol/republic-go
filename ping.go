package network

import (
	"log"

	"github.com/republicprotocol/go-dht"
	do "github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-network/rpc"
	"golang.org/x/net/context"
)

// Ping is used to test the connection to the Node and exchange MultiAddresses.
// If the Node does not respond, or it responds with an error, then the
// connection is considered unhealthy.
func (node *Node) Ping(ctx context.Context, peer *rpc.MultiAddress) (*rpc.MultiAddress, error) {
	log.Println("received call to Ping")

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	ch := do.Process(node.handlePing)
	defer close(ch)

	select {
	case multiAddressOpt := <-ch:
		if multiAddress, ok := multiAddressOpt.Ok.(identity.MultiAddress); ok {
			return &rpc.MultiAddress{Multi: multiAddress.String()}, multiAddressOpt.Err
		}
		return nil, multiAddressOpt.Err

	// Select the timeout from the context.
	case <-ctx.Done():
		return &rpc.MultiAddress{Multi: node.MultiAddress.String()}, ctx.Err()
	}
}

func (node *Node) handlePing(peer *rpc.MultiAddress) error {
	multi, err := identity.NewMultiAddressFromString(peer.Multi)
	if err != nil {
		return err
	}
	node.Delegate.OnPingReceived(multi)

	// Attempt to update the DHT.
	err = node.DHT.Update(multi)
	if err == dht.ErrFullBucket {
		// If the DHT is full then try and prune disconnected peers.
		address, err := multi.Address()
		if err != nil {
			return err
		}
		pruned, err := node.Prune(address)
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
