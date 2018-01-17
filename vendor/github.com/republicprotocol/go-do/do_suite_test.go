package do_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoAsync(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Concurrency Suite")
}
