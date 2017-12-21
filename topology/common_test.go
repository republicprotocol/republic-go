package topology_test

import (
	"sync"
	"github.com/republicprotocol/go-identity"
	"math/rand"
)

// μ prevents multiple topology tests running in parallel. This is needed to
// protect overlapping ports during tests.
var μ *sync.Mutex = new(sync.Mutex)

// The number of peers that should be included in each topology test.
var numberOfPeers = 100

// The number of messages that will be sent through the topology.
var numberOfMessages = 100

func generatePeers() ([]*Peer, error){
	peers := make([]*Peer, numberOfPeers)
	for i := 0; i < numberOfPeers; i++ {
		keyPair, err := identity.NewKeyPair()
		if err != nil {
			return nil, err
		}
		config := Config {
			Address: keyPair.PublicAddress(),
			Host: "0.0.0.0",
			Port: 3000 + i,
			Peers: make([]Config, 0, numberOfPeers-1),
		}
		peers[i] = NewPeer(keyPair, config)
	}
	return peers, nil
}

func sendMessages(peers []*Peer) error {
	for i := 0; i < numberOfMessages; i++ {
		left, right := randomPeers(peers)
		if err := sendMessage(left, right); err != nil {
			return err
		}
	}
	return nil
}

func sendMessage(to *Peer, from *Peer) error {
	_, err := from.SendOrderFragment(&rpc.OrderFragment{
		To: to.Config.Address,
	}, to.Config.Address)
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