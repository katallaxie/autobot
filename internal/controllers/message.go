package controllers

import (
	"context"
	"errors"
	"time"

	"github.com/katallaxie/autobot/internal/network"
	"github.com/katallaxie/autobot/internal/ports"
	pb "github.com/katallaxie/autobot/pkg/proto/v1"

	gocache "github.com/patrickmn/go-cache"
)

// MessageController ...
type MessageController struct {
	store    ports.MessageStore
	commands ports.CommandStore
	analyzer ports.Analyzer
	nn       map[string]network.Network
}

// NewMessageController ...
func NewMessageController(store ports.MessageStore, commands ports.CommandStore, analyzer ports.Analyzer, nn map[string]network.Network) *MessageController {
	return &MessageController{
		store:    store,
		commands: commands,
		analyzer: analyzer,
		nn:       nn,
	}
}

// Subscribe ...
func (m *MessageController) Subscribe(ctx context.Context) error {
	queue := make(chan *pb.Event)
	cache := gocache.New(5*time.Minute, 5*time.Minute)

	_, err := m.store.Subscribe(ports.SubscribeRequest{}, queue)
	if err != nil {
		return err
	}

	for {
		msg, ok := <-queue
		if !ok {
			return errors.New("adapter error: no new messages")
		}

		var reply *pb.Reply
		if msg.GetCommand() != nil {
			reply, err = m.replyCommands(ctx, msg)
		} else {
			reply, err = m.replyMentioned(ctx, msg, cache)
		}

		if err != nil {
			continue // do not reply
		}

		_, err := m.store.Publish(ports.PublishRequest{Reply: reply})
		if err != nil {
			return err
		}
	}
}

func (m *MessageController) replyCommands(ctx context.Context, e *pb.Event) (*pb.Reply, error) {
	return m.commands.Do(ctx, e)
}

//nolint:unparam
func (m *MessageController) replyMentioned(_ context.Context, e *pb.Event, cache *gocache.Cache) (*pb.Reply, error) {
	_, responseSentence := m.analyzer.Analyze("en", e.GetMentioned().GetMessage().GetTextMessage().GetText()).Calculate(*cache, m.nn["en"], "demo", 0.004)

	reply := &pb.Reply{
		Room: &pb.Room{
			Name: e.GetMentioned().GetRoom().GetName(),
		},
		Reply: &pb.Reply_Message{
			Message: &pb.Message{
				Message: &pb.Message_TextMessage{
					TextMessage: &pb.TextMessage{
						Text:     responseSentence,
						ThreadId: e.GetMentioned().GetMessage().GetTextMessage().GetThreadId(),
					},
				},
			},
		},
	}

	return reply, nil
}
