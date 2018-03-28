package rpc

import (
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/smpc"
)

// MarshalMultiAddress into an RPC protobuf object.
func MarshalMultiAddress(multiAddress *identity.MultiAddress) *MultiAddress {
	return &MultiAddress{
		Signature:    []byte{},
		MultiAddress: multiAddress.String(),
	}
}

// UnmarshalMultiAddress from an RPC protobuf object.
func UnmarshalMultiAddress(multiAddress *MultiAddress) (identity.MultiAddress, error) {
	return identity.NewMultiAddressFromString(multiAddress.MultiAddress)
}

func MarshalDeltaFragment(deltaFragment *smpc.DeltaFragment) *DeltaFragment {
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

func UnmarshalDeltaFragment(deltaFragment *DeltaFragment) (smpc.DeltaFragment, error) {
	val := smpc.DeltaFragment{
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
		return smpc.DeltaFragment{}, err
	}
	val.SndCodeShare, err = shamir.FromBytes(deltaFragment.SndCodeShare)
	if err != nil {
		return smpc.DeltaFragment{}, err
	}
	val.PriceShare, err = shamir.FromBytes(deltaFragment.PriceShare)
	if err != nil {
		return smpc.DeltaFragment{}, err
	}
	val.MaxVolumeShare, err = shamir.FromBytes(deltaFragment.MaxVolumeShare)
	if err != nil {
		return smpc.DeltaFragment{}, err
	}
	val.MinVolumeShare, err = shamir.FromBytes(deltaFragment.MinVolumeShare)
	if err != nil {
		return smpc.DeltaFragment{}, err
	}
	return val, nil
}

func MarshalDeltaFragments(deltaFragments smpc.DeltaFragments) *DeltaFragments {
	val := make([]*DeltaFragment, len(deltaFragments))
	for i := range deltaFragments {
		val[i] = MarshalDeltaFragment(&deltaFragments[i])
	}
	return &DeltaFragments{
		DeltaFragments: val,
	}
}

func UnmarshalDeltaFragments(deltaFragments *DeltaFragments) (smpc.DeltaFragments, error) {
	val := make(smpc.DeltaFragments, 0, len(deltaFragments.DeltaFragments))
	for i := range deltaFragments.DeltaFragments {
		deltaFragment, err := UnmarshalDeltaFragment(deltaFragments.DeltaFragments[i])
		if err != nil {
			return val, err
		}
		val = append(val, deltaFragment)
	}
	return val, nil
}
