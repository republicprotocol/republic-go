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

The Relay RESTful API provides HTTP POST, GET and DELETE commands to communicate with a relay. These methods will allow users to send orders and order fragments, retrieve orders and get real-time updates regarding order statuses, and cancel orders that are have not been matched.

The following sections will describe ways to connect to the Relay and use these commands to handle orders using the relay.

### Starting a new Relay

```go
relay := relay.NewRelay(configurations, darknodeRegistry, orderbook, relayerClient, smpcerClient, swarmerClient)
port := 18515 //Specify the port address. 
bind := 127.0.0.1 // Specify the bind address. 
if err := relay.ListenAndServe(bind, fmt.Sprintf("%d", port+1)); err != nil {
    log.Fatalf("error serving http: %v", err)
}
```

### Opening orders

The OpenOrdersHandler method can send full orders or fragmented orders to specific dark pools.

**Full orders**

```go
sendOrder := relay.OpenOrderRequest{}
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
sendOrder.Order = order
sendOrder.OrderFragments = relay.OrderFragments{}
s, _ := json.Marshal(sendOrder)
body := bytes.NewBuffer(s)
r := httptest.NewRequest("POST", "http://127.0.0.1:18515/orders", body)
w := httptest.NewRecorder()
relayNode := relay.Relay{}
relayNode.Config.Token = ""
handler := relay.RecoveryHandler(relayNode.AuthorizationHandler(relayNode.OpenOrdersHandler()))
handler.ServeHTTP(w, r)
```

**Order Fragments**

```go
sendOrder := relay.OpenOrderRequest{}
n := int64(17)
k := int64(12)
prime, err := stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")
if err != nil {
    log.Fatal(err)
}
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
fragments, err := order.Split(n, k, &prime)
if err != nil {
    log.Fatal(err)
}
// Create an OrderFragments object that stores details of the order
// along with the fragments for specific dark pools
fragmentedOrder := OrderFragments{}
fragmentsForPool := map[string][]*order.Fragment{}
fragmentsForPool["vrZhWU3VV9LRIM="] = fragments
fragmentedOrder.ID = []byte("vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM=")
fragmentedOrder.Type = 2
fragmentedOrder.Parity = 1
fragmentedOrder.Expiry = time.Time{}
fragmentedOrder.DarkPools = fragmentsForPool
sendOrder.Order = relay.Order{}
sendOrder.OrderFragments = fragmentedOrder
s, _ := json.Marshal(sendOrder)
body := bytes.NewBuffer(s)
r := httptest.NewRequest("POST", "http://127.0.0.1:18515/orders", body)
w := httptest.NewRecorder()
relayNode := relay.Relay{}
relayNode.Config.Token = "test"
handler := relay.RecoveryHandler(relayNode.AuthorizationHandler(relayNode.OpenOrdersHandler()))
handler.ServeHTTP(w, r)
```

## Canceling orders

```go
cancelRequest := relay.CancelOrderRequest{}
cancelRequest.ID = []byte("vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM=")
s, _ := json.Marshal(cancelRequest)
body := bytes.NewBuffer(s)
r := httptest.NewRequest("POST", "http://127.0.0.1:18515/orders", body)
w := httptest.NewRecorder()
relayNode := relay.Relay{}
relayNode.Config.Token = "test"
handler := relay.RecoveryHandler(relayNode.AuthorizationHandler(relayNode.CancelOrderHandler()))
handler.ServeHTTP(w, r)
```

## Getting an order

```go
// Retrieve an order with id vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM
r := httptest.NewRequest("GET", "http://127.0.0.1:18515/orders/vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM", nil)
w := httptest.NewRecorder()
relayNode := relay.Relay{}
relayNode.Config.Token = "test"
handler := relay.RecoveryHandler(relayNode.AuthorizationHandler(relayNode.CancelOrderHandler()))
handler.ServeHTTP(w, r)
order := new(order.Order)
if err := json.Unmarshal(w.Body.Bytes(), order); err != nil {
    log.Fatal(err)
}
```