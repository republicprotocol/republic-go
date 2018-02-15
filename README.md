# Dark Network

[![Build Status](https://travis-ci.org/republicprotocol/go-dark-network.svg?branch=master)](https://travis-ci.org/republicprotocol/go-xing)
[![Coverage Status](https://coveralls.io/repos/github/republicprotocol/go-dark-network/badge.svg?branch=master)](https://coveralls.io/github/republicprotocol/go-xing?branch=master)

> This library is a work in progress.

The Xing library is an official reference implementation of the Xing overlay network for the Republic Protocol, written in Go. It supports the calculation of the number of classes in the Xing Network, and the number of M Networks in the Xing Network. It also supports the assignment of classes and M Networks, based on an Epoch hash.

More details on the inner workings of the Dark Network will be made available on the Republic Protocol Wiki in the future.

## Assigning the Dark Network

Each miner is assigned to a class, and M Network, within the Dark Network. The class of a miner determines which order fragments it is authorized to receive, and the M Networks define partitioned groups of miners that work together to match orders.

The assignments are based on the current Epoch hash, the miner's commitment hash, and the total number of miners in the current Epoch. All of this data must be loaded from the Ethereum smart contract for the Republic Protocol. Once it has been loaded, the `AssignXOverlay` function can be used to assign a class and M Network to each miner.

```go
epochHash := []byte{} // Load from Ethereum.
miners := []x.Miner{} // Load from Ethereum.
numberOfMNetworks := x.NumberOfMNetworks(len(miners))
x.AssignXOverlay(miners, epochHash, numberOfMNetworks)
```

Every time the Epoch hash is changed, the Dark Network must be reassigned.

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