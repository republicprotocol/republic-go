package rpc

import "time"

// ClientPoolOptions change the behavior of the ClientPool. By default, the
// ClientPool will use the DefaultClientPoolOptions. The ClientPool has methods
// that will configure it to use other options in a chaining style.
type ClientPoolOptions struct {
	Timeout        time.Duration `json:"timeout"`
	TimeoutBackoff time.Duration `json:"timeoutBackoff"`
	TimeoutRetries int           `json:"timeoutRetries"`
	CacheLimit     int           `json:"cacheLimit"`
}

// DefaultClientPoolOptions returns the ClientPoolOptions that the ClientPool
// uses when it is first created.
func DefaultClientPoolOptions() ClientPoolOptions {
	return ClientPoolOptions{
		Timeout:        30 * time.Second,
		TimeoutBackoff: 0 * time.Second,
		TimeoutRetries: 3,
		CacheLimit:     100,
	}
}

// WithTimeout returns a ClientPool that uses the given timeout for all
// Clients.
func (pool *ClientPool) WithTimeout(timeout time.Duration) *ClientPool {
	pool.Enter(nil)
	defer pool.Exit()
	pool.options.Timeout = timeout
	return pool
}

// WithTimeoutBackoff returns a ClientPool that uses the given timeout backoff
// to increment the timeout after failed RPCs for all Clients.
func (pool *ClientPool) WithTimeoutBackoff(backoff time.Duration) *ClientPool {
	pool.Enter(nil)
	defer pool.Exit()
	pool.options.TimeoutBackoff = backoff
	return pool
}

// WithTimeoutRetries returns a ClientPool that uses the given timeout retries
// to determine how many times it will attempt an RPC before failing proper for
// all Clients.
func (pool *ClientPool) WithTimeoutRetries(retries int) *ClientPool {
	pool.Enter(nil)
	defer pool.Exit()
	pool.options.TimeoutRetries = retries
	return pool
}

// WithCacheLimit returns a ClientPool that uses the given cache limit to
// determine how many Clients it will maintain in the cache.
func (pool *ClientPool) WithCacheLimit(cacheLimit int) *ClientPool {
	pool.Enter(nil)
	defer pool.Exit()
	pool.options.CacheLimit = cacheLimit
	return pool
}

// ClientOptions change the behavior of the Client. By default, the Client will
// use the DefaultClientOptions. The Client has methods that will configure it
// to use other options in a chaining style.
type ClientOptions struct {
	Timeout        time.Duration `json:"timeout"`
	TimeoutBackoff time.Duration `json:"timeoutBackoff"`
	TimeoutRetries int           `json:"timeoutRetries"`
}

// DefaultClientOptions returns the ClientOptions that the Client uses when it
// is first created.
func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		Timeout:        30 * time.Second,
		TimeoutBackoff: 0 * time.Second,
		TimeoutRetries: 3,
	}
}

// WithTimeout returns a Client that uses the given timeout for all RPCs.
func (client *Client) WithTimeout(timeout time.Duration) *Client {
	client.Options.Timeout = timeout
	return client
}

// WithTimeoutBackoff returns a Client that uses the given timeout backoff to
// increment the timeout after a every failed RPC.
func (client *Client) WithTimeoutBackoff(backoff time.Duration) *Client {
	client.Options.TimeoutBackoff = backoff
	return client
}

// WithTimeoutRetries returns a Client that uses the given timeout retries to
// determine how many times it will attempt every RPC before failing proper.
func (client *Client) WithTimeoutRetries(retries int) *Client {
	client.Options.TimeoutRetries = retries
	return client
}
