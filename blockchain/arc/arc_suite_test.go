package arc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestArc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Arc Suite")
}
