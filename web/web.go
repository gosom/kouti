package web

import (
	"net/http"

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
}

func NewRouter(cfg RouterConfig) *chi.Mux {
	cfg = setupDefaults(cfg)
	r := chi.NewRouter()
	for i := range cfg.Middlewares {
		r.Use(cfg.Middlewares[i])
	}
	r.NotFound(cfg.NotFoundHandler)
	r.MethodNotAllowed(cfg.MethodNotAllowedHandler)
	return r
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
