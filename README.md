# X

[![Build Status](https://travis-ci.org/republicprotocol/go-x.svg?branch=master)](https://travis-ci.org/republicprotocol/go-x)
[![Coverage Status](https://coveralls.io/repos/github/republicprotocol/go-x/badge.svg?branch=master)](https://coveralls.io/github/republicprotocol/go-x?branch=master)

> This library is a work in progress.

The X library is an official reference implementation of X Network overlay functions in the Republic Protocol, written in Go. It supports the calculation of the number of classes in the X Network, and the number of M Networks in the X Network. It supports the sorting of miners into classes and M Networks based on an Epoch hash, and commitment hashes.

More details on the inner workings of the X Network will be made available on the Republic Protocol Wiki in the future.

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

The X library was developed by the Republic Protocol team and is available under the MIT license. For more information, see our website https://republicprotocol.com.