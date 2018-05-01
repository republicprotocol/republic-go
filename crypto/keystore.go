package crypto

// A Keystore stores an EcdsaKey and an RsaKey. It exists primarily to couple
// the keys together to form one unified identity, capable of signing,
// verifying, encrypting, and decrypting. It also exists to expose an easy
// interface for storing/loading keys to/from persistent storagae, optionally
// encrypted.
type Keystore struct {
	EcdsaKey `json:"ecdsa"`
	RsaKey   `json:"rsa"`
}
