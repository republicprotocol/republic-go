package hyperdrive

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
				h := proposal.Hash()
				if _, ok := proposals[h]; ok {
					continue
				}
				chSetOut.Proposals <- proposal
				proposals[h] = struct{}{}

			case prepare, ok := <-chSetIn.Prepares:
				if !ok {
					return
				}
				h := prepare.Hash()
				if _, ok := prepares[h]; ok {
					continue
				}
				chSetOut.Prepares <- prepare
				prepares[h] = struct{}{}

			case commit, ok := <-chSetIn.Commits:
				if !ok {
					return
				}
				h := commit.Hash()
				if _, ok := commits[h]; ok {
					continue
				}
				chSetOut.Commits <- commit
				commits[h] = struct{}{}

			case block, ok := <-chSetIn.Blocks:
				if !ok {
					return
				}
				h := block.Hash()
				if _, ok := faults[h]; ok {
					continue
				}
				chSetOut.Blocks <- block
				faults[h] = struct{}{}

			case fault, ok := <-chSetIn.Faults:
				if !ok {
					return
				}
				h := fault.Hash()
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

// HeightCeiling defines the maximum height, with respective to the current
// height, that will be buffered. Above this limit, messages will be dropped.
const HeightCeiling = 1000

// FilterHeight consumes a ChannelSet and filters all messages that are not
// for the current height. The filtered messages are buffered and reconsidered
// when the height is changed.
func FilterHeight(chSetIn ChannelSet, height <-chan Height, capacity int) ChannelSet {
	chSetOut := NewChannelSet(capacity)

	go func() {
		defer chSetOut.Close()

		buffer := map[Height]ChannelSet{}
		defer func() {
			for _, chSet := range buffer {
				chSet.Close()
			}
		}()

		h := <-height

		for {
			select {
			case nextH, ok := <-height:
				if !ok {
					return
				}
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
				if proposal.Height < h || proposal.Height > h+HeightCeiling {
					continue
				}
				if proposal.Height == h {
					chSetOut.Proposals <- proposal
					continue
				}
				if _, ok := buffer[proposal.Height]; !ok {
					buffer[proposal.Height] = NewChannelSet(capacity)
				}
				buffer[proposal.Height].Proposals <- proposal

			case prepare, ok := <-chSetIn.Prepares:
				if !ok {
					return
				}
				if prepare.Height < h || prepare.Height > h+HeightCeiling {
					continue
				}
				if prepare.Height == h {
					chSetOut.Prepares <- prepare
					continue
				}
				if _, ok := buffer[prepare.Height]; !ok {
					buffer[prepare.Height] = NewChannelSet(capacity)
				}
				buffer[prepare.Height].Prepares <- prepare

			case commit, ok := <-chSetIn.Commits:
				if !ok {
					return
				}
				if commit.Height < h || commit.Height > h+HeightCeiling {
					continue
				}
				if commit.Height == h {
					chSetOut.Commits <- commit
					continue
				}
				if _, ok := buffer[commit.Height]; !ok {
					buffer[commit.Height] = NewChannelSet(capacity)
				}
				buffer[commit.Height].Commits <- commit

			case block, ok := <-chSetIn.Blocks:
				if !ok {
					return
				}
				if block.Height < h || block.Height > h+HeightCeiling {
					continue
				}
				if block.Height == h {
					chSetOut.Blocks <- block
					continue
				}
				if _, ok := buffer[block.Height]; !ok {
					buffer[block.Height] = NewChannelSet(capacity)
				}
				buffer[block.Height].Blocks <- block

			case fault, ok := <-chSetIn.Faults:
				if !ok {
					return
				}
				if fault.Height < h || fault.Height > h+HeightCeiling {
					continue
				}
				if fault.Height == h {
					chSetOut.Faults <- fault
					continue
				}
				if _, ok := buffer[fault.Height]; !ok {
					buffer[fault.Height] = NewChannelSet(capacity)
				}
				buffer[fault.Height].Faults <- fault
			}
		}
	}()

	return chSetOut
}
