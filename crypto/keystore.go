package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/randentropy"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

// ErrUnsupportedKeystoreVersion is returned when the version of the Keystore
// is unsupported. Currently, only version 3 is supported.
var ErrUnsupportedKeystoreVersion = errors.New("unsupported keystore version")

// A Keystore stores an EcdsaKey and an RsaKey. It exists primarily to couple
// the keys together to form one unified identity, capable of signing,
// verifying, encrypting, and decrypting. It also exists to expose an easy
// interface for storing/loading keys to/from persistent storagae, optionally
// encrypted.
type Keystore struct {
	ID      uuid.UUID `json:"id"`
	Version int       `json:"version"`

	EcdsaKey `json:"ecdsa"`
	RsaKey   `json:"rsa"`
}

// NewKeystore returns a new Keystore that stores an EcdsaKey and an RsaKey and
// has a randomly generated UUID.
func (keystore *Keystore) NewKeystore(ecdsaKey EcdsaKey, rsaKey RsaKey) Keystore {
	return Keystore{
		ID:       uuid.NewRandom(),
		Version:  3,
		EcdsaKey: ecdsaKey,
		RsaKey:   rsaKey,
	}
}

// Encrypt the EcdsaKey and RsaKey in the Keystore. Return the encrypted
// Keystore marshaled as a JSON object.
func (keystore *Keystore) Encrypt(passphrase string, scryptN, scryptP int) ([]byte, error) {
	ecdsaKeyEncrypted, err := keystore.encryptEcdsaKey(passphrase, scryptN, scryptP)
	if err != nil {
		return nil, err
	}
	rsaKeyEncrypted, err := []byte{}, nil // keystore.encryptRsaKey(passphrase, scryptN, scryptP)
	if err != nil {
		return nil, err
	}
	keystoreEncrypted := map[string]interface{}{
		"id":      keystore.ID.String(),
		"version": keystore.Version,
		"ecdsa":   ecdsaKeyEncrypted,
		"rsa":     rsaKeyEncrypted,
	}
	return json.Marshal(keystoreEncrypted)
}

// From: https://github.com/ethereum/go-ethereum/accounts/keystore
func (keystore *Keystore) encryptEcdsaKey(passphrase string, scryptN, scryptP int) (encryptedKeyJSONV3, error) {
	salt := randentropy.GetEntropyCSPRNG(32)
	derivedKey, err := scrypt.Key([]byte(passphrase), salt, scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		return encryptedKeyJSONV3{}, err
	}
	encryptKey := derivedKey[:16]

	// FIXME: Marshal into JSON data
	keyBytes := math.PaddedBigBytes(keystore.EcdsaKey.PrivateKey.D, 32)

	iv := randentropy.GetEntropyCSPRNG(aes.BlockSize) // 16
	cipherText, err := aesCTRXOR(encryptKey, keyBytes, iv)
	if err != nil {
		return encryptedKeyJSONV3{}, err
	}
	mac := crypto.Keccak256(derivedKey[16:32], cipherText)

	scryptParamsJSON := make(map[string]interface{}, 5)
	scryptParamsJSON["n"] = scryptN
	scryptParamsJSON["r"] = scryptR
	scryptParamsJSON["p"] = scryptP
	scryptParamsJSON["dkLength"] = scryptDKLen
	scryptParamsJSON["salt"] = hex.EncodeToString(salt)

	cipherParamsJSON := cipherparamsJSON{
		IV: hex.EncodeToString(iv),
	}

	cryptoStruct := cryptoJSON{
		Cipher:       "aes-128-ctr",
		CipherText:   hex.EncodeToString(cipherText),
		CipherParams: cipherParamsJSON,
		KDF:          keyHeaderKDF,
		KDFParams:    scryptParamsJSON,
		MAC:          hex.EncodeToString(mac),
	}
	return encryptedKeyJSONV3{
		Address: hex.EncodeToString([]byte(keystore.EcdsaKey.Address())),
		Crypto:  cryptoStruct,
	}, nil
}

// From: https://github.com/ethereum/go-ethereum/accounts/keystore
func decryptEcdsaKey(keyjson []byte, auth string) (EcdsaKey, error) {
	// Parse the json into a simple map to fetch the key version
	m := make(map[string]interface{})
	if err := json.Unmarshal(keyjson, &m); err != nil {
		return EcdsaKey{}, err
	}
	// Depending on the version try to parse one way or another
	var (
		keyBytes, keyID []byte
		err             error
	)
	if version, ok := m["version"].(string); ok && version == "3" {
		k := new(encryptedKeyJSONV3)
		if err := json.Unmarshal(keyjson, k); err != nil {
			return EcdsaKey{}, err
		}
		keyBytes, keyID, err = decryptKeyV3(k, auth)
	} else {
		return EcdsaKey{}, ErrUnsupportedKeystoreVersion
	}
	// Handle any decryption errors and return the key
	if err != nil {
		return EcdsaKey{}, err
	}
	key := crypto.ToECDSAUnsafe(keyBytes)

	return NewEcdsaKey(key), nil
}

// Adapted from https://github.com/ethereum/go-ethereum/accounts/keystore
func (keystore *Keystore) encryptRsaKey(passphrase string, scryptN, scryptP int) (encryptedKeyJSONV3, error) {
	salt := randentropy.GetEntropyCSPRNG(32)
	derivedKey, err := scrypt.Key([]byte(passphrase), salt, scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		return encryptedKeyJSONV3{}, err
	}
	encryptKey := derivedKey[:16]
	// FIXME: Marshal into JSON data
	keyBytes := math.PaddedBigBytes(keystore.RsaKey.PrivateKey.D, 32)

	iv := randentropy.GetEntropyCSPRNG(aes.BlockSize) // 16
	cipherText, err := aesCTRXOR(encryptKey, keyBytes, iv)
	if err != nil {
		return encryptedKeyJSONV3{}, err
	}
	mac := crypto.Keccak256(derivedKey[16:32], cipherText)

	scryptParamsJSON := make(map[string]interface{}, 5)
	scryptParamsJSON["n"] = scryptN
	scryptParamsJSON["r"] = scryptR
	scryptParamsJSON["p"] = scryptP
	scryptParamsJSON["dkLength"] = scryptDKLen
	scryptParamsJSON["salt"] = hex.EncodeToString(salt)

	cipherParamsJSON := cipherparamsJSON{
		IV: hex.EncodeToString(iv),
	}

	cryptoStruct := cryptoJSON{
		Cipher:       "aes-128-ctr",
		CipherText:   hex.EncodeToString(cipherText),
		CipherParams: cipherParamsJSON,
		KDF:          keyHeaderKDF,
		KDFParams:    scryptParamsJSON,
		MAC:          hex.EncodeToString(mac),
	}
	return encryptedKeyJSONV3{
		Address: hex.EncodeToString([]byte(keystore.EcdsaKey.Address())),
		Crypto:  cryptoStruct,
	}, nil
}

// Adapted from https://github.com/ethereum/go-ethereum/accounts/keystore
func decryptRsaKey(keyjson []byte, auth string) (EcdsaKey, error) {
	// Parse the json into a simple map to fetch the key version
	m := make(map[string]interface{})
	if err := json.Unmarshal(keyjson, &m); err != nil {
		return EcdsaKey{}, err
	}
	// Depending on the version try to parse one way or another
	var (
		keyBytes, keyID []byte
		err             error
	)
	if version, ok := m["version"].(string); ok && version == "3" {
		k := new(encryptedKeyJSONV3)
		if err := json.Unmarshal(keyjson, k); err != nil {
			return EcdsaKey{}, err
		}
		keyBytes, keyID, err = decryptKeyV3(k, auth)
	} else {
		return EcdsaKey{}, ErrUnsupportedKeystoreVersion
	}
	// Handle any decryption errors and return the key
	if err != nil {
		return EcdsaKey{}, err
	}
	key := crypto.ToECDSAUnsafe(keyBytes)

	return NewEcdsaKey(key), nil
}

// From: https://github.com/ethereum/go-ethereum/accounts/keystore
func decryptKeyV3(keyProtected *encryptedKeyJSONV3, auth string) (keyBytes []byte, keyId []byte, err error) {
	if keyProtected.Version != 3 {
		return nil, nil, fmt.Errorf("Version not supported: %v", keyProtected.Version)
	}

	if keyProtected.Crypto.Cipher != "aes-128-ctr" {
		return nil, nil, fmt.Errorf("Cipher not supported: %v", keyProtected.Crypto.Cipher)
	}

	keyId = uuid.Parse(keyProtected.Id)
	mac, err := hex.DecodeString(keyProtected.Crypto.MAC)
	if err != nil {
		return nil, nil, err
	}

	iv, err := hex.DecodeString(keyProtected.Crypto.CipherParams.IV)
	if err != nil {
		return nil, nil, err
	}

	cipherText, err := hex.DecodeString(keyProtected.Crypto.CipherText)
	if err != nil {
		return nil, nil, err
	}

	derivedKey, err := getKDFKey(keyProtected.Crypto, auth)
	if err != nil {
		return nil, nil, err
	}

	calculatedMAC := crypto.Keccak256(derivedKey[16:32], cipherText)
	if !bytes.Equal(calculatedMAC, mac) {
		return nil, nil, ErrDecrypt
	}

	plainText, err := aesCTRXOR(derivedKey[:16], cipherText, iv)
	if err != nil {
		return nil, nil, err
	}
	return plainText, keyId, err
}

// From: https://github.com/ethereum/go-ethereum/accounts/keystore
func getKDFKey(cryptoJSON cryptoJSON, auth string) ([]byte, error) {
	authArray := []byte(auth)
	salt, err := hex.DecodeString(cryptoJSON.KDFParams["salt"].(string))
	if err != nil {
		return nil, err
	}
	dkLen := ensureInt(cryptoJSON.KDFParams["dklen"])

	if cryptoJSON.KDF == keyHeaderKDF {
		n := ensureInt(cryptoJSON.KDFParams["n"])
		r := ensureInt(cryptoJSON.KDFParams["r"])
		p := ensureInt(cryptoJSON.KDFParams["p"])
		return scrypt.Key(authArray, salt, n, r, p, dkLen)

	} else if cryptoJSON.KDF == "pbkdf2" {
		c := ensureInt(cryptoJSON.KDFParams["c"])
		prf := cryptoJSON.KDFParams["prf"].(string)
		if prf != "hmac-sha256" {
			return nil, fmt.Errorf("Unsupported PBKDF2 PRF: %s", prf)
		}
		key := pbkdf2.Key(authArray, salt, c, dkLen, sha256.New)
		return key, nil
	}

	return nil, fmt.Errorf("Unsupported KDF: %s", cryptoJSON.KDF)
}

// From: https://github.com/ethereum/go-ethereum/accounts/keystore
func aesCTRXOR(key, inText, iv []byte) ([]byte, error) {
	// AES-128 is selected due to size of encryptKey.
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(aesBlock, iv)
	outText := make([]byte, len(inText))
	stream.XORKeyStream(outText, inText)
	return outText, err
}

// From: https://github.com/ethereum/go-ethereum/accounts/keystore
func aesCBCDecrypt(key, cipherText, iv []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decrypter := cipher.NewCBCDecrypter(aesBlock, iv)
	paddedPlaintext := make([]byte, len(cipherText))
	decrypter.CryptBlocks(paddedPlaintext, cipherText)
	plaintext := pkcs7Unpad(paddedPlaintext)
	if plaintext == nil {
		return nil, ErrDecrypt
	}
	return plaintext, err
}

// From https://github.com/ethereum/go-ethereum/accounts/keystore
// From https://leanpub.com/gocrypto/read#leanpub-auto-block-cipher-modes
func pkcs7Unpad(in []byte) []byte {
	if len(in) == 0 {
		return nil
	}

	padding := in[len(in)-1]
	if int(padding) > len(in) || padding > aes.BlockSize {
		return nil
	} else if padding == 0 {
		return nil
	}

	for i := len(in) - 1; i > len(in)-int(padding)-1; i-- {
		if in[i] != padding {
			return nil
		}
	}
	return in[:len(in)-int(padding)]
}

// From: https://github.com/ethereum/go-ethereum/accounts/keystore
func ensureInt(x interface{}) int {
	res, ok := x.(int)
	if !ok {
		res = int(x.(float64))
	}
	return res
}

const (
	keyHeaderKDF = "scrypt"

	// StandardScryptN is the N parameter of Scrypt encryption algorithm, using 256MB
	// memory and taking approximately 1s CPU time on a modern processor.
	StandardScryptN = 1 << 18

	// StandardScryptP is the P parameter of Scrypt encryption algorithm, using 256MB
	// memory and taking approximately 1s CPU time on a modern processor.
	StandardScryptP = 1

	// LightScryptN is the N parameter of Scrypt encryption algorithm, using 4MB
	// memory and taking approximately 100ms CPU time on a modern processor.
	LightScryptN = 1 << 12

	// LightScryptP is the P parameter of Scrypt encryption algorithm, using 4MB
	// memory and taking approximately 100ms CPU time on a modern processor.
	LightScryptP = 6

	scryptR     = 8
	scryptDKLen = 32
)

type encryptedKeyJSONV3 struct {
	Address   *string    `json:"address,omitempty"`
	PublicKey *string    `json:"publicKey,omitempty"`
	Crypto    cryptoJSON `json:"crypto"`
}

type cryptoJSON struct {
	Cipher       string                 `json:"cipher"`
	CipherText   string                 `json:"cipherText"`
	CipherParams cipherparamsJSON       `json:"cipherParams"`
	KDF          string                 `json:"kdf"`
	KDFParams    map[string]interface{} `json:"kdfParams"`
	MAC          string                 `json:"mac"`
}

type cipherparamsJSON struct {
	IV string `json:"iv"`
}
