package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/katallaxie/autobot/internal/controllers"
	srv "github.com/katallaxie/pkg/server"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/google/go-github/v47/github"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

// Opts ...
type Opts struct {
	addr         string
	clientID     string
	clientSecret string
	redirectURL  string
}

// Opt ...
type Opt func(*Opts)

// Configure ...
func (o *Opts) Configure(opts ...Opt) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithAddr ...
func WithAddr(addr string) Opt {
	return func(o *Opts) {
		o.addr = addr
	}
}

// WithClientID ...
func WithClientID(clientId string) Opt {
	return func(o *Opts) {
		o.clientID = clientId
	}
}

// WithClientSecret ...
func WithClientSecret(clientSecret string) Opt {
	return func(o *Opts) {
		o.clientSecret = clientSecret
	}
}

// WithRedirectURL ...
func WithRedirectURL(redirectURL string) Opt {
	return func(o *Opts) {
		o.redirectURL = redirectURL
	}
}

type service struct {
	s      *controllers.SessionController
	config *oauth2.Config

	opts *Opts
}

// New ...
func New(s *controllers.SessionController, opts ...Opt) *service {
	options := new(Opts)
	options.Configure(opts...)

	api := new(service)
	api.opts = options
	api.s = s

	api.config = &oauth2.Config{
		ClientID:     api.opts.clientID,
		ClientSecret: api.opts.clientSecret,
		Scopes:       []string{"user:email", "repo"},
		RedirectURL:  fmt.Sprintf("%s/auth/callback", api.opts.redirectURL),
		Endpoint:     endpoints.GitHub,
	}

	return api
}

// Start ...
func (s *service) Start(ctx context.Context, ready srv.ReadyFunc, run srv.RunFunc) func() error {
	return func() error {
		engine := html.NewFileSystem(http.FS(files), ".html")
		engine.Debug(true)

		app := fiber.New(fiber.Config{
			Views: engine,
		})

		app.Get("/auth/login", func(c *fiber.Ctx) error {
			url := s.config.AuthCodeURL("pseudo-random", oauth2.AccessTypeOnline)
			return c.Redirect(url, http.StatusTemporaryRedirect)
		})

		app.Get("/auth/callback", func(c *fiber.Ctx) error {

			code := c.FormValue("code")
			token, err := s.config.Exchange(c.Context(), code)
			if err != nil {
				return err
			}

			state := c.FormValue("state")
			if state != "pseudo-random" {
				return c.Redirect("/", http.StatusTemporaryRedirect)
			}

			oauthClient := s.config.Client(c.Context(), token)
			client := github.NewClient(oauthClient)
			emails, _, err := client.Users.ListEmails(c.Context(), &github.ListOptions{})
			if err != nil {
				return c.Redirect("/", http.StatusTemporaryRedirect)
			}

			for _, email := range emails {
				log.Print(*email.Email)
			}

			user, _, err := client.Users.Get(c.Context(), "")
			if err != nil {
				return c.Redirect("/", http.StatusTemporaryRedirect)
			}

			log.Print(user.Email)

			return nil
		})

		app.Get("/", func(c *fiber.Ctx) error {
			return c.Render("templates/index", fiber.Map{})
		})

		if err := app.Listen(s.opts.addr); err != nil {
			return err
		}

		return nil
	}
}
