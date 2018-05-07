# Relay

## SDK

A software development kit (SDK) for providing an abstraction over the Relay RESTful API. This includes creating orders, canceling orders, getting an order, and getting updates on the status of an order.

### Using the SDK

To use the Relay SDK, first you need to install the SDK.

```sh
go get -u github.com/republicprotocol/republic-go/relay
```

Now you can import it into whichever files need to use it.

```go
import "github.com/republicprotocol/republic-go/relay"
```

### Connecting to a Relay

Connect to the Relay with a call to `relay.NewAPI`, passing the connection parameters. These parameters will vary depending on the configuration of the Relay.

```go
url := "localhost:18515"
token := "mySecretToken"
api := relay.NewAPI(url, token)
```

- `url` – The host, and port, of the Relay.
- `token` – The Authorization token required by the Relay.

You can connect to a Relay that does not require a token using `relay.Insecure`. It is _not_ recommended to configure a Relay without an Authorization token, unless you are in a testing environment.

```go
url := "localhost:18515"
api := relay.NewAPI(url, relay.Insecure())
```

### Using the Relay API

**Opening orders**

```go
order := relay.Order{
    Type:      1,
    Parity:    1,
    Expiry:    time.Now().AddDate(0, 0, 7),
    FstCode:   1,
    SndCode:   2,
    Price:     100,
    MaxVolume: 100,
    MinVolume: 100,
}
orderID, err := api.OpenOrder(order)
if err != nil {
    log.Fatal(err)
}
```

**Closing orders**

```go
orderID := "vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM"
err := api.CloseOrder(orderID)
if err != nil {
    log.Fatal(err)
}
```

**Getting an order**

You can get the status information of an order by specifying its ID.

```go
orderID := "vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM"
order, err := api.GetOrder(orderID)
if err != nil {
    log.Fatal(err)
}
```

**Getting updates on an order status**

You can get a channel of updates for order statuses. This connects to a WebSocket exposed by the Relay, and reads all updates as they reach the Relay. A `relay.Filter` can be used to filter updates for a specific order (by ID), or specific statuses (use `relay.StatusAll` for all statuses).

```go
filter := relay.Filter{
    ID:     "vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM",
    Status: "confirmed, settled",
}
updates, errs := api.GetOrderbookUpdates(filter)
go func() {
    for {
        select {
        case update, ok := <-updates:
            // ...
        case err, ok := <-errs:
            // ...
    	}
    }
}()
```

## RESTful API

The Relay RESTful API supports the POST, GET and DELETE verbs to communicate with a Relay over HTTP. These methods will allow users to send orders and order fragments, retrieve orders and get real-time updates regarding order statuses, and cancel orders that are have not been matched.

The following sections will describe ways to connect to the Relay and use these commands to handle orders using the relay.

### Certificates

SSL certificates are not supported. TLS is not needed when transporting pre-signed, and pre-encrypted, orders to the Relay since the security of Republic Protocol assumes that all communication is easily observed by an adversary. In future, when Relays are self-hosted and can themselves sign and encrypt orders, self-signed certificates will be supported.

### Creating an Order

**Request**

```sh
HTTP/1.1 POST /orders
```

```json
{
}
```

**Response**

```json
{
}
```

### Canceling an order

**Request**

```sh
HTTP/1.1 DELETE /orders/{id}
```

**Response**

```json
{
}
```

### Getting an order

**Request**

```sh
HTTP/1.1 GET /orders/{id}
```

**Response**

```json
{
}
```

### Getting updates for orders

**Request**

Open a WebSocket with the Relay. Updates are streamed over the WebSocket from the server to the client. Using optional query paramters, the client can filter the stream.

```sh
HTTP/1.1 GET /orders?id={string}&status={string}&trader={string}
```

- `id` A base64 encoding of the order hash that will filter the stream for updates to this order only. Optional. No default. 
- `status` One or more of 'open', 'unconfirmed', 'confirmed', 'settled', 'canceled' , comma separated, will filter the stream for updates on these statuses only. Optional. No default.
- `trader` A Republic Protocol address that will filter the stream for updates for this trader only. Optional. No default.

**Response**

```json
{
}
```