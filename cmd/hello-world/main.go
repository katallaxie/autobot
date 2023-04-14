package main

import (
	"context"
	"fmt"
	"log"

	"github.com/katallaxie/autobot/pkg/plugin/v1"
	pb "github.com/katallaxie/autobot/pkg/proto/v1"
)

type handler struct {
	plugin.Handler
}

func (h *handler) Handle(ctx context.Context, msg *pb.Message) (*pb.Reply, error) {
	return &pb.Reply{
		Reply: &pb.Reply_Message{
			Message: &pb.Message{
				Message: &pb.Message_TextMessage{
					TextMessage: &pb.TextMessage{
						Text:     fmt.Sprintf("Hello %s!", msg.GetTextMessage().GetFrom().GetUser().GetDisplayName()),
						ThreadId: msg.GetTextMessage().GetThreadId(),
					},
				},
			},
		},
	}, nil
}

func main() {
	filters := []plugin.FilterFunc{
		plugin.FilterArgumentText(`(?i)(hello)( me)? (.+)`),
	}

	err := plugin.Serve(
		new(handler),
		plugin.WithFilters(filters),
		plugin.WithID("autobot.plugin.hello-world"),
	)

	if err != nil {
		log.Fatal(err)
	}
}
