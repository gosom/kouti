package main

import "github.com/gosom/kouti/utils"

func main() {
	defer utils.ExitRecover()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return nil
}
