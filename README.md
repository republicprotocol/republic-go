# Distributed Hash Table

The Distributed Hash Table (DHT) library is a Go implementation of a Kademlia DHT structure. The library does not implement the Kademlia logic, only the underlying data structure and helper functions. We recommend the Wikipedia page as the starting point for readers that are not familiar with Distributed Hash Tables, or Kademlia.

* [Distributed Hash Table](https://en.wikipedia.org/wiki/Distributed_hash_table)
* [Kademlia](https://en.wikipedia.org/wiki/Kademlia)

## How it works

You can create an empty DHT using the `NewDHT` function, and a Republic address. For more information about Republic addresses, see the [Identity](https://github.com/republicprotocol/go-identity) library.

```go
address, _, err := identity.NewAddress()
table := dht.NewDHT(address)
```

### Updating the DHT

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

These are the two methods used to add, update, and remove entries in the DHT. Remember, you should always check that the value of `err` is `nil`.

### Finding entries in the DHT

To find an exact entry in the DHT we can use the `FindMultiAddress` function. Given a target Republic address, this will return a multi-address that contains that same Republic address but also includes networking information.

```go
target, _, err := identity.NewAddress()
multi, err := table.FindMultiAddress(target)
```

Entries in the DHT are organized into buckets of addresses that are a similar distance from the DHT address. Distance is calculated by performing a bitwise XOR between addresses. To get the bucket for an address, we can use the `FindBucket` function.

```go
target, _, err := identity.NewAddress()
bucket, err := table.FindBucket(multi)
```

This returns a pointer to a bucket. This is useful for implementing the logic behind Kademlia, and a pointer is used so that the bucket can be updated if necessary.

By default, buckets are sorted by how recently an entry has been added, with newer entries at the end of the bucket. However, we can use the `Sort` method to explicitly sort the bucket, just in case.

```go
bucket.Sort()
```

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