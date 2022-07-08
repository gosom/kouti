package main

import (
	"net/http"

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

func main() {
	defer utils.ExitRecover()

	log := logger.New(logger.Config{Debug: true})

	router := web.NewRouter(web.RouterConfig{
		Log: log,
	})

	{
		helloHandler := NewHelloHandler(log)
		router.Get("/", helloHandler.SayHello)
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
