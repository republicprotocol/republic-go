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