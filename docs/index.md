# Republic Protocol

* [Darknodes](./pages/01-darknodes.md)
* [Dark Pools](./pages/02-third-party-dark-pools.md)
* [Secure Order Matcher](./pages/03-secure-order-matcher.md)
* [Byzantine Fault Tolerance](./pages/04-byzantine-fault-tolerance.md)

Republic Protocol is an open-source protocol powering decentralised dark pools. It is run by a decentralised network of Darknodes that using [secure multi-party computation](https://en.wikipedia.org/wiki/Secure_multi-party_computation) to match orders without exposing the price, or volume, of the orders. Dark pools powered by Republic Protocol can support large volume trades, with minimal price slippage and market impact, whilst guaranteeing that the rules of the dark pool cannot be broken.

## How it works

![Overview](./assets/images/00-index-diagram.jpg "Overview")

**Traders** open orders by submitting a transaction to the Orderbook and sending order fragments to the network of Darknodes. On the Orderbook, the trader designates a [dark pool](./pages/02-third-party-dark-pools.md) for the order and the Orderbook verifies that the order has been signed by the required brokers.

**Darknodes** synchronise orders from the Orderbook and receive order fragments directly from the traders. After receiving its order fragment, a Darknode runs a secure multi-party computation with other Darknodes to find a matching order without exposing the price and volume. Unless the majority of Darknodes are acting maliciously, the security of orders and the integrity of the order matching rules cannot be corrupted.

**Orderbook** stores the priority, and the state, of orders. It stores this data against the order hash of the order, exposing nothing about the price and volume.

**Settlemnet** is defined differently by different dark pools, but all adhere to a common interface defined by Republic Protocol.