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

type RegisterService struct {
	logger   *slog.Logger
	appender ports.UserAppender
}

func NewRegisterService(logger *slog.Logger, appender ports.UserAppender) *RegisterService {
	return &RegisterService{
		logger:   logger,
		appender: appender,
	}
}

func (svc *RegisterService) Register(ctx context.Context, username string, password string, email string) error {
	svc.logger.DebugContext(ctx, "RegisterService.Register called")

	usernameObj, usernameErr := entities.NewUsername(username)
	emailObj, emailErr := entities.NewEmail(email)
	passwordObj, passwordErr := entities.NewPassword(password)

	if passwordErr != nil {
		var vErr *entities.ValidationError
		if !errors.As(passwordErr, &vErr) {
			svc.logger.ErrorContext(ctx, "Failed to generate a password", slog.Any("error", passwordErr))
			return ErrInternal
		}
	}

	err := errors.Join(usernameErr, emailErr, passwordErr)
	if err != nil {
		svc.logger.WarnContext(ctx, "Input validation failed", slog.Any("error", err))
		return err
	}

	user := entities.NewUser(usernameObj, emailObj, passwordObj)
	err = svc.appender.AppendUser(ctx, user)
	if err != nil {
		var duplicationErr *ports.DuplicationError
		if errors.As(err, &duplicationErr) {
			switch duplicationErr.Field {
			case "email":
				svc.logger.InfoContext(
					ctx, "User registration failed: email taken",
					slog.Any("email", user.Email()),
				)
				return ErrEmailTaken
			case "username":
				svc.logger.InfoContext(
					ctx, "User registration failed: username taken",
					slog.Any("username", user.Email()),
				)
				return ErrUsernameTaken
			default:
				svc.logger.ErrorContext(
					ctx, "User registration failed: unhandled duplication field",
					slog.String("field", duplicationErr.Field),
					slog.Any("error", err),
				)
				return ErrInternal
			}
		}

		svc.logger.ErrorContext(
			ctx, "Failed to append user to repository",
			slog.Any("error", err),
		)
		return ErrInternal
	}

	svc.logger.InfoContext(ctx, "Successfully registered a user")

	return nil
}
