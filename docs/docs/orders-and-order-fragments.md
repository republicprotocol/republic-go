# Orders and Order Fragments

## Orders

### Signature

The signature is an ECDSA secp256k1 s256 signature of the order ID, stored as bytes. The Republic Protocol address of the trader can be recover from this signature.

### ID

The ID is the keccak256 hash of the order, stored as bytes.

### Type

The order type is an enumerated value.

**Values**

* `0` — A mid-point order. The order will be matched at the mid-point between the best bid, and best offer, for the specified tokens. The price will be ignored.
* `1` — A limit order. The price must be specified.

### Parity

The order partity is an enumerated value.

**Values**

* `0` — The order is buying the higher token.
* `1` — The order is selling the higher token.

### Expiry

The order expiry is a 64-bit unsigned integer.

**Values**

Unix time, in seconds, after which the order will be automatically canceled.

### Tokens

The tokens are a 64-bit unsigned integer. It is the concatenation of two of 32-bit unsigned integers. The smaller value is always listed before the larger value, and is referred to as the higher token.

**Values**

* REN — 0
* ETH — 1
* BTC — 2
* DGX — 3

### Price

The price is a 16-bit unsigned integer. The first 8-bits are the co-efficient, _c_, and the second 8-bits are the exponent, _q_. The price is encoded as _0.1*c*10^(q-25.0)_. The value of _c_ must be an integer in the range 1 to 99 inclusively. The value of _q_ must be an integer in the range 0 to 52 inclusively.

The price is interpreted as the cost of 1 standard unit of the higher token, in 10e-12 units of the lower token. For example, if the price of REN is 0.0000095BTC, then the price should be listed as 9.5*10e6, where _c = 95_ and _q = 31_. 

### Volume

The volume is a 16-bit unsigned integer. The first 8-bits are the co-efficient, _c_, and the second 8-bits are the exponent, _q_. The price is encoded as _0.2*c*10^q_. The value of _c_ must be an integer in the range 1 to 49 inclusively. The value of _q_ must be an integer in the range 0 to 52 inclusively.

The price is interpreted as the maximum number of 10e-12 units of the higher token that can be traded by this order.

### Minimum Volume


The volume is a 16-bit unsigned integer. The first 8-bits are the co-efficient, _c_, and the second 8-bits are the exponent, _q_. The price is encoded as _0.2*c*10^q_. The value of _c_ must be an integer in the range 1 to 49 inclusively. The value of _q_ must be an integer in the range 0 to 52 inclusively.

The price is interpreted as the minimum number of 10e-12 units of the higher token that can be traded by this order.

### Nonce

The nonce is a 64-bit unsigned integer.

### Protobuf

The Protobuf format encodes an order, using definitions from the specification.

```proto
enum OrderType {
    MIDPOINT = 0;
    LIMIT = 1;
}

enum OrderParity {
    BUY = 0;
    SELL = 1;
}

message Order {
    bytes signature = 1;
    bytes id = 2;
    OrderType type = 3;
    OrderParity parity = 4;
    int64 expiry = 5;

    int64 tokens = 6;
    int16 price = 7;
    int16 volume = 8;
    int16 minimumVolume = 9;
    int64 nonce = 10;
}
```

### JSON

The JSON format encodes an order, using definitions from the specification with minor modifications.

```json
{
    "signature": "RmDYPYqqBTd2YBsInNPWySQUlMhPBMKduTqqhMRy0+w=",
    "id": "KthRO2cp4hwS+egr5xWYpdIPeyMEe1uHsPKp6Cut3co=",
    "type": 0,
    "parity": 0,
    "expiry": 1523238476,
    "tokens": 1,
    "price": [95, 31],
    "volume": [5, 13],
    "minimumVolume": [5, 12],
    "nonce": 1
}
```

- `signature` is encoded as a base64 string.
- `id` is encoded as a base64 string.
- `price` is encoded as a tuple of the co-efficient and exponent.
- `volume` is encoded as a tuple of the co-efficient and exponent.
- `minimumVolume` is encoded as a tuple of the co-efficient and exponent.

## OrderFragments

### Signature

The signature is an ECDSA secp256k1 s256 signature of the order fragment ID, stored as bytes. The Republic Protocol address of the trader can be recover from this signature.

### ID

The ID is the keccak256 hash of the order fragment, stored as bytes.

### Order Signature

The signature of the associated order. See definitions in Order.

### Order ID

The ID of the associated order. See definitions in Order.

### Order Type

The type of the associated order. See definitions in Order.

### Order Parity

The partity of the associated order. See definitions in Order.

### Expiry

The expiry of the associated order. See definitions in Order.

### Tokens

The tokens of the associated order, shared using shamir secret sharing, encoded into bytes, and encrypted by an 2048-bit RSA public key. See definitions in Order.

### Price

The price of the associated order, shared using shamir secret sharing, encoded into bytes, and encrypted by an 2048-bit RSA public key. See definitions in Order.

### Volume

The price of the associated order, shared using shamir secret sharing, encoded into bytes, and encrypted by an 2048-bit RSA public key. See definitions in Order.

### Minimum Volume

The price of the associated order, shared using shamir secret sharing, encoded into bytes, and encrypted by an 2048-bit RSA public key. See definitions in Order.

### Protobuf

The Protobuf format encodes an order fragment, using definitions from the specification.

```proto
enum OrderType {
    MIDPOINT = 0;
    LIMIT = 1;
}

enum OrderParity {
    BUY = 0;
    SELL = 1;
}

message OrderFragment {
    bytes signature = 1;
    bytes id = 2;
    bytes orderSignature = 3;
    bytes orderId = 4;
    OrderType orderType = 5;
    OrderParity orderParity = 6;
    int64 orderExpiry = 7;

    bytes tokens = 8;
    bytes price = 9;
    bytes volume = 10;
    bytes minimumVolume = 11;
}
```

### JSON

The JSON format encodes an order fragment, using definitions from the specification with minor modifications.

```json
{
    "signature": "RmDYPYqqBTd2YBsInNPWySQUlMhPBMKduTqqhMRy0+w=",
    "id": "KthRO2cp4hwS+egr5xWYpdIPeyMEe1uHsPKp6Cut3co=",
    "orderSignature": "Td2YBy0MRYPYqqBduRmDsIhTySQUlMhPBM+wnNPWKqq=",
    "orderId": "h1uHs+egr5xWYpwSdIPeyt36PKpKthROcoCMEe2cp4u=",
    "orderType": 0,
    "orderParity": 0,
    "orderExpiry": 1523238476,
    "tokens": "qQUhRMuTqlRmYqqPBMKdP0+whYMDIynNPWySBTd2YBs=",
    "price": ["YqqMDYP0+wBTd2YBsIqhMRynNPWySQUhPBMKduTqlRm=", "YqqMDYP0+wBTd2YBsIqhMRynNPWySQUhPBMKduTqlRm="],
    "volume": ["MhPRmDYPYWySQUlBsInNPduTqqqBTBMKqhMRy0+wd2Y=", "MhPRmDYPYWySQUlBsInNPduTqqqBTBMKqhMRy0+wd2Y="],
    "minimumVolume": ["hMRlMdBsInNPWyS2y0+mDYPYqKduTqqRqYQUBwhPBMT=", "hMRlMdBsInNPWyS2y0+mDYPYqKduTqqRqYQUBwhPBMT="]
}
```

- `signature` is encoded as a base64 string.
- `id` is encoded as a base64 string.
- `orderSignature` is encoded as a base64 string.
- `orderId` is encoded as a base64 string.
- `tokens` is an RSA encrypted number encoded as a base64 string.
- `price` is an RSA encrypted tuple of the co-efficient and exponent encoded as base64 strings.
- `volume` is an RSA encrypted a tuple of the co-efficient and exponent encoded as base64 strings.
- `minimumVolume` is an RSA encrypted a tuple of the co-efficient and exponent encoded as base64 strings.