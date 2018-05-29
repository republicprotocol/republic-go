package smpc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSMPC(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "sMPC's Test Suite")
}
