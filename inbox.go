package dark

import (
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-order-compute"
)

// Inbox stores the results for a certain trader address
// It keeps a record of unread results. The guardedObject
// ensures it is safe to use concurrently
type Inbox struct {
	do.GuardedObject

	r            int
	results      []*compute.Result
	hasNewResult *do.Guard
}

// NewInbox creates an empty Inbox.
func NewInbox() *Inbox {
	inbox := new(Inbox)
	inbox.GuardedObject = do.NewGuardedObject()
	inbox.r = 0
	inbox.results = make([]*compute.Result, 0)
	inbox.hasNewResult = inbox.Guard(func() bool { return inbox.r < len(inbox.results) })
	return inbox
}

// AddNewResult adds new result into the inbox
func (inbox *Inbox) AddNewResult(result *compute.Result) {
	inbox.Enter(nil)
	defer inbox.Exit()
	inbox.results = append(inbox.results, result)
}

// GetAllNewResults returns all the unread results
func (inbox *Inbox) GetAllNewResults() []*compute.Result {
	inbox.Enter(inbox.hasNewResult)
	defer inbox.Exit()
	ret := make([]*compute.Result, 0, len(inbox.results)-inbox.r)
	ret = append(ret, inbox.results[inbox.r:]...)
	inbox.r = len(inbox.results)
	return ret
}

// GetAllResults returns all the results in the inbox
func (inbox *Inbox) GetAllResults() []*compute.Result {
	inbox.Enter(nil)
	defer inbox.Exit()
	ret := make([]*compute.Result, 0, len(inbox.results))
	ret = append(ret, inbox.results...)
	inbox.r = len(inbox.results)
	return ret
}
