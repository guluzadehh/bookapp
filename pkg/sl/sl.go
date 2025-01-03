package sl

import (
	"log/slog"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func Init(log *slog.Logger, op string, requestId string) *slog.Logger {
	log = log.With(slog.String("request_id", requestId))
	return log.With("op", op)
}
