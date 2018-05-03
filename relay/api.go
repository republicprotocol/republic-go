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
		Token: token,
	}
}

// OpenOrder opens a new order using the HTTP API.
func (api *API) OpenOrder(ty order.Type, parity order.Parity, expiry time.Time, fstCode, sndCode order.CurrencyCode, price, maxVolume, minVolume uint) (order.ID, error) {
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
		return nil, fmt.Errorf("cannot marshal open order request: %s", err)
	}

	// Create a request and set the authorization header
	req, err := http.NewRequest("POST", api.URL+"/orders", bytes.NewBuffer(json))
	req.Header.Set("Authorization", "Bearer "+api.Token)

	// Create a client and fetch request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot complete request: %s", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body: %s", err)
	}

	// TODO: Handle response
	fmt.Println("Status: ", resp.Status)
	fmt.Println("Headers: ", resp.Header)
	fmt.Println("Body: ", string(body))

	return ord.ID, nil
}

// CancelOrder cancels an existing order using the HTTP API.
func (api *API) CancelOrder(orderID order.ID) error {
	// Create a request and set the authorization header
	req, err := http.NewRequest("DELETE", api.URL+"/orders/"+orderID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+api.Token)

	// Create a client and fetch request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot complete request: %s", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %s", err)
	}

	// TODO: Handle response
	fmt.Println("Status: ", resp.Status)
	fmt.Println("Headers: ", resp.Header)
	fmt.Println("Body: ", string(body))

	return nil
}

// GetOrder gets an existing order using the HTTP API.
func (api *API) GetOrder(orderID order.ID) (order.Order, error) {
	// Create a request and set the authorization header
	req, err := http.NewRequest("GET", api.URL+"/orders/"+orderID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+api.Token)

	// Create a client and fetch request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return order.Order{}, fmt.Errorf("cannot complete request: %s", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return order.Order{}, fmt.Errorf("cannot read response body: %s", err)
	}

	// TODO: Handle response
	fmt.Println("Status: ", resp.Status)
	fmt.Println("Headers: ", resp.Header)
	fmt.Println("Body: ", string(body))

	var ord order.Order
	if err := json.Unmarshal(body, ord); err != nil {
		return order.Order{}, fmt.Errorf("cannot unmarshal message: %s", err)
	}

	return ord, nil
}

// GetOrders retrieves updates to an existing order using the HTTP API. This
// function returns a channel which will contain all the updates that are
// received.
func (api *API) GetOrders(filter Filter) (<-chan order.Order, <-chan error) {
	orders := make(chan order.Order)
	errs := make(chan error)

	// Construct WebSocket URL using filters
	u := url.URL{Scheme: "ws", Host: api.URL, Path: "/orders"}
	query := u.String() + "?id=" + filter.ID
	if filter.Status != "" {
		query = query + "&status=" + filter.Status
	}

	// Set authorization header and connect to WebSocket
	var header http.Header
	header.Set("Authorization", "Bearer "+api.Token)
	c, _, err := websocket.DefaultDialer.Dial(query, header)
	if err != nil {
		errs <- err
		close(errs) // TODO: Check this
		return orders, errs
	}
	defer c.Close()

	// Write any WebSocket messages to channel
	go func() {
		defer close(orders)
		defer close(errs)
		for {
			var ord order.Order
			_, message, err := c.ReadMessage()
			if err != nil {
				errs <- err
				continue // TODO: Confirm we want to continue after a read error
			}
			json.Unmarshal(message, ord)
			orders <- ord
		}
	}()

	return orders, errs
}

// Insecure returns an empty string used for when when instantiating a new API
// object without providing a token.
func Insecure() string {
	return ""
}
