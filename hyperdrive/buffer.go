package hyper

import (
	"context"
	"sync"
	"time"
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
	doneCh := make(chan struct{})
	sb := validator.SharedBlocks()
	buffer := NewBuffer()
	go func() {
		defer close(doneCh)
		defer time.Sleep(10 * time.Second)
		for {
			h := sb.ReadHeight()
			select {
			case proposal, ok := <-chanSetIn.Proposal:
				if !ok {
					return
				}
				if proposal.Height < h {
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
					buffer.chanSets[proposal.Height] = EmptyChannelSet(ctx, validator.Threshold())
				}
				buffer.chanSets[proposal.Height].Proposal <- proposal
				buffer.channelSetsMu.Unlock()

			case prepare, ok := <-chanSetIn.Prepare:
				if !ok {
					return
				}
				if prepare.Height < h {
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
					buffer.chanSets[prepare.Height] = EmptyChannelSet(ctx, validator.Threshold())
				}
				buffer.chanSets[prepare.Height].Prepare <- prepare
				buffer.channelSetsMu.Unlock()

			case commit, ok := <-chanSetIn.Commit:
				if !ok {
					return
				}
				if commit.Height < h {
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
					buffer.chanSets[commit.Height] = EmptyChannelSet(ctx, validator.Threshold())
				}
				buffer.chanSets[commit.Height].Commit <- commit
				buffer.channelSetsMu.Unlock()

			case fault, ok := <-chanSetIn.Fault:
				if !ok {
					return
				}
				if fault.Height < h {
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
					buffer.chanSets[fault.Height] = EmptyChannelSet(ctx, validator.Threshold())
				}
				buffer.chanSets[fault.Height].Fault <- fault
				buffer.channelSetsMu.Unlock()
			}
		}
	}()
	return buffer, doneCh
}

func ConsumeBuffer(buffer Buffer, doneCh chan struct{}, validator Validator) ChannelSet {

	sb := validator.SharedBlocks()

	ctx, cancel := context.WithCancel(context.Background())
	chanSetOut := EmptyChannelSet(ctx, validator.Threshold())
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
		buffer.chanSets[height] = EmptyChannelSet(ctx, validator.Threshold())
	}
	go chanSetOut.Copy(buffer.chanSets[0])
	buffer.channelSetsMu.Unlock()

	go func() {
		defer cancel()
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
					go chanSetOut.Copy(buffer.chanSets[newHeight])
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
