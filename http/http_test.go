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
	return errors.New("cannot cancel order")
}

var _ = Describe("HTTP handlers", func() {

	Context("when opening orders", func() {

		It("should return status 201 for a valid request", func() {

			mockOrder := new(OpenOrderRequest)
			data, err := json.Marshal(mockOrder)
			Expect(err).ShouldNot(HaveOccurred())

			body := bytes.NewBuffer(data)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http://localhost/orders", body)

			adapter := weakAdapter{}
			server := NewServer(&adapter, &adapter)
			server.ServeHTTP(w, r)

			Expect(w.Code).To(Equal(http.StatusCreated))
			Expect(atomic.LoadInt64(&adapter.numOpened)).To(Equal(int64(1)))
		})

		It("should return status 400 for an invalid request", func() {

			mockOrder := ""
			data, err := json.Marshal(mockOrder)
			Expect(err).ShouldNot(HaveOccurred())

			body := bytes.NewBuffer(data)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http://localhost/orders", body)

			adapter := weakAdapter{}
			server := NewServer(&adapter, &adapter)
			server.ServeHTTP(w, r)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(atomic.LoadInt64(&adapter.numOpened)).To(Equal(int64(0)))
		})

		It("should return status 500 for adapter errors", func() {

			mockOrder := new(OpenOrderRequest)
			data, err := json.Marshal(mockOrder)
			Expect(err).ShouldNot(HaveOccurred())

			body := bytes.NewBuffer(data)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http://localhost/orders", body)

			adapter := errAdapter{}
			server := NewServer(&adapter, &adapter)
			server.ServeHTTP(w, r)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("when canceling orders", func() {

		It("should return status 200 for a valid request", func() {

			orderID := "vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQ"
			signature := "Td2YBy0MRYPYqqBduRmDsIhTySQUlMhPBM+wnNPWKqq="
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "http://localhost/orders?id="+orderID+"&signature="+signature, nil)

			adapter := weakAdapter{}
			server := NewServer(&adapter, &adapter)
			server.ServeHTTP(w, r)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(atomic.LoadInt64(&adapter.numCanceled)).To(Equal(int64(1)))
		})

		It("should return status 400 for an invalid signature", func() {

			orderID := "vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQ"
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "http://localhost/orders?id="+orderID, nil)

			adapter := weakAdapter{}
			server := NewServer(&adapter, &adapter)
			server.ServeHTTP(w, r)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(atomic.LoadInt64(&adapter.numOpened)).To(Equal(int64(0)))
		})

		It("should return status 500 for adapter errors", func() {

			orderID := "vrZhWU3VV9LRIriRvuzT9CbVc57wQhbQ"
			signature := "Td2YBy0MRYPYqqBduRmDsIhTySQUlMhPBM+wnNPWKqq="
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "http://localhost/orders?id="+orderID+"&signature="+signature, nil)

			adapter := errAdapter{}
			server := NewServer(&adapter, &adapter)
			server.ServeHTTP(w, r)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
