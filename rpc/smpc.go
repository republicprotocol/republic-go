package rpc

import (
	"context"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
)

type Smpc interface {
}

type SmpcClient struct {
	smpc Smpc
}

func (client *SmpcClient) OpenOrder(ctx context.Context, multiAddr identity.MultiAddress, orderFragment order.Fragment) error {
	panic("unimplemented")
}
