package grpc

import (
	"fmt"
	"log"
	"time"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"golang.org/x/net/context"
)

type OrderbookService struct {
	server orderbook.Server
}

// NewOrderbookService returns a gRPC service that unmarshals
// EncryptedOrderFragments and delegates control to an orderbook.Service.
func NewOrderbookService(server orderbook.Server) OrderbookService {
	return OrderbookService{
		server: server,
	}
}

// Register the OrderbookService to a Server.
func (service *OrderbookService) Register(server *Server) {
	RegisterOrderbookServiceServer(server.Server, service)
}

// OpenOrder implements the gRPC service for receiving EncryptedOrderFragments.
func (service *OrderbookService) OpenOrder(ctx context.Context, request *OpenOrderRequest) (*OpenOrderResponse, error) {
	log.Printf("opening order using order fragment...")
	return &OpenOrderResponse{}, service.server.OpenOrder(ctx, unmarshalEncryptedOrderFragment(request.OrderFragment))
}

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
func (client *orderbookClient) OpenOrder(ctx context.Context, multiAddr identity.MultiAddress, orderFragment order.EncryptedFragment) error {
	conn, err := client.connPool.Dial(ctx, multiAddr)
	if err != nil {
		return fmt.Errorf("cannot dial %v: %v", multiAddr, err)
	}
	defer conn.Close()

	request := &OpenOrderRequest{
		OrderFragment: marshalEncryptedOrderFragment(orderFragment),
	}
	_, err = NewOrderbookServiceClient(conn.ClientConn).OpenOrder(ctx, request)

	return err
}

func marshalEncryptedOrderFragment(orderFragmentIn order.EncryptedFragment) *EncryptedOrderFragment {
	return &EncryptedOrderFragment{
		OrderId:     orderFragmentIn.OrderID[:],
		OrderType:   OrderType(orderFragmentIn.OrderType),
		OrderParity: OrderParity(orderFragmentIn.OrderParity),
		OrderExpiry: orderFragmentIn.OrderExpiry.Unix(),

		Id:            orderFragmentIn.ID[:],
		Tokens:        orderFragmentIn.Tokens,
		Price:         marshalEncryptedCoExpShare(orderFragmentIn.Price),
		Volume:        marshalEncryptedCoExpShare(orderFragmentIn.Volume),
		MinimumVolume: marshalEncryptedCoExpShare(orderFragmentIn.MinimumVolume),
	}
}

func unmarshalEncryptedOrderFragment(orderFragmentIn *EncryptedOrderFragment) order.EncryptedFragment {
	orderFragment := order.EncryptedFragment{
		OrderType:   order.Type(orderFragmentIn.OrderType),
		OrderParity: order.Parity(orderFragmentIn.OrderParity),
		OrderExpiry: time.Unix(orderFragmentIn.OrderExpiry, 0),

		Tokens:        orderFragmentIn.Tokens,
		Price:         unmarshalEncryptedCoExpShare(orderFragmentIn.Price),
		Volume:        unmarshalEncryptedCoExpShare(orderFragmentIn.Volume),
		MinimumVolume: unmarshalEncryptedCoExpShare(orderFragmentIn.MinimumVolume),
	}
	copy(orderFragment.OrderID[:], orderFragmentIn.OrderId)
	copy(orderFragment.ID[:], orderFragmentIn.Id)
	return orderFragment
}

func marshalEncryptedCoExpShare(value order.EncryptedCoExpShare) *EncryptedCoExpShare {
	return &EncryptedCoExpShare{
		Co:  value.Co,
		Exp: value.Exp,
	}
}

func unmarshalEncryptedCoExpShare(value *EncryptedCoExpShare) order.EncryptedCoExpShare {
	return order.EncryptedCoExpShare{
		Co:  value.Co,
		Exp: value.Exp,
	}
}
