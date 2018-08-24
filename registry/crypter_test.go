package registry_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"runtime"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/registry"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
)

var _ = Describe("Crypter", func() {

	var contract *registryBinder
	var darknodeKeystores map[identity.Address]crypto.Keystore

	var crypter Crypter
	var message []byte

	BeforeEach(func() {
		var err error
		contract, darknodeKeystores, err = newRegistryBinder()
		Expect(err).ShouldNot(HaveOccurred())
		keystore, err := crypto.RandomKeystore()
		Expect(err).ShouldNot(HaveOccurred())
		darknodes, err := contract.Darknodes()
		Expect(err).ShouldNot(HaveOccurred())
		crypter = NewCrypter(keystore, contract, len(darknodes)/2, time.Second)
		message = []byte("REN")
	})

	Context("when signing", func() {

		It("should produce valid signatures", func() {
			signature, err := crypter.Sign(crypto.Keccak256(message))
			Expect(err).ShouldNot(HaveOccurred())
			addr, err := crypto.RecoverAddress(crypto.Keccak256(message), signature)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(addr).Should(Equal(crypter.Keystore().Address()))
		})

	})

	Context("when verifying signatures", func() {

		It("should return an error for unregistered addresses", func() {
			dispatch.CoForAll(runtime.NumCPU(), func(i int) {
				keystore, err := crypto.RandomKeystore()
				Expect(err).ShouldNot(HaveOccurred())

				signature, err := keystore.Sign(crypto.Keccak256(message))
				Expect(err).ShouldNot(HaveOccurred())
				err = crypter.Verify(crypto.Keccak256(message), signature)
				Expect(err).Should(HaveOccurred())
			})
		})

		It("should not return an error for registered addresses", func() {
			darknodes, err := contract.Darknodes()
			Expect(err).ShouldNot(HaveOccurred())
			for _, darknode := range darknodes {
				keystore := darknodeKeystores[darknode]
				signature, err := keystore.Sign(crypto.Keccak256(message))
				Expect(err).ShouldNot(HaveOccurred())
				err = crypter.Verify(crypto.Keccak256(message), signature)
				Expect(err).ShouldNot(HaveOccurred())
			}
		})

	})

	Context("when encrypting", func() {

		It("should encrypt messages for registered addresses", func() {
			darknodes, err := contract.Darknodes()
			Expect(err).ShouldNot(HaveOccurred())
			for _, darknode := range darknodes {
				keystore := darknodeKeystores[darknode]
				cipherText, err := crypter.Encrypt(darknode.String(), message[:])
				Expect(err).ShouldNot(HaveOccurred())

				plainText, err := keystore.Decrypt(cipherText)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(bytes.Equal(plainText, message[:])).Should(BeTrue())
			}
		})

		It("should not encrypt messages for unregistered addresses", func() {
			dispatch.CoForAll(runtime.NumCPU(), func(i int) {
				keystore, err := crypto.RandomKeystore()
				Expect(err).ShouldNot(HaveOccurred())

				_, err = crypter.Encrypt(keystore.Address(), message[:])
				Expect(err).Should(HaveOccurred())
			})
		})

	})

	Context("when decrypting", func() {

		It("should produce the original plain text", func() {
			cipherText, err := crypter.Keystore().Encrypt(message[:])
			Expect(err).ShouldNot(HaveOccurred())
			plainText, err := crypter.Decrypt(cipherText)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(bytes.Equal(plainText, message[:])).Should(BeTrue())
		})

	})
})

// ErrPublicKeyNotFound is returned when an rsa.PublicKey cannot be found for a
// given identity.Address. This happens when an identity.Address is not registered in
// the current Epoch.
var ErrPublicKeyNotFound = errors.New("public key not found")

type registryBinder struct {
	darknodes map[identity.Address]crypto.Keystore
	pods      []Pod
}

func newRegistryBinder() (*registryBinder, map[identity.Address]crypto.Keystore, error) {
	binder := registryBinder{
		darknodes: map[identity.Address]crypto.Keystore{},
		pods:      []Pod{},
	}
	pod := Pod{
		Hash:      [32]byte{},
		Darknodes: []identity.Address{},
	}
	rand.Read(pod.Hash[:])
	for i := 0; i < 6; i++ {
		keystore, err := crypto.RandomKeystore()
		if err != nil {
			return &binder, binder.darknodes, err
		}
		binder.darknodes[identity.Address(keystore.Address())] = keystore
		pod.Darknodes = append(pod.Darknodes, identity.Address(keystore.Address()))
	}
	return &binder, binder.darknodes, nil
}

func (binder *registryBinder) Darknodes() (identity.Addresses, error) {
	darknodes := identity.Addresses{}
	for _, pod := range binder.pods {
		darknodes = append(darknodes, pod.Darknodes...)
	}
	return darknodes, nil
}

func (binder *registryBinder) PublicKey(addr identity.Address) (rsa.PublicKey, error) {
	if keystore, ok := binder.darknodes[addr]; ok {
		return keystore.RsaKey.PublicKey, nil
	}
	return rsa.PublicKey{}, ErrPublicKeyNotFound
}

func (binder *registryBinder) IsRegistered(addr identity.Address) (bool, error) {
	_, ok := binder.darknodes[addr]
	return ok, nil
}
