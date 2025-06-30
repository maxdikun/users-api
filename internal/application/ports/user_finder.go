package ports

import (
	"context"

	"github.com/maxdikun/users-api/internal/entities"
)

type UserFinder interface {
	FindByUsername(ctx context.Context, username entities.Username) (entities.User, error)
	FindByEmail(ctx context.Context, email entities.Email) (entities.User, error)
}
