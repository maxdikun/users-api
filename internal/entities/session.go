package entities

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	id          uuid.UUID
	user        uuid.UUID
	token       string
	createdAt   time.Time
	refreshedAt time.Time
	expiresAt   time.Time
}

func (s Session) ExpiresAt() time.Time {
	return s.expiresAt
}

func (s Session) Token() string {
	return s.token
}

func (s Session) User() uuid.UUID {
	return s.user
}

func (s *Session) Refresh(newToken string, duration time.Duration) {
	s.token = newToken
	s.refreshedAt = time.Now()
	s.expiresAt = time.Now().Add(duration)
}

func NewSession(user uuid.UUID, token string, duration time.Duration) Session {
	return Session{
		id:          uuid.New(),
		user:        user,
		token:       token,
		createdAt:   time.Now(),
		refreshedAt: time.Now(),
		expiresAt:   time.Now().Add(duration),
	}
}

func LoadSession(
	id uuid.UUID,
	user uuid.UUID,
	token string,
	createdAt time.Time,
	refreshedAt time.Time,
	expiresAt time.Time,
) Session {
	return Session{
		id:          id,
		user:        user,
		token:       token,
		createdAt:   createdAt,
		refreshedAt: refreshedAt,
		expiresAt:   expiresAt,
	}
}
