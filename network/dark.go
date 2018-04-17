package network

import (
	"context"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/compute"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
	"google.golang.org/grpc"
)

// A DarkDelegate is used as a callback interface to inject behavior into the
// DarkService service.
type DarkDelegate interface {
	// OnSync(from identity.MultiAddress)

	// OnSignOrderFragment(from identity.MultiAddress)
	OnOpenOrder(from identity.MultiAddress, orderFragment *order.Fragment)
	// OnCancelOrder(from identity.MultiAddress)

	// OnRandomFragmentShares(from identity.MultiAddress)
	// OnResidueFragmentShares(from identity.MultiAddress)
	// OnComputeResidueFragment(from identity.MultiAddress)
	// OnBroadcastAlphaBetaFragment(from identity.MultiAddress)
	OnBroadcastDeltaFragment(from identity.MultiAddress, deltaFragment *compute.DeltaFragment)
}

// DarkService implements the gRPC Dark service.
type DarkService struct {
	DarkDelegate
	Options
	Logger *logger.Logger
}

// NewDarkService returns a new DarkService with the provided delegate, options and logger
func NewDarkService(delegate DarkDelegate, options Options, logger *logger.Logger) *DarkService {
	return &DarkService{
		DarkDelegate: delegate,
		Options:      options,
		Logger:       logger,
	}
}

// Register the gRPC service.
func (service *DarkService) Register(server *grpc.Server) {
	rpc.RegisterDarkServer(server, service)
}

// Sync handles an rpc.SyncRequest
func (service *DarkService) Sync(syncRequest *rpc.SyncRequest, stream rpc.Dark_SyncServer) error {
	wait := do.Process(func() do.Option {
		return do.Err(service.sync(syncRequest, stream))
	})

	select {
	case val := <-wait:
		return val.Err

	case <-stream.Context().Done():
		return stream.Context().Err()
	}
}

func (service *DarkService) sync(syncRequest *rpc.SyncRequest, stream rpc.Dark_SyncServer) error {
	// todo : unimplemented
	return nil
}

// SignOrderFragment handles an rpc.SignOrderFragmentRequest
func (service *DarkService) SignOrderFragment(ctx context.Context, signOrderFragmentRequest *rpc.SignOrderFragmentRequest) (*rpc.OrderFragmentSignature, error) {
	wait := do.Process(func() do.Option {
		orderFragmentSignature, err := service.signOrderFragment(signOrderFragmentRequest)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(orderFragmentSignature)
	})

	select {
	case val := <-wait:
		if val, ok := val.Ok.(*rpc.OrderFragmentSignature); ok {
			return val, nil
		}
		return &rpc.OrderFragmentSignature{}, val.Err
	case <-ctx.Done():
		return &rpc.OrderFragmentSignature{}, ctx.Err()
	}
}

func (service *DarkService) signOrderFragment(signOrderFragmentRequest *rpc.SignOrderFragmentRequest) (*rpc.OrderFragmentSignature, error) {
	// todo : unimplemented
	return &rpc.OrderFragmentSignature{}, nil
}

// OpenOrder handles an rpc.OpenOrderRequest
func (service *DarkService) OpenOrder(ctx context.Context, openOrderRequest *rpc.OpenOrderRequest) (*rpc.Nothing, error) {
	wait := do.Process(func() do.Option {
		nothing, err := service.openOrder(openOrderRequest)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(nothing)
	})

	select {
	case val := <-wait:
		if val, ok := val.Ok.(*rpc.Nothing); ok {
			return val, nil
		}
		return &rpc.Nothing{}, val.Err

	case <-ctx.Done():
		return &rpc.Nothing{}, ctx.Err()
	}
}

func (service *DarkService) openOrder(openOrderRequest *rpc.OpenOrderRequest) (*rpc.Nothing, error) {
	from, sig, err := rpc.DeserializeMultiAddress(openOrderRequest.From)
	if err != nil {
		return &rpc.Nothing{}, err
	}
	err = from.VerifySignature(sig)
	if err != nil {
		return &rpc.Nothing{}, err
	}
	orderFragment, err := rpc.DeserializeOrderFragment(openOrderRequest.OrderFragment)
	if err != nil {
		return &rpc.Nothing{}, err
	}
	// Verify fragment signature
	err = orderFragment.VerifySignature(from.ID())
	if err != nil {
		return &rpc.Nothing{}, err
	}
	service.OnOpenOrder(from, orderFragment)
	return &rpc.Nothing{}, nil
}

// CancelOrder handles an rpc.CancelOrderRequest
func (service *DarkService) CancelOrder(ctx context.Context, cancelOrderRequest *rpc.CancelOrderRequest) (*rpc.Nothing, error) {
	wait := do.Process(func() do.Option {
		nothing, err := service.cancelOrder(cancelOrderRequest)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(nothing)
	})

	select {
	case val := <-wait:
		if val, ok := val.Ok.(*rpc.Nothing); ok {
			return val, nil
		}
		return &rpc.Nothing{}, val.Err

	case <-ctx.Done():
		return &rpc.Nothing{}, ctx.Err()
	}
}

func (service *DarkService) cancelOrder(cancelOrderRequest *rpc.CancelOrderRequest) (*rpc.Nothing, error) {
	// todo : unimplemented
	return &rpc.Nothing{}, nil
}

// RandomFragmentShares handles an rpc.RandomFragmentSharesRequest
func (service *DarkService) RandomFragmentShares(ctx context.Context, randomFragmentSharesRequest *rpc.RandomFragmentSharesRequest) (*rpc.RandomFragments, error) {
	wait := do.Process(func() do.Option {
		nothing, err := service.randomFragmentShares(randomFragmentSharesRequest)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(nothing)
	})

	select {
	case val := <-wait:
		if val, ok := val.Ok.(*rpc.RandomFragments); ok {
			return val, nil
		}
		return &rpc.RandomFragments{}, val.Err

	case <-ctx.Done():
		return &rpc.RandomFragments{}, ctx.Err()
	}
}

func (service *DarkService) randomFragmentShares(randomFragmentSharesRequest *rpc.RandomFragmentSharesRequest) (*rpc.RandomFragments, error) {
	// todo : unimplemented
	return &rpc.RandomFragments{}, nil
}

// ResidueFragmentShares handles an rpc.ResidueFragmentSharesRequest
func (service *DarkService) ResidueFragmentShares(ctx context.Context, residueFragmentSharesRequest *rpc.ResidueFragmentSharesRequest) (*rpc.ResidueFragments, error) {
	wait := do.Process(func() do.Option {
		nothing, err := service.residueFragmentShares(residueFragmentSharesRequest)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(nothing)
	})

	select {
	case val := <-wait:
		if val, ok := val.Ok.(*rpc.ResidueFragments); ok {
			return val, nil
		}
		return &rpc.ResidueFragments{}, val.Err

	case <-ctx.Done():
		return &rpc.ResidueFragments{}, ctx.Err()
	}
}

func (service *DarkService) residueFragmentShares(residueFragmentSharesRequest *rpc.ResidueFragmentSharesRequest) (*rpc.ResidueFragments, error) {
	// todo : unimplemented
	return &rpc.ResidueFragments{}, nil
}

// ComputeResidueFragment handles  an rpc.ComputeResidueFragmentRequest
func (service *DarkService) ComputeResidueFragment(ctx context.Context, computeResidueFragmentRequest *rpc.ComputeResidueFragmentRequest) (*rpc.Nothing, error) {
	wait := do.Process(func() do.Option {
		nothing, err := service.computeResidueFragment(computeResidueFragmentRequest)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(nothing)
	})

	select {
	case val := <-wait:
		if val, ok := val.Ok.(*rpc.Nothing); ok {
			return val, nil
		}
		return &rpc.Nothing{}, val.Err

	case <-ctx.Done():
		return &rpc.Nothing{}, ctx.Err()
	}
}

func (service *DarkService) computeResidueFragment(computeResidueFragmentRequest *rpc.ComputeResidueFragmentRequest) (*rpc.Nothing, error) {
	// todo : unimplemented
	return &rpc.Nothing{}, nil
}

// BroadcastAlphaBetaFragment handles an rpc.BroadcastAlphaBetaFragmentRequest
func (service *DarkService) BroadcastAlphaBetaFragment(ctx context.Context, broadcastAlphaBetaFragmentRequest *rpc.BroadcastAlphaBetaFragmentRequest) (*rpc.AlphaBetaFragment, error) {
	wait := do.Process(func() do.Option {
		nothing, err := service.broadcastAlphaBetaFragment(broadcastAlphaBetaFragmentRequest)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(nothing)
	})

	select {
	case val := <-wait:
		if val, ok := val.Ok.(*rpc.AlphaBetaFragment); ok {
			return val, nil
		}
		return &rpc.AlphaBetaFragment{}, val.Err

	case <-ctx.Done():
		return &rpc.AlphaBetaFragment{}, ctx.Err()
	}
}

func (service *DarkService) broadcastAlphaBetaFragment(broadcastAlphaBetaFragmentRequest *rpc.BroadcastAlphaBetaFragmentRequest) (*rpc.AlphaBetaFragment, error) {
	// todo : unimplemented
	return &rpc.AlphaBetaFragment{}, nil
}

// BroadcastDeltaFragment handles an rpc.BroadcastDeltaFragmentRequest
func (service *DarkService) BroadcastDeltaFragment(ctx context.Context, broadcastDeltaFragmentRequest *rpc.BroadcastDeltaFragmentRequest) (*rpc.DeltaFragment, error) {
	wait := do.Process(func() do.Option {
		deltaFragment, err := service.broadcastDeltaFragment(broadcastDeltaFragmentRequest)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(deltaFragment)
	})

	select {
	case val := <-wait:
		if val, ok := val.Ok.(*rpc.DeltaFragment); ok {
			return val, nil
		}
		return &rpc.DeltaFragment{}, nil

	case <-ctx.Done():
		return &rpc.DeltaFragment{}, ctx.Err()
	}
}

func (service *DarkService) broadcastDeltaFragment(broadcastDeltaFragmentRequest *rpc.BroadcastDeltaFragmentRequest) (*rpc.DeltaFragment, error) {
	from, _, err := rpc.DeserializeMultiAddress(broadcastDeltaFragmentRequest.From)
	if err != nil {
		return &rpc.DeltaFragment{}, err
	}
	deltaFragment, err := rpc.DeserializeDeltaFragment(broadcastDeltaFragmentRequest.DeltaFragment)
	if err != nil {
		return &rpc.DeltaFragment{}, err
	}
	service.OnBroadcastDeltaFragment(from, deltaFragment)
	// FIXME: Return the respective delta fragment.
	return &rpc.DeltaFragment{}, nil
}
