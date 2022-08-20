package main

import (
	"context"
	"embed"

	"github.com/go-chi/chi/v5"
	"github.com/gosom/kouti/dbdriver"
	"github.com/gosom/kouti/httpserver"
	"github.com/gosom/kouti/logger"
	"github.com/gosom/kouti/utils"
	"github.com/gosom/kouti/web"
	"github.com/rs/zerolog"
)

//go:generate swag init -g main.go --parseDependency

//go:embed docs/swagger.json
var specFs embed.FS

// @title Todo API based on kouti
// @version 0.1
// @description This is a sample server todo server.
// @description You can visit the GitHub repository at https://github.com/gosom/kouti

// @contact.name Giorgos
// @contact.url https://github.com/gosom/kouti/issues

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @accept json
// @produce json
// @query.collection.format multi
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
		SwaggerUI: &web.SwaggerUIConfig{
			SpecName: "TODO API",
			SpecFile: "/docs/swagger.json",
			Path:     "/docs",
			SpecFS:   specFs,
		},
	})
	if err != nil {
		return err
	}

	dsn := "postgres://postgres:secret@localhost:5432/todo?sslmode=disable&pool_max_conns=10"
	db, err := dbdriver.New(ctx, dbdriver.Config{
		Logger:     logger.NewSubLogger(log, "db"),
		ConnString: dsn,
	})
	if err != nil {
		return err
	}
	store := NewPostgresStore(logger.NewSubLogger(log, "store"), db)
	h := NewToDoHandler(log, store)
	router.Route("/todos", func(r chi.Router) {
		r.Post("/", h.Post)
		r.Get("/{id:[a-z0-9-]{36}}", h.Get)
		r.Put("/{id:[a-z0-9-]{36}}", h.Put)
		r.Delete("/{id:[a-z0-9-]{36}}", h.Delete)
		r.Get("/", h.Select)
	})

	return httpserver.Run(ctx, httpserver.Config{
		Host:   "127.0.0.1:8080",
		Router: router,
	})
}
