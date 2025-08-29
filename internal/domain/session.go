package domain

import "github.com/google/uuid"

// for redis
type Session struct {
	Uuid      uuid.UUID
	Name      string
	AvatarURL string
	Saved     bool
}
