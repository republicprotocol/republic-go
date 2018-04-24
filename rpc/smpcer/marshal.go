package smpcer

import (
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/delta"
	"github.com/republicprotocol/republic-go/order"
)

// MarshalOrderFragment converts an order.Fragment into its network
// representation.
func MarshalOrderFragment(crypter crypto.Crypter, orderFragment *order.Fragment) (*OrderFragment, error) {
	val := &OrderFragment{
		OrderFragmenId: orderFragment.ID,
		OrderId:        orderFragment.OrderID,
		Expiry:         orderFragment.OrderExpiry.Unix(),
		Type:           int32(orderFragment.OrderType),
		Tokens:         int32(0), // FIXME: Real token pairing
		PriceShare: &Share{
			Index: orderFragment.PriceShare.Key,
		},
		VolumeShare: &Share{
			Index: orderFragment.MaxVolumeShare.Key, // FIXME: Unify volumes
		},
	}

	var err error
	val.PriceShare.Value, err = crypter.Encrypt(orderFragment.PriceShare.Value.Bytes())
	if err != nil {
		return nil, err
	}
	val.VolumeShare.Value, err = crypter.Encrypt(orderFragment.MaxVolumeShare.Value.Bytes()) // FIXME: Unify volumes
	if err != nil {
		return nil, err
	}

	return val, nil
}

// UnmarshalOrderFragment converts a network representation of an
// OrderFragment into an order.Fragment. An error is returned if the network
// representation is malformed.
func UnmarshalOrderFragment(orderFragment *OrderFragment) (order.Fragment, error) {
	// var err error

	// val := &order.Fragment{
	// 	Signature: orderFragment.Id.Signature,
	// 	ID:        orderFragment.Id.OrderFragmentId,

	// 	OrderID:     order.ID(orderFragment.Order.Id.OrderId),
	// 	OrderType:   order.Type(orderFragment.Order.Type),
	// 	OrderParity: order.Parity(orderFragment.Order.Parity),
	// 	OrderExpiry: time.Unix(orderFragment.Order.Expiry, 0),
	// }

	// fstCodeShare, err := crypto.Decrypt(privKey, orderFragment.FstCodeShare)
	// if err != nil {
	// 	return nil, err
	// }
	// val.FstCodeShare, err = shamir.FromBytes(fstCodeShare)
	// if err != nil {
	// 	return nil, err
	// }

	// sndCodeShare, err := crypto.Decrypt(privKey, orderFragment.SndCodeShare)
	// if err != nil {
	// 	return nil, err
	// }
	// val.SndCodeShare, err = shamir.FromBytes(sndCodeShare)
	// if err != nil {
	// 	return nil, err
	// }

	// priceShare, err := crypto.Decrypt(privKey, orderFragment.PriceShare)
	// if err != nil {
	// 	return nil, err
	// }
	// val.PriceShare, err = shamir.FromBytes(priceShare)
	// if err != nil {
	// 	return nil, err
	// }

	// maxVolumeShare, err := crypto.Decrypt(privKey, orderFragment.MaxVolumeShare)
	// if err != nil {
	// 	return nil, err
	// }
	// val.MaxVolumeShare, err = shamir.FromBytes(maxVolumeShare)
	// if err != nil {
	// 	return nil, err
	// }

	// minVolumeShare, err := crypto.Decrypt(privKey, orderFragment.MinVolumeShare)
	// if err != nil {
	// 	return nil, err
	// }
	// val.MinVolumeShare, err = shamir.FromBytes(minVolumeShare)
	// if err != nil {
	// 	return nil, err
	// }

	// return val, nil
	panic("unimplemented")
}

// MarshalDeltaFragment into a RPC protobuf object.
func MarshalDeltaFragment(deltaFragment *delta.Fragment) *DeltaFragment {
	return &DeltaFragment{
		DeltaFragmentId: deltaFragment.ID,
		BuyOrderId:      deltaFragment.BuyOrderID,
		SellOrderId:     deltaFragment.SellOrderID,
		PriceShare: &Share{
			Index: deltaFragment.PriceShare.Key,
			Value: deltaFragment.PriceShare.Value.Bytes(),
		},
		VolumeShare: &Share{
			Index: deltaFragment.MaxVolumeShare.Key,           // FIXME: Unify volumes
			Value: deltaFragment.MaxVolumeShare.Value.Bytes(), // FIXME: Unify volumes
		},
	}
}

// UnmarshalDeltaFragment from a RPC protobuf object.
func UnmarshalDeltaFragment(deltaFragment *DeltaFragment) (delta.Fragment, error) {
	// val := smpc.DeltaFragment{
	// 	ID:                  deltaFragment.Id,
	// 	DeltaID:             deltaFragment.DeltaId,
	// 	BuyOrderID:          deltaFragment.BuyOrderId,
	// 	SellOrderID:         deltaFragment.SellOrderId,
	// 	BuyOrderFragmentID:  deltaFragment.BuyOrderFragmentId,
	// 	SellOrderFragmentID: deltaFragment.SellOrderFragmentId,
	// }
	// var err error
	// val.FstCodeShare, err = shamir.FromBytes(deltaFragment.FstCodeShare)
	// if err != nil {
	// 	return smpc.DeltaFragment{}, err
	// }
	// val.SndCodeShare, err = shamir.FromBytes(deltaFragment.SndCodeShare)
	// if err != nil {
	// 	return smpc.DeltaFragment{}, err
	// }
	// val.PriceShare, err = shamir.FromBytes(deltaFragment.PriceShare)
	// if err != nil {
	// 	return smpc.DeltaFragment{}, err
	// }
	// val.MaxVolumeShare, err = shamir.FromBytes(deltaFragment.MaxVolumeShare)
	// if err != nil {
	// 	return smpc.DeltaFragment{}, err
	// }
	// val.MinVolumeShare, err = shamir.FromBytes(deltaFragment.MinVolumeShare)
	// if err != nil {
	// 	return smpc.DeltaFragment{}, err
	// }
	// return val, nil
	panic("unimplemented")
}
