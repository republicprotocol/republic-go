package network

import (
	"context"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/compute"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
	"google.golang.org/grpc"
)

// A DarkOceanDelegate is used as a callback interface to inject behavior into
// the DarkOceanService service.
type DarkOceanDelegate interface {
	OnLog(from identity.MultiAddress)
	OnSync(from identity.MultiAddress)

	OnSignOrderFragment(from identity.MultiAddress)
	OnOpenOrder(from identity.MultiAddress, orderFragment *order.Fragment)
	OnCancelOrder(from identity.MultiAddress)

	OnRandomFragmentShares(from identity.MultiAddress)
	OnResidueFragmentShares(from identity.MultiAddress)
	OnComputeResidueFragment(from identity.MultiAddress)
	OnBroadcastAlphaBetaFragment(from identity.MultiAddress)
	OnBroadcastDeltaFragment(from identity.MultiAddress, deltaFragment *compute.DeltaFragment)
}

// DarkOceanService implements the gRPC DarkOceanService service.
type DarkOceanService struct {
	DarkOceanDelegate
	Options
}

func NewDarkOcean(delegate DarkOceanDelegate, options Options) *DarkOceanService {
	return &DarkOceanService{
		Delegate: delegate,
		Options:  options,
	}
}

// Register the gRPC service.
func (service *DarkOceanService) Register(*grpc.Server) {
	rpc.RegisterDarkNodeServer(server, service)
}

func (service *DarkOceanService) Log(logRequest *rpc.LogRequest, stream rpc.DarkOcean_LogServer) error {
	wait := do.Process(func() do.Option {
		return do.Err(node.log(logRequest, stream))
	})

	select {
	case val := <-wait:
		return val.Err

	case <-stream.Context().Done():
		return stream.Context().Err()
	}
}

func (service *DarkOceanService) log(logRequest *rpc.LogRequest, stream rpc.DarkOcean_LogServer) error {
	from, err := rpc.DeserializeMultiAddress(logRequest.From)
	if err != nil {
		return err
	}
	service.DarkOceanDelegate.OnLog(from)
	panic("unimplemented")
}

func (service *DarkOceanService) Sync(syncRequest *rpc.SyncRequest, stream rpc.DarkOcean_SyncServer) error {
	wait := do.Process(func() do.Option {
		return do.Err(node.sync(syncRequest, stream))
	})

	select {
	case val := <-wait:
		return val.Err

	case <-stream.Context().Done():
		return stream.Context().Err()
	}
}

func (service *DarkOceanService) sync(syncRequest *rpc.SyncRequest, stream rpc.DarkOcean_SyncServer) error {
	from, err := rpc.DeserializeMultiAddress(syncRequest.From)
	if err != nil {
		return err
	}
	service.DarkOceanDelegate.OnSync(from)
	panic("unimplemented")
}

func (service *DarkOceanService) SignOrderFragment(signOrderFragmentRequest *rpc.SignOrderFragmentRequest) (*rpc.OrderFragmentSignature, error) {
	wait := do.Process(func() do.Option {
		orderFragmentSignature, err := node.signOrderFragment(signOrderFragmentRequest)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(orderFragmentSignature)
	})

	select {
	case val := <-wait:
		if val, ok := val.Ok.(*rpc.OrderFragmentSignature); ok {
			return val, val.Err
		}
		return &rpc.OrderFragmentSignature{}, nil
	case <-stream.Context().Done():
		return stream.Context().Err()
	}
}

func (service *DarkOceanService) signOrderFragment(signOrderFragmentRequest *rpc.SignOrderFragmentRequest) (*rpc.OrderFragmentSignature, error) {
	from, err := rpc.DeserializeMultiAddress(syncRequest.From)
	if err != nil {
		return nil, err
	}
	service.DarkOceanDelegate.OnSignOrderFragment(from)
	panic("unimplemented")
}

func (service *DarkOceanService) OpenOrder(ctx context.Context, openOrderRequest *rpc.OpenOrderRequest) (*rpc.Nothing, error) {
	wait := do.Process(func() do.Option {
		nothing, err := node.openOrder(openOrderRequest)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(nothing)
	})

	select {
	case val := <-wait:
		if nothing, ok := val.Ok.(*rpc.Nothing); ok {
			return nothing, val.Err
		}
		return &rpc.Nothing{}, val.Err

	case <-ctx.Done():
		return &rpc.Nothing{}, ctx.Err()
	}
}

func (service *DarkOceanService) openOrder(openOrderRequest *rpc.OpenOrderRequest) (*rpc.Nothing, error) {
	from, err := rpc.DeserializeMultiAddress(openOrderRequest.From)
	if err != nil {
		return &rpc.Nothing{}, err
	}
	orderFragment, err := rpc.DeserializeOrderFragment(openOrderRequest.OrderFragment)
	if err != nil {
		return &rpc.Nothing{}, err
	}
	node.OnOpenOrder(from, orderFragment)
	return &rpc.Nothing{}, nil
}

func (service *DarkOceanService) CancelOrder(ctx context.Context, cancelOrderRequest *rpc.CancelOrderRequest) (*rpc.Nothing, error) {
	wait := do.Process(func() do.Option {
		nothing, err := node.cancelOrder(cancelOrderRequest)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(nothing)
	})

	select {
	case val := <-wait:
		if nothing, ok := val.Ok.(*rpc.Nothing); ok {
			return nothing, val.Err
		}
		return &rpc.Nothing{}, val.Err

	case <-ctx.Done():
		return &rpc.Nothing{}, ctx.Err()
	}
}

func (service *DarkOceanService) cancelOrder(cancelOrderRequest *rpc.CancelOrderRequest) (*rpc.Nothing, error) {
	from, err := rpc.DeserializeMultiAddress(cancelOrderRequest.From)
	if err != nil {
		return &rpc.Nothing{}, err
	}
	node.OnCancelOrder(from, orderFragment)
	panic("unimplemented")
}

func (service *DarkOceanService) RandomFragmentShares(ctx context.Context, randomFragmentSharesRequest *rpc.RandomFragmentSharesRequest) (*rpc.RandomFragments, error) {
	wait := do.Process(func() do.Option {
		nothing, err := node.randomFragmentShares(randomFragmentSharesRequest)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(nothing)
	})

	select {
	case val := <-wait:
		if val, ok := val.Ok.(*rpc.RandomFragments); ok {
			return val, val.Err
		}
		return &rpc.RandomFragments{}, nil

	case <-ctx.Done():
		return &rpc.RandomFragments{}, ctx.Err()
	}
}

func (service *DarkOceanService) randomFragmentShares(randomFragmentSharesRequest *rpc.RandomFragmentSharesRequest) (*rpc.RandomFragments, error) {
	from, err := rpc.DeserializeMultiAddress(randomFragmentSharesRequest.From)
	if err != nil {
		return &rpc.RandomFragments{}, err
	}
	node.OnRandomFragmentShares(from)
	panic("unimplemented")
}

func (service *DarkOceanService) ResidueFragmentShares(ctx context.Context, residueFragmentSharesRequest *rpc.ResidueFragmentSharesRequest) (*rpc.ResidueFragments, error) {
	wait := do.Process(func() do.Option {
		nothing, err := node.residueFragmentShares(residueFragmentSharesRequest)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(nothing)
	})

	select {
	case val := <-wait:
		if val, ok := val.Ok.(*rpc.ResidueFragments); ok {
			return val, val.Err
		}
		return &rpc.ResidueFragments{}, val.Err

	case <-ctx.Done():
		return &rpc.ResidueFragments{}, ctx.Err()
	}
}

func (service *DarkOceanService) residueFragmentShares(residueFragmentSharesRequest *rpc.ResidueFragmentSharesRequest) (*rpc.ResidueFragments, error) {
	from, err := rpc.DeserializeMultiAddress(residueFragmentSharesRequest.From)
	if err != nil {
		return &rpc.ResidueFragments{}, err
	}
	node.OnResidueFragmentShares(from)
	panic("unimplemented")
}

func (service *DarkOceanService) ComputeResidueFragment(ctx context.Context, computeResidueFragmentRequest *rpc.ComputeResidueFragmentRequest) (*rpc.Nothing, error) {
	wait := do.Process(func() do.Option {
		nothing, err := node.computeResidueFragment(computeResidueFragmentRequest)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(nothing)
	})

	select {
	case val := <-wait:
		if nothing, ok := val.Ok.(*rpc.Nothing); ok {
			return nothing, val.Err
		}
		return &rpc.Nothing{}, val.Err

	case <-ctx.Done():
		return &rpc.Nothing{}, ctx.Err()
	}
}

func (service *DarkOceanService) computeResidueFragment(computeResidueFragmentRequest *rpc.ComputeResidueFragmentRequest) (*rpc.Nothing, error) {
	from, err := rpc.DeserializeMultiAddress(computeResidueFragmentRequest.From)
	if err != nil {
		return &rpc.Nothing{}, err
	}
	node.OnComputeResidueFragment(from)
	panic("unimplemented")
}

func (service *DarkOceanService) BroadcastAlphaBetaFragment(ctx context.Context, broadcastAlphaBetaFragmentRequest *rpc.BroadcastAlphaBetaFragmentRequest) (*rpc.AlphaBetaFragment, error) {
	wait := do.Process(func() do.Option {
		nothing, err := node.broadcastAlphaBetaFragment(broadcastAlphaBetaFragmentRequest)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(nothing)
	})

	select {
	case val := <-wait:
		if nothing, ok := val.Ok.(*rpc.Nothing); ok {
			return nothing, val.Err
		}
		return &rpc.Nothing{}, val.Err

	case <-ctx.Done():
		return &rpc.Nothing{}, ctx.Err()
	}
}

func (service *DarkOceanService) broadcastAlphaBetaFragment(broadcastAlphaBetaFragmentRequest *rpc.BroadcastAlphaBetaFragmentRequest) (*rpc.AlphaBetaFragment, error) {
	from, err := rpc.DeserializeMultiAddress(broadcastAlphaBetaFragmentRequest.From)
	if err != nil {
		return &rpc.AlphaBetaFragment{}, err
	}
	service.OnBroadcastAlphaBetaFragment(from)
	panic("unimplemented")
}

func (service *DarkOceanService) BroadcastDeltaFragment(ctx context.Context, broadcastDeltaFragmentRequest *rpc.BroadcastDeltaFragmentRequest) (*rpc.DeltaFragment, error) {
	wait := do.Process(func() do.Option {
		deltaFragment, err := node.broadcastDeltaFragment(broadcastDeltaFragmentRequest)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(deltaFragment)
	})

	select {
	case val := <-wait:
		if val, ok := val.Ok.(*rpc.DeltaFragment); ok {
			return val, val.Err
		}
		return &rpc.DeltaFragment{}, nil

	case <-ctx.Done():
		return &rpc.DeltaFragment{}, ctx.Err()
	}
}

func (service *DarkOceanService) broadcastDeltaFragment(broadcastDeltaFragmentRequest *rpc.BroadcastDeltaFragmentRequest) (*rpc.DeltaFragment, error) {
	from, err := rpc.DeserializeMultiAddress(broadcastDeltaFragmentRequest.From)
	if err != nil {
		return &rpc.DeltaFragment{}, err
	}
	deltaFragment, err := rpc.DeserializeDeltaFragment(broadcastDeltaFragmentRequest.DeltaFragment)
	if err != nil {
		return &rpc.DeltaFragment{}, err
	}
	node.OnBroadcastDeltaFragment(from, deltaFragment)
	panic("unimplemented")
}
