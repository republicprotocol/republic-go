package compute_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoFragment(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Order Compute Suite")
}
