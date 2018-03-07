package node

import (
	"fmt"
	"log"

	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-rpc"
)

// OnPingReceived implements the swarm.Delegate interface. It is used by the
// underlying swarm.Node whenever the Miner has handled a Ping RPC.
func (node *DarkNode) OnPingReceived(peer identity.MultiAddress) {
	// TODO: Log metrics for the ping.
}

// OnQueryCloserPeersReceived implements the swarm.Delegate interface. It is
// used by the underlying swarm.Node whenever the Miner has handled a
// QueryCloserPeers RPC.
func (node *DarkNode) OnQueryCloserPeersReceived(peer identity.MultiAddress) {
	// TODO: Log metrics for the query.
}

// OnQueryCloserPeersOnFrontierReceived implements the swarm.Delegate
// interface. It is called by the underlying swarm.Node whenever the Miner
// has handled a QueryCloserPeersOnFrontier RPC.
func (node *DarkNode) OnQueryCloserPeersOnFrontierReceived(peer identity.MultiAddress) {
	// TODO: Log metrics for the deep query.
}

func (node *DarkNode) OnSync(from identity.MultiAddress) chan do.Option {
	// FIXME: Synchronize.
	panic("uninmplemented")
}

// OnOrderFragmentReceived implements the dark.Delegate interface. It is called
// by the underlying dark.Node whenever it receives a compute.OrderFragment
// that it must process.
func (node *DarkNode) OnOrderFragmentReceived(from identity.MultiAddress, orderFragment *compute.OrderFragment) {
	go func() {
		deltaFragments, err := node.DeltaFragmentMatrix.InsertOrderFragment(orderFragment)
		if err != nil {
			log.Println(err)
		}
		for _, deltaFragment := range deltaFragments {
			node.DeltaBuilder.InsertDeltaFragment(deltaFragment)
		}
		node.DarkPool.EnterReadOnly(nil)
		defer node.DarkPool.ExitReadOnly()
		for _, multiAddress := range node.DarkPool.Nodes {
			client, err := rpc.NewClient(multiAddress, node.Swarm.MultiAddress())
			if err != nil {
				log.Println(err)
				continue
			}
			for _, deltaFragment := range deltaFragments {
				_, err := client.BroadcastDeltaFragment(deltaFragment)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
}

// OnBroadcastDeltaFragment implements the dark.Delegate interface. It is
// called by the underlying dark.Node whenever it receives a
// compute.DeltaFragment that it must process.
func (node *DarkNode) OnBroadcastDeltaFragment(from identity.MultiAddress, deltaFragment *compute.DeltaFragment) {
	go func() {
		delta, err := node.DeltaBuilder.InsertDeltaFragment(deltaFragment)
		if err != nil {
			log.Println(err)
			return
		}
		if delta == nil {
			return
		}
		if delta.IsMatch(node.Configuration.Prime) {
			log.Printf("%v compared (%v, %v) <- MATCH\n", node.Swarm.Address(), base58.Encode(deltaFragment.BuyOrderID), base58.Encode(deltaFragment.SellOrderID))
			node.log(
				"match",
				fmt.Sprintf(
					`{"id": "%s", "buyID": "%s", "sellID": "%s"}`,
					node.Configuration.MultiAddress.String(),
					base58.Encode(deltaFragment.BuyOrderID),
					base58.Encode(deltaFragment.SellOrderID),
				),
			)
			// TODO: Attempt to get consensus on the match and then mark the orders
			// handled if the consensus is won. If the consensus is not won take
			// either the buy, or sell (or both), orders and mark them as completed
			// (this depends on which ones conflicted).
		} else {
			log.Printf("%v compared (%v, %v)\n", node.Swarm.Address(), base58.Encode(deltaFragment.BuyOrderID), base58.Encode(deltaFragment.SellOrderID))
		}
	}()
}

// SubscribeToLogs will start sending log events to logChannel
func (node *DarkNode) SubscribeToLogs(logChannel chan do.Option) {
	node.logQueue.Subscribe(logChannel)
}

// UnsubscribeFromLogs will stop sending log events to logChannel
func (node *DarkNode) UnsubscribeFromLogs(logChannel chan do.Option) {
	node.logQueue.Unsubscribe(logChannel)
}
