package rest

import (
	"context"
	"embed"

	"github.com/go-chi/chi/v5"
	"github.com/gosom/kouti/dbdriver"
	"github.com/gosom/kouti/logger"
	"github.com/gosom/kouti/um"
	"github.com/gosom/kouti/web"

	"github.com/gosom/kouti/examples/simple-todo/config"
	"github.com/gosom/kouti/examples/simple-todo/rest/todos"
	"github.com/gosom/kouti/examples/simple-todo/rest/users"
)

type Services struct {
	DB      dbdriver.DB
	WebConf web.RouterConfig
	UserSrv *users.UserSrv
	TodoSrv *todos.TodoSrv
}

func InitializeServices(specFs embed.FS, cfg *config.Config) (*Services, error) {
	log := logger.New(logger.Config{Debug: true})
	var (
		ans Services
		err error
	)
	ans.DB, err = dbdriver.New(context.Background(), dbdriver.Config{
		Logger:     logger.NewSubLogger(log, "db"),
		ConnString: cfg.Dsn,
	})
	if err != nil {
		return nil, err
	}
	ans.WebConf = web.RouterConfig{
		Log: logger.NewSubLogger(log, "router"),
		SwaggerUI: &web.SwaggerUIConfig{
			SpecName: "TODO API",
			SpecFile: "/docs/swagger.json",
			Path:     "/docs",
			SpecFS:   specFs,
		},
	}

	ans.UserSrv = &users.UserSrv{
		Log: logger.NewSubLogger(log, "users"),
		Srv: um.NewService(um.Config{
			Log:         logger.NewSubLogger(log, "um"),
			DB:          ans.DB,
			SystemRoles: cfg.Roles,
		}),
	}

	ans.TodoSrv = &todos.TodoSrv{
		Log: logger.NewSubLogger(log, "todos"),
	}

	return &ans, nil
}

// @title Todo API based on kouti
// @version 0.1
// @description This is a sample server todo server.
// @description You can visit the GitHub repository at https://github.com/gosom/kouti

// @contact.name Giorgos
// @contact.url https://github.com/gosom/kouti/issues

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1/
// @accept json
// @produce json
// @query.collection.format multi
func NewRouter(srv *Services) (*chi.Mux, error) {
	router, err := web.NewRouter(srv.WebConf)
	if err != nil {
		return nil, err
	}

	routes(router, srv)

	return router, nil
}

func routes(router *chi.Mux, srv *Services) {

	uh := users.NewUserHandler(srv.WebConf.Log, srv.UserSrv)
	lh := users.NewAuthHandler(srv.WebConf.Log, srv.UserSrv)

	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			// public
			r.Group(func(r chi.Router) {
				r.Post("/", uh.Post)
				r.Post("/login", lh.Login)
			})

		})
	})
}
