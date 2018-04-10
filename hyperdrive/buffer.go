package hyper

func ProcessBuffer(chanSetIn ChannelSet, chanSetOut ChannelSet, sb *SharedBlocks) {

	height := sb.Height
	proposals := map[Height][]Proposal{}
	prepares := map[Height][]Prepare{}
	commits := map[Height][]Commit{}
	faults := map[Height][]Fault{}

	go func() {
		for {
			select {
			case proposal, ok := <-chanSetIn.Proposal:
				if !ok {
					return
				}
				if proposal.Height < sb.Height {
					continue
				}

				if proposal.Height == sb.Height {
					chanSetOut.Proposal <- proposal
				} else {
					proposals[proposal.Height] = append(proposals[proposal.Height], proposal)
				}
			case prepare, ok := <-chanSetIn.Prepare:
				if !ok {
					return
				}
				if prepare.Height < sb.Height {
					continue
				}
				if prepare.Height == sb.Height {
					chanSetOut.Prepare <- prepare
				} else {
					prepares[prepare.Height] = append(prepares[prepare.Height], prepare)
				}

			case commit, ok := <-chanSetIn.Commit:
				if !ok {
					return
				}
				if commit.Height < sb.Height {
					continue
				}
				if commit.Height == sb.Height {
					chanSetOut.Commit <- commit
				} else {
					commits[commit.Height] = append(commits[commit.Height], commit)
				}

			case fault, ok := <-chanSetIn.Fault:
				if !ok {
					return
				}
				if fault.Height < sb.Height {
					continue
				}
				if fault.Height == sb.Height {
					chanSetOut.Fault <- fault
				} else {
					faults[fault.Height] = append(faults[fault.Height], fault)
				}
			}
		}
	}()

	go func() {
		for {
			newHeight := sb.ReadHeight()
			if height == newHeight {
				continue
			} else {
				for _, proposal := range proposals[newHeight] {
					chanSetOut.Proposal <- proposal
				}
				for _, prepare := range prepares[newHeight] {
					chanSetOut.Prepare <- prepare
				}
				for _, commit := range commits[newHeight] {
					chanSetOut.Commit <- commit
				}
				for _, fault := range faults[newHeight] {
					chanSetOut.Fault <- fault
				}

				for i := height + 1; i <= newHeight; i++ {
					delete(proposals, i)
					delete(prepares, i)
					delete(commits, i)
					delete(faults, i)
				}
				height = newHeight
			}
		}
	}()
}
