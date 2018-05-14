package order

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/republicprotocol/republic-go/shamir"
)

// ErrUnexpectedValueRange is returned when a Value is encoded with an integer
// that is outside of the restricted range.
var ErrUnexpectedValueRange = errors.New("unexpected value range")

// A Value represented by the equation `co * 10 ^ exp`. The coefficient is
// restricted to the range 1-1999, for values 0.005 to 9.995. The exponent is
// restricted to the range 0-52, for values -26 to 25.
type Value struct {
	Co  uint32
	Exp uint32
}

// MarshalJSON implements the json.Marshaler interface and marshals the Value
// into a tuple of numbers.
func (val Value) MarshalJSON() ([]byte, error) {
	return json.Marshal([]uint32{val.Co, val.Exp})
}

// UnmarshalJSON implements the json.Unmarshaler interface and unmarshals the
// Value from a tuple of numbers.
func (val *Value) UnmarshalJSON(data []byte) error {
	tuple := [2]uint32{0, 0}
	if err := json.Unmarshal(data, &tuple); err != nil {
		return err
	}
	val.Co = tuple[0]
	val.Exp = tuple[1]
	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (val Value) MarshalBinary() (data []byte, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, val.Co)
	binary.Write(buf, binary.BigEndian, val.Exp)
	return buf.Bytes(), nil
}

// A FragmentValue is a Value that has been shared using Shamir's secret
// sharing scheme. A FragmentValue is unencrypted.
type FragmentValue struct {
	Co  shamir.Share
	Exp shamir.Share
}

// Equal returns true when both FragmeValues are equal. Otherwise, it returns
// false.
func (val *FragmentValue) Equal(other *FragmentValue) bool {
	return val.Co.Equal(&other.Co) &&
		val.Exp.Equal(&other.Exp)
}

// MarshalJSON implements the json.Marshaler interface and marshals the
// FragmentValue into a tuple of bytes.
func (val FragmentValue) MarshalJSON() ([]byte, error) {
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
// FragmentValue from a tuple of strings.
func (val *FragmentValue) UnmarshalJSON(data []byte) error {
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
	if err = val.Exp.UnmarshalBinary(expBytes); err != nil {
		return err
	}

	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface and
// marshals the FragmentValue using encoding.BigEndian.
func (val FragmentValue) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, val.Co.Index)
	binary.Write(buf, binary.BigEndian, val.Co.Value)
	binary.Write(buf, binary.BigEndian, val.Exp.Index)
	binary.Write(buf, binary.BigEndian, val.Exp.Value)
	return buf.Bytes(), nil
}

// An EncryptedFragmentValue is a FragmentValue that has been encrypted using
// an RSA public key.
type EncryptedFragmentValue struct {
	Co  []byte
	Exp []byte
}

// Equal returns true when both EncryptedFragmentValue are equal. Otherwise, it
// returns false.
func (val *EncryptedFragmentValue) Equal(other *EncryptedFragmentValue) bool {
	return bytes.Equal(val.Co[:], other.Co[:]) &&
		bytes.Equal(val.Exp[:], other.Exp[:])
}

// MarshalJSON implements the json.Marshaler interface and marshals the
// EncryptedFragmentValue into a tuple of bytes.
func (val EncryptedFragmentValue) MarshalJSON() ([]byte, error) {
	tuple := [2]string{"", ""}
	tuple[0] = base64.StdEncoding.EncodeToString(val.Co)
	tuple[1] = base64.StdEncoding.EncodeToString(val.Exp)
	return json.Marshal(tuple)
}

// UnmarshalJSON implements the json.Unmarshaler interface and unmarshals the
// EncryptedFragmentValue from a tuple of strings.
func (val *EncryptedFragmentValue) UnmarshalJSON(data []byte) error {
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
// marshals the EncryptedFragmentValue using encoding.BigEndian.
func (val EncryptedFragmentValue) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, val.Co)
	binary.Write(buf, binary.BigEndian, val.Exp)
	return buf.Bytes(), nil
}
