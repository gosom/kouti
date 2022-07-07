package utils

import (
	"os"
)

//
type Exit struct {
	Code int
	Msg  string
}

// ExitRecover should be used with a defer statement
// See examples/todo/main.go for a sample usage
func ExitRecover() {
	if e := recover(); e != nil {
		if exit, ok := e.(Exit); ok == true {
			os.Exit(exit.Code)
		}
		panic(e)
	}
	os.Exit(0)
}
