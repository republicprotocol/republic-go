package identity

import (
	ma "github.com/multiformats/go-multiaddr"
	mh "github.com/multiformats/go-multihash"

	"fmt"
	"errors"
	"encoding/binary"
)

const (
	P_IP4   	= 	0x0004
	P_TCP   	= 	0x0006
	P_UDP   	= 	0x0111
	P_DCCP  	= 	0x0021
	P_IP6   	= 	0x0029
	P_QUIC  	= 	0x01CC
	P_SCTP  	= 	0x0084
	P_UDT   	= 	0x012D
	P_UTP   	= 	0x012E
	P_UNIX  	= 	0x0190
	P_IPFS  	= 	0x01A5
	P_HTTP  	= 	0x01E0
	P_HTTPS 	= 	0x01BB
	P_ONION 	= 	0x01BC
	P_REPUBLIC  = 	0x0065

	LengthPrefixedVarSize = -1
)

// Add the republic protocol
func init() {
	republic := ma.Protocol{
		Code:       P_REPUBLIC,
		Size:       LengthPrefixedVarSize,
		Name:       "republic",
		Path:       false,
		Transcoder: ma.NewTranscoderFromFunctions(republicStB,republicBtS),
	}
	ma.AddProtocol(republic)
}

// convert republic address from string to bytes
func republicStB(s string) ([]byte, error) {
	// the address is a varint prefixed multihash string representation
	m, err := mh.FromB58String(s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse republic addr: %s %s", s, err)
	}
	size := CodeToVarint(len(m))
	b := append(size, m...)
	return b, nil
}

// convert republic address from bytes to string
func republicBtS(b []byte) (string, error) {
	// the address is a variant-prefixed multihash string representation
	size, n, err := ReadVarintCode(b)
	if err != nil {
		return "", err
	}

	b = b[n:]
	if len(b) != size {
		return "", errors.New("inconsistent lengths")
	}
	m, err := mh.Cast(b)
	if err != nil {
		return "", err
	}
	return m.B58String(), nil
}

// NewMultiaddr parses and validates an input string, returning a *Multiaddr
func NewMultiaddr(s string) (ma.Multiaddr, error) {
	return ma.NewMultiaddr(s)
}

// NewMultiaddrBytes initializes a Multiaddr from a byte representation.
// It validates it as an input string.
func NewMultiaddrBytes(b []byte) (ma.Multiaddr, error) {
	return ma.NewMultiaddrBytes(b)
}

// CodeToVarint converts an integer to a varint-encoded []byte
func CodeToVarint(num int) []byte {
	buf := make([]byte, (num/7)+1) // varint package is uint64
	n := binary.PutUvarint(buf, uint64(num))
	return buf[:n]
}

// ReadVarintCode reads a varint code from the beginning of buf.
// returns the code, and the number of bytes read.
func ReadVarintCode(buf []byte) (int, int, error) {
	num, n := binary.Uvarint(buf)
	if n < 0 {
		return 0, 0, fmt.Errorf("varints larger than uint64 not yet supported")
	}
	return int(num), n, nil
}

