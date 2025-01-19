package authhttp

import (
	"net/http"

	"github.com/guluzadehh/bookapp/pkg/http/api"
	"github.com/guluzadehh/bookapp/pkg/http/httputils"
	"github.com/guluzadehh/bookapp/pkg/http/middlewares/requestidmdw"
	"github.com/guluzadehh/bookapp/pkg/sl"
)

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.logout.New"

	log := sl.Init(h.Log, op, requestidmdw.GetId(r.Context()))

	var refresh string
	cookie, err := r.Cookie(h.config.JWT.Refresh.CookieName)
	if err == nil {
		refresh = cookie.Value
	}

	access, _ := httputils.BearerToken(r)

	h.srvc.Logout(r.Context(), access, refresh)

	http.SetCookie(w, &http.Cookie{
		Name:     h.config.JWT.Refresh.CookieName,
		Value:    "",
		SameSite: http.SameSiteNoneMode,
		Path:     "/api/refresh",
		HttpOnly: true,
		MaxAge:   -1,
	})
	log.Info("refresh cookie has been deleted")

	h.JSON(w, http.StatusNoContent, api.Ok())
}
