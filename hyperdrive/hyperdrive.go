package hyper

import (
	"bytes"
	"context"
	"encoding/binary"
	"log"

	"github.com/republicprotocol/republic-go/dispatch"
	"golang.org/x/crypto/sha3"
)

type OrderID string
type Rank uint64

type Replica struct {
	ingress        ChannelSet
	internalEgress ChannelSet
	validator      Validator
}

func NewReplica(
	validator Validator,
	ingress ChannelSet,
) Replica {
	return Replica{
		ingress:        ingress,
		validator:      validator,
		internalEgress: EmptyChannelSet(validator.Threshold()),
	}
}

func (r *Replica) Run(ctx context.Context) ChannelSet {
	egress := EmptyChannelSet(r.validator.Threshold())
	go func() {
		internalIngress := EmptyChannelSet(r.validator.Threshold())
		go func() {
			internalIngress.Copy(ProcessBuffer(r.ingress, r.validator))
		}()
		go func() {
			egress.Copy(ProcessBroadcast(r.internalEgress, r.validator))
		}()
		dispatch.Wait(r.HandleProposals(ctx, internalIngress), r.HandlePrepares(ctx, internalIngress), r.HandleCommits(ctx, internalIngress))
	}()
	return egress
}

func (r *Replica) HandleProposals(ctx context.Context, ingress ChannelSet) chan struct{} {
	doneCh := make(chan struct{})
	prepCh, faultCh, errCh := ProcessProposal(ctx, ingress.Proposal, r.validator)
	go func() {
		defer close(doneCh)
		for {
			select {
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
				r.internalEgress.Prepare <- prepare
			case fault, ok := <-faultCh:
				if !ok {
					return
				}
				r.internalEgress.Fault <- fault
			}
		}
	}()
	return doneCh
}

func (r *Replica) HandlePrepares(ctx context.Context, ingress ChannelSet) chan struct{} {
	doneCh := make(chan struct{})
	commCh, faultCh, errCh := ProcessPreparation(ctx, ingress.Prepare, r.validator)
	go func() {
		defer close(doneCh)
		for {
			select {
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
				r.internalEgress.Commit <- commit
			case fault, ok := <-faultCh:
				if !ok {
					return
				}
				r.internalEgress.Fault <- fault
			}
		}
	}()
	return doneCh
}

func (r *Replica) HandleCommits(ctx context.Context, ingress ChannelSet) chan struct{} {
	doneCh := make(chan struct{})
	blockCh := make(chan Block)
	counter := 0
	commCh, blockCh, faultCh, errCh := ProcessCommit(ctx, ingress.Commit, r.validator)
	go func() {
		defer close(doneCh)
		for {
			select {
			case <-ctx.Done():
				return
			case _ = <-errCh:
				return
			case block, ok := <-blockCh:
				if !ok {
					return
				}
				counter++
				log.Println("Finality reached for block number", counter, "on", r.validator.Sign())
				r.internalEgress.Block <- block
			case commit, ok := <-commCh:
				if !ok {
					return
				}
				r.internalEgress.Commit <- commit
			case fault, ok := <-faultCh:
				if !ok {
					return
				}
				log.Println("Fault block reached on", r.validator.Sign())
				r.internalEgress.Fault <- fault
			}
		}
	}()
	return doneCh
}

func (r *Replica) HandleBlocks(ctx context.Context) chan struct{} {
	log.Println("In handle blocks")
	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		ConsumeCertifiedBlocks(r.ingress.Block, r.validator.SharedBlocks())
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
