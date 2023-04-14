package plugin

import (
	"context"
	"regexp"

	cfg "github.com/katallaxie/autobot/pkg/config"
	pb "github.com/katallaxie/autobot/pkg/proto/v1"

	s "github.com/katallaxie/pkg/server"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

// ID ...
type ID string

// HandlerFunc ...
type HandlerFunc func(ctx context.Context, msg *pb.Message) (*pb.Reply, error)

// FilterFunc ...
type FilterFunc func(*pb.Message) (bool, error)

// Filters ...
type Filters []FilterFunc

// FilterArgumentText ...
func FilterArgumentText(pattern string) FilterFunc {
	return func(m *pb.Message) (bool, error) {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return false, err
		}

		return !re.MatchString(m.GetTextMessage().ArgumentText), nil
	}
}

// Serve ...
func Serve(handler Handler, opts ...Opt) error {
	p := new(plugin)
	p.handler = handler

	for _, o := range opts {
		o(p)
	}

	return serve(p)
}

// Handler ...
type Handler interface {
	Handle(ctx context.Context, msg *pb.Message) (*pb.Reply, error)
}

type plugin struct {
	id      string
	filters Filters
	url     string
	handler Handler

	s.Listener
}

// Opt ...
type Opt func(*plugin)

// WithWithURLUrl ...
func WithURL(url string) Opt {
	return func(p *plugin) {
		p.url = url
	}
}

// WithID ...
func WithID(id string) Opt {
	return func(p *plugin) {
		p.id = id
	}
}

// WithFilters ...
func WithFilters(filters Filters) Opt {
	return func(p *plugin) {
		p.filters = append(p.filters, filters...)
	}
}

func serve(p *plugin) error {
	root, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, _ := s.WithContext(root)
	srv.Listen(p, false)

	if err := srv.Wait(); err != nil {
		return err
	}

	return nil
}

// Start ...
func (p *plugin) Start(ctx context.Context, ready s.ReadyFunc, run s.RunFunc) func() error {
	return func() error {
		nc, err := nats.Connect(p.url)
		if err != nil {
			return err
		}
		defer nc.Close()

		ch := make(chan *nats.Msg, 64)

		sub, err := nc.ChanQueueSubscribe(cfg.DefaultMessagesSubject.String(), p.id, ch)
		if err != nil {
			return err
		}

		defer func() { _ = sub.Unsubscribe() }()

	OUTTER:
		for {
			select {
			case msg := <-ch:
				m := &pb.Message{}
				err := proto.Unmarshal(msg.Data, m)
				if err != nil {
					return err
				}

				for _, f := range p.filters {
					ok, err := f(m)
					if err != nil {
						return err
					}

					if ok {
						continue OUTTER
					}
				}

				reply, err := p.handler.Handle(ctx, m)
				if err != nil {
					return err
				}

				bb, err := proto.Marshal(reply)
				if err != nil {
					return err
				}

				err = nc.Publish(msg.Reply, bb)
				if err != nil {
					return err
				}

				_ = msg.Ack()
			case <-ctx.Done():
				return nil
			}
		}
	}
}
