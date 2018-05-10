package adapter_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAdapter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Adapter Suite")
}
