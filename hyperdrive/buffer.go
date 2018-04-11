package hyper

import (
	"log"
	"sync"
)

func ProcessBuffer(chanSetIn ChannelSet, validator Validator) ChannelSet {

	sb := validator.SharedBlocks()
	threshold := validator.Threshold()
	height := sb.ReadHeight()
	chanSetOut := EmptyChannelSet(threshold)
	proposals := map[uint64][]Proposal{}
	prepares := map[uint64][]Prepare{}
	commits := map[uint64][]Commit{}
	faults := map[uint64][]Fault{}
	var wg sync.WaitGroup

	go func() {
		defer chanSetOut.Close()
		// log.Println("Hello")
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer log.Println("Closing buffer")
			for {
				select {
				case proposal, ok := <-chanSetIn.Proposal:
					log.Println("Buffering proposals")
					if !ok {
						return
					}
					if proposal.Height < sb.ReadHeight() {
						continue
					}
					// log.Println("Sorting proposals")
					if proposal.Height == sb.ReadHeight() {
						// log.Println("Writing proposals")
						chanSetOut.Proposal <- proposal
					} else {
						proposals[proposal.Height] = append(proposals[proposal.Height], proposal)
					}
				case prepare, ok := <-chanSetIn.Prepare:
					if !ok {
						return
					}
					if prepare.Height < sb.ReadHeight() {
						continue
					}
					// log.Println("Sorting prepares")
					if prepare.Height == sb.ReadHeight() {
						chanSetOut.Prepare <- prepare
					} else {
						prepares[prepare.Height] = append(prepares[prepare.Height], prepare)
					}

				case commit, ok := <-chanSetIn.Commit:
					if !ok {
						return
					}
					if commit.Height < sb.ReadHeight() {
						continue
					}
					if commit.Height == sb.ReadHeight() {
						chanSetOut.Commit <- commit
					} else {
						commits[commit.Height] = append(commits[commit.Height], commit)
					}

				case fault, ok := <-chanSetIn.Fault:
					if !ok {
						return
					}
					if fault.Height < sb.ReadHeight() {
						continue
					}
					if fault.Height == sb.ReadHeight() {
						chanSetOut.Fault <- fault
					} else {
						faults[fault.Height] = append(faults[fault.Height], fault)
					}
				}
			}
		}()

		wg.Add(1)
		go func() {
			log.Println("Hello2")
			defer wg.Done()
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
		wg.Wait()
	}()

	return chanSetOut
}
