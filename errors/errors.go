package db

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/pkg/errors"
)

// WrapErr for creating errors and wrapping them along
// with callers info
func WrapErr(e interface{}, p ...interface{}) error {
	if e == nil {
		return nil
	}

	var err error

	switch e := e.(type) {
	case string:
		err = fmt.Errorf(e, p...)
	case error:
		err = e
	}

	pc, _, no, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return errors.Wrap(err, fmt.Sprintf("%s#%s\n", details.Name(), strconv.Itoa(no)))
	}

	return err
}
