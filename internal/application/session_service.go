package application

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/maxdikun/users-api/internal/application/ports"
	"github.com/maxdikun/users-api/internal/entities"
)

var (
	ErrUserAlreadyLoggedIn = errors.New("user already logged in")
	ErrInvalidToken        = errors.New("invalid token was provided")
)

type SessionService struct {
	logger *slog.Logger

	sessionAppender ports.SessionAppender
	sessionFinder   ports.SessionFinder
	sessionUpdater  ports.SessionUpdater

	maxTokenRetries     int
	sessionDuration     time.Duration
	accessTokenDuration time.Duration
	tokenSecret         string
}

type TokenSet struct {
	Access           string
	Refresh          string
	RefreshExpiresAt time.Time
}

func (svc *SessionService) CreateSession(ctx context.Context, user uuid.UUID) (TokenSet, error) {
	svc.logger.DebugContext(ctx, "Attempting to create new session", slog.String("user_id", user.String()))

	for i := 0; i < svc.maxTokenRetries; i++ {
		token, err := generateRandomString(32)
		if err != nil {
			svc.logger.ErrorContext(ctx, "Generating refresh token failed", slog.String("user_id", user.String()), slog.Any("error", err))
			return TokenSet{}, ErrInternal
		}

		session := entities.NewSession(user, token, svc.sessionDuration)
		if err := svc.sessionAppender.AppendSession(ctx, session); err != nil {
			var duplicationErr *ports.DuplicationError
			if errors.As(err, &duplicationErr) {
				if duplicationErr.Field == "token" {
					svc.logger.WarnContext(
						ctx, "Session token collision detected, retrying...",
						slog.String("user_id", user.String()),
						slog.String("attempted_token_prefix", token[:4]),
						slog.Int("retry_count", i+1),
					)
					continue
				}
				svc.logger.WarnContext(
					ctx, "User already has an active session (duplication on non-token field)",
					slog.String("user_id", user.String()),
					slog.String("duplication_field", duplicationErr.Field),
					slog.Any("error", err),
				)
				return TokenSet{}, ErrUserAlreadyLoggedIn
			}
			svc.logger.ErrorContext(
				ctx, "Failed to append session to repository",
				slog.String("user_id", user.String()),
				slog.Any("error", err),
			)
			return TokenSet{}, ErrInternal
		}

		accessToken, err := svc.generateAccessToken(session.User())
		if err != nil {
			svc.logger.ErrorContext(
				ctx, "Failed to generate access token after session creation",
				slog.String("user_id", user.String()),
				slog.Any("error", err),
			)
			return TokenSet{}, ErrInternal
		}

		svc.logger.InfoContext(
			ctx, "Session created successfully",
			slog.String("user_id", user.String()),
			slog.Time("refresh_expires_at", session.ExpiresAt()),
			slog.String("refresh_token_prefix", session.Token()[:4]),
		)

		return TokenSet{
			Access:           accessToken,
			Refresh:          session.Token(),
			RefreshExpiresAt: session.ExpiresAt(),
		}, nil
	}

	svc.logger.ErrorContext(
		ctx, "Failed to create unique session token after multiple retries",
		slog.String("user_id", user.String()),
		slog.Int("max_retries", svc.maxTokenRetries),
	)

	return TokenSet{}, ErrInternal
}

func (svc *SessionService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	svc.logger.DebugContext(
		ctx, "Attempting to refresh access token",
		slog.String("refresh_token_prefix", refreshToken[:4]),
	)

	session, err := svc.findSession(ctx, refreshToken)
	if err != nil {
		return "", err
	}

	accessToken, err := svc.generateAccessToken(session.User())
	if err != nil {
		svc.logger.WarnContext(
			ctx, "Attempted to generate access token",
			slog.String("user_id", session.User().String()),
			slog.Any("error", err),
		)
		return "", ErrInternal
	}

	return accessToken, nil
}

func (svc *SessionService) RefreshSession(ctx context.Context, refreshToken string) (TokenSet, error) {
	session, err := svc.sessionFinder.Find(ctx, refreshToken)
	if err != nil {
		var notFound *ports.NotFoundError
		if errors.As(err, &notFound) {
			return TokenSet{}, ErrInvalidToken
		}

		return TokenSet{}, ErrInternal
	}

	for i := 0; i < svc.maxTokenRetries; i++ {
		token, err := generateRandomString(32)
		if err != nil {
			return TokenSet{}, ErrInternal
		}
		session.Refresh(token, svc.sessionDuration)

		if err := svc.sessionUpdater.UpdateSession(ctx, session); err != nil {
			var duplicationErr *ports.DuplicationError
			if errors.As(err, &duplicationErr) {
				if duplicationErr.Field == "token" {
					continue
				}
			}

			return TokenSet{}, ErrInternal
		}

		accessToken, err := svc.generateAccessToken(session.User())
		if err != nil {
			return TokenSet{}, ErrInternal
		}

		return TokenSet{
			Access:           accessToken,
			Refresh:          session.Token(),
			RefreshExpiresAt: session.ExpiresAt(),
		}, nil
	}

	return TokenSet{}, ErrInternal
}

func (svc *SessionService) findSession(ctx context.Context, token string) (entities.Session, error) {
	session, err := svc.sessionFinder.Find(ctx, token)
	if err != nil {
		var notFound *ports.NotFoundError
		if errors.As(err, &notFound) {
			return entities.Session{}, ErrInvalidToken
		}

		return entities.Session{}, ErrInternal
	}

	return session, nil
}

func (svc *SessionService) generateAccessToken(user uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(svc.accessTokenDuration),
		"sub": user.String(),
	})

	tokenStr, err := token.SignedString(svc.tokenSecret)
	if err != nil {

	}
	return tokenStr, nil
}
