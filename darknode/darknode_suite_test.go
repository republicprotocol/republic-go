package darknode_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDarknode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Darknode Suite")
}
