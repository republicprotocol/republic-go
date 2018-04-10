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
type Height uint64

type Replica struct {
	ingress         ChannelSet
	egress          ChannelSet
	internalIngress ChannelSet
	internalEgress  ChannelSet
	validator       Validator
	signer          Signer
	threshold       uint8
	sharedBlocks    SharedBlocks
}

func NewReplica(
	validator Validator,
	signer Signer,
	ingress ChannelSet,
	egress ChannelSet,
	threshold uint8,
) Replica {
	return Replica{
		ingress:         ingress,
		egress:          egress,
		internalIngress: EmptyChannelSet(threshold),
		internalEgress:  EmptyChannelSet(threshold),
		validator:       validator,
		signer:          signer,
		threshold:       threshold,
		sharedBlocks:    NewSharedBlocks(0, 0),
	}
}

func (r *Replica) Run(ctx context.Context) {
	go func() {
		defer r.internalIngress.Close()
		defer r.internalEgress.Close()
		go func() {
			ProcessBuffer(r.ingress, r.internalIngress, &r.sharedBlocks)
		}()
		go func() {
			ProcessBroadcast(r.internalEgress, r.egress)
		}()
		dispatch.Wait(r.HandleProposals(ctx), r.HandlePrepares(ctx), r.HandleCommits(ctx))
	}()
}

func (r *Replica) HandleProposals(ctx context.Context) chan struct{} {
	doneCh := make(chan struct{})
	prepCh, faultCh, errCh := ProcessProposal(ctx, r.internalIngress.Proposal, r.signer, r.validator)
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

func (r *Replica) HandlePrepares(ctx context.Context) chan struct{} {
	doneCh := make(chan struct{})
	commCh, faultCh, errCh := ProcessPreparation(ctx, r.internalIngress.Prepare, r.signer, r.validator, r.threshold)
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

func (r *Replica) HandleCommits(ctx context.Context) chan struct{} {
	doneCh := make(chan struct{})
	blockCh := make(chan Block)
	counter := 0
	commCh, blockCh, faultCh, errCh := ProcessCommit(ctx, r.internalIngress.Commit, r.signer, r.validator, &r.sharedBlocks, r.threshold)
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
				log.Println("Finality reached for block number", counter, "on", r.signer.Sign())
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
				log.Println("Fault block reached on", r.signer.Sign())
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
		ConsumeCertifiedBlocks(r.ingress.Block, &r.sharedBlocks)
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
