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

// An OrderFragment is a secret share of an Order. Is is created using Shamir
// secret sharing where the secret is an Order encoded as a big.Int.
type OrderFragment struct {
	// Public data.
	ID       OrderFragmentID
	OrderIDs OrderIDs

	// Private data.
	LittleCodeShare sss.Share
	BigCodeShare    sss.Share
	PriceShare      sss.Share
	MaxVolumeShare  sss.Share
	MinVolumeShare  sss.Share
}

// Add two OrderFragments together and return the resulting output
// OrderFragment. The output OrderFragment will have its ID computed.
func (orderFragment *OrderFragment) Add(rhs *OrderFragment, prime *big.Int) (*OrderFragment, error) {
	if err := orderFragment.IsCompatible(rhs); err != nil {
		return nil, err
	}
	outputFragment := &OrderFragment{
		OrderIDs: append(orderFragment.OrderIDs, rhs.OrderIDs...),
		LittleCodeShare: sss.Share{
			Key:   orderFragment.LittleCodeShare.Key,
			Value: big.NewInt(0).Add(orderFragment.LittleCodeShare.Value, rhs.LittleCodeShare.Value),
		},
		BigCodeShare: sss.Share{
			Key:   orderFragment.BigCodeShare.Key,
			Value: big.NewInt(0).Add(orderFragment.BigCodeShare.Value, rhs.BigCodeShare.Value),
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
	outputFragment.LittleCodeShare.Value.Mod(outputFragment.LittleCodeShare.Value, prime)
	outputFragment.BigCodeShare.Value.Mod(outputFragment.BigCodeShare.Value, prime)
	outputFragment.PriceShare.Value.Mod(outputFragment.PriceShare.Value, prime)
	outputFragment.MaxVolumeShare.Value.Mod(outputFragment.MaxVolumeShare.Value, prime)
	outputFragment.MinVolumeShare.Value.Mod(outputFragment.MinVolumeShare.Value, prime)
	outputFragment.ID = OrderFragmentID(crypto.Keccak256(outputFragment.Bytes()))
	return outputFragment, nil
}

// Bytes returns an Order serialized into a bytes.
func (orderFragment *OrderFragment) Bytes() []byte {
	buf := new(bytes.Buffer)
	for _, orderID := range orderFragment.OrderIDs {
		binary.Write(buf, binary.LittleEndian, orderID)
	}
	binary.Write(buf, binary.LittleEndian, orderFragment.LittleCodeShare.Key)
	binary.Write(buf, binary.LittleEndian, orderFragment.LittleCodeShare.Value)
	binary.Write(buf, binary.LittleEndian, orderFragment.BigCodeShare.Key)
	binary.Write(buf, binary.LittleEndian, orderFragment.BigCodeShare.Value)
	binary.Write(buf, binary.LittleEndian, orderFragment.PriceShare.Key)
	binary.Write(buf, binary.LittleEndian, orderFragment.PriceShare.Value)
	binary.Write(buf, binary.LittleEndian, orderFragment.MaxVolumeShare.Key)
	binary.Write(buf, binary.LittleEndian, orderFragment.MaxVolumeShare.Value)
	binary.Write(buf, binary.LittleEndian, orderFragment.MinVolumeShare.Key)
	binary.Write(buf, binary.LittleEndian, orderFragment.MinVolumeShare.Value)
	return buf.Bytes()
}

// IsCompatible returns an error when the two OrderFragments do not have
// the same share indices.
func (orderFragment *OrderFragment) IsCompatible(rhs *OrderFragment) error {
	if orderFragment.LittleCodeShare.Key != rhs.LittleCodeShare.Key {
		return NewOrderFragmentationError(orderFragment.LittleCodeShare.Key, rhs.LittleCodeShare.Key)
	}
	if orderFragment.BigCodeShare.Key != rhs.LittleCodeShare.Key {
		return NewOrderFragmentationError(orderFragment.BigCodeShare.Key, rhs.BigCodeShare.Key)
	}
	if orderFragment.PriceShare.Key != rhs.LittleCodeShare.Key {
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
