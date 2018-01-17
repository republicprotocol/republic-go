# Order Compute Engine

[![Build Status](https://travis-ci.org/republicprotocol/go-order-compute.svg?branch=master)](https://travis-ci.org/republicprotocol/go-order-compute)
[![Coverage Status](https://coveralls.io/repos/github/republicprotocol/go-order-compute/badge.svg?branch=master)](https://coveralls.io/github/republicprotocol/go-order-compute?branch=master)

The Order Compute Engine library is an official reference implementation of the computation engine for orders and order fragments in the Republic Protocol, written in Go. This library supports the construction of orders and order fragments, as well as the computations required for order fragments.

## Tests

To run the test suite, install Ginkgo.

```
go get github.com/onsi/ginkgo/ginkgo
```

Now we can run the tests.

```
ginkgo -v --trace --cover --coverprofile coverprofile.out
```

## License

The Order Compute Engine library was developed by the Republic Protocol team and is available under the MIT license. For more information, see our website https://republicprotocol.com.