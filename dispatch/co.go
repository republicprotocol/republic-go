package dispatch

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
)

// CoBegin multiple functions onto background goroutines. This function blocks
// until all goroutines have terminated.
func CoBegin(fs ...func()) {
	var wg sync.WaitGroup
	for _, f := range fs {
		wg.Add(1)
		go func(f func()) {
			defer wg.Done()
			f()
		}(f)
	}
	wg.Wait()
}

// CoForAll uses an iterator to execute an iterand function on each value
// returned by the iterator, using a background goroutine for each iteration.
// An iterator can be an array, a slice, a map, or an int. For arrays, and
// slices, the iterand function must accept an int index as the only argument.
// For maps, the iterand function must accept a key as the only argument. For
// ints, the iterand function must accept an int, in the range [0, n), as the
// only argument. This function blocks until all goroutines have terminated.
func CoForAll(iter interface{}, f interface{}) {
	funTy := reflect.TypeOf(f)
	if funTy.Kind() != reflect.Func {
		panic(fmt.Sprintf("coforall error: expected iterator got %T", iter))
	}
	fun := reflect.ValueOf(f)

	switch reflect.TypeOf(iter).Kind() {

	case reflect.Array, reflect.Slice:
		it := reflect.ValueOf(iter)
		numGoroutines := it.Len()

		var wg sync.WaitGroup
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(i int) {
				defer wg.Done()
				fun.Call([]reflect.Value{reflect.ValueOf(i)})
			}(i)
		}
		wg.Wait()

	case reflect.Map:
		it := reflect.ValueOf(iter)
		keys := it.MapKeys()
		numGoroutines := len(keys)

		var wg sync.WaitGroup
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(i int) {
				defer wg.Done()
				fun.Call([]reflect.Value{keys[i]})
			}(i)
		}
		wg.Wait()

	case reflect.Int:
		numGoroutines := int(reflect.ValueOf(iter).Int())

		var wg sync.WaitGroup
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(i int) {
				defer wg.Done()
				fun.Call([]reflect.Value{reflect.ValueOf(i)})
			}(i)
		}
		wg.Wait()

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		panic(fmt.Sprintf("coforall error: expected int got %T", iter))

	default:
		panic(fmt.Sprintf("coforall error: expected iterator got %T", iter))
	}

}

// ForAll uses an iterator to execute an iterand function on each value
// returned by the iterator, using a background goroutine for CPU and
// distributing the iterands evenly across each goroutine. An iterator can be
// an array, a slice, a map, or an int. For arrays, and slices, the iterand
// function must accept an int index as the only argument. For maps, the
// iterand function must accept a key as the only argument. For ints, the
// iterand function must accept an int, in the range [0, n), as the only
// argument. This function blocks until all goroutines have terminated.
func ForAll(iter interface{}, f interface{}) {
	funTy := reflect.TypeOf(f)
	if funTy.Kind() != reflect.Func {
		panic(fmt.Sprintf("coforall error: expected iterator got %T", iter))
	}
	fun := reflect.ValueOf(f)

	switch reflect.TypeOf(iter).Kind() {

	case reflect.Array, reflect.Slice:
		it := reflect.ValueOf(iter)
		num := it.Len()
		numGoroutines := runtime.NumCPU()
		numPerGoroutine := num/numGoroutines + 1

		var wg sync.WaitGroup
		for i := 0; i < num; i += numPerGoroutine {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				for j := i; j < i+numPerGoroutine && j < num; j++ {
					fun.Call([]reflect.Value{reflect.ValueOf(j)})
				}
			}(i)
		}
		wg.Wait()

	case reflect.Map:
		it := reflect.ValueOf(iter)
		keys := it.MapKeys()
		num := len(keys)
		numGoroutines := runtime.NumCPU()
		numPerGoroutine := num/numGoroutines + 1

		var wg sync.WaitGroup
		for i := 0; i < num; i += numPerGoroutine {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				for j := i; j < i+numPerGoroutine && j < num; j++ {
					fun.Call([]reflect.Value{keys[j]})
				}
			}(i)
		}
		wg.Wait()

	case reflect.Int:
		num := int(reflect.ValueOf(iter).Int())
		numGoroutines := runtime.NumCPU()
		numPerGoroutine := num/numGoroutines + 1

		var wg sync.WaitGroup
		for i := 0; i < num; i += numPerGoroutine {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				for j := i; j < i+numPerGoroutine && j < num; j++ {
					fun.Call([]reflect.Value{reflect.ValueOf(j)})
				}
			}(i)
		}
		wg.Wait()

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		panic(fmt.Sprintf("forall error: expected int got %T", iter))

	default:
		panic(fmt.Sprintf("forall error: expected iterator got %T", iter))
	}

}
