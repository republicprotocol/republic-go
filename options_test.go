package network_test

import (
	"time"

	"github.com/republicprotocol/go-network"
)

const (
	DefaultOptionsDebug           = network.DebugOff
	DefaultOptionsAlpha           = 3
	DefaultOptionsMaxBucketLength = 20
	DefaultOptionsTimeout         = 3 * time.Second
	DefaultOptionsTimeoutStep     = 3 * time.Second
	DefaultOptionsConcurrent      = false
)
