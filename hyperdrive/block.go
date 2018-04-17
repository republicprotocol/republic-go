package hyper

// The Epoch in which a Block was proposed.
type Epoch [32]byte

// The Rank in which a Block was proposed. The Epoch defines the set of
// validating Replicas and the Rank defines which Replica is responsible for
// proposing Blocks.
type Rank int

// The Height of the Hyperdrive determines which Block is expected to be
// proposed next.
type Height int

// Blocks must not store any Nonce more than once within any Block.
type Blocks []Block

// A Block of Txs on which consensus will be reached. A Block must not store
// any Nonce more than once within any Tx.
type Block struct {
	Epoch
	Rank
	Height
	Txs

	// Signature of the Replica that proposed this Block
	Signature
}

// Txs must not store any Nonce more than once within any Tx.
type Txs []Tx

// A Tx stores Nonces alongside a Keccak256 Hash of the Nonces. A valid Tx must
// not store any Nonce more than once.
type Tx struct {
	Hash [32]byte
	Nonces
}

// Nonces must not store any Nonce more than once.
type Nonces []Nonce

// A Nonce is a unique 256-bit value that makes up a Tx. It must be unique
// within the entire Hyperdrive blockchain.
type Nonce [32]byte
