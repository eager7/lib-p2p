package errors

import (
	"github.com/eager7/go/mlog"
	"errors"
)

func New(log mlog.Logger, err string) error {
	log.ErrStack()
	return errors.New(err)
}

func CheckErrorPanic(err error) {
	if err != nil {
		mlog.L.Panic(err)
	}
}

func CheckEqualPanic(b bool) {
	if !b {
		mlog.L.Panic("not equal")
	}
}
