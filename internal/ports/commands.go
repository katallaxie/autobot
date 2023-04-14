package ports

import (
	"context"

	pb "github.com/katallaxie/autobot/pkg/proto/v1"
)

// Command ...
type Command interface {
	Do(context.Context, *pb.Event) (*pb.Reply, error)
}

// CommandStore ...
type CommandStore interface {
	// Do ...
	Do(context.Context, *pb.Event) (*pb.Reply, error)
}
