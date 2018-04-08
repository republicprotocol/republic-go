package hyper

import (
	"log"
	"time"
)

func ProcessHeightBuffer(proposalChIn <-chan Proposal, prepareChIn <-chan Prepare, commitChIn <-chan Commit, faultChIn <-chan Fault, v Validator) (<-chan Proposal, <-chan Prepare, <-chan Commit, <-chan Fault) {

	sb := v.GetSharedBlocks()
	height := sb.Height
	proposalChs := map[Height]chan Proposal{}
	proposalChOut := make(chan Proposal)
	prepareChs := map[Height]chan Prepare{}
	prepareChOut := make(chan Prepare)
	commitChs := map[Height]chan Commit{}
	commitChOut := make(chan Commit)
	faultChs := map[Height]chan Fault{}
	faultChOut := make(chan Fault)

	proposalChs[height] = make(chan Proposal)
	prepareChs[height] = make(chan Prepare)
	commitChs[height] = make(chan Commit)
	faultChs[height] = make(chan Fault)

	go func() {
		for {
			select {
			case proposal, ok := <-proposalChIn:
				if !ok {
					return
				}
				log.Println("reading from proposal channel", proposal.Height, sb.Height)
				proposalChs[proposal.Height] <- proposal

			case prepare, ok := <-prepareChIn:
				if !ok {
					return
				}
				prepareChs[prepare.Height] <- prepare

			case commit, ok := <-commitChIn:
				if !ok {
					return
				}
				commitChs[commit.Height] <- commit

			case fault, ok := <-faultChIn:
				if !ok {
					return
				}
				faultChs[fault.Height] <- fault
			}
		}
	}()

	go func() {
		for {
			select {
			case proposal, ok := <-proposalChs[sb.Height]:
				log.Println("writing to proposal channel", sb.Height)
				if !ok {
					return
				}
				proposalChOut <- proposal
			case prepare, ok := <-prepareChs[sb.Height]:
				if !ok {
					return
				}
				prepareChOut <- prepare

			case commit, ok := <-commitChs[sb.Height]:
				if !ok {
					return
				}
				commitChOut <- commit

			case fault, ok := <-faultChs[sb.Height]:
				if !ok {
					return
				}
				faultChOut <- fault
			}
		}
	}()

	go func() {
		for {
			if height == sb.Height {
				time.Sleep(1 * time.Minute)
			} else {
				close(proposalChs[height])
				close(prepareChs[height])
				close(commitChs[height])
				close(faultChs[height])

				delete(proposalChs, height)
				delete(prepareChs, height)
				delete(commitChs, height)
				delete(faultChs, height)
				height++

				proposalChs[height] = make(chan Proposal)
				prepareChs[height] = make(chan Prepare)
				commitChs[height] = make(chan Commit)
				faultChs[height] = make(chan Fault)
			}
		}
	}()

	return proposalChOut, prepareChOut, commitChOut, faultChOut
}
