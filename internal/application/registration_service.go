package application

import (
	"context"
	"errors"
	"log/slog"

	"github.com/maxdikun/users-api/internal/application/ports"
	"github.com/maxdikun/users-api/internal/entities"
)

var (
	ErrUsernameTaken = errors.New("provided username is taken")
	ErrEmailTaken    = errors.New("email is taken")
	ErrInternal      = errors.New("internal service error")
)

type RegistrationService struct {
	logger   *slog.Logger
	appender ports.UserAppender
}

func (svc *RegistrationService) RegisterUser(
	ctx context.Context,
	usernameValue string,
	emailValue string,
	passwordValue string,
) error {
	username, usernameErr := entities.NewUsername(usernameValue)
	email, emailErr := entities.NewEmail(emailValue)
	password, passwordErr := entities.NewPassword(passwordValue)
	if passwordErr != nil {
		var vErr *entities.ValidationError
		if !errors.As(passwordErr, &vErr) {
			svc.logger.ErrorContext(ctx, "Unexcpeted error creating password", slog.Any("error", passwordErr))
			return ErrInternal
		}
	}

	if err := errors.Join(usernameErr, emailErr, passwordErr); err != nil {
		svc.logger.WarnContext(ctx, "User registration input validation failed", slog.Any("error", err))
		return err
	}

	user := entities.NewUser(username, email, password)

	if err := svc.appender.AppendUser(ctx, user); err != nil {
		var duplicationErr *ports.DuplicationError
		if errors.As(err, &duplicationErr) {
			switch duplicationErr.Field {
			case "email":
				svc.logger.InfoContext(ctx, "User registration failed: email taken", slog.String("email", emailValue))
				return ErrEmailTaken
			case "username":
				svc.logger.InfoContext(ctx, "User registration failed: username taken", slog.String("username", usernameValue))
				return ErrUsernameTaken
			default:
				svc.logger.ErrorContext(ctx, "User registration failed: unhandled duplication field",
					slog.String("field", duplicationErr.Field), slog.Any("error", err))
				return ErrInternal
			}
		}

		svc.logger.ErrorContext(ctx, "Failed to append user to repository", slog.Any("error", err))
		return ErrInternal
	}

	svc.logger.InfoContext(ctx, "User registered successfully", slog.String("username", usernameValue), slog.String("email", emailValue))
	return nil
}
