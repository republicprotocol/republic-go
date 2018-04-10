package dispatch

import (
	"sync"
)

// MergeErrors from multiple channels into one channel.
func MergeErrors(errChs ...(<-chan error)) <-chan error {
	errCh := make(chan error)
	go func() {
		defer close(errCh)

		var wg sync.WaitGroup
		wg.Add(len(errChs))
		for _, errChIn := range errChs {
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

// FilterErrors from a channel.
func FilterErrors(errChIn <-chan error, f func(error) bool) <-chan error {
	errCh := make(chan error)
	go func() {
		defer close(errCh)

		for err := range errChIn {
			if f(err) {
				errCh <- err
			}
		}
	}()
	return errCh
}

// ConsumeErrors from an error channel using a list of error handling
// functions.
func ConsumeErrors(errCh <-chan error, f ...func(err error)) {
	for err := range errCh {
		if f != nil && len(f) > 0 {
			for i := range f {
				f[i](err)
			}
		}
	}
}
