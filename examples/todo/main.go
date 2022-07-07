package main

import (
	"github.com/gosom/kouti/httpserver"
	"github.com/gosom/kouti/utils"
)

func main() {
	defer utils.ExitRecover()
	cfg := httpserver.Config{
		Host:   "localhost:8080",
		UseTLS: false,
	}

	if err := httpserver.Run(cfg); err != nil {
		panic(utils.Exit{1, err.Error()})
	}
}
