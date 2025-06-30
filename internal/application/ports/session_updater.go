package ports

import (
	"context"

	"github.com/maxdikun/users-api/internal/entities"
)

type SessionUpdater interface {
	UpdateSession(ctx context.Context, session entities.Session) error
}
