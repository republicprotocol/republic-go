# Swarm Network

[![Build Status](https://travis-ci.org/republicprotocol/go-swarm-network.svg?branch=master)](https://travis-ci.org/republicprotocol/go-swarm-network)
[![Coverage Status](https://coveralls.io/repos/github/republicprotocol/go-swarm-network/badge.svg?branch=master)](https://coveralls.io/github/republicprotocol/go-swarm-network?branch=master)

> This library is a work in progress.

The Swarm Network library is an official reference implementation of the Swarm Network in the Republic Protocol, written in Go. It is a P2P network that supports pinging nodes, searching for nodes, and bootstrapping.

## Installation

There are several components that need to be installed before you can run tests.

### Install Proto3

Protobuf is a binary communication protocol developed by Google. The Swarm Network uses it to perform remote procedure calls between nodes. To install it, run the commands given below to install the required version of Protobuf.

```
curl -OL https://github.com/google/protobuf/releases/download/v3.2.0/protoc-3.2.0-linux-x86_64.zip
unzip protoc-3.2.0-linux-x86_64.zip -d protoc3

sudo mv protoc3/bin/* /usr/local/bin/
sudo mv protoc3/include/* /usr/local/include/

sudo chown $USER /usr/local/bin/protoc
sudo chown -R $USER /usr/local/include/google
```

The Swarm Network is written in Go, so you will also need to install the Go plugin. If you are building an implementation in another language, you will need to install the Protobuf plugin for that language.

```
go get -u github.com/golang/protobuf/protoc-gen-go
```

### Install gRPC

gRPC is a remote procedure calling library developed by Google, built on top of Protobuf. Make sure you have followed the instructions for installing Protobuf, and then run the command below.

```
go get -u google.golang.org/grpc
```

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

The Swarm Network library was developed by the Republic Protocol team and is available under the MIT license. For more information, see our website https://republicprotocol.com.