package http_test

import (
	"bytes"
	netHttp "net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/http"

	"github.com/republicprotocol/republic-go/http/adapter"
	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Status handler", func() {

	sendRequestAndAssertStatusCode := func(isFaulty bool, statusCode int) {

		body := bytes.NewBuffer([]byte{})
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://localhost/status", body)

		reader := testutils.NewMockReader(isFaulty)
		statusAdapter := adapter.NewStatusAdapter(&reader)

		server := NewStatusServer(statusAdapter)
		server.ServeHTTP(w, r)

		Expect(w.Code).To(Equal(statusCode))
	}

	Context("when status returns an error", func() {

		It("should return a 400 (StatusBadRequest) status code", func() {

			sendRequestAndAssertStatusCode(true, netHttp.StatusBadRequest)
		})
	})

	Context("when status does not return an error", func() {

		It("should return a 200 (StatusOK) status code", func() {

			sendRequestAndAssertStatusCode(false, netHttp.StatusOK)
		})
	})
})
