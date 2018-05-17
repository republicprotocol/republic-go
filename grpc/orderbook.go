package grpc

import (
	"fmt"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"golang.org/x/net/context"
)

type orderbookClient struct {
	connPool *ConnPool
}

// NewOrderbookClient returns an implementation of the orderbook.Client
// interface that uses gRPC and a recycled connection pool.
func NewOrderbookClient(connPool *ConnPool) orderbook.Client {
	return &orderbookClient{
		connPool: connPool,
	}
}

// OpenOrder implements the orderbook.Client interface.
func (client *orderbookClient) OpenOrder(ctx context.Context, multiAddr identity.MultiAddress, orderFragmentIn order.EncryptedFragment) error {
	conn, err := client.connPool.Dial(ctx, multiAddr)
	if err != nil {
		return fmt.Errorf("cannot dial %v: %v", multiAddr, err)
	}
	defer conn.Close()

	orderbookServiceClient := NewOrderbookServiceClient(conn.ClientConn)
	_, err = orderbookServiceClient.OpenOrder(ctx, adaptEncryptedOrderFragment(orderFragmentIn))
	return err
}

func adaptEncryptedOrderFragment(orderFragmentIn order.EncryptedFragment) *EncryptedOrderFragment {
	return &EncryptedOrderFragment{
		OrderId:     orderFragmentIn.OrderID[:],
		OrderType:   OrderType(orderFragmentIn.OrderType),
		OrderParity: OrderParity(orderFragmentIn.OrderParity),
		OrderExpiry: orderFragmentIn.OrderExpiry.Unix(),

		Id:            orderFragmentIn.ID[:],
		Tokens:        orderFragmentIn.Tokens,
		Price:         adapterEncryptedCoExpShare(orderFragmentIn.Price),
		Volume:        adapterEncryptedCoExpShare(orderFragmentIn.Volume),
		MinimumVolume: adapterEncryptedCoExpShare(orderFragmentIn.MinimumVolume),
	}
}

func adapterEncryptedCoExpShare(value order.EncryptedCoExpShare) *EncryptedCoExpShare {
	return &EncryptedCoExpShare{
		Co:  value.Co,
		Exp: value.Exp,
	}
}
