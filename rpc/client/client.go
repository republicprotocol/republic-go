// Package client implements common functions used by the difference RPC
// service clients. It implements an atomically reference counted connection
// pool for safe concurrent gRPC connection reuse.
package client
