package x_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGoXNetwork(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "X Network Suite")
}
