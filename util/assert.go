package util

import (
	"fmt"
)

// Panics with "illegal argument" if test is false.
func Argument(test bool) {
	if !test {
		Log.PanicIfError(fmt.Errorf("illegal argument"))
	}
}

// Panics with msg if test is false
func AssertMsg(test bool, msg interface{}) {
	if !test {
		Log.PanicIfError(fmt.Errorf("%v", msg))
	}
}

// Panics with "assertion failed" if test is false.
func Assert(test bool) {
	if !test {
		Log.PanicIfError(fmt.Errorf("assertion failed"))
	}
}

// Special error that is not fatal when paniced on and called from GUI
// E.g.: try to set bad grid size: panic on UserErr, recover, print error, carry on.
type UserErr string

func (e UserErr) Error() string { return string(e) }

func CheckRecoverable(err error) {
	if err != nil {
		panic(UserErr(err.Error()))
	}
}
