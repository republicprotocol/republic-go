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

type Connector interface {
	Connect(ctx context.Context, networkID NetworkID, to identity.MultiAddress, receiver Receiver) (Sender, error)
}

type Listener interface {
	Listen(ctx context.Context, networkID NetworkID, to identity.Address, receiver Receiver) (Sender, error)
}

type ConnectorListener interface {
	Connector
	Listener
}

// Network provides an abstraction for message passing over multiple networks.
type Network interface {

	// Connect to a new network of addresses.
	Connect(networkID NetworkID, addrs identity.Addresses)

	// Disconnect from an existing network.
	Disconnect(networkID NetworkID)

	// Send a message to all addresses in a connected network.
	Send(networkID NetworkID, message Message)

	// Send a message to a specific address on a specific network.
	SendTo(networkID NetworkID, to identity.Address, message Message)
}

type network struct {
	conn     ConnectorListener
	receiver Receiver
	swarmer  swarm.Swarmer

	networkMu      *sync.RWMutex
	networkSenders map[NetworkID](map[identity.Address]Sender)
	networkCancels map[NetworkID](map[identity.Address]context.CancelFunc)
}

func NewNetwork(conn ConnectorListener, receiver Receiver, swarmer swarm.Swarmer) Network {
	return &network{
		conn:     conn,
		receiver: receiver,
		swarmer:  swarmer,

		networkMu:      new(sync.RWMutex),
		networkSenders: map[NetworkID](map[identity.Address]Sender){},
		networkCancels: map[NetworkID](map[identity.Address]context.CancelFunc){},
	}
}

// Connect implements the Network interface.
func (network *network) Connect(networkID NetworkID, addrs identity.Addresses) {

	k := int64(2 * (len(addrs) + 1) / 3)
	log.Printf("[info] connecting to network %v with thresold = (%v, %v)", networkID, len(addrs), k)

	func() {
		network.networkMu.Lock()
		defer network.networkMu.Unlock()
		network.networkSenders[networkID] = map[identity.Address]Sender{}
		network.networkCancels[networkID] = map[identity.Address]context.CancelFunc{}
	}()

	go dispatch.CoForAll(addrs, func(i int) {
		addr := addrs[i]
		if addr == network.swarmer.MultiAddress().Address() {
			// Skip trying to connect to ourself
			return
		}

		// Create and store a context for connections
		ctx, cancel := context.WithCancel(context.Background())
		func() {
			network.networkMu.Lock()
			defer network.networkMu.Unlock()
			network.networkCancels[networkID][addr] = cancel
		}()

		// Connect, or listen for a connection, and store the sending handle
		sender := network.connectOrListen(ctx, networkID, addr)
		if sender == nil {
			return
		}

		func() {
			network.networkMu.Lock()
			defer network.networkMu.Unlock()
			if _, ok := network.networkSenders[networkID]; ok {
				network.networkSenders[networkID][addr] = sender
			}
		}()
	})
}

// Disconnect implements the Network interface.
func (network *network) Disconnect(networkID NetworkID) {
	log.Printf("[info] disconnecting from network %v", networkID)

	network.networkMu.Lock()
	defer network.networkMu.Unlock()

	cancels, ok := network.networkCancels[networkID]
	if ok {
		for _, cancel := range cancels {
			cancel()
		}
	}
	delete(network.networkSenders, networkID)
	delete(network.networkCancels, networkID)
}

// Send implements the Network interface.
func (network *network) Send(networkID NetworkID, message Message) {
	network.networkMu.RLock()
	defer network.networkMu.RUnlock()

	senders, ok := network.networkSenders[networkID]
	if !ok {
		log.Printf("[error] cannot send message to unknown network %v", networkID)
		return
	}

	go dispatch.CoForAll(senders, func(addr identity.Address) {
		sender := senders[addr]
		if err := sender.Send(message); err != nil {
			// These logs are disabled to prevent verbose output
			// log.Printf("[error] cannot send message to %v on network %v: %v", addr, networkID, err)
		}
	})
}

// SendTo implements the Network interface.
func (network *network) SendTo(networkID NetworkID, to identity.Address, message Message) {
	network.networkMu.RLock()
	defer network.networkMu.RUnlock()

	senders, ok := network.networkSenders[networkID]
	if !ok {
		log.Printf("[error] cannot send message to unknown network %v", networkID)
		return
	}
	sender, ok := senders[to]
	if !ok {
		log.Printf("[error] cannot send message to unknown peer %v", to)
		return
	}

	go func() {
		if err := sender.Send(message); err != nil {
			// These logs are disabled to prevent verbose output
			// log.Printf("[error] cannot send message to %v on network %v: %v", addr, networkID, err)
		}
	}()
}

func (network *network) query(q identity.Address) (identity.MultiAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	multiAddr, err := network.swarmer.Query(ctx, q, -1)
	if err != nil {
		return multiAddr, fmt.Errorf("cannot query peer %v: %v", q, err)
	}
	return multiAddr, nil
}

func (network *network) connectOrListen(ctx context.Context, networkID NetworkID, addr identity.Address) Sender {
	if addr < network.swarmer.MultiAddress().Address() {

		// Query for the multi-address
		log.Printf("[debug] querying peer %v on network %v", addr, networkID)
		multiAddr, err := network.query(addr)
		if err != nil {
			log.Printf("[error] cannot connect to peer %v on network %v: %v", addr, networkID, err)
			if addr < network.swarmer.MultiAddress().Address() {
				return nil
			}
		}

		// Connect to the remote server
		log.Printf("[debug] connecting to peer %v on network %v", addr, networkID)
		sender, err := network.conn.Connect(ctx, networkID, multiAddr, network.receiver)
		if err != nil {
			log.Printf("[error] cannot connect to peer %v on network %v: %v", addr, networkID, err)
			return nil
		}
		log.Printf("[debug] 🔗 connected to peer %v on network %v", addr, networkID)
		return sender
	}

	// Wait for the client to connect to us
	log.Printf("[debug] listening for peer %v on network %v", addr, networkID)
	sender, err := network.conn.Listen(ctx, networkID, addr, network.receiver)
	if err != nil {
		log.Printf("[error] cannot listen for peer %v on network %v: %v", addr, networkID, err)
		return nil
	}
	log.Printf("[debug] 🔗 accepted peer %v on network %v", addr, networkID)
	return sender
}
