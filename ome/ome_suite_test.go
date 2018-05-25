package ome_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOrderbook(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Order Matching Engine Suite")
}
