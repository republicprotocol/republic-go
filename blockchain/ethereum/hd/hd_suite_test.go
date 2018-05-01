package hd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hd Suite")
}
