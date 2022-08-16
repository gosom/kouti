package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

var (
	ErrInvalidAuthToken = errors.New("invalid authentication token")
)

type AuthenticatorConfig struct {
	// When set it uses JWT authentication
	JwtSignKey          string
	AccessTokenDuration time.Duration
	Issuer              string
	// Log
	Log zerolog.Logger
}

type Authenticator interface {
	Authenticate(r *http.Request) (string, error)
	GetAccessToken(u string) (string, error)
}

func New(cfg AuthenticatorConfig) (Authenticator, error) {
	if len(cfg.JwtSignKey) > 0 {
		if cfg.AccessTokenDuration == 0 {
			cfg.AccessTokenDuration = time.Minute * 10
		}
		return newJwtAuthenticator(cfg.Log, cfg.JwtSignKey, cfg.Issuer, cfg.AccessTokenDuration), nil
	}
	return nil, errors.New("unsuppoted auth method")
}
