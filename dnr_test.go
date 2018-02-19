package dnr_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-atom/ethereum"
	. "github.com/republicprotocol/go-dark-node-registrar"

	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/go-identity"
)

var ether = big.NewInt(1000000000000000000)

const key = `{"version":3,"id":"7844982f-abe7-4690-8c15-34f75f847c66","address":"db205ea9d35d8c01652263d58351af75cfbcbf07","Crypto":{"ciphertext":"378dce3c1279b36b071e1c7e2540ac1271581bff0bbe36b94f919cb73c491d3a","cipherparams":{"iv":"2eb92da55cc2aa62b7ffddba891f5d35"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"80d3341678f83a14024ba9c3edab072e6bd2eea6aa0fbc9e0a33bae27ffa3d6d","n":8192,"r":8,"p":1},"mac":"3d07502ea6cd6b96a508138d8b8cd2e46c3966240ff276ce288059ba4235cb0d"}}`

var _ = Describe("Dark Node Registrar", func() {

	auth, err := bind.NewTransactor(strings.NewReader(key), "password1")
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}

	client := ethereum.Ropsten("https://ropsten.infura.io/")

	contractAddress := common.HexToAddress("0x32Dad9E9Fe2A3eA2C2c643675A7d2A56814F554f")

	UserConnection := NewDarkNodeRegistrar(context.Background(), client, auth, &bind.CallOpts{}, contractAddress, nil)

	keyPair, err := identity.NewKeyPair()
	Ω(err).Should(BeNil())
	publicKey := append(keyPair.PublicKey.X.Bytes(), keyPair.PublicKey.Y.Bytes()...)
	darkNodeID := [20]byte(keyPair.ID())
	It("Can register a dark node", func() {
		err := UserConnection.Register(darkNodeID, publicKey)
		Ω(err).Should(BeNil())
	})

	It("Can get bond of a registered dark node", func() {
		err := UserConnection.GetBond(darkNodeID)
		Ω(err).Should(BeNil())
	})

	It("Can check if a dark node is registered", func() {
		status, err := UserConnection.IsDarkNodeRegistered(darkNodeID)
		Ω(err).Should(BeNil())
		Ω(status).Should(BeTrue())
	})

	It("Can get the current epoch", func() {

	})

	It("Can get the commitment of a dark node", func() {

	})

	It("Can get the owner of a dark node", func() {

	})

	It("Can get the public key of a dark node", func() {

	})

	It("Can get the xing overlay network", func() {
		XingOverlay, err := UserConnection.GetXingOverlay()
		Ω(err).Should(BeNil())
		Ω(status).Should(BeTrue())
	})

	It("Can get the minimum bond", func() {

	})

	It("Can get the minimum epoch interval", func() {

	})

	It("Can get the  pending refunds", func() {

	})

	It("Can deregister a dark node", func() {

	})

	It("Can get refund", func() {

	})

})
