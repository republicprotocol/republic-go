package network_test

import (
	"time"

	"github.com/republicprotocol/go-network"
)

const (
	DefaultOptionsDebug           = network.DebugLow
	DefaultOptionsAlpha           = 3
	DefaultOptionsMaxBucketLength = 128
	DefaultOptionsTimeout         = 30 * time.Second
	DefaultOptionsTimeoutStep     = 30 * time.Second
	DefaultOptionsTimeoutRetries  = 1
	DefaultOptionsConcurrent      = false
)
