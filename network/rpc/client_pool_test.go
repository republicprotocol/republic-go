package rpc_test

import (
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/identity"
)

type ClientPool struct {
	do.GuardedObject

	from       identity.Multiaddress
	cache      map[string]*Client
	cacheLRU   []*Client
	cacheLimit int
}

func NewClientPool(from identity.Multiaddress, cacheLimit int) *ClientPool {
	pool := new(ClientPool)
	pool.GuardedObject = do.NewGuardedObject()
	pool.cacheLimit = cacheLimit
	return pool
}

func (pool *ClientPool) FindOrCreateClientConnection(to identity.Multiaddress) (*Client, error) {
	pool.Enter(nil)
	defer pool.Exit()
	return pool.findOrCreateClientConnection(to)
}

func (pool *ClientPool) findOrCreateClientConnection(to identity.Multiaddress) (*Client, error) {
	client, ok := pool.cache[to.String()]
	if !ok {
		client, err := NewClient(to, from)
		if err != nil {
			return client, err
		}
	}
	return client, nil
}
