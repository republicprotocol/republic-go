package xing

import (
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-order-compute"
)

type Inbox struct {
	do.GuardedObject

	results      []*compute.Result
	newResults   []*compute.Result
	hasNewResult *do.Guard
}

func NewInbox() *Inbox {
	inbox := new(Inbox)
	inbox.GuardedObject = do.NewGuardedObject()
	inbox.results = make([]*compute.Result, 0)
	inbox.newResults = make([]*compute.Result, 0)
	inbox.hasNewResult = inbox.Guard(func() bool { return len(inbox.newResults) > 0 })
	return inbox
}

func (inbox *Inbox) AddNewResult(result *compute.Result) {
	inbox.Enter(nil)
	defer inbox.Exit()
	inbox.results = append(inbox.results, result)
	inbox.newResults = append(inbox.newResults, result)
}

func (inbox *Inbox) GetAllNewResults() []*compute.Result {
	inbox.Enter(inbox.hasNewResult)
	defer inbox.Exit()
	ret := make([]*compute.Result, 0, len(inbox.newResults))
	ret = append(ret, inbox.newResults...)
	inbox.newResults = inbox.newResults[:0]
	return ret
}

func (inbox *Inbox) GetAllResults() []*compute.Result{
	inbox.Enter(nil)
	defer inbox.Exit()
	ret := make([]*compute.Result, 0, len(inbox.results))
	ret = append(ret, inbox.results...)
	inbox.newResults = inbox.newResults[:0]
	return ret
}
