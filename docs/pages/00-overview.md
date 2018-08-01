# Overview

Republic Protocol is an open-source protocol powering decentralised dark pools. It is run by a decentralised network of [Darknodes](./01-darknodes.md) that using [secure multi-party computation](https://en.wikipedia.org/wiki/Secure_multi-party_computation) to match orders without exposing the price, or volume, of the orders. Dark pools powered by Republic Protocol can support large volume trades, with minimal price slippage and market impact, whilst guaranteeing that the rules of the dark pool cannot be broken.

### [Darknodes](./01-darknodes.md)

[Darknodes](./01-darknodes.md) are nodes that communicate with each other to run the Secure Order Matcher. By distributing orders to the Darknodes using [Shamir's secret sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing), the Darknodes are able to match orders without knowing their price or volume. Once matching orders are found, the orders are settled on the respective dark pool and a reward is paid to the Darknodes.

### [Dark pools](./02-third-party-dark-pools.md)

[Dark pools](./02-third-party-dark-pools.md) are third-party exchanges, built on top of Republic Protocol, that use the Secure Order Matcher to implement a decentralised order matcher that does not expose the price or volume of orders. Dark pools are responsible for implementing trading interfaces, settlement, brokerage, and adhering to necessary regulation. Dark pools can be centralised or decentralised, as long as the interface required by Republic Protocol is exposed.