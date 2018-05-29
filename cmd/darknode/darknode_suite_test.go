package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDarknodeCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Darknode Cmd Suite")
}
