# Ingress

An Ingress service exposes APIs that abstract over Republic Protocol. This can be useful for those that want to interact with Republic Protocol but are unfamiliar with gRPC, Ethereum, or do not want to continually acquire REN.

The REN Terminal is an example of an Ingress service. It provides a web application, and RESTful API, for managing orders. One of the core services it provides is the payment of REN fees on behalf of the user. This reduces the barrier for users, since they do not need access to REN to access Republic Protocol.

## RESTful API

The RESTful API supports the POST and DELETE verbs. These methods allow users to open and cancel orders over HTTP, without needing to spend REN.

### Certificates

SSL certificates are not supported. TLS is not needed when transporting pre-signed, and pre-encrypted, orders to the Broker since the security of Republic Protocol assumes that all communication is easily observed by an adversary. In future, when Brokers are self-hosted and can themselves sign and encrypt orders, self-signed certificates will be supported.

### Creating an Order

**Request**

```
HTTP/1.1 POST /orders
```

```json
{
    "signature": "Td2YBy0MRYPYqqBduRmDsIhTySQUlMhPBM+wnNPWKqq=",
    "orderFragmentMapping": {
        "d2YDsRUlMhPBMMRYIhTy+nNPTQKqqBwqq+duPYWm0yS=": [{
            "orderSignature": "Td2YBy0MRYPYqqBduRmDsIhTySQUlMhPBM+wnNPWKqq=",
            "orderId": "h1uHs+egr5xWYpwSdIPeyt36PKpKthROcoCMEe2cp4u=",
            "orderType": 0,
            "orderParity": 0,
            "orderExpiry": 1523238476,
            "id": "KthRO2cp4hwS+egr5xWYpdIPeyMEe1uHsPKp6Cut3co=",
            "tokens": "qQUhRMuTqlRmYqqPBMKdP0+whYMDIynNPWySBTd2YBs=",
            "price": ["YqqMDYP0+wBTd2YBsIqhMRynNPWySQUhPBMKduTqlRm=", "YqqMDYP0+wBTd2YBsIqhMRynNPWySQUhPBMKduTqlRm="],
            "volume": ["MhPRmDYPYWySQUlBsInNPduTqqqBTBMKqhMRy0+wd2Y=", "MhPRmDYPYWySQUlBsInNPduTqqqBTBMKqhMRy0+wd2Y="],
            "minimumVolume": ["hMRlMdBsInNPWyS2y0+mDYPYqKduTqqRqYQUBwhPBMT=", "hMRlMdBsInNPWyS2y0+mDYPYqKduTqqRqYQUBwhPBMT="]
        }]
    }
}
```

- `signature` A base64 encoded string of the signature that authenticates the open order request.
- `orderFragmentMapping` A mapping of pods to lists of encrypted order fragments.

**Response**

- `201` The order was succecssfully opened.
- `400` The JSON was malformed.
- `500` An internal error was encountered — e.g. the Broker has insufficient REN.

```json
{
}
```

### Canceling an order

**Request**

```
HTTP/1.1 DELETE /orders/?id={id}&signature={signature}
```

- `id` A base64 encoded keccak256 hash of the order.
- `signature` A base64 encoded string of the signature that authenticates the cancel order request.

**Response**

- `200` The order was successfully canceled.
- `400` The order ID, or signature, was malformed.
- `500` An internal error was encountered — e.g. the Broker has insufficient REN.

```json
{
}
```