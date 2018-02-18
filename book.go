package node

import (
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-order-compute"
)

type HiddenOrderBook struct {
	do.GuardedObject

	comparisons map[string]map[string]*compute.OrderFragment
}
