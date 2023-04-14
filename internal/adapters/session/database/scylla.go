package database

import (
	"context"
	"time"

	"github.com/katallaxie/autobot/internal/models"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
)

// Record ...
type Record struct {
	// UUID ...
	UUID string

	// Provider ...
	Provider string

	// Username ...
	Username string

	// Email ...
	Email string

	// Token ...
	Token string

	// ExpiresAt ...
	ExpiresAt time.Time

	// CreatedAt ...
	CreatedAt time.Time

	// UpdatedAt ...
	UpdatedAt time.Time

	// DeletedAt ...
	DeletedAt time.Time
}

// Opts ...
type Opts struct {
	Hosts []string
}

// Configure ...
func (o *Opts) Configure(opts ...Opt) {
	for _, opt := range opts {
		opt(o)
	}
}

// Opt ...
type Opt func(*Opts)

type scylla struct {
	opts *Opts

	session gocqlx.Session
}

// New ...
func New(session gocqlx.Session, opts ...Opt) *scylla {
	options := new(Opts)
	options.Configure(opts...)

	s := new(scylla)
	s.opts = options
	s.session = session

	return s
}

// Get ...
func (s *scylla) Get() (*models.Session, error) {
	return nil, nil
}

// Insert ...
func (s *scylla) Insert(ctx context.Context, session *models.Session) error {
	insertStmt := qb.
		Insert("session").
		Columns("uuid", "provider", "username", "created_at").
		Query(s.session).
		WithContext(ctx)

	record := &Record{
		UUID:      session.UUID,
		Provider:  session.Provider,
		Username:  session.Username,
		CreatedAt: time.Now(),
	}

	err := insertStmt.BindStruct(record).ExecRelease() //nolint:contextcheck
	if err != nil {
		return err
	}

	return nil
}
