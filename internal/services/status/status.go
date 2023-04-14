package status

import (
	"context"

	srv "github.com/katallaxie/pkg/server"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type status struct {
	addr string
	srv.Listener
}

// Opt ...
type Opt func(*status)

// WithAddr ...
func WithAddr(addr string) Opt {
	return func(s *status) {
		s.addr = addr
	}
}

// New ...
func New(opts ...Opt) *status {
	a := new(status)

	for _, opt := range opts {
		opt(a)
	}

	return a
}

// Start ...
func (s *status) Start(ctx context.Context, ready srv.ReadyFunc, run srv.RunFunc) func() error {
	return func() error {
		app := fiber.New()

		app.Use(recover.New())
		app.Use(requestid.New())
		app.Use(logger.New())

		app.Get("/", func(c *fiber.Ctx) error {
			return c.SendString("Hello, World 🤖!")
		})

		app.Get("/health", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		go func() {
			<-ctx.Done()
			_ = app.Shutdown()
		}()

		err := app.Listen(s.addr)
		if err != nil {
			return err
		}

		return err
	}
}
