package main

import (
	"embed"

	"github.com/gosom/kouti/examples/simple-todo/config"
	"github.com/gosom/kouti/examples/simple-todo/rest"
	"github.com/gosom/kouti/httpserver"
	"github.com/gosom/kouti/utils"
)

//go:embed docs/swagger.json
var specFs embed.FS

func main() {
	defer utils.ExitRecover()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	srv, err := rest.InitializeServices(specFs, cfg)
	if err != nil {
		panic(err)
	}

	router, err := rest.NewRouter(srv)
	if err != nil {
		return err
	}

	srvCfg := httpserver.Config{
		Host:   cfg.ServerAddr,
		Router: router,
		UseTLS: cfg.UseTLS,
	}

	return httpserver.Run(srvCfg)
}
