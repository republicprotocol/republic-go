package xing

import (
	"context"
	"log"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-rpc"
	"google.golang.org/grpc"
)

// A Delegate is used as a callback interface to inject behavior into the
// different RPCs.
type Delegate interface {
	OnOrderFragmentReceived(from identity.MultiAddress, orderFragment *compute.OrderFragment)
	OnResultFragmentReceived(from identity.MultiAddress, resultFragment *compute.ResultFragment)
	OnOrderFragmentForwarding(to identity.Address, from identity.MultiAddress, orderFragment *compute.OrderFragment)
	OnResultFragmentForwarding(to identity.Address, from identity.MultiAddress, resultFragment *compute.ResultFragment)
}

// Node implements the gRPC Node service.
type Node struct {
	Delegate
	Server  *grpc.Server
	Options Options
}

// NewNode returns a Node that delegates the responsibility of handling RPCs to
// a Delegate.
func NewNode(server *grpc.Server, delegate Delegate, options Options) *Node {
	return &Node{
		Delegate: delegate,
		Server:   server,
		Options:  options,
	}
}

// Register the gRPC service.
func (node *Node) Register() {
	rpc.RegisterXingNodeServer(node.Server, node)
}

// Address returns the identity.Address of the Node.
func (node *Node) Address() identity.Address {
	return node.Options.Address
}

// SendOrderFragment to the Node. If the rpc.OrderFragment is not destined for
// this Node then it will be forwarded on to the correct destination.
func (node *Node) SendOrderFragment(ctx context.Context, orderFragment *rpc.OrderFragment) (*rpc.Nothing, error) {
	if node.Options.Debug >= DebugHigh {
		log.Printf("%v is receiving an order fragment...\n", node.Address())
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		nothing, err := node.sendOrderFragment(orderFragment)
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

// SendResultFragment to the Node. If the rpc.ResultFragment is not destined
// for this Node then it will be forwarded on to the correct destination.
func (node *Node) SendResultFragment(ctx context.Context, resultFragment *rpc.ResultFragment) (*rpc.Nothing, error) {
	if node.Options.Debug >= DebugHigh {
		log.Printf("%v is receiving a result fragment...\n", node.Address())
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		nothing, err := node.sendResultFragment(resultFragment)
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
		return nil, val.Err

	case <-ctx.Done():
		return &rpc.Nothing{}, ctx.Err()
	}
}

func (node *Node) sendOrderFragment(orderFragment *rpc.OrderFragment) (*rpc.Nothing, error) {
	deserializedTo, err := rpc.DeserializeAddress(orderFragment.To)
	if err != nil {
		return &rpc.Nothing{}, err
	}
	deserializedFrom, err := rpc.DeserializeMultiAddress(orderFragment.From)
	if err != nil {
		return &rpc.Nothing{}, err
	}
	deserializedOrderFragment, err := rpc.DeserializeOrderFragment(orderFragment)
	if err != nil {
		return &rpc.Nothing{}, err
	}

	// If the compute.OrderFragment needs to be forwarded.
	if deserializedTo != node.Address() {
		node.OnOrderFragmentForwarding(deserializedTo, deserializedFrom, deserializedOrderFragment)
		return &rpc.Nothing{}, nil
	}

	// Otherwise it has reached its destination.
	node.OnOrderFragmentReceived(deserializedFrom, deserializedOrderFragment)
	return &rpc.Nothing{}, nil
}

func (node *Node) sendResultFragment(resultFragment *rpc.ResultFragment) (*rpc.Nothing, error) {
	deserializedTo, err := rpc.DeserializeAddress(resultFragment.To)
	if err != nil {
		return &rpc.Nothing{}, err
	}
	deserializedFrom, err := rpc.DeserializeMultiAddress(resultFragment.From)
	if err != nil {
		return &rpc.Nothing{}, err
	}
	deserializedResultFragment, err := rpc.DeserializeResultFragment(resultFragment)
	if err != nil {
		return &rpc.Nothing{}, err
	}

	// If the compute.ResultFragment needs to be forwarded.
	if deserializedTo != node.Address() {
		node.OnResultFragmentForwarding(deserializedTo, deserializedFrom, deserializedResultFragment)
		return &rpc.Nothing{}, nil
	}

	// Otherwise it has reached its destination.
	node.OnResultFragmentReceived(deserializedFrom, deserializedResultFragment)
	return &rpc.Nothing{}, nil
}
