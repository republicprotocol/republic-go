package do

import (
	"reflect"
	"runtime"
	"sync"
)

// Option types are returned from Process functions. They are either a value or
// an error. The error should be checked before using the value.
type Option struct {
	Ok  interface{}
	Err error
}

// Ok returns an Option struct with a value and no error.
func Ok(ok interface{}) Option {
	return Option{
		Ok: ok,
	}
}

// Err returns an Option struct with an error and no value.
func Err(err error) Option {
	return Option{
		Err: err,
	}
}

// Then runs the function on the Option value, if there is no error in the
// Option. In all other cases it returns the Option unmodified.
func (option Option) Then(f func(ok interface{}) Option) Option {
	if option.Err != nil {
		return option
	}
	return f(option.Ok)
}

// Catch runs the function on the Option error, if there is an error in the
// Option. In all other cases it returns the Option unmodified.
func (option Option) Catch(f func(err error) Option) Option {
	if option.Err != nil {
		return f(option.Err)
	}
	return option
}

// ForAll items in the data set, apply the function. The function accepts the
// index of the item to which is it being applied. One goroutine is launched
// for each CPU, so the given function must be safe to use concurrently.
func ForAll(data interface{}, f func(i int)) {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Array, reflect.Slice:
		// Calculate workload size per CPU.
		length := reflect.ValueOf(data).Len()
		numCPUs := runtime.NumCPU()
		numIterationsPerCPU := (length / numCPUs) + 1
		// Apply the function in parallel over the data.
		var wg sync.WaitGroup
		for offset := 0; offset < length; offset += numIterationsPerCPU {
			wg.Add(1)
			go func(offset int) {
				defer wg.Done()
				for i := offset; i < offset+numIterationsPerCPU && i < length; i++ {
					f(i)
				}
			}(offset)
		}
		wg.Wait()
	}
}

// CoForAll is the same as ForAll expect it uses one goroutine per item in the
// data set.
func CoForAll(data interface{}, f func(i int)) {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Array, reflect.Slice:
		// Apply the function in parallel over the data.
		length := reflect.ValueOf(data).Len()
		var wg sync.WaitGroup
		wg.Add(length)
		for i := 0; i < length; i++ {
			go func(i int) {
				defer wg.Done()
				f(i)
			}(i)
		}
		wg.Wait()
	}
}

// Process runs each function in a goroutine and writes the return values to a
// channel. The values will be written to the channel in the same order that
// the process functions are provided.
func Process(fs ...func() Option) chan Option {
	guards := make([]*sync.Mutex, len(fs))
	for i := range fs {
		guards[i] = new(sync.Mutex)
		guards[i].Lock()
	}
	ch := make(chan Option, len(fs))

	// Create a wait group for all processes.
	var wg sync.WaitGroup
	wg.Add(len(fs))
	for i, f := range fs {
		go func(i int, f func() Option) {
			defer wg.Done()
			defer guards[i].Unlock()
			if i > 0 {
				guards[i-1].Lock()
			}
			ch <- f()
		}(i, f)
	}

	// Run a goroutine that will close the channel when all processes have
	// terminated.
	go func() {
		defer close(ch)
		wg.Wait()
	}()

	return ch
}

// Begin runs each function and blocks until all functions have terminated. It
// distributes the functions evenly over a number of goroutines equal to the
// number of CPUs.
func Begin(fs ...func() Option) []Option {
	// Calculate workload size per CPU.
	length := len(fs)
	numCPUs := runtime.NumCPU()
	numIterationsPerCPU := (length / numCPUs) + 1
	// Apply the functions in parallel.
	output := make([]Option, length)
	var wg sync.WaitGroup
	for offset := 0; offset < length; offset += numIterationsPerCPU {
		wg.Add(1)
		go func(offset int) {
			defer wg.Done()
			for i := offset; i < offset+numIterationsPerCPU && i < length; i++ {
				output[i] = fs[i]()
			}
		}(offset)
	}
	wg.Wait()
	return output
}

// CoBegin is the same as Begin expect it uses one goroutine per function.
func CoBegin(fs ...func() Option) []Option {
	output := make([]Option, len(fs))
	var wg sync.WaitGroup
	wg.Add(len(fs))
	for i, f := range fs {
		go func(i int, f func() Option) {
			defer wg.Done()
			output[i] = f()
		}(i, f)
	}
	wg.Wait()
	return output
}
