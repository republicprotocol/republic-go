package smpcer

import (
	"fmt"
	"time"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/delta"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
)

// MarshalOrderFragment converts an order.Fragment into its network
// representation.
func MarshalOrderFragment(addr string, crypter crypto.Crypter, orderFragment *order.Fragment) (*OrderFragment, error) {
	val := &OrderFragment{
		OrderFragmentId: orderFragment.ID,
		OrderId:         orderFragment.OrderID,
		Expiry:          orderFragment.OrderExpiry.Unix(),
		Parity:          int64(orderFragment.OrderParity),
		Type:            int64(orderFragment.OrderType),
	}

	var err error
	val.FstCodeShare, err = crypter.Encrypt(addr, shamir.ToBytes(orderFragment.FstCodeShare))
	if err != nil {
		return nil, fmt.Errorf("cannot encrypt fst code: %v", err)
	}
	val.SndCodeShare, err = crypter.Encrypt(addr, shamir.ToBytes(orderFragment.SndCodeShare))
	if err != nil {
		return nil, fmt.Errorf("cannot encrypt snd code: %v", err)
	}
	val.PriceShare, err = crypter.Encrypt(addr, shamir.ToBytes(orderFragment.PriceShare))
	if err != nil {
		return nil, fmt.Errorf("cannot encrypt price share: %v", err)
	}
	val.MaxVolumeShare, err = crypter.Encrypt(addr, shamir.ToBytes(orderFragment.MaxVolumeShare)) // FIXME: Unify volumes
	if err != nil {
		return nil, fmt.Errorf("cannot encrypt max volume share: %v", err)
	}
	val.MinVolumeShare, err = crypter.Encrypt(addr, shamir.ToBytes(orderFragment.MinVolumeShare)) // FIXME: Unify volumes
	if err != nil {
		return nil, fmt.Errorf("cannot encrypt min volume: %v", err)
	}

	return val, nil
}

// UnmarshalOrderFragment converts a network representation of an
// OrderFragment into an order.Fragment. An error is returned if the network
// representation is malformed.
func UnmarshalOrderFragment(crypter crypto.Crypter, orderFragment *OrderFragment) (order.Fragment, error) {
	var err error

	val := order.Fragment{
		Signature: []byte{}, // FIXME: Verify signature from the RPC order fragment
		ID:        orderFragment.GetOrderFragmentId(),

		OrderID:     order.ID(orderFragment.GetOrderId()),
		OrderType:   order.Type(orderFragment.Type),
		OrderParity: order.Parity(orderFragment.Parity),
		OrderExpiry: time.Unix(orderFragment.Expiry, 0),
	}

	fstCodeShare, err := crypter.Decrypt(orderFragment.FstCodeShare)
	if err != nil {
		return val, err
	}
	val.FstCodeShare, err = shamir.FromBytes(fstCodeShare)
	if err != nil {
		return val, err
	}

	sndCodeShare, err := crypter.Decrypt(orderFragment.SndCodeShare)
	if err != nil {
		return val, err
	}
	val.SndCodeShare, err = shamir.FromBytes(sndCodeShare)
	if err != nil {
		return val, err
	}

	priceShare, err := crypter.Decrypt(orderFragment.PriceShare)
	if err != nil {
		return val, err
	}
	val.PriceShare, err = shamir.FromBytes(priceShare)
	if err != nil {
		return val, err
	}

	maxVolumeShare, err := crypter.Decrypt(orderFragment.MaxVolumeShare)
	if err != nil {
		return val, err
	}
	val.MaxVolumeShare, err = shamir.FromBytes(maxVolumeShare)
	if err != nil {
		return val, err
	}

	minVolumeShare, err := crypter.Decrypt(orderFragment.MinVolumeShare)
	if err != nil {
		return val, err
	}
	val.MinVolumeShare, err = shamir.FromBytes(minVolumeShare)
	if err != nil {
		return val, err
	}

	return val, nil
}

// MarshalDeltaFragment into a RPC protobuf object.
func MarshalDeltaFragment(deltaFragment *delta.Fragment) *DeltaFragment {
	return &DeltaFragment{
		DeltaFragmentId:     deltaFragment.ID,
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

// UnmarshalDeltaFragment from a RPC protobuf object.
func UnmarshalDeltaFragment(deltaFragment *DeltaFragment) (delta.Fragment, error) {
	val := delta.Fragment{
		ID:                  deltaFragment.DeltaFragmentId,
		DeltaID:             deltaFragment.DeltaId,
		BuyOrderID:          deltaFragment.BuyOrderId,
		SellOrderID:         deltaFragment.SellOrderId,
		BuyOrderFragmentID:  deltaFragment.BuyOrderFragmentId,
		SellOrderFragmentID: deltaFragment.SellOrderFragmentId,
	}
	var err error
	val.FstCodeShare, err = shamir.FromBytes(deltaFragment.FstCodeShare)
	if err != nil {
		return val, err
	}
	val.SndCodeShare, err = shamir.FromBytes(deltaFragment.SndCodeShare)
	if err != nil {
		return val, err
	}
	val.PriceShare, err = shamir.FromBytes(deltaFragment.PriceShare)
	if err != nil {
		return val, err
	}
	val.MaxVolumeShare, err = shamir.FromBytes(deltaFragment.MaxVolumeShare)
	if err != nil {
		return val, err
	}
	val.MinVolumeShare, err = shamir.FromBytes(deltaFragment.MinVolumeShare)
	if err != nil {
		return val, err
	}
	return val, nil
}
