package relay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/stackint"
)

// API objects are constructed to interact with the Relay HTTP layer.
type API struct {
	URL   string
	Token string
}

// Order objects are construct for opening orders using the OpenOrder function.
type Order struct {
	Type      order.Type
	Parity    order.Parity
	Expiry    time.Time
	FstCode   order.CurrencyCode
	SndCode   order.CurrencyCode
	Price     uint // TODO: Accept int64 instead of uint
	MaxVolume uint
	MinVolume uint
}

// Filter objects are constructed to filter updates when retrieving updates
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

// OpenOrder opens a new order using the Relay API.
func (api *API) OpenOrder(ord Order) (order.ID, error) {
	// Construct order request using parameters
	nonce := stackint.Zero() // TODO: Set nonce
	newOrder := order.NewOrder(ord.Type, ord.Parity, ord.Expiry, ord.FstCode, ord.SndCode, stackint.FromUint(ord.Price), stackint.FromUint(ord.MaxVolume), stackint.FromUint(ord.MinVolume), nonce)
	orderRequest := OpenOrderRequest{
		Order:          *newOrder,
		OrderFragments: OrderFragments{},
	}
	json, err := json.Marshal(orderRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal open order request: %s", err)
	}

	// Create a request and set the authorization header
	req, err := http.NewRequest("POST", api.URL+"/orders", bytes.NewBuffer(json))
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %s", err)
	}
	req.Header.Set("Authorization", "Bearer "+api.Token)

	// Create a client and fetch request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot send request: %s", err)
	}
	defer resp.Body.Close()

	// Verify response
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf(fmt.Sprintf("invalid response status code: %d", resp.StatusCode))
	}

	return newOrder.ID, nil
}

// CancelOrder cancels an existing order using the Relay API.
func (api *API) CancelOrder(orderID string) error {
	// Create a request and set the authorization header
	req, err := http.NewRequest("DELETE", api.URL+"/orders/"+orderID, nil)
	if err != nil {
		return fmt.Errorf("cannot create request: %s", err)
	}
	req.Header.Set("Authorization", "Bearer "+api.Token)

	// Create a client and fetch request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot send request: %s", err)
	}
	defer resp.Body.Close()

	// Verify response
	if resp.StatusCode != http.StatusGone {
		return fmt.Errorf(fmt.Sprintf("invalid response status code: %d", resp.StatusCode))
	}

	return nil
}

// GetOrder gets an existing order using the Relay API.
func (api *API) GetOrder(orderID string) (order.Order, error) {
	// Create a request and set the authorization header
	req, err := http.NewRequest("GET", api.URL+"/orders/"+orderID, nil)
	if err != nil {
		return order.Order{}, fmt.Errorf("cannot create request: %s", err)
	}
	req.Header.Set("Authorization", "Bearer "+api.Token)

	// Create a client and fetch request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return order.Order{}, fmt.Errorf("cannot send request: %s", err)
	}
	defer resp.Body.Close()

	// Verify response
	if resp.StatusCode != http.StatusOK {
		return order.Order{}, fmt.Errorf(fmt.Sprintf("invalid response status code: %d", resp.StatusCode))
	}

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return order.Order{}, fmt.Errorf("cannot read response body: %s", err)
	}

	var ord order.Order
	if err := json.Unmarshal(body, ord); err != nil {
		return order.Order{}, fmt.Errorf("cannot unmarshal message: %s", err)
	}

	return ord, nil
}

// GetOrderbookUpdates retrieves updates to an existing order using the Relay API. This
// function returns a channel which will contain all the updates that are
// received.
func (api *API) GetOrderbookUpdates(filter Filter) (<-chan orderbook.Entry, <-chan error) {
	entries := make(chan orderbook.Entry)
	errs := make(chan error)

	// Construct WebSocket URL using filters
	u := url.URL{Scheme: "ws", Host: api.URL, Path: "/orders"}
	query := u.String() + "?id=" + filter.ID
	if filter.Status != "" {
		query = query + "&status=" + strings.Replace(filter.Status, " ", "", -1)
	}

	// Set authorization header and connect to WebSocket
	var header http.Header
	header.Set("Authorization", "Bearer "+api.Token)
	c, _, err := websocket.DefaultDialer.Dial(query, header)
	if err != nil {
		errs <- err
		close(errs) // TODO: Check this
		return entries, errs
	}
	defer c.Close()

	// Write WebSocket messages to the channels
	go func() {
		defer close(entries)
		defer close(errs)
		for {
			var entry orderbook.Entry
			_, message, err := c.ReadMessage()
			if err != nil {
				errs <- err
				continue // TODO: Confirm we want to continue after a read error
			}
			json.Unmarshal(message, entry)
			entries <- entry
		}
	}()

	return entries, errs
}

// Insecure returns an empty string used for when instantiating a new API
// object without providing a token.
func Insecure() string {
	return ""
}

// StatusAll returns an empty string used for when the user wishes to receive
// all updates to an order.
func StatusAll() string {
	return ""
}
