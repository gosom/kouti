package main

import (
	"context"

	"github.com/gosom/kouti/dbdriver"
	"github.com/gosom/kouti/examples/todo/db"
	"github.com/gosom/kouti/examples/todo/rest"
	"github.com/gosom/kouti/logger"
	"github.com/gosom/kouti/utils"
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

	us := rest.UserSrv{
		Log: logger.NewSubLogger(log, "UserSrv"),
		DB:  dbconn,
	}

	if err := us.InsertRandomUsers(ctx, 10000); err != nil {
		return err
	}
	return nil
}
