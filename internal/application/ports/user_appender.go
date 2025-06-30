package ports

import (
	"context"

	"github.com/maxdikun/users-api/internal/entities"
)

type UserAppender interface {
	AppendUser(ctx context.Context, user entities.User) error
}
