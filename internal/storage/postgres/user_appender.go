package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maxdikun/users-api/internal/application/ports"
	"github.com/maxdikun/users-api/internal/entities"
	"github.com/maxdikun/users-api/internal/storage/postgres/gen"
)

type UserAppender struct {
	pool *pgxpool.Pool
}

var _ ports.UserAppender = (*UserAppender)(nil)

func (u UserAppender) AppendUser(ctx context.Context, user entities.User) error {
	queries := gen.New(u.pool)

	err := queries.InsertUser(ctx, gen.InsertUserParams{
		ID:               user.Id(),
		Username:         string(user.Username()),
		Email:            string(user.Email()),
		Password:         string(user.Password()),
		EmailConfirmedAt: user.EmailConfirmedAt(),
		CreatedAt:        user.CreatedAt(),
		UpdatedAt:        user.UpdatedAt(),
		IsDeleted:        user.IsDeleted(),
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return &ports.DuplicationError{
					Source: "postgres.UserAppender",
					Object: "user",
					Field:  pgErr.ColumnName,
				}
			}
		}

		return err
	}

	return nil
}

func NewUserAppender(p *pgxpool.Pool) *UserAppender {
	return &UserAppender{
		pool: p,
	}
}
