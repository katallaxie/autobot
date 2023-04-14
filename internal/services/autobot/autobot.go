package autobot

import (
	"context"

	"github.com/katallaxie/autobot/internal/controllers"

	s "github.com/katallaxie/pkg/server"
)

type autobot struct {
	m *controllers.MessageController

	s.Listener
}

// Opt ...
type Opt func(*autobot)

// New ...
func New(m *controllers.MessageController, opts ...Opt) s.Listener {
	r := new(autobot)
	r.m = m

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// Start ...
func (a *autobot) Start(ctx context.Context, ready s.ReadyFunc, run s.RunFunc) func() error {
	return func() error {
		err := a.m.Subscribe(ctx)
		if err != nil {
			return err
		}

		return nil

		// m := &adapter.Meta{
		// 	Path: a.adapter,
		// }
		// f := m.Factory(ctx)

		// adapter, err := f()
		// if err != nil {
		// 	return err
		// }
		// defer func() { _ = adapter.Close() }()

		// nn := make(map[string]network.Network)

		// for _, locale := range locales.Locales {
		// 	util.SerializeMessages(locale.Tag)

		// 	nn[locale.Tag] = training.CreateNeuralNetwork(locale.Tag, false)
		// }

		// _, err = adapter.Subscribe(pb.Subscribe_Request{}, c.Queue())
		// if err != nil {
		// 	return err
		// }

		// return nil

		// ch := make(chan *nats.Msg, 64)

		// sub, err := a.conn.ChanQueueSubscribe(cfg.DefaultInboxSubject.String(), "autobot", ch)
		// if err != nil {
		// 	return err
		// }
		// defer func() { _ = sub.Unsubscribe() }()

		// timer := time.NewTicker(1 * time.Second)

		// for {
		// 	select {
		// 	case msg := <-ch:
		// 		run(a.handleMessage(ctx, msg, nn))
		// 	case <-timer.C:
		// 	case <-ctx.Done():
		// 		return nil
		// 	}
		// }
	}
}

// func (a *autobot) handleMessage(ctx context.Context, msg *nats.Msg, nn map[string]network.Network) func() error {
// 	return func() error {
// 		_, cancel := context.WithTimeout(ctx, time.Second*10)
// 		defer cancel()

// 		m := &pb.Message{}
// 		err := proto.Unmarshal(msg.Data, m)
// 		if err != nil {
// 			return err
// 		}

// 		_, responseSentence :=
// 			analysis.NewSentence("en", m.GetTextMessage().GetArgumentText()).Calculate(*a.cache, nn["en"], "demo")

// 		reply := &pb.Reply{
// 			Reply: &pb.Reply_Message{
// 				Message: &pb.Message{
// 					Message: &pb.Message_TextMessage{
// 						TextMessage: &pb.TextMessage{
// 							ThreadId: m.GetTextMessage().GetThreadId(),
// 							Text:     responseSentence,
// 						},
// 					},
// 				},
// 			},
// 		}

// 		bb, err := proto.Marshal(reply)
// 		if err != nil {
// 			return err
// 		}

// 		// reply, err := a.conn.RequestWithContext(ctx, cfg.DefaultMessagesSubject.String(), msg.Data)
// 		// if err != nil && !errors.Is(context.DeadlineExceeded, err) {
// 		// 	return err
// 		// }

// 		// if errors.Is(context.DeadlineExceeded, err) {
// 		// 	return nil
// 		// }

// 		err = a.conn.Publish(msg.Reply, bb)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	}
// }
