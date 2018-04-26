package hyperdrive

import (
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/republic-go/delta"
)

// Txs must not store any Nonce more than once within any Tx.
type Txs []Tx

// A Tx stores Nonces alongside a Keccak256 Hash of the Nonces. A valid Tx must
// not store any Nonce more than once.
type Tx struct {
	Hash []byte
	Nonces
}

func NewTx(nonces ...Nonce) Tx {
	return Tx{
		Hash:   crypto.Keccak256(nonces...),
		Nonces: nonces,
	}
}

func NewTxFromByteSlices(nonces ...[]byte) Tx {
	return Tx{
		Hash:   crypto.Keccak256(nonces...),
		Nonces: nonces,
	}
}

type NonceWithTimestamp struct {
	Nonce     Nonce
	Timestamp time.Time
	DeltaID   delta.ID
}

func NewNonceWithTimestamp(nonce Nonce, t time.Time, deltaID delta.ID) NonceWithTimestamp {
	return NonceWithTimestamp{
		Nonce:     nonce,
		Timestamp: t,
		DeltaID:   deltaID,
	}
}

// Nonces must not store any Nonce more than once.
type Nonces []Nonce

// A Nonce is a unique 256-bit value that makes up a Tx. It must be unique
// within the entire Hyperdrive blockchain.
type Nonce = []byte
