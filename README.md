# Distributed Hash Table

The Distributed Hash Table (DHT) library is a Go implementation of a Kademlia DHT structure. The library does not implement the Kademlia logic, only the underlying data structure and helper functions. We recommend the Wikipedia page as the starting point for readers that are not familiar with Distributed Hash Tables, or Kademlia.

* [Distributed Hash Table](https://en.wikipedia.org/wiki/Distributed_hash_table)
* [Kademlia](https://en.wikipedia.org/wiki/Kademlia)

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