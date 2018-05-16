package smpc

import (
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/smpc/delta"
)

type Inst struct {
	*InstConnect
	*InstCompute
}

type InstConnect struct {
	Epoch []byte
	Peers []identity.Address
	N     int
	K     int
}

type InstCompute struct {
	Epoch []byte
	Buy   order.Fragment
	Sell  order.Fragment
}

type Result struct {
	*ResultCompute
}

type ResultCompute struct {
	Delta delta.Delta
	Err   error
}
