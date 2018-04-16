package hyper

import (
	"bytes"
	"context"
	"encoding/binary"
	"log"
	"time"

	"github.com/republicprotocol/republic-go/dispatch"
	"golang.org/x/crypto/sha3"
)

type OrderID string
type Rank uint64

type Replica struct {
	ctx            context.Context
	ingress        ChannelSet
	internalEgress ChannelSet
	validator      Validator
}

func NewReplica(ctx context.Context, validator Validator, ingress ChannelSet) Replica {
	return Replica{
		ctx:            ctx,
		ingress:        ingress,
		validator:      validator,
		internalEgress: NewChannelSet(validator.Threshold()),
	}
}

func (r *Replica) Run() ChannelSet {
	egress := NewChannelSet(r.validator.Threshold())
	go func() {
		internalIngress := NewChannelSet(r.validator.Threshold())
		go internalIngress.Pipe(ProcessBuffer(r.ingress, r.validator))
		go egress.Pipe(FilterDuplicates(r.internalEgress, r.validator.Threshold()))
		dispatch.Wait(r.HandleProposals(r.ctx, internalIngress), r.HandlePrepares(r.ctx, internalIngress), r.HandleCommits(r.ctx, internalIngress))
	}()
	return egress
}

func (r *Replica) HandleProposals(ctx context.Context, ingress ChannelSet) chan struct{} {
	doneCh := make(chan struct{})
	prepCh, faultCh, errCh := ProcessProposal(ctx, ingress.Proposals, r.validator)
	go func() {
		defer close(doneCh)
		for {
			select {
			case <-time.After(10 * time.Second):
				log.Println("Replica timedout")
				return
			case <-ctx.Done():
				return
			case _, ok := <-errCh:
				if !ok {
					return
				}
				return
			case prepare, ok := <-prepCh:
				if !ok {
					return
				}
				r.internalEgress.Prepares <- prepare
			case fault, ok := <-faultCh:
				if !ok {
					return
				}
				r.internalEgress.Faults <- fault
			}
		}
	}()
	return doneCh
}

func (r *Replica) HandlePrepares(ctx context.Context, ingress ChannelSet) chan struct{} {
	doneCh := make(chan struct{}, r.validator.Threshold())
	commCh, faultCh, errCh := ProcessPreparation(ctx, ingress.Prepares, r.validator)
	go func() {
		defer close(doneCh)
		for {
			select {
			case <-time.After(10 * time.Second):
				log.Println("Replica timedout")
				return
			case <-ctx.Done():
				return
			case _, ok := <-errCh:
				if !ok {
					return
				}
				return
			case commit, ok := <-commCh:
				if !ok {
					return
				}
				r.internalEgress.Commits <- commit
			case fault, ok := <-faultCh:
				if !ok {
					return
				}
				r.internalEgress.Faults <- fault
			}
		}
	}()
	return doneCh
}

func (r *Replica) HandleCommits(ctx context.Context, ingress ChannelSet) chan struct{} {
	doneCh := make(chan struct{})
	blockCh := make(chan Block, r.validator.Threshold())
	counter := 0
	commCh, blockCh, faultCh, errCh := ProcessCommit(ctx, ingress.Commits, r.validator)
	go func() {
		defer close(doneCh)
		for {
			select {
			case <-time.After(10 * time.Second):
				log.Println("Replica timedout")
				return
			case <-ctx.Done():
				return
			case _ = <-errCh:
				return
			case block, ok := <-blockCh:
				if !ok {
					return
				}
				counter++
				// log.Printf("%sFinality reached on block%s\n", "\x1b[32;1m", r.validator.Sign())
				r.internalEgress.Blocks <- block
			case commit, ok := <-commCh:
				if !ok {
					return
				}
				r.internalEgress.Commits <- commit
			case fault, ok := <-faultCh:
				if !ok {
					return
				}
				log.Println("Fault block reached on", r.validator.Sign())
				r.internalEgress.Faults <- fault
			}
		}
	}()
	return doneCh
}

func PrepareHash(p Prepare) [32]byte {
	var prepareBuf bytes.Buffer
	binary.Write(&prepareBuf, binary.BigEndian, p.Block)
	binary.Write(&prepareBuf, binary.BigEndian, p.Height)
	binary.Write(&prepareBuf, binary.BigEndian, p.Rank)
	return sha3.Sum256(prepareBuf.Bytes())
}

func ProposalHash(p Proposal) [32]byte {
	var proposalBuf bytes.Buffer
	binary.Write(&proposalBuf, binary.BigEndian, p.Block)
	binary.Write(&proposalBuf, binary.BigEndian, p.Height)
	binary.Write(&proposalBuf, binary.BigEndian, p.Rank)
	return sha3.Sum256(proposalBuf.Bytes())
}

func CommitHash(c Commit) [32]byte {
	var commitBuf bytes.Buffer
	binary.Write(&commitBuf, binary.BigEndian, c.Block)
	binary.Write(&commitBuf, binary.BigEndian, c.Height)
	binary.Write(&commitBuf, binary.BigEndian, c.Rank)
	binary.Write(&commitBuf, binary.BigEndian, c.ThresholdSignature)
	return sha3.Sum256(commitBuf.Bytes())
}

func FaultHash(f Fault) [32]byte {
	var faultBuf bytes.Buffer
	binary.Write(&faultBuf, binary.BigEndian, f.Height)
	binary.Write(&faultBuf, binary.BigEndian, f.Rank)
	return sha3.Sum256(faultBuf.Bytes())
}

func BlockHash(b Block) [32]byte {
	blockBuffer := new(bytes.Buffer)
	for i := range b.Tuples {
		binary.Write(blockBuffer, binary.BigEndian, b.Tuples[i].ID)
	}
	// log.Println("--------------------------------------", blockBuffer.Bytes())
	return sha3.Sum256(blockBuffer.Bytes())
}
