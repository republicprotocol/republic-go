package relay_test

import (
	"net/http"
    "net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/relay"
)

var _ = Describe("HTTP handlers", func() {
	Context("when posting orders", func() {

		It("should return 400 for empty bodies", func() {
    
			r := httptest.NewRequest("POST", "http://localhost/orders", nil)
			w := httptest.NewRecorder()

			var trader, _ = identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f")
			handler := relay.RecoveryHandler(relay.PostOrdersHandler(trader, nil))
			handler.ServeHTTP(w, r)

			// Check the status code is what we expect.
			Ω(w.Code).Should(Equal(http.StatusOK))
		})
	})
})
