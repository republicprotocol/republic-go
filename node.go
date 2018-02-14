package xing

import (
	"context"
	"log"
	"sync"

	"github.com/jbenet/go-base58"
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

	resultReceived bool
	resultsMu      *sync.RWMutex
	results        map[identity.Address]*Inbox
}

// NewNode returns a Node that delegates the responsibility of handling RPCs to
// a Delegate.
func NewNode(server *grpc.Server, delegate Delegate, options Options) *Node {
	return &Node{
		Delegate:       delegate,
		Server:         server,
		Options:        options,
		resultReceived: false,
		resultsMu:      new(sync.RWMutex),
		results:        make(map[identity.Address]*Inbox),
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
		log.Printf("%v received order fragment %v [%v]\n", node.Address(), base58.Encode(orderFragment.Id), base58.Encode(orderFragment.OrderId))
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
		log.Printf("%v received result fragment %v [%v, %v]\n", node.Address(), base58.Encode(resultFragment.Id), base58.Encode(resultFragment.BuyOrderId), base58.Encode(resultFragment.SellOrderId))
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

func (node *Node) Notifications(traderAddress *rpc.MultiAddress, stream rpc.XingNode_NotificationsServer) error {
	if node.Options.Debug >= DebugHigh {
		log.Printf("%v received a query for notifications of [%v]\n", node.Address(), traderAddress.Multi)
	}
	if err := stream.Context().Err(); err != nil {
		return err
	}

	wait := do.Process(func() do.Option {
		return do.Err(node.notifications(traderAddress, stream))
	})

	select {
	case val := <-wait:
		return val.Err

	case <-stream.Context().Done():
		return stream.Context().Err()
	}
}

func (node *Node) GetResults(traderAddress *rpc.MultiAddress, stream rpc.XingNode_GetResultsServer) error {
	if node.Options.Debug >= DebugHigh {
		log.Printf("%v received a query for all results of [%v]\n", node.Address(), traderAddress.Multi)
	}

	if err := stream.Context().Err(); err != nil {
		return err
	}

	wait := do.Process(func() do.Option {
		return do.Err(node.getResults(traderAddress, stream))
	})

	for {
		select {
		case val := <-wait:
			return val.Err

		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}

func (node *Node) sendOrderFragment(orderFragment *rpc.OrderFragment) (*rpc.Nothing, error) {
	deserializedTo := rpc.DeserializeAddress(orderFragment.To)
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
	deserializedTo := rpc.DeserializeAddress(resultFragment.To)
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

type Notifications struct {
	do.GuardedObject

	results         []*compute.Result
	resultsNotEmpty *do.Guard
}

func (node *Node) notifications(traderAddress *rpc.MultiAddress, stream rpc.XingNode_NotificationsServer) error {
	multiAddress, err := rpc.DeserializeMultiAddress(traderAddress)
	if err != nil {
		return err
	}
	address := identity.Address(multiAddress.Address())

	node.resultsMu.RLock()
	results, ok := node.results[address]
	node.resultsMu.RUnlock()
	if !ok {
		return nil
	}
	for {
		result := results.GetNewResult()
		stream.Send(rpc.SerializeResult(result))
	}

	return nil
}

func (node *Node) getResults(traderAddress *rpc.MultiAddress, stream rpc.XingNode_GetResultsServer) error {
	multiAddress, err := rpc.DeserializeMultiAddress(traderAddress)
	if err != nil {
		return err
	}
	address := identity.Address(multiAddress.Address())
	node.resultsMu.RLock()
	notifications, ok := node.results[address]
	node.resultsMu.RUnlock()
	if !ok {
		return nil
	}
	notifications.Enter(nil)
	defer notifications.Exit()
	for _, result := range notifications.GetAllResults() {
		stream.Send(rpc.SerializeResult(result))
	}
	return nil
}

// Notify ...
func (node *Node) Notify(traderAddress identity.Address, result *compute.Result) {
	node.resultsMu.RLock()
	results, ok := node.results[traderAddress]
	node.resultsMu.RUnlock()
	if !ok {
		node.resultsMu.Lock()
		defer node.resultsMu.Unlock()
		node.results[traderAddress] = NewInbox()
	} else {
		node.resultsMu.RLock()
		defer node.resultsMu.RUnlock()
		results.AddNewResult(result)
	}
}
