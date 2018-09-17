package grpc

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/shamir"
	"golang.org/x/net/context"
)

// ErrOpenOrderRequestIsNil is returned when a gRPC request is nil or has nil
// fields.
var ErrOpenOrderRequestIsNil = errors.New("open order request is nil")

// ErrOrderFragmentIsNil is returned when the order fragment contains
// nil fields.
var ErrOrderFragmentIsNil = errors.New("order fragment is nil")

// ErrEncryptedCoExpShareIsNil is returned when the encrypted CoExp is nil.
var ErrEncryptedCoExpShareIsNil = errors.New("encrypted coexp is nil")

// ErrEncryptedOrderFragmentIsNil is returned when the encrypted order fragmetn
// is nil.
var ErrEncryptedOrderFragmentIsNil = errors.New("encrypted order fragment is nil")

type orderbookClient struct {
}

// NewOrderbookClient returns an implementation of the orderbook.Client
// interface that uses gRPC.
func NewOrderbookClient() orderbook.Client {
	return &orderbookClient{}
}

// OpenOrder implements the orderbook.Client interface.
func (client *orderbookClient) OpenOrder(ctx context.Context, multiAddr identity.MultiAddress, orderFragment order.EncryptedFragment) error {
	if orderFragment.IsNil() {
		return ErrOrderFragmentIsNil
	}
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
	if server == nil {
		logger.Network(logger.LevelError, "server is nil")
		return
	}
	RegisterOrderbookServiceServer(server.Server, service)
}

// OpenOrder implements the gRPC service for receiving EncryptedOrderFragments
// defined in protobuf.
func (service *OrderbookService) OpenOrder(ctx context.Context, request *OpenOrderRequest) (*OpenOrderResponse, error) {
	// Check for empty or invalid request fields.
	if request == nil {
		return nil, ErrOpenOrderRequestIsNil
	}
	if request.OrderFragment == nil {
		return nil, ErrOrderFragmentIsNil
	}

	fragment, err := unmarshalEncryptedOrderFragment(request.OrderFragment)
	if err != nil {
		return nil, err
	}
	return &OpenOrderResponse{}, service.server.OpenOrder(ctx, fragment)
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

		Blinding:    []byte(orderFragmentIn.Blinding),
		Commitments: marshalCommitments(orderFragmentIn.Commitments),
	}
}

func unmarshalEncryptedOrderFragment(orderFragmentIn *EncryptedOrderFragment) (order.EncryptedFragment, error) {
	if orderFragmentIn == nil {
		return order.EncryptedFragment{}, ErrEncryptedOrderFragmentIsNil
	}

	price, err := unmarshalEncryptedCoExpShare(orderFragmentIn.Price)
	if err != nil {
		return order.EncryptedFragment{}, err
	}
	volume, err := unmarshalEncryptedCoExpShare(orderFragmentIn.Volume)
	if err != nil {
		return order.EncryptedFragment{}, err
	}
	minVolume, err := unmarshalEncryptedCoExpShare(orderFragmentIn.MinimumVolume)
	if err != nil {
		return order.EncryptedFragment{}, err
	}

	orderFragment := order.EncryptedFragment{
		OrderType:       order.Type(orderFragmentIn.OrderType),
		OrderParity:     order.Parity(orderFragmentIn.OrderParity),
		OrderSettlement: order.Settlement(orderFragmentIn.OrderSettlement),
		OrderExpiry:     time.Unix(orderFragmentIn.OrderExpiry, 0),

		EpochDepth:    order.FragmentEpochDepth(orderFragmentIn.EpochDepth),
		Tokens:        orderFragmentIn.Tokens,
		Price:         price,
		Volume:        volume,
		MinimumVolume: minVolume,
		Nonce:         orderFragmentIn.Nonce,

		Blinding:    orderFragmentIn.Blinding,
		Commitments: unmarshalCommitments(orderFragmentIn.Commitments),
	}
	copy(orderFragment.OrderID[:], orderFragmentIn.OrderId)
	copy(orderFragment.ID[:], orderFragmentIn.Id)

	return orderFragment, nil
}

func marshalEncryptedCoExpShare(value order.EncryptedCoExpShare) *EncryptedCoExpShare {
	return &EncryptedCoExpShare{
		Co:  value.Co,
		Exp: value.Exp,
	}
}

func unmarshalEncryptedCoExpShare(value *EncryptedCoExpShare) (order.EncryptedCoExpShare, error) {
	if value == nil {
		return order.EncryptedCoExpShare{}, ErrEncryptedCoExpShareIsNil
	}
	return order.EncryptedCoExpShare{
		Co:  value.Co,
		Exp: value.Exp,
	}, nil
}

func marshalCommitments(values order.FragmentCommitments) map[uint64]*OrderFragmentCommitment {
	commitments := map[uint64]*OrderFragmentCommitment{}
	for i, value := range values {
		commitments[i] = &OrderFragmentCommitment{
			PriceCo:          value.PriceCo.Bytes(),
			PriceExp:         value.PriceExp.Bytes(),
			VolumeCo:         value.VolumeCo.Bytes(),
			VolumeExp:        value.VolumeExp.Bytes(),
			MinimumVolumeCo:  value.MinimumVolumeCo.Bytes(),
			MinimumVolumeExp: value.MinimumVolumeExp.Bytes(),
		}
	}
	return commitments
}

func unmarshalCommitments(values map[uint64]*OrderFragmentCommitment) order.FragmentCommitments {
	commitments := order.FragmentCommitments{}
	for i, value := range values {
		if value == nil {
			continue
		}
		commitments[i] = order.FragmentCommitment{
			PriceCo:          shamir.Commitment{Int: big.NewInt(0).SetBytes(value.PriceCo)},
			PriceExp:         shamir.Commitment{Int: big.NewInt(0).SetBytes(value.PriceExp)},
			VolumeCo:         shamir.Commitment{Int: big.NewInt(0).SetBytes(value.VolumeCo)},
			VolumeExp:        shamir.Commitment{Int: big.NewInt(0).SetBytes(value.VolumeExp)},
			MinimumVolumeCo:  shamir.Commitment{Int: big.NewInt(0).SetBytes(value.MinimumVolumeCo)},
			MinimumVolumeExp: shamir.Commitment{Int: big.NewInt(0).SetBytes(value.MinimumVolumeExp)},
		}
	}
	return commitments
}
