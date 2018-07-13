package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/republicprotocol/republic-go/crypto"
)

func main() {
	from := flag.String("from", "", "Input Ethereum keystore that will be converted")
	fileName := flag.String("out", "keystore.json", "Output keystore file")
	passphrase := flag.String("passphrase", "", "Passphrase used to encrypt the keystore file")
	flag.Parse()

	var keystoreJSON []byte
	if *from == "" {
		keystore, err := crypto.RandomKeystore()
		if err != nil {
			log.Fatalf("cannot generate random keystore: %v", err)
		}
		if *passphrase == "" {
			keystoreJSON, err = json.MarshalIndent(keystore, "", "  ")
		} else {
			keystoreJSON, err = keystore.EncryptToJSON(*passphrase, crypto.StandardScryptN, crypto.StandardScryptP)
		}
		if err != nil {
			log.Fatalf("cannot marshal keystore: %v", err)
		}

		if err := ioutil.WriteFile(*fileName, keystoreJSON, 0640); err != nil {
			log.Fatal("cannot write to keystore file:", err)
		}
		return
	}

	if *passphrase == "" {
		log.Fatalf("cannot read keystore from %v in plain-text", *from)
	}

	file, err := os.Open(*from)
	if err != nil {
		log.Fatalf("cannot open keystore from %v: %v", *from, err)
	}
	defer file.Close()
	keyData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read keystore from %v: %v", *from, err)
	}

	key, err := keystore.DecryptKey(keyData, *passphrase)
	if err != nil {
		log.Fatalf("cannot decrypt keystore from %v: %v", *from, err)
	}
	ecdsaKey := crypto.NewEcdsaKey(key.PrivateKey)
	keystore := crypto.Keystore{
		EcdsaKey: ecdsaKey,
	}
	keystoreJSON, err = keystore.EncryptToJSON(*passphrase, crypto.StandardScryptN, crypto.StandardScryptP)
	if err != nil {
		log.Fatalf("cannot marshal keystore: %v", err)
	}

	if err := ioutil.WriteFile(*fileName, keystoreJSON, 0640); err != nil {
		log.Fatal("cannot write to keystore file:", err)
	}
	return
}
