package grpc

import (
	"fmt"
	"time"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"golang.org/x/net/context"
)

type orderbookClient struct {
}

// NewOrderbookClient returns an implementation of the orderbook.Client
// interface that uses gRPC.
func NewOrderbookClient() orderbook.Client {
	return &orderbookClient{}
}

// OpenOrder implements the orderbook.Client interface.
func (client *orderbookClient) OpenOrder(ctx context.Context, multiAddr identity.MultiAddress, orderFragment order.EncryptedFragment) error {
	conn, err := Dial(ctx, multiAddr)
	if err != nil {
		return fmt.Errorf("cannot dial %v: %v", multiAddr, err)
	}
	defer conn.Close()

	request := &OpenOrderRequest{
		OrderFragment: marshalEncryptedOrderFragment(orderFragment),
	}
	return Backoff(ctx, func() error {
		_, err := NewOrderbookServiceClient(conn).OpenOrder(ctx, request)
		return err
	})
}

// OrderbookService is a Service that implements the gRPC OrderbookService
// defined in protobuf. It exposes an RPC that accepts OpenOrderRequests and
// delegates control to an orderbook.Server.
type OrderbookService struct {
	server orderbook.Server
}

// NewOrderbookService returns a gRPC service that unmarshals OpenOrderRequests
// defined in protobuf, and delegates control of the RPC to an
// orderbook.Server.
func NewOrderbookService(server orderbook.Server) OrderbookService {
	return OrderbookService{
		server: server,
	}
}

// Register implements the Service interface.
func (service *OrderbookService) Register(server *Server) {
	RegisterOrderbookServiceServer(server.Server, service)
}

// OpenOrder implements the gRPC service for receiving EncryptedOrderFragments
// defined in protobuf.
func (service *OrderbookService) OpenOrder(ctx context.Context, request *OpenOrderRequest) (*OpenOrderResponse, error) {
	return &OpenOrderResponse{}, service.server.OpenOrder(ctx, unmarshalEncryptedOrderFragment(request.OrderFragment))
}

func marshalEncryptedOrderFragment(orderFragmentIn order.EncryptedFragment) *EncryptedOrderFragment {
	return &EncryptedOrderFragment{
		OrderId:         orderFragmentIn.OrderID[:],
		OrderType:       OrderType(orderFragmentIn.OrderType),
		OrderParity:     OrderParity(orderFragmentIn.OrderParity),
		OrderSettlement: OrderSettlement(orderFragmentIn.OrderSettlement),
		OrderExpiry:     orderFragmentIn.OrderExpiry.Unix(),

		Id:            orderFragmentIn.ID[:],
		EpochDepth:    int32(orderFragmentIn.EpochDepth),
		Tokens:        orderFragmentIn.Tokens,
		Price:         marshalEncryptedCoExpShare(orderFragmentIn.Price),
		Volume:        marshalEncryptedCoExpShare(orderFragmentIn.Volume),
		MinimumVolume: marshalEncryptedCoExpShare(orderFragmentIn.MinimumVolume),
		Nonce:         orderFragmentIn.Nonce,
	}
}

func unmarshalEncryptedOrderFragment(orderFragmentIn *EncryptedOrderFragment) order.EncryptedFragment {
	orderFragment := order.EncryptedFragment{
		OrderType:       order.Type(orderFragmentIn.OrderType),
		OrderParity:     order.Parity(orderFragmentIn.OrderParity),
		OrderSettlement: order.Settlement(orderFragmentIn.OrderSettlement),
		OrderExpiry:     time.Unix(orderFragmentIn.OrderExpiry, 0),

		EpochDepth:    order.FragmentEpochDepth(orderFragmentIn.EpochDepth),
		Tokens:        orderFragmentIn.Tokens,
		Price:         unmarshalEncryptedCoExpShare(orderFragmentIn.Price),
		Volume:        unmarshalEncryptedCoExpShare(orderFragmentIn.Volume),
		MinimumVolume: unmarshalEncryptedCoExpShare(orderFragmentIn.MinimumVolume),
		Nonce:         orderFragmentIn.Nonce,
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
