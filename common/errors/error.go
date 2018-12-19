package errors

import (
	"errors"
	"runtime/debug"
)

func New(err string) error {
	debug.PrintStack()
	return errors.New(err)
}

func CheckErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckEqualPanic(b bool) {
	if !b {
		panic("not equal")
	}
}
