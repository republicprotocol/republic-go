# Ethereum

## Regenerating contract bindings

#### Clone the contracts git submodule

```sh
git submodule update --init --recursive
```

#### Install go dependencies

```sh
go get
```

#### Install abigen

```sh
cd $GOPATH/src/github.com/ethereum/go-ethereum
go install ./cmd/abigen
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
