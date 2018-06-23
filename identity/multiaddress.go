package identity

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
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
		Transcoder: multiaddr.NewTranscoderFromFunctions(republicStB, republicBtS, nil),
	}
	multiaddr.AddProtocol(republic)
}

// MultiAddress is an alias.
type MultiAddress struct {
	Signature []byte

	address          Address
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
	addressAsMultiAddress, err := multiaddr.NewMultiaddr("/republic/" + address)
	if err != nil {
		return MultiAddress{}, err
	}
	baseMultiAddress := multiAddress.Decapsulate(addressAsMultiAddress)

	return MultiAddress{[]byte{}, Address(address), baseMultiAddress}, err
}

// ValueForProtocol returns the value of the specific protocol in the MultiAddress
func (multiAddress MultiAddress) ValueForProtocol(code int) (string, error) {
	if code == RepublicCode {
		return multiAddress.address.String(), nil
	}
	return multiAddress.baseMultiAddress.ValueForProtocol(code)
}

// Address returns the Republic address of a MultiAddress.
func (multiAddress MultiAddress) Address() Address {
	return multiAddress.address
}

// ID returns the Republic ID of a MultiAddress.
func (multiAddress MultiAddress) ID() ID {
	return multiAddress.address.ID()
}

// String returns the MultiAddress as a plain string.
func (multiAddress MultiAddress) String() string {
	return fmt.Sprintf("%s/republic/%s", multiAddress.baseMultiAddress.String(), multiAddress.address.String())
}

// Hash returns the Keccak256 hash of a multiaddrfess. This hash is used to create
// signatures for a multiaddress.
func (multiAddress MultiAddress) Hash() []byte {
	return crypto.Keccak256([]byte(multiAddress.String()))
}

// MarshalJSON implements the json.Marshaler interface.
func (multiAddress MultiAddress) MarshalJSON() ([]byte, error) {
	if multiAddress.address.String() == "" {
		return json.Marshal("")
	}
	return json.Marshal(multiAddress.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (multiAddress *MultiAddress) UnmarshalJSON(data []byte) error {
	multiAddressAsString := ""
	if err := json.Unmarshal(data, &multiAddressAsString); err != nil {
		return err
	}
	if multiAddressAsString == "" {
		return nil
	}
	newMultiAddress, err := NewMultiAddressFromString(multiAddressAsString)
	if err != nil {
		return err
	}
	multiAddress.baseMultiAddress = newMultiAddress.baseMultiAddress
	multiAddress.address = newMultiAddress.address
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
