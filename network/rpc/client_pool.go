package rpc

import (
	"time"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/identity"
)

// A ClientCacheEntry is pointer to a Client that is stored in a cache. It is
// coupled with the timestamp at which the Client was last accessed from the
// cach.
type ClientCacheEntry struct {
	*Client
	Timestamp time.Time
}

// A ClientPool maintains a concurrency-safe pool of Clients. It caches all
// Clients and when it has reached its limit, it will remove the least recently
// used Client.
type ClientPool struct {
	do.GuardedObject

	from    identity.MultiAddress
	cache   map[string]ClientCacheEntry
	options ClientPoolOptions
}

// NewClientPool returns a new ClientPool with the given cache limit. All
// Clients that are created in this pool will identify themselves using the
// given MultiAddress.
func NewClientPool(from identity.MultiAddress) *ClientPool {
	pool := new(ClientPool)
	pool.GuardedObject = do.NewGuardedObject()
	pool.from = from
	pool.cache = map[string]ClientCacheEntry{}
	pool.options = DefaultClientPoolOptions()
	return pool
}

// FindOrCreateClient will return a Client that is connected to the given
// MultiAddress. It will first try to find an existing Client in the cache. If
// it cannot find one in the cache, it will create a new one and add it to the
// cache.
func (pool *ClientPool) FindOrCreateClient(to identity.MultiAddress) (*Client, error) {
	pool.Enter(nil)
	defer pool.Exit()
	return pool.findOrCreateClient(to)
}

func (pool *ClientPool) findOrCreateClient(to identity.MultiAddress) (*Client, error) {
	clientCacheEntry, ok := pool.cache[to.String()]
	if ok {
		clientCacheEntry.Timestamp = time.Now()
		return clientCacheEntry.Client, nil
	}

	client, err := NewClient(to, pool.from)
	if err != nil {
		return client, err
	}

	if len(pool.cache) >= pool.options.CacheLimit {
		var k string
		for multiAddress := range pool.cache {
			if k == "" || pool.cache[multiAddress].Timestamp.After(pool.cache[multiAddress].Timestamp) {
				k = multiAddress
			}
		}
		delete(pool.cache, k)
	}

	pool.cache[to.String()] = ClientCacheEntry{Client: client, Timestamp: time.Now()}
	return client, nil
}
