package swarm_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	. "github.com/republicprotocol/go-swarm"

	"fmt"
	"math/rand"
)

const (
	Number_Of_Peers = 50
	Localhost       = "127.0.0.1"
)

var _ = Describe("Node", func() {

	Describe("Star topology", func() {
		Ω(false).Should(Equal(true))
	})
})
