package util

// Logging and error reporting utility functions

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/fatih/color"
)

func Fatal(msg ...interface{}) {
	log.Fatal(msg...)
}

func Fatalf(format string, msg ...interface{}) {
	log.Fatalf(format, msg...)
}

func FatalErr(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		color.Red(fmt.Sprint("// ", file, ":", line, err))
		os.Exit(1)
	}
}

func PanicErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func LogErr(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		color.Red(fmt.Sprint("// ", file, ":", line, err))
	}
}

func LogThenExit(msg string) {
	color.Red(msg)
	os.Exit(1)
}

func Log(msg ...interface{}) {
	color.Green("// " + fmt.Sprint(msg...))
}

// Panics with "illegal argument" if test is false.
func Argument(test bool) {
	if !test {
		log.Panic("illegal argument")
	}
}

// Panics with msg if test is false
func AssertMsg(test bool, msg interface{}) {
	if !test {
		log.Panic(msg)
	}
}

// Panics with "assertion failed" if test is false.
func Assert(test bool) {
	if !test {
		log.Panic("assertion failed")
	}
}

// Hack to avoid cyclic dependency on engine.
var (
	progress_ func(int, int, string) = PrintProgress
	progLock  sync.Mutex
)

// Set progress bar to progress/total and display msg
// if GUI is up and running.
func Progress(progress, total int, msg string) {
	progLock.Lock()
	defer progLock.Unlock()
	if progress_ != nil {
		progress_(progress, total, msg)
	}
}

var (
	lastPct   = -1      // last progress percentage shown
	lastProgT time.Time // last time we showed progress percentage
)

func PrintProgress(prog, total int, msg string) {
	pct := (prog * 100) / total
	if pct != lastPct { // only print percentage if changed
		if (time.Since(lastProgT) > time.Second) || pct == 100 { // only print percentage once/second unless finished
			fmt.Println("//", msg, pct, "%")
			lastPct = pct
			lastProgT = time.Now()
		}
	}
}

// Sets the function to be used internally by Progress.
// Avoids cyclic dependency on engine.
func SetProgress(f func(int, int, string)) {
	progLock.Lock()
	defer progLock.Unlock()
	progress_ = f
}
