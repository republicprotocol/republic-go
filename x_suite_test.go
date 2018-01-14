package x_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoX(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "X Suite")
}
