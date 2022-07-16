package web

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

type RouterConfig struct {
	// The middlewares to use
	// set this not an empty []func(http.Handler) http.Hander{} to not
	// use any middleware. Otherwise it uses the defaults
	// (see setupDefaults)
	Middlewares []func(http.Handler) http.Handler
	// NotFoundHandler
	NotFoundHandler func(w http.ResponseWriter, r *http.Request)
	// MethodNotAllowedHandler
	MethodNotAllowedHandler func(w http.ResponseWriter, r *http.Request)
	// Log
	Log zerolog.Logger
	// SwaggerUI
	SwaggerUI *SwaggerUIConfig
}

func NewRouter(cfg RouterConfig) (*chi.Mux, error) {
	cfg = setupDefaults(cfg)
	r := chi.NewRouter()
	for i := range cfg.Middlewares {
		r.Use(cfg.Middlewares[i])
	}
	r.NotFound(cfg.NotFoundHandler)
	r.MethodNotAllowed(cfg.MethodNotAllowedHandler)
	if cfg.SwaggerUI != nil {
		h, err := NewSwaggerUI(cfg.SwaggerUI)
		if err != nil {
			return nil, err
		}
		sp := cfg.SwaggerUI.Path
		if !strings.HasSuffix(sp, "/") {
			sp += "/"
		}
		sp += "*"
		r.Handle(sp, http.StripPrefix(cfg.SwaggerUI.Path, h))
	}
	return r, nil
}

func setupDefaults(cfg RouterConfig) RouterConfig {
	if cfg.NotFoundHandler == nil || cfg.MethodNotAllowedHandler == nil {
		h := DefaultHandler{}
		if cfg.NotFoundHandler == nil {
			cfg.NotFoundHandler = h.NotFound
		}
		if cfg.MethodNotAllowedHandler == nil {
			cfg.MethodNotAllowedHandler = h.MethodNotAllowed
		}
	}
	if cfg.Middlewares == nil {
		cfg.Middlewares = append(cfg.Middlewares,
			middleware.RequestID,
			middleware.RealIP,
			RequestLogger(cfg.Log),
			middleware.Recoverer,
		)
	}
	return cfg
}

func renderJson(w http.ResponseWriter, statusCode int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if value == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(value); err != nil {
		panic(err)
	}
}
