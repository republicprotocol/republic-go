package hyper

// A Message is passed between Replicas to advanced the state of the
// Hyperdrive blockchain.
type Message interface {
	Hash() Hash
	Fault() Fault
	Verify() error
	SetSignatures(Signatures)
	Signatures() Signatures
}

// A MessageStore stores and loads Messages.
type MessageStore interface {
	Store(Message)
	Message(Hash) Message
}

// VerifyAndSignMessage using a MessageStore to keep the state of previously
// seen Messages, and their Signatures. A Signer is used to Sign the Message,
// and a threshold is used to define the number of Signatures needed. Returns
// the Message, if it has reached its threshold of Signatures. Returns a Fault
// if the Message is invalid. Otherwise, returns nil, and any error
// encountered.
func VerifyAndSignMessage(message Message, store MessageStore, signer Signer, threshold int) (interface{}, error) {

	// If the message is invalid then return a signed Fault
	if err := message.Verify(); err != nil {
		fault := message.Fault()
		signature, err := signer.Sign(fault)
		if err != nil {
			return nil, err
		}
		fault.Signatures = Signatures{signature}
		return fault, nil
	}

	if !store.Message(message) {
		// If this message is not in the store, then sign it and store it
		signature, err := signer.Sign(message.Hash())
		if err != nil {
			return nil, err
		}
		message.SetSignatures(message.Signatures().Merge(Signatures{signature}))
		store.Store(message)
		// If the required threshold is reached, then return the message
		if len(message.Signatures()) >= threshold {
			return message, nil
		}
		return nil, nil
	}

	stored := store.Message(message)
	if len(stored.Signatures()) >= threshold {
		// The stored message has already reached the threshold
		return nil, nil
	}
	// Merge the signatures and if the required threshold is now reached, then
	// return the message
	message.SetSignatures(message.Signatures().Merge(stored.Signatures()))
	store.Store(message)
	if len(message.Signatures()) >= threshold {
		return message, nil
	}
	return nil, nil
}
