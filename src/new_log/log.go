package new_log

// Logging and error reporting utility functions

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/MathieuMoalic/amumax/src/new_fsutil"
	"github.com/fatih/color"
)

type Logs struct {
	Hist   string // console history for GUI
	debug  bool
	writer *bufio.Writer
	file   *os.File
}

func NewLogs(zarrPath string, fs *new_fsutil.FileSystem, debug bool) *Logs {
	l := &Logs{}
	writer, file, err := fs.Create(zarrPath + "/log.txt")
	if err != nil {
		color.Red(fmt.Sprintf("Error creating the log file: %v", err))
	}
	l.writer = writer
	l.file = file
	l.debug = debug
	return l
}

func (l *Logs) Close() {
	l.FlushToFile()
	l.file.Close()
}
func (l *Logs) AutoFlushToFile() {
	for {
		l.FlushToFile()
		time.Sleep(5 * time.Second)
	}
}

func (l *Logs) FlushToFile() {
	l.writer.Flush()
}

func (l *Logs) writeToFile(msg string) {
	n, err := l.writer.WriteString(msg)
	if err != nil {
		color.Red(fmt.Sprintf("Error writing to the log file: %v", err))
	}
	if n != len(msg) {
		color.Red(fmt.Sprintf("Error writing to the log file: %v", err))
	}
}

func (l *Logs) addAndWrite(msg string) {
	l.Hist += msg
	l.writeToFile(msg)
}

func (l *Logs) Command(msg ...interface{}) {
	fmt.Println(fmt.Sprint(msg...))
	l.addAndWrite(fmt.Sprint(msg...) + "\n")
}

func (l *Logs) Info(msg string, args ...interface{}) {
	formattedMsg := "// " + fmt.Sprintf(msg, args...) + "\n"
	color.Green(formattedMsg)
	l.addAndWrite(formattedMsg)
}

func (l *Logs) Warn(msg string, args ...interface{}) {
	formattedMsg := "// " + fmt.Sprintf(msg, args...) + "\n"
	color.Yellow(formattedMsg)
	l.addAndWrite(formattedMsg)
}

func (l *Logs) Debug(msg string, args ...interface{}) {
	if l.debug {
		formattedMsg := "// " + fmt.Sprintf(msg, args...) + "\n"
		color.Blue(formattedMsg)
		l.addAndWrite(formattedMsg)
	}
}

func (l *Logs) Err(msg string, args ...interface{}) {
	formattedMsg := "// " + fmt.Sprintf(msg, args...) + "\n"
	color.Red(formattedMsg)
	l.addAndWrite(formattedMsg)
}

func (l *Logs) PanicIfError(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		color.Red(fmt.Sprint("// ", file, ":", line, err) + "\n")
		panic(err)
	}
}

func (l *Logs) ErrAndExit(msg string, args ...interface{}) {
	l.Err(msg, args...)
	os.Exit(1)
}

// Panics with msg if test is false
func (l *Logs) AssertMsg(test bool, msg interface{}) {
	if !test {
		l.ErrAndExit("%v", msg)
	}
}
