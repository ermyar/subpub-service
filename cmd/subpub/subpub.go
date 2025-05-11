package subpub

import "context"

// MessageHandler is a callback function that processes messages delivered to subscribers.
type MessageHandler func(msg interface{})

type Subscription interface {
	// Unsubscribe will remove interest in the given subject.
	Unsubscribe()
}

type SubPub interface {
	// Subscribe creates an asynchronous queue subscriber on the given subject.
	Subscribe(subject string, cb MessageHandler) (Subscription, error)

	// Publish publishes the msg argument to the given subject.
	Publish(subj string, msg interface{}) error

	// Close will shutdown pub-sub system.
	// May be blocked by data delivery until the context is canceled.
	Close(ctx context.Context) error
}

// func NewSubPub() SubPub
