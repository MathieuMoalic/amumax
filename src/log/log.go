package log

// Logging and error reporting utility functions

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/MathieuMoalic/amumax/src/fsutil"
	"github.com/fatih/color"
)

type Logs struct {
	Hist   string // console history for GUI
	debug  bool
	writer *bufio.Writer
	file   *os.File
}

func NewLogs(debug bool) *Logs {
	return &Logs{debug: debug}
}

func (l *Logs) InitLogs(fs *fsutil.FileSystem, debug bool) {
	writer, file, err := fs.Create("log.txt")
	if err != nil {
		color.Red(fmt.Sprintf("Error creating the log file: %v", err))
	}
	l.writer = writer
	l.file = file
	l.debug = debug
}

// print version to stdout
func (l *Logs) PrintVersion(version string, cudaInfo [6]string) {
	cudaVersion, cudaCC, GPUName, GPUMem, DriverVersion, GPUCC := cudaInfo[0], cudaInfo[1], cudaInfo[2], cudaInfo[3], cudaInfo[4], cudaInfo[5]
	l.Info("Version:         %s", version)
	l.Info("Platform:        %s_%s", runtime.GOOS, runtime.GOARCH)
	l.Info("Go Version:      %s (%s)", runtime.Version(), runtime.Compiler)
	l.Info("CUDA Version:    %s (CC=%s PTX)", cudaVersion, cudaCC)
	l.Info("GPU Information: %s(%s), CUDA Driver %s, cc=%s", GPUName, GPUMem, DriverVersion, GPUCC)
}
func (l *Logs) Close() error {
	err := l.FlushToFile()
	if err != nil {
		return fmt.Errorf("error flushing the log file: %v", err)
	}
	l.file.Close()
	return nil
}

func (l *Logs) FlushToFile() error {
	if l.writer == nil {
		return nil
	}
	err := l.writer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing the log file: %v", err)
	}
	return nil
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
	// We can write logs before the writer is initialized
	if l.writer != nil {
		l.writeToFile(msg)
	}
}

func (l *Logs) Command(cmd string) {
	// Find the position of the comment marker "//"
	commentIndex := strings.Index(cmd, "//")

	if commentIndex != -1 {
		// Split the string into command and comment parts
		commandPart := cmd[:commentIndex]
		commentPart := cmd[commentIndex:]

		// Print the command part normally
		fmt.Print(commandPart)

		// Print the comment part in green
		color.Green(commentPart)
	} else {
		// Print the whole command if no comment exists
		fmt.Println(cmd)
	}

	// Log the full command
	l.addAndWrite(cmd + "\n")
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
