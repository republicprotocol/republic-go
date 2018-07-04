package order

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/republicprotocol/republic-go/shamir"
)

// ErrUnexpectedCoExpRange is returned when a CoExp is encoded with an integer
// that is outside of the restricted range.
var ErrUnexpectedCoExpRange = errors.New("unexpected value range")

// A CoExp represented by the equation `co * 10 ^ exp`. The coefficient is
// restricted to the range 1-1999, for values 0.005 to 9.995. The exponent is
// restricted to the range 0-52, for values -26 to 25.
type CoExp struct {
	Co  uint64
	Exp uint64
}

// NewCoExp returns a CoExp with the given coefficient and exponent.
func NewCoExp(co uint64, exp uint64) CoExp {
	return CoExp{
		Co:  co,
		Exp: exp,
	}
}

// MarshalJSON implements the json.Marshaler interface and marshals the CoExp
// into a tuple of numbers.
func (val CoExp) MarshalJSON() ([]byte, error) {
	return json.Marshal([]uint64{val.Co, val.Exp})
}

// UnmarshalJSON implements the json.Unmarshaler interface and unmarshals the
// CoExp from a tuple of numbers.
func (val *CoExp) UnmarshalJSON(data []byte) error {
	tuple := [2]uint64{0, 0}
	if err := json.Unmarshal(data, &tuple); err != nil {
		return err
	}
	val.Co = tuple[0]
	val.Exp = tuple[1]
	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (val CoExp) MarshalBinary() (data []byte, err error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, val.Co); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, val.Exp); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// A CoExpShare is a CoExp that has been shared using Shamir's secret sharing
// scheme. A CoExpShare is unencrypted.
type CoExpShare struct {
	Co  shamir.Share
	Exp shamir.Share
}

// Equal returns true when both CoExpShares are equal. Otherwise, it returns
// false.
func (val *CoExpShare) Equal(other *CoExpShare) bool {
	return val.Co.Equal(&other.Co) &&
		val.Exp.Equal(&other.Exp)
}

// MarshalJSON implements the json.Marshaler interface and marshals the
// CoExpShare into a tuple of bytes.
func (val CoExpShare) MarshalJSON() ([]byte, error) {
	tuple := [2]string{"", ""}

	coBytes, err := val.Co.MarshalBinary()
	if err != nil {
		return nil, err
	}
	tuple[0] = base64.StdEncoding.EncodeToString(coBytes)

	expBytes, err := val.Exp.MarshalBinary()
	if err != nil {
		return nil, err
	}
	tuple[1] = base64.StdEncoding.EncodeToString(expBytes)

	return json.Marshal(tuple)
}

// UnmarshalJSON implements the json.Unmarshaler interface and unmarshals the
// CoExpShare from a tuple of strings.
func (val *CoExpShare) UnmarshalJSON(data []byte) error {
	tuple := [2]string{"", ""}
	if err := json.Unmarshal(data, &tuple); err != nil {
		return err
	}

	coBytes, err := base64.StdEncoding.DecodeString(tuple[0])
	if err != nil {
		return err
	}
	if err = val.Co.UnmarshalBinary(coBytes); err != nil {
		return err
	}

	expBytes, err := base64.StdEncoding.DecodeString(tuple[1])
	if err != nil {
		return err
	}

	return val.Exp.UnmarshalBinary(expBytes)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface and
// marshals the CoExpShare using encoding.BigEndian.
func (val CoExpShare) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, val.Co.Index); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, val.Co.Value); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, val.Exp.Index); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, val.Exp.Value); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Encrypt a CoExpShare using an rsa.PublicKey.
func (val *CoExpShare) Encrypt(pubKey rsa.PublicKey) (EncryptedCoExpShare, error) {
	var err error
	encryptedVal := EncryptedCoExpShare{}
	encryptedVal.Co, err = val.Co.Encrypt(pubKey)
	if err != nil {
		return encryptedVal, err
	}
	encryptedVal.Exp, err = val.Exp.Encrypt(pubKey)
	if err != nil {
		return encryptedVal, err
	}
	return encryptedVal, nil
}

// An EncryptedCoExpShare is a CoExpShare that has been encrypted using an RSA
// public key.
type EncryptedCoExpShare struct {
	Co  []byte
	Exp []byte
}

// Equal returns true when both EncryptedCoExpShare are equal. Otherwise, it
// returns false.
func (val *EncryptedCoExpShare) Equal(other *EncryptedCoExpShare) bool {
	return bytes.Equal(val.Co[:], other.Co[:]) &&
		bytes.Equal(val.Exp[:], other.Exp[:])
}

// MarshalJSON implements the json.Marshaler interface and marshals the
// EncryptedCoExpShare into a tuple of bytes.
func (val EncryptedCoExpShare) MarshalJSON() ([]byte, error) {
	tuple := [2]string{"", ""}
	tuple[0] = base64.StdEncoding.EncodeToString(val.Co)
	tuple[1] = base64.StdEncoding.EncodeToString(val.Exp)
	return json.Marshal(tuple)
}

// UnmarshalJSON implements the json.Unmarshaler interface and unmarshals the
// EncryptedCoExpShare from a tuple of strings.
func (val *EncryptedCoExpShare) UnmarshalJSON(data []byte) error {
	var err error

	tuple := [2]string{"", ""}
	if err = json.Unmarshal(data, &tuple); err != nil {
		return err
	}

	val.Co, err = base64.StdEncoding.DecodeString(tuple[0])
	if err != nil {
		return err
	}

	val.Exp, err = base64.StdEncoding.DecodeString(tuple[1])
	if err != nil {
		return err
	}

	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface and
// marshals the EncryptedCoExpShare using encoding.BigEndian.
func (val EncryptedCoExpShare) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, val.Co); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, val.Exp); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decrypt an EncryptedCoExpShare using an rsa.PrivateKey.
func (val *EncryptedCoExpShare) Decrypt(privKey *rsa.PrivateKey) (CoExpShare, error) {
	decryptedVal := CoExpShare{Co: shamir.Share{}, Exp: shamir.Share{}}
	if err := decryptedVal.Co.Decrypt(privKey, val.Co); err != nil {
		return decryptedVal, err
	}
	if err := decryptedVal.Exp.Decrypt(privKey, val.Exp); err != nil {
		return decryptedVal, err
	}
	return decryptedVal, nil
}
