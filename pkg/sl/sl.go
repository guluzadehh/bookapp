package sl

import (
	"log/slog"
	"net/http"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

type loggerSetter interface {
	SetLog(log *slog.Logger)
}

func HandlerJob(log *slog.Logger, op string, r *http.Request, ls loggerSetter) *slog.Logger {
	log = log.With(slog.String("request_id", r.Header.Get("X-Request-Id")))
	ls.SetLog(log)
	return log.With("op", op)
}
