package auth

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

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

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
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
func (a *casbinAuthorizator) Authorize(identity any, r *http.Request) error {
	// TODO may use select with type ?
	sub := fmt.Sprintf("%v", identity)
	a.log.Info().Msgf("authorize %s", sub)
	// TODO Replace with filtered policy based on request ?
	if err := a.e.LoadPolicy(); err != nil {
		a.log.Error().Err(err)
		return err
	}
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
	return nil
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
