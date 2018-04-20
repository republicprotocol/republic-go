package hyperdrive

import (
	"bytes"
	"context"
	"encoding/binary"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/republicprotocol/republic-go/identity"
	"golang.org/x/crypto/sha3"
)

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

// BlockHeader distinguishes Blocks from other message types that have the
// same content.
const BlockHeader = 4

// A Block of Txs on which consensus will be reached. A Block must not store
// any Nonce more than once within any Tx.
type Block struct {
	Epoch
	Rank
	Height
	Txs

	// Signature of the Replica that proposed this Block
	identity.Signature
}

// Hash returns the SHA3-256 hash of the block.
func (block *Block) Hash() identity.Hash {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, BlockHeader)
	binary.Write(&buf, binary.BigEndian, block.Epoch)
	binary.Write(&buf, binary.BigEndian, block.Rank)
	binary.Write(&buf, binary.BigEndian, block.Height)
	for i := range block.Txs {
		binary.Write(&buf, binary.BigEndian, block.Txs[i])
	}
	binary.Write(&buf, binary.BigEndian, block.Signature)

	return sha3.Sum256(buf.Bytes())
}

// Verify the Block message. Returns an error if the message is invalid,
// otherwise nil.
func (block *Block) VerifyBlock(verifier identity.Verifier) error {
	return verifier.VerifySignature(block.Signature)
}

type HyperChain struct {
	epochChange *sync.Cond
	signer      *identity.Signer
	verifier    *identity.Verifier

	proposalMu *sync.RWMutex
	proposals  map[[32]byte]struct{}

	prepareMu *sync.RWMutex
	prepares  map[[32]byte]*int64

	commitMu *sync.RWMutex
	commit   map[[32]byte]*int64

	blockMu *sync.RWMutex
	blocks  []Block

	faultMu *sync.RWMutex
	faults  map[[32]byte]*int64

	nonceMu *sync.RWMutex
	nonces  map[[32]byte]struct{}
}

func NewHyperChain() HyperChain {
	return HyperChain{
		epochChange: sync.NewCond(new(sync.Mutex)),

		proposalMu: new(sync.RWMutex),
		proposals:  map[[32]byte]struct{}{},

		prepareMu: new(sync.RWMutex),
		prepares:  map[[32]byte]*int64{},

		commitMu: new(sync.RWMutex),
		commit:   map[[32]byte]*int64{},

		blockMu: new(sync.RWMutex),
		blocks:  []Block{},

		faultMu: new(sync.RWMutex),
		faults:  map[[32]byte]*int64{},

		nonceMu: new(sync.RWMutex),
		nonces:  map[[32]byte]struct{}{},
	}
}

func (chain *HyperChain) AddProposal(proposal Proposal) bool {
	if chain.VerifyBlock(proposal.Block) {
		chain.proposalMu.Lock()
		defer chain.proposalMu.Unlock()

		if _, ok := chain.proposals[proposal.Hash()]; !ok {
			chain.proposals[proposal.Hash()] = struct{}{}
			return true
		}
	}

	return false
}

func (chain *HyperChain) AddPrepare(prepare Prepare, threshold int) bool {
	if chain.VerifyBlock(prepare.Proposal.Block) {
		if _, ok := chain.prepares[prepare.Hash()]; !ok {
			chain.prepareMu.Lock()
			defer chain.prepareMu.Unlock()

			count := int64(1)
			chain.prepares[prepare.Hash()] = &count
		} else {
			chain.prepareMu.RLock()
			defer chain.prepareMu.RUnlock()

			atomic.AddInt64(chain.prepares[prepare.Hash()], 1)
			if atomic.LoadInt64(chain.prepares[prepare.Hash()]) > int64(threshold) {
				return true
			}
		}
	}

	return false
}

func (chain *HyperChain) AddCommit(commit Commit, threshold int) bool {
	if chain.VerifyBlock(commit.Proposal.Block) {
		if _, ok := chain.commit[commit.Hash()]; !ok {
			chain.commitMu.Lock()
			defer chain.commitMu.Unlock()

			count := int64(1)
			chain.commit[commit.Hash()] = &count
		} else {
			chain.commitMu.RLock()
			defer chain.commitMu.RUnlock()

			atomic.AddInt64(chain.commit[commit.Hash()], 1)
			if atomic.LoadInt64(chain.commit[commit.Hash()]) > int64(threshold) {
				return true
			}
		}
	}

	return false
}

func (chain *HyperChain) AddBlock(block Block) error {
	if !chain.VerifyBlock(block) {
		return errors.New("conflicts with previous blocks")
	}

	// Check the height is right
	if len(chain.blocks) != int(block.Height) {
		return errors.New("invalid block to add to the chain")
	}

	// Finalize the block in the hyperChain
	chain.blockMu.Lock()
	chain.blocks = append(chain.blocks, block)

	// When we reach certain number of blocks, trigger an epoch.
	if len(chain.blocks)%10 == 0 {
		chain.epochChange.L.Lock()
		chain.epochChange.Broadcast()
		chain.epochChange.L.Unlock()
	}
	chain.blockMu.Unlock()

	// Update the nonces
	chain.nonceMu.Lock()
	for _, tx := range block.Txs {
		for _, nonce := range tx.Nonces {
			chain.nonces[nonce] = struct{}{}
		}
	}
	chain.nonceMu.Unlock()

	return nil
}

func (chain *HyperChain) ListenForEpochChange(ctx context.Context) <-chan struct{} {
	epochChange := make(chan struct{})
	go func() {
		// todo : does this goroutine get closed gracefully?
		defer close(epochChange)

		for {
			chain.epochChange.L.Lock()
			chain.epochChange.Wait()
			chain.epochChange.L.Unlock()

			select {
			case <-ctx.Done():
				return
			default:
				epochChange <- struct{}{}
			}
		}
	}()

	return epochChange
}

func (chain *HyperChain) VerifyBlock(block Block) bool {
	chain.nonceMu.RLock()
	defer chain.nonceMu.RUnlock()

	// Check block has no conflicts with nonces in the past blocks
	for _, tx := range block.Txs {
		for _, nonce := range tx.Nonces {
			if _, ok := chain.nonces[nonce]; ok {
				return false
			}
		}
	}

	return true
}

func (chain *HyperChain) VerifyTx(tx Tx) bool {
	chain.nonceMu.RLock()
	defer chain.nonceMu.RUnlock()

	// Check block has no conflicts with nonces in the past blocks
	for _, nonce := range tx.Nonces {
		if _, ok := chain.nonces[nonce]; ok {
			return false
		}
	}

	return true
}
