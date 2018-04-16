package hyper

// The Rank of the Hyperdrive determines the Replica that is responsible for
// proposing Blocks.
type Rank int

// The Height of the Hyperdrive determines which Block is expected to be
// proposed next.
type Height int

// A Block of Tuples on which consensus needs to be reached.
type Block struct {
	Rank
	Height
	Tuples

	Signature
}
