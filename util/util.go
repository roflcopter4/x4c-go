package util

// /*
// #include <stdint.h>
// static double
// ldouble_wrap(const uint64_t nsec1, const uint64_t nsec2)
// {
//         static const long double sec = 1000000000.0L;
//         long double              x   = ((long double)nsec2) -
//                                        ((long double)nsec1);
//         return (double)(x / sec);
// }
// */
// import "C"

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

var Logfiles map[string]*os.File

func FuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	elems := strings.Split(fn.Name(), ".")
	return elems[len(elems)-1]
}

func Eprint(a ...interface{})                 { fmt.Fprint(os.Stderr, a...) }
func Eprintln(a ...interface{})               { fmt.Fprintln(os.Stderr, a...) }
func Eprintf(format string, a ...interface{}) { fmt.Fprintf(os.Stderr, format, a...) }

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func BoolInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

// func Assert(cond bool, mes string, a ...interface{}) {
func Assert(cond bool, a ...interface{}) {
	if !cond {
		var msg string
		if _, fileName, fileLine, ok := runtime.Caller(1); ok {
			msg = fmt.Sprintf("Assertion failed at (%s: %d): ", fileName, fileLine)
		} else {
			msg = "Assertion failed: "
		}
		panic(msg + fmt.Sprint(a...))
	}
}

func SafeFopen(fname string, mode, perm int) *os.File {
	file, e := os.OpenFile(fname, mode, os.FileMode(perm))
	if e != nil {
		panic(e)
	}
	return file
}

func QuickRead(filename string) []byte {
	var (
		buf  bytes.Buffer
		file *os.File
		e    error
	)

	if file, e = os.Open(filename); e != nil {
		return nil
	}

	if _, e = buf.ReadFrom(file); e != nil {
		log.Panicf("Unexpected read error: %v\n", e)
	}

	file.Close()
	return buf.Bytes()
}

func UniqueStr(strlist []string) []string {
	keys := make(map[string]bool)
	ret := []string{}

	for _, entry := range strlist {
		if entry == "" {
			continue
		}
		if _, value := keys[entry]; !value {
			keys[entry] = true
			ret = append(ret, entry)
		}
	}

	return ret
}

func StrEqAny(cmp string, lst ...string) bool {
	for _, s := range lst {
		if cmp == s {
			return true
		}
	}
	return false
}

func StrEqAll(cmp string, lst ...string) bool {
	for _, s := range lst {
		if cmp != s {
			return false
		}
	}
	return true
}

func IntIsBetween(i, low, high int) bool {
	return i >= low && i < high
}
