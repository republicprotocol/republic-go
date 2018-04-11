package rpc

import (
	"time"

	"github.com/republicprotocol/republic-go/identity"
)

// Options that parameterize the behavior of Nodes.
type Options struct {
	SwarmOptions  `json:"swarm"`
	SyncerOptions `json:"syncer"`

	MultiAddress      identity.MultiAddress `json:"multiAddress"`
	Timeout           time.Duration         `json:"timeout"`
	TimeoutBackoff    time.Duration         `json:"timeoutBackoff"`
	TimeoutRetries    int                   `json:"timeoutRetries"`
	MessageQueueLimit int                   `json:"messageQueueLimit"`
}
