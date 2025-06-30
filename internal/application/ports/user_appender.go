package ports

import (
	"context"
	"fmt"

	"github.com/maxdikun/users-api/internal/entities"
)

type DuplicationError struct {
	Source string
	Object string
	Field  string
	Value  any
}

// Error implements error.
func (err *DuplicationError) Error() string {
	return fmt.Sprintf(
		"duplication error occurred in %s storage: object '%s' with value '%v' in field '%s' already exists",
		err.Source,
		err.Object,
		err.Value,
		err.Field,
	)
}

type UserAppender interface {
	AppendUser(ctx context.Context, user entities.User) error
}
