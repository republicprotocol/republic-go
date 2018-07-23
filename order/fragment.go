package order

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"time"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/shamir"
)

// An FragmentID is the Keccak256 hash of a Fragment.
type FragmentID [32]byte

// Equal returns an equality check between two FragmentIDs.
func (id FragmentID) Equal(other FragmentID) bool {
	return bytes.Equal(id[:], other[:])
}

// String returns a truncated base64 encoding of the FragmentID.
func (id FragmentID) String() string {
	return base64.StdEncoding.EncodeToString(id[:8])
}

// FragmentEpochDepth is the number of epochs in the passed for which a
// Fragment has been built.
type FragmentEpochDepth uint32

// A Fragment is a secret share of an Order, created using Shamir's secret
// sharing on the secure fields in an Order.
type Fragment struct {
	OrderID         ID                 `json:"orderID"`
	OrderType       Type               `json:"orderType"`
	OrderParity     Parity             `json:"orderParity"`
	OrderSettlement Settlement         `json:"orderSettlement"`
	OrderExpiry     time.Time          `json:"orderExpiry"`
	ID              FragmentID         `json:"id"`
	EpochDepth      FragmentEpochDepth `json:"epochDepth"`

	Tokens        shamir.Share `json:"tokens"`
	Price         CoExpShare   `json:"price"`
	Volume        CoExpShare   `json:"volume"`
	MinimumVolume CoExpShare   `json:"minimumVolume"`
	Nonce         shamir.Share `json:"nonce"`
}

// NewFragment returns a new Fragment and computes the FragmentID.
func NewFragment(orderID ID, orderType Type, orderParity Parity, orderSettlement Settlement, orderExpiry time.Time, tokens shamir.Share, price, volume, minimumVolume CoExpShare, nonce shamir.Share) (Fragment, error) {
	fragment := Fragment{
		OrderID:         orderID,
		OrderType:       orderType,
		OrderParity:     orderParity,
		OrderSettlement: orderSettlement,
		OrderExpiry:     orderExpiry,

		Tokens:        tokens,
		Price:         price,
		Volume:        volume,
		MinimumVolume: minimumVolume,
		Nonce:         nonce,
	}
	fragmentHash, err := fragment.Hash()
	if err != nil {
		return Fragment{}, err
	}

	fragment.ID = FragmentID(fragmentHash)
	return fragment, nil
}

// Hash returns the Keccak256 hash of a Fragment. This hash is used to create
// the FragmentID and signature for a Fragment.
func (fragment *Fragment) Hash() ([32]byte, error) {
	fragmentBytes, err := fragment.Bytes()
	if err != nil {
		return [32]byte{}, err
	}
	hash := crypto.Keccak256(fragmentBytes)
	hash32 := [32]byte{}
	copy(hash32[:], hash)
	return hash32, nil
}

// Bytes returns a Fragment serialized into a bytes.
func (fragment *Fragment) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, fragment.OrderType); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, fragment.OrderID); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, fragment.OrderParity); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, fragment.OrderSettlement); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, fragment.OrderExpiry.Unix()); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, fragment.Tokens); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, fragment.Price); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, fragment.Volume); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, fragment.MinimumVolume); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, fragment.Nonce); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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
		fragment.MinimumVolume.Equal(&other.MinimumVolume) &&
		fragment.Nonce.Equal(&other.Nonce)
}

// Encrypt a Fragment using an rsa.PublicKey.
func (fragment *Fragment) Encrypt(pubKey rsa.PublicKey) (EncryptedFragment, error) {
	var err error
	encryptedFragment := EncryptedFragment{
		OrderID:         fragment.OrderID,
		OrderType:       fragment.OrderType,
		OrderParity:     fragment.OrderParity,
		OrderSettlement: fragment.OrderSettlement,
		OrderExpiry:     fragment.OrderExpiry,
		ID:              fragment.ID,
		EpochDepth:      fragment.EpochDepth,
	}
	encryptedFragment.Tokens, err = fragment.Tokens.Encrypt(pubKey)
	if err != nil {
		return encryptedFragment, err
	}
	encryptedFragment.Price, err = fragment.Price.Encrypt(pubKey)
	if err != nil {
		return encryptedFragment, err
	}
	encryptedFragment.Volume, err = fragment.Volume.Encrypt(pubKey)
	if err != nil {
		return encryptedFragment, err
	}
	encryptedFragment.MinimumVolume, err = fragment.MinimumVolume.Encrypt(pubKey)
	if err != nil {
		return encryptedFragment, err
	}
	encryptedFragment.Nonce, err = fragment.Nonce.Encrypt(pubKey)
	if err != nil {
		return encryptedFragment, err
	}
	return encryptedFragment, nil
}

// An EncryptedFragment is a Fragment that has been encrypted by an RSA public
// key.
type EncryptedFragment struct {
	OrderID         ID                  `json:"orderId"`
	OrderType       Type                `json:"orderType"`
	OrderParity     Parity              `json:"orderParity"`
	OrderSettlement Settlement          `json:"orderSettlement"`
	OrderExpiry     time.Time           `json:"orderExpiry"`
	ID              FragmentID          `json:"id"`
	EpochDepth      FragmentEpochDepth  `json:"epochDepth"`
	Tokens          []byte              `json:"tokens"`
	Price           EncryptedCoExpShare `json:"price"`
	Volume          EncryptedCoExpShare `json:"volume"`
	MinimumVolume   EncryptedCoExpShare `json:"minimumVolume"`
	Nonce           []byte              `json:"nonce"`
}

// Decrypt an EncryptedFragment using an rsa.PrivateKey.
func (fragment *EncryptedFragment) Decrypt(privKey *rsa.PrivateKey) (Fragment, error) {
	var err error
	decryptedFragment := Fragment{
		OrderID:         fragment.OrderID,
		OrderType:       fragment.OrderType,
		OrderParity:     fragment.OrderParity,
		OrderSettlement: fragment.OrderSettlement,
		OrderExpiry:     fragment.OrderExpiry,
		ID:              fragment.ID,
		EpochDepth:      fragment.EpochDepth,
	}
	if err := decryptedFragment.Tokens.Decrypt(privKey, fragment.Tokens); err != nil {
		return decryptedFragment, err
	}
	decryptedFragment.Price, err = fragment.Price.Decrypt(privKey)
	if err != nil {
		return decryptedFragment, err
	}
	decryptedFragment.Volume, err = fragment.Volume.Decrypt(privKey)
	if err != nil {
		return decryptedFragment, err
	}
	decryptedFragment.MinimumVolume, err = fragment.MinimumVolume.Decrypt(privKey)
	if err != nil {
		return decryptedFragment, err
	}
	if err := decryptedFragment.Nonce.Decrypt(privKey, fragment.Nonce); err != nil {
		return decryptedFragment, err
	}
	return decryptedFragment, nil
}
