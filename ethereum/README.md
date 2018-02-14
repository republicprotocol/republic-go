## Testing

In the `ethereum` directory:

```sh
ginkgo -v --trace --cover --coverprofile coverprofile.out
```

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