package main

import (
	"flag"
	"log"
	"os"
	"path"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	dataParam := flag.String("data", path.Join(os.Getenv("HOME"), ".darknode/data"), "Data directory")

	orderFragments, err := leveldb.OpenFile(path.Join(*dataParam, "orderFragments"), nil)
	if err != nil {
		log.Fatalf("cannot open leveldb: %v", err)
	}

	iter := orderFragments.NewIterator(nil, nil)
	fragmentsCount := 0
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		orderID := hexutil.Encode(key)
		log.Println(orderID)
		fragmentsCount++
	}
	log.Printf("Have %d order fragments", fragmentsCount)
	iter.Release()
	err = iter.Error()

	computations, err := leveldb.OpenFile(path.Join(*dataParam, "computation"), nil)
	if err != nil {
		log.Fatalf("cannot open leveldb: %v", err)
	}

	iter = computations.NewIterator(nil, nil)
	computationCount := 0
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		computeID := hexutil.Encode(key)
		log.Println(computeID)
		computationCount++
	}
	log.Printf("Have %d computations", computationCount)
	iter.Release()
	err = iter.Error()

}
