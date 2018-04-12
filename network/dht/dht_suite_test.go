package dht_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoDht(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DHT Suite")
}
