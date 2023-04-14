package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/katallaxie/autobot/internal/ports"
)

// WebhookController ...
type WebhookController struct {
	message ports.MessageStore
}

// NewWebhookController ...
func NewWebhookController(message ports.MessageStore) *WebhookController {
	return &WebhookController{
		message: message,
	}
}

// Publish ...
func (w *WebhookController) Handler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}
