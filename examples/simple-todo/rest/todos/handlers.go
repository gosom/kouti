package todos

import (
	"github.com/gosom/kouti/web"
	"github.com/rs/zerolog"
)

type TodoQueryParams struct {
}

type Todo struct {
}

type TodoCreate struct {
}

func NewTodoHandler(log zerolog.Logger, srv *TodoSrv) web.ResourceHandler[TodoCreate, TodoQueryParams, Todo] {
	h := web.ResourceHandler[TodoCreate, TodoQueryParams, Todo]{}
	h.Logger = log
	h.Srv = srv
	return h
}
