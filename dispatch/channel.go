package dispatch

import (
	"fmt"
	"reflect"
)

func Wait(chs ...chan struct{}) {
	for _, ch := range chs {
		for range ch {
		}
	}
}

func Close(chs ...interface{}) {
	for _, ch := range chs {
		if reflect.TypeOf(ch).Kind() == reflect.Chan {
			reflect.ValueOf(ch).Close()
		}
	}
}

func Split(ch interface{}, chs ...interface{}) {
	if reflect.TypeOf(ch).Kind() != reflect.Chan {
		return
	}
	for {
		msg, ok := reflect.ValueOf(ch).Recv()
		if !ok {
			return
		}
		for _, c := range chs {
			switch reflect.TypeOf(c).Kind() {
			case reflect.Array, reflect.Slice:
				for i := 0; i < reflect.ValueOf(c).Len(); i++ {
					if reflect.ValueOf(c).Index(i).Kind() == reflect.Chan {
						reflect.ValueOf(c).Index(i).Send(msg)
					}
				}
			case reflect.Chan:
				reflect.ValueOf(c).Send(msg)
			default:
				panic(fmt.Sprintf("cannot split value of type %T", ch))
			}
		}
	}
}

// func Merge(ch chan interface{}, chs ...chan interface{}) {
// 	for _, channel := range chs {
// 		go func(channel chan interface{}) {
// 			for msg := range channel {
// 				ch <- msg
// 			}
// 		}(channel)
// 	}
// }

func Filter(chIn interface{}, chOut interface{}, predicate func(p interface{}) bool) {
	if reflect.TypeOf(chIn).Kind() != reflect.Chan || reflect.TypeOf(chOut).Kind() != reflect.Chan {
		return
	}
	for {
		msg, ok := reflect.ValueOf(chIn).Recv()
		if !ok {
			return
		}
		if !predicate(msg) {
			continue
		}
		reflect.ValueOf(chOut).Send(msg)
	}
}
