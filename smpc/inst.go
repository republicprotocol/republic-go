package smpc

import (
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/smpc/delta"
)

type Inst struct {
	*InstConnect
	*InstCompute
	*InstJoin
}

type InstConnect struct {
	PeersID []byte
	Peers   []identity.Address
	N       int
	K       int
}

type InstCompute struct {
	PeersID []byte
	Buy     order.Fragment
	Sell    order.Fragment
}

type InstJoin struct {
	PeersID []byte
	Buy     order.ID
	Sell    order.ID
}

type Result struct {
	*ResultCompute
	*ResultJoin
}

type ResultCompute struct {
	Delta delta.Delta
	Err   error
}

type ResultJoin struct {
	Buy  order.Order
	Sell order.Order
	Err  error
}
