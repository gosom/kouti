package main

import (
	"context"
	"net/http"

	"github.com/gosom/kouti/httpserver"
	"github.com/gosom/kouti/logger"
	"github.com/gosom/kouti/utils"
	"github.com/gosom/kouti/web"
	"github.com/rs/zerolog"
)

func main() {
	log := logger.New(logger.Config{
		Debug: true,
	})
	fn := func(msg string) {
		log.Error().Msg(msg)
	}
	defer utils.ExitRecover(fn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := run(ctx, log); err != nil {
		panic(err)
	}
}

func run(ctx context.Context, log zerolog.Logger) error {
	router, err := web.NewRouter(web.RouterConfig{
		Log: logger.NewSubLogger(log, "router"),
	})
	if err != nil {
		return err
	}

	helloHn := NewHelloWorldHander(log)
	router.Get("/", helloHn.sayHello)
	router.Get("/panic", helloHn.sayHelloWithPanic)
	return httpserver.Run(ctx, httpserver.Config{
		Host:   "127.0.0.1:8080",
		Router: router,
	})
}

// HelloWorldHandler this is our Handler
type HelloWorldHandler struct {
	web.BaseHandler
}

// NewHelloWorldHander we just initialize the HelloWorldHandler and return an instance
func NewHelloWorldHander(log zerolog.Logger) HelloWorldHandler {
	ans := HelloWorldHandler{}
	ans.Logger = logger.NewSubLogger(log, "helloWorld")
	return ans
}

// sayHello just returns status code 200 and a message
func (h HelloWorldHandler) sayHello(w http.ResponseWriter, r *http.Request) {
	h.Logger.Debug().Msg("saying hello")
	h.Json(w, http.StatusOK, map[string]string{"message": "hello world"})
}

// sayHelloWithPanic the code below panics. Kouti will handle the panic and will print the stacktrace
func (h HelloWorldHandler) sayHelloWithPanic(w http.ResponseWriter, r *http.Request) {
	h.Logger.Debug().Msg("saying helloWithPanic")
	var ans map[string]string
	ans["message"] = "this will panic"
	h.Json(w, http.StatusOK, ans)
}
