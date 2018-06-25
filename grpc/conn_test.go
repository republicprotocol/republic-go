package grpc_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/grpc"

	_ "github.com/republicprotocol/republic-go/identity" // initialise the protocol
	"google.golang.org/grpc"
)

var _ = Describe("Connections", func() {

	var server *grpc.Server

	BeforeSuite(func(done Done) {
		defer close(done)

		server = grpc.NewServer()
		listener, err := net.Listen("tcp", "127.0.0.1:3000")
		Expect(err).ShouldNot(HaveOccurred())
		go func() {
			defer GinkgoRecover()
			err := server.Serve(listener)
			Expect(err).ShouldNot(HaveOccurred())
		}()
	})

	AfterSuite(func() {
		server.Stop()
	})

	Context("when dialing", func() {
		It("should return a connection for valid multiaddresses", func() {

			ecdsaKey, err := crypto.RandomEcdsaKey()
			Expect(err).ShouldNot(HaveOccurred())
			multiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/3000/republic/%s", ecdsaKey.Address()))
			Expect(err).ShouldNot(HaveOccurred())

			conn, err := Dial(context.Background(), multiAddr)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conn).ShouldNot(BeNil())

		})
	})

	Context("when backing off", func() {
		It("should retry with a longer wait time until the operation succeeds", func() {
			attempts := 0

			start := time.Now()
			err := Backoff(context.Background(), func() error {
				if attempts < 2 {
					attempts++
					return errors.New("error")
				}
				return nil
			})
			end := time.Now()

			Expect(err).ShouldNot(HaveOccurred())
			Expect(attempts).Should(Equal(2))
			Expect(end.Sub(start).Seconds()).Should(BeNumerically(">=", 2.0))
			Expect(end.Sub(start).Seconds()).Should(BeNumerically("<=", 4.0))
		})

		It("should timeout with an error after the context is done", func() {
			attempts := 0

			start := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			err := Backoff(ctx, func() error {
				attempts++
				return errors.New("error")
			})
			end := time.Now()

			Expect(err).Should(HaveOccurred())
			Expect(attempts).Should(Equal(2))
			Expect(end.Sub(start).Seconds()).Should(BeNumerically(">=", 2.0))
			Expect(end.Sub(start).Seconds()).Should(BeNumerically("<=", 4.0))
		})
	})

})
