package rpc

import (
	"context"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
	"google.golang.org/grpc"
)

type RelayDelegate interface {
	OnOpenOrder(from identity.MultiAddress, orderFragment *order.Fragment)
}

type RelayService struct {
	MultiAddress      identity.MultiAddress
	Logger            *logger.Logger
	Delegate          RelayDelegate
	MessageQueueLimit int
}

func NewRelayService(multiAddress identity.MultiAddress, logger *logger.Logger, delegate RelayDelegate, messageQueueLimit int) *RelayService {
	return &RelayService{
		MultiAddress:      multiAddress,
		Logger:            logger,
		Delegate:          delegate,
		MessageQueueLimit: messageQueueLimit,
	}
}

// Register the RelayService with a gRPC server.
func (service *RelayService) Register(server *grpc.Server) {
	RegisterRelayServer(server, service)
}

func (service *RelayService) OpenOrder(ctx context.Context, req *OpenOrderRequest) (*Nothing, error) {
	wait := do.Process(func() do.Option {
		return do.Err(service.openOrder(req))
	})

	select {
	case val := <-wait:
		return &Nothing{}, val.Err
	case <-ctx.Done():
		return &Nothing{}, ctx.Err()
	}
}

func (service *RelayService) openOrder(req *OpenOrderRequest) error {
	multi , err := identity.NewMultiAddressFromString(req.From.MultiAddress)
	if err != nil {
		return err
	}

	fragment, err := UnmarshalOrderFragment(req.OrderFragment)
	if err != nil {
		return err
	}
	service.Delegate.OnOpenOrder(multi, fragment)
	return nil
}

func (service *RelayService) SignOrderFragment(ctx context.Context, id *OrderFragmentId) (*OrderFragmentId, error) {
	wait := do.Process(func() do.Option {
		return do.Err(service.signOrderFragment(id))
	})

	select {
	case val := <-wait:
		return &OrderFragmentId{}, val.Err
	case <-ctx.Done():
		return &OrderFragmentId{}, ctx.Err()
	}
}

func (service *RelayService) signOrderFragment(id *OrderFragmentId) error {
	return nil
}

func (service *RelayService) CancelOrder(ctx context.Context, req *CancelOrderRequest) (*Nothing, error) {
	wait := do.Process(func() do.Option {
		return do.Err(service.cancelOrder(req))
	})

	select {
	case val := <-wait:
		return &Nothing{}, val.Err
	case <-ctx.Done():
		return &Nothing{}, ctx.Err()
	}
}

func (service *RelayService) cancelOrder(req *CancelOrderRequest) error {
	return nil
}
