package main

import (
	"context"
	"log"
	"time"

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
						Text:     time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"),
						ThreadId: msg.GetTextMessage().GetThreadId(),
					},
				},
			},
		},
	}, nil
}

func main() {
	filters := []plugin.FilterFunc{
		plugin.FilterArgumentText(`(?i)(time)( now)? (.+)`),
	}

	err := plugin.Serve(
		new(handler),
		plugin.WithFilters(filters),
		plugin.WithID("autobot.plugin.time-now"),
	)
	if err != nil {
		log.Fatal(err)
	}
}
