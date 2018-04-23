package contracts_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestContracts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Contracts Suite")
}
