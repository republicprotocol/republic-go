package stackint_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStackint(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stackint Suite")
}
