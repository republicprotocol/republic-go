package dispatch

// The Message interface is type-only interface.
type Message interface{}

// The MessageQueue interface defines a set of expected functionality for a
// queue to be integrated with the Dispatcher.
type MessageQueue interface {

	// Run the MessageQueue, processing all messages that are sent and
	// received.
	Run() error

	// Shutdown the MessageQueue gracefully.
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
