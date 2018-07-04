package dispatch

import (
	"fmt"
	"reflect"
)

func Merge(done <-chan struct{}, out interface{}, in interface{}) {
	if reflect.TypeOf(out).Kind() != reflect.Chan {
		panic(fmt.Sprintf("cannot merge into type %T", out))
	}
	elemKind := reflect.TypeOf(out).Elem().Kind()

	if reflect.TypeOf(in).Kind() != reflect.Chan {
		panic(fmt.Sprintf("cannot merge from type %T", out))
	}
	if reflect.TypeOf(in).Elem().Kind() != reflect.Chan {
		panic(fmt.Sprintf("cannot merge from type %T", out))
	}
	if reflect.TypeOf(in).Elem().Elem().Kind() != elemKind {
		panic(fmt.Sprintf("cannot merge from type %T", out))
	}

	for {
		chosen, ch, ok := reflect.Select([]reflect.SelectCase{
			reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(done)},
			reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(in)},
		})
		if chosen == 0 || !ok {
			return
		}

		go func() {
			for {
				chosen, val, ok := reflect.Select([]reflect.SelectCase{
					reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(done)},
					reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)},
				})
				if chosen == 0 || !ok {
					return
				}

				chosen, val, ok = reflect.Select([]reflect.SelectCase{
					reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(done)},
					reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(out), Send: val},
				})
				if chosen == 0 || !ok {
					return
				}
			}
		}()
	}
}
