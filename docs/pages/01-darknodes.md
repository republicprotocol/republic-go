# Darknodes

[Darknodes](./01-darknodes.md) are nodes that communicate with each other to run the [Secure Order Matcher](./03-secure-order-matcher.md). By distributing orders to the Darknodes using [Shamir's secret sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing), the Darknodes are able to match orders without knowing their price or volume. Once a matching pair of orders is found, the Darknodes claim a reward fee. There is only one network of Darknodes in Republic Protocol, matching orders for all [dark pools](./02-third-party-dark-pools.md).

## Darknode Registry

The Darknode Registry is an Ethereum smart contract used to register, and deregister, Darknodes. Before a Darknode is accepted into the network by others, it must submit a bond of 100,000REN to the Darknode Registry. After deregistration, the Darknode can refund this bond. It also defines a network topology that partitions Darknodes into Darknode Pods.

See https://github.com/republicprotocol/republic-sol for more information.

## Syncing the Orderbook

The Orderbook is an Ethereum smart contract used to store the state of orders. Before an honest Darknode admits an order into the [Secure Order Matcher](./03-secure-order-matcher.md), it will verify that the order has been signed by the required brokers and is currently in the `Open` state. The required brokers are defined by the dark pool on which the order was opened.

After finding matching orders, Darknodes will change the state of the orders to `Confirmed` in the Orderbook. The Orderbook will reject the confirmation of orders that are not currently in the `Open` state — preventing the asynchronous nature of the decentralised network from resulting in conflicting order matches. Honest Darknodes watch the Orderbook for state changes to `Open` orders. Orders can transition to the `Confirmed` state when a match is found, or to the `Canceled` state when the trader cancels the order. Expiration is managed off-chain by the Darknodes.

The Darknode Slasher is used to protect against malicious confirmations. See https://github.com/republicprotocol/republic-sol for more information.

## Syncing order fragments

Before an honest Darknode admits an order into the [Secure Order Matcher](./03-secure-order-matcher.md), it needs to receive its order fragment for that order. Order fragments store the public data of an order (buy/sell, id, type, expiry, and settlement layer), and Shamir's secret shares of the private data (price, and volume). Merkle proofs — built from the hashes of order fragments — are signed by traders, and required brokers, and are verified by honest Darknodes. Order fragments are required because the Shamir's secret shares are needed when running the Secure Order Matcher. They are transferred directly to the Darknodes to prevent public discovery and to reduce gas fees.

Once an honest Darknodes has synchronised an order from the Orderbook, and it has received its appropriate order fragment, it admits the order into the Secure Order Matcher and begins the process of finding a match.

## Settlement

After finding matching orders, and confirming the match in the Orderbook, Darknodes will use the order fragments to reconstruct the private data of the orders. At this stage, everything about the matched orders is known to the Darknodes and the matched orders are submitted to the Settlement Layer of the appropriate dark pool. The Settlement Layer is an Ethereum smart contract that exposes the ABI required by Republic Protocol. Third-party dark pools must define *at least one* Settlement Layer. A Settlement layer accepts matched orders, verifies that they are confirmed in the Orderbook, and pays a reward to the Darknodes that initiated the settlement.

The first Darknode to acquire a sufficient number of order fragments for the matched orders will receive the reward, assuming it initiates the settlement in a timely fashion. To ensure faireness, honest Darknodes follow a round-robin schedule when distributing their order fragments; forwarding to the first Darknode in the schedule, then the second, and so on. The result is that the Darknode initiating the settlement, and earning the reward, alternates each time matching orders are found.

See https://github.com/republicprotocol/republic-sol for more information.

See https://github.com/republicprotocol/renex-sol for an example.

*Note: The actual exchange of cryptographic assets from one trader to another is defined by third-party dark pools. There is no strict requirement defined by Republic Protocol, beyond the payment of fees to the Darknodes.*