# Relay SDK

A software development kit (SDK) for providing an abstraction over the Relay RESTful API. This includes creating orders, canceling orders, getting an order, and getting updates on the status of an order.

## Using the SDK

### Import the Relay package

```go
import "github.com/republicprotocol/republic-go/relay"
```

### Construct a new Relay API object

```go
api := relay.NewAPI(url, token)
```

- `url` – The address your Relay is serving from.
- `token` – The authentication token used to verify your Relay.

> Note: If your Relay does not require an authentication token use `relay.Insecure()`.

### Interacting with the Relay API

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
    // Handle error
}
```

> Note: For information regarding the `relay.Order` fields, please refer to [...]

**Closing orders**

```go
orderID := "vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM"
err := api.CloseOrder(orderID)
if err != nil {
    // Handle error
}
```

**Getting an order**

```go
orderID := "vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM"
order, err := api.GetOrder(orderID)
if err != nil {
    // Handle error
}
```

**Getting updates on the status of an order**

```go
filter := relay.Filter{
    ID:     "vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQyV6ryi1wDSM",
    Status: "confirmed, settled",
}
updates, errs := GetOrderbookUpdates(filter)
go func() {
    for {
        select {
            case entry, ok := <-updates:
                if !ok {
                    return
                }
                status := entry.Status
                // ...
    		case <-errs:
    			// Handle errors
    	}
    }
}()
```

> Note: If you wish to receive updates regarding all statuses, use `relay.StatusAll()`.
