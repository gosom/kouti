package main

import (
	"os"

	"github.com/gosom/kouti/um"
)

func main() {
	if err := um.CreateUmMigrationsFile(os.Args[1]); err != nil {
		panic(err)
	}
}
