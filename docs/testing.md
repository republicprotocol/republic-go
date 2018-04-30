# Testing

All packages use the [Ginkgo](https://github.com/onsi/ginkgo) testing framework to encourage behavior testing.

To run tests, you will need to install Ginkgo.

```sh
go get github.com/onsi/ginkgo/ginkgo
```

Once Ginkgo is installed, you can run tests of a package by typing ```ginkgo``` within the package directory's terminal.

To write tests, you will need to install Gomega.

```sh
go get github.com/onsi/gomega/...
```

To run a local testnet for testing, you will need to install [Ganache-cli](https://github.com/trufflesuite/ganache-cli/blob/master/README.md).

```sh
npm install -g ganache-cli
```

## Darknode Local Testnet

The [testnet.go](https://github.com/republicprotocol/republic-go/blob/darknodetest/darknode/testnet.go) file inside the darknode package provides the _NewTestnet_ function which will set up a local testnet with a user specified number of dark nodes and bootstrap nodes.

The following is an example of how you can use it for your tests.

```go
var err error
var env darknode.TestnetEnv

BeforeEach(func() {
    env, err = darknode.NewTestnet(numberOfDarknodes, numberOfBootstrapNodes)
    Expect(err).ShouldNot(HaveOccurred())
})

AfterEach(func() {
    env.Teardown()
})
```

NewTestnet will create dark nodes, register them to a darknode registry and connect to a local ganache server. It will then turn on their RPCs and deploy their smart contracts.

To use the local tesnet, you must first import the _darknode_ package in your test file and pass the following arguments to NewTestnet function:

* Number of darknodes that your test will require in the local testnet
* Number of bootstrap nodes in the local testnet

NewTestnet currently **does not** support a Bitcoin connection and cannot control smart contract variables.

The testnet prevents multiple tests from runnning simultaneously because it needs to control local ports. It uses ganache snapshots to keep a clean local ethereum testnet between tests.
_See the darknode/testnet\_test.go for an example_