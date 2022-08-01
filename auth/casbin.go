package auth

import (
	"errors"
	"io"
	"io/ioutil"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/rs/zerolog"
)

// See example for REST policy here
// https://github.com/casbin/casbin/blob/master/examples/keymatch_policy.csv

// https://github.com/casbin/casbin/blob/master/examples/keymatch_model.conf
const defaultCasbinModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = (g(r.sub, p.sub) || keyMatch(r.sub, p.sub)) && keyMatch(r.obj, p.obj) && keyMatch(r.act, p.act)
`

func NewEnforcer(r io.Reader, a persist.Adapter) (*casbin.Enforcer, error) {
	m, err := getCasbinModel(r)
	if err != nil {
		return nil, err
	}
	enforcer, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, err
	}
	return enforcer, err
}

func getCasbinModel(r io.Reader) (model.Model, error) {
	var casbinPolicy string
	if r == nil {
		casbinPolicy = defaultCasbinModel
	} else {
		policyBytes, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		casbinPolicy = string(policyBytes)
	}
	return model.NewModelFromString(casbinPolicy)
}

type casbinAuthorizator struct {
	log zerolog.Logger
	e   *casbin.Enforcer
}

// TODO perhaps use generics
func (a *casbinAuthorizator) HasPermission(identity any, action string, asset string) error {
	roles := []string{"visitor"}
	if identity != nil {
		roles = append(roles, "member")
	}
	a.log.Info().Msgf("authorize %v %s %s %v", identity, action, asset, roles)
	if err := a.e.LoadPolicy(); err != nil {
		a.log.Error().Err(err)
		return err
	}
	for _, role := range roles {
		ok, err := a.e.Enforce(role, asset, action)
		if err != nil {
			a.log.Error().Err(err).Msg("err")
			return err
		}
		if ok {
			return nil
		}
	}
	a.log.Info().Msgf("not authorized")
	return errors.New("unauthorized")
	// TODO may use select with type ?
	/*
		sub := fmt.Sprintf("%v", identity)
		a.log.Info().Msgf("authorize %s", sub)
		// TODO Replace with filtered policy based on request ?
		ok, err := a.e.Enforce(sub, r.URL.Path, r.Method)
		if err != nil {
			a.log.Error().Err(err).Msg("err")
			return err
		}
		if !ok {
			err := fmt.Errorf("%s not authorized to access %s %s",
				sub, r.URL.Path, r.Method,
			)
			a.log.Error().Err(err)
			return err
		}
	*/
}

func newCasbinAuthorizator(log zerolog.Logger, a persist.Adapter, r io.Reader) (*casbinAuthorizator, error) {
	ans := casbinAuthorizator{
		log: log,
	}
	var err error
	ans.e, err = NewEnforcer(r, a)
	if err != nil {
		return &ans, err
	}
	return &ans, nil
}
