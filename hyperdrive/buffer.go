package hyper

import (
	"context"
	"sync"

	"golang.org/x/crypto/sha3"
)

type HeightContext struct {
	ctx    context.Context
	cancel context.CancelFunc
}

type Buffer struct {
	channelSetsMu *sync.RWMutex
	chanSets      map[uint64]ChannelSet

	HeightContextsMu *sync.RWMutex
	HeightContexts   map[uint64]HeightContext
}

func NewBuffer() Buffer {
	return Buffer{
		channelSetsMu:    &sync.RWMutex{},
		chanSets:         map[uint64]ChannelSet{},
		HeightContextsMu: &sync.RWMutex{},
		HeightContexts:   map[uint64]HeightContext{},
	}
}

func ProcessBuffer(chanSetIn ChannelSet, validator Validator) ChannelSet {
	buffer, doneCh := ProduceBuffer(chanSetIn, validator)
	return ConsumeBuffer(buffer, doneCh, validator)
}

func ProduceBuffer(chanSetIn ChannelSet, validator Validator) (Buffer, chan struct{}) {
	duplicateProposals := map[[32]byte]bool{}
	duplicatePrepares := map[[32]byte]bool{}
	duplicateCommits := map[[32]byte]bool{}
	duplicateFaults := map[[32]byte]bool{}

	doneCh := make(chan struct{})
	sb := validator.SharedBlocks()
	buffer := NewBuffer()
	go func() {
		defer close(doneCh)
		for {
			h := sb.ReadHeight()
			select {
			case proposal, ok := <-chanSetIn.Proposals:
				if !ok {
					return
				}
				proposalHash := ProposalHash(proposal)
				proposalID := sha3.Sum256(append([]byte(proposal.Signature), proposalHash[:]...))
				if proposal.Height < h || duplicateProposals[proposalID] {
					continue
				}
				buffer.channelSetsMu.Lock()
				if _, ok := buffer.chanSets[proposal.Height]; !ok {
					ctx, cancel := context.WithCancel(context.Background())
					buffer.HeightContextsMu.Lock()
					if _, ok := buffer.HeightContexts[proposal.Height]; !ok {
						buffer.HeightContexts[proposal.Height] = HeightContext{
							ctx:    ctx,
							cancel: cancel,
						}
					}
					buffer.HeightContextsMu.Unlock()
					buffer.chanSets[proposal.Height] = NewChannelSet(validator.Threshold())
				}
				// log.Println("proposal", validator.Sign(), proposal.Height)
				buffer.chanSets[proposal.Height].Proposals <- proposal
				duplicateProposals[proposalID] = true
				buffer.channelSetsMu.Unlock()

			case prepare, ok := <-chanSetIn.Prepares:
				if !ok {
					return
				}
				prepareHash := PrepareHash(prepare)
				prepareID := sha3.Sum256(append([]byte(prepare.Signature), prepareHash[:]...))
				if prepare.Height < h || duplicatePrepares[prepareID] {
					continue
				}

				buffer.channelSetsMu.Lock()
				if _, ok := buffer.chanSets[prepare.Height]; !ok {
					ctx, cancel := context.WithCancel(context.Background())
					buffer.HeightContextsMu.Lock()
					if _, ok := buffer.HeightContexts[prepare.Height]; !ok {
						buffer.HeightContexts[prepare.Height] = HeightContext{
							ctx:    ctx,
							cancel: cancel,
						}
					}
					buffer.HeightContextsMu.Unlock()
					buffer.chanSets[prepare.Height] = NewChannelSet(validator.Threshold())
				}
				buffer.chanSets[prepare.Height].Prepares <- prepare
				duplicatePrepares[prepareID] = true
				buffer.channelSetsMu.Unlock()

			case commit, ok := <-chanSetIn.Commits:
				if !ok {
					return
				}
				commitHash := CommitHash(commit)
				commitID := sha3.Sum256(append([]byte(commit.Signature), commitHash[:]...))
				if commit.Height < h || duplicateCommits[commitID] {
					continue
				}
				buffer.channelSetsMu.Lock()
				if _, ok := buffer.chanSets[commit.Height]; !ok {
					ctx, cancel := context.WithCancel(context.Background())
					buffer.HeightContextsMu.Lock()
					if _, ok := buffer.HeightContexts[commit.Height]; !ok {
						buffer.HeightContexts[commit.Height] = HeightContext{
							ctx:    ctx,
							cancel: cancel,
						}
					}
					buffer.HeightContextsMu.Unlock()
					buffer.chanSets[commit.Height] = NewChannelSet(validator.Threshold())
				}
				buffer.chanSets[commit.Height].Commits <- commit
				duplicateCommits[commitID] = true
				buffer.channelSetsMu.Unlock()

			case fault, ok := <-chanSetIn.Faults:
				if !ok {
					return
				}
				faultHash := FaultHash(fault)
				faultID := sha3.Sum256(append([]byte(fault.Signature), faultHash[:]...))
				if fault.Height < h || duplicateFaults[faultID] {
					continue
				}
				buffer.channelSetsMu.Lock()
				if _, ok := buffer.chanSets[fault.Height]; !ok {
					ctx, cancel := context.WithCancel(context.Background())
					buffer.HeightContextsMu.Lock()
					if _, ok := buffer.HeightContexts[fault.Height]; !ok {
						buffer.HeightContexts[fault.Height] = HeightContext{
							ctx:    ctx,
							cancel: cancel,
						}
					}
					buffer.HeightContextsMu.Unlock()
					buffer.chanSets[fault.Height] = NewChannelSet(validator.Threshold())
				}
				buffer.chanSets[fault.Height].Faults <- fault
				duplicateFaults[faultID] = true
				buffer.channelSetsMu.Unlock()
			}
		}
	}()
	return buffer, doneCh
}

func ConsumeBuffer(buffer Buffer, doneCh chan struct{}, validator Validator) ChannelSet {

	sb := validator.SharedBlocks()

	chanSetOut := NewChannelSet(validator.Threshold())
	height := sb.ReadHeight()

	buffer.channelSetsMu.Lock()
	if _, ok := buffer.chanSets[height]; !ok {
		hctx, hcancel := context.WithCancel(context.Background())
		buffer.HeightContextsMu.Lock()
		if _, ok := buffer.HeightContexts[height]; !ok {
			buffer.HeightContexts[height] = HeightContext{
				ctx:    hctx,
				cancel: hcancel,
			}
		}
		buffer.HeightContextsMu.Unlock()
		buffer.chanSets[height] = NewChannelSet(validator.Threshold())
	}
	go chanSetOut.Pipe(buffer.chanSets[0])
	buffer.channelSetsMu.Unlock()

	go func() {
		for {
			select {
			case <-doneCh:
				if _, ok := buffer.HeightContexts[height]; ok {
					buffer.HeightContexts[height].cancel()
				}
				return
			default:
				if height == sb.ReadHeight() {
					continue
				} else {
					newHeight := sb.ReadHeight()

					buffer.HeightContextsMu.Lock()
					if _, ok := buffer.HeightContexts[height]; ok {
						buffer.HeightContexts[height].cancel()
					}
					buffer.HeightContextsMu.Unlock()

					buffer.channelSetsMu.RLock()
					go chanSetOut.Pipe(buffer.chanSets[newHeight])
					buffer.channelSetsMu.RUnlock()

					for i := height; i < newHeight; i++ {

						buffer.channelSetsMu.Lock()
						if _, ok := buffer.chanSets[i]; ok {
							delete(buffer.chanSets, i)
						}
						buffer.channelSetsMu.Unlock()

						buffer.HeightContextsMu.Lock()
						if _, ok := buffer.HeightContexts[i]; ok {
							delete(buffer.HeightContexts, i)
						}
						buffer.HeightContextsMu.Unlock()

					}
					height = newHeight
				}
			}
		}
	}()

	return chanSetOut
}
