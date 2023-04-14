package cmd

import (
	"context"
	"sync"

	"github.com/katallaxie/autobot/internal/ports"

	pb "github.com/katallaxie/autobot/pkg/proto/v1"
)

var (
	registry     = map[int32]ports.Command{}
	registerSync sync.Mutex
)

// Register ...
func Register(id int32, cmd ports.Command) {
	registerSync.Lock()
	defer registerSync.Unlock()

	registry[id] = cmd
}

type commands struct {
	ports.CommandStore
}

// New ...
func New() *commands {
	c := new(commands)

	return c
}

// Do ...
func (c *commands) Do(ctx context.Context, e *pb.Event) (*pb.Reply, error) {
	registerSync.Lock()
	defer registerSync.Unlock()

	cmd, ok := registry[e.GetCommand().GetId()]
	if !ok {
		return &pb.Reply{
			Room: &pb.Room{
				Name: e.GetCommand().GetRoom().GetName(),
			},
			Reply: &pb.Reply_Message{
				Message: &pb.Message{
					Message: &pb.Message_TextMessage{
						TextMessage: &pb.TextMessage{
							Text:     "Sorry, this command is unknown.",
							ThreadId: e.GetCommand().GetMessage().GetTextMessage().GetThreadId(),
						},
					},
				},
			},
		}, nil
	}

	reply, err := cmd.Do(ctx, e)
	if err != nil {
		return nil, err
	}

	return reply, nil
}
