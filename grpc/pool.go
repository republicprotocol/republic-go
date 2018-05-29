package grpc

import (
	"context"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/identity"
)

// ConnPoolEntry is an entry in the ConnPool that stores a timestamp alongside
// a connection. This timestamp is updated whenever the ConnPoolEntry is
// accessed and is used by the ConnPool to implement an LRU cache.
type ConnPoolEntry struct {
	Conn      *Conn
	Timestamp time.Time
}

// ConnPool is a pool of connections that can be reused when creating clients.
type ConnPool struct {
	cacheMu    *sync.Mutex
	cache      map[string]ConnPoolEntry
	cacheLimit int
}

// NewConnPool returns an empty ConnPool. The cache limit defines the maximum
// number of connections that the ConnPool will keep alive.
func NewConnPool(cacheLimit int) ConnPool {
	return ConnPool{
		cacheMu:    new(sync.Mutex),
		cache:      map[string]ConnPoolEntry{},
		cacheLimit: cacheLimit,
	}
}

// Dial creates a client connection to the given multiaddress (see Dial for
// creating a Conn) if the client connection does not already exist in the
// ConnPool. If the client connection already exists, it is cloned (see
// Conn.Clone) and returned. All calls to ConnPool.Dial must be accompanied by
// exactly one call to Conn.Close on the returned Conn.
func (pool *ConnPool) Dial(ctx context.Context, multiAddress identity.MultiAddress) (*Conn, error) {
	pool.cacheMu.Lock()
	defer pool.cacheMu.Unlock()

	// Check the cache for an existing connection and update the timestamp
	entry, ok := pool.cache[multiAddress.String()]
	if ok {
		entry.Timestamp = time.Now()

		// Clone the Conn so that closing it in the ConnPool will not
		// necessarily free resources if the Conn is still being used outside
		// the ConnPool
		return entry.Conn.Clone(), nil
	}

	// A new connection must be created so check if the cache has reached its
	// limit
	if len(pool.cache) >= pool.cacheLimit {
		var k string
		for multiAddress := range pool.cache {
			if k == "" || pool.cache[k].Timestamp.After(pool.cache[multiAddress].Timestamp) {
				k = multiAddress
			}
		}
		// Close the connection (this will safely decrement a reference counter
		// and only free resource if no other reference to the Conn is still
		// alive) and delete it from the cache
		pool.cache[k].Conn.Close()
		delete(pool.cache, k)
	}

	// Create the new connection and store it in the cache
	// FIXME: The context of the connection should be the combined cnotext of
	// all users... somehow...
	conn, err := Dial(context.Background(), multiAddress)
	if err != nil {
		return nil, err
	}
	pool.cache[multiAddress.String()] = ConnPoolEntry{
		Conn:      conn,
		Timestamp: time.Now(),
	}

	// Clone the Conn so that closing it in the ConnPool will not necessarily
	// free resources if the Conn is still being used outside the ConnPool
	return conn.Clone(), nil
}

// Close all connections in the ConnPool. Any Conn that is not referenced
// outside the ConnPool will be closed and all pending operations will be
// terminated (see Conn.Close). The first error that occurs will be returned,
// but it will not stop other connections from being closed.
func (pool *ConnPool) Close() error {
	pool.cacheMu.Lock()
	defer pool.cacheMu.Unlock()

	// Close all connections and store the first error
	var err error
	for multiAddress := range pool.cache {
		if innerErr := pool.cache[multiAddress].Conn.Close(); innerErr != nil && err == nil {
			err = innerErr
		}
	}

	// Clear the cache and return the first error
	pool.cache = map[string]ConnPoolEntry{}
	return err
}
