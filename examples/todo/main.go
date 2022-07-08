package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gosom/kouti/auth"
	"github.com/gosom/kouti/httpserver"
	"github.com/gosom/kouti/logger"
	"github.com/gosom/kouti/utils"
	"github.com/gosom/kouti/web"
	"github.com/rs/zerolog"
)

type HelloHandler struct {
	web.BaseHandler
}

func (h HelloHandler) SayHello(w http.ResponseWriter, r *http.Request) {
	h.Log(r, "we are saying hello %s", "world")
	h.Json(w, http.StatusOK, map[string]string{"msg": "hello world"})
}

func NewHelloHandler(log zerolog.Logger) HelloHandler {
	h := HelloHandler{}
	h.Logger = logger.NewSubLogger(log, "HelloHandler")
	return h
}

type TodoHandler struct {
	web.BaseHandler
}

func (h TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
}

func (h TodoHandler) List(w http.ResponseWriter, r *http.Request) {
}

func NewTodoHandler(log zerolog.Logger) TodoHandler {
	h := TodoHandler{}
	h.Logger = logger.NewSubLogger(log, "TodoHandler")
	return h
}

type UserHandler struct {
	web.BaseHandler
}

func (h UserHandler) Create(w http.ResponseWriter, r *http.Request) {
}

func NewUserHandler(log zerolog.Logger) UserHandler {
	h := UserHandler{}
	h.Logger = logger.NewSubLogger(log, "UserHandler")
	return h
}

// ==========================================================================
func main() {
	defer utils.ExitRecover()

	log := logger.New(logger.Config{Debug: true})

	router := web.NewRouter(web.RouterConfig{
		Log: log,
	})

	{
		helloHandler := NewHelloHandler(log)
		router.Get("/", helloHandler.SayHello)

		authenticator, err := auth.New(auth.AuthenticatorConfig{
			JwtSignKey: "secret",
			Log:        logger.NewSubLogger(log, "Authenticator"),
		})
		if err != nil {
			panic(err)
		}

		router.Route("/api/v1/user", func(r chi.Router) {
			h := NewUserHandler(log)
			r.Post("/", h.Create)
		})

		router.Route("/api/v1/todo", func(r chi.Router) {
			r.Use(web.Authentication(authenticator))
			h := NewTodoHandler(log)
			r.Post("/", h.Create)
			r.Get("/", h.List)
		})
	}

	cfg := httpserver.Config{
		Host:   "localhost:8080",
		Router: router,
		UseTLS: false,
	}

	if err := httpserver.Run(cfg); err != nil {
		panic(utils.Exit{1, err.Error()})
	}
}
