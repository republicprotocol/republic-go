# go-eth

[![Build Status](https://travis-ci.org/republicprotocol/go-eth.svg?branch=master)](https://travis-ci.org/republicprotocol/go-eth)
[![Coverage Status](https://coveralls.io/repos/github/republicprotocol/go-eth/badge.svg?branch=master)](https://coveralls.io/github/republicprotocol/go-eth?branch=master)

> This library is a work in progress.


## Regenerating contract bindings

#### Clone the contracts git submodule

```sh
git submodule update --init --recursive
```

#### Install godep

```sh
go get github.com/tools/godep
```

#### Use godep to install abigen
(TODO: Figure out how to install all dependencies in one go - `go get` doesn't)

```sh
go get github.com/ethereum/go-ethereum
cd $GOPATH/src/github.com/ethereum/go-ethereum
go get github.com/btcsuite/btcd/btcec
go get github.com/go-stack/stack
gi get github.com/golang/snappy
go get github.com/rcrowley/go-metrics
go get github.com/syndtr/goleveldb/leveldb
go get github.com/ethereum/go-ethereum/logger
go get gopkg.in/karalabe/cookiejar.v2/collections/prque
godep save
godep go install ./cmd/abigen
```

#### Install solc

Ubuntu:

```sh
sudo add-apt-repository ppa:ethereum/ethereum
sudo apt-get update
sudo apt-get install solc
```

MacOS:

```sh
brew update
brew upgrade
brew tap ethereum/ethereum
brew install solidity
brew linkapps solidity
```

#### Generate bindings

```sh
cd contracts
./generate.sh
cd ../
```

## Testing

```sh
ginkgo -v --trace --cover --coverprofile coverprofile.out
```

## License

The go-eth library was developed by the Republic Protocol team and is available under the MIT license. For more information, see our website https://republicprotocol.com.
