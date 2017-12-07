package dht

import "container/list"

// IDLength is the number of bytes needed to store an ID.
const (
	IDLength       = 30
	IDLengthInBits = IDLength * 8
)

// ID is the public address used to identify Nodes, and other entities, in the
// overlay network. It is generated from the public key of a key pair.
type ID string

// RoutingTable is a k-bucket routing table, where each bucket is a list of
// multiaddress strings, identifying peers by their network address as well as
// their ID.
type RoutingTable struct {
	Buckets [IDLengthInBits]list.List
}
