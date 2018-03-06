package swarm_test

import (
	"time"

	"github.com/republicprotocol/go-swarm-network"
)

const (
	DefaultOptionsDebug           = swarm.DebugOff
	DefaultOptionsAlpha           = 3
	DefaultOptionsMaxBucketLength = 10
	DefaultOptionsTimeout         = 30 * time.Second
	DefaultOptionsTimeoutStep     = 30 * time.Second
	DefaultOptionsTimeoutRetries  = 1
	DefaultOptionsConcurrent      = false
)
