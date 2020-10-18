package util

import (
	"fmt"

	"github.com/pkg/errors"
)

type StackTracer interface {
	StackTrace() errors.StackTrace
}

func Unwrap(err_in error) (ret string) {
	if err, ok := err_in.(StackTracer); ok {
		for _, f := range err.StackTrace() {
			ret += fmt.Sprintf("%+s:%d\n", f, f)
		}
		ret = ret[:len(ret)-1]
	} else {
		ret = fmt.Sprint(err_in)
	}
	return
}
