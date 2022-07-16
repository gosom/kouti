package rest

import (
	"github.com/go-chi/chi/v5"

	"github.com/gosom/kouti/auth"
	"github.com/gosom/kouti/casbinpgx"
	"github.com/gosom/kouti/logger"
	"github.com/gosom/kouti/web"

	"github.com/gosom/kouti/examples/todo/db"
)

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
func NewRouter(db db.DB, cfg web.RouterConfig) (*chi.Mux, error) {
	router, err := web.NewRouter(cfg)
	if err != nil {
		return nil, err
	}
	authenticator, err := auth.New(auth.AuthenticatorConfig{
		JwtSignKey: "secret",
		Log:        logger.NewSubLogger(cfg.Log, "Authenticator"),
	})
	if err != nil {
		return nil, err
	}
	ad, err := casbinpgx.NewAdapter(casbinpgx.Config{
		DbConn: db.RawConn(),
	})
	if err != nil {
		return nil, err
	}
	authorizator, err := auth.NewAuthorizator(
		auth.AuthorizatorConfig{
			CasbinModelReader: nil,
			CasbinAdapter:     ad,
			Log:               logger.NewSubLogger(cfg.Log, "Authorizator"),
		},
	)
	if err != nil {
		return nil, err
	}
	{
		us := UserSrv{
			Log:  logger.NewSubLogger(cfg.Log, "UserSrv"),
			DB:   db,
			Auth: authenticator,
		}
		h := NewUserHandler(
			logger.NewSubLogger(cfg.Log, "UserHandler"),
			&us,
		)
		lh := NewAuthHandler(
			logger.NewSubLogger(cfg.Log, "AuthHandler"),
			&us,
		)

		// User Routes
		router.Route("/api/v1/users", func(r chi.Router) {
			// public
			r.Post("/", h.Post)
			r.Post("/login", lh.Login)

			// logged in
			r.Group(func(r chi.Router) {
				r.Use(web.Authentication(authenticator))
				r.Use(web.Authorization(authorizator))

				r.Get("/", h.Select)
				r.Get(`/{id:\d+}`, h.Get)
				r.Get("/search", h.Search)
				r.Delete(`/{id:\d+}`, h.Delete)

			})

		})

		router.Route("/api/v1/todo", func(r chi.Router) {
		})
	}
	return router, nil
}
