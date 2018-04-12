package dark_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDark(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dark Suite")
}
