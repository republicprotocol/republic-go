package hyper

// FilterDuplicates from a ChannelSet and produce a new ChannelSet where all
// messages, on all channels, are guaranteed to be unique.
func FilterDuplicates(chSetIn ChannelSet, capacity int) ChannelSet {
	chSetOut := NewChannelSet(capacity)

	go func() {
		defer chSetOut.Close()

		proposals := map[[32]byte]struct{}{}
		prepares := map[[32]byte]struct{}{}
		commits := map[[32]byte]struct{}{}
		blocks := map[[32]byte]struct{}{}
		faults := map[[32]byte]struct{}{}

		for {
			select {

			case proposal, ok := <-chSetIn.Proposals:
				if !ok {
					return
				}
				h := ProposalHash(proposal)
				if _, ok := proposals[h]; ok {
					continue
				}
				chSetOut.Proposals <- proposal
				proposals[h] = struct{}{}

			case prepare, ok := <-chSetIn.Prepares:
				if !ok {
					return
				}
				h := PrepareHash(prepare)
				if _, ok := prepares[h]; ok {
					continue
				}
				chSetOut.Prepares <- prepare
				prepares[h] = struct{}{}

			case commit, ok := <-chSetIn.Commits:
				if !ok {
					return
				}
				h := CommitHash(commit)
				if _, ok := commits[h]; ok {
					continue
				}
				chSetOut.Commits <- commit
				commits[h] = struct{}{}

			case block, ok := <-chSetIn.Blocks:
				if !ok {
					return
				}
				h := BlockHash(block)
				if _, ok := faults[h]; ok {
					continue
				}
				chSetOut.Blocks <- block
				faults[h] = struct{}{}

			case fault, ok := <-chSetIn.Faults:
				if !ok {
					return
				}
				h := FaultHash(fault)
				if _, ok := blocks[h]; ok {
					continue
				}
				chSetOut.Faults <- fault
				blocks[h] = struct{}{}
			}
		}
	}()

	return chSetOut
}

// FilterHeight consumes a ChannelSet and filters all messages that are not
// for the current height. The filtered messages are buffered and reconsidered
// when the height is changed.
func FilterHeight(chSetIn ChannelSet, height <-chan int, capacity int) ChannelSet {
	chSetOut := NewChannelSet(capacity)

	go func() {
		defer chSetOut.Close()

		buffer := map[int]ChannelSet{}
		defer func() {
			for _, chSet := range buffer {
				chSet.Close()
			}
		}()

		h := <-height

		for {
			select {
			case nextH := <-height:
				h = nextH
				if bufferedChSet, ok := buffer[h]; ok {
					delete(buffer, h)
					bufferedChSet.Close()
					bufferedChSet.Pipe(chSetOut)
				}

			case proposal, ok := <-chSetIn.Proposals:
				if !ok {
					return
				}
				if proposal.Height < h {
					continue
				}
				if proposal.Height == h {
					chSetOut.Proposals <- proposal
					continue
				}
				if _, ok := buffer[h]; !ok {
					buffer[h] = NewChannelSet(capacity)
				}
				buffer[h].Proposals <- proposal

			case prepare, ok := <-chSetIn.Prepares:
				if !ok {
					return
				}
				if prepare.Height < h {
					continue
				}
				if prepare.Height == h {
					chSetOut.Prepares <- prepare
					continue
				}
				if _, ok := buffer[h]; !ok {
					buffer[h] = NewChannelSet(capacity)
				}
				buffer[h].Prepares <- prepare

			case commit, ok := <-chSetIn.Commits:
				if !ok {
					return
				}
				if commit.Height < h {
					continue
				}
				if commit.Height == h {
					chSetOut.Commits <- commit
					continue
				}
				if _, ok := buffer[h]; !ok {
					buffer[h] = NewChannelSet(capacity)
				}
				buffer[h].Commits <- commit

			case block, ok := <-chSetIn.Blocks:
				if !ok {
					return
				}
				if block.Height < h {
					continue
				}
				if block.Height == h {
					chSetOut.Blocks <- block
					continue
				}
				if _, ok := buffer[h]; !ok {
					buffer[h] = NewChannelSet(capacity)
				}
				buffer[h].Blocks <- block

			case fault, ok := <-chSetIn.Faults:
				if !ok {
					return
				}
				if fault.Height < h {
					continue
				}
				if fault.Height == h {
					chSetOut.Faults <- fault
					continue
				}
				if _, ok := buffer[h]; !ok {
					buffer[h] = NewChannelSet(capacity)
				}
				buffer[h].Faults <- fault
			}
		}
	}()

	return chSetOut
}
