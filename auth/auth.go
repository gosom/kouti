package auth

import (
	"errors"
	"net/http"

	"github.com/rs/zerolog"
)

type AuthenticatorConfig struct {
	// When set it uses JWT authentication
	JwtSignKey string
	// Log
	Log zerolog.Logger
}

type Authenticator interface {
	Authenticate(r *http.Request) (any, error)
}

func New(cfg AuthenticatorConfig) (Authenticator, error) {
	if len(cfg.JwtSignKey) > 0 {
		return newJwtAuthenticator(cfg.Log, cfg.JwtSignKey), nil
	}
	return nil, errors.New("unsuppoted auth method")
}
