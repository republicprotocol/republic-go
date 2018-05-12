# Broker

## RESTful API

The RESTful API supports the POST and DELETE verbs to communicate with a Broker over HTTP. These methods allow traders to open and cancel orders .

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
            "id": "KthRO2cp4hwS+egr5xWYpdIPeyMEe1uHsPKp6Cut3co=",
            "orderSignature": "Td2YBy0MRYPYqqBduRmDsIhTySQUlMhPBM+wnNPWKqq=",
            "orderId": "h1uHs+egr5xWYpwSdIPeyt36PKpKthROcoCMEe2cp4u=",
            "orderType": 0,
            "orderParity": 0,
            "orderExpiry": 1523238476,
            "tokens": "qQUhRMuTqlRmYqqPBMKdP0+whYMDIynNPWySBTd2YBs=",
            "price": "YqqMDYP0+wBTd2YBsIqhMRynNPWySQUhPBMKduTqlRm=",
            "volume": "MhPRmDYPYWySQUlBsInNPduTqqqBTBMKqhMRy0+wd2Y=",
            "minimumVolume": "hMRlMdBsInNPWyS2y0+mDYPYqKduTqqRqYQUBwhPBMT="
        }]
    }
}
```

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
HTTP/1.1 DELETE /orders/{id}?signature={signature}
```

- `id` An order ID (e.g. `h1uHs+egr5xWYpwSdIPeyt36PKpKthROcoCMEe2cp4u=`) of an open order.
- `signature` A signature from the trader that opened the order. The signed message must be `Republic Protocol: cancel order: ` prefixed to the `id`.

**Response**

- `200` The order was successfully canceled.
- `400` The order ID, or signature, was malformed.
- `500` An internal error was encountered — e.g. the Broker has insufficient REN.

```json
{
}
```