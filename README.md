# Dark Network

[![Build Status](https://travis-ci.org/republicprotocol/go-dark-network.svg?branch=master)](https://travis-ci.org/republicprotocol/go-dark-network)
[![Coverage Status](https://coveralls.io/repos/github/republicprotocol/go-dark-network/badge.svg?branch=master)](https://coveralls.io/github/republicprotocol/go-dark-network?branch=master)

> This library is a work in progress.

The Dark Network library is an official reference implementation of the Dark Network for the Republic Protocol, written in Go. It provides the required gRPC interfaces for the secure multi-party computations that compare orders.

More details on the inner workings of the Dark Network will be made available on the Republic Protocol Wiki in the future.

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

The Dark library was developed by the Republic Protocol team and is available under the MIT license. For more information, see our website https://republicprotocol.com.
