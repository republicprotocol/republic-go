package rpc

import "time"

// ClientOptions change the behavior of the Client. By default, the Client will
// use the DefaultClientOptions. The Client has methods that will configure it
// to use other options in a chaining style.
type ClientOptions struct {
	Timeout        time.Duration
	TimeoutBackoff time.Duration
	TimeoutRetries int
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
