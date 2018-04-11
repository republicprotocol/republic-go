package smpc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSmpc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Smpc Suite")
}
