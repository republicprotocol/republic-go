package rpc

import (
	"time"

	"github.com/republicprotocol/republic-go/identity"
)

// Options that parameterize the behavior of Nodes.
type Options struct {
	BootstrapMultiAddresses identity.MultiAddresses `json:"bootstrapMultiAddresses"`
	Address                 identity.Address        `json:"address"`
	MultiAddress            identity.MultiAddress   `json:"multiAddress"`
	Timeout                 time.Duration           `json:"timeout"`
	TimeoutBackoff          time.Duration           `json:"timeoutBackoff"`
	TimeoutRetries          int                     `json:"timeoutRetries"`
	MessageQueueLimit       int                     `json:"messageQueueLimit"`
	Concurrent              bool                    `json:"concurrent"`
	Alpha                   int                     `json:"alpha"`
	MaxBucketLength         int                     `json:"maxBucketLength"`
	ClientPoolCacheLimit    int                     `json:"clientPoolCacheLimit"`
	Debug                   DebugLevel              `json:"debug"`
}

type Provider struct {
	swarm Swarm
}

func NewProvider(swarm Swarm) {

}
