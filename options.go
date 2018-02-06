package network

import (
	"time"

	"github.com/republicprotocol/go-identity"
)

// Constants for different options.
const (
	DebugOff    = 0
	DebugLow    = 1
	DebugMedium = 2
	DebugHigh   = 3
)

// Options that parameterize the behavior of Nodes.
type Options struct {
	Host                    string
	Port                    string
	MultiAddress            identity.MultiAddress
	BootstrapMultiAddresses identity.MultiAddresses

	Debug           int
	Alpha           int
	MaxBucketLength int
	Timeout         time.Duration
	TimeoutStep     time.Duration
	TimeoutRetries  int
	Concurrent      bool
}
