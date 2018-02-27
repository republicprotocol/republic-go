package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/republicprotocol/go-atom"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-sss"
	"google.golang.org/grpc"
)

// Dial the identity.MultiAddress using a background context.Context. Returns a
// grpc.ClientConn, or an error. The grpc.ClientConn must be closed before it
// exits scope.
func Dial(multiAddress identity.MultiAddress, timeout time.Duration) (*grpc.ClientConn, error) {
	host, err := multiAddress.ValueForProtocol(identity.IP4Code)
	if err != nil {
		return nil, err
	}
	port, err := multiAddress.ValueForProtocol(identity.TCPCode)
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
func DeserializeAddress(address *Address) identity.Address {
	return identity.Address(address.Address)
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
		To:          &Address{Address: ""},
		From:        &MultiAddress{Multi: ""},
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

//// SerializeFinal converts a compute.Final into its network representation.
//func SerializeFinal(input *compute.) *Final {
//	result := &Final{
//		Id:          []byte(input.ID),
//		BuyOrderId:  []byte(input.BuyOrderID),
//		SellOrderId: []byte(input.SellOrderID),
//	}
//
//	result.FstCode = input.FstCode.Bytes()
//	result.SndCode = input.SndCode.Bytes()
//	result.Price = input.Price.Bytes()
//	result.MaxVolume = input.MaxVolume.Bytes()
//	result.MinVolume = input.MinVolume.Bytes()
//	return result
//}
//
//// DeserializeFinal converts a network representation of a Final into a
//// compute.FinalFragment. An error is returned if the network representation
//// is malformed.
//func DeserializeFinal(input *Final) *compute.Final {
//	result := &compute.Final{
//		ID:          compute.FinalID(input.Id),
//		BuyOrderID:  []byte(input.BuyOrderId),
//		SellOrderID: []byte(input.SellOrderId),
//	}
//	result.FstCode = big.NewInt(0).SetBytes(input.FstCode)
//	result.SndCode = big.NewInt(0).SetBytes(input.SndCode)
//	result.Price = big.NewInt(0).SetBytes(input.Price)
//	result.MaxVolume = big.NewInt(0).SetBytes(input.MaxVolume)
//	result.MinVolume = big.NewInt(0).SetBytes(input.MinVolume)
//	return result
//}

// SerializeAtom converts an atomic.Atom into its network representation.
func SerializeAtom(a atom.Atom) *Atom {
	return &Atom{
		Ledger:    int64(a.Ledger),
		Data:      a.LedgerData,
		Signature: a.Signature,
	}
}

// DeserializeAtom converts a network representation of an Atom into an
// atom.Atom. An error is returned if the network representation is malformed.
func DeserializeAtom(a *Atom) atom.Atom {
	return atom.Atom{
		Ledger:     atom.Ledger(a.Ledger),
		LedgerData: a.Data,
		Signature:  a.Signature,
	}
}

func SerializeShard(shard compute.Shard) *Shard {
	serializedDelats := make([]*DeltaFragment, len(shard.Deltas))
	for i ,j  := range shard.Deltas{
		serializedDelats[i] = SerializeDeltaFragment(j)
	}
	//todo
	//serializedResidues := make([]*ResidueFragment, len(shard.Deltas))
	//for i ,j  := range shard.Deltas{
	//	serializedResidues[i] = SerializeReidueFragment(j)
	//}

	return &Shard{
		Signature: shard.Signature,
		Deltas:    serializedDelats,
		Residues:  [][]byte{},
	}
}

func DeserializeShard(shard *Shard) (*compute.Shard, error) {
	var err error
	deltas := make ([]*compute.DeltaFragment, len(shard.Deltas))
	for i,j := range shard.Deltas{
		deltas[i], err = DeserializeDeltaFragment(j)
		if err != nil {
			return nil, err
		}
	}
	return &compute.Shard{
		Signature: shard.Signature,
		Deltas:    deltas,
		Residues:  []*compute.ResidueFragment{},//todo
	},nil
}

// todo: serialize the deltas and residues to bytes
func SerializeFinalShard(shard compute.DeltaShard) *DeltaShard {
	return &DeltaShard{
		Signature: shard.Signature,
		Finals:    [][]byte{},
	}
}

// todo: deserialize deltas and residues
func DeserializeFinalShard(shard *DeltaShard) *compute.DeltaShard {
	return &compute.DeltaShard{
		Signature: shard.Signature,
		DeltaFragments:    []*compute.DeltaFragment{},
	}
}

func SerializeDeltaFragment(fragment *compute.DeltaFragment) *DeltaFragment{
	return &DeltaFragment{
		Id:  fragment.ID,
		BuyOrderId:fragment.BuyOrderID,
		SellOrderId: fragment.SellOrderID,
		BuyOrderFragmentId: fragment.BuyOrderFragmentID,
		SellOrderFragmentId: fragment.SellOrderFragmentID,
		FstCodeShare: sss.ToBytes(fragment.FstCodeShare),
		SndCodeShare: sss.ToBytes(fragment.SndCodeShare),
		PriceShare: sss.ToBytes(fragment.PriceShare),
		MaxVolumeShare: sss.ToBytes(fragment.MaxVolumeShare),
		MinVolumeShare: sss.ToBytes(fragment.MinVolumeShare),
	}
}

func DeserializeDeltaFragment(fragment *DeltaFragment) (*compute.DeltaFragment ,error){
	deltaFragment := &compute.DeltaFragment{
		ID:  fragment.Id,
		BuyOrderID:fragment.BuyOrderId,
		SellOrderID: fragment.SellOrderId,
		BuyOrderFragmentID: fragment.BuyOrderFragmentId,
		SellOrderFragmentID: fragment.SellOrderFragmentId,
	}
	var err error
	deltaFragment.FstCodeShare , err = sss.FromBytes(fragment.FstCodeShare)
	if err != nil {
		return nil, err
	}
	deltaFragment.SndCodeShare, err = sss.FromBytes(fragment.SndCodeShare)
	if err != nil {
		return nil, err
	}
	deltaFragment.PriceShare, err = sss.FromBytes(fragment.PriceShare)
	if err != nil {
		return nil, err
	}
	deltaFragment.MaxVolumeShare, err = sss.FromBytes(fragment.MaxVolumeShare)
	if err != nil {
		return nil, err
	}
	deltaFragment.MinVolumeShare, err = sss.FromBytes(fragment.MinVolumeShare)
	if err != nil {
		return nil, err
	}
	return deltaFragment, nil
}