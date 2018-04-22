# Republic Protocol

[![Build Status](https://travis-ci.org/republicprotocol/republic-go.svg?branch=master)](https://travis-ci.org/republicprotocol/republic-go)
[![Coverage Status](https://coveralls.io/repos/github/republicprotocol/republic-go/badge.svg?branch=master)](https://coveralls.io/github/republicprotocol/republic-go?branch=master)

An official reference implementation of Republic Protocol, written in Go.

## Developing

### Protobuf

#### macOS

```sh
brew install protobuf
```

#### Linux

```sh
PROTOC_ZIP=protoc-3.3.0-linux-x86_64.zip
curl -OL https://github.com/google/protobuf/releases/download/v3.3.0/$PROTOC_ZIP
sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
rm -f $PROTOC_ZIP
```

### gRPC

```sh
go get -u google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go
```

## Testing

### Ganache

Install the Ganache CLI.

```sh
npm i -g ganache-cli
```

Before running tests that use the Ethereum network, run Ganache.

```sh
go run ./cmd/ganache/ganache.go
```