# Distributed Hash Table

The Distributed Hash Table (DHT) library is a Go implementation of a Kademlia DHT structure. The library does not implement the Kademlia logic, only the underlying data structure and helper functions. We recommend the Wikipedia page as the starting point for readers that are not familiar with Distributed Hash Tables, or Kademlia.

* [Distributed Hash Table](https://en.wikipedia.org/wiki/Distributed_hash_table)
* [Kademlia](https://en.wikipedia.org/wiki/Kademlia)

## How it works

You can create an empty DHT using the `NewDHT` function, and a Republic address. For more information about Republic addresses, see the (Identity)[https://github.com/republicprotocol/go-identity] library.

```go
address, _, err := identity.NewAddress()
table := dht.NewDHT(address)
```

You can add new entries to the DHT by using the `Update` function. If the entry already exists, its associated multi-address and timestamp will be updated.

```go
target, _, err := identity.NewAddress()
multi, err := target.MultiAddress()
err = table.Update(multi)
```

You can remove an entry from the DHT by using the `Remove` function. This function does nothing if the target is not in the DHT.

```go
target, _, err := identity.NewAddress()
multi, err := target.MultiAddress()
err = table.Remove(multi)
```

That's all there is to it! Remember, you should always check that the value of `err` is `nil`.

## Tests

To run the test suite, install Ginkgo.

```sh
go get github.com/onsi/ginkgo/ginkgo
```

Now we can run the tests.

```sh
ginkgo -v
```

## Republic

The DHT library was developed by the Republic Protocol team and is available under the MIT license. For more information, see our website https://republicprotocol.com.

## Contributors

* Loong
* Yunshi