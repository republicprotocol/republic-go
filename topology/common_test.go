package topology_test

import (
	"fmt"
	"github.com/republicprotocol/go-identity"
	. "github.com/republicprotocol/go-swarm"
	"math/rand"
	"sync"
	"time"
)

// μ prevents multiple topology tests running in parallel. This is needed to
// protect overlapping ports during tests.
var μ *sync.Mutex = new(sync.Mutex)

// The number of peers that should be included in each topology test.
var numberOfPeers = 100

// The number of messages that will be sent through the topology.
var numberOfMessages = 100

// The duration to wait for peers to start listening for RPCs.
var startTimeDelay = time.Second

func generatePeers() ([]*Peer, error) {
	peers := make([]*Peer, numberOfPeers)
	for i := 0; i < numberOfPeers; i++ {
		keyPair, err := identity.NewKeyPair()
		if err != nil {
			return nil, err
		}
		multiAddress, err := identity.NewMultiAddress(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/republic/%s", 3000+i, keyPair.PublicAddress()))
		if err != nil {
			return nil, err
		}
		peers[i] = NewPeer(&Config{
			KeyPair:      keyPair,
			MultiAddress: multiAddress,
			Peers:        make([]identity.MultiAddress, 0, numberOfPeers-1),
		})
	}
	return peers, nil
}

func sendMessages(peers []*Peer) error {
	for i := 0; i < numberOfMessages; i++ {
		left, right := randomPeers(peers)
		if err := sendMessage(left, right.Config.MultiAddress); err != nil {
			return err
		}
	}
	return nil
}

func sendMessage(from *Peer, to identity.MultiAddress) error {
	_, err := from.SendFragment(to)
	return err
}

func randomPeers(peers []*Peer) (*Peer, *Peer) {
	left := rand.Intn(len(peers))
	right := rand.Intn(len(peers))
	for left == right {
		right = rand.Intn(len(peers))
	}
	return peers[left], peers[right]
}
