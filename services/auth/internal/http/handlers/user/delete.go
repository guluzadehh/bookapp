package userhttp

import (
	"net/http"

	"github.com/guluzadehh/bookapp/pkg/http/api"
	"github.com/guluzadehh/bookapp/services/auth/internal/http/middlewares/authmdw"
)

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := authmdw.User(r.Context())

	if err := h.srvc.Delete(r.Context(), user.Email); err != nil {
		h.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
		return
	}

	h.JSON(w, http.StatusNoContent, api.Ok())
}
