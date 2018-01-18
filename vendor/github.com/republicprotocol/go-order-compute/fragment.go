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

func (id OrderFragmentID) Equals(other OrderFragmentID) bool {
	return bytes.Equal(id, other)
}

// An OrderFragment is a secret share of an Order. Is is created using Shamir
// secret sharing where the secret is an Order encoded as a big.Int.
type OrderFragment struct {
	// Public data.
	ID       OrderFragmentID
	OrderIDs OrderIDs

	// Private data.
	FstCodeShare   sss.Share
	SndCodeShare   sss.Share
	PriceShare     sss.Share
	MaxVolumeShare sss.Share
	MinVolumeShare sss.Share
}

// NewOrderFragment returns a new OrderFragment and computes the
// OrderFragmentID for the OrderFragment.
func NewOrderFragment(orderIDs OrderIDs, fstCodeShare, sndCodeShare, priceShare, maxVolumeShare, minVolumeShare sss.Share) *OrderFragment {
	orderFragment := &OrderFragment{
		OrderIDs:       orderIDs,
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
func (orderFragment *OrderFragment) Add(rhs *OrderFragment, prime *big.Int) (*OrderFragment, error) {
	if err := orderFragment.IsCompatible(rhs); err != nil {
		return nil, err
	}
	outputFragment := &OrderFragment{
		OrderIDs: append(orderFragment.OrderIDs, rhs.OrderIDs...),
		FstCodeShare: sss.Share{
			Key:   orderFragment.FstCodeShare.Key,
			Value: big.NewInt(0).Add(orderFragment.FstCodeShare.Value, rhs.FstCodeShare.Value),
		},
		SndCodeShare: sss.Share{
			Key:   orderFragment.SndCodeShare.Key,
			Value: big.NewInt(0).Add(orderFragment.SndCodeShare.Value, rhs.SndCodeShare.Value),
		},
		PriceShare: sss.Share{
			Key:   orderFragment.PriceShare.Key,
			Value: big.NewInt(0).Add(orderFragment.PriceShare.Value, rhs.PriceShare.Value),
		},
		MaxVolumeShare: sss.Share{
			Key:   orderFragment.MaxVolumeShare.Key,
			Value: big.NewInt(0).Add(orderFragment.MaxVolumeShare.Value, rhs.MaxVolumeShare.Value),
		},
		MinVolumeShare: sss.Share{
			Key:   orderFragment.MinVolumeShare.Key,
			Value: big.NewInt(0).Add(orderFragment.MinVolumeShare.Value, rhs.MinVolumeShare.Value),
		},
	}
	outputFragment.FstCodeShare.Value.Mod(outputFragment.FstCodeShare.Value, prime)
	outputFragment.SndCodeShare.Value.Mod(outputFragment.SndCodeShare.Value, prime)
	outputFragment.PriceShare.Value.Mod(outputFragment.PriceShare.Value, prime)
	outputFragment.MaxVolumeShare.Value.Mod(outputFragment.MaxVolumeShare.Value, prime)
	outputFragment.MinVolumeShare.Value.Mod(outputFragment.MinVolumeShare.Value, prime)
	orderFragment.ID = OrderFragmentID(crypto.Keccak256(orderFragment.Bytes()))
	return outputFragment, nil
}

// Sub two OrderFragments from one another and return the resulting output
// OrderFragment. The output OrderFragment will have its ID computed.
func (orderFragment *OrderFragment) Sub(rhs *OrderFragment, prime *big.Int) (*OrderFragment, error) {
	if err := orderFragment.IsCompatible(rhs); err != nil {
		return nil, err
	}
	outputFragment := &OrderFragment{
		OrderIDs: append(orderFragment.OrderIDs, rhs.OrderIDs...),
		FstCodeShare: sss.Share{
			Key:   orderFragment.FstCodeShare.Key,
			Value: big.NewInt(0).Add(orderFragment.FstCodeShare.Value, big.NewInt(0).Set(prime).Sub(prime, rhs.FstCodeShare.Value)),
		},
		SndCodeShare: sss.Share{
			Key:   orderFragment.SndCodeShare.Key,
			Value: big.NewInt(0).Add(orderFragment.SndCodeShare.Value, big.NewInt(0).Set(prime).Sub(prime, rhs.SndCodeShare.Value)),
		},
		PriceShare: sss.Share{
			Key:   orderFragment.PriceShare.Key,
			Value: big.NewInt(0).Add(orderFragment.PriceShare.Value, big.NewInt(0).Set(prime).Sub(prime, rhs.PriceShare.Value)),
		},
		MaxVolumeShare: sss.Share{
			Key:   orderFragment.MaxVolumeShare.Key,
			Value: big.NewInt(0).Add(orderFragment.MaxVolumeShare.Value, big.NewInt(0).Set(prime).Sub(prime, rhs.MaxVolumeShare.Value)),
		},
		MinVolumeShare: sss.Share{
			Key:   orderFragment.MinVolumeShare.Key,
			Value: big.NewInt(0).Add(orderFragment.MinVolumeShare.Value, big.NewInt(0).Set(prime).Sub(prime, rhs.MinVolumeShare.Value)),
		},
	}
	outputFragment.FstCodeShare.Value.Mod(outputFragment.FstCodeShare.Value, prime)
	outputFragment.SndCodeShare.Value.Mod(outputFragment.SndCodeShare.Value, prime)
	outputFragment.PriceShare.Value.Mod(outputFragment.PriceShare.Value, prime)
	outputFragment.MaxVolumeShare.Value.Mod(outputFragment.MaxVolumeShare.Value, prime)
	outputFragment.MinVolumeShare.Value.Mod(outputFragment.MinVolumeShare.Value, prime)
	orderFragment.ID = OrderFragmentID(crypto.Keccak256(orderFragment.Bytes()))
	return outputFragment, nil
}

// Bytes returns an Order serialized into a bytes.
func (orderFragment *OrderFragment) Bytes() []byte {
	buf := new(bytes.Buffer)
	for _, orderID := range orderFragment.OrderIDs {
		binary.Write(buf, binary.LittleEndian, orderID)
	}
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
