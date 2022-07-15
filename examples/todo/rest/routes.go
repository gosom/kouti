package rest

import (
	"github.com/go-chi/chi/v5"

	"github.com/gosom/kouti/auth"
	"github.com/gosom/kouti/casbinpgx"
	"github.com/gosom/kouti/logger"
	"github.com/gosom/kouti/web"

	"github.com/gosom/kouti/examples/todo/db"
)

func NewRouter(db db.DB, cfg web.RouterConfig) (*chi.Mux, error) {
	router := web.NewRouter(cfg)
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
			Log: logger.NewSubLogger(cfg.Log, "UserSrv"),
			DB:  db,
		}
		h := NewUserHandler(
			logger.NewSubLogger(cfg.Log, "UserHandler"),
			&us,
		)

		// Admin routes
		router.Route("/api/v1/admin", func(r chi.Router) {
			router.Route("/users", func(r chi.Router) {
				r.Get("/", h.Select)
				r.Get("/search", h.Search)
			})
		})

		// User Routes
		router.Route("/api/v1/users", func(r chi.Router) {
			// public
			r.Post("/", h.Post)

			// logged in
			r.Group(func(r chi.Router) {
				r.Use(web.Authentication(authenticator))
				r.Use(web.Authorization(authorizator))
				r.Get(`/{id:\d+}`, h.Get)
				r.Delete(`/{id:\d+}`, h.Delete)
			})

		})

		router.Route("/api/v1/todo", func(r chi.Router) {
		})
	}
	return router, nil
}
