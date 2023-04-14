package ports

import (
	pb "github.com/katallaxie/autobot/pkg/proto/v1"
)

// SubscribeRequest ...
type SubscribeRequest struct {
	Vars      map[string]string
	Arguments []string
}

// SubscribeResponse ...
type SubscribeResponse struct{}

// SetupRequest ...
type SetupRequest struct{}

// SetupResponse ...
type SetupResponse struct{}

// PublishRequest ...
type PublishRequest struct {
	Reply *pb.Reply
}

// PublishResponse ...
type PublishResponse struct {
}

// MessageStore ...
type MessageStore interface {
	// Subscribe ...
	Subscribe(SubscribeRequest, chan<- *pb.Event) (SubscribeResponse, error)

	// Publish ...
	Publish(PublishRequest) (PublishResponse, error)

	// Close ...
	Close() error
}
