// Package log provides logging and error reporting utility functions.
package log

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/fatih/color"

	"github.com/MathieuMoalic/amumax/src/fsutil"
)

var Log Logs

type Logs struct {
	Hist    string                   // console history for GUI
	logfile fsutil.WriteCloseFlusher // saves history of input commands +  output
	debug   bool
	path    string
}

func (l *Logs) AutoFlushToFile() {
	for {
		l.FlushToFile()
		time.Sleep(5 * time.Second)
	}
}

func (l *Logs) FlushToFile() {
	if l.logfile != nil {
		err := l.logfile.Flush()
		if err != nil {
			color.Red(fmt.Sprintf("Error flushing log file: %v", err))
		}
	}
}

func (l *Logs) SetDebug(debug bool) {
	l.debug = debug
}

func (l *Logs) Init(zarrPath string) {
	l.path = zarrPath + "/log.txt"
	l.createLogFile()
	l.writeToFile(l.Hist)
}

func (l *Logs) createLogFile() {
	var err error
	l.logfile, err = fsutil.Create(l.path)
	if err != nil {
		color.Red(fmt.Sprintf("Error creating the log file: %v", err))
	}
}

func (l *Logs) writeToFile(msg string) {
	if l.logfile == nil {
		return
	}
	_, err := l.logfile.Write([]byte(msg))
	if err != nil {
		if err.Error() == "short write" {
			color.Yellow("Error writing to log file, trying to recreate it...")
			l.createLogFile()
			_, _ = l.logfile.Write([]byte(msg))
		} else {
			color.Red(fmt.Sprintf("Error writing to log file: %v", err))
		}
	}
}

func (l *Logs) addAndWrite(msg string) {
	l.Hist += msg
	l.writeToFile(msg)
}

func (l *Logs) Command(msg ...any) {
	fmt.Println(fmt.Sprint(msg...))
	l.addAndWrite(fmt.Sprint(msg...) + "\n")
}

func (l *Logs) Info(msg string, args ...any) {
	formattedMsg := "// " + fmt.Sprintf(msg, args...) + "\n"
	color.Green(formattedMsg)
	l.addAndWrite(formattedMsg)
}

func (l *Logs) Warn(msg string, args ...any) {
	formattedMsg := "// " + fmt.Sprintf(msg, args...) + "\n"
	color.Yellow(formattedMsg)
	l.addAndWrite(formattedMsg)
}

func (l *Logs) Debug(msg string, args ...any) {
	if l.debug {
		formattedMsg := "// " + fmt.Sprintf(msg, args...) + "\n"
		color.Blue(formattedMsg)
		l.addAndWrite(formattedMsg)
	}
}

// Err prints an error message in red and adds it to the log history, does not exit or panic
func (l *Logs) Err(msg string, args ...any) {
	formattedMsg := "// " + fmt.Sprintf(msg, args...) + "\n"
	color.Red(formattedMsg)
	l.addAndWrite(formattedMsg)
}

// PanicIfError panics if err != nil, printing also the file and line number of the caller
func (l *Logs) PanicIfError(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		color.Red(fmt.Sprint("// ", file, ":", line, err) + "\n")
		panic(err)
	}
}

// ErrAndExit prints an error message in red, adds it to the log history, and exits with code 1
func (l *Logs) ErrAndExit(msg string, args ...any) {
	l.Err(msg, args...)
	os.Exit(1)
}

// AssertMsg Panics with msg if test is false
func (l *Logs) AssertMsg(test bool, msg any) {
	if !test {
		l.ErrAndExit("%v", msg)
	}
}

// AssertMsg Panics with msg if test is false
func AssertMsg(test bool, msg any) {
	if !test {
		Log.ErrAndExit("%v", msg)
	}
}
