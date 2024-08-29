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
	dev     bool
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

func (l *Logs) Init(zarrPath string, dev bool) {
	f, err := httpfs.Create(zarrPath + "/log.txt")
	if err != nil {
		color.Red(fmt.Sprintf("Error creating the log file: %v", err))
	}
	l.logfile = f // otherwise f gets dropped
	_, err = l.logfile.Write([]byte(l.Hist))
	if err != nil {
		color.Red(fmt.Sprintf("Error writing to log file: %v", err))
	}
	l.dev = dev
}

func (l *Logs) writeToFile(msg string) {
	if l.logfile != nil {
		_, err := l.logfile.Write([]byte(msg + "\n"))
		if err != nil {
			color.Red(fmt.Sprintf("Error writing to log file: %v", err))
		}
		l.Hist += msg + "\n"
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
	if l.dev {
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
