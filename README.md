# Republic Protocol

[![Build Status](https://travis-ci.org/republicprotocol/republic-go.svg?branch=master)](https://travis-ci.org/republicprotocol/republic-go)
[![Coverage Status](https://coveralls.io/repos/github/republicprotocol/republic-go/badge.svg?branch=master)](https://coveralls.io/github/republicprotocol/republic-go?branch=master)

An official reference implementation of Republic Protocol, written in Go.

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