package hyper

import (
	"sync"

	"github.com/republicprotocol/republic-go/dispatch"
)

type ChannelSet struct {
	BufferSize uint8
	Proposal   chan Proposal
	Prepare    chan Prepare
	Fault      chan Fault
	Commit     chan Commit
	Err        chan error
	Block      chan Block
}

func NewChannelSet(size uint8, proposal chan Proposal, prepare chan Prepare, commit chan Commit, fault chan Fault, block chan Block, err chan error) ChannelSet {
	return ChannelSet{
		BufferSize: size,
		Proposal:   proposal,
		Prepare:    prepare,
		Commit:     commit,
		Fault:      fault,
		Block:      block,
		Err:        err,
	}
}

func EmptyChannelSet(size uint8) ChannelSet {
	return ChannelSet{
		Proposal: make(chan Proposal, size),
		Prepare:  make(chan Prepare, size),
		Fault:    make(chan Fault, size),
		Commit:   make(chan Commit, size),
		Err:      make(chan error, size),
		Block:    make(chan Block, size),
	}
}

func (c *ChannelSet) Close() {
	close(c.Proposal)
	close(c.Prepare)
	close(c.Fault)
	close(c.Commit)
	close(c.Err)
	close(c.Block)
}

func (c *ChannelSet) Split(cs []ChannelSet) {
	var wg sync.WaitGroup

	proposals := make([]chan Proposal, len(cs))
	prepares := make([]chan Prepare, len(cs))
	commits := make([]chan Commit, len(cs))
	faults := make([]chan Fault, len(cs))
	errs := make([]chan error, len(cs))
	blocks := make([]chan Block, len(cs))

	for i, chset := range cs {
		proposals[i] = chset.Proposal
		prepares[i] = chset.Prepare
		commits[i] = chset.Commit
		faults[i] = chset.Fault
		errs[i] = chset.Err
		blocks[i] = chset.Block
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		dispatch.Split(c.Proposal, proposals)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		dispatch.Split(c.Prepare, prepares)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		dispatch.Split(c.Commit, commits)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		dispatch.Split(c.Fault, faults)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		dispatch.Split(c.Block, blocks)
	}()

	wg.Wait()
}

func (c *ChannelSet) Copy(cs ChannelSet) {
	go func() {
		defer c.Close()
		for {
			select {
			case proposal, ok := <-cs.Proposal:
				if !ok {
					return
				}
				c.Proposal <- proposal
			case prepare, ok := <-cs.Prepare:
				if !ok {
					return
				}
				c.Prepare <- prepare
			case commit, ok := <-cs.Commit:
				if !ok {
					return
				}
				c.Commit <- commit
			case fault, ok := <-cs.Fault:
				if !ok {
					return
				}
				c.Fault <- fault

			case block, ok := <-cs.Block:
				if !ok {
					return
				}
				c.Block <- block
			}
		}
	}()
}
