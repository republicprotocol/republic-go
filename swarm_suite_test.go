package x_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGoX(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "X network test Suite")
}
