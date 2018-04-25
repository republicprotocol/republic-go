package smpcer_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/rpc/smpcer"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/rpc/client"
)

var _ = Describe("Smpcer Client", func() {

	Context("NewClient method", func() {

		It("should return a Client object", func() {
			client, _, err := createNewClient("127.0.0.1", "3000")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(client).ShouldNot(BeNil())
		})
	})
})

func createMultiAddrAndConnPool(bind, port string) (identity.MultiAddress, client.ConnPool, error) {
	addr, _, err := identity.NewAddress()
	if err != nil {
		return identity.MultiAddress{}, client.ConnPool{}, err
	}
	multiaddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%s/tcp/%s/republic/%v", bind, port, addr))
	if err != nil {
		return identity.MultiAddress{}, client.ConnPool{}, err
	}
	connPool := client.NewConnPool(5)
	return multiaddress, connPool, nil
}

func createNewClient(bind, port string) (Client, identity.MultiAddress, error) {
	multiaddr, connPool, err := createMultiAddrAndConnPool(bind, port)
	if err != nil {
		return Client{}, identity.MultiAddress{}, err
	}
	crypter := crypto.NewWeakCrypter()
	return NewClient(&crypter, multiaddr, &connPool), multiaddr, nil
}
