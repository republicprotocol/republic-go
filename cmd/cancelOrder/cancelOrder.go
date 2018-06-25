package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	netHttp "net/http"
	"os"

	"github.com/republicprotocol/republic-go/contract"
	"github.com/republicprotocol/republic-go/crypto"
)

func main() {
	relayParam := flag.String("relay", "http://127.0.0.1:18515", "Binding and port of the relay")
	keystoreParam := flag.String("keystore", "", "Optionally encrypted keystore file")
	passphraseParam := flag.String("passphrase", "", "Optional passphrase to decrypt the keystore file")
	orderIDParam := flag.String("order", "", "ID of the order to be canceled")

	flag.Parse()

	keystore, err := loadKeystore(*keystoreParam, *passphraseParam)
	if err != nil {
		log.Fatalf("cannot load keystore: %v", err)
	}

	if *orderIDParam == "" {
		log.Fatal("Please provide the order ID in base 64 encoding.")
	}
	id, err := base64.StdEncoding.DecodeString(*orderIDParam)
	if err != nil {
		log.Fatal("wrong base64 encoding.")
	}

	signatureData := crypto.Keccak256([]byte("Republic Protocol: open: "), id)
	signatureData = crypto.Keccak256([]byte("\x19Ethereum Signed Message:\n32"), signatureData)
	signature, err := keystore.Sign(signatureData)

	log.Printf("id = %v; sig = %v", base64.StdEncoding.EncodeToString(id), base64.StdEncoding.EncodeToString(signature))

	if err != nil {
		log.Fatalf("cannot sign order: %v", err)
	}

	client := &netHttp.Client{}
	encodedSignature := base64.StdEncoding.EncodeToString(signature)
	req, err := netHttp.NewRequest("DELETE", fmt.Sprintf("%v/orders/?id=%v&signature=%v", *relayParam, *orderIDParam, encodedSignature), nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("cannot send request: %v", err)
	}
	defer res.Body.Close()

	log.Printf("status: %v", res.StatusCode)
	resText, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("cannot read response body: %v", err)
	}
	log.Printf("body: %v", string(resText))

}

func loadKeystore(keystoreFile, passphrase string) (crypto.Keystore, error) {
	file, err := os.Open(keystoreFile)
	if err != nil {
		return crypto.Keystore{}, err
	}
	defer file.Close()

	if passphrase == "" {
		keystore := crypto.Keystore{}
		if err := json.NewDecoder(file).Decode(&keystore); err != nil {
			return keystore, err
		}
		return keystore, nil
	}

	keystore := crypto.Keystore{}
	keystoreData, err := ioutil.ReadAll(file)
	if err != nil {
		return keystore, err
	}
	if err := keystore.DecryptFromJSON(keystoreData, passphrase); err != nil {
		return keystore, err
	}
	return keystore, nil
}

func loadConfig(configFile string) (contract.Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return contract.Config{}, err
	}
	defer file.Close()
	config := contract.Config{}
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return contract.Config{}, err
	}
	return config, nil
}
