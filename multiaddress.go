package identity

import (
	"encoding/binary"
	"errors"
	"fmt"

	multiaddr "github.com/multiformats/go-multiaddr"
	multihash "github.com/multiformats/go-multihash"
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

// NewMultiaddr parses and validates an input string, returning a *Multiaddr
func NewMultiaddr(s string) (multiaddr.Multiaddr, error) {
	return multiaddr.NewMultiaddr(s)
}

// NewMultiaddrBytes initializes a Multiaddr from a byte representation. It
// validates it as an input string.
func NewMultiaddrBytes(b []byte) (multiaddr.Multiaddr, error) {
	return multiaddr.NewMultiaddrBytes(b)
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
