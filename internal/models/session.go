package models

import "time"

// Session ...
type Session struct {
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

	// CreatedAt ...
	CreatedAt time.Time

	// UpdatedAt ...
	UpdatedAt time.Time

	// DeletedAt ...
	DeletedAt time.Time
}
