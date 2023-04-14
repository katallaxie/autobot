package controllers

import (
	"context"

	"github.com/katallaxie/autobot/internal/models"
	"github.com/katallaxie/autobot/internal/ports"
)

// SessionController ...
type SessionController struct {
	store ports.SessionStore
}

// NewSessionController ...
func NewSessionController(store ports.SessionStore) *SessionController {
	return &SessionController{store: store}
}

// Create ...
func (s *SessionController) Create(ctx context.Context, session *models.Session) error {
	return s.store.Insert(ctx, session)
}
