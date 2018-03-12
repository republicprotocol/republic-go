package dark

import (
	"bytes"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/identity"
)

type Node struct {
	identity.ID
	*identity.MultiAddress
}

type Nodes []Node

type PoolID []identity.ID

type Pools []*Pool

// A Pool is a list of nodes, identified by their Multiaddresses.
type Pool struct {
	do.GuardedObject
	nodes Nodes
}

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

// CoForAll loop over all nodes. The applied function must not modify the Pool.
func (pool *Pool) CoForAll(f func(node *Node)) {
	pool.EnterReadOnly(nil)
	defer pool.ExitReadOnly()
	do.CoForAll(pool.nodes, func(i int) {
		f(&pool.nodes[i])
	})
}

// ForAll loop over all nodes. The applied function must not modify the Pool.
func (pool *Pool) ForAll(f func(node *Node)) {
	pool.EnterReadOnly(nil)
	defer pool.ExitReadOnly()
	do.ForAll(pool.nodes, func(i int) {
		f(&pool.nodes[i])
	})
}
