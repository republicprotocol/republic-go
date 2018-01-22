package network

import (
	"log"
	"time"

	identity "github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-network/rpc"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-sss"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// PingTarget uses a new grpc.ClientConn to make a Ping RPC to a target
// identity.MultiAddress. It uses the from identity.MultiAddress to identify
// the sender.
func PingTarget(to identity.MultiAddress, from identity.MultiAddress) error {
	conn, err := Dial(to)
	if err != nil {
		return nil
	}
	defer conn.Close()
	client := rpc.NewNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err = client.Ping(ctx, SerializeMultiAddress(from), grpc.FailFast(false))
	return err
}

// GetPeersFromTarget uses a new grpc.ClientConn to make a Peers RPC to a
// target identity.MultiAddress. It uses the from identity.MultiAddress to
// identify the sender. It returns the identity.MultiAddress of the peers
// connected directly to the target, or an error.
func GetPeersFromTarget(to identity.MultiAddress, from identity.MultiAddress) (identity.MultiAddresses, error) {
	conn, err := Dial(to)
	if err != nil {
		return identity.MultiAddresses{}, nil
	}
	defer conn.Close()
	client := rpc.NewNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	peers, err := client.Peers(ctx, SerializeMultiAddress(from), grpc.FailFast(false))
	if err != nil {
		return identity.MultiAddresses{}, err
	}
	return DeserializeMultiAddresses(peers)
}

func SendOrderFragmentToTarget(to identity.MultiAddress, from identity.MultiAddress, orderFragment *compute.OrderFragment) error {
	address, err := to.Address()
	if err != nil {
		return err
	}
	serializedOrderFragment := SerializeOrderFragment(orderFragment)
	serializedOrderFragment.From = SerializeMultiAddress(from)
	serializedOrderFragment.To = &rpc.Address{Address: address.String()}

	conn, err := Dial(to)
	if err != nil {
		return nil
	}
	defer conn.Close()
	client := rpc.NewNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err = client.SendOrderFragment(ctx, serializedOrderFragment, grpc.FailFast(false))
	return err
}

func SendResultFragmentToTarget(to identity.MultiAddress, from identity.MultiAddress, resultFragment *compute.ResultFragment) error {
	address, err := to.Address()
	if err != nil {
		return err
	}
	serializedResultFragment := SerializeResultFragment(resultFragment)
	serializedResultFragment.From = SerializeMultiAddress(from)
	serializedResultFragment.To = &rpc.Address{Address: address.String()}

	conn, err := Dial(to)
	if err != nil {
		return nil
	}
	defer conn.Close()
	client := rpc.NewNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err = client.SendResultFragment(ctx, serializedResultFragment, grpc.FailFast(false))
	return err
}

func SerializeMultiAddress(multiAddress identity.MultiAddress) *rpc.MultiAddress {
	return &rpc.MultiAddress{Multi: multiAddress.String()}
}

func DeserializeMultiAddress(multiAddress *rpc.MultiAddress) (identity.MultiAddress, error) {
	return identity.NewMultiAddressFromString(multiAddress.Multi)
}

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

// SerializeOrderFragment converts a compute.OrderFragment into its network
// representation.
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
		log.Println("FstCodeShare =", err)
		return nil, err
	}
	orderFragment.SndCodeShare, err = sss.FromBytes(input.SndCodeShare)
	if err != nil {
		log.Println("SndCodeShare =", err)
		return nil, err
	}
	orderFragment.PriceShare, err = sss.FromBytes(input.PriceShare)
	if err != nil {
		log.Println("PriceShare =", err)
		return nil, err
	}
	orderFragment.MaxVolumeShare, err = sss.FromBytes(input.MaxVolumeShare)
	if err != nil {
		log.Println("MaxVolumeShare =", err)
		return nil, err
	}
	orderFragment.MinVolumeShare, err = sss.FromBytes(input.MinVolumeShare)
	if err != nil {
		log.Println("MinVolumeShare =", err)
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
