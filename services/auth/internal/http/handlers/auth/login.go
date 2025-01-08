package authhttp

import (
	"errors"
	"net/http"

	"github.com/guluzadehh/bookapp/pkg/http/api"
	"github.com/guluzadehh/bookapp/pkg/http/middlewares/requestidmdw"
	"github.com/guluzadehh/bookapp/pkg/http/render"
	"github.com/guluzadehh/bookapp/pkg/sl"
	"github.com/guluzadehh/bookapp/services/auth/internal/services"
)

type AuthenticateReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticateRes struct {
	api.Response
	Data *AuthenticateData `json:"data,omitempty"`
}

type AuthenticateData struct {
	Token string `json:"access_token"`
}

func (h *AuthHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.Authenticate"

	log := sl.Init(h.Log, op, requestidmdw.GetId(r.Context()))

	var req AuthenticateReq
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("can't decode json", sl.Err(err))
		h.JSON(w, http.StatusBadRequest, api.Err("failed to parse request body"))
		return
	}

	access, refresh, err := h.srvc.Authenticate(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			h.JSON(w, http.StatusUnauthorized, api.Err("invalid credentials"))
			return
		}

		h.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
		return
	}

	http.SetCookie(
		w,
		&http.Cookie{
			Name:     h.config.JWT.Refresh.CookieName,
			Value:    refresh,
			SameSite: http.SameSiteNoneMode,
			HttpOnly: true,
			Path:     h.config.JWT.Refresh.Uri,
			MaxAge:   int(h.config.JWT.Refresh.Expire.Seconds()),
			Secure: func(env string) bool {
				if env == "prod" {
					return true
				} else {
					return false
				}
			}(h.config.Env),
		},
	)
	log.Info("refresh cookie has been set")

	h.JSON(w, http.StatusOK, AuthenticateRes{
		Response: api.Ok(),
		Data: &AuthenticateData{
			Token: access,
		},
	})
}
