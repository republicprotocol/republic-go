package relay

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/republicprotocol/republic-go/dark"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/stackint"
)

var prime, _ = stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")

const reset = "\x1b[0m"
const red = "\x1b[31;1m"

// The HTTPDelete object
type HTTPDelete struct {
	signature []byte   `json:"signature"`
	ID        order.ID `json:"id"`
}

// Fragments will store a list of Fragments with their order details
type Fragments struct {
	// TODO: Confirm this . .
	DarkPool []byte `json:"darkPool"`

	Fragment []*order.Fragment `json:"fragment"`
}

// OrderFragments will store a list of Fragment Sets with their order details
type OrderFragments struct {
	Signature []byte   `json:"signature"`
	ID        order.ID `json:"id"`

	Type   order.Type   `json:"type"`
	Parity order.Parity `json:"parity"`
	Expiry time.Time    `json:"expiry"`

	FragmentSet []Fragments `json:"fragmentSet"`
}

// RecoveryHandler handles errors while processing the requests and populates the errors in the response
func RecoveryHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("%v", r)))
			}
		}()
		h.ServeHTTP(w, r)
	})
}

// PostOrdersHandler handles all HTTP Post requests
func PostOrdersHandler(multiAddress *identity.MultiAddress, darkPools dark.Pools) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Get this checked . .
		postOrder := order.Order{}
		if err := json.NewDecoder(r.Body).Decode(&postOrder); err != nil {
			postOrder := OrderFragments{}
			if err1 := json.NewDecoder(r.Body).Decode(&postOrder); err1 != nil {
				panic(fmt.Sprintf("cannot decode json into an order or a list of order fragments: %v %v", err, err1))
				return
			}
			SendOrderFragmentsToDarkOcean(postOrder, multiAddress, darkPools)
		}
		SendOrderToDarkOcean(postOrder, multiAddress, darkPools)
	})
}

func HandleGetOrder() http.Handler {
	// To-do: Add authentication.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["orderID"]
		// status := getOrderStatus(id)
		status := "confirmed"
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{} {
			"id":     id,
			"status": status,
		}); err != nil {
			fmt.Sprintf("cannot encode json object: %v", err)
		}
	})
}

// DeleteOrderHandler handles HTTP Delete Requests
func DeleteOrderHandler(multiAddress *identity.MultiAddress, darkPools dark.Pools) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cancelOrder := HTTPDelete{}
		err := json.NewDecoder(r.Body).Decode(&cancelOrder)
		if err != nil {
			panic(fmt.Sprintf("cannot decode json: %v", err))
			return
		}
		CancelOrder(cancelOrder.ID, multiAddress, darkPools)
	})
}
