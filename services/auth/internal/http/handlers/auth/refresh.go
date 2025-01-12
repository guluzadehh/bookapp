package authhttp

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/guluzadehh/bookapp/pkg/http/api"
	"github.com/guluzadehh/bookapp/pkg/http/middlewares/requestidmdw"
	"github.com/guluzadehh/bookapp/pkg/sl"
	"github.com/guluzadehh/bookapp/services/auth/internal/services"
)

type RefreshResponse struct {
	api.Response
	Data *RefreshData `json:"data"`
}

type RefreshData struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.refresh.New"

	log := sl.Init(h.Log, op, requestidmdw.GetId(r.Context()))

	cookie, err := r.Cookie(h.config.JWT.Refresh.CookieName)
	if err == http.ErrNoCookie {
		log.Info("refresh cookie doesn't exist", slog.String("refresh_cookie_name", h.config.JWT.Refresh.CookieName))
		h.JSON(w, http.StatusUnauthorized, refreshInvalidResponse())
		return
	}

	refresh := cookie.Value

	var oldAccess string

	if authHeader := r.Header.Get("Authorization"); strings.HasPrefix(authHeader, "Bearer ") {
		oldAccess = strings.TrimPrefix(authHeader, "Bearer ")
	}

	access, err := h.srvc.RefreshToken(r.Context(), refresh, oldAccess)
	if err != nil {
		if errors.Is(err, services.ErrInvalidToken) {
			h.JSON(w, http.StatusUnauthorized, refreshInvalidResponse())
			return
		}

		h.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
		return
	}

	h.JSON(w, http.StatusOK, RefreshResponse{
		Response: api.Ok(),
		Data: &RefreshData{
			Token: access,
		},
	})
}

func refreshInvalidResponse() api.Response {
	return api.Err("refresh token is invalid")
}