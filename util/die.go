package util

import (
	"fmt"
	"os"
)

var prefix string = os.Args[0] + ": "

func Die(status int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, prefix+format+"\n", args...)
	os.Exit(status)
}

func Warn(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, prefix+format+"\n", args...)
}

func DieE(status int, err error) {
	Die(status, "%+v", err)
}

func PanicE(err error) {
	panic(Unwrap(err))
}

func GetPrefix() string {
	return prefix
}

func SetPrefix(s string) {
	prefix = s
}
