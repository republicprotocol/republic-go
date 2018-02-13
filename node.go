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

	resultMu  *sync.RWMutex
	Results   map[identity.Address]Notifications
	newResult map[identity.Address]chan *compute.Result
}

// NewNode returns a Node that delegates the responsibility of handling RPCs to
// a Delegate.
func NewNode(server *grpc.Server, delegate Delegate, options Options) *Node {
	return &Node{
		Delegate:  delegate,
		Server:    server,
		Options:   options,
		resultMu:  new(sync.RWMutex),
		Results:   make(map[identity.Address]Notifications),
		newResult: make(map[identity.Address]chan *compute.Result),
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

// Notify ...
func (node *Node) Notify(result *compute.Result) {
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

func (node *Node) Notifications(traderAddress *rpc.Address, stream rpc.XingNode_NotificationsServer) error {
	if node.Options.Debug >= DebugHigh {
		log.Printf("%v received a query for all results of [%v]\n", node.Address(), traderAddress.Address)
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

func (node *Node) GetResults(ctx context.Context, traderAddress *rpc.Address) (*rpc.Results, error) {
	if node.Options.Debug >= DebugHigh {
		log.Printf("%v registered a trader for notifications [%v]\n", node.Address(), traderAddress.Address)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		results, err := node.getResults(traderAddress)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(results)
	})

	select {
	case val := <-wait:
		if nothing, ok := val.Ok.(*rpc.Results); ok {
			return nothing, val.Err
		}
		return nil, val.Err

	case <-ctx.Done():
		return nil, ctx.Err()
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

type Notifications struct {
	resultRWMu *sync.RWMutex
	resultMu   *sync.Mutex
	results    []*compute.Result
	//cond     *sync.Cond
}

func NewNotifications(result *compute.Result) Notifications {
	return Notifications{
		resultRWMu: new(sync.RWMutex),
		resultMu:   new(sync.Mutex),
		results:    []*compute.Result{result},
		//cond : &sync.Cond{
		//	L:new(sync.Mutex),
		//},
	}
}

func (node *Node) notifications(traderAddress *rpc.Address, stream rpc.XingNode_NotificationsServer) error {
	address := identity.Address(traderAddress.Address)
	node.resultMu.RLock()
	res, ok := node.Results[address]
	node.resultMu.RUnlock()
	if ok {
		for _, j := range res.results {
			stream.Send(rpc.SerializeResult(j))
		}
	}

	go func() {
		for {
			select {
			case new, ok  := <-node.newResult[address]:
				if !ok{
					break
				}
				stream.Send(rpc.SerializeResult(new))
			case stream.Context().Done():

				break
			}
		}
	}()

	return nil
}

func (node *Node) getResults(traderAddress *rpc.Address) (*rpc.Results, error) {
	address := identity.Address(traderAddress.Address)
	node.resultMu.RLock()
	defer node.resultMu.RUnlock()
	notifications, ok := node.Results[address]
	if !ok {
		return *rpc.Results{}, nil
	}
	notifications.resultRWMu.RLock()
	defer notifications.resultRWMu.RUnlock()
	ret := *rpc.Results{}

	for _, j := range notifications.results {
		ret = append(ret, j)
	}
	return ret, nil
}

func (node *Node) UpdateResult(result *compute.Result, traderAddress identity.Address) {
	node.resultMu.RLock()
	notifications, ok := node.Results[traderAddress]
	node.resultMu.RUnlock()
	if !ok {
		node.resultMu.Lock()
		defer node.resultMu.Unlock()
		node.Results[traderAddress] = NewNotifications(result)
	} else {
		notifications.resultMu.Lock()
		defer notifications.resultMu.Unlock()
		notifications.results = append(notifications.results, result)
	}
	go func() {
		ch := make(chan *compute.Result)
		node.resultMu.Lock()
		node.newResult[traderAddress] = ch
		node.resultMu.Unlock()
		ch <- result
	}()
}
