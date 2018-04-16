package hyper

func ProcessBroadcast(chanSetIn ChannelSet, validator Validator) ChannelSet {
	chanSetOut := NewChannelSet(validator.Threshold())
	broadcastedProposals := map[[32]byte]bool{}
	broadcastedPrepares := map[[32]byte]bool{}
	broadcastedCommits := map[[32]byte]bool{}
	broadcastedFaults := map[[32]byte]bool{}
	broadcastedBlocks := map[[32]byte]bool{}
	go func() {

		for {
			select {
			case proposal, ok := <-chanSetIn.Proposals:
				if !ok {
					return
				}
				h := ProposalHash(proposal)
				if broadcastedProposals[h] {
					continue
				}
				chanSetOut.Proposals <- proposal
				broadcastedProposals[h] = true
			case prepare, ok := <-chanSetIn.Prepares:
				if !ok {
					return
				}
				h := PrepareHash(prepare)
				if broadcastedPrepares[h] {
					continue
				}
				chanSetOut.Prepares <- prepare
				broadcastedPrepares[h] = true
			case commit, ok := <-chanSetIn.Commits:
				if !ok {
					return
				}
				h := CommitHash(commit)
				if broadcastedCommits[h] {
					continue
				}
				chanSetOut.Commits <- commit
				broadcastedCommits[h] = true
			case fault, ok := <-chanSetIn.Faults:
				if !ok {
					return
				}
				h := FaultHash(fault)
				if broadcastedFaults[h] {
					continue
				}
				chanSetOut.Faults <- fault
				broadcastedFaults[h] = true
			case block, ok := <-chanSetIn.Blocks:
				if !ok {
					return
				}
				h := BlockHash(block)
				if broadcastedBlocks[h] {
					continue
				}
				chanSetOut.Blocks <- block
				broadcastedBlocks[h] = true
			}
		}
	}()
	return chanSetOut
}
