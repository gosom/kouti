package casbinpgx

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	DbConn          *pgxpool.Pool
	TableName       string
	Timeout         time.Duration
	SkipCreateTable bool
}

func (o *Config) setDefaults() {
	if o.TableName == "" {
		o.TableName = "casbin"
	}
	if o.Timeout == 0 {
		o.Timeout = time.Second * 5
	}
}

// https://github.com/casbin/xorm-adapter/blob/master/adapter.go#L251
type CasbinRule struct {
	PType string
	V0    string
	V1    string
	V2    string
	V3    string
	V4    string
	V5    string
}

type Adapter struct {
	dbconn    *pgxpool.Pool
	tableName string
	timeout   time.Duration
	loadQ     string
	truncQ    string
}

func NewAdapter(cfg Config) (*Adapter, error) {
	cfg.setDefaults()
	ans := Adapter{
		dbconn:    cfg.DbConn,
		tableName: cfg.TableName,
		timeout:   cfg.Timeout,
		loadQ: fmt.Sprintf(
			`SELECT ptype, v0, v1, v2, v3, v4, v5 FROM %s`,
			cfg.TableName,
		),
		truncQ: fmt.Sprintf(
			`TRUNCATE %s WITH RESTART IDENTITY`, cfg.TableName,
		),
	}
	if !cfg.SkipCreateTable {
		if err := ans.createTable(); err != nil {
			return &ans, err
		}
	}

	return &ans, nil
}

func (a *Adapter) LoadPolicy(model model.Model) error {
	// unfortunately due to the interface we can't pass context
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()
	var rule CasbinRule
	_, err := a.dbconn.QueryFunc(
		ctx,
		a.loadQ,
		nil,
		[]interface{}{&rule.PType, &rule.V0, &rule.V1, &rule.V2, &rule.V3, &rule.V4, &rule.V5},
		func(pgx.QueryFuncRow) error {
			loadPolicyLine(&rule, model)
			return nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (a *Adapter) SavePolicy(model model.Model) error {
	lines := make([][]any, 0, 64)
	id := 0
	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			id++
			line := a.genPolicyLine(id, ptype, rule)
			lines = append(lines, line)
		}
	}
	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			id++
			line := a.genPolicyLine(id, ptype, rule)
			lines = append(lines, line)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()
	tx, err := a.dbconn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, a.truncQ); err != nil {
		return err
	}

	if len(lines) > 0 {
		_, err = tx.CopyFrom(
			ctx,
			pgx.Identifier{a.tableName},
			[]string{"id", "ptype", "v0", "v1", "v2", "v3", "v4", "v5"},
			pgx.CopyFromRows(lines),
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// AddPolicy adds a policy rule to the storage.
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemovePolicy removes a policy rule from the storage.
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return errors.New("not implemented")
}

func (a *Adapter) createTable() error {
	q := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (
				id int PRIMARY KEY,
				ptype text NOT NULL,
				v0 text NOT NULL,
				v1 text NOT NULL,
				v2 text NOT NULL,
				v3 text NOT NULL,
				v4 text NOT NULL,
				v5 text NOT NULL
			)`, a.tableName)
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()
	_, err := a.dbconn.Exec(ctx, q)
	if err != nil {
		return err
	}

	return nil
}

// from https://github.com/casbin/xorm-adapter
func (a *Adapter) genPolicyLine(id int, ptype string, rule []string) []any {
	line := make([]any, 8)
	line[0] = id
	line[1] = ptype
	l := len(rule)
	if l > 0 {
		line[2] = rule[0]
	}
	if l > 1 {
		line[3] = rule[1]
	}
	if l > 2 {
		line[4] = rule[2]
	}
	if l > 3 {
		line[5] = rule[3]
	}
	if l > 4 {
		line[6] = rule[4]
	}
	if l > 5 {
		line[7] = rule[5]
	}

	return line
}

// from https://github.com/casbin/xorm-adapter
func loadPolicyLine(line *CasbinRule, model model.Model) {
	var p = []string{line.PType,
		line.V0, line.V1, line.V2, line.V3, line.V4, line.V5}
	var lineText string
	if line.V5 != "" {
		lineText = strings.Join(p, ", ")
	} else if line.V4 != "" {
		lineText = strings.Join(p[:6], ", ")
	} else if line.V3 != "" {
		lineText = strings.Join(p[:5], ", ")
	} else if line.V2 != "" {
		lineText = strings.Join(p[:4], ", ")
	} else if line.V1 != "" {
		lineText = strings.Join(p[:3], ", ")
	} else if line.V0 != "" {
		lineText = strings.Join(p[:2], ", ")
	}

	persist.LoadPolicyLine(lineText, model)
}
