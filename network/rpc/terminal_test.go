package rpc_test

import (
	"context"
	"fmt"
	"net"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-atom"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-rpc"
	"google.golang.org/grpc"
)

func (s *mockServer) SendAtom(ctx context.Context, a *rpc.Atom) (*rpc.Atom, error) {
	return a, nil
}

var _ = Describe("Terminals", func() {
	var server *grpc.Server
	var rpcServer mockServer
	var rpcClient mockClient
	var defaultTimeout = time.Second
	var badTargetAddress identity.MultiAddress
	var err error
	var a atom.Atom

	createServe := func() {
		keyPair, err := identity.NewKeyPair()
		Ω(err).ShouldNot(HaveOccurred())
		multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", 3000, keyPair.Address()))
		Ω(err).ShouldNot(HaveOccurred())
		server = grpc.NewServer()
		rpcServer = mockServer{MultiAddress: multiAddress}
		rpc.RegisterTerminalNodeServer(server, &rpcServer)
	}

	createClient := func() {
		keyPair, err := identity.NewKeyPair()
		Ω(err).ShouldNot(HaveOccurred())
		multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", 4000, keyPair.Address()))
		Ω(err).ShouldNot(HaveOccurred())
		rpcClient = mockClient{MultiAddress: multiAddress}
		badTargetAddress, err = identity.NewMultiAddressFromString("/ip4/192.168.0.1/republic/8MHzQ7ZQDvvT8Nqo3HLQQDZvfcHJYB")
		Ω(err).ShouldNot(HaveOccurred())
	}

	createAtom := func() {
		a = atom.Atom{
			Ledger:     atom.Ledger(0),
			LedgerData: []byte{},
			Signature:  []byte{},
		}
	}

	BeforeEach(func() {
		createClient()
		createServe()
		createAtom()
	})

	Context("sending trading atoms", func() {
		keyPair, _ := identity.NewKeyPair()
		to := keyPair.Address()
		from := rpcClient.MultiAddress

		It("should return a valid trading atom", func() {
			lis, err := net.Listen("tcp", ":3000")
			Ω(err).ShouldNot(HaveOccurred())
			go func(server *grpc.Server) {
				defer GinkgoRecover()
				Ω(server.Serve(lis)).ShouldNot(HaveOccurred())
			}(server)
			defer server.Stop()

			err = rpc.SendAtomToTarget(rpcServer.MultiAddress, to, from, a, defaultTimeout)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should return an error for bad multi-addresses", func() {
			err = rpc.SendAtomToTarget(badTargetAddress, to, from, a, defaultTimeout)
			Ω(err).Should(HaveOccurred())
		})

		It("should return a timeout error when there is no response within the timeout duration", func() {
			err := rpc.SendAtomToTarget(rpcServer.MultiAddress, to, from, a, defaultTimeout)
			Ω(err).Should(HaveOccurred())
		})
	})
})
