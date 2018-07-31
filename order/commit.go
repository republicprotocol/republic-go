package order

import (
	"math/big"
)

type Commitment struct {
	Index uint64 `json:"index"`

	Price         CoExpCommitment `json:"priceCommit"`
	Volume        CoExpCommitment `json:"volumeCommit"`
	MinimumVolume CoExpCommitment `json:"minimumVolumeCommit"`
}

type CoExpCommitment struct {
	Co  *big.Int `json:"co"`
	Exp *big.Int `json:"exp"`
}
