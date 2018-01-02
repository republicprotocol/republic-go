package identity

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multihash"
)

// Codes for extracting specific protocol values from a multiaddress.
const (
	IP4Code      = 0x0004
	TCPCode      = 0x0006
	IP6Code      = 0x0029
	RepublicCode = 0x0065
)

// Add the republic protocol
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
	multiaddr.Multiaddr
}

// MultiAddresses is an alias.
type MultiAddresses []MultiAddress

// Address returns the Republic Address of a MultiAddress, or an error.
func (multiAddr MultiAddress) Address() (Address, error) {
	addr, err := multiAddr.ValueForProtocol(RepublicCode)
	return Address(addr), err
}

// MarshalJSON implements the json.Marshaler interface.
func (multiAddr MultiAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(multiAddr.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (multiAddr MultiAddress) UnmarshalJSON(data []byte) error {
	addr, err := NewMultiAddressBytes(data)
	if err != nil {
		return err
	}
	multiAddr.Multiaddr = addr.Multiaddr
	return nil
}

// NewMultiAddress parses and validates an input string, returning a
// MultiAddress. It returns a MultiAddress or an error.
func NewMultiAddress(s string) (MultiAddress, error) {
	multiAddr, err := multiaddr.NewMultiaddr(s)
	return MultiAddress{multiAddr}, err
}

// NewMultiAddressBytes initializes a MultiAddress from a byte representation.
// It returns a MultiAddress or an error.
func NewMultiAddressBytes(b []byte) (MultiAddress, error) {
	multiAddr, err := multiaddr.NewMultiaddrBytes(b)
	return MultiAddress{multiAddr}, err
}

// ProtocolWithName returns the Protocol description with given string name.
func ProtocolWithName(s string) multiaddr.Protocol {
	return multiaddr.ProtocolWithName(s)
}

// ProtocolWithCode returns the Protocol description with given string name.
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
	size := codeToVarint(len(m))
	b := append(size, m...)
	return b, nil
}

// republicBtS converts a republic address from bytes to a string.
func republicBtS(b []byte) (string, error) {
	// The address is a variant-prefixed multihash string representation.
	size, n, err := readVarintCode(b)
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
	return m.B58String(), nil
}

// codeToVarint converts an integer to a varint-encoded []byte.
func codeToVarint(num int) []byte {
	buf := make([]byte, (num/7)+1) // varint package is uint64
	n := binary.PutUvarint(buf, uint64(num))
	return buf[:n]
}

// readVarintCode reads a varint code from the beginning of buf. It returns the
// code, and the number of bytes read.
func readVarintCode(buf []byte) (int, int, error) {
	num, n := binary.Uvarint(buf)
	if n < 0 {
		return 0, 0, fmt.Errorf("varints larger than uint64 not yet supported")
	}
	return int(num), n, nil
}
