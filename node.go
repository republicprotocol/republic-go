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
	OnSync(from identity.MultiAddress) chan do.Option
	OnOrderFragmentReceived(from identity.MultiAddress, orderFragment *compute.OrderFragment)
	OnBroadcastDeltaFragment(from identity.MultiAddress, deltaFragment *compute.DeltaFragment)
	SubscribeToLogs(chan do.Option)
	UnsubscribeFromLogs(chan do.Option)
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
		return do.Err(node.sync(syncRequest, stream))
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
	panic("unimplemented")
}

// ElectShard will returns availability of the shards listed in the request.
func (node *Node) ElectShard(ctx context.Context, electShardRequest *rpc.ElectShardRequest) (*rpc.Shard, error) {
	panic("unimplemented")
}

// FinalizeShard returns finalized shards.
func (node *Node) FinalizeShard(ctx context.Context, finaliseShardRequest *rpc.FinalizeShardRequest) (*rpc.Nothing, error) {
	panic("unimplemented")
}

// SendOrderFragmentCommitment is sent before sending the order fragment.
// The request contained the signature of the sender and we'll return
// a commitment with our signature.
func (node *Node) SendOrderFragmentCommitment(ctx context.Context, orderFragmentCommitment *rpc.OrderFragmentCommitment) (*rpc.OrderFragmentCommitment, error) {
	panic("unimplemented")
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

// Logs will create a logs channel with the delegate and send any received logs
// through the RPC stream.
func (node *Node) Logs(logRequest *rpc.LogRequest, stream rpc.DarkNode_LogsServer) error {
	if node.Options.Debug >= DebugHigh {
		log.Printf("[%v] received a log query", node.Address())
	}
	if err := stream.Context().Err(); err != nil {
		return err
	}

	wait := do.Process(func() do.Option {
		return do.Err(node.logs(logRequest, stream))
	})

	select {
	case val := <-wait:
		return val.Err

	case <-stream.Context().Done():
		return stream.Context().Err()
	}
}

// BroadcastDeltaFragment receives the delta fragment from the broadcast.
func (node *Node) BroadcastDeltaFragment(ctx context.Context, broadcastDeltaFragmentRequest *rpc.BroadcastDeltaFragmentRequest) (*rpc.DeltaFragment, error) {
	if node.Options.Debug >= DebugHigh {
		log.Printf("[%v] received delta fragment [%v]\n", node.Address(), base58.Encode(broadcastDeltaFragmentRequest.DeltaFragment.Id))
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		nothing, err := node.broadcastDeltaFragment(broadcastDeltaFragmentRequest)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(nothing)
	})

	select {
	case val := <-wait:
		if shard, ok := val.Ok.(*rpc.DeltaFragment); ok {
			return shard, val.Err
		}
		return &rpc.DeltaFragment{}, val.Err

	case <-ctx.Done():
		return &rpc.DeltaFragment{}, ctx.Err()
	}
}

func (node *Node) broadcastDeltaFragment(broadcastDeltaFragmentRequest *rpc.BroadcastDeltaFragmentRequest) (*rpc.DeltaFragment, error) {
	from, err := rpc.DeserializeMultiAddress(broadcastDeltaFragmentRequest.From)
	if err != nil {
		return &rpc.DeltaFragment{}, err
	}
	deltaFragment, err := rpc.DeserializeDeltaFragment(broadcastDeltaFragmentRequest.DeltaFragment)
	if err != nil {
		return nil, err
	}
	node.Delegate.OnBroadcastDeltaFragment(from, deltaFragment)
	return &rpc.DeltaFragment{}, nil
}

//// Notifications will connect the rpc client with a channel and send all
//// unread results to the client via a stream
//func (node *Node) Notifications(traderAddress *rpc.MultiAddress, stream rpc.DarkNode_NotificationsServer) error {
//	if node.Options.Debug >= DebugHigh {
//		log.Printf("[%v] received a query for notifications of [%v]\n", node.Address(), traderAddress.Multi)
//	}
//	if err := stream.Context().Err(); err != nil {
//		return err
//	}
//
//	wait := do.Process(func() do.Option {
//		return do.Err(node.notifications(traderAddress, stream))
//	})
//
//	select {
//	case val := <-wait:
//		return val.Err
//
//	case <-stream.Context().Done():
//		return stream.Context().Err()
//	}
//}
//
//// GetFinals will connect the rpc client with a channel and send all
//// related results to the client via a stream
//func (node *Node) GetFinals(traderAddress *rpc.MultiAddress, stream rpc.DarkNode_GetFinalsServer) error {
//	if node.Options.Debug >= DebugHigh {
//		log.Printf("[%v] received a query for all results of [%v]\n", node.Address(), traderAddress.Multi)
//	}
//
//	if err := stream.Context().Err(); err != nil {
//		return err
//	}
//
//	wait := do.Process(func() do.Option {
//		return do.Err(node.getFinals(traderAddress, stream))
//	})
//
//	for {
//		select {
//		case val := <-wait:
//			return val.Err
//
//		case <-stream.Context().Done():
//			return stream.Context().Err()
//		}
//	}
//}

//// Notify will store new result in the node.
//func (node *Node) Notify(traderAddress identity.Address, result *compute.Final) {
//	results, ok := node.results.Load(traderAddress)
//	if !ok {
//		newInbox := NewInbox()
//		newInbox.AddNewResult(result)
//		node.results.Store(traderAddress, newInbox)
//	} else {
//		results.(*Inbox).AddNewResult(result)
//	}
//}

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
		// This is not meant for us. Do nothing.
		return &rpc.Nothing{}, nil
	}

	// Otherwise it has reached its destination.
	node.OnOrderFragmentReceived(deserializedFrom, deserializedOrderFragment)
	return &rpc.Nothing{}, nil
}

//func (node *Node) notifications(traderAddress *rpc.MultiAddress, stream rpc.DarkNode_NotificationsServer) error {
//	multiAddress, err := rpc.DeserializeMultiAddress(traderAddress)
//	if err != nil {
//		return err
//	}
//	address := identity.Address(multiAddress.Address())
//
//	results, ok := node.results.Load(address)
//	if !ok {
//		return nil
//	}
//	for {
//		results := results.(*Inbox).GetAllNewResults()
//		for i := range results {
//			err := stream.Send(rpc.SerializeFinal(results[i]))
//			if err != nil {
//				return err
//			}
//		}
//	}
//}
//
//func (node *Node) getFinals(traderAddress *rpc.MultiAddress, stream rpc.DarkNode_GetFinalsServer) error {
//	multiAddress, err := rpc.DeserializeMultiAddress(traderAddress)
//	if err != nil {
//		return err
//	}
//	address := multiAddress.Address()
//	notifications, ok := node.results.Load(address)
//	if !ok {
//		return nil
//	}
//	results := notifications.(*Inbox).GetAllResults()
//	for _, result := range results {
//		err = stream.Send(rpc.SerializeFinal(result))
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

func (node *Node) sync(syncRequest *rpc.SyncRequest, stream rpc.DarkNode_SyncServer) error {
	from, err := rpc.DeserializeMultiAddress(syncRequest.From)
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

func (node *Node) logs(logsRequest *rpc.LogRequest, stream rpc.DarkNode_LogsServer) error {
	logChannel := make(chan do.Option)
	node.Delegate.SubscribeToLogs(logChannel)
	defer node.Delegate.UnsubscribeFromLogs(logChannel)
	for event := range logChannel {
		// TODO: need to serialize data into the network representation
		if err := stream.Send(event.Ok.(*rpc.LogEvent)); err != nil {
			return err
		}
	}
	return nil
}
