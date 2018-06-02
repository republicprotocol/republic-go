package darknode_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"runtime"
	"time"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/darknode"

	"github.com/republicprotocol/republic-go/crypto"
)

var _ = Describe("Crypter", func() {

	var darkpool cal.Darkpool
	var darknodeKeystores map[identity.Address]crypto.Keystore

	var crypter Crypter
	var message []byte

	BeforeEach(func() {
		var err error
		darkpool, darknodeKeystores, err = newMockDarkpool()
		Expect(err).ShouldNot(HaveOccurred())
		keystore, err := crypto.RandomKeystore()
		Expect(err).ShouldNot(HaveOccurred())
		darknodes, err := darkpool.Darknodes()
		Expect(err).ShouldNot(HaveOccurred())
		crypter = NewCrypter(keystore, darkpool, len(darknodes)/2, time.Second)
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
			darknodes, err := darkpool.Darknodes()
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
			darknodes, err := darkpool.Darknodes()
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

type mockDarkpool struct {
	darknodes map[identity.Address]crypto.Keystore
	pods      []cal.Pod
}

func newMockDarkpool() (cal.Darkpool, map[identity.Address]crypto.Keystore, error) {
	darkpool := mockDarkpool{
		darknodes: map[identity.Address]crypto.Keystore{},
		pods:      []cal.Pod{},
	}
	pod := cal.Pod{
		Hash:      [32]byte{},
		Darknodes: []identity.Address{},
	}
	rand.Read(pod.Hash[:])
	for i := 0; i < 6; i++ {
		keystore, err := crypto.RandomKeystore()
		if err != nil {
			return &darkpool, darkpool.darknodes, err
		}
		darkpool.darknodes[identity.Address(keystore.Address())] = keystore
		pod.Darknodes = append(pod.Darknodes, identity.Address(keystore.Address()))
	}
	return &darkpool, darkpool.darknodes, nil
}

func (darkpool *mockDarkpool) Darknodes() (identity.Addresses, error) {
	darknodes := identity.Addresses{}
	for _, pod := range darkpool.pods {
		darknodes = append(darknodes, pod.Darknodes...)
	}
	return darknodes, nil
}

func (darkpool *mockDarkpool) Epoch() (cal.Epoch, error) {
	darknodes, err := darkpool.Darknodes()
	if err != nil {
		return cal.Epoch{}, err
	}
	return cal.Epoch{
		Hash:      [32]byte{},
		Pods:      darkpool.pods,
		Darknodes: darknodes,
	}, nil
}

func (darkpool *mockDarkpool) Pods() ([]cal.Pod, error) {
	return darkpool.pods, nil
}

func (darkpool *mockDarkpool) Pod(addr identity.Address) (cal.Pod, error) {
	for _, pod := range darkpool.pods {
		for _, darknode := range pod.Darknodes {
			if addr == darknode {
				return pod, nil
			}
		}
	}
	return cal.Pod{}, cal.ErrPodNotFound
}

func (darkpool *mockDarkpool) PublicKey(addr identity.Address) (rsa.PublicKey, error) {
	if keystore, ok := darkpool.darknodes[addr]; ok {
		return keystore.RsaKey.PublicKey, nil
	}
	return rsa.PublicKey{}, cal.ErrPublicKeyNotFound
}

func (darkpool *mockDarkpool) IsRegistered(addr identity.Address) (bool, error) {
	_, ok := darkpool.darknodes[addr]
	return ok, nil
}
