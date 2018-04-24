package rpc

import (
	"crypto/rsa"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/rpc/smpcer"
	"github.com/republicprotocol/republic-go/smpc"
)

// MarshalOrderFragment converts an order.Fragment into its network
// representation.
func MarshalOrderFragment(pubKey *rsa.PublicKey, orderFragment *order.Fragment) (*smpcer.OrderFragment, error) {
	// val := &smpcer.OrderFragment{
	// 	Id: &OrderFragmentId{
	// 		Signature:       orderFragment.Signature,
	// 		OrderFragmentId: orderFragment.ID,
	// 	},
	// 	Order: &Order{
	// 		Id: &OrderId{
	// 			Signature: orderFragment.Signature,
	// 			OrderId:   orderFragment.OrderID,
	// 		},
	// 		Type:   int64(orderFragment.OrderType),
	// 		Parity: int64(orderFragment.OrderParity),
	// 		Expiry: orderFragment.OrderExpiry.Unix(),
	// 	},
	// }

	// var err error
	// val.FstCodeShare, err = crypto.Encrypt(pubKey, shamir.ToBytes(orderFragment.FstCodeShare))
	// if err != nil {
	// 	return nil, err
	// }
	// val.SndCodeShare, err = crypto.Encrypt(pubKey, shamir.ToBytes(orderFragment.SndCodeShare))
	// if err != nil {
	// 	return nil, err
	// }
	// val.PriceShare, err = crypto.Encrypt(pubKey, shamir.ToBytes(orderFragment.PriceShare))
	// if err != nil {
	// 	return nil, err
	// }
	// val.MaxVolumeShare, err = crypto.Encrypt(pubKey, shamir.ToBytes(orderFragment.MaxVolumeShare))
	// if err != nil {
	// 	return nil, err
	// }
	// val.MinVolumeShare, err = crypto.Encrypt(pubKey, shamir.ToBytes(orderFragment.MinVolumeShare))
	// if err != nil {
	// 	return nil, err
	// }

	// return val, nil
	panic("unimplemented")
}

// UnmarshalOrderFragment converts a network representation of an
// OrderFragment into an order.Fragment. An error is returned if the network
// representation is malformed.
func UnmarshalOrderFragment(privKey *rsa.PrivateKey, orderFragment *smpcer.OrderFragment) (*order.Fragment, error) {
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
func MarshalDeltaFragment(deltaFragment *smpc.DeltaFragment) *smpcer.DeltaFragment {
	// return &smpcer.DeltaFragment{
	// 	Id:                  deltaFragment.ID,
	// 	DeltaId:             deltaFragment.DeltaID,
	// 	BuyOrderId:          deltaFragment.BuyOrderID,
	// 	SellOrderId:         deltaFragment.SellOrderID,
	// 	BuyOrderFragmentId:  deltaFragment.BuyOrderFragmentID,
	// 	SellOrderFragmentId: deltaFragment.SellOrderFragmentID,
	// 	FstCodeShare:        shamir.ToBytes(deltaFragment.FstCodeShare),
	// 	SndCodeShare:        shamir.ToBytes(deltaFragment.SndCodeShare),
	// 	PriceShare:          shamir.ToBytes(deltaFragment.PriceShare),
	// 	MaxVolumeShare:      shamir.ToBytes(deltaFragment.MaxVolumeShare),
	// 	MinVolumeShare:      shamir.ToBytes(deltaFragment.MinVolumeShare),
	// }
	panic("unimplemented")
}

// UnmarshalDeltaFragment from a RPC protobuf object.
func UnmarshalDeltaFragment(deltaFragment *smpcer.DeltaFragment) (smpc.DeltaFragment, error) {
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
