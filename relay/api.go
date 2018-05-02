package relay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/stackint"
)

// API objects are constructed to interact with the Relay HTTP layer.
type API struct {
	URL   string
	Token string
}

// Filter objects are constructed to filter updates when retrieving data
// regarding an order using the GetOrders function.
type Filter struct {
	ID     string
	Status string
}

// NewAPI returns a new API object.
func NewAPI(url string, token string) API {
	return API{
		URL:   url,
		Token: token, // TODO: Handle token
	}
}

// OpenOrder opens a new order using the HTTP API.
func (api *API) OpenOrder(ty order.Type, parity order.Parity, expiry time.Time, fstCode, sndCode order.CurrencyCode, price, maxVolume, minVolume uint) {
	// Construct order request using parameters
	nonce := stackint.FromUint(0) // TODO: Set nonce
	// TODO: Accept int64 instead of uint?
	ord := order.NewOrder(ty, parity, expiry, fstCode, sndCode, stackint.FromUint(price), stackint.FromUint(maxVolume), stackint.FromUint(minVolume), nonce)
	orderRequest := OpenOrderRequest{
		Order:          *ord,
		OrderFragments: OrderFragments{},
	}
	json, err := json.Marshal(orderRequest)
	if err != nil {
		// TODO: Handle error
	}

	// Create a request
	req, err := http.NewRequest("POST", api.URL+"/orders", bytes.NewBuffer(json))

	// Create a client and fetch request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// TODO: Handle error
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: Handle response
	fmt.Println("Status: ", resp.Status)
	fmt.Println("Headers: ", resp.Header)
	fmt.Println("Body: ", string(body))
}

// CancelOrder cancels an existing order using the HTTP API.
func (api *API) CancelOrder(orderID order.ID) {
	// Construct cancel request using parameters
	cancelRequest := CancelOrderRequest{
		ID:        orderID,
		Signature: []byte{}, // TODO: Handle signature
	}
	json, err := json.Marshal(cancelRequest)
	if err != nil {
		// TODO: Handle error
	}

	// Create a request
	// TODO: Check use of DELETE
	req, err := http.NewRequest("POST", api.URL+"/orders", bytes.NewBuffer(json))

	// Create a client and fetch request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// TODO: Handle error
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: Handle response
	fmt.Println("Status: ", resp.Status)
	fmt.Println("Headers: ", resp.Header)
	fmt.Println("Body: ", string(body))
}

// GetOrder gets an existing order using the HTTP API.
func (api *API) GetOrder(orderID order.ID) {
	// Fetch request
	resp, err := http.Get(api.URL + "/orders/" + orderID.String())
	if err != nil {
		// TODO: Handle error
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: Handle response
	fmt.Println("Status: ", resp.Status)
	fmt.Println("Headers: ", resp.Header)
	fmt.Println("Body: ", string(body))
}

// GetOrders retrieves updates to an existing order using the HTTP API. This
// function returns a channel which will contain all the updates that are
// received.
func (api *API) GetOrders(filter Filter) <-chan order.Order {
	// Construct WebSocket URL using filters
	u := url.URL{Scheme: "ws", Host: api.URL, Path: "/orders"}
	query := u.String() + "?id=" + filter.ID
	if filter.Status != "" {
		query = query + "&status=" + filter.Status
	}

	// Connect to WebSocket
	c, _, err := websocket.DefaultDialer.Dial(query, nil)
	if err != nil {
		// TODO: Handle error
	}
	defer c.Close()

	// Write any WebSocket messages to channel
	orders := make(chan order.Order)
	go func() {
		defer close(orders)
		for {
			var ord order.Order
			_, message, err := c.ReadMessage()
			if err != nil {
				// TODO: Handle error
			}
			json.Unmarshal(message, ord)
			orders <- ord
		}
	}()

	return orders
}

// Insecure returns an empty string used for when when instantiating a new API
// object without providing a token.
func Insecure() string {
	return ""
}
