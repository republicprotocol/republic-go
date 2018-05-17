package orderbook_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOrderbook(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Orderbook Suite")
}
