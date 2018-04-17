package hyper

type Hash [32]byte

type Hasher interface {
	Hash() Hash
}

type Signatures = []Signature

func (signatures Signatures) Merge(others Signatures) Signatures {
	merger := map[Signature]struct{}{}
	for i := range signatures {
		merger[signatures[i]] = struct{}
	}
	for i := range others {
		merger[others[i]] = struct{}
	}

	i := 0
	mergedSignatures := make(Signatures, len(merger))
	for key := range merger {
		mergedSignatures[i] = key
		i++
	}
	return mergedSignatures
}

type Signature [32]byte

type Signer interface {
	Sign(Hasher) (Signature, error)
	Threshold() int
}
