package dispatch

import (
	"fmt"
	"reflect"
	"sync"
)

// CoForAll uses an iterator to execute a function on each element in a
// goroutine. An iterator can be an integer, an array, or a slice. This
// function blocks until all goroutines have terminated.
func CoForAll(iter interface{}, f func(int)) {

	numGoroutines := 0
	switch reflect.TypeOf(iter).Kind() {
	case reflect.Array, reflect.Slice:
		numGoroutines = reflect.ValueOf(iter).Len()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		numGoroutines = int(reflect.ValueOf(iter).Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		numGoroutines = int(reflect.ValueOf(iter).Uint())

	default:
		panic(fmt.Sprintf("coforall error: expected iterator got %T", iter))
	}

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			f(i)
		}(i)
	}
	wg.Wait()
}
