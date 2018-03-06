package identity

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multihash"
)

// Codes for extracting specific protocol values from a Multiaddress.
const (
	IP4Code      = 0x0004
	IP6Code      = 0x0029
	TCPCode      = 0x0006
	RepublicCode = 0x0065
)

// Add the Republic Protocol when the package is initialized.
func init() {
	republic := multiaddr.Protocol{
		Code:       RepublicCode,
		Size:       multiaddr.LengthPrefixedVarSize,
		Name:       "republic",
		Path:       false,
		Transcoder: multiaddr.NewTranscoderFromFunctions(republicStB, republicBtS),
	}
	multiaddr.AddProtocol(republic)
}

// Multiaddress is an alias.
type Multiaddress struct {
	address          Address
	baseMultiaddress multiaddr.Multiaddr
}

// Multiaddresses is an alias.
type Multiaddresses []Multiaddress

// NewMultiaddressFromString parses and validates an input string. It returns a
// Multiaddress, or an error.
func NewMultiaddressFromString(s string) (Multiaddress, error) {
	multiaddress, err := multiaddr.NewMultiaddr(s)
	if err != nil {
		return Multiaddress{}, err
	}
	address, err := multiaddress.ValueForProtocol(RepublicCode)
	if err != nil {
		return Multiaddress{}, err
	}
	addressAsMultiaddress, err := multiaddr.NewMultiaddr("/republic/" + address)
	if err != nil {
		return Multiaddress{}, err
	}
	baseMultiaddress := multiaddress.Decapsulate(addressAsMultiaddress)

	return Multiaddress{Address(address), baseMultiaddress}, err
}

// ValueForProtocol returns the value of the specific protocol in the Multiaddress
func (multiaddress Multiaddress) ValueForProtocol(code int) (string, error) {
	if code == RepublicCode {
		return multiaddress.address.String(), nil
	}
	return multiaddress.baseMultiaddress.ValueForProtocol(code)
}

// Address returns the Republic address of a Multiaddress.
func (multiaddress Multiaddress) Address() Address {
	return multiaddress.address
}

// ID returns the Republic ID of a Multiaddress.
func (multiaddress Multiaddress) ID() ID {
	return multiaddress.address.ID()
}

// String returns the Multiaddress as a plain string.
func (multiaddress Multiaddress) String() string {
	return fmt.Sprintf("%s/republic/%s", multiaddress.baseMultiaddress, multiaddress.address)
}

// MarshalJSON implements the json.Marshaler interface.
func (multiaddress Multiaddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(multiaddress.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (multiaddress *Multiaddress) UnmarshalJSON(data []byte) error {
	multiaddressAsString := ""
	if err := json.Unmarshal(data, &multiaddressAsString); err != nil {
		return err
	}
	newMultiaddress, err := NewMultiaddressFromString(multiaddressAsString)
	if err != nil {
		return err
	}
	multiaddress.baseMultiaddress = newMultiaddress.baseMultiaddress
	multiaddress.address = newMultiaddress.address
	return nil
}

// ProtocolWithName returns the Protocol description with the given name.
func ProtocolWithName(s string) multiaddr.Protocol {
	return multiaddr.ProtocolWithName(s)
}

// ProtocolWithCode returns the Protocol description with the given code.
func ProtocolWithCode(c int) multiaddr.Protocol {
	return multiaddr.ProtocolWithCode(c)
}

// republicStB converts a republic address from a string to bytes.
func republicStB(s string) ([]byte, error) {
	// The address is a varint prefixed multihash string representation.
	m, err := multihash.FromB58String(s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse republic addr: %s %s", s, err)
	}
	size := multiaddr.CodeToVarint(len(m))
	b := append(size, m...)
	return b, nil
}

// republicBtS converts a Republic address, encoded as bytes, to a string.
func republicBtS(b []byte) (string, error) {
	size, n, err := multiaddr.ReadVarintCode(b)
	if err != nil {
		return "", err
	}
	b = b[n:]
	if len(b) != size {
		return "", errors.New("inconsistent lengths")
	}
	m, err := multihash.Cast(b)
	if err != nil {
		return "", err
	}
	// This uses the default Bitcoin alphabet for Base58 encoding.
	return m.B58String(), nil
}
