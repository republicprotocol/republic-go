package identity

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multihash"
)

// Codes for extracting specific protocol values from a MultiAddress.
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

// MultiAddress is an alias.
type MultiAddress struct {
	address Address
	baseMultiAddress multiaddr.Multiaddr
}

// MultiAddresses is an alias.
type MultiAddresses []MultiAddress

// NewMultiAddressFromString parses and validates an input string. It returns a
// MultiAddress, or an error.
func NewMultiAddressFromString(s string) (MultiAddress, error) {
	multiAddress, err := multiaddr.NewMultiaddr(s)
	if err != nil {
		return MultiAddress{}, err
	}

	address, err := multiAddress.ValueForProtocol(RepublicCode)
	if err != nil {
		return MultiAddress{}, err
	}
	addressAsMultiAddress, err := multiaddr.NewMultiaddr(address)
	if err != nil {
		return MultiAddress{}, err
	}
	baseMultiAddress := multiAddress.Decapsulate(addressAsMultiAddress)

	return MultiAddress{Address(address), baseMultiAddress }, err
}

// NewMultiAddressFromBytes parses and validates an input byte slice. It
// returns a MultiAddress, or an error.
func NewMultiAddressFromBytes(b []byte) (MultiAddress, error) {
	multiAddress, err := multiaddr.NewMultiaddrBytes(b)

	address, err := multiAddress.ValueForProtocol(RepublicCode)
	if err != nil {
		return MultiAddress{}, err
	}
	addressAsMultiAddress, err := multiaddr.NewMultiaddr(address)
	if err != nil {
		return MultiAddress{}, err
	}
	baseMultiAddress := multiAddress.Decapsulate(addressAsMultiAddress)

	return MultiAddress{Address(address), baseMultiAddress }, err
}

// Address returns the Republic address of a MultiAddress, or an error.
func (multiAddress MultiAddress) Address() (Address) {
	return multiAddress.address
}

// MarshalJSON implements the json.Marshaler interface.
func (multiAddress MultiAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(multiAddress.baseMultiAddress.String() + "/republic/" + multiAddress.address.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (multiAddress *MultiAddress) UnmarshalJSON(data []byte) (error) {
	newMultiAddress, err := NewMultiAddressFromBytes(data)
	multiAddress.baseMultiAddress = newMultiAddress.baseMultiAddress
	multiAddress.address = newMultiAddress.address
	return err
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
