package hyperdrive

import (
	"bytes"
	"context"
	"encoding/binary"
	"sync"

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
	return sha3.Sum256(buf.Bytes())
}

// Verify the Block message. Returns an error if the message is invalid,
// otherwise nil.
func (block *Block) Verify(verifier identity.Verifier) error {
	return verifier.VerifySignature(block.Signature)
}

type HyperChain struct {
	mu          *sync.RWMutex
	epochChange *sync.Cond

	blocks []Block
	nonces map[[32]byte]struct{}
}

func NewHyperChain() HyperChain {
	return HyperChain{
		mu:          new(sync.RWMutex),
		epochChange: sync.NewCond(new(sync.Mutex)),

		blocks: make([]Block, 0),
		nonces: map[[32]byte]struct{}{},
	}
}

func (chain *HyperChain) FinalizeBlock(block Block) error {
	chain.mu.Lock()
	defer chain.mu.Unlock()

	// Check block has no conflicts with nonces in the past blocks
	for _, tx := range block.Txs {
		for _, nonce := range tx.Nonces {
			if _, ok := chain.nonces[nonce]; ok {
				return errors.New("nonces already recorded in the hyper chain")
			}
		}
	}

	// Check the height is right
	if len(chain.blocks) != int(block.Height) {
		return errors.New("invalid block to add to the chain")
	}

	// Finalize the block in the hyperChain
	chain.blocks = append(chain.blocks, block)
	for _, tx := range block.Txs {
		for _, nonce := range tx.Nonces {
			chain.nonces[nonce] = struct{}{}
		}
	}

	// When we reach certain blocks, trigger an epoch.
	if len(chain.blocks)%10 == 0 {
		chain.epochChange.L.Lock()
		chain.epochChange.Broadcast()
		chain.epochChange.L.Unlock()
	}

	return nil
}

func (chain *HyperChain) Block(i int) (Block, error) {
	chain.mu.RLock()
	defer chain.mu.RUnlock()

	if len(chain.blocks) < i {
		return Block{}, errors.New("invalid block height")
	}

	return chain.blocks[i], nil
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
