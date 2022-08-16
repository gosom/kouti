package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog"
)

type JwtClaims struct {
	User  string                 `json:"user`
	Extra map[string]interface{} `json:"extra,omitempty"`
	jwt.StandardClaims
}

type jwtAuthenticator struct {
	signkey   string
	issuer    string
	aduration time.Duration
	log       zerolog.Logger
}

func newJwtAuthenticator(log zerolog.Logger, signkey, issuer string, ad time.Duration) *jwtAuthenticator {
	ans := jwtAuthenticator{
		signkey:   signkey,
		issuer:    issuer,
		aduration: ad,
		log:       log,
	}
	return &ans
}

func (o *jwtAuthenticator) Authenticate(r *http.Request) (string, error) {
	token, err := o.getBearerTokenFromHeader(r)
	if err != nil {
		return "", err
	}
	return o.validateAccessToken(o.signkey, token)
}

func (o *jwtAuthenticator) GetAccessToken(u string) (string, error) {
	claims := JwtClaims{
		User:  u,
		Extra: nil,
		StandardClaims: jwt.StandardClaims{
			Issuer:    o.issuer,
			Subject:   u,
			ExpiresAt: time.Now().UTC().Add(o.aduration).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(o.signkey))
}

func (o *jwtAuthenticator) getBearerTokenFromHeader(r *http.Request) (string, error) {
	const (
		headerAuthorization = "Authorization"
		headerPrefixBearer  = "BEARER"
	)
	bearer := r.Header.Get(headerAuthorization)
	size := len(headerPrefixBearer) + 1
	if len(bearer) > size && strings.ToUpper(bearer[0:size-1]) == headerPrefixBearer {
		return bearer[size:], nil
	}
	return "", fmt.Errorf("%w : invalid %s header", ErrInvalidAuthToken, headerAuthorization)
}

func (o *jwtAuthenticator) validateAccessToken(signingKey string, accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return "", fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("%w : %s", ErrInvalidAuthToken, err.Error())
	}

	payload, ok := token.Claims.(*JwtClaims)
	if ok && token.Valid {
		return payload.User, nil
	}

	return payload.User, ErrInvalidAuthToken
}
