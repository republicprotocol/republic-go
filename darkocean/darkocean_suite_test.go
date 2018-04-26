package darkocean_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDarkocean(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Darkocean Suite")
}
