package rpc

import (
	"time"

	"github.com/republicprotocol/republic-go/compute"
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

func SerializeDeltaFragment(deltaFragment *compute.DeltaFragment) *DeltaFragment {
	return &DeltaFragment{
		Id:                  deltaFragment.ID,
		DeltaId:             deltaFragment.DeltaID,
		BuyOrderId:          deltaFragment.BuyOrderID,
		SellOrderId:         deltaFragment.SellOrderID,
		BuyOrderFragmentId:  deltaFragment.BuyOrderFragmentID,
		SellOrderFragmentId: deltaFragment.SellOrderFragmentID,
		FstCodeShare:        shamir.ToBytes(deltaFragment.FstCodeShare),
		SndCodeShare:        shamir.ToBytes(deltaFragment.SndCodeShare),
		PriceShare:          shamir.ToBytes(deltaFragment.PriceShare),
		MaxVolumeShare:      shamir.ToBytes(deltaFragment.MaxVolumeShare),
		MinVolumeShare:      shamir.ToBytes(deltaFragment.MinVolumeShare),
	}
}

func DeserializeDeltaFragment(deltaFragment *DeltaFragment) (*compute.DeltaFragment, error) {
	val := &compute.DeltaFragment{
		ID:                  deltaFragment.Id,
		DeltaID:             deltaFragment.DeltaId,
		BuyOrderID:          deltaFragment.BuyOrderId,
		SellOrderID:         deltaFragment.SellOrderId,
		BuyOrderFragmentID:  deltaFragment.BuyOrderFragmentId,
		SellOrderFragmentID: deltaFragment.SellOrderFragmentId,
	}
	var err error
	val.FstCodeShare, err = shamir.FromBytes(deltaFragment.FstCodeShare)
	if err != nil {
		return nil, err
	}
	val.SndCodeShare, err = shamir.FromBytes(deltaFragment.SndCodeShare)
	if err != nil {
		return nil, err
	}
	val.PriceShare, err = shamir.FromBytes(deltaFragment.PriceShare)
	if err != nil {
		return nil, err
	}
	val.MaxVolumeShare, err = shamir.FromBytes(deltaFragment.MaxVolumeShare)
	if err != nil {
		return nil, err
	}
	val.MinVolumeShare, err = shamir.FromBytes(deltaFragment.MinVolumeShare)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func SerializeOrder(ord *order.Order) *Order {
	orderRpc := new(Order)
	orderRpc.Id = ord.ID
	orderRpc.Signature = ord.Signature
	orderRpc.Type = int64(ord.Type)
	orderRpc.Parity = int64(ord.Parity)
	orderRpc.Expiry = ord.Expiry.Unix()

	return orderRpc
}

func DeserializeOrder(orderRpc *Order) *order.Order {
	ord := new(order.Order)
	ord.ID = orderRpc.Id
	ord.Signature = orderRpc.Signature
	ord.Type = order.Type(orderRpc.Type)
	ord.Parity = order.Parity(orderRpc.Parity)
	ord.Expiry = time.Unix(orderRpc.Expiry, 0)

	return ord
}
