package do

import (
	"sync"
)

// A Guard is a conditional variable that has be specialized for the
// GuardedObject. It should only ever be created by using a GuardedObject and
// should only ever be used when entering a GuardedObject.
type Guard struct {
	condition       func() bool
	conditionSignal chan struct{}
}

func newGuard(condition func() bool) Guard {
	guard := Guard{
		condition:       condition,
		conditionSignal: make(chan struct{}, 1),
	}
	if guard.condition() {
		guard.conditionSignal <- struct{}{}
	}
	return guard
}

func (guard *Guard) open() {
	select {
	case guard.conditionSignal <- struct{}{}:
	default:
	}
}

func (guard *Guard) close() {
	select {
	case <-guard.conditionSignal:
	default:
	}
}

func (guard *Guard) wait() {
	<-guard.conditionSignal
}

// A GuardedObject is used to provide safe concurrent read and write access to
// an object that is accessed concurrently. To be protected, an object should
// call Enter at the beginning of every method that accesses it, and call Exit
// at the end of every method that accesses it. A GuardedObject can also define
// conditional Guards that will block at a call to Enter, until a given
// condition is met.
type GuardedObject struct {
	mu     *sync.RWMutex
	guards []Guard
}

// NewGuardedObject returns a new GuardedObject with no Guards defined.
func NewGuardedObject() GuardedObject {
	return GuardedObject{
		mu:     new(sync.RWMutex),
		guards: make([]Guard, 0),
	}
}

// Guard this GuardedObject with a certain condition. The conditional lambda
// function must not modify any non-local variables. The Guard returned can be
// passed to a call to the Enter function, blocking the call until the
// condition is met. The conditions will automatically be re-evaluated. A Guard
// must only be used with the GuardedObject that created it.
func (object *GuardedObject) Guard(condition func() bool) *Guard {
	guard := newGuard(condition)
	object.guards = append(object.guards, guard)
	return &object.guards[len(object.guards)-1]
}

// Enter the GuardedObject. This will acquire a write lock to the
// GuardedObject. If a Guard is passed, then the write lock will not be
// acquired unless the condition is met. Exit must always be called after the
// lock is no longer needed. Between the calls Enter and Exit no other
// goroutine will have access to the GuardedObject. The condition on the Guard
// is guaranteed to hold immediately after a goroutine has passed the Enter
// function.
func (object *GuardedObject) Enter(guard *Guard) {
	// If there is no Guard to wait on, acquire a write lock on the
	// GuardedObject.
	if guard == nil {
		object.mu.Lock()
		return
	}
	for {
		// Wait for the Guard to be signaled on the condition.
		guard.wait()
		// Acquire a write lock to the GuardedObject.
		object.mu.Lock()
		// Check that the condition is still true. It is possible that between
		// receiving the signal on the condition, and acquiring the write lock,
		// the conidition has changed.
		if guard.condition() {
			break
		}
		// If the condition is no longer true, release the write lock and try
		// again.
		object.mu.Unlock()
	}
}

// EnterReadOnly is similar to Enter, but it acquires a read lock. The
// protected object, and any value used by a conditional Guard, must not be
// modified. ExitReadOnly must be called instead of Exit.
func (object *GuardedObject) EnterReadOnly(guard *Guard) {
	if guard == nil {
		object.mu.RLock()
		return
	}
	for {
		guard.wait()
		object.mu.RLock()
		if guard.condition() {
			break
		}
		object.mu.RUnlock()
	}
}

// Exit the GuardedObject. All Guards attached to the GuardedObject will be
// re-evaluated before the GuardedObject can be entered again. This must not be
// called unless a call to Enter has already been made.
func (object *GuardedObject) Exit() {
	object.resolveGuards()
	object.mu.Unlock()
}

// ExitReadOnly is similar to Exit, but it must not be called unless a call to
// EnterReadOnly has already been made.
func (object *GuardedObject) ExitReadOnly() {
	object.resolveGuards()
	object.mu.RUnlock()
}

func (object *GuardedObject) resolveGuards() {
	for i := range object.guards {
		if object.guards[i].condition() {
			object.guards[i].open()
			continue
		}
		object.guards[i].close()
	}
}
