package auth

import (
	"io"

	"github.com/casbin/casbin/v2/persist"
	"github.com/rs/zerolog"
)

type AuthorizatorConfig struct {
	CasbinModelReader io.Reader
	CasbinAdapter     persist.Adapter
	Log               zerolog.Logger
}

type Authorizator interface {
	//Authorize(identity any, r *http.Request) error
	HasPermission(identity any, action string, asset string) error
}

func NewAuthorizator(cfg AuthorizatorConfig) (Authorizator, error) {
	return newCasbinAuthorizator(cfg.Log, cfg.CasbinAdapter, cfg.CasbinModelReader)
}
