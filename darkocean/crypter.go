package darkocean

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
)

var ErrInvalidRegistration = errors.New("invalid registration")

const maxCacheSize = 256

type registryCacheEntry struct {
	timestamp    time.Time
	isRegistered bool
}

type publicKeyCacheEntry struct {
	timestamp time.Time
	publicKey rsa.PublicKey
}

// Crypter is an implementation of the crypto.Crypter interface. In addition to
// standard signature verification, the Crypter uses a dnr.DarknodeRegister to
// verify that the signatory is correctly registered to the network. It also
// uses the dnr.DarknodeRegister to lazily acquire the necessary rsa.PublicKeys
// for encryption.
type Crypter struct {
	keystore         crypto.Keystore
	darknodeRegistry dnr.DarknodeRegistry

	registryCache  map[string]registryCacheEntry
	publicKeyCache map[string]publicKeyCacheEntry
}

// NewCrypter returns a new Crypter that uses a crypto.Keystore to identify
// itself when signing and decrypting messages. It uses a dnr.DarknodeRegistry
// to identify others when verifying and encrypting messages.
func NewCrypter(keystore crypto.Keystore, darknodeRegistry dnr.DarknodeRegistry) Crypter {
	return Crypter{
		keystore:         keystore,
		darknodeRegistry: darknodeRegistry,
	}
}

// Sign using the crypto.Keystore that identifies the Crypter.
func (crypter *Crypter) Sign(hasher crypto.Hasher) ([]byte, error) {
	return crypter.keystore.Sign(hasher)
}

func (crypter *Crypter) Verify(hasher crypto.Hasher, signature []byte) error {
	addr, err := crypto.RecoverAddress(hasher, signature)
	if err != nil {
		return err
	}
	if err := crypter.verifyAddress(addr); err != nil {
		return err
	}
	return nil
}

func (crypter *Crypter) Encrypt(addr string, plainText []byte) ([]byte, error) {
	return crypter.encryptToAddress(addr, plainText)
}

func (crypter *Crypter) Decrypt(cipherText []byte) ([]byte, error) {
	return crypter.keystore.Decrypt(cipherText)
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
	if !ok || entry.timestamp.Add(time.Minute).Before(time.Now()) {
		isRegistered, err := crypter.darknodeRegistry.IsRegistered(identity.Address(addr).ID())
		if err != nil {
			return err
		}
		entry = registryCacheEntry{isRegistered: isRegistered}
	}
	entry.timestamp = time.Now()
	crypter.registryCache[addr] = entry

	// Ensure the cache has not exceeded its limit
	if len(crypter.registryCache) > maxCacheSize {
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
	rsaKey := crypto.RsaKey{}
	rsaKey.PublicKey = crypter.publicKeyCache[addr].publicKey
	return rsaKey.Encrypt(plainText)
}

func (crypter *Crypter) updatePublicKeyCache(addr string) error {

	// Update the entry in the cache
	entry, ok := crypter.publicKeyCache[addr]
	if !ok || entry.timestamp.Add(time.Minute).Before(time.Now()) {
		publicKeyBytes, err := crypter.darknodeRegistry.GetPublicKey(identity.Address(addr).ID())
		if err != nil {
			return err
		}
		publicKey, err := crypto.RsaPublicKeyFromBytes(publicKeyBytes)
		if err != nil {
			return err
		}
		entry = publicKeyCacheEntry{publicKey: publicKey}
	}
	entry.timestamp = time.Now()
	crypter.publicKeyCache[addr] = entry

	// Ensure the cache has not exceeded its limit
	if len(crypter.publicKeyCache) > maxCacheSize {
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
