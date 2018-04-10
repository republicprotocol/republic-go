package hyper_test

import (
	"context"
	"sync"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/hyperdrive"
	"github.com/republicprotocol/republic-go/identity"
)

func TestHyperdrive(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hyperdrive Suite")
}

type TestValidator struct {
	SharedBlocks
	// prepareSigs map[Signature]bool
	// commitSigs  map[Signature]bool
}

func NewTestValidator(sb SharedBlocks) (Validator, error) {
	return &TestValidator{
		sb,
		// map[Signature]bool{},
		// map[Signature]bool{},
	}, nil
}

func (v *TestValidator) Block(b Block) bool {
	for _, tuple := range b.Tuples {
		if !v.SharedBlocks.ValidateTuple(tuple) {
			return false
		}
	}
	return true
}

func (v *TestValidator) Signatures(p Proposal) bool {
	return (p.Signature == p.Block.Signature)
}

func (v *TestValidator) Rank(r Rank) bool {
	return (r == v.SharedBlocks.Rank)
}

func (v *TestValidator) Height(h Height) bool {
	return (h == v.SharedBlocks.Height)
}

func (v *TestValidator) GetSharedBlocks() *SharedBlocks {
	return v.SharedBlocks.GetSharedBlocks()
}

func (v *TestValidator) Proposal(p Proposal) bool {
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

func (v *TestValidator) Prepare(p Prepare) bool {
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

func (v *TestValidator) Commit(c Commit) bool {
	// if v.commitSigs[c.Signature] {
	// 	return false
	// }
	// v.commitSigs[c.Signature] = true
	return (c.ThresholdSignature == ThresholdSignature("Threshold_BLS"))
}

type TestSigner struct {
	id identity.ID
	kp identity.KeyPair
}

func NewTestSigner() (Signer, error) {
	id, kp, err := identity.NewID()
	if err != nil {
		return &TestSigner{}, err
	}
	return &TestSigner{
		id: id,
		kp: kp,
	}, nil
}

func (t *TestSigner) Sign() Signature {
	return Signature(t.id.String())
}

func (t *TestSigner) ID() identity.ID {
	return t.id
}

type TestNetwork struct {
	commanderCount uint64
	Ingress        []ChannelSet
	Egress         []ChannelSet
}

func NewTestNetwork(commanderCount uint64) *TestNetwork {
	return &TestNetwork{
		commanderCount: commanderCount,
		Ingress:        make([]ChannelSet, commanderCount),
		Egress:         make([]ChannelSet, commanderCount),
	}
}

func (t *TestNetwork) init() {
	for i := uint64(0); i < t.commanderCount; i++ {
		t.Ingress[i] = EmptyChannelSet()
		t.Egress[i] = EmptyChannelSet()
	}
}

func (t *TestNetwork) propose(p Proposal) {
	for i := uint64(0); i < t.commanderCount; i++ {
		go func(j uint64) {
			t.Ingress[j].Proposal <- p
		}(i)
	}
}

func (t *TestNetwork) proposeMultiple(proposals []Proposal) {
	var wg sync.WaitGroup
	for i := uint64(0); i < t.commanderCount; i++ {
		wg.Add(1)
		go func(i uint64) {
			defer wg.Done()
			for _, proposal := range proposals {
				t.Ingress[i].Proposal <- proposal
			}
		}(i)
	}
	wg.Wait()
}

func (t *TestNetwork) run(ctx context.Context) {
	var wg sync.WaitGroup
	for i := uint64(0); i < t.commanderCount; i++ {
		wg.Add(1)
		go func(i uint64) {
			defer wg.Done()
			t.Egress[i].Split(t.Ingress)
		}(i)
	}
	wg.Wait()
}

type HyperDrive struct {
	commanderCount uint64
	threshold      uint8
	network        TestNetwork
	validator      Validator
	replicas       []Replica
}

func NewHyperDrive(commanderCount uint64, validator Validator, threshold uint8) *HyperDrive {
	replicas := make([]Replica, commanderCount)
	network := NewTestNetwork(commanderCount)
	network.init()
	return &HyperDrive{
		commanderCount: commanderCount,
		network:        *network,
		replicas:       replicas,
		validator:      validator,
		threshold:      threshold,
	}
}

func (h *HyperDrive) init() {
	for i := uint64(0); i < h.commanderCount; i++ {
		signer, _ := NewTestSigner()
		h.replicas[i] = NewReplica(h.validator, signer, h.network.Ingress[i], h.network.Egress[i], h.threshold)
	}
}

func (h *HyperDrive) run(ctx context.Context) {
	for i := uint64(0); i < h.commanderCount; i++ {
		h.replicas[i].Run(ctx)
	}
}
