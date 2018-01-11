# Identity

[![Build Status](https://travis-ci.org/republicprotocol/go-identity.svg?branch=master)](https://travis-ci.org/republicprotocol/go-identity)
[![Coverage Status](https://coveralls.io/repos/github/republicprotocol/go-identity/badge.svg?branch=master)](https://coveralls.io/github/republicprotocol/go-identity?branch=master)

The Identity library is an official reference implementation of identities and addresses in the Republic Protocol, written in Go. It supports the generation and use of Republic key pairs, Republic IDs, Republic addresses, and Republic multi-addresses.

## Key Pairs

To create an identity, or an address, in the Republic Protocol an ECDSA key pair is required. This key pair is the same as those used by Bitcoin and Ethereum.

```go
keyPair, err := identity.NewKeyPair()
```

## ID

Republic IDs are used to identity miners and traders in the Republic Protocol, and are required when registering with the Registrar. To generate an ID, take the first 20 bytes of the Keccak-256 hash of your public key. 

```go
keyPair, err := identity.NewKeyPair()
id := keyPair.ID()
fmt.Println(id)
```

An example output of this code snippet is given below.

```
[247 107 192 130 208 33 69 182 183 127 237 24 22 83 157 99 207 221 9 243]
```

## Address

Republic addresses are URL compatible encoding of an ID. To generate an address, take the Base58 encoding of the multi-hash of an ID. The Bitcoin alphabet for Base58 encoding must be used. For more information about multi-hashes, see https://multiformats.io/multihash.


```go
keyPair, err := identity.NewKeyPair()
address := keyPair.Address()
fmt.Println(address)
```

An example output of this code snippet is given below.

```
8MGHyHWxQ5zUzXXnK5b6LkWGFCWSyX
```

## Multi-address

Republic multi-addresses are a form of network address that can represent multiple different networking layers in a single address. In the Republic Protocol, they are used to hold network addresses and Republic addresses. At the moment, the network address is excepted to be an IPv4 address and a TCP port. In future, it may change to be an I2P address that provides a greater level of anonymity. For more information about multi-addresses, see https://multiformats.io/multiaddr.

```go
keyPair, err := identity.NewKeyPair()
multi, err := keyPair.MultiAddress()
multi, err = identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/80/republic/%s", multi.String())
fmt.Println(multi.String())
fmt.Println(multi.ValueForProtocol(identity.RepublicCode))
fmt.Println(multi.ValueForProtocol(identity.IP4Code))
fmt.Println(multi.ValueForProtocol(identity.TCPCode))
```
 
An example output of this code snippet is given below.

```
/ip4/127.0.0.1/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f
8MGfbzAMS59Gb4cSjpm34soGNYsM2f <nil>
127.0.0.1 <nil>
80 <nil>
```

## Tests

To run the test suite, install Ginkgo.

```
go get github.com/onsi/ginkgo/ginkgo
```

Now we can run the tests.

```
ginkgo -v
```

## Republic

The Identity library was developed by the Republic Protocol team and is available under the MIT license. For more information, see our website https://republicprotocol.com. 

## Contributors

* Loong
* Yunshi