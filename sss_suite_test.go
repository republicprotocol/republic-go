package sss_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSss(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shamir's Secret Sharing Suite")
}
