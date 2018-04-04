package hyper

func FilterPrepare(prepareChIn chan Prepare, predicate func(*Prepare) bool) chan Prepare {
	prepareCh := make(chan Prepare)
	go func() {
		defer close(prepareCh)
		for prepare := range prepareChIn {
			if !predicate(&prepare) {
				prepareCh <- prepare
			}
		}
	}()
	return prepareCh
}

func FilterCommit(commitChIn chan Commit, predicate func(*Commit) bool) chan Commit {
	commitCh := make(chan Commit)
	go func() {
		defer close(commitCh)
		for commit := range commitChIn {
			if !predicate(&commit) {
				commitCh <- commit
			}
		}
	}()
	return commitCh
}

func FilterFault(faultChIn chan Fault, predicate func(*Fault) bool) chan Fault {
	faultCh := make(chan Fault)
	go func() {
		defer close(faultCh)
		for fault := range faultChIn {
			if !predicate(&fault) {
				faultCh <- fault
			}
		}
	}()
	return faultCh
}
