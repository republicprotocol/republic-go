package rpc

import (
	"context"
	"errors"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
)

var ErrUnexpectedMessageType = errors.New("unexpected message type")

type SmpcClient interface {
	OpenOrder(ctx context.Context, multiAddr identity.MultiAddress, orderFragment order.Fragment) error
}
