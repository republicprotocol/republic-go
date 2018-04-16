package hyper

import (
	"sync"

	"github.com/republicprotocol/republic-go/dispatch"
)

// The ChannelSet struct aggregates a set of channels that are commonly used
// together. All channels in the ChannelSet must be closed together using the
// Close method.
type ChannelSet struct {
	Proposals chan Proposal
	Prepares  chan Prepare
	Commits   chan Commit
	Blocks    chan Block
	Faults    chan Fault
}

// NewChannelSet returns Channels with the given capacity.
func NewChannelSet(capacity int) ChannelSet {
	return ChannelSet{
		Proposals: make(chan Proposal, capacity),
		Prepares:  make(chan Prepare, capacity),
		Commits:   make(chan Commit, capacity),
		Blocks:    make(chan Block, capacity),
		Faults:    make(chan Fault, capacity),
	}
}

// Close all channels in the ChannelSet.
func (chSet *ChannelSet) Close() {
	close(chSet.Proposals)
	close(chSet.Prepares)
	close(chSet.Commits)
	close(chSet.Blocks)
	close(chSet.Faults)
}

// Split the ChannelSet into multipl ChannelSets. The output ChannelSets must
// not be closed before the input ChannelSet. Once the input ChannelSet is
// closed, this method will eventually terminate after all values have been
// piped to the output ChannelSet. This method will block the current
// goroutine.
func (chSet *ChannelSet) Split(chSetsIn ...ChannelSet) {
	proposals := make([]chan Proposal, len(chSetsIn))
	prepares := make([]chan Prepare, len(chSetsIn))
	commits := make([]chan Commit, len(chSetsIn))
	blocks := make([]chan Block, len(chSetsIn))
	faults := make([]chan Fault, len(chSetsIn))
	for i := range chSetsIn {
		proposals[i] = chSetsIn[i].Proposals
		prepares[i] = chSetsIn[i].Prepares
		commits[i] = chSetsIn[i].Commits
		blocks[i] = chSetsIn[i].Blocks
		faults[i] = chSetsIn[i].Faults
	}

	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		dispatch.Split(chSet.Proposals, proposals)
	}()
	go func() {
		defer wg.Done()
		dispatch.Split(chSet.Prepares, prepares)
	}()
	go func() {
		defer wg.Done()
		dispatch.Split(chSet.Commits, commits)
	}()
	go func() {
		defer wg.Done()
		dispatch.Split(chSet.Blocks, blocks)
	}()
	dispatch.Split(chSet.Faults, faults)
	wg.Wait()
}

// Pipe all message from the ChannelSet to another ChannelSet. The output
// ChannelSet must not be closed before the input ChannelSet. Once the input
// ChannelSet is closed, this method will eventually terminate after all values
// have been piped to the output ChannelSet. This method will block the current
// goroutine.
func (chSet *ChannelSet) Pipe(chSetOut ChannelSet) {
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		for proposal := range chSet.Proposals {
			chSetOut.Proposals <- proposal
		}
	}()
	go func() {
		defer wg.Done()
		for prepare := range chSet.Prepares {
			chSetOut.Prepares <- prepare
		}
	}()
	go func() {
		defer wg.Done()
		for commit := range chSet.Commits {
			chSetOut.Commits <- commit
		}
	}()
	go func() {
		defer wg.Done()
		for block := range chSet.Blocks {
			chSetOut.Blocks <- block
		}
	}()
	for fault := range chSet.Faults {
		chSetOut.Faults <- fault
	}
	wg.Wait()
}
