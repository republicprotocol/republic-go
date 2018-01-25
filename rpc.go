package network

import (
	"fmt"
	"time"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-network/rpc"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-sss"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Dial the target rpc.MultiAddress using a background context.Context. Returns
// a grpc.ClientConn, or an error. The grpc.ClientConn must be closed before it
// exists scope.
func Dial(target *rpc.MultiAddress) (*grpc.ClientConn, error) {
	targetMultiAddress, err := DeserializeMultiAddress(target)
	if err != nil {
		return nil, err
	}
	host, err := targetMultiAddress.ValueForProtocol(identity.IP4Code)
	if err != nil {
		return nil, err
	}
	port, err := targetMultiAddress.ValueForProtocol(identity.TCPCode)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// PingTarget using a new grpc.ClientConn to make a Ping RPC to a target
// rpc.MultiAddress.
func PingTarget(to *rpc.MultiAddress, from *rpc.MultiAddress) (*rpc.Nothing, error) {
	conn, err := Dial(to)
	if err != nil {
		return &rpc.Nothing{}, err
	}
	defer conn.Close()
	client := rpc.NewNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return client.Ping(ctx, from, grpc.FailFast(false))
}

// GetPeersFromTarget using a new grpc.ClientConn to make a Peers RPC to a
// target rpc.MultiAddress.
func GetPeersFromTarget(to *rpc.MultiAddress, from *rpc.MultiAddress) (*rpc.MultiAddresses, error) {
	conn, err := Dial(to)
	if err != nil {
		return &rpc.MultiAddresses{Multis: []*rpc.MultiAddress{}}, nil
	}
	defer conn.Close()
	client := rpc.NewNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return client.Peers(ctx, from, grpc.FailFast(false))
}

// FindPeerFromTarget using a new grpc.ClientConn to make a FindPeer RPC to a
// target rpc.MultiAddress.
func FindPeerFromTarget(to *rpc.MultiAddress, from *rpc.MultiAddress, peer *rpc.Address) (*rpc.MultiAddresses, error) {
	conn, err := Dial(to)
	if err != nil {
		return &rpc.MultiAddresses{Multis: []*rpc.MultiAddress{}}, nil
	}
	defer conn.Close()
	client := rpc.NewNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return client.FindPeer(ctx, &rpc.Finder{Peer: peer, From: from}, grpc.FailFast(false))
}

// SendOrderFragmentToTarget using a new grpc.ClientConn to make a
// SendOrderFragment RPC to a target rpc.MultiAddress.
func SendOrderFragmentToTarget(to *rpc.MultiAddress, orderFragment *rpc.OrderFragment) (*rpc.Nothing, error) {
	conn, err := Dial(to)
	if err != nil {
		return &rpc.Nothing{}, err
	}
	defer conn.Close()
	client := rpc.NewNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return client.SendOrderFragment(ctx, orderFragment, grpc.FailFast(false))
}

// SendResultFragmentToTarget using a new grpc.ClientConn to make a
// SendResultFragment RPC to a target rpc.MultiAddress.
func SendResultFragmentToTarget(to *rpc.MultiAddress, resultFragment *rpc.ResultFragment) (*rpc.Nothing, error) {
	conn, err := Dial(to)
	if err != nil {
		return &rpc.Nothing{}, err
	}
	defer conn.Close()
	client := rpc.NewNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return client.SendResultFragment(ctx, resultFragment, grpc.FailFast(false))
}

// SerializeMultiAddress converts an identity.MultiAddress into its
// rpc.MultiAddress network representation.
func SerializeMultiAddress(multiAddress identity.MultiAddress) *rpc.MultiAddress {
	return &rpc.MultiAddress{Multi: multiAddress.String()}
}

func DeserializeMultiAddress(multiAddress *rpc.MultiAddress) (identity.MultiAddress, error) {
	return identity.NewMultiAddressFromString(multiAddress.Multi)
}

// SerializeMultiAddresses converts identity.MultiAddresses into their
// rpc.MultiAddresses network representation.
func SerializeMultiAddresses(multiAddresses identity.MultiAddresses) *rpc.MultiAddresses {
	serializedMultiAddresses := make([]*rpc.MultiAddress, len(multiAddresses))
	for i, multiAddress := range multiAddresses {
		serializedMultiAddresses[i] = SerializeMultiAddress(multiAddress)
	}
	return &rpc.MultiAddresses{Multis: serializedMultiAddresses}
}

func DeserializeMultiAddresses(multiAddresses *rpc.MultiAddresses) (identity.MultiAddresses, error) {
	deserializedMultiAddresses := make(identity.MultiAddresses, 0, len(multiAddresses.Multis))
	for _, multiAddress := range multiAddresses.Multis {
		deserializedMultiAddress, err := DeserializeMultiAddress(multiAddress)
		if err != nil {
			return deserializedMultiAddresses, err
		}
		deserializedMultiAddresses = append(deserializedMultiAddresses, deserializedMultiAddress)
	}
	return deserializedMultiAddresses, nil
}

// SerializeOrderFragment converts a compute.OrderFragment into its
// rpc.OrderFragment network representation.
func SerializeOrderFragment(input *compute.OrderFragment) *rpc.OrderFragment {
	orderFragment := &rpc.OrderFragment{
		Id:           []byte(input.ID),
		OrderId:      []byte(input.OrderID),
		OrderType:    int64(input.OrderType),
		OrderBuySell: int64(input.OrderBuySell),
	}
	orderFragment.FstCodeShare = sss.ToBytes(input.FstCodeShare)
	orderFragment.SndCodeShare = sss.ToBytes(input.SndCodeShare)
	orderFragment.PriceShare = sss.ToBytes(input.PriceShare)
	orderFragment.MaxVolumeShare = sss.ToBytes(input.MaxVolumeShare)
	orderFragment.MinVolumeShare = sss.ToBytes(input.MinVolumeShare)
	return orderFragment
}

// DeserializeOrderFragment converts a network representation of a
// OrderFragment into a compute.OrderFragment. An error is returned if the
// network representation is malformed.
func DeserializeOrderFragment(input *rpc.OrderFragment) (*compute.OrderFragment, error) {
	orderFragment := &compute.OrderFragment{
		ID:           compute.OrderFragmentID(input.Id),
		OrderID:      compute.OrderID(input.OrderId),
		OrderType:    compute.OrderType(input.OrderType),
		OrderBuySell: compute.OrderBuySell(input.OrderBuySell),
	}

	var err error
	orderFragment.FstCodeShare, err = sss.FromBytes(input.FstCodeShare)
	if err != nil {
		return nil, err
	}
	orderFragment.SndCodeShare, err = sss.FromBytes(input.SndCodeShare)
	if err != nil {
		return nil, err
	}
	orderFragment.PriceShare, err = sss.FromBytes(input.PriceShare)
	if err != nil {
		return nil, err
	}
	orderFragment.MaxVolumeShare, err = sss.FromBytes(input.MaxVolumeShare)
	if err != nil {
		return nil, err
	}
	orderFragment.MinVolumeShare, err = sss.FromBytes(input.MinVolumeShare)
	if err != nil {
		return nil, err
	}
	return orderFragment, nil
}

// SerializeResultFragment converts a compute.ResultFragment into its network
// representation.
func SerializeResultFragment(input *compute.ResultFragment) *rpc.ResultFragment {
	resultFragment := &rpc.ResultFragment{
		Id:                  []byte(input.ID),
		BuyOrderId:          []byte(input.BuyOrderID),
		SellOrderId:         []byte(input.SellOrderID),
		BuyOrderFragmentId:  []byte(input.BuyOrderFragmentID),
		SellOrderFragmentId: []byte(input.SellOrderFragmentID),
	}
	resultFragment.FstCodeShare = sss.ToBytes(input.FstCodeShare)
	resultFragment.SndCodeShare = sss.ToBytes(input.SndCodeShare)
	resultFragment.PriceShare = sss.ToBytes(input.PriceShare)
	resultFragment.MaxVolumeShare = sss.ToBytes(input.MaxVolumeShare)
	resultFragment.MinVolumeShare = sss.ToBytes(input.MinVolumeShare)
	return resultFragment
}

// DeserializeResultFragment converts a network representation of a
// ResultFragment into a compute.ResultFragment. An error is returned if the
// network representation is malformed.
func DeserializeResultFragment(input *rpc.ResultFragment) (*compute.ResultFragment, error) {
	resultFragment := &compute.ResultFragment{
		ID:                  compute.ResultFragmentID(input.Id),
		BuyOrderID:          compute.OrderID(input.BuyOrderId),
		SellOrderID:         compute.OrderID(input.SellOrderId),
		BuyOrderFragmentID:  compute.OrderFragmentID(input.BuyOrderFragmentId),
		SellOrderFragmentID: compute.OrderFragmentID(input.SellOrderFragmentId),
	}

	var err error
	resultFragment.FstCodeShare, err = sss.FromBytes(input.FstCodeShare)
	if err != nil {
		return nil, err
	}
	resultFragment.SndCodeShare, err = sss.FromBytes(input.SndCodeShare)
	if err != nil {
		return nil, err
	}
	resultFragment.PriceShare, err = sss.FromBytes(input.PriceShare)
	if err != nil {
		return nil, err
	}
	resultFragment.MaxVolumeShare, err = sss.FromBytes(input.MaxVolumeShare)
	if err != nil {
		return nil, err
	}
	resultFragment.MinVolumeShare, err = sss.FromBytes(input.MinVolumeShare)
	if err != nil {
		return nil, err
	}
	return resultFragment, nil
}
