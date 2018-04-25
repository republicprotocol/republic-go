package order

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/shamir"
)

// An FragmentID is the Keccak256 hash of a Fragment.
type FragmentID []byte

// Equal returns an equality check between two FragmentIDs.
func (id FragmentID) Equal(other FragmentID) bool {
	return bytes.Equal(id, other)
}

// String returns a FragmentID as a Base58 encoded string.
func (id FragmentID) String() string {
	return base58.Encode(id)
}

// A Fragment is a secret share of an Order, created using Shamir's secret
// sharing on the secure fields in an Order.
type Fragment struct {
	Signature identity.Signature `json:"signature"`
	ID        FragmentID         `json:"id"`

	OrderID     ID        `json:"orderID"`
	OrderType   Type      `json:"orderType"`
	OrderParity Parity    `json:"orderParity"`
	OrderExpiry time.Time `json:"orderExpiry"`

	FstCodeShare   shamir.Share `json:"fstCodeShare"`
	SndCodeShare   shamir.Share `json:"sndCodeShare"`
	PriceShare     shamir.Share `json:"priceShare"`
	MaxVolumeShare shamir.Share `json:"maxVolumeShare"`
	MinVolumeShare shamir.Share `json:"minVolumeShare"`
}

// NewFragment returns a new Fragment and computes the FragmentID.
func NewFragment(orderID ID, orderType Type, orderParity Parity, fstCodeShare, sndCodeShare, priceShare, maxVolumeShare, minVolumeShare shamir.Share) *Fragment {
	fragment := &Fragment{
		OrderID:        orderID,
		OrderType:      orderType,
		OrderParity:    orderParity,
		FstCodeShare:   fstCodeShare,
		SndCodeShare:   sndCodeShare,
		PriceShare:     priceShare,
		MaxVolumeShare: maxVolumeShare,
		MinVolumeShare: minVolumeShare,
	}
	fragment.ID = FragmentID(fragment.Hash())
	return fragment
}

// Hash returns the Keccak256 hash of a Fragment. This hash is used to create
// the FragmentID and signature for a Fragment.
func (fragment *Fragment) Hash() []byte {
	return crypto.Keccak256(fragment.Bytes())
}

// Sign signs the fragment using the provided keypair, and assigns it the the fragments's
// Signature field.
func (fragment *Fragment) Sign(keyPair identity.KeyPair) error {
	var err error
	fragment.Signature, err = keyPair.Sign(fragment)
	return err
}

// SignFragments maps over an array of fragments, calling Sign on each one
func SignFragments(keyPair identity.KeyPair, fragments []*Fragment) error {
	for _, fragment := range fragments {
		if err := fragment.Sign(keyPair); err != nil {
			return err
		}
	}
	return nil
}

// VerifySignature verifies that the Signature field has been signed by the provided
// ID's private key, returning an error if the signature is invalid
func (fragment *Fragment) VerifySignature(ID identity.ID) error {
	return identity.VerifySignature(fragment, fragment.Signature, ID)
}

// VerifyFragmentSignatures maps over an array of fragments,
// calling VerifySignature on each one
func VerifyFragmentSignatures(ID identity.ID, fragments []*Fragment) error {
	for _, fragment := range fragments {
		if err := fragment.VerifySignature(ID); err != nil {
			return err
		}
	}
	return nil
}

// Bytes returns a Fragment serialized into a bytes.
func (fragment *Fragment) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, fragment.OrderID)
	binary.Write(buf, binary.LittleEndian, fragment.OrderType)
	binary.Write(buf, binary.LittleEndian, fragment.OrderParity)
	binary.Write(buf, binary.LittleEndian, fragment.OrderExpiry)

	binary.Write(buf, binary.LittleEndian, fragment.FstCodeShare.Key)
	binary.Write(buf, binary.LittleEndian, fragment.FstCodeShare.Value.Bytes())
	binary.Write(buf, binary.LittleEndian, fragment.SndCodeShare.Key)
	binary.Write(buf, binary.LittleEndian, fragment.SndCodeShare.Value.Bytes())
	binary.Write(buf, binary.LittleEndian, fragment.PriceShare.Key)
	binary.Write(buf, binary.LittleEndian, fragment.PriceShare.Value.Bytes())
	binary.Write(buf, binary.LittleEndian, fragment.MaxVolumeShare.Key)
	binary.Write(buf, binary.LittleEndian, fragment.MaxVolumeShare.Value.Bytes())
	binary.Write(buf, binary.LittleEndian, fragment.MinVolumeShare.Key)
	binary.Write(buf, binary.LittleEndian, fragment.MinVolumeShare.Value.Bytes())

	return buf.Bytes()
}

// Equal returns an equality check between two Orders.
func (fragment *Fragment) Equal(other *Fragment) bool {
	return fragment.ID.Equal(other.ID) &&
		fragment.OrderID.Equal(other.OrderID) &&
		fragment.OrderType == other.OrderType &&
		fragment.OrderParity == other.OrderParity &&
		fragment.OrderExpiry.Equal(other.OrderExpiry) &&
		fragment.FstCodeShare.Value.Cmp(&other.FstCodeShare.Value) == 0 &&
		fragment.SndCodeShare.Value.Cmp(&other.SndCodeShare.Value) == 0 &&
		fragment.PriceShare.Value.Cmp(&other.PriceShare.Value) == 0 &&
		fragment.MaxVolumeShare.Value.Cmp(&other.MaxVolumeShare.Value) == 0 &&
		fragment.MinVolumeShare.Value.Cmp(&other.MinVolumeShare.Value) == 0
}

// IsCompatible returns true when two Fragments are compatible for a
// computation, otherwise it returns false. For a Fragment to be compatible
// with another Fragment it must have a diferrent ID, it must have a different
// order ID, it must have a different parity, it must have a different owner,
// and all secret sharing fields must have the same secret sharing index.
func (fragment *Fragment) IsCompatible(other *Fragment) bool {
	// TODO: Check that signatories are different
	return !fragment.ID.Equal(other.ID) &&
		!fragment.OrderID.Equal(other.OrderID) &&
		fragment.OrderParity != other.OrderParity &&
		fragment.FstCodeShare.Key == other.FstCodeShare.Key &&
		fragment.SndCodeShare.Key == other.SndCodeShare.Key &&
		fragment.PriceShare.Key == other.PriceShare.Key &&
		fragment.MaxVolumeShare.Key == other.MaxVolumeShare.Key &&
		fragment.MinVolumeShare.Key == other.MinVolumeShare.Key
}
