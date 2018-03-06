package order

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	sss "github.com/republicprotocol/go-sss"
	"github.com/republicprotocol/republic-go/identity"
)

// An FragmentID is the Keccak256 hash of an OrderFragment.
type FragmentID []byte

// Equals checks if two OrderFragmentIDs are equal in value.
func (id OrderFragmentID) Equals(other OrderFragmentID) bool {
	return bytes.Equal(id, other)
}

func (id OrderFragmentID) String() string {
	return string(id)
}

// An OrderFragment is a secret share of an Order. Is is created using Shamir
// secret sharing where the secret is an Order encoded as a big.Int.
type OrderFragment struct {
	// Signature
	Owner     identity.ID
	Signature []byte

	// Public
	ID          OrderFragmentID
	OrderID     OrderID
	OrderType   OrderType
	OrderParity OrderParity
	OrderExpiry time.Time

	// Secure
	FstCodeShare   sss.Share
	SndCodeShare   sss.Share
	PriceShare     sss.Share
	MaxVolumeShare sss.Share
	MinVolumeShare sss.Share
}

// NewOrderFragment returns a new OrderFragment and computes the
// OrderFragmentID for the OrderFragment.
func NewOrderFragment(orderID OrderID, orderType OrderType, orderParity OrderParity, fstCodeShare, sndCodeShare, priceShare, maxVolumeShare, minVolumeShare sss.Share) *OrderFragment {
	orderFragment := &OrderFragment{
		OrderID:        orderID,
		OrderType:      orderType,
		OrderParity:    orderParity,
		FstCodeShare:   fstCodeShare,
		SndCodeShare:   sndCodeShare,
		PriceShare:     priceShare,
		MaxVolumeShare: maxVolumeShare,
		MinVolumeShare: minVolumeShare,
	}
	orderFragment.ID = OrderFragmentID(crypto.Keccak256(orderFragment.Bytes()))
	return orderFragment
}

// Sub two OrderFragments from one another and return the resulting output
// ResultFragment. The output ResultFragment will have its ID computed.
func (orderFragment *OrderFragment) Sub(other *OrderFragment, prime *big.Int) (*DeltaFragment, error) {
	// Check that the OrderFragments have compatible sss.Shares, and that one
	// of them is an OrderBuy and the other is an OrderSell.
	if err := orderFragment.IsCompatible(other); err != nil {
		return nil, err
	}

	// Label the OrderFragments appropriately.
	var buyOrderFragment, sellOrderFragment *OrderFragment
	if orderFragment.OrderParity == OrderParityBuy {
		buyOrderFragment = orderFragment
		sellOrderFragment = other
	} else {
		buyOrderFragment = other
		sellOrderFragment = orderFragment
	}

	// Perform the addition using the buyOrderFragment as the LHS and the
	// sellOrderFragment as the RHS.
	fstCodeShare := sss.Share{
		Key:   buyOrderFragment.FstCodeShare.Key,
		Value: big.NewInt(0).Add(buyOrderFragment.FstCodeShare.Value, big.NewInt(0).Sub(prime, sellOrderFragment.FstCodeShare.Value)),
	}
	sndCodeShare := sss.Share{
		Key:   buyOrderFragment.SndCodeShare.Key,
		Value: big.NewInt(0).Add(buyOrderFragment.SndCodeShare.Value, big.NewInt(0).Sub(prime, sellOrderFragment.SndCodeShare.Value)),
	}
	priceShare := sss.Share{
		Key:   buyOrderFragment.PriceShare.Key,
		Value: big.NewInt(0).Add(buyOrderFragment.PriceShare.Value, big.NewInt(0).Sub(prime, sellOrderFragment.PriceShare.Value)),
	}
	maxVolumeShare := sss.Share{
		Key:   buyOrderFragment.MaxVolumeShare.Key,
		Value: big.NewInt(0).Add(buyOrderFragment.MaxVolumeShare.Value, big.NewInt(0).Sub(prime, sellOrderFragment.MinVolumeShare.Value)),
	}
	minVolumeShare := sss.Share{
		Key:   buyOrderFragment.MinVolumeShare.Key,
		Value: big.NewInt(0).Add(sellOrderFragment.MaxVolumeShare.Value, big.NewInt(0).Sub(prime, buyOrderFragment.MinVolumeShare.Value)),
	}
	fstCodeShare.Value.Mod(fstCodeShare.Value, prime)
	sndCodeShare.Value.Mod(sndCodeShare.Value, prime)
	priceShare.Value.Mod(priceShare.Value, prime)
	maxVolumeShare.Value.Mod(maxVolumeShare.Value, prime)
	minVolumeShare.Value.Mod(minVolumeShare.Value, prime)
	deltaFragment := &DeltaFragment{
		nil,
		buyOrderFragment.OrderID,
		sellOrderFragment.OrderID,
		buyOrderFragment.ID,
		sellOrderFragment.ID,
		fstCodeShare,
		sndCodeShare,
		priceShare,
		maxVolumeShare,
		minVolumeShare,
	}
	deltaFragment.ID = DeltaFragmentID(crypto.Keccak256(deltaFragment.BuyOrderFragmentID[:], deltaFragment.SellOrderFragmentID[:]))
	return deltaFragment, nil
}

// Bytes returns an OrderFragment serialized into a bytes.
func (orderFragment *OrderFragment) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, orderFragment.OrderID)
	binary.Write(buf, binary.LittleEndian, orderFragment.OrderType)
	binary.Write(buf, binary.LittleEndian, orderFragment.OrderParity)

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
	if orderFragment.OrderParity == rhs.OrderParity {
		return NewOrderParityError(orderFragment.OrderParity)
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
