package order

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/shamir"
)

// An FragmentID is the Keccak256 hash of a Fragment.
type FragmentID [32]byte

// A Fragment is a secret share of an Order, created using Shamir's secret
// sharing on the secure fields in an Order.
type Fragment struct {
	OrderID       ID            `json:"orderID"`
	OrderType     Type          `json:"orderType"`
	OrderParity   Parity        `json:"orderParity"`
	OrderExpiry   time.Time     `json:"orderExpiry"`
	ID            FragmentID    `json:"id"`
	Tokens        shamir.Share  `json:"tokens"`
	Price         FragmentValue `json:"price"`
	Volume        FragmentValue `json:"volume"`
	MinimumVolume FragmentValue `json:"minimumVolume"`
}

// NewFragment returns a new Fragment and computes the FragmentID.
func NewFragment(orderID ID, orderType Type, orderParity Parity, tokens shamir.Share, price, volume, minimumVolume FragmentValue) Fragment {
	fragment := Fragment{
		OrderID:       orderID,
		OrderType:     orderType,
		OrderParity:   orderParity,
		Tokens:        tokens,
		Price:         price,
		Volume:        volume,
		MinimumVolume: minimumVolume,
	}
	fragment.ID = FragmentID(fragment.Hash())
	return fragment
}

// Hash returns the Keccak256 hash of a Fragment. This hash is used to create
// the FragmentID and signature for a Fragment.
func (fragment *Fragment) Hash() [32]byte {
	hash := crypto.Keccak256(fragment.Bytes())
	hash32 := [32]byte{}
	copy(hash32[:], hash)
	return hash32
}

// Bytes returns a Fragment serialized into a bytes.
func (fragment *Fragment) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, fragment.OrderID)
	binary.Write(buf, binary.BigEndian, fragment.OrderType)
	binary.Write(buf, binary.BigEndian, fragment.OrderParity)
	binary.Write(buf, binary.BigEndian, fragment.OrderExpiry)
	binary.Write(buf, binary.BigEndian, fragment.Tokens)
	binary.Write(buf, binary.BigEndian, fragment.Price)
	binary.Write(buf, binary.BigEndian, fragment.Volume)
	binary.Write(buf, binary.BigEndian, fragment.MinimumVolume)
	return buf.Bytes()
}

// Equal returns an equality check between two Orders.
func (fragment *Fragment) Equal(other *Fragment) bool {
	return bytes.Equal(fragment.OrderID[:], other.OrderID[:]) &&
		fragment.OrderType == other.OrderType &&
		fragment.OrderParity == other.OrderParity &&
		fragment.OrderExpiry.Equal(other.OrderExpiry) &&
		bytes.Equal(fragment.ID[:], other.ID[:]) &&
		fragment.Tokens.Equal(&other.Tokens) &&
		fragment.Price.Equal(&other.Price) &&
		fragment.Volume.Equal(&other.Volume) &&
		fragment.MinimumVolume.Equal(&other.MinimumVolume)
}

// IsCompatible returns true when two Fragments are compatible for a
// computation, otherwise it returns false. For a Fragment to be compatible
// with another Fragment it must have a diferrent ID, it must have a different
// order ID, it must have a different parity, it must have a different owner,
// and all secret sharing fields must have the same secret sharing index.
func (fragment *Fragment) IsCompatible(other *Fragment) bool {
	// TODO: Check that signatories are different
	return !bytes.Equal(fragment.ID[:], other.ID[:]) &&
		!bytes.Equal(fragment.OrderID[:], other.OrderID[:]) &&
		fragment.OrderParity != other.OrderParity &&
		fragment.Tokens.Index == other.Tokens.Index &&
		fragment.Price.Co.Index == other.Price.Co.Index &&
		fragment.Price.Exp.Index == other.Price.Exp.Index &&
		fragment.Volume.Co.Index == other.Volume.Co.Index &&
		fragment.Volume.Exp.Index == other.Volume.Exp.Index &&
		fragment.MinimumVolume.Co.Index == other.MinimumVolume.Co.Index &&
		fragment.MinimumVolume.Exp.Index == other.MinimumVolume.Exp.Index
}
