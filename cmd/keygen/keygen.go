package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/republicprotocol/republic-go/crypto"
)

func main() {
	fileName := flag.String("out", "keystore.json", "Output keystore file")
	passphrase := flag.String("passphrase", "", "Passphrase used to encrypt the keystore file")
	flag.Parse()

	keystore, err := crypto.RandomKeystore()
	if err != nil {
		log.Fatalf("cannot generate random keystore: %v", err)
	}

	var keystoreJSON []byte
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
}
