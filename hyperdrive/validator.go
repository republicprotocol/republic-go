package hyper

type Validator interface {
	Sign() Signature
	SharedBlocks() *SharedBlocks
	Threshold() uint8
	ValidateProposal(Proposal) bool
	ValidatePrepare(Prepare) bool
	ValidateCommit(Commit) bool
	ValidateBlock(Block) bool
	ValidateRank(Rank) bool
	ValidateHeight(uint64) bool
	ValidateTuple(Tuple) bool
}
