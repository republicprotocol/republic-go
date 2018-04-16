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
			SendToInterface(chOut, msg)
		}
	}
}

func SendToInterface(chOut interface{}, msg interface{}) {
	msgValue := reflect.ValueOf(msg)
	switch reflect.TypeOf(chOut).Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < reflect.ValueOf(chOut).Len(); i++ {
			if reflect.ValueOf(chOut).Index(i).Kind() != reflect.Chan {
				panic(fmt.Sprintf("cannot split to value of type %T", chOut))
			}
			reflect.ValueOf(chOut).Index(i).Send(msgValue)
		}
	case reflect.Chan:
		reflect.ValueOf(chOut).Send(msgValue)
	default:
		panic(fmt.Sprintf("cannot split to value of type %T", chOut))
	}
}

type Splitter struct {
	mu          *sync.RWMutex
	subscribers map[interface{}]struct{}

	maxConnections int
}

func NewSplitter(maxConnections int) Splitter {
	return Splitter{
		mu:          &sync.RWMutex{},
		subscribers: make(map[interface{}]struct{}),

		maxConnections: maxConnections,
	}
}

func (splitter *Splitter) Subscribe(ch interface{}) error {
	splitter.mu.Lock()
	defer splitter.mu.Unlock()

	if len(splitter.subscribers) >= splitter.maxConnections {
		return fmt.Errorf("cannot run message queue: max connections reached")
	}

	splitter.subscribers[ch] = struct{}{}
	return nil
}

func (splitter *Splitter) Unsubscribe(ch interface{}) {
	splitter.mu.Lock()
	defer splitter.mu.Unlock()
	delete(splitter.subscribers, ch)
}

func (splitter *Splitter) Split(chIn interface{}) {
	if reflect.TypeOf(chIn).Kind() != reflect.Chan {
		panic(fmt.Sprintf("cannot split from value of type %T", chIn))
	}

	for {
		msg, ok := reflect.ValueOf(chIn).Recv()
		if !ok {
			return
		}

		func() {
			splitter.mu.RLock()
			defer splitter.mu.RUnlock()
			for chOut := range splitter.subscribers {
				SendToInterface(chOut, msg)
			}
		}()
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
