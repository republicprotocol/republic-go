# X

[![Build Status](https://travis-ci.org/republicprotocol/go-x.svg?branch=master)](https://travis-ci.org/republicprotocol/go-x)
[![Coverage Status](https://coveralls.io/repos/github/republicprotocol/go-x/badge.svg?branch=master)](https://coveralls.io/github/republicprotocol/go-x?branch=master)

> This library is a work in progress.

The X library is an official reference implementation of the X Network overlay for the Republic Protocol, written in Go. It supports the calculation of the number of classes in the X Network, and the number of M Networks in the X Network. It also supports the assignment of classes and M Networks, based on an Epoch hash.

More details on the inner workings of the X Network will be made available on the Republic Protocol Wiki in the future.

## Assigning the X Overlay

Each miner is assigned to a class, and M Network, within the X Network. This assignment is based on the current Epoch hash, the miner's commitment hash, and the total number of miners in the current Epoch. All of this data must be loaded from the Ethereum smart contract for the Republic Protocol. Once it has been loaded, the `AssignXOverlay` function can be used to assign a class and M Network to each miner.

```go
epochHash := []byte{} // Load from Ethereum.
miners := []x.Miner{} // Load from Ethereum.
numberOfMNetworks := x.NumberOfMNetworks(len(miners))
x.AssignXOverlay(miners, epochHash, numberOfMNetworks)
```

Every time the Epoch hash is changed, the X Overlay must be reassigned.

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