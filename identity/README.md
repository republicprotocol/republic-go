# Identity

The Identity library is an official reference implementation of identities and addresses in the Republic Protocol. It supports the generation and use of ECDSA SECP256K private keys, Republic Protocol IDs, Republic Protocol addresses, and multi-addresses.

## Key Pairs

An identity in Republic Protocol is defined by an ECDSA SECP256K private key, the same kind of private key used by both Bitcoin and Ethereum.

```go
keyPair, err := identity.NewKeyPair()
```

## Republic Protocol IDs

Republic Protocol IDs are generated from your private key, and are used to identify traders and nodes in Republic Protocol. To generate a Republic Protocol ID, take the last 20 bytes of the Keccak256 hash of your ECDSA public key. 

```go
keyPair, err := identity.NewKeyPair()
id := keyPair.ID()
fmt.Println(id)
```

An example output of this code snippet is given below.

```
[247 107 192 130 208 33 69 182 183 127 237 24 22 83 157 99 207 221 9 243]
```

## Republic Protocol addresses

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

## License

The Identity library was developed by the Republic Protocol team and is available under the MIT license. For more information, see our website https://republicprotocol.com.
>>>>>>> c0ee5d01eece4eab67cd4e16179098407a1f32ae
