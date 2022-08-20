package utils

import (
	"os"
)

// Exit a struct that contains the ExitCode and message
type Exit struct {
	Code int
	Msg  string
}

// ExitRecover should be used with a defer statement
// fn is a function that accept a string. It should be used to print
// the error message
func ExitRecover(fn func(msg string)) {
	if e := recover(); e != nil {
		if exit, ok := e.(Exit); ok == true {
			if fn != nil {
				fn(exit.Msg)
			}
			os.Exit(exit.Code)
		}
		panic(e)
	}
	os.Exit(0)
}
