package util

// [>
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

func FuncName() []byte {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	elems := strings.Split(fn.Name(), ".")
	return []byte(elems[len(elems)-1])
}

func Eprintf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

// func Warn(format string, a ...interface{}) {
//       fmt.Fprintf(os.Stderr, "WARNING: "+format, a...)
// }

func Max_Int(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min_Int(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Safe_Fopen(fname string, mode, perm int) *os.File {
	file, e := os.OpenFile(fname, mode, os.FileMode(perm))
	if e != nil {
		panic(e)
	}
	return file
}

func Assert(cond bool, mes string, a ...interface{}) {
	if !cond {
		panic(fmt.Sprintf(mes, a...))
	}
}

func Boolint(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func Quick_Read(filename string) []byte {
	// st, e := os.Stat(filename)
	// if e != nil {
	//         return nil
	// }
	//
	// var (
	//         ret  = make([]byte, 0, st.Size())
	//         file *os.File
	//         n    int
	// )

	var (
		buf  bytes.Buffer
		file *os.File
		e    error
	)

	if file, e = os.Open(filename); e != nil {
		return nil
	}
	/* if n, e = file.Read(ret); e != nil || int64(n) != st.Size() {
		log.Panicf("Unexpected io error: %s, (n=%d, size=%d)", e, n, st.Size())
	} */

	if _, e = buf.ReadFrom(file); e != nil {
		log.Panicf("Unexpected read error: %v\n", e)
	}

	file.Close()
	return buf.Bytes()
}

func Unique_Str(strlist []string) []string {
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
