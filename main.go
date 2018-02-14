package main

import (
	"fmt"
)

func main() {
	err, usage := Open("msZrrQTX1mh5x2iCNsuTCKbqyep1wc97qT", 0.1, "testnet", "testuser", "testpassword");
	if (err != nil){
		 fmt.Println(err);
	}
	fmt.Println(usage);
}