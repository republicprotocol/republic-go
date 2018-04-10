package hyper

type Validator interface {
	Sign([32]byte) []byte
	SharedBlocks() *SharedBlocks
	Threshold() uint8
	Proposal(Proposal) bool
	Prepare(Prepare) bool
	Commit(Commit) bool
	Block(Block) bool
	Rank(Rank) bool
	Height(Height) bool
}

type EthereumValidator struct {
}

// func NewEthereumValidator() Validator {
// }
