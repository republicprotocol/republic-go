package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/randentropy"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/scrypt"
)

// ErrUnsupportedKeystoreVersion is returned when the version of the Keystore
// is unsupported. Currently, only version 3 is supported.
var ErrUnsupportedKeystoreVersion = errors.New("unsupported keystore version")

// ErrUnsupportedKDF is returned when the KDF of an encrypted Keystore is
// unsupported. Currently, only scrypt is supported.
var ErrUnsupportedKDF = errors.New("unsupported kdf")

// ErrPassphraseCannotDecryptKey is returned when the given passphrase cannot
// be used to decrypted a Keystore.
var ErrPassphraseCannotDecryptKey = errors.New("passphrase cannot decrypt key")

// A Keystore stores an EcdsaKey and an RsaKey. It exists primarily to couple
// the keys together to form one unified identity, capable of signing,
// verifying, encrypting, and decrypting. It also exists to expose an easy
// interface for storing/loading keys to/from persistent storagae, optionally
// encrypted.
type Keystore struct {
	ID      uuid.UUID `json:"id"`
	Version string    `json:"version"`

	EcdsaKey `json:"ecdsa"`
	RsaKey   `json:"rsa"`
}

// RandomKeystore returns a new Keystore that stores a randomly generated
// EcdsaKey, RsaKey, and UUID.
func RandomKeystore() (Keystore, error) {
	var err error
	keystore := Keystore{}
	keystore.EcdsaKey, err = RandomEcdsaKey()
	if err != nil {
		return keystore, err
	}
	keystore.RsaKey, err = RandomRsaKey()
	if err != nil {
		return keystore, err
	}
	keystore.ID = uuid.NewRandom()
	keystore.Version = "3"
	return keystore, nil
}

// NewKeystore returns a new Keystore that stores an EcdsaKey and an RsaKey and
// has a randomly generated UUID.
func (keystore *Keystore) NewKeystore(ecdsaKey EcdsaKey, rsaKey RsaKey) Keystore {
	return Keystore{
		ID:       uuid.NewRandom(),
		Version:  "3",
		EcdsaKey: ecdsaKey,
		RsaKey:   rsaKey,
	}
}

// EncryptToJSON will encrypt the EcdsaKey and RsaKey in the Keystore. It
// returns the encrypted Keystore, marshaled as a JSON object.
func (keystore *Keystore) EncryptToJSON(passphrase string, scryptN, scryptP int) ([]byte, error) {
	ecdsaKeyEncrypted, err := keystore.encryptEcdsaKey(passphrase, scryptN, scryptP)
	if err != nil {
		return nil, err
	}
	rsaKeyEncrypted, err := keystore.encryptRsaKey(passphrase, scryptN, scryptP)
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

// DecryptFromJSON will decrypt the EcdsaKey and RsaKey from the JSON object
// into a Keystore.
func (keystore *Keystore) DecryptFromJSON(data []byte, passphrase string) error {
	// Parse the json into a map to fetch the key version
	val := make(map[string]interface{})
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	if version, ok := val["version"].(string); !ok || version != "3" {
		return ErrUnsupportedKeystoreVersion
	}

	keystoreEncrypted := encryptedKeystoreJSONV3{}
	if err := json.Unmarshal(data, &keystoreEncrypted); err != nil {
		return err
	}
	if keystoreEncrypted.EcdsaKey != nil {
		ecdsaDecrypted, err := decryptEcdsaKeyV3(keystoreEncrypted.EcdsaKey, passphrase)
		if err != nil {
			return err
		}
		keystore.EcdsaKey = ecdsaDecrypted
	}
	if keystoreEncrypted.RsaKey != nil {
		rsaDecrypted, err := decryptRsaKeyV3(keystoreEncrypted.RsaKey, passphrase)
		if err != nil {
			return err
		}
		keystore.RsaKey = rsaDecrypted
	}

	keystore.ID = uuid.Parse(keystoreEncrypted.ID)
	keystore.Version = "3"
	return nil
}

// From https://github.com/ethereum/go-ethereum/accounts/keystore
const (
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

	keyHeaderKDF = "scrypt"
	scryptR      = 8
	scryptDKLen  = 32
)

type encryptedKeystoreJSONV3 struct {
	ID       string              `json:"id"`
	Version  string              `json:"version"`
	EcdsaKey *encryptedKeyJSONV3 `json:"ecdsa,omitempty"`
	RsaKey   *encryptedKeyJSONV3 `json:"rsa,omitempty"`
}

// Adapted from https://github.com/ethereum/go-ethereum/accounts/keystore
type encryptedKeyJSONV3 struct {
	Address   *string    `json:"address,omitempty"`
	PublicKey []byte     `json:"publicKey,omitempty"`
	Crypto    cryptoJSON `json:"crypto"`
}

// Adapted from https://github.com/ethereum/go-ethereum/accounts/keystore
type cryptoJSON struct {
	Cipher       string                 `json:"cipher"`
	CipherText   string                 `json:"cipherText"`
	CipherParams cipherparamsJSON       `json:"cipherParams"`
	KDF          string                 `json:"kdf"`
	KDFParams    map[string]interface{} `json:"kdfParams"`
	MAC          string                 `json:"mac"`
}

// From https://github.com/ethereum/go-ethereum/accounts/keystore
type cipherparamsJSON struct {
	IV string `json:"iv"`
}

// Adapted from https://github.com/ethereum/go-ethereum/accounts/keystore
func (keystore *Keystore) encryptEcdsaKey(passphrase string, scryptN, scryptP int) (encryptedKeyJSONV3, error) {

	salt := randentropy.GetEntropyCSPRNG(32)
	keyDerived, err := scrypt.Key([]byte(passphrase), salt, scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		return encryptedKeyJSONV3{}, err
	}
	keyEncrypted := keyDerived[:16]
	keyBytes, err := keystore.EcdsaKey.MarshalJSON()
	if err != nil {
		return encryptedKeyJSONV3{}, err
	}
	iv := randentropy.GetEntropyCSPRNG(aes.BlockSize) // 16
	cipherText, err := aesCTRXOR(keyEncrypted, keyBytes, iv)
	if err != nil {
		return encryptedKeyJSONV3{}, err
	}
	mac := crypto.Keccak256(keyDerived[16:32], cipherText)

	scryptParamsJSON := make(map[string]interface{}, 5)
	scryptParamsJSON["n"] = scryptN
	scryptParamsJSON["r"] = scryptR
	scryptParamsJSON["p"] = scryptP
	scryptParamsJSON["dkLength"] = scryptDKLen
	scryptParamsJSON["salt"] = hex.EncodeToString(salt)
	cipherParamsJSON := cipherparamsJSON{
		IV: hex.EncodeToString(iv),
	}

	addr := keystore.EcdsaKey.Address()
	cryptoStruct := cryptoJSON{
		Cipher:       "aes-128-ctr",
		CipherText:   hex.EncodeToString(cipherText),
		CipherParams: cipherParamsJSON,
		KDF:          keyHeaderKDF,
		KDFParams:    scryptParamsJSON,
		MAC:          hex.EncodeToString(mac),
	}
	return encryptedKeyJSONV3{
		Address: &addr,
		Crypto:  cryptoStruct,
	}, nil
}

// Adapted from https://github.com/ethereum/go-ethereum/accounts/keystore
func (keystore *Keystore) encryptRsaKey(passphrase string, scryptN, scryptP int) (encryptedKeyJSONV3, error) {

	salt := randentropy.GetEntropyCSPRNG(32)
	keyDerived, err := scrypt.Key([]byte(passphrase), salt, scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		return encryptedKeyJSONV3{}, err
	}
	keyEncrypted := keyDerived[:16]
	keyBytes, err := keystore.RsaKey.MarshalJSON()
	if err != nil {
		return encryptedKeyJSONV3{}, err
	}
	iv := randentropy.GetEntropyCSPRNG(aes.BlockSize) // 16
	cipherText, err := aesCTRXOR(keyEncrypted, keyBytes, iv)
	if err != nil {
		return encryptedKeyJSONV3{}, err
	}
	mac := crypto.Keccak256(keyDerived[16:32], cipherText)

	scryptParamsJSON := make(map[string]interface{}, 5)
	scryptParamsJSON["n"] = scryptN
	scryptParamsJSON["r"] = scryptR
	scryptParamsJSON["p"] = scryptP
	scryptParamsJSON["dkLength"] = scryptDKLen
	scryptParamsJSON["salt"] = hex.EncodeToString(salt)
	cipherParamsJSON := cipherparamsJSON{
		IV: hex.EncodeToString(iv),
	}

	publicKey, err := BytesFromRsaPublicKey(&keystore.RsaKey.PublicKey)
	if err != nil {
		return encryptedKeyJSONV3{}, fmt.Errorf("cannot read rsa.PublicKey from bytes: %v", err)
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
		PublicKey: publicKey,
		Crypto:    cryptoStruct,
	}, nil
}

func decryptEcdsaKeyV3(ecdsaKeyEncrypted *encryptedKeyJSONV3, passphrase string) (EcdsaKey, error) {
	ecdsaKeyBytes, err := decryptKeyV3(ecdsaKeyEncrypted, passphrase)
	if err != nil {
		return EcdsaKey{}, err
	}
	ecdsaKey := EcdsaKey{}
	if err := ecdsaKey.UnmarshalJSON(ecdsaKeyBytes); err != nil {
		return ecdsaKey, err
	}
	return ecdsaKey, nil
}

func decryptRsaKeyV3(esaKeyEncrypted *encryptedKeyJSONV3, passphrase string) (RsaKey, error) {
	rsaKeyBytes, err := decryptKeyV3(esaKeyEncrypted, passphrase)
	if err != nil {
		return RsaKey{}, err
	}
	rsaKey := RsaKey{}
	if err := rsaKey.UnmarshalJSON(rsaKeyBytes); err != nil {
		return rsaKey, err
	}
	return rsaKey, nil
}

// Adapted from https://github.com/ethereum/go-ethereum/accounts/keystore
func decryptKeyV3(keyEncrypted *encryptedKeyJSONV3, auth string) (keyBytes []byte, err error) {
	if keyEncrypted.Crypto.Cipher != "aes-128-ctr" {
		return nil, fmt.Errorf("Cipher not supported: %v", keyEncrypted.Crypto.Cipher)
	}

	mac, err := hex.DecodeString(keyEncrypted.Crypto.MAC)
	if err != nil {
		return nil, err
	}

	iv, err := hex.DecodeString(keyEncrypted.Crypto.CipherParams.IV)
	if err != nil {
		return nil, err
	}

	cipherText, err := hex.DecodeString(keyEncrypted.Crypto.CipherText)
	if err != nil {
		return nil, err
	}

	derivedKey, err := getKDFKey(keyEncrypted.Crypto, auth)
	if err != nil {
		return nil, err
	}

	calculatedMAC := crypto.Keccak256(derivedKey[16:32], cipherText)
	if !bytes.Equal(calculatedMAC, mac) {
		return nil, ErrPassphraseCannotDecryptKey
	}

	plainText, err := aesCTRXOR(derivedKey[:16], cipherText, iv)
	if err != nil {
		return nil, err
	}
	return plainText, err
}

// From https://github.com/ethereum/go-ethereum/accounts/keystore
func getKDFKey(cryptoJSON cryptoJSON, auth string) ([]byte, error) {
	authArray := []byte(auth)
	salt, err := hex.DecodeString(cryptoJSON.KDFParams["salt"].(string))
	if err != nil {
		return nil, err
	}
	dkLen := ensureInt(cryptoJSON.KDFParams["dkLength"])

	if cryptoJSON.KDF == keyHeaderKDF {
		n := ensureInt(cryptoJSON.KDFParams["n"])
		r := ensureInt(cryptoJSON.KDFParams["r"])
		p := ensureInt(cryptoJSON.KDFParams["p"])
		return scrypt.Key(authArray, salt, n, r, p, dkLen)

	}
	return nil, ErrUnsupportedKDF
}

// From https://github.com/ethereum/go-ethereum/accounts/keystore
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

// From https://github.com/ethereum/go-ethereum/accounts/keystore
func ensureInt(x interface{}) int {
	res, ok := x.(int)
	if !ok {
		res = int(x.(float64))
	}
	return res
}
