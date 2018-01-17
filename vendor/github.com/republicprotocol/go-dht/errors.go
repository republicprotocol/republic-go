package dht

import "fmt"

// ErrFullBucket is used when a peer is inserted into a Bucket that already has
// the maximum number of Entries.
var (
	ErrFullBucket = fmt.Errorf("cannot add entry to a full bucket")
	ErrDHTAddress = fmt.Errorf("cannot add entry for the DHT address")
)
