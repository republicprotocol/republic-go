package dispatch

import (
	"sync"
)

// MergeErrors from multiple channels into one channel. The merged channel has
// a capacity equal to the number of channels being merged.
func MergeErrors(errChsIn ...(<-chan error)) <-chan error {
	errCh := make(chan error, len(errChsIn))
	go func() {
		defer close(errCh)

		var wg sync.WaitGroup
		wg.Add(len(errChsIn))
		for _, errChIn := range errChsIn {
			go func(errChIn <-chan error) {
				defer wg.Done()

				for err := range errChIn {
					errCh <- err
				}
			}(errChIn)
		}
		wg.Wait()
	}()
	return errCh
}

// FilterErrors from a channel. The filtered channel has the same capacity as
// the input channel.
func FilterErrors(errChIn <-chan error, predicate func(error) bool) <-chan error {
	errCh := make(chan error, cap(errChIn))
	go func() {
		defer close(errCh)

		for err := range errChIn {
			if predicate(err) {
				errCh <- err
			}
		}
	}()
	return errCh
}

// ConsumeErrors from an error channel using a list of error handling
// functions. This function blocks the current goroutine until the error
// channel is closed, and empty.
func ConsumeErrors(errCh <-chan error, f ...func(err error)) {
	for err := range errCh {
		if f != nil && len(f) > 0 {
			for i := range f {
				f[i](err)
			}
		}
	}
}
