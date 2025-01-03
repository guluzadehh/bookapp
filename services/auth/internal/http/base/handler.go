package httpbase

import (
	"log/slog"

	"github.com/guluzadehh/bookapp/pkg/http/render"
)

type Handler struct {
	*render.Responder
	Log *slog.Logger
}

func New(log *slog.Logger) *Handler {
	return &Handler{
		Responder: render.NewResponder(log),
		Log:       log,
	}
}
