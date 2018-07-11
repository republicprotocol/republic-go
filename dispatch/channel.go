package dispatch

import (
	"fmt"
	"reflect"
	"sync"
)

// Forward all values from an input channels into an output channel. Forward is
// blocking and panics when the input channel type do no match the output
// channel.
func Forward(done <-chan struct{}, in interface{}, out interface{}) {

	// Ensure that all arguments are compatible types
	if reflect.TypeOf(out).Kind() != reflect.Chan {
		panic(fmt.Sprintf("cannot forward into type %v", reflect.TypeOf(out)))
	}
	if reflect.TypeOf(in).Kind() != reflect.Chan {
		panic(fmt.Sprintf("cannot forward from type %v", reflect.TypeOf(in)))
	}
	if reflect.TypeOf(in).Elem().Kind() != reflect.TypeOf(out).Elem().Kind() {
		panic(fmt.Sprintf("cannot forward from type %v to type %v", reflect.TypeOf(in), reflect.TypeOf(out)))
	}

	for {
		// select {
		// case <-done:
		// case val, ok := <-in:
		// }
		chosen, val, ok := reflect.Select([]reflect.SelectCase{
			reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(done)},
			reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(in)},
		})
		if chosen == 0 || !ok {
			return
		}

		// select {
		// case <-done:
		// case out <- val:
		// }
		chosen, _, _ = reflect.Select([]reflect.SelectCase{
			reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(done)},
			reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(out), Send: val},
		})
		if chosen == 0 {
			return
		}
	}
}

// Merge multiple input channels into an output channel. Merge accepts a
// channel of channels as input. For each of the channel read from the channel
// of channels, all values are consumed and produced onto the output channel.
// Merge is blocking and panics when the input channel types do no match the
// output channel.
func Merge(done <-chan struct{}, in interface{}, out interface{}) {

	// Ensure that all arguments are compatible types
	if reflect.TypeOf(out).Kind() != reflect.Chan {
		panic(fmt.Sprintf("cannot merge into type %v", reflect.TypeOf(out)))
	}
	if reflect.TypeOf(in).Kind() != reflect.Chan {
		panic(fmt.Sprintf("cannot merge from type %v", reflect.TypeOf(in)))
	}
	if reflect.TypeOf(in).Elem().Kind() != reflect.Chan {
		panic(fmt.Sprintf("cannot merge from type %v", reflect.TypeOf(in).Elem()))
	}
	if reflect.TypeOf(in).Elem().Elem().Kind() != reflect.TypeOf(out).Elem().Kind() {
		panic(fmt.Sprintf("cannot merge from type %T with elements of type", reflect.TypeOf(in).Elem().Elem()))
	}

	var wg sync.WaitGroup
	for {

		// select {
		// case <-done:
		// case ch, ok := <-in:
		// }
		chosen, ch, ok := reflect.Select([]reflect.SelectCase{
			reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(done)},
			reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(in)},
		})
		if chosen == 0 || !ok {
			break
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				// select {
				// case <-done:
				// case val, ok := <-ch:
				// }
				chosen, val, ok := reflect.Select([]reflect.SelectCase{
					reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(done)},
					reflect.SelectCase{Dir: reflect.SelectRecv, Chan: ch},
				})
				if chosen == 0 || !ok {
					return
				}

				// select {
				// case <-done:
				// case out <- val:
				// }
				chosen, _, _ = reflect.Select([]reflect.SelectCase{
					reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(done)},
					reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(out), Send: val},
				})
				if chosen == 0 {
					return
				}
			}
		}()
	}

	wg.Wait()
}
