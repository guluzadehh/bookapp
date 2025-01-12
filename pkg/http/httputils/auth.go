package httputils

import (
	"fmt"
	"net/http"
	"strings"
)

func BearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid Authorization header format")
	}

	return strings.TrimPrefix(authHeader, "Bearer "), nil
}
