package dark

import (
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-order-compute"
)

type Inbox struct {
	do.GuardedObject

	r            int
	results      []*compute.Result
	hasNewResult *do.Guard
}

func NewInbox() *Inbox {
	inbox := new(Inbox)
	inbox.GuardedObject = do.NewGuardedObject()
	inbox.r = 0
	inbox.results = make([]*compute.Result, 0)
	inbox.hasNewResult = inbox.Guard(func() bool { return r < len(inbox.results) })
	return inbox
}

func (inbox *Inbox) AddNewResult(result *compute.Result) {
	inbox.Enter(nil)
	defer inbox.Exit()
	inbox.results = append(inbox.results, result)
}

func (inbox *Inbox) GetAllNewResults() []*compute.Result {
	inbox.Enter(inbox.hasNewResult)
	defer inbox.Exit()
	ret := make([]*compute.Result, 0, len(inbox.results)-r)
	ret = append(ret, inbox.results[r:]...)
	inbox.r = len(inbox.results)
	return ret
}

func (inbox *Inbox) GetAllResults() []*compute.Result {
	inbox.Enter(nil)
	defer inbox.Exit()
	ret := make([]*compute.Result, 0, len(inbox.results))
	ret = append(ret, inbox.results...)
	inbox.r = len(inbox.results)
	return ret
}
