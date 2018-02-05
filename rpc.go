package rpc

import (
	"fmt"

	"context"
	"time"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-sss"
	"google.golang.org/grpc"
)

// Dial the target identity.MultiAddress using a background context.Context.
// Returns a grpc.ClientConn, or an error. The grpc.ClientConn must be closed
// before it exits scope.
func Dial(target identity.MultiAddress, timeout time.Duration) (*grpc.ClientConn, error) {
	host, err := target.ValueForProtocol(identity.IP4Code)
	if err != nil {
		return nil, err
	}
	port, err := target.ValueForProtocol(identity.TCPCode)
	if err != nil {
		return nil, err
	}
	timeoutContext, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	conn, err := grpc.DialContext(timeoutContext, fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// SerializeAddress converts an identity.MultiAddress into its network
// representation.
func SerializeAddress(address identity.Address) *Address {
	return &Address{Address: address.String()}
}

// DeserializeAddress converts a network representation of an Address into an
// an identity.Address. An error is returned if the network representation is
// malformed.
func DeserializeAddress(address *Address) (identity.Address, error) {
	return identity.Address(address.Address), nil
}

// SerializeMultiAddress converts an identity.MultiAddress into its network
// representation.
func SerializeMultiAddress(multiAddress identity.MultiAddress) *MultiAddress {
	return &MultiAddress{Multi: multiAddress.String()}
}

// DeserializeMultiAddress converts a network representation of a MultiAddress
// into an identity.MultiAddress. An error is returned if the network
// representation is malformed.
func DeserializeMultiAddress(multiAddress *MultiAddress) (identity.MultiAddress, error) {
	return identity.NewMultiAddressFromString(multiAddress.Multi)
}

// SerializeMultiAddresses converts identity.MultiAddresses into their network
// representation.
func SerializeMultiAddresses(multiAddresses identity.MultiAddresses) *MultiAddresses {
	serializedMultiAddresses := make([]*MultiAddress, len(multiAddresses))
	for i, multiAddress := range multiAddresses {
		serializedMultiAddresses[i] = SerializeMultiAddress(multiAddress)
	}
	return &MultiAddresses{Multis: serializedMultiAddresses}
}

// DeserializeMultiAddresses converts a network representation of
// MultiAddresses into identity.MultiAddresses. An error is returned if the
// network representation is malformed.
func DeserializeMultiAddresses(multiAddresses *MultiAddresses) (identity.MultiAddresses, error) {
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
func SerializeOrderFragment(input *compute.OrderFragment) *OrderFragment {
	orderFragment := &OrderFragment{
		Id:          []byte(input.ID),
		OrderId:     []byte(input.OrderID),
		OrderType:   int64(input.OrderType),
		OrderParity: int64(input.OrderParity),
	}
	orderFragment.FstCodeShare = sss.ToBytes(input.FstCodeShare)
	orderFragment.SndCodeShare = sss.ToBytes(input.SndCodeShare)
	orderFragment.PriceShare = sss.ToBytes(input.PriceShare)
	orderFragment.MaxVolumeShare = sss.ToBytes(input.MaxVolumeShare)
	orderFragment.MinVolumeShare = sss.ToBytes(input.MinVolumeShare)
	return orderFragment
}

// DeserializeOrderFragment converts a network representation of an
// OrderFragment into a compute.OrderFragment. An error is returned if the
// network representation is malformed.
func DeserializeOrderFragment(input *OrderFragment) (*compute.OrderFragment, error) {
	orderFragment := &compute.OrderFragment{
		ID:          compute.OrderFragmentID(input.Id),
		OrderID:     compute.OrderID(input.OrderId),
		OrderType:   compute.OrderType(input.OrderType),
		OrderParity: compute.OrderParity(input.OrderParity),
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
func SerializeResultFragment(input *compute.ResultFragment) *ResultFragment {
	resultFragment := &ResultFragment{
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
func DeserializeResultFragment(input *ResultFragment) (*compute.ResultFragment, error) {
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

// SerializeTradingAtom converts an atomic.TradingAtom into its network
// representation.
func SerializeTradingAtom(tradingAtom struct{}) *TradingAtom {
	return &TradingAtom{}
}

// DeserializeTradingAtom converts a network representation of a TradingAtom
// into an atomic.TradingAtom. An error is returned if the network
// representation is malformed.
func DeserializeTradingAtom(tradingAtom *TradingAtom) (struct{}, error) {
	return struct{}{}, nil
}
