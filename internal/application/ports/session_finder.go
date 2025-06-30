package ports

import (
	"context"

	"github.com/maxdikun/users-api/internal/entities"
)

type SessionFinder interface {
	Find(ctx context.Context, token string) (entities.Session, error)
}
