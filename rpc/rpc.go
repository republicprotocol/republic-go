package rpc

import (
	"time"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
)

// Options that parameterize the behavior of Nodes.
type Options struct {
	SwarmOptions `json:"swarm"`

	Logger            *logger.Logger
	MultiAddress      identity.MultiAddress `json:"multiAddress"`
	Timeout           time.Duration         `json:"timeout"`
	TimeoutBackoff    time.Duration         `json:"timeoutBackoff"`
	TimeoutRetries    int                   `json:"timeoutRetries"`
	MessageQueueLimit int                   `json:"messageQueueLimit"`
}
