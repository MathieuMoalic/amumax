package util

// Logging and error reporting utility functions

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/MathieuMoalic/amumax/httpfs"
	"github.com/fatih/color"
)

var Log Logs

type Logs struct {
	Hist    string                   // console history for GUI
	logfile httpfs.WriteCloseFlusher // saves history of input commands +  output
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
		l.logfile.Flush()

	}
}

func (l *Logs) Init(zarrPath string, debug bool) {
	l.path = zarrPath + "/log.txt"
	l.debug = debug
	l.createLogFile()
	l.writeToFile(l.Hist)
}

func (l *Logs) createLogFile() {
	var err error
	l.logfile, err = httpfs.Create(l.path)
	if err != nil {
		color.Red(fmt.Sprintf("Error creating the log file: %v", err))
	}
}

func (l *Logs) writeToFile(msg string) {
	l.Hist += msg + "\n"
	if l.logfile == nil {
		return
	}
	_, err := l.logfile.Write([]byte(msg + "\n"))
	if err != nil {
		if err.Error() == "short write" {
			color.Yellow("Error writing to log file, trying to recreate it...")
			l.createLogFile()
			_, _ = l.logfile.Write([]byte(msg + "\n"))
		} else {
			color.Red(fmt.Sprintf("Error writing to log file: %v", err))
		}
	}
}

func (l *Logs) Command(msg ...interface{}) {
	fmt.Println(fmt.Sprint(msg...))
	l.writeToFile(fmt.Sprint(msg...))
}

func (l *Logs) Comment(msg string, args ...interface{}) {
	formattedMsg := "// " + fmt.Sprintf(msg, args...)
	color.Green(formattedMsg)
	l.writeToFile(formattedMsg)
}

func (l *Logs) Warn(msg string, args ...interface{}) {
	formattedMsg := "// " + fmt.Sprintf(msg, args...)
	color.Yellow(formattedMsg)
	l.writeToFile(formattedMsg)
}

func (l *Logs) Debug(msg string, args ...interface{}) {
	if l.debug {
		formattedMsg := "// " + fmt.Sprintf(msg, args...)
		color.Blue(formattedMsg)
		l.writeToFile(formattedMsg)
	}
}

func (l *Logs) Err(msg string, args ...interface{}) {
	formattedMsg := "// " + fmt.Sprintf(msg, args...)
	color.Red(formattedMsg)
	l.writeToFile(formattedMsg)
}

func (l *Logs) PanicIfError(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		color.Red(fmt.Sprint("// ", file, ":", line, err))
		panic(err)
	}
}

func (l *Logs) ErrAndExit(msg string, args ...interface{}) {
	l.Err(msg, args...)
	os.Exit(1)
}
