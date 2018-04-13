package dispatch

import (
	"sync"
)

// MergeErrors from multiple channels into one channel.
func MergeErrors(errChIns ...(<-chan error)) <-chan error {
	errCh := make(chan error, len(errChIns))

	go func() {
		defer close(errCh)

		var wg sync.WaitGroup
		wg.Add(len(errChIns))
		for _, errChIn := range errChIns {
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

// FilterErrors from a channel using a predicate function.
func FilterErrors(errChIn <-chan error, f func(err error) bool) <-chan error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)

		for err := range errChIn {
			errCh <- err
		}
	}()

	return errCh
}
