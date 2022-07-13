package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/gosom/kouti/examples/todo/db"
	"github.com/gosom/kouti/logger"
	"github.com/gosom/kouti/web"
)

func NewRouter(db db.DB, cfg web.RouterConfig) (*chi.Mux, error) {
	router := web.NewRouter(cfg)
	{
		router.Route("/api/v1/users", func(r chi.Router) {
			us := UserSrv{
				Log: logger.NewSubLogger(cfg.Log, "UserSrv"),
				DB:  db,
			}
			h := NewUserHandler(
				logger.NewSubLogger(cfg.Log, "UserHandler"),
				&us,
			)

			r.Post("/", h.Post)
			r.Get("/", h.Select)
			r.Get("/search", h.Search)
			r.Get(`/{id:\d+}`, h.Get)
			r.Delete(`/{id:\d+}`, h.Delete)
		})

		router.Route("/api/v1/todo", func(r chi.Router) {
		})
	}
	return router, nil
}
