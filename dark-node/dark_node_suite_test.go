package node_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoMiner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dark Node Suite")
}
