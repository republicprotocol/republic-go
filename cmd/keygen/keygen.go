package main

import (
	"crypto/rand"
	"flag"
	"io/ioutil"
	"log"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func main() {
	fileName := flag.String("out", "keystore.json", "Output keystore file")
	passphrase := flag.String("passphrase", "", "Passphrase used to encrypt the keystore file")
	flag.Parse()

	if *passphrase == "" {
		log.Fatal("cannot encrypt key: empty passphrase")
	}

	key := keystore.NewKeyForDirectICAP(rand.Reader)
	keyJSON, err := keystore.EncryptKey(key, *passphrase, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		log.Fatal("cannot encrypt key:", err)
	}

	if err := ioutil.WriteFile(*fileName, keyJSON, 0600); err != nil {
		log.Fatal("cannot create keystore file:", err)
	}
}
