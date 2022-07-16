package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog"
)

type JwtClaims struct {
	User  any                    `json:"user`
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

func (o *jwtAuthenticator) Authenticate(r *http.Request) (any, error) {
	token, err := o.getBearerTokenFromHeader(r)
	if err != nil {
		return nil, err
	}
	return o.validateAccessToken(o.signkey, token)
}

func (o *jwtAuthenticator) GetAccessToken(u any) (string, error) {
	claims := JwtClaims{
		User:  u,
		Extra: nil,
		StandardClaims: jwt.StandardClaims{
			Issuer:    o.issuer,
			Subject:   fmt.Sprintf("%v", u),
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
	return "", errors.New("invalid " + headerAuthorization + " header")
}

func (o *jwtAuthenticator) validateAccessToken(signingKey string, accessToken string) (any, error) {
	token, err := jwt.ParseWithClaims(accessToken, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, err
	}
	payload, ok := token.Claims.(*JwtClaims)
	if ok && token.Valid {
		return payload.User, nil
	}

	return payload.User, errors.New("invalid token")
}
