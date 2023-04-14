package ports

import (
	"context"

	"github.com/katallaxie/autobot/internal/models"
)

// SessionStore ...
type SessionStore interface {
	// Insert ...
	Insert(context.Context, *models.Session) error

	// Delete ...
	Delete(context.Context, *models.Session) error

	// Update ...
	Update(context.Context, *models.Session) error
}
