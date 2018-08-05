package registry

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
)

// ErrInvalidRegistration is returned when an address is not registered in the
// DarknodeRegsitry. It is possible that the address recently registered, but
// the Crypter has already cached it as unregistered. In these cases, the cache
// will be updated periodically, so a secondary attempt can be made slightly
// later.
var ErrInvalidRegistration = errors.New("invalid registration")

// ErrPodNotFound is returned when a Pod cannot be found for a given
// identity.Address. This happens when an identity.Address is not registered in
// the current Epoch.
var ErrPodNotFound = errors.New("pod not found")

// Crypter is an implementation of the crypto.Crypter interface. In addition to
// standard signature verification, the Crypter uses a cal.Darkpool to verify
// that the signatory is correctly registered to the network. It also uses the
// cal.Darkpool to lazily acquire the necessary rsa.PublicKeys for encryption.
// The cache will be updated periodically, to ensure up-to-date information.
type Crypter struct {
	keystore crypto.Keystore
	contract ContractBinder

	registryCacheMu *sync.Mutex
	registryCache   map[string]registryCacheEntry

	publicKeyCacheMu *sync.Mutex
	publicKeyCache   map[string]publicKeyCacheEntry

	cacheLimit        int
	cacheUpdatePeriod time.Duration
}

// NewCrypter returns a new Crypter that uses a crypto.Keystore to identify
// itself when signing and decrypting messages. It uses a cal.Darkpool to
// identify others when verifying and encrypting messages.
func NewCrypter(keystore crypto.Keystore, contract ContractBinder, cacheLimit int, cacheUpdatePeriod time.Duration) Crypter {
	return Crypter{
		keystore:          keystore,
		contract:          contract,
		registryCacheMu:   new(sync.Mutex),
		registryCache:     map[string]registryCacheEntry{},
		publicKeyCacheMu:  new(sync.Mutex),
		publicKeyCache:    map[string]publicKeyCacheEntry{},
		cacheLimit:        cacheLimit,
		cacheUpdatePeriod: cacheUpdatePeriod,
	}
}

// Sign using the crypto.Keystore that identifies the Crypter.
func (crypter *Crypter) Sign(data []byte) ([]byte, error) {
	return crypter.keystore.Sign(data)
}

// Verify a signature and ensure that the signatory is a registered Darknode.
func (crypter *Crypter) Verify(data []byte, signature []byte) error {
	addr, err := crypto.RecoverAddress(data, signature)
	if err != nil {
		return err
	}
	crypter.registryCacheMu.Lock()
	defer crypter.registryCacheMu.Unlock()
	return crypter.verifyAddress(addr)
}

// Encrypt plain text so that is can be securely sent to a specific address.
// The address will be used to lookup the required rsa.PublicKey in the
// DarknodeRegistry. The address registration is verified before encryption is
// attempted. Returns the cipher text, or an error.
func (crypter *Crypter) Encrypt(addr string, plainText []byte) ([]byte, error) {
	crypter.registryCacheMu.Lock()
	crypter.publicKeyCacheMu.Lock()
	defer crypter.registryCacheMu.Unlock()
	defer crypter.publicKeyCacheMu.Unlock()
	if err := crypter.verifyAddress(addr); err != nil {
		return nil, fmt.Errorf("cannot verify address %v", err)
	}
	return crypter.encryptToAddress(addr, plainText)
}

// Decrypt a cipher text that was sent to the identity defined by the
// crypto.Keystore in the Crypter. Returns the plain text, or an error.
func (crypter *Crypter) Decrypt(cipherText []byte) ([]byte, error) {
	return crypter.keystore.Decrypt(cipherText)
}

// Keystore used to identify the Crypter.
func (crypter *Crypter) Keystore() *crypto.Keystore {
	return &crypter.keystore
}

type registryCacheEntry struct {
	timestamp    time.Time
	isRegistered bool
}

type publicKeyCacheEntry struct {
	timestamp time.Time
	publicKey rsa.PublicKey
}

func (crypter *Crypter) verifyAddress(addr string) error {
	if err := crypter.updateRegistryCache(addr); err != nil {
		return err
	}
	if entry, ok := crypter.registryCache[addr]; ok && entry.isRegistered {
		return nil
	}
	return ErrInvalidRegistration
}

func (crypter *Crypter) updateRegistryCache(addr string) error {
	// Update the entry in the cache
	entry, ok := crypter.registryCache[addr]
	if !ok || entry.timestamp.Add(crypter.cacheUpdatePeriod).Before(time.Now()) {
		isRegistered, err := crypter.contract.IsRegistered(identity.Address(addr))
		if err != nil {
			return err
		}
		entry = registryCacheEntry{isRegistered: isRegistered}
	}
	entry.timestamp = time.Now()
	crypter.registryCache[addr] = entry

	// Ensure the cache has not exceeded its limit
	if len(crypter.registryCache) > crypter.cacheLimit {
		var oldest time.Time
		var oldestK string
		for k := range crypter.registryCache {
			if oldestK == "" || crypter.registryCache[k].timestamp.Before(oldest) {
				oldest = crypter.registryCache[k].timestamp
				oldestK = k
			}
		}
		delete(crypter.registryCache, oldestK)
	}
	return nil
}

func (crypter *Crypter) encryptToAddress(addr string, plainText []byte) ([]byte, error) {
	if err := crypter.updatePublicKeyCache(addr); err != nil {
		return nil, err
	}

	rsaKey := crypto.RsaKey{PrivateKey: &rsa.PrivateKey{}}
	rsaKey.PublicKey = crypter.publicKeyCache[addr].publicKey
	return rsaKey.Encrypt(plainText)
}

func (crypter *Crypter) updatePublicKeyCache(addr string) error {

	// Update the entry in the cache
	entry, ok := crypter.publicKeyCache[addr]
	if !ok || entry.timestamp.Add(crypter.cacheUpdatePeriod).Before(time.Now()) {
		publicKey, err := crypter.contract.PublicKey(identity.Address(addr))
		if err != nil {
			return err
		}
		entry = publicKeyCacheEntry{publicKey: publicKey}
	}
	entry.timestamp = time.Now()
	crypter.publicKeyCache[addr] = entry

	// Ensure the cache has not exceeded its limit
	if len(crypter.publicKeyCache) > crypter.cacheLimit {
		var oldest time.Time
		var oldestK string
		for k := range crypter.publicKeyCache {
			if oldestK == "" || crypter.publicKeyCache[k].timestamp.Before(oldest) {
				oldest = crypter.publicKeyCache[k].timestamp
				oldestK = k
			}
		}
		delete(crypter.publicKeyCache, oldestK)
	}
	return nil
}
