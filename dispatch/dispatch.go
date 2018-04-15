package dispatch

import (
	"fmt"
	"reflect"
	"sync"
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

func Split(chIn interface{}, chsOut ...interface{}) {
	if reflect.TypeOf(chIn).Kind() != reflect.Chan {
		panic(fmt.Sprintf("cannot split from value of type %T", chIn))
	}
	for {
		msg, ok := reflect.ValueOf(chIn).Recv()
		if !ok {
			return
		}
		for _, chOut := range chsOut {
			switch reflect.TypeOf(chOut).Kind() {
			case reflect.Array, reflect.Slice:
				for i := 0; i < reflect.ValueOf(chOut).Len(); i++ {
					if reflect.ValueOf(chOut).Index(i).Kind() != reflect.Chan {
						panic(fmt.Sprintf("cannot split to value of type %T", chOut))
					}
					reflect.ValueOf(chOut).Index(i).Send(msg)
				}
			case reflect.Chan:
				reflect.ValueOf(chOut).Send(msg)
			default:
				panic(fmt.Sprintf("cannot split to value of type %T", chOut))
			}
		}
	}
}

func Merge(chOut interface{}, chsIn ...interface{}) {
	if reflect.TypeOf(chOut).Kind() != reflect.Chan {
		panic(fmt.Sprintf("cannot merge to value of type %T", chOut))
	}

	var wg sync.WaitGroup

	mergeCh := func(chIn interface{}) {
		defer wg.Done()
		for {
			msg, ok := reflect.ValueOf(chIn).Recv()
			if !ok {
				return
			}
			reflect.ValueOf(chOut).Send(msg)
		}
	}

	for _, chIn := range chsIn {
		switch reflect.TypeOf(chIn).Kind() {
		case reflect.Array, reflect.Slice:
			for i := 0; i < reflect.ValueOf(chIn).Len(); i++ {
				if reflect.ValueOf(chIn).Index(i).Kind() != reflect.Chan {
					panic(fmt.Sprintf("cannot merge from value of type %T", chIn))
				}
				wg.Add(1)
				go mergeCh(reflect.ValueOf(chIn).Index(i))
			}
		case reflect.Chan:
			wg.Add(1)
			go mergeCh(chIn)
		default:
			panic(fmt.Sprintf("cannot merge from value of type %T", chOut))
		}
	}

	wg.Wait()
}
