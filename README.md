# Dark Node

[![Build Status](https://travis-ci.org/republicprotocol/go-dark-node.svg?branch=master)](https://travis-ci.org/republicprotocol/go-dark-node)
[![Coverage Status](https://coveralls.io/repos/github/republicprotocol/go-dark-node/badge.svg?branch=master)](https://coveralls.io/github/republicprotocol/go-dark-node?branch=master)

The Dark Node library is an official reference implementation of Dark Nodes for the Republic Protocol, written in Go. It supports the bootstrapping of nodes into the network, and into their assigned Dark Pools. It provides a delegate implementation for the Swarm and Dark networks. Custom binary command-line tools are available for testing, and booting a Dark Node.

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

The Dark Node library was developed by the Republic Protocol team and is available under the MIT license. For more information, see our website https://republicprotocol.com.
