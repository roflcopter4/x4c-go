package util

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/roflcopter4/x4c-go/util/color"
)

var prefix string

func init() {
	progname := filepath.Base(os.Args[0])
	if progname != "" {
		prefix = color.BYellow(progname+":") + " "
	}
}

func Die(status int, format string, args ...interface{}) {
	format = prefix + color.BRed("Error: ") + format + "\n"
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(status)
}

func Warn(format string, args ...interface{}) {
	format = prefix + color.BOrange("Warning: ") + format + "\n"
	fmt.Fprintf(os.Stderr, format, args...)
}

func DieE(status int, err error) {
	Die(status, "%+v", err)
}

func PanicUnwrap(err error) {
	panic(Unwrap(err))
}

func GetPrefix() string {
	return prefix
}

func SetPrefix(s string) {
	prefix = s
}
