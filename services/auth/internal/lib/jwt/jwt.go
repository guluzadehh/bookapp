package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/guluzadehh/bookapp/services/auth/internal/config"
)

type TokenType int

const (
	REFRESH TokenType = iota
	ACCESS
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type AuthClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func MakeRefresh(email, role string, config *config.Config) (string, error) {
	return MakeToken(email, role, config, REFRESH)
}

func MakeAccess(email, role string, config *config.Config) (string, error) {
	return MakeToken(email, role, config, ACCESS)
}

func MakeToken(email, role string, config *config.Config, tokenType TokenType) (string, error) {
	var expire time.Duration

	if tokenType == REFRESH {
		expire = config.JWT.Refresh.Expire
	} else if tokenType == ACCESS {
		expire = config.JWT.Access.Expire
	}

	claims := AuthClaims{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(config.JWT.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenStr, err
}

func Verify(tokenStr string, config *config.Config) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return token, nil
}
