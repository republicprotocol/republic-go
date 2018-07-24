package smpc

import (
	"encoding/base64"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/swarm"
	"golang.org/x/net/context"
)

type Sender interface {
	Send(message Message) error
}

type Receiver interface {
	Receive(from identity.Address, message Message)
}

// NetworkID for a network of Smpcer nodes. Using a NetworkID allows nodes to
// be involved in multiple distinct computation networks in parallel.
type NetworkID [32]byte

// String returns a human-readable representation of a NetworkID.
func (id NetworkID) String() string {
	return base64.StdEncoding.EncodeToString(id[:8])
}

type NetworkConnector interface {
	Connect(ctx context.Context, networkID NetworkID, to identity.MultiAddress, receiver Receiver) (Sender, error)
	Listen(ctx context.Context, networkID NetworkID, to identity.MultiAddress, receiver Receiver) (Sender, error)
}

// NetworkOverlay provides an abstraction for message passing over multiple
// networks.
type NetworkOverlay interface {
	Receiver

	// Connect to a new network of addresses.
	Connect(networkID NetworkID, addrs identity.Addresses)

	// Disconnect from an existing network.
	Disconnect(networkID NetworkID)

	// Send a message to all addresses in a connected network.
	Send(networkID NetworkID, message Message)

	// Send a message to a specific address on a specific network.
	SendTo(networkID NetworkID, to identity.Address, message Message)
}

type networkOverlay struct {
	connector NetworkConnector
	receiver  Receiver
	swarmer   swarm.Swarmer

	networkMu      *sync.RWMutex
	network        map[NetworkID](map[identity.Address]Sender)
	networkCancels map[NetworkID](map[identity.Address]context.CancelFunc)

	lookupMu *sync.RWMutex
	lookup   map[identity.Address]identity.MultiAddress
}

func NewNetworkOverlay(connector NetworkConnector, receiver Receiver, swarmer swarm.Swarmer) NetworkOverlay {
	panic("unimplemented")
}

// Connect implements the NetworkOverlay interface.
func (overlay *networkOverlay) Connect(networkID NetworkID, addrs identity.Addresses) {
	k := int64(2 * (len(addrs) + 1) / 3)

	log.Printf("[info] connecting to network %v with thresold = (%v, %v)", networkID, len(addrs), k)

	overlay.networkMu.Lock()
	defer overlay.networkMu.Unlock()

	overlay.network[networkID] = map[identity.Address]Sender{}
	overlay.networkCancels[networkID] = map[identity.Address]context.CancelFunc{}

	go dispatch.CoForAll(addrs, func(i int) {
		addr := addrs[i]
		if addr == overlay.swarmer.MultiAddress().Address() {
			// Skip trying to connect to ourself
			return
		}

		log.Printf("[debug] querying peer %v on network %v", addr, networkID)
		multiAddr, err := overlay.query(addr)
		if err != nil {
			log.Printf("[error] cannot connect to peer %v on network %v: %v", addr, networkID, err)
			return
		}

		// Store the identity.Identity to identity.MultiAddress mapping
		overlay.lookupMu.Lock()
		overlay.lookup[addr] = multiAddr
		overlay.lookupMu.Unlock()

		var sender Sender
		ctx, cancel := context.WithCancel(context.Background())

		if addr < overlay.swarmer.MultiAddress().Address() {
			log.Printf("[debug] connecting to peer %v on network %v", addr, networkID)
			sender, err = overlay.connector.Connect(ctx, networkID, multiAddr, overlay)
			if err != nil {
				log.Printf("[error] cannot connect to peer %v on network %v: %v", addr, networkID, err)
				return
			}
			log.Printf("[debug] connected to peer %v on network %v", addr, networkID)
		} else {
			log.Printf("[debug] listening for peer %v on network %v", addr, networkID)
			sender, err = overlay.connector.Listen(ctx, networkID, multiAddr, overlay)
			if err != nil {
				log.Printf("[error] cannot listen for peer %v on network %v: %v", addr, networkID, err)
				return
			}
			log.Printf("[debug] accepted peer %v on network %v", addr, networkID)
		}

		overlay.networkMu.Lock()
		defer overlay.networkMu.Unlock()

		overlay.network[networkID][addr] = sender
		overlay.networkCancels[networkID][addr] = cancel
	})
}

// Disconnect implements the NetworkOverlay interface.
func (overlay *networkOverlay) Disconnect(networkID NetworkID) {
	log.Printf("[info] disconnecting from network %v", networkID)

	overlay.networkMu.Lock()
	defer overlay.networkMu.Unlock()

	cancels, ok := overlay.networkCancels[networkID]
	if ok {
		for _, cancel := range cancels {
			cancel()
		}
	}
	delete(overlay.network, networkID)
	delete(overlay.networkCancels, networkID)
}

// Send implements the NetworkOverlay interface.
func (overlay *networkOverlay) Send(networkID NetworkID, message Message) {
	overlay.networkMu.RLock()
	defer overlay.networkMu.RUnlock()

	peers, ok := overlay.network[networkID]
	if !ok {
		log.Printf("[error] cannot send message to unknown network %v", networkID)
		return
	}

	go dispatch.CoForAll(peers, func(peer identity.Address) {
		peers[peer].Send(message)
	})
}

// SendTo implements the NetworkOverlay interface.
func (overlay *networkOverlay) SendTo(networkID NetworkID, to identity.Address, message Message) {
	overlay.networkMu.RLock()
	defer overlay.networkMu.RUnlock()

	peers, ok := overlay.network[networkID]
	if !ok {
		log.Printf("[error] cannot send message to unknown network %v", networkID)
		return
	}
	peer, ok := peers[to]
	if !ok {
		log.Printf("[error] cannot send message to unknown peer %v", to)
		return
	}

	go peer.Send(message)
}

// OnReceiveMessage implements the NetworkOverlay interface.
func (overlay *networkOverlay) Receive(from identity.Address, message Message) {
	overlay.receiver.Receive(from, message)
}

func (overlay *networkOverlay) query(q identity.Address) (identity.MultiAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	multiAddr, err := overlay.swarmer.Query(ctx, q, -1)
	if err != nil {
		return multiAddr, fmt.Errorf("cannot query peer %v: %v", q, err)
	}
	return multiAddr, nil
}
