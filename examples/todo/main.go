package main

import (
	"context"

	"github.com/gosom/kouti/dbdriver"
	"github.com/gosom/kouti/httpserver"
	"github.com/gosom/kouti/logger"
	"github.com/gosom/kouti/utils"
	"github.com/gosom/kouti/web"

	"github.com/gosom/kouti/examples/todo/db"
	"github.com/gosom/kouti/examples/todo/rest"
)

func main() {
	defer utils.ExitRecover()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	ctx := context.Background()
	log := logger.New(logger.Config{Debug: true})

	dsn := "postgres://postgres:secret@localhost:5432/todo?sslmode=disable&pool_max_conns=10"
	dbconn, err := db.New(ctx, dbdriver.Config{
		Logger:     logger.NewSubLogger(log, "DB"),
		ConnString: dsn,
	})
	if err != nil {
		return err
	}
	defer dbconn.Close()

	router, err := rest.NewRouter(dbconn, web.RouterConfig{
		Log: log,
	})
	if err != nil {
		return err
	}
	cfg := httpserver.Config{
		Host:   "localhost:8080",
		Router: router,
		UseTLS: false,
	}

	return httpserver.Run(cfg)
}
