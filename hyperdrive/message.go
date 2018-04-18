package hyperdrive

import (
	"github.com/republicprotocol/republic-go/identity"
)

// A Message is passed between Replicas to advanced the state of the
// Hyperdrive blockchain.
type Message interface {
	Hash() identity.Hash
	Fault() *Fault
	Verify(identity.Verifier) error
	SetSignatures(identity.Signatures)
	GetSignatures() identity.Signatures
}

// A MessageStore stores and loads Messages.
type MessageStore interface {
	Load(identity.Hash) Message
	Store(Message)
}

// MessageMapStore implements the MessageStore interface using an in-memory
// map.
type MessageMapStore struct {
	memory map[identity.Hash]Message
}

// NewMessageMapStore returns a new MessageMapStore.
func NewMessageMapStore() MessageMapStore {
	return MessageMapStore{
		memory: map[identity.Hash]Message{},
	}
}

// Load implements the MessageStore interface.
func (store *MessageMapStore) Load(hash identity.Hash) Message {
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
func VerifyAndSignMessage(message Message, store MessageStore, signer identity.Signer, verifier identity.Verifier, threshold int) (interface{}, error) {

	// If the message is invalid then return a signed Fault
	if err := message.Verify(verifier); err != nil {
		fault := message.Fault()
		signature, err := signer.Sign(fault.Hash())
		if err != nil {
			return nil, err
		}
		fault.Signatures = identity.Signatures{signature}
		return fault, nil
	}

	if store.Load(message.Hash()) == nil {
		// If this message is not in the store, then sign it and store it
		signature, err := signer.Sign(message.Hash())
		if err != nil {
			return nil, err
		}
		message.SetSignatures(message.GetSignatures().Merge(identity.Signatures{signature}))
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
