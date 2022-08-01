package main

import (
	"context"
	"embed"
	"fmt"

	"github.com/gosom/kouti/dbdriver"
	"github.com/gosom/kouti/logger"
	"github.com/gosom/kouti/um"
	"github.com/gosom/kouti/utils"

	"github.com/gosom/kouti/examples/todo/db"
)

//go:embed docs/swagger.json
var specFs embed.FS

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

	// we first setup user mananagement service

	users := um.NewService(um.Config{
		Log:         logger.NewSubLogger(log, "Users"),
		DB:          dbconn.DB,
		SystemRoles: []string{"admin", "member"},
	})

	if err := users.InitSchema(ctx); err != nil {
		return err
	}

	rp := um.RegisterUserOpts{
		Identity: "giorgos@example.com",
		Passwd:   "123abc!",
		Roles:    []string{"admin"},

		BeforeCommit: func(ctx context.Context, conn dbdriver.DBTX, u um.User) error {
			fmt.Println("before commit")
			return nil
		},
		AfterCommit: func(ctx context.Context, conn dbdriver.DBTX, u um.User) error {
			fmt.Println("after commit")
			return nil
		},
	}

	_, err, afterHookErr := users.RegisterUser(ctx, rp)
	if err != nil {
		return err
	}
	if afterHookErr != nil {
		return err
	}

	return nil

	/*
		router, err := rest.NewRouter(dbconn, web.RouterConfig{
			Log: log,
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
		cfg := httpserver.Config{
			Host:   "localhost:8080",
			Router: router,
			UseTLS: false,
		}

		return httpserver.Run(cfg)
	*/
}
