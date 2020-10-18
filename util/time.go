package util

import (
	"fmt"
	"time"
)

type echofunc func(s string, a ...interface{})
type Timer struct {
	tv1, tv2 time.Time
}

var echo echofunc = nil

const sec = float64(1000000000)

//========================================================================================

func Fsleep(length float64) {
	ilen := time.Duration(length * float64(time.Second))
	time.Sleep(ilen)
}

func Tdiff(tv1, tv2 *time.Time) float64 {
	nflt := (float64(tv2.Nanosecond()) - float64(tv1.Nanosecond())) / sec
	return (nflt + float64(tv2.Second()) - float64(tv1.Second()))
}

func SetEcho(f echofunc) {
	echo = f
}

//========================================================================================

func NewTimer() *Timer {
	var t Timer
	t.tv1 = time.Now()
	return &t
}

func (t *Timer) Update() {
	t.tv2 = time.Now()
}

func (t *Timer) Report(s string) string {
	t.tv2 = time.Now()
	return fmt.Sprintf("Time for \"%s\": %.10f", s, Tdiff(&t.tv1, &t.tv2))
}

func (t *Timer) Show(s string) string {
	if t.tv2.IsZero() {
		return "Timer: error, not updated"
	}
	return fmt.Sprintf("Running time for \"%s\": %.10fs", s, Tdiff(&t.tv1, &t.tv2))
}

func (t *Timer) Reset() {
	t.tv1 = time.Now()
	t.tv2 = time.Time{}
}

// Nvim specific

func (t *Timer) EchoReport(s string) {
	// if echo != nil {
	Eprintf("%s\n", t.Report(s))
	// }
}

func (t *Timer) EchoShow(s string) {
	// if echo != nil {
	Eprintf("%s\n", t.Show(s))
	// }
}

//========================================================================================

func TimeRoutine(s string, fn func(a ...interface{}) interface{}, args ...interface{}) interface{} {
	if echo == nil {
		panic("Cannot use TimeRoutine() when echo is nil.")
	}

	t := NewTimer()
	ret := fn(args...)
	echo("%s", t.Report(s))
	return ret
}
