package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/http"
	"github.com/republicprotocol/republic-go/http/adapter"
)

type weakAdapter struct {
	numOpened   int64
	numCanceled int64
}

func (adapter *weakAdapter) OpenOrder(signature string, orderFragmentMapping adapter.OrderFragmentMapping) error {
	atomic.AddInt64(&adapter.numOpened, 1)
	return nil
}

func (adapter *weakAdapter) CancelOrder(signature string, orderID string) error {
	atomic.AddInt64(&adapter.numCanceled, 1)
	return nil
}

type errAdapter struct {
}

func (adapter *errAdapter) OpenOrder(signature string, orderFragmentMapping adapter.OrderFragmentMapping) error {
	return errors.New("cannot open order")
}

func (adapter *errAdapter) CancelOrder(signature string, orderID string) error {
	return errors.New("cannot close order")
}

var _ = Describe("HTTP handlers", func() {

	Context("when opening orders", func() {

		It("should return a StatusCreated response for a valid request", func() {
			openAdapter := weakAdapter{}
			handler := OpenOrderHandler(&openAdapter)
			mockOrder := new(OpenOrderRequest)
			s, _ := json.Marshal(mockOrder)
			body := bytes.NewBuffer(s)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			handler.ServeHTTP(w, r)

			Expect(w.Code).To(Equal(http.StatusCreated))
			Expect(atomic.LoadInt64(&openAdapter.numOpened)).To(Equal(int64(1)))
		})

		It("should return a StatusBadRequest response for an invalid request", func() {
			openAdapter := weakAdapter{}
			handler := OpenOrderHandler(&openAdapter)
			mockOrder := ""
			s, _ := json.Marshal(mockOrder)
			body := bytes.NewBuffer(s)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			handler.ServeHTTP(w, r)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("cannot decode json into an order or a list of order fragments:"))
			Expect(atomic.LoadInt64(&openAdapter.numOpened)).To(Equal(int64(0)))
		})

		It("should return a StatusBadRequest response for an error in opening orders", func() {
			openAdapter := errAdapter{}
			handler := OpenOrderHandler(&openAdapter)
			mockOrder := new(OpenOrderRequest)
			s, _ := json.Marshal(mockOrder)
			body := bytes.NewBuffer(s)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http://localhost/orders", body)
			handler.ServeHTTP(w, r)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("cannot open order"))
		})
	})

	Context("when canceling orders", func() {

		It("should return a StatusGone response for a valid request", func() {
			cancelAdapter := weakAdapter{}
			handler := CancelOrderHandler(&cancelAdapter)
			cancelOrderID := "vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQ"
			signature := "Td2YBy0MRYPYqqBduRmDsIhTySQUlMhPBM+wnNPWKqq="
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http://localhost/orders/"+cancelOrderID+"?signature="+signature, nil)
			handler.ServeHTTP(w, r)

			Expect(w.Code).To(Equal(http.StatusGone))
			Expect(atomic.LoadInt64(&cancelAdapter.numCanceled)).To(Equal(int64(1)))
		})

		It("should return a StatusBadRequest response for an invalid signature", func() {
			cancelAdapter := weakAdapter{}
			handler := CancelOrderHandler(&cancelAdapter)
			cancelOrderID := "vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQ"
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http://localhost/orders/"+cancelOrderID, nil)
			handler.ServeHTTP(w, r)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("nil signature"))
			Expect(atomic.LoadInt64(&cancelAdapter.numOpened)).To(Equal(int64(0)))
		})

		It("should return a StatusBadRequest response for an empty order id", func() {
			cancelAdapter := weakAdapter{}
			handler := CancelOrderHandler(&cancelAdapter)
			signature := "Td2YBy0MRYPYqqBduRmDsIhTySQUlMhPBM+wnNPWKqq="
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http://localhost/orders/?signature="+signature, nil)
			handler.ServeHTTP(w, r)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("nil id"))
			Expect(atomic.LoadInt64(&cancelAdapter.numOpened)).To(Equal(int64(0)))
		})

		It("should return a StatusBadRequest response for an error in canceling orders", func() {
			cancelAdapter := errAdapter{}
			handler := CancelOrderHandler(&cancelAdapter)
			mockOrder := new(OpenOrderRequest)
			s, _ := json.Marshal(mockOrder)
			body := bytes.NewBuffer(s)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http://localhost/orders", body)
			handler.ServeHTTP(w, r)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("cannot cancel order"))
		})
	})
})