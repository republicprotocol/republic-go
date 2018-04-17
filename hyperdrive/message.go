package hyperdrive

// A Message is passed between Replicas to advanced the state of the
// Hyperdrive blockchain.
type Message interface {
	Hash() Hash
	Fault() Fault
	Verify(Verifier) error
	SetSignatures(Signatures)
	GetSignatures() Signatures
}

type Verifier struct {

}

// A MessageStore stores and loads Messages.
type MessageStore interface {
	Load(Hash) Message
	Store(Message)
}

// MessageMapStore implements the MessageStore interface using an in-memory
// map.
type MessageMapStore struct {
	memory map[Hash]Message
}

// NewMessageMapStore returns a new MessageMapStore.
func NewMessageMapStore() MessageMapStore {
	return MessageMapStore{}
}

// Load implements the MessageStore interface.
func (store *MessageMapStore) Load(hash Hash) Message {
	if message, ok := store.memory[hash]; ok {
		return message
	}
	return nil
}

// Store implements the MessageStore interface.
func (store *MessageMapStore) Store(message Message) {
	store.memory[message.Hash()] = message
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
		signature, err := signer.Sign(fault.Hash())
		if err != nil {
			return nil, err
		}
		fault.Signatures = Signatures{signature}
		return fault, nil
	}

	if store.Load(message.Hash()) == nil {
		// If this message is not in the store, then sign it and store it
		signature, err := signer.Sign(message.Hash())
		if err != nil {
			return nil, err
		}
		message.SetSignatures(message.GetSignatures().Merge(Signatures{signature}))
		store.Store(message)
		// If the required threshold is reached, then return the message
		if len(message.GetSignatures()) >= threshold {
			return message, nil
		}
		return nil, nil
	}

	loaded := store.Load(message.Hash())
	if len(loaded.GetSignatures()) >= threshold {
		// The stored message has already reached the threshold
		return nil, nil
	}
	// Merge the signatures and if the required threshold is now reached, then
	// return the message
	message.SetSignatures(message.GetSignatures().Merge(loaded.GetSignatures()))
	store.Store(message)
	if len(message.GetSignatures()) >= threshold {
		return message, nil
	}
	return nil, nil
}
