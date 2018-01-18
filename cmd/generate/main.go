package main

import (
	// "bufio"
	"fmt"
	"strconv"
	"math"
	// "io/ioutil"
	// "os"

	"github.com/republicprotocol/go-identity"
	// "github.com/republicprotocol/go-miner"
)


func main() {
	err := genereteMiners(3);
	if err != nil {
		fmt.Println("Something unexpected happened");
	}
}

// func generateMiners(uint count) error {
// 	d1 := []byte("hello\ngo\n")
// 	err := ioutil.WriteFile("/tmp/dat1", d1, 0644)
// 	check(err)


// 	f, err := os.Create("/tmp/dat2")
// 	check(err)


// 	defer f.Close()


// 	d2 := []byte{115, 111, 109, 101, 10}
// 	n2, err := f.Write(d2)
// 	check(err)
// 	fmt.Printf("wrote %d bytes\n", n2)


// 	n3, err := f.WriteString("writes\n")
// 	fmt.Printf("wrote %d bytes\n", n3)


// 	f.Sync()


// 	w := bufio.NewWriter(f)
// 	n4, err := w.WriteString("buffered\n")
// 	fmt.Printf("wrote %d bytes\n", n4)


// 	w.Flush()

// }

func genereteMiners() error {
	port := 120935;
	var MultiAddresses [3]string
	for i := 0; i < 3; i++ {
		keyPair, err := identity.NewKeyPair();
		if err != nil {
			return err;
		}
		address := keyPair.Address();
		// publicKey := keyPair.PublicKey;
		// privateKey := keyPair.PrivateKey;
		MultiAddresses[i] = "/ip4/127.0.0.1/tcp/"+strconv.Itoa(port)+"/republic/"+string(address);
		port ++;
	}

	for i := 0; i < 3; i++ {
		minerConfig :=  {
			"multi" : MultiAddresses[i],
			"bootstrap_multis" : [MultiAddresses[math.Mod(i+1, 3)], MultiAddresses[math.Mod(i+2, 3)]],
		}
		fmt.Print(minerConfig)
	}
	return nil;
}
