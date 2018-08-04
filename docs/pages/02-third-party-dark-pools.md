# Third-Party Dark Pools

> Work in progress.

Dark pools are built on top of Republic Protocol and use the decentralised Secure Order Matcher to implement an order matcher that does not expose the price or volume of orders. Dark pools are responsible for implementing settlement and brokerage. Dark pools can be centralised or decentralised, as long as the interfaces required by Republic Protocol are exposed.

## Settlement

See https://github.com/republicprotocol/republic-sol for more information.

## Brokerage

```sol
function verifyOrderSignature(bytes32 _order, bytes _signature) returns (bool) { /* ... */ }
```

```sol
function verifyOrderSignatureAndPay(bytes32 _order, bytes _signature, uint256 _value) returns (bool) { /* ... */ }
```

## RenEx

RenEx is an official third-party dark pool. It is built, and maintained, by the Republic Protcol team.

See https://github.com/republicprotoco/renex-sol for more information.