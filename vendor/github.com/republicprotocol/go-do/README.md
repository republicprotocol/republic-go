# Do Concurrently

[![Build Status](https://travis-ci.org/republicprotocol/go-do.svg?branch=master)](https://travis-ci.org/republicprotocol/go-do)
[![Coverage Status](https://coveralls.io/repos/github/republicprotocol/go-do/badge.svg?branch=master)](https://coveralls.io/github/republicprotocol/go-do?branch=master)

The Do Concurrently library is a Go implementation of high level concurrent features. It provides a simple API for common task parallel and data parallel constructs. Using goroutines, the parallelism provided is actually a form of concurrency since goroutines are not guaranteed to run strict simultaneity.

## ForAll and CoForAll

The `ForAll` and `CoForAll` loops are data parallel loops that distributes iterations evenly across several goroutines. The `ForAll` loop will use one goroutine per CPU and the `CoForAll` loop will use one goroutine per data item. Both loops work over arrays and slices. Both loops will block until all iterations have terminated.

```go
xs := []int{1,2,3,4,5,6,7,8,10}
do.ForAll(xs, func(i int) {
    xs[i] *= 2
})
do.CoForAll(xs, func(i int) {
    xs[i] *= 3
})
log.Println("6 =", xs[0])
```

It is the responsibility of the programmer to ensure that the function being used is safe for concurrent environments. A simple way of ensuring this is checking that the function will never mutate any object other than the object accessible using the `i` index. You can also use the go tools to check for race conditions during testing.

## Process

A `Process` accepts a set of process functions that are executed concurrently and all return values are written to a channel. The of order of the return values being written to the channel is the same as the order of the set of functions. Using channels to handle return values is simpler and safer than trying to synchronize using share memory. The official Go documentation recommends the use of channels over shared memory.

```go
ret := do.Process(func() do.Option {
    return do.Ok(1)
}, func() do.Option {
    return do.Ok(2)
})
log.Println("1 =", <- ret)
log.Println("2 =", <- ret)
```

### Begin and CoBegin

The `Begin` and `CoBegin` functions accept a set of functions that are executed concurrently. Both functions will store the return values in an output list, where the order of outputs matches the order of the functions. Both functions will block until all functions have executed in full.

```go
ret := do.Begin(func() do.Option {
    return do.Ok(1)
}, func() do.Option {
    return do.Ok(2)
})
ret = do.CoBegin(func() do.Option {
    return do.Ok(ret[0].Ok + ret[1].Ok*2)
}, func() do.Option {
    return do.Ok(ret[0].Ok + ret[1].Ok*3)
})
log.Println("5 =", ret[0])
log.Println("7 =", ret[1])
```

`Begin` function will use one goroutine per CPU and distribute the functions evenly over each goroutine. `CoBegin` will always use one goroutine per function.

### Guarded Objects

Using the `GuardedObject` you can provide safe concurrent access to an object. Using the `Enter` and `Exit` functions you can guarantee that the object is being accessed mutually exclusively. In addition to this, you can define `Guards` and pass them to `Enter` to make sure that the object is only accessed if a certain boolean condition is true.

Below is an example of a concurrent `Stack`. If it is full, a call to `Push` function will not overflow. Instead, it will wait until the `Stack` is not full and then continue. Similarly, the `Pop` function will not underflow.

```go
type Stack struct {
    do.GuardedObject

    i             int
    items         []interface{}
    itemsNotEmpty *do.Guard
    itemsNotFull  *do.Guard
}

func NewStack() *Stack {
    s := new(Stack)
    s.GuardedObject = do.NewGuardedObject()
    s.i = 0
    s.items = make([]interface{}, 10)
    s.itemsNotEmpty = q.Guard(func() bool { return len(q.items) > 0 })
    s.itemsNotFull = q.Guard(func() bool { return len(q.items) < 10 })
    return q
}

func (s *Stack) Push(item interface{}) {
    s.Enter(s.isListNotFull)
    defer s.Exit()
    s.items[i] = item
    s.i++
}

func (s *Stack) Pop() (item interface{}) {
    s.Enter(s.isListNotEmpty)
    defer s.Exit()
    item = s.items[i]
    s.i--
}
```

## Tests

To run the test suite, install Ginkgo.

```
go get github.com/onsi/ginkgo/ginkgo
```

Now we can run the tests.

```
ginkgo -v --race --trace --cover --coverprofile coverprofile.out
```

## License

The Do Concurrently library was developed by the Republic Protocol team and is available under the MIT license. For more information, see our website https://republicprotocol.com.