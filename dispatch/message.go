package dispatch

// The Message interface is type-only interface.
type Message interface{}

// MessageQueues is a slice of MessageQueue interfaces.
type MessageQueues []MessageQueue

// The MessageQueue interface defines a set of expected functionality for a
// queue to be integrated with the Dispatcher.
type MessageQueue interface {

	// Run the MessageQueue, processing all messages that are sent and
	// received. Run must only be called once.
	Run() error

	// Shutdown the MessageQueue gracefully. Shutdown implementations should be
	// idempotent, and handle multiple calls without panicking, or returning an
	// error.
	Shutdown() error

	// Send a Message to the MessageQueue. The implementation should throw a
	// type error if it receives a concrete type that it does not recognize.
	// This method should block if the MessageQueue is full.
	Send(Message) error

	// Recv a message from the MessageQueue. Return nil, and false, if the
	// MessageQueue has been shutdown, otherwise return a Message, and true.
	// This method should block if the MessageQueue is empty.
	Recv() (Message, bool)
}
