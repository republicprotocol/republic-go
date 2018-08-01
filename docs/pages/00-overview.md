# Overview

Republic Protocol is an open-source protocol powering decentralised dark pool exchanges. It is run by a decentralised network of [Darknodes](./01-darknodes.md) that using [secure multi-party computation](https://en.wikipedia.org/wiki/Secure_multi-party_computation) to match orders without exposing the price, or volume, of the orders.

Dark pool exchanges powered by Republic Protocol can support large volume trades, with minimal price slippage and market impact, whilst guaranteeing that the rules of the dark pool cannot be broken.

### Darknodes

[Darknodes](./01-darknodes.md) are nodes that communicate with each other to run the [Secure Order Matcher](./03-secure-order-matcher.md). By distributing orders to the Darknodes using [Shamir's secret sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing), the Darknodes are able to match orders without knowing their price or volume. Once a matching pair of orders is found, the Darknodes claim a reward fee.

There is only one network of Darknodes in Republic Protocol, matching orders for all [dark pools](./02-third-party-dark-pools.md).

### Dark pools

[Dark pools](./02-third-party-dark-pools.md) are third-party exchanges, built on top of Republic Protocol, that use the [Secure Order Matcher](./03-secure-order-matcher.md) to implement a decentralised order matcher that does not expose the price or volume of orders. Dark pools can be centralised or decentralised, as long as they respect the required interface of the protocol.

Dark pool exchanges are responsible for implementing trading interfaces, settlement, brokerage, and adhering to regulation. The core respos
