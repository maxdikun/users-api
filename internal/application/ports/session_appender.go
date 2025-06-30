package ports

import (
	"context"

	"github.com/maxdikun/users-api/internal/entities"
)

type SessionAppender interface {
	AppendSession(ctx context.Context, session entities.Session) error
}
