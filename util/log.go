package util

// Logging and error reporting utility functions

import (
	"fmt"
	"log"
	"os"
	"runtime"

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

func LogThenExit(msg ...interface{}) {
	color.Red("// " + fmt.Sprint(msg...))
	os.Exit(1)
}

func LogErr(msg ...interface{}) {
	color.Red("// " + fmt.Sprint(msg...))
}

func Log(msg ...interface{}) {
	color.Green("// " + fmt.Sprint(msg...))
}

func LogWarn(msg ...interface{}) {
	color.Yellow("// " + fmt.Sprint(msg...))
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
