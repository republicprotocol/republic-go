package rpc

import (
	"context"
	"time"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/identity"
)

// A ClientCacheEntry is pointer to a Client that is stored in a cache. It is
// coupled with the timestamp at which the Client was last accessed from the
// cache.
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

	client, err := NewClient(context.Background(), to, pool.from)
	if err != nil {
		return client, err
	}
	client = client.
		WithTimeout(pool.options.Timeout).
		WithTimeoutBackoff(pool.options.TimeoutBackoff).
		WithTimeoutRetries(pool.options.TimeoutRetries)

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
func (pool *ClientPool) Ping(ctx context.Context, to identity.MultiAddress) error {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		return err
	}

	return client.Ping(ctx)
}

// QueryPeers RPC.
func (pool *ClientPool) QueryPeers(ctx context.Context, to identity.MultiAddress, target *Address) (<-chan *MultiAddress, <-chan error) {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		errCh := make(chan error, 1)
		errCh <- err
		return nil, errCh
	}

	return client.QueryPeers(ctx, target)
}

// QueryPeersDeep RPC.
func (pool *ClientPool) QueryPeersDeep(ctx context.Context, to identity.MultiAddress, target *Address) (<-chan *MultiAddress, <-chan error) {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		errCh := make(chan error, 1)
		errCh <- err
		return nil, errCh
	}

	return client.QueryPeersDeep(ctx, target)
}

// Sync RPC.
func (pool *ClientPool) Sync(ctx context.Context, to identity.MultiAddress) (<-chan *SyncBlock, <-chan error) {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		errCh := make(chan error, 1)
		defer close(errCh)
		errCh <- err
		return nil, errCh
	}

	return client.Sync(ctx)
}

// SignOrderFragment RPC.
func (pool *ClientPool) SignOrderFragment(ctx context.Context, to identity.MultiAddress, orderFragmentId *OrderFragmentId) (*OrderFragmentId, error) {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		return nil, err
	}

	return client.SignOrderFragment(ctx, orderFragmentId)
}

// OpenOrder RPC.
func (pool *ClientPool) OpenOrder(ctx context.Context, to identity.MultiAddress, openOrderRequest *OpenOrderRequest) error {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		return err
	}

	return client.OpenOrder(ctx, openOrderRequest)
}

// CancelOrder RPC.
func (pool *ClientPool) CancelOrder(ctx context.Context, to identity.MultiAddress, cancelOrderRequest *CancelOrderRequest) error {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		return err
	}

	return client.CancelOrder(ctx, cancelOrderRequest)
}

// Compute RPC.
func (pool *ClientPool) Compute(ctx context.Context, to identity.MultiAddress, computationChIn <-chan *Computation) (<-chan *Computation, <-chan error) {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		errCh := make(chan error, 1)
		defer close(errCh)
		errCh <- err
		return nil, errCh
	}

	return client.Compute(ctx, computationChIn)
}

// OrderMatch RPC.
func (pool *ClientPool) SendTx(ctx context.Context, to identity.MultiAddress, tx *Tx) error {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		return err
	}

	return client.SendTx(ctx, tx)
}

// FinalizedBlock RPC.
func (pool *ClientPool) SyncBlocks(ctx context.Context, to identity.MultiAddress) (<-chan *Block, <-chan error) {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		errCh := make(chan error, 1)
		defer close(errCh)
		errCh <- err
		return nil, errCh
	}

	return client.SyncBlock(ctx)
}

// Drive RPC.
func (pool *ClientPool) Drive(ctx context.Context, to identity.MultiAddress, driveMessages <-chan *DriveMessage) (<-chan *DriveMessage, <-chan error) {
	client, err := pool.FindOrCreateClient(to)
	if err != nil {
		errCh := make(chan error, 1)
		defer close(errCh)
		errCh <- err
		return nil, errCh
	}

	return client.Drive(ctx, driveMessages)
}
