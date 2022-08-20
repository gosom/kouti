package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gosom/kouti/logger"
)

func TestGETSayHello(t *testing.T) {
	log := logger.New(logger.Config{Debug: true})
	h := NewHelloWorldHander(log)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Error(err)
		return
	}
	response := httptest.NewRecorder()
	h.sayHello(response, req)

	if response.Code != http.StatusOK {
		t.Errorf("got status %d want %d", response.Code, http.StatusOK)
		return
	}

	want := map[string]string{
		"message": "hello world",
	}
	var got map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &got); err != nil {
		t.Error(err)
		return
	}
	if got["message"] != want["message"] {
		t.Errorf("got message %s want %s", got["message"], want["message"])
	}
}
