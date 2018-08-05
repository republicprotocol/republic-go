package oracle_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOracle(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Oracle Suite")
}
