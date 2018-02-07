package swarm_test

import (
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

// testMu is used to serialize tests that cannot be run in parallel. This is
// useful for network tests that open ports.
var testMu = new(sync.Mutex)

func TestGoXNetwork(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Swarm Network Suite")
}
