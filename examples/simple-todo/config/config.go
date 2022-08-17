package config

import "github.com/gosom/kouti/utils"

type Config struct {
	ServerAddr string   `default:"localhost:8080"`
	UseTLS     bool     `default:"false"`
	Debug      bool     `default:"false"`
	Dsn        string   `default:"postgres://postgres:secret@localhost:5432/todo?sslmode=disable&pool_max_conns=10"`
	Roles      []string `default:"admin,member,visitor"`
	JwtSignKey string   `default:"secret"`
}

func New() (*Config, error) {
	cfg, err := utils.NewConfig[Config]("")
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
