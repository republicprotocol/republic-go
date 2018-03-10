package darkocean

import (
	"errors"
	"time"

	do "github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/contracts/dnr"
)

// WatchForDarkOceanChanges returns a channel through which it will send an update every epoch
// Will check if a new epoch has been triggered and then sleep for 5 minutes
// Blocking function
func WatchForDarkOceanChanges(registrar dnr.DarkNodeRegistrarInterface, channel chan do.Option) {

	// This function runs until the channel is closed
	defer func() { recover() }()

	var currentBlockhash [32]byte

	// TODO loop until an epoch, call calculateDarkOcean()
	for {
		epoch, err := registrar.CurrentEpoch()
		if err != nil {
			channel <- do.Err(errors.New("Couldn't retrieve current epoch"))
		}

		if epoch.Blockhash != currentBlockhash {
			currentBlockhash = epoch.Blockhash
			darkOceanOverlay, err := GetDarkPools(registrar)
			if err != nil {
				channel <- do.Err(errors.New("Couldn't retrieve dark ocean overlay"))
			} else {
				channel <- do.Ok(darkOceanOverlay)
			}
		}
		time.Sleep(5 * 60 * time.Second)
	}
}
