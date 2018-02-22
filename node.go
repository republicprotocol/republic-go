package dark

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
	OnOrderFragmentForwarding(to identity.Address, from identity.MultiAddress, orderFragment *compute.OrderFragment)
	OnSync(from identity.MultiAddress) chan do.Option
	OnElectShard(from identity.MultiAddress, shard compute.Shard) compute.Shard
	OnComputeShard(from identity.MultiAddress, shard compute.Shard)
	OnFinalizeShard(from identity.MultiAddress, shard compute.Shard)
}

// Node implements the gRPC Node service.
type Node struct {
	Delegate
	Server  *grpc.Server
	Options Options
	results *sync.Map
}

// NewNode returns a Node that delegates the responsibility of handling RPCs to
// a Delegate.
func NewNode(server *grpc.Server, delegate Delegate, options Options) *Node {
	return &Node{
		Delegate: delegate,
		Server:   server,
		Options:  options,
		results:  new(sync.Map),
	}
}

// Register the gRPC service.
func (node *Node) Register() {
	rpc.RegisterDarkNodeServer(node.Server, node)
}

// Address returns the identity.Address of the Node.
func (node *Node) Address() identity.Address {
	return node.Options.Address
}

// Sync returns all deltaFragments and residueFragments the node have.
func (node *Node) Sync(syncRequest *rpc.SyncRequest, stream rpc.DarkNode_SyncServer) error {
	if node.Options.Debug >= DebugHigh {
		log.Printf("[%v] received a sync query from [%v]\n", node.Address(), syncRequest.From.Multi)
	}
	if err := stream.Context().Err(); err != nil {
		return err
	}

	wait := do.Process(func() do.Option {
		return do.Err(node.Sync(syncRequest, stream))
	})

	select {
	case val := <-wait:
		return val.Err

	case <-stream.Context().Done():
		return stream.Context().Err()
	}
}

// ComputeShard will start compute the shards on receiving the request.
func (node *Node) ComputeShard(ctx context.Context, computeShardRequest *rpc.ComputeShardRequest) (*rpc.Nothing, error) {
	if node.Options.Debug >= DebugHigh {
		log.Printf("[%v] received a compute shard query from [%v]\n", node.Address(), computeShardRequest.From.Multi)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		nothing, err := node.computeShard(computeShardRequest)
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

// ElectShard will returns availability of the shards listed in the request.
func (node *Node) ElectShard(ctx context.Context, electShardRequest *rpc.ElectShardRequest) (*rpc.Shard, error) {
	if node.Options.Debug >= DebugHigh {
		log.Printf("[%v] received a elect shard query from [%v]\n", node.Address(), electShardRequest.From.Multi)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		shard, err := node.electShard(electShardRequest)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(shard)
	})

	select {
	case val := <-wait:
		if shard, ok := val.Ok.(*rpc.Shard); ok {
			return shard, val.Err
		}
		return &rpc.Shard{}, val.Err

	case <-ctx.Done():
		return &rpc.Shard{}, ctx.Err()
	}
}

// FinalizeShard
func (node *Node) FinalizeShard(ctx context.Context, finaliseShardRequest *rpc.FinalizeShardRequest) (*rpc.Nothing, error) {
	if node.Options.Debug >= DebugHigh {
		log.Printf("[%v] received a finalize shard request from [%v]\n", node.Address(), finaliseShardRequest.From.Multi)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		nothing, err := node.finalizeShard(finaliseShardRequest)
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

// SendOrderFragmentCommitment is sent before sending the order fragment.
// The request contained the signature of the sender and we'll return
// a commitment with our signature.
func (node *Node) SendOrderFragmentCommitment(ctx context.Context, orderFragmentCommitment *rpc.OrderFragmentCommitment) (*rpc.OrderFragmentCommitment, error) {
	if node.Options.Debug >= DebugHigh {
		log.Printf("%v received a order commitment from %v\n", node.Address(), orderFragmentCommitment.From.Multi)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		commitment, err := node.sendOrderFragmentCommitment(orderFragmentCommitment)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(commitment)
	})

	select {
	case val := <-wait:
		if commitment, ok := val.Ok.(*rpc.OrderFragmentCommitment); ok {
			return commitment, val.Err
		}
		return &rpc.OrderFragmentCommitment{}, val.Err

	case <-ctx.Done():
		return &rpc.OrderFragmentCommitment{}, ctx.Err()
	}
}

// SendOrderFragment to the Node. If the rpc.OrderFragment is not destined for
// this Node then it will be forwarded on to the correct destination.
func (node *Node) SendOrderFragment(ctx context.Context, orderFragment *rpc.OrderFragment) (*rpc.Nothing, error) {
	if node.Options.Debug >= DebugHigh {
		log.Printf("[%v] received order fragment %v [%v]\n", node.Address(), base58.Encode(orderFragment.Id), base58.Encode(orderFragment.OrderId))
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

// Notifications will connect the rpc client with a channel and send all
// unread results to the client via a stream
func (node *Node) Notifications(traderAddress *rpc.MultiAddress, stream rpc.DarkNode_NotificationsServer) error {
	if node.Options.Debug >= DebugHigh {
		log.Printf("[%v] received a query for notifications of [%v]\n", node.Address(), traderAddress.Multi)
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

// GetFinals will connect the rpc client with a channel and send all
// related results to the client via a stream
func (node *Node) GetFinals(traderAddress *rpc.MultiAddress, stream rpc.DarkNode_GetFinalsServer) error {
	if node.Options.Debug >= DebugHigh {
		log.Printf("[%v] received a query for all results of [%v]\n", node.Address(), traderAddress.Multi)
	}

	if err := stream.Context().Err(); err != nil {
		return err
	}

	wait := do.Process(func() do.Option {
		return do.Err(node.getFinals(traderAddress, stream))
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

// Notify will store new result in the node.
func (node *Node) Notify(traderAddress identity.Address, result *compute.Final) {
	results, ok := node.results.Load(traderAddress)
	if !ok {
		newInbox := NewInbox()
		newInbox.AddNewResult(result)
		node.results.Store(traderAddress, newInbox)
	} else {
		results.(*Inbox).AddNewResult(result)
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

func (node *Node) notifications(traderAddress *rpc.MultiAddress, stream rpc.DarkNode_NotificationsServer) error {
	multiAddress, err := rpc.DeserializeMultiAddress(traderAddress)
	if err != nil {
		return err
	}
	address := identity.Address(multiAddress.Address())

	results, ok := node.results.Load(address)
	if !ok {
		return nil
	}
	for {
		results := results.(*Inbox).GetAllNewResults()
		for i := range results {
			err := stream.Send(rpc.SerializeFinal(results[i]))
			if err != nil {
				return err
			}
		}
	}
}

func (node *Node) getFinals(traderAddress *rpc.MultiAddress, stream rpc.DarkNode_GetFinalsServer) error {
	multiAddress, err := rpc.DeserializeMultiAddress(traderAddress)
	if err != nil {
		return err
	}
	address := multiAddress.Address()
	notifications, ok := node.results.Load(address)
	if !ok {
		return nil
	}
	results := notifications.(*Inbox).GetAllResults()
	for _, result := range results {
		err = stream.Send(rpc.SerializeFinal(result))
		if err != nil {
			return err
		}
	}
	return nil
}

func (node *Node) sync(syncRequest *rpc.SyncRequest, stream rpc.DarkNode_SyncServer) error {
	from, err := identity.NewMultiAddressFromString(syncRequest.From.String())
	if err != nil {
		return err
	}
	blockChan := node.Delegate.OnSync(from)
	for data := range blockChan {
		//todo : need to serialize data into the network representation
		stream.Send(data.Ok.(*rpc.SyncBlock))
	}
	return nil
}

func (node *Node) computeShard(computeShardRequest *rpc.ComputeShardRequest) (*rpc.Nothing, error) {
	from, err := identity.NewMultiAddressFromString(computeShardRequest.From.String())
	if err != nil {
		return &rpc.Nothing{}, err
	}
	shard := rpc.DeserializeShard(computeShardRequest.Shard)
	node.Delegate.OnComputeShard(from, *shard)
	return &rpc.Nothing{}, nil
}

func (node *Node) electShard(electShardRequest *rpc.ElectShardRequest) (*rpc.Shard, error) {
	from, err := identity.NewMultiAddressFromString(electShardRequest.From.String())
	if err != nil {
		return &rpc.Shard{}, err
	}
	shard := rpc.DeserializeShard(electShardRequest.Shard)
	shardReturn := node.Delegate.OnElectShard(from, *shard)
	return rpc.SerializeShard(shardReturn), nil
}

func (node *Node) finalizeShard(finaliseShardRequest *rpc.FinalizeShardRequest) (*rpc.Nothing, error) {
	from, err := identity.NewMultiAddressFromString(finaliseShardRequest.From.String())
	if err != nil {
		return &rpc.Nothing{}, err
	}
	shard := rpc.DeserializeShard(finaliseShardRequest.Shard)
	node.Delegate.OnFinalizeShard(from, *shard)
	return &rpc.Nothing{}, nil
}

func (node *Node) sendOrderFragmentCommitment(orderFragmentCommitment *rpc.OrderFragmentCommitment) (*rpc.OrderFragmentCommitment, error) {
	// todo :
	return orderFragmentCommitment, nil
}
