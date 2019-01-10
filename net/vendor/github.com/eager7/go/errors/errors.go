package errors

import (
	"github.com/eager7/go/log"
	"errors"
)

func New(log log.Logger, err string) error {
	log.ErrStack(err)
	return errors.New(err)
}

func CheckErrorPanic(err error) {
	if err != nil {
		log.L.Panic(err)
	}
}

func CheckEqualPanic(b bool) {
	if !b {
		log.L.Panic("not equal")
	}
}