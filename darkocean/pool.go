package darkocean

import (
	"bytes"
	"sync"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/identity"
)

// A Node represents a dark node in a dark ocean pool
type Node struct {
	identity.ID
	mu           *sync.RWMutex
	multiAddress *identity.MultiAddress
}

// NewNode returns a new Node representing the dark pool with the given ID
func NewNode(id identity.ID) Node {
	return Node{
		ID:           id,
		mu:           new(sync.RWMutex),
		multiAddress: nil,
	}
}

// SetMultiAddress sets the multiaddress of the Node
func (node *Node) SetMultiAddress(multiAddress identity.MultiAddress) {
	node.mu.Lock()
	defer node.mu.Unlock()
	node.multiAddress = &multiAddress
}

// MultiAddress returns the multiaddress of the Node
func (node *Node) MultiAddress() *identity.MultiAddress {
	node.mu.RLock()
	defer node.mu.RUnlock()
	return node.multiAddress
}

// Nodes stores a list of Node structs
type Nodes []Node

// A Pool is a list of nodes, identified by their Multiaddresses.
type Pool struct {
	do.GuardedObject
	nodes Nodes
}

// Pools stores a list of Pool structs
type Pools []*Pool

// NewPool returns a new empty Pool.
func NewPool() *Pool {
	pool := new(Pool)
	pool.GuardedObject = do.NewGuardedObject()
	pool.nodes = Nodes{}
	return pool
}

// Append nodes to the DarkPool.
func (pool *Pool) Append(nodes ...Node) {
	pool.Enter(nil)
	defer pool.Exit()
	pool.nodes = append(pool.nodes, nodes...)
}

// Has returns a Node if an ID is held by the Pool, otherwise nil.
func (pool *Pool) Has(nodeID identity.ID) *Node {
	pool.EnterReadOnly(nil)
	defer pool.ExitReadOnly()
	for _, node := range pool.nodes {
		if bytes.Equal([]byte(nodeID), []byte(node.ID)) {
			return &node
		}
	}
	return nil
}

// Size returns the number of Nodes in the Pool.
func (pool *Pool) Size() int {
	pool.EnterReadOnly(nil)
	defer pool.ExitReadOnly()
	return len(pool.nodes)
}

// RemoveAll Nodes from the Pool.
func (pool *Pool) RemoveAll() {
	pool.Enter(nil)
	defer pool.Exit()
	pool.nodes = pool.nodes[:0]
}

// CoForAll loop over all Nodes. The applied function must not modify the Pool
// but it can modify the Node.
func (pool *Pool) CoForAll(f func(node *Node)) {
	pool.EnterReadOnly(nil)
	defer pool.ExitReadOnly()
	do.CoForAll(pool.nodes, func(i int) {
		f(&pool.nodes[i])
	})
}

// ForAll loop over all Nodes. The applied function must not modify the Pool
// but it can modify the Node.
func (pool *Pool) ForAll(f func(node *Node)) {
	pool.EnterReadOnly(nil)
	defer pool.ExitReadOnly()
	do.ForAll(pool.nodes, func(i int) {
		f(&pool.nodes[i])
	})
}

// For loop over all Nodes. The applied function must not modify the Pool
// but it can modify the Node.
func (pool *Pool) For(f func(node *Node)) {
	pool.EnterReadOnly(nil)
	defer pool.ExitReadOnly()
	for i := range pool.nodes {
		f(&pool.nodes[i])
	}
}
