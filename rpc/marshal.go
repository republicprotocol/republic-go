package rpc

import (
	"time"

	"github.com/republicprotocol/republic-go/hyperdrive"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/smpc"
)

// MarshalAddress into a RPC protobuf object.
func MarshalAddress(address identity.Address) *Address {
	return &Address{
		Address: address.String(),
	}
}

// UnmarshalAddress into a RPC protobuf object.
func UnmarshalAddress(address *Address) identity.Address {
	return identity.Address(address.Address)
}

// MarshalMultiAddress into a RPC protobuf object.
func MarshalMultiAddress(multiAddress *identity.MultiAddress) *MultiAddress {
	return &MultiAddress{
		Signature:    []byte{},
		MultiAddress: multiAddress.String(),
	}
}

// UnmarshalMultiAddress from a RPC protobuf object.
func UnmarshalMultiAddress(multiAddress *MultiAddress) (identity.MultiAddress, []byte, error) {
	multi, err := identity.NewMultiAddressFromString(multiAddress.MultiAddress)
	return multi, multiAddress.Signature, err
}

// MarshalOrderFragment converts an order.Fragment into its network
// representation.
func MarshalOrderFragment(orderFragment *order.Fragment) *OrderFragment {
	val := &OrderFragment{
		Id: &OrderFragmentId{
			Signature:       orderFragment.Signature,
			OrderFragmentId: orderFragment.ID,
		},
		Order: &Order{
			Id: &OrderId{
				Signature: orderFragment.Signature,
				OrderId:   orderFragment.OrderID,
			},
			Type:   int64(orderFragment.OrderType),
			Parity: int64(orderFragment.OrderParity),
			Expiry: orderFragment.OrderExpiry.Unix(),
		},
		FstCodeShare:   shamir.ToBytes(orderFragment.FstCodeShare),
		SndCodeShare:   shamir.ToBytes(orderFragment.SndCodeShare),
		PriceShare:     shamir.ToBytes(orderFragment.PriceShare),
		MaxVolumeShare: shamir.ToBytes(orderFragment.MaxVolumeShare),
		MinVolumeShare: shamir.ToBytes(orderFragment.MinVolumeShare),
	}

	return val
}

// UnmarshalOrderFragment converts a network representation of an
// OrderFragment into an order.Fragment. An error is returned if the network
// representation is malformed.
func UnmarshalOrderFragment(orderFragment *OrderFragment) (*order.Fragment, error) {
	var err error

	val := &order.Fragment{
		Signature: orderFragment.Id.Signature,
		ID:        orderFragment.Id.OrderFragmentId,

		OrderID:     order.ID(orderFragment.Order.Id.OrderId),
		OrderType:   order.Type(orderFragment.Order.Type),
		OrderParity: order.Parity(orderFragment.Order.Parity),
		OrderExpiry: time.Unix(orderFragment.Order.Expiry, 0),
	}

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

// MarshalDeltaFragment into a RPC protobuf object.
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

// UnmarshalDeltaFragment from a RPC protobuf object.
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

//// MarshalDeltaFragments into an RPC protobuf object.
//func MarshalDeltaFragments(deltaFragments smpc.DeltaFragments) *smpc.DeltaFragments {
//	val := make([]*DeltaFragment, len(deltaFragments))
//	for i := range deltaFragments {
//		val[i] = MarshalDeltaFragment(&deltaFragments[i])
//	}
//	return &DeltaFragments{
//		DeltaFragments: val,
//	}
//}
//
//// UnmarshalDeltaFragments from an RPC protobuf object.
//func UnmarshalDeltaFragments(deltaFragments *DeltaFragments) (smpc.DeltaFragments, error) {
//	val := make(smpc.DeltaFragments, 0, len(deltaFragments.DeltaFragments))
//	for i := range deltaFragments.DeltaFragments {
//		deltaFragment, err := UnmarshalDeltaFragment(deltaFragments.DeltaFragments[i])
//		if err != nil {
//			return val, err
//		}
//		val = append(val, deltaFragment)
//	}
//	return val, nil
//}

// MarshalOrder into an RPC protobuf object
func MarshalOrder(ord *order.Order) *Order {
	rpcOrder := new(Order)
	rpcOrder.Id = &OrderId{
		OrderId:   ord.ID,
		Signature: ord.Signature,
	}
	rpcOrder.Type = int64(ord.Type)
	rpcOrder.Parity = int64(ord.Parity)
	rpcOrder.Expiry = ord.Expiry.Unix()

	return rpcOrder
}

// UnmarshalOrder from an RPC protobuf object.
func UnmarshalOrder(rpcOrder *Order) order.Order {
	ord := order.Order{}
	ord.ID = rpcOrder.Id.OrderId
	ord.Signature = rpcOrder.Id.Signature
	ord.Type = order.Type(rpcOrder.Type)
	ord.Parity = order.Parity(rpcOrder.Parity)
	ord.Expiry = time.Unix(rpcOrder.Expiry, 0)

	return ord
}

func MarshalTx(tx hyperdrive.Tx) *Tx {
	nonces := make([]*Nonce, len(tx.Nonces))
	for i := range nonces {
		nonces[i] = &Nonce{
			Nonce: tx.Nonces[i][:],
		}
	}

	return &Tx{
		Hash:   tx.Hash[:],
		Nonces: nonces,
	}
}

func UnmarshalTx(tx *Tx) hyperdrive.Tx {
	var hash identity.Hash
	copy(hash[:], tx.Hash)

	nonces := make(hyperdrive.Nonces, len(tx.Nonces))
	for i := range nonces {
		copy(nonces[i][:], tx.Nonces[i].Nonce)
	}

	return hyperdrive.Tx{
		Hash:   hash,
		Nonces: nonces,
	}
}

func MarshalBlock(block hyperdrive.Block) *Block {
	txs := make([]*Tx, len(block.Txs))
	for i := range block.Txs {
		txs[i] = MarshalTx(block.Txs[i])
	}

	return &Block{
		Height:    int64(block.Height),
		Rank:      int64(block.Rank),
		EpocHash:  block.Epoch[:],
		Txs:       txs,
		Signature: block.Signature[:],
	}
}

func UnmarshalBlock(block *Block) *hyperdrive.Block {
	txs := make([]hyperdrive.Tx, len(block.Txs))
	for i := range block.Txs {
		txs[i] = UnmarshalTx(block.Txs[i])
	}

	var epochHash [32]byte
	copy(epochHash[:], block.EpocHash)
	var signature [65]byte
	copy(signature[:], block.Signature)

	return &hyperdrive.Block{
		Height:    hyperdrive.Height(block.Height),
		Rank:      hyperdrive.Rank(block.Rank),
		Epoch:     epochHash,
		Txs:       txs,
		Signature: signature,
	}
}

func MarshalProposal(tx hyperdrive.Proposal) *Proposal {
	return &Proposal{
		Block:     MarshalBlock(tx.Block),
		Signature: tx.Signature[:],
	}
}

func UnmarshalProposal(proposal *Proposal) *hyperdrive.Proposal {
	var signature [65]byte
	copy(signature[:], proposal.Signature)
	return &hyperdrive.Proposal{
		Block:     *UnmarshalBlock(proposal.Block),
		Signature: signature,
	}
}

func MarshalPrepare(prepare hyperdrive.Prepare) *Prepare {
	signatures := Signatures{Signature: make([][]byte, len(prepare.Signature))}
	for i := range signatures.Signature {
		signatures.Signature[i] = prepare.Signatures[i][:]
	}
	return &Prepare{
		Proposal:   MarshalProposal(prepare.Proposal),
		Signatures: &signatures,
	}
}

func UnmarshalPrepare(prepare *Prepare) *hyperdrive.Prepare {
	signature := make(identity.Signatures, len(prepare.Signatures.Signature))
	for i := range signature {
		signature[i] = [65]byte{}
		copy(signature[i][:], prepare.Signatures.Signature[i])
	}

	return &hyperdrive.Prepare{
		Proposal:   *UnmarshalProposal(prepare.Proposal),
		Signatures: signature,
	}
}
