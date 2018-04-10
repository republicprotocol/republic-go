package relay_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/relay"
)

var _ = Describe("HTTP handlers", func() {
	Context("when posting orders", func() {

		It("should return 400 for empty bodies", func() {
			r := httptest.NewRequest("POST", "http://localhost/orders", nil)
			w := httptest.NewRecorder()

			handler := PostOrdersHandler(nil, nil)
			err := handler(w, r)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(w.Code).Should(Equal(http.StatusBadRequest))
		})
	})
})
