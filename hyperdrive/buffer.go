package hyper

import (
	"context"
	"log"
	"sync"
	"time"
)

type HeightContext struct {
	ctx    context.Context
	cancel context.CancelFunc
}

var HeightContexts = map[uint64]HeightContext{}

type Buffer struct {
	mu       sync.RWMutex
	chanSets map[uint64]ChannelSet
}

func NewBuffer()
func ProcessBuffer(chanSetIn ChannelSet, validator Validator) ChannelSet {
	buffer, doneCh := ProduceBuffer(chanSetIn, validator)
	return ConsumeBuffer(buffer, doneCh, validator)
}

func ProduceBuffer(chanSetIn ChannelSet, validator Validator) (Buffer, chan struct{}) {
	doneCh := make(chan struct{})
	sb := validator.SharedBlocks()

	buffer.Store
	go func() {
		defer close(doneCh)
		defer time.Sleep(10 * time.Second)
		defer log.Println("Closing done channel")
		for {
			h := sb.ReadHeight()
			select {
			case proposal, ok := <-chanSetIn.Proposal:
				if !ok {
					log.Println("returning because of closed proposal channel")
					return
				}
				if proposal.Height < h {
					continue
				}
				if _, ok := buffer.proposals[h]; !ok {
					buffer.proposals[h] = make(chan Proposal, validator.Threshold())
				}
				buffer.proposals[h] <- proposal

			case prepare, ok := <-chanSetIn.Prepare:
				if !ok {
					return
				}
				if prepare.Height < h {
					continue
				}
				if _, ok := buffer.prepares[h]; !ok {
					buffer.prepares[h] = make(chan Prepare, validator.Threshold())
				}
				buffer.prepares[h] <- prepare

			case commit, ok := <-chanSetIn.Commit:
				if !ok {
					return
				}
				if commit.Height < h {
					continue
				}
				if _, ok := buffer.commits[h]; !ok {
					buffer.commits[h] = make(chan Commit, validator.Threshold())
				}
				buffer.commits[h] <- commit

			case fault, ok := <-chanSetIn.Fault:
				if !ok {
					return
				}
				if fault.Height < h {
					continue
				}
				if _, ok := buffer.faults[h]; !ok {
					buffer.faults[h] = make(chan Fault, validator.Threshold())
				}
				buffer.faults[h] <- fault
			}
		}
	}()
	return buffer, doneCh
}

func ConsumeBuffer(buffer Buffer, doneCh chan struct{}, validator Validator) ChannelSet {

	sb := validator.SharedBlocks()
	chanSetOut := EmptyChannelSet(validator.Threshold())
	height := sb.ReadHeight()

	chanSetOut.Proposal = buffer.proposals[height]
	chanSetOut.Prepare = buffer.prepares[height]
	chanSetOut.Commit = buffer.commits[height]
	chanSetOut.Fault = buffer.faults[height]

	go func() {
		for {
			select {
			case <-doneCh:
				if _, ok := HeightContexts[height]; ok {
					HeightContexts[height].cancel()
				}
				return
			default:

				if height == sb.ReadHeight() {
					continue
				} else {
					newHeight := sb.ReadHeight()
					ctx, cancel := context.WithCancel(context.Background())

					hctx := HeightContext{
						ctx:    ctx,
						cancel: cancel,
					}

					HeightContexts[newHeight] = hctx
					if _, ok := HeightContexts[height]; ok {
						HeightContexts[height].cancel()
					}

					chanSet := NewChannelSet(validator.Threshold(), buffer.proposals[newHeight], buffer.prepares[newHeight], buffer.commits[newHeight], buffer.faults[newHeight], make(chan Block, validator.Threshold()), make(chan error, validator.Threshold()))
					go chanSetOut.Copy(ctx, chanSet)

					for i := height; i < newHeight; i++ {
						delete(buffer.proposals, i)
						delete(buffer.prepares, i)
						delete(buffer.commits, i)
						delete(buffer.faults, i)
					}
					height = newHeight
				}
			}
		}
	}()

	return chanSetOut
}
