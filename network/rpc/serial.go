package rpc

import (
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
)

// SerializeAddress converts an identity.MultiAddress into its network
// representation.
func SerializeAddress(address identity.Address) *Address {
	return &Address{Address: address.String()}
}

// DeserializeAddress converts a network representation of an Address into an
// an identity.Address. An error is returned if the network representation is
// malformed.
func DeserializeAddress(address *Address) identity.Address {
	return identity.Address(address.Address)
}

// SerializeMultiAddress converts an identity.MultiAddress into its network
// representation.
func SerializeMultiAddress(multiAddress identity.MultiAddress) *MultiAddress {
	return &MultiAddress{MultiAddress: multiAddress.String()}
}

// DeserializeMultiAddress converts a network representation of a MultiAddress
// into an identity.MultiAddress. An error is returned if the network
// representation is malformed.
func DeserializeMultiAddress(multiAddress *MultiAddress) (identity.MultiAddress, error) {
	return identity.NewMultiAddressFromString(multiAddress.MultiAddress)
}

// SerializeOrderFragment converts an order.Fragment into its network
// representation.
func SerializeOrderFragment(orderFragment *order.Fragment) *OrderFragment {
	val := &OrderFragment{
		Id:          []byte(orderFragment.ID),
		OrderId:     []byte(orderFragment.OrderID),
		OrderType:   int64(orderFragment.OrderType),
		OrderParity: int64(orderFragment.OrderParity),
	}
	val.FstCodeShare = shamir.ToBytes(orderFragment.FstCodeShare)
	val.SndCodeShare = shamir.ToBytes(orderFragment.SndCodeShare)
	val.PriceShare = shamir.ToBytes(orderFragment.PriceShare)
	val.MaxVolumeShare = shamir.ToBytes(orderFragment.MaxVolumeShare)
	val.MinVolumeShare = shamir.ToBytes(orderFragment.MinVolumeShare)
	return val
}

// DeserializeOrderFragment converts a network representation of an
// OrderFragment into an order.Fragment. An error is returned if the network
// representation is malformed.
func DeserializeOrderFragment(orderFragment *OrderFragment) (*order.Fragment, error) {
	val := &order.Fragment{
		ID:          order.FragmentID(orderFragment.Id),
		OrderID:     order.ID(orderFragment.OrderId),
		OrderType:   order.Type(orderFragment.OrderType),
		OrderParity: order.Parity(orderFragment.OrderParity),
	}
	var err error
	val.FstCodeShare, err = shamir.FromBytes(orderFragment.FstCodeShare)
	if err != nil {
		return nil, err
	}
	val.SndCodeShare, err = shamir.FromBytes(orderFragment.SndCodeShare)
	if err != nil {
		return nil, err
	}
	val.PriceShare, err = shamir.FromBytes(orderFragment.PriceShare)
	if err != nil {
		return nil, err
	}
	val.MaxVolumeShare, err = shamir.FromBytes(orderFragment.MaxVolumeShare)
	if err != nil {
		return nil, err
	}
	val.MinVolumeShare, err = shamir.FromBytes(orderFragment.MinVolumeShare)
	if err != nil {
		return nil, err
	}
	return val, nil
}
