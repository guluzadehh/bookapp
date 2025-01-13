package authmdw

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/guluzadehh/bookapp/pkg/http/api"
	"github.com/guluzadehh/bookapp/pkg/http/httputils"
	"github.com/guluzadehh/bookapp/pkg/http/middlewares/requestidmdw"
	"github.com/guluzadehh/bookapp/pkg/http/render"
	"github.com/guluzadehh/bookapp/pkg/sl"
	"github.com/guluzadehh/bookapp/services/auth/internal/config"
	"github.com/guluzadehh/bookapp/services/auth/internal/domain/models"
	"github.com/guluzadehh/bookapp/services/auth/internal/lib/jwt"
	"github.com/guluzadehh/bookapp/services/auth/internal/services"
	"github.com/guluzadehh/bookapp/services/auth/internal/storage"
)

type contextKey string

const userContextKey contextKey = "user"

type AuthService interface {
	VerifyToken(ctx context.Context, token string) (*jwt.AuthClaims, error)
}

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

func Authorize(
	log *slog.Logger,
	config *config.Config,
	authService AuthService,
	userService UserService,
) func(next http.Handler) http.Handler {
	log = log.With(slog.String("component", "middleware/auth"))
	log.Info("auth middleware is enabled")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middlewares.authMdw.Authorize"

			log := log.With(
				slog.String("op", op),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", requestidmdw.GetId(r.Context())),
			)

			render := render.NewResponder(log)

			log.Info("authorizing the user")

			tokenStr, err := httputils.BearerToken(r)
			if err != nil {
				log.Error("failed to get access token from header", sl.Err(err))
				render.JSON(w, http.StatusUnauthorized, authFailResponse())
				return
			}

			claims, err := authService.VerifyToken(r.Context(), tokenStr)
			if err != nil {
				if errors.Is(err, services.ErrInvalidToken) {
					log.Warn("invalid access token")
					render.JSON(w, http.StatusUnauthorized, authFailResponse())
					return
				}

				log.Error("failed to verify token", sl.Err(err))
				render.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
				return
			}

			user, err := userService.GetUserByEmail(r.Context(), claims.Email)
			if err != nil {
				if errors.Is(err, storage.UserNotFound) {
					log.Warn("access token user not found in storage")
					render.JSON(w, http.StatusUnauthorized, authFailResponse())
					return
				}

				log.Error("failed to get user by username from storage", sl.Err(err))
				render.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func User(ctx context.Context) *models.User {
	user, ok := ctx.Value(userContextKey).(*models.User)
	if !ok || user == nil {
		return nil
	}
	return user
}

func authFailResponse() api.Response {
	return api.Err("you are not authorized")
}
