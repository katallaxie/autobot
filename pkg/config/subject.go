package config

import "fmt"

// Subject ...
type Subject string

func (s Subject) String() string {
	return string(s)
}

var (
	// DefaultMessagesSubject ...
	DefaultMessagesSubject Subject = "autobot.messages"

	// DefaultInboxSubject ...
	DefaultInboxSubject Subject = "autobot.inbox"

	// DefaultOutboxSubject ...
	DefaultOutboxSubject Subject = "autobot.outbox.>"
)

// NewInboxSubject ...
func NewInboxSubject(id string) Subject {
	return Subject(fmt.Sprintf("autobot.inbox.%s", id))
}

// NewOutboxSubject ...
func NewOutboxSubject(id string) Subject {
	return Subject(fmt.Sprintf("autobot.outbox.%s", id))
}
