# Dark Ocean

> This library is a work in progress.

The Dark Ocean library is an official reference implementation of the Dark Ocean for the Republic Protocol, written in Go. It provides the required gRPC interfaces for the secure multi-party computations that compare orders.

More details on the inner workings of the Dark Ocean will be made available on the Republic Protocol Wiki in the future.

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

The Dark Ocean library was developed by the Republic Protocol team and is available under the MIT license. For more information, see our website https://republicprotocol.com.