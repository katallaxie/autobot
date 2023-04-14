package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/katallaxie/autobot/internal/ports"

	pb "github.com/katallaxie/autobot/pkg/proto/v1"
)

type debug struct {
	ports.Command
}

// NewDebug ...
func NewDebug() *debug {
	d := new(debug)

	return d
}

// Do
func (d *debug) Do(ctx context.Context, event *pb.Event) (*pb.Reply, error) {
	b, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	json, err := prettyPrint(b)
	if err != nil {
		return nil, err
	}

	reply := &pb.Reply{
		Room: &pb.Room{
			Name: event.GetCommand().GetRoom().GetName(),
		},
		Reply: &pb.Reply_Message{
			Message: &pb.Message{
				Message: &pb.Message_TextMessage{
					TextMessage: &pb.TextMessage{
						Text:     fmt.Sprintf("```%s```", json),
						ThreadId: event.GetCommand().GetMessage().GetTextMessage().GetThreadId(),
					},
				},
			},
		},
	}

	return reply, nil
}

func prettyPrint(b []byte) (string, error) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, b, "", "\t")
	if err != nil {
		return "", err
	}

	return prettyJSON.String(), nil
}
