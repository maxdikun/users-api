package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxdikun/users-api/internal/application/ports"
	"github.com/maxdikun/users-api/internal/entities"
	"github.com/maxdikun/users-api/internal/storage/postgres/gen"
)

type UserFinder struct {
	pool *pgxpool.Pool
}

var _ ports.UserFinder = (*UserFinder)(nil)

func NewUserFinder(p *pgxpool.Pool) *UserFinder {
	return &UserFinder{
		pool: p,
	}
}

func (u UserFinder) FindByUsername(ctx context.Context, username entities.Username) (entities.User, error) {
	queries := gen.New(u.pool)

	res, err := queries.SelectUserByUsername(ctx, string(username))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.User{}, &ports.NotFoundError{
				Source: "postgres.UserFinder",
				Object: "user",
				Field:  "username",
			}
		}

		return entities.User{}, err
	}

	return u.convert(res), nil
}

func (u UserFinder) FindByEmail(ctx context.Context, email entities.Email) (entities.User, error) {
	queries := gen.New(u.pool)

	res, err := queries.SelectUserByEmail(ctx, string(email))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.User{}, &ports.NotFoundError{
				Source: "postgres.UserFinder",
				Object: "user",
				Field:  "email",
			}
		}

		return entities.User{}, err
	}

	return u.convert(res), nil
}

func (u UserFinder) convert(user gen.User) entities.User {
	username, _ := entities.NewUsername(user.Username)
	email, _ := entities.NewEmail(user.Email)

	return entities.LoadUser(
		user.ID,
		username,
		email,
		entities.RawPassword(user.Password),
		user.CreatedAt,
		user.EmailConfirmedAt,
		user.UpdatedAt,
		user.IsDeleted,
	)
}
