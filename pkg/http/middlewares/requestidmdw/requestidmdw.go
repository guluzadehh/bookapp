package requestidmdw

import (
	"context"
	"net/http"
)

type requestIdKey string

var requestId requestIdKey = "requestId"

func GetId(ctx context.Context) string {
	val, _ := ctx.Value(requestId).(string)
	return val
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), requestId, r.Header.Get("X-Request-Id"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
