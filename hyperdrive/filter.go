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
			}
		}
	}()

	return chSetOut
}
