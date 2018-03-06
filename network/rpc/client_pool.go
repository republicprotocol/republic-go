package rpc

import (
	"time"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/identity"
)

type ClientCacheEntry struct {
	*Client
	Timestamp time.Time
}

type ClientPool struct {
	do.GuardedObject

	from       identity.Multiaddress
	cache      map[string]ClientCacheEntry
	cacheLimit int
}

func NewClientPool(from identity.Multiaddress, cacheLimit int) *ClientPool {
	pool := new(ClientPool)
	pool.GuardedObject = do.NewGuardedObject()
	pool.from = from
	pool.cache = map[string]ClientCacheEntry{}
	pool.cacheLimit = cacheLimit
	return pool
}

func (pool *ClientPool) FindOrCreateClientConnection(to identity.Multiaddress) (*Client, error) {
	pool.Enter(nil)
	defer pool.Exit()
	return pool.findOrCreateClientConnection(to)
}

func (pool *ClientPool) findOrCreateClientConnection(to identity.Multiaddress) (*Client, error) {
	clientCacheEntry, ok := pool.cache[to.String()]
	if ok {
		clientCacheEntry.Timestamp = time.Now()
		return clientCacheEntry.Client, nil
	}

	client, err := NewClient(to, from)
	if err != nil {
		return client, err
	}

	if len(pool.cache) >= pool.cacheLimit {
		var k string
		for multiaddress := range pool.cache {
			if k == "" || pool.cache[multiaddress].Timestamp.After(pool.cache[multiaddress]) {
				k = multiaddress
			}
		}
		delete(pool.cache, k)
	}

	pool.cache[to.String()] = ClientCacheEntry{Client: client, Timestamp: time.Now()}
	return client, nil
}
