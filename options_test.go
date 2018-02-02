package network_test

import (
	"time"

	"github.com/republicprotocol/go-network"
)

const (
	DefaultOptionsDebug           = network.DebugOff
	DefaultOptionsAlpha           = 3
	DefaultOptionsMaxBucketLength = 20
	DefaultOptionsTimeout         = 2 * time.Second
	DefaultOptionsTimeoutStep     = 0 * time.Second
	DefaultOptionsTimeoutRetries  = 1
	DefaultOptionsConcurrent      = false
)
