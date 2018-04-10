package hyper

func ProcessBroadcast(chanSetIn ChannelSet, chanSetOut ChannelSet) {
	broadcastedProposals := map[[32]byte]bool{}
	broadcastedPrepares := map[[32]byte]bool{}
	broadcastedCommits := map[[32]byte]bool{}
	broadcastedFaults := map[[32]byte]bool{}
	broadcastedBlocks := map[[32]byte]bool{}
	go func() {
		defer chanSetOut.Close()

		for {
			select {
			case proposal, ok := <-chanSetIn.Proposal:
				if !ok {
					return
				}
				h := ProposalHash(proposal)
				if broadcastedProposals[h] {
					continue
				}
				chanSetOut.Proposal <- proposal
				broadcastedProposals[h] = true
			case prepare, ok := <-chanSetIn.Prepare:
				if !ok {
					return
				}
				h := PrepareHash(prepare)
				if broadcastedPrepares[h] {
					continue
				}
				chanSetOut.Prepare <- prepare
				broadcastedPrepares[h] = true
			case commit, ok := <-chanSetIn.Commit:
				if !ok {
					return
				}
				h := CommitHash(commit)
				if broadcastedCommits[h] {
					continue
				}
				chanSetOut.Commit <- commit
				broadcastedCommits[h] = true
			case fault, ok := <-chanSetIn.Fault:
				if !ok {
					return
				}
				h := FaultHash(fault)
				if broadcastedFaults[h] {
					continue
				}
				chanSetOut.Fault <- fault
				broadcastedFaults[h] = true
			case block, ok := <-chanSetIn.Block:
				if !ok {
					return
				}
				h := BlockHash(block)
				if broadcastedBlocks[h] {
					continue
				}
				chanSetOut.Block <- block
				broadcastedBlocks[h] = true
			}
		}
	}()
}
