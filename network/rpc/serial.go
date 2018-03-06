package rpc

import (
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/shamir"
)

// SerializeAddress converts an identity.Multiaddress into its network
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

// SerializeMultiaddress converts an identity.Multiaddress into its network
// representation.
func SerializeMultiaddress(multiAddress identity.Multiaddress) *Multiaddress {
	return &Multiaddress{Multiaddress: multiAddress.String()}
}

// DeserializeMultiaddress converts a network representation of a Multiaddress
// into an identity.Multiaddress. An error is returned if the network
// representation is malformed.
func DeserializeMultiaddress(multiAddress *Multiaddress) (identity.Multiaddress, error) {
	return identity.NewMultiaddressFromString(multiAddress.Multiaddress)
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
	orderFragment.FstCodeShare = shamir.ToBytes(input.FstCodeShare)
	orderFragment.SndCodeShare = shamir.ToBytes(input.SndCodeShare)
	orderFragment.PriceShare = shamir.ToBytes(input.PriceShare)
	orderFragment.MaxVolumeShare = shamir.ToBytes(input.MaxVolumeShare)
	orderFragment.MinVolumeShare = shamir.ToBytes(input.MinVolumeShare)
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
	orderFragment.FstCodeShare, err = shamir.FromBytes(input.FstCodeShare)
	if err != nil {
		return nil, err
	}
	orderFragment.SndCodeShare, err = shamir.FromBytes(input.SndCodeShare)
	if err != nil {
		return nil, err
	}
	orderFragment.PriceShare, err = shamir.FromBytes(input.PriceShare)
	if err != nil {
		return nil, err
	}
	orderFragment.MaxVolumeShare, err = shamir.FromBytes(input.MaxVolumeShare)
	if err != nil {
		return nil, err
	}
	orderFragment.MinVolumeShare, err = shamir.FromBytes(input.MinVolumeShare)
	if err != nil {
		return nil, err
	}
	return orderFragment, nil
}

func SerializeDeltaFragment(fragment *compute.DeltaFragment) *DeltaFragment {
	return &DeltaFragment{
		Id:                  fragment.ID,
		BuyOrderID:          fragment.BuyOrderID,
		SellOrderID:         fragment.SellOrderID,
		BuyOrderFragmentID:  fragment.BuyOrderFragmentID,
		SellOrderFragmentID: fragment.SellOrderFragmentID,
		FstCodeShare:        shamir.ToBytes(fragment.FstCodeShare),
		SndCodeShare:        shamir.ToBytes(fragment.SndCodeShare),
		PriceShare:          shamir.ToBytes(fragment.PriceShare),
		MaxVolumeShare:      shamir.ToBytes(fragment.MaxVolumeShare),
		MinVolumeShare:      shamir.ToBytes(fragment.MinVolumeShare),
	}
}

func DeserializeDeltaFragment(fragment *DeltaFragment) (*compute.DeltaFragment, error) {
	deltaFragment := &compute.DeltaFragment{
		ID:                  fragment.Id,
		BuyOrderID:          fragment.BuyOrderID,
		SellOrderID:         fragment.SellOrderID,
		BuyOrderFragmentID:  fragment.BuyOrderFragmentID,
		SellOrderFragmentID: fragment.SellOrderFragmentID,
	}
	var err error
	deltaFragment.FstCodeShare, err = shamir.FromBytes(fragment.FstCodeShare)
	if err != nil {
		return nil, err
	}
	deltaFragment.SndCodeShare, err = shamir.FromBytes(fragment.SndCodeShare)
	if err != nil {
		return nil, err
	}
	deltaFragment.PriceShare, err = shamir.FromBytes(fragment.PriceShare)
	if err != nil {
		return nil, err
	}
	deltaFragment.MaxVolumeShare, err = shamir.FromBytes(fragment.MaxVolumeShare)
	if err != nil {
		return nil, err
	}
	deltaFragment.MinVolumeShare, err = shamir.FromBytes(fragment.MinVolumeShare)
	if err != nil {
		return nil, err
	}
	return deltaFragment, nil
}
