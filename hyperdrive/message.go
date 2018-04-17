package hyper

type Message interface {
	Hasher

	Fault() Fault
	Verify() error
	SetSignatures(Signatures)
	Signatures() Signatures
}

func Foo(message Message, messageStorage map[Hash]struct{}, signer Signer) (interface{}, error) {
	hash := message.Hash()

	if err := message.Verify(); err != nil {
		fault := message.Fault()
		fault.Signatures = Signatures{signer.Sign(fault)}
	}

	if _, ok := messages[hash]; !ok {
		messages[hash] = message
		signature, err := signer.Sign(message)
		if err != nil {
			return
		}
		messages[hash].SetSignatures(messages[hash].Signatures().Merge(Signatures{signature}))
		return messages[hash], nil
	}

	messages[hash].SetSignatures(messages[hash].Signatures().Merge(message.Signatures()))
	if len(messages[hash].Signatures()) >= threshold {
		return messages[hash], nil
	}
	return nil, nil
}
