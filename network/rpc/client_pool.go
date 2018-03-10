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

// Ping RPC.
func (pool *ClientPool) Ping(to identity.MultiAddress) error {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		return err
	}
	return client.Ping()
}

// QueryPeers RPC.
func (pool *ClientPool) QueryPeers(to identity.MultiAddress, target *Address) (chan *MultiAddress, error) {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		return nil, err
	}
	return client.QueryPeers(target)
}

// QueryPeersDeep RPC.
func (pool *ClientPool) QueryPeersDeep(to identity.MultiAddress, target *Address) (chan *MultiAddress, error) {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		return nil, err
	}
	return client.QueryPeersDeep(target)
}

// Sync RPC.
func (pool *ClientPool) Sync(to identity.MultiAddress) (chan *SyncBlock, error) {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		return nil, err
	}
	return client.Sync()
}

// SignOrderFragment RPC.
func (pool *ClientPool) SignOrderFragment(to identity.MultiAddress, orderFragmentSignature *OrderFragmentSignature) (*OrderFragmentSignature, error) {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		return nil, err
	}
	return client.SignOrderFragment(orderFragmentSignature)
}

// OpenOrder RPC.
func (pool *ClientPool) OpenOrder(to identity.MultiAddress, orderSignature *OrderSignature, orderFragment *OrderFragment) error {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		return err
	}
	return client.OpenOrder(orderSignature, orderFragment)
}

// CancelOrder RPC.
func (pool *ClientPool) CancelOrder(to identity.MultiAddress, orderSignature *OrderSignature) error {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		return err
	}
	return client.CancelOrder(orderSignature)
}

// RandomFragmentShares RPC.
func (pool *ClientPool) RandomFragmentShares(to identity.MultiAddress) (*RandomFragments, error) {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		return nil, err
	}
	return client.RandomFragmentShares()
}

// ResidueFragmentShares RPC.
func (pool *ClientPool) ResidueFragmentShares(to identity.MultiAddress, randomFragments *RandomFragments) (*ResidueFragments, error) {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		return nil, err
	}
	return client.ResidueFragmentShares(randomFragments)
}

// ComputeResidueFragment RPC.
func (pool *ClientPool) ComputeResidueFragment(to identity.MultiAddress, residueFragments *ResidueFragments) error {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		return err
	}
	return client.ComputeResidueFragment(residueFragments)
}

// BroadcastAlphaBetaFragment RPC.
func (pool *ClientPool) BroadcastAlphaBetaFragment(to identity.MultiAddress, alphaBetaFragment *AlphaBetaFragment) (*AlphaBetaFragment, error) {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		return nil, err
	}
	return client.BroadcastAlphaBetaFragment(alphaBetaFragment)
}

// BroadcastDeltaFragment RPC.
func (pool *ClientPool) BroadcastDeltaFragment(to identity.MultiAddress, deltaFragment *DeltaFragment) (*DeltaFragment, error) {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		return nil, err
	}
	return client.BroadcastDeltaFragment(deltaFragment)
}
