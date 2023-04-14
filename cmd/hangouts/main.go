package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/katallaxie/autobot/pkg/adapter"
	pb "github.com/katallaxie/autobot/pkg/proto/v1"

	ps "cloud.google.com/go/pubsub"
	"github.com/caarlos0/env/v6"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/chat/v1"
	"google.golang.org/api/option"
)

// config ...
type config struct {
	Credentials  string `env:"AUTOBOT_HANGOUT_CREDENTIALS"`
	Subscription string `env:"AUTOBOT_HANGOUT_SUBSCRIPTION"`
	Project      string `env:"AUTOBOT_HANGOUT_PROJECT"`
}

var cfg *config

// Event ...
type Event struct {
	Type      string
	EventTime time.Time
	Space     *Space
	Message   *Message
}

// Space ...
type Space struct {
	Name        string
	DisplayName string
	Type        string
}

// Message ...
type Message struct {
	Name         string
	CreateTime   time.Time
	Text         string
	ArgumentText string
	Sender       *Sender
	Thread       *Thread
	Annotations  []Annotation
	SlashCommand *SlashCommand `json:"slashCommand,omitempty"`
}

// SlashCommand ...
type SlashCommand struct {
	Bot         *Bot
	CommandId   string
	CommandName string
	Type        string
}

// Bot ...
type Bot struct {
	Name        string
	DisplayName string
	AvatarUrl   string
	Type        string
}

// Sender ...
type Sender struct {
	Name        string
	DisplayName string
	AvatarUrl   string
	Email       string
}

// Thread ...
type Thread struct {
	Name string
}

// Annotation ...
type Annotation struct {
	Length       int
	StartIndex   int
	Type         string
	UserMention  *UserMention
	SlashCommand *SlashCommand
}

// UserMention ...
type UserMention struct {
	Type string
	User *User
}

// User ...
type User struct {
	AvatarUrl   string
	DisplayName string
	Name        string
	Type        string
}

func init() {
	cfg = &config{}
	if err := env.Parse(cfg); err != nil {
		log.Panic(err)
	}
}

type server struct {
	pb.UnimplementedAdapterServer
}

// FromSlashCommand ...
func FromSlashCommand(e Event) (*pb.Subscribe_Response, error) {
	id, err := strconv.ParseUint(e.Message.SlashCommand.CommandId, 10, 16)
	if err != nil {
		return nil, err
	}

	args := strings.Split(e.Message.ArgumentText, " ")
	for i, arg := range args {
		args[i] = strings.TrimSpace(arg)
	}

	return &pb.Subscribe_Response{
		Event: &pb.Event{
			Event: &pb.Event_Command{
				Command: &pb.Command{
					Id:          int32(id),
					Name:        e.Message.SlashCommand.CommandName,
					RawArgument: e.Message.ArgumentText,
					Arguments:   args,
					Room: &pb.Room{
						Name: e.Space.Name,
					},
					Message: &pb.Message{
						Message: &pb.Message_TextMessage{
							TextMessage: &pb.TextMessage{
								Text:     e.Message.Text,
								ThreadId: e.Message.Thread.Name,
							},
						},
					},
				},
			},
		},
	}, nil
}

// Subscribe ...
func (s *server) Subscribe(req *pb.Subscribe_Request, srv pb.Adapter_SubscribeServer) error {
	client, err := ps.NewClient(srv.Context(), cfg.Project, option.WithCredentialsJSON([]byte(cfg.Credentials)))
	if err != nil {
		return err
	}
	defer client.Close()

	msgs := make(chan []byte)

	g, ctx := errgroup.WithContext(srv.Context())

	g.Go(func() error {
		sub := client.Subscription(cfg.Subscription)
		err = sub.Receive(ctx, func(ctx context.Context, m *ps.Message) {
			msgs <- m.Data

			m.Ack()
		})
		if err != nil {
			return err
		}

		return nil
	})

	g.Go(func() error {
		for msg := range msgs {
			var e Event
			err := json.Unmarshal(msg, &e)
			if err != nil {
				return err
			}

			var res *pb.Subscribe_Response
			if e.Message.SlashCommand != nil {
				res, err = FromSlashCommand(e)
			} else {
				res = &pb.Subscribe_Response{
					Event: &pb.Event{
						Event: &pb.Event_Mentioned{
							Mentioned: &pb.Mentioned{
								Room: &pb.Room{
									Name: e.Space.Name,
								},
								Message: &pb.Message{
									Message: &pb.Message_TextMessage{
										TextMessage: &pb.TextMessage{
											Text:     e.Message.Text,
											ThreadId: e.Message.Thread.Name,
										},
									},
								},
							},
						},
					},
				}
			}

			if err != nil {
				continue
			}

			err = srv.Send(res)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

// Publish ...
func (s *server) Publish(ctx context.Context, req *pb.Publish_Request) (*pb.Publish_Response, error) {
	hc, err := getOauthClient(ctx, []byte(cfg.Credentials))
	if err != nil {
		return nil, err
	}

	chatClient, err := chat.NewService(ctx, option.WithHTTPClient(hc))
	if err != nil {
		return nil, err
	}

	c := chat.NewSpacesMessagesService(chatClient)

	cm := &chat.Message{
		Text: req.GetReply().GetMessage().GetTextMessage().GetText(),
		Thread: &chat.Thread{
			Name: req.GetReply().GetMessage().GetTextMessage().GetThreadId(),
		},
	}

	res := c.Create(req.GetReply().GetRoom().GetName(), cm)
	res = res.Context(ctx)

	_, err = res.Do()
	if err != nil {
		return nil, err
	}

	return &pb.Publish_Response{}, nil
}

func main() {
	adapter.Serve(&adapter.ServeOpts{
		GRPCAdapterFunc: func() pb.AdapterServer {
			return &server{}
		},
	})
}

func getOauthClient(ctx context.Context, credentials []byte) (*http.Client, error) {
	creds, err := google.CredentialsFromJSON(
		ctx,
		credentials,
		"https://www.googleapis.com/auth/chat.bot",
	)
	if err != nil {
		return nil, err
	}

	return oauth2.NewClient(ctx, creds.TokenSource), nil
}
