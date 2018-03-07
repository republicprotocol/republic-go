package network

import (
	"time"

	"github.com/republicprotocol/republic-go/identity"
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
	MultiAddress            identity.MultiAddress   `json:"multiAddress"`
	BootstrapMultiAddresses identity.MultiAddresses `json:"bootstrapMultiAddresses"`

	Debug                int           `json:"debug"`
	Alpha                int           `json:"alpha"`
	MaxBucketLength      int           `json:"maxBucketLength"`
	ClientPoolCacheLimit int           `json:"clientPoolCacheLimit"`
	Timeout              time.Duration `json:"timeout"`
	TimeoutBackoff       time.Duration `json:"timeoutBackoff"`
	TimeoutRetries       int           `json:"timeoutRetries"`
	Concurrent           bool
}
