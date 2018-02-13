package go_eth_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoEth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoEth Suite")
}
