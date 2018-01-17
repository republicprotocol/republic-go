package compute

import (
	"bytes"
	"encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/go-sss"
)

// An OrderFragmentID is the Keccak256 hash of an OrderFragment.
type OrderFragmentID []byte

// Equals checks if two OrderFragmentIDs are equal in value.
func (id OrderFragmentID) Equals(other OrderFragmentID) bool {
	return bytes.Equal(id, other)
}

// An OrderFragment is a secret share of an Order. Is is created using Shamir
// secret sharing where the secret is an Order encoded as a big.Int.
type OrderFragment struct {
	// Public data.
	ID           OrderFragmentID
	OrderID      OrderID
	OrderType    OrderType
	OrderBuySell OrderBuySell

	// Private data.
	FstCodeShare   sss.Share
	SndCodeShare   sss.Share
	PriceShare     sss.Share
	MaxVolumeShare sss.Share
	MinVolumeShare sss.Share
}

// NewOrderFragment returns a new OrderFragment and computes the
// OrderFragmentID for the OrderFragment.
func NewOrderFragment(orderID OrderID, orderType OrderType, orderBuySell OrderBuySell, fstCodeShare, sndCodeShare, priceShare, maxVolumeShare, minVolumeShare sss.Share) *OrderFragment {
	orderFragment := &OrderFragment{
		OrderID:        orderID,
		OrderType:      orderType,
		OrderBuySell:   orderBuySell,
		FstCodeShare:   fstCodeShare,
		SndCodeShare:   sndCodeShare,
		PriceShare:     priceShare,
		MaxVolumeShare: maxVolumeShare,
		MinVolumeShare: minVolumeShare,
	}
	orderFragment.ID = OrderFragmentID(crypto.Keccak256(orderFragment.Bytes()))
	return orderFragment
}

// Add two OrderFragments together and return the resulting output
// OrderFragment. The output OrderFragment will have its ID computed.
func (orderFragment *OrderFragment) Add(rhs *OrderFragment, prime *big.Int) (*ComputedOrderFragment, error) {
	// Check that the OrderFragments have compatible sss.Shares, and that one
	// of them is an OrderBuy and the other is an OrderSell.
	if err := orderFragment.IsCompatible(rhs); err != nil {
		return nil, err
	}

	// Label the OrderFragments appropriately.
	var buyOrderFragment, sellOrderFragment *OrderFragment
	if orderFragment.OrderBuySell == OrderBuy {
		buyOrderFragment = orderFragment
		sellOrderFragment = rhs
	} else {
		buyOrderFragment = rhs
		sellOrderFragment = orderFragment
	}

	// Perform the addition using the buyOrderFragment as the LHS and the
	// sellOrderFragment as the RHS.
	computed := &ComputedOrderFragment{
		BuyOrderID:          buyOrderFragment.OrderID,
		BuyOrderFragmentID:  buyOrderFragment.ID,
		SellOrderID:         sellOrderFragment.OrderID,
		SellOrderFragmentID: sellOrderFragment.ID,
		FstCodeShare: sss.Share{
			Key:   buyOrderFragment.FstCodeShare.Key,
			Value: big.NewInt(0).Add(buyOrderFragment.FstCodeShare.Value, sellOrderFragment.FstCodeShare.Value),
		},
		SndCodeShare: sss.Share{
			Key:   buyOrderFragment.SndCodeShare.Key,
			Value: big.NewInt(0).Add(buyOrderFragment.SndCodeShare.Value, sellOrderFragment.SndCodeShare.Value),
		},
		PriceShare: sss.Share{
			Key:   buyOrderFragment.PriceShare.Key,
			Value: big.NewInt(0).Add(buyOrderFragment.PriceShare.Value, sellOrderFragment.PriceShare.Value),
		},
		MaxVolumeShare: sss.Share{
			Key:   buyOrderFragment.MaxVolumeShare.Key,
			Value: big.NewInt(0).Add(buyOrderFragment.MaxVolumeShare.Value, sellOrderFragment.MaxVolumeShare.Value),
		},
		MinVolumeShare: sss.Share{
			Key:   buyOrderFragment.MinVolumeShare.Key,
			Value: big.NewInt(0).Add(buyOrderFragment.MinVolumeShare.Value, sellOrderFragment.MinVolumeShare.Value),
		},
	}
	computed.FstCodeShare.Value.Mod(computed.FstCodeShare.Value, prime)
	computed.SndCodeShare.Value.Mod(computed.SndCodeShare.Value, prime)
	computed.PriceShare.Value.Mod(computed.PriceShare.Value, prime)
	computed.MaxVolumeShare.Value.Mod(computed.MaxVolumeShare.Value, prime)
	computed.MinVolumeShare.Value.Mod(computed.MinVolumeShare.Value, prime)
	computed.ID = OrderFragmentID(crypto.Keccak256(computed.Bytes()))
	return computed, nil
}

// Sub two OrderFragments from one another and return the resulting output
// OrderFragment. The output OrderFragment will have its ID computed.
func (orderFragment *OrderFragment) Sub(rhs *OrderFragment, prime *big.Int) (*OrderFragment, error) {
	// Check that the OrderFragments have compatible sss.Shares, and that one
	// of them is an OrderBuy and the other is an OrderSell.
	if err := orderFragment.IsCompatible(rhs); err != nil {
		return nil, err
	}

	// Label the OrderFragments appropriately.
	var buyOrderFragment, sellOrderFragment *OrderFragment
	if orderFragment.OrderBuySell == OrderBuy {
		buyOrderFragment = orderFragment
		sellOrderFragment = rhs
	} else {
		buyOrderFragment = rhs
		sellOrderFragment = orderFragment
	}

	// Perform the addition using the buyOrderFragment as the LHS and the
	// sellOrderFragment as the RHS.
	computed := &ComputedOrderFragment{
		BuyOrderID:          buyOrderFragment.OrderID,
		BuyOrderFragmentID:  buyOrderFragment.ID,
		SellOrderID:         sellOrderFragment.OrderID,
		SellOrderFragmentID: sellOrderFragment.ID,
		FstCodeShare: sss.Share{
			Key:   buyOrderFragment.FstCodeShare.Key,
			Value: big.NewInt(0).Add(buyOrderFragment.FstCodeShare.Value, big.NewInt(0).Set(prime).Sub(prime, sellOrderFragment.FstCodeShare.Value)),
		},
		SndCodeShare: sss.Share{
			Key:   buyOrderFragment.SndCodeShare.Key,
			Value: big.NewInt(0).Add(buyOrderFragment.SndCodeShare.Value, big.NewInt(0).Set(prime).Sub(prime, sellOrderFragment.SndCodeShare.Value)),
		},
		PriceShare: sss.Share{
			Key:   buyOrderFragment.PriceShare.Key,
			Value: big.NewInt(0).Add(buyOrderFragment.PriceShare.Value, big.NewInt(0).Set(prime).Sub(prime, sellOrderFragment.PriceShare.Value)),
		},
		MaxVolumeShare: sss.Share{
			Key:   buyOrderFragment.MaxVolumeShare.Key,
			Value: big.NewInt(0).Add(buyOrderFragment.MaxVolumeShare.Value, big.NewInt(0).Set(prime).Sub(prime, sellOrderFragment.MaxVolumeShare.Value)),
		},
		MinVolumeShare: sss.Share{
			Key:   buyOrderFragment.MinVolumeShare.Key,
			Value: big.NewInt(0).Add(buyOrderFragment.MinVolumeShare.Value, big.NewInt(0).Set(prime).Sub(prime, sellOrderFragment.MinVolumeShare.Value)),
		},
	}
	computed.FstCodeShare.Value.Mod(computed.FstCodeShare.Value, prime)
	computed.SndCodeShare.Value.Mod(computed.SndCodeShare.Value, prime)
	computed.PriceShare.Value.Mod(computed.PriceShare.Value, prime)
	computed.MaxVolumeShare.Value.Mod(computed.MaxVolumeShare.Value, prime)
	computed.MinVolumeShare.Value.Mod(computed.MinVolumeShare.Value, prime)
	computed.ID = OrderFragmentID(crypto.Keccak256(computed.Bytes()))
	return computed, nil
}

// Bytes returns an Order serialized into a bytes.
func (orderFragment *OrderFragment) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, orderFragment.OrderID)
	binary.Write(buf, binary.LittleEndian, orderFragment.OrderType)
	binary.Write(buf, binary.LittleEndian, orderFragment.OrderBuySell)

	binary.Write(buf, binary.LittleEndian, orderFragment.FstCodeShare.Key)
	binary.Write(buf, binary.LittleEndian, orderFragment.FstCodeShare.Value.Bytes())
	binary.Write(buf, binary.LittleEndian, orderFragment.SndCodeShare.Key)
	binary.Write(buf, binary.LittleEndian, orderFragment.SndCodeShare.Value.Bytes())
	binary.Write(buf, binary.LittleEndian, orderFragment.PriceShare.Key)
	binary.Write(buf, binary.LittleEndian, orderFragment.PriceShare.Value.Bytes())
	binary.Write(buf, binary.LittleEndian, orderFragment.MaxVolumeShare.Key)
	binary.Write(buf, binary.LittleEndian, orderFragment.MaxVolumeShare.Value.Bytes())
	binary.Write(buf, binary.LittleEndian, orderFragment.MinVolumeShare.Key)
	binary.Write(buf, binary.LittleEndian, orderFragment.MinVolumeShare.Value.Bytes())

	return buf.Bytes()
}

// IsCompatible returns an error when the two OrderFragments do not have
// the same share indices.
func (orderFragment *OrderFragment) IsCompatible(rhs *OrderFragment) error {
	if orderFragment.OrderBuySell == rhs.OrderBuySell {
		return NewOrderComputationError(orderFragment.OrderBuySell)
	}
	if orderFragment.FstCodeShare.Key != rhs.FstCodeShare.Key {
		return NewOrderFragmentationError(orderFragment.FstCodeShare.Key, rhs.FstCodeShare.Key)
	}
	if orderFragment.SndCodeShare.Key != rhs.SndCodeShare.Key {
		return NewOrderFragmentationError(orderFragment.SndCodeShare.Key, rhs.SndCodeShare.Key)
	}
	if orderFragment.PriceShare.Key != rhs.PriceShare.Key {
		return NewOrderFragmentationError(orderFragment.PriceShare.Key, rhs.PriceShare.Key)
	}
	if orderFragment.MaxVolumeShare.Key != rhs.MaxVolumeShare.Key {
		return NewOrderFragmentationError(orderFragment.MaxVolumeShare.Key, rhs.MaxVolumeShare.Key)
	}
	if orderFragment.MinVolumeShare.Key != rhs.MinVolumeShare.Key {
		return NewOrderFragmentationError(orderFragment.MinVolumeShare.Key, rhs.MinVolumeShare.Key)
	}
	return nil
}

// A ComputedOrderFragment is the result of a computation over two
// OrderFragments.
type ComputedOrderFragment struct {
	// Public data.
	ID                  OrderFragmentID
	BuyOrderID          OrderID
	BuyOrderFragmentID  OrderFragmentID
	SellOrderID         OrderID
	SellOrderFragmentID OrderFragmentID

	// Private data.
	FstCodeShare   sss.Share
	SndCodeShare   sss.Share
	PriceShare     sss.Share
	MaxVolumeShare sss.Share
	MinVolumeShare sss.Share
}

// Bytes returns an Order serialized into a bytes.
func (computed *ComputedOrderFragment) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, computed.BuyOrderID)
	binary.Write(buf, binary.LittleEndian, computed.BuyOrderFragmentID)
	binary.Write(buf, binary.LittleEndian, computed.SellOrderID)
	binary.Write(buf, binary.LittleEndian, computed.SellOrderFragmentID)

	binary.Write(buf, binary.LittleEndian, computed.FstCodeShare.Key)
	binary.Write(buf, binary.LittleEndian, computed.FstCodeShare.Value.Bytes())
	binary.Write(buf, binary.LittleEndian, computed.SndCodeShare.Key)
	binary.Write(buf, binary.LittleEndian, computed.SndCodeShare.Value.Bytes())
	binary.Write(buf, binary.LittleEndian, computed.PriceShare.Key)
	binary.Write(buf, binary.LittleEndian, computed.PriceShare.Value.Bytes())
	binary.Write(buf, binary.LittleEndian, computed.MaxVolumeShare.Key)
	binary.Write(buf, binary.LittleEndian, computed.MaxVolumeShare.Value.Bytes())
	binary.Write(buf, binary.LittleEndian, computed.MinVolumeShare.Key)
	binary.Write(buf, binary.LittleEndian, computed.MinVolumeShare.Value.Bytes())

	return buf.Bytes()
}
