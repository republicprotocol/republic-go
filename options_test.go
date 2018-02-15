package dark_test

import (
	"time"

	"github.com/republicprotocol/go-dark-network"
)

const (
	DefaultOptionsDebug           = dark.DebugOff
	DefaultOptionsAlpha           = 3
	DefaultOptionsMaxBucketLength = 20
	DefaultOptionsTimeout         = 30 * time.Second
	DefaultOptionsTimeoutStep     = 30 * time.Second
	DefaultOptionsTimeoutRetries  = 1
	DefaultOptionsConcurrent      = false
)
