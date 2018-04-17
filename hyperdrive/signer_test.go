package hyper_test

import (
	. "github.com/republicprotocol/republic-go/hyperdrive"
	"github.com/republicprotocol/republic-go/identity"
)

type TestValidator struct {
	id           identity.ID
	kp           identity.KeyPair
	sharedBlocks SharedBlocks
	t            uint8
}

func NewTestValidator(sb SharedBlocks, t uint8) (Validator, error) {
	id, kp, err := identity.NewID()
	if err != nil {
		return TestValidator{}, err
	}
	return TestValidator{
		id:           id,
		kp:           kp,
		sharedBlocks: sb,
		t:            t,
	}, nil
}

func (v *TestValidator) ValidateBlock(b Block) bool {
	for _, tuple := range b.Tuples {
		if !v.ValidateTuple(tuple) {
			return false
		}
	}
	return true
}

func (v *TestValidator) ValidateRank(r Rank) bool {
	return (r == v.sharedBlocks.Rank)
}

func (v *TestValidator) ValidateHeight(h uint64) bool {
	return (h == v.sharedBlocks.ReadHeight())
}

func (v *TestValidator) GetSharedBlocks() *SharedBlocks {
	return &v.sharedBlocks
}

func (v *TestValidator) ValidateProposal(p Proposal) bool {
	// if !v.Block(p.Block) {
	// 	return false
	// }
	// if !v.Signatures(p) {
	// 	return false
	// }
	// if !v.Rank(p.Rank) {
	// 	return false
	// }
	// if !v.Height(p.Height) {
	// 	return false
	// }
	return true
}

func (v *TestValidator) ValidatePrepare(p Prepare) bool {
	// if v.prepareSigs[p.Signature] {
	// 	return false
	// }
	// v.prepareSigs[p.Signature] = true
	// if !v.Block(p.Block) {
	// 	return false
	// }
	// if !v.Rank(p.Rank) {
	// 	return false
	// }
	// if !v.Height(p.Height) {
	// 	return false
	// }
	return true
}

func (v *TestValidator) ValidateCommit(c Commit) bool {
	// if v.commitSigs[c.Signature] {
	// 	return false
	// }
	// v.commitSigs[c.Signature] = true
	return (c.ThresholdSignature == ThresholdSignature("Threshold_BLS"))
}

func (v *TestValidator) SharedBlocks() *SharedBlocks {
	return &v.sharedBlocks
}

func (v *TestValidator) Sign() Signature {
	return Signature(v.id.String())
}

func (v *TestValidator) Threshold() uint8 {
	return v.t
}

func (v *TestValidator) ValidateTuple(t Tuple) bool {
	return true
}
