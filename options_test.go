package xing_test

import (
	"time"

	"github.com/republicprotocol/go-xing"
)

const (
	DefaultOptionsDebug           = xing.DebugOff
	DefaultOptionsAlpha           = 3
	DefaultOptionsMaxBucketLength = 20
	DefaultOptionsTimeout         = 30 * time.Second
	DefaultOptionsTimeoutStep     = 30 * time.Second
	DefaultOptionsTimeoutRetries  = 1
	DefaultOptionsConcurrent      = false
)
