package hyper

type Validator interface {
	Proposal(Proposal) bool
	Prepare(Prepare) bool
	Commit(Commit) bool
	Block(Block) bool
	Rank(Rank) bool
	Height(Height) bool
	GetSharedBlocks() *SharedBlocks
}

type EthereumValidator struct {
}

// func NewEthereumValidator() Validator {
// }
