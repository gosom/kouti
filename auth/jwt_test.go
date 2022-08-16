package auth

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJwtAuthenticator(t *testing.T) {
	a, err := New(AuthenticatorConfig{
		JwtSignKey: "secret",
		Issuer:     "test_app",
	})

	if err != nil {
		t.Error(err.Error())
		return
	}

	uid := uuid.New().String()

	accessToken, err := a.GetAccessToken(uid)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(accessToken) == 0 {
		t.Error("expected non empty access token")
		return
	}

	req, err := http.NewRequest(http.MethodPost, "http://example.com", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	retrieved, err := a.Authenticate(req)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if retrieved != uid {
		t.Errorf("expected retrieved uid to be %s but got %s", uid, retrieved)
		return
	}
}

func TestJwtAuthenticatorExpired(t *testing.T) {
	a, err := New(AuthenticatorConfig{
		JwtSignKey:          "secret",
		Issuer:              "test_app",
		AccessTokenDuration: time.Second,
	})

	if err != nil {
		t.Error(err.Error())
		return
	}

	uid := uuid.New().String()

	accessToken, err := a.GetAccessToken(uid)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(accessToken) == 0 {
		t.Error("expected non empty access token")
		return
	}

	time.Sleep(2 * time.Second)
	req, err := http.NewRequest(http.MethodPost, "http://example.com", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	_, err = a.Authenticate(req)
	if err == nil {
		t.Error("expected authentication error")
		return
	}
	if !errors.Is(err, ErrInvalidAuthToken) {
		t.Error("expected error to be ErrInvalidAuthToken")
		return
	}
}
