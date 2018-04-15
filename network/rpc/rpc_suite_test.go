package rpc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGoRPCNetwork(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Go RPC Suite")
}
