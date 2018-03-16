package consensus_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConsensus(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Consensus Suite")
}
