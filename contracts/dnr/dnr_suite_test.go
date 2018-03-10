package dnr_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDarkNodeRegistrar(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dark Node Registrar Suite")
}
