package engine

import (
	"github.com/MathieuMoalic/amumax/httpfs"
)

var (
	Hist    string                   // console history for GUI
	logfile httpfs.WriteCloseFlusher // saves history of input commands +  output
)

// Special error that is not fatal when paniced on and called from GUI
// E.g.: try to set bad grid size: panic on UserErr, recover, print error, carry on.
type UserErr string

func (e UserErr) Error() string { return string(e) }

func CheckRecoverable(err error) {
	if err != nil {
		panic(UserErr(err.Error()))
	}
}

// func LogExpression(msg ...interface{}) {
// 	str := sprint(msg...)
// 	log2GUI(str)
// 	log2File(str)
// 	fmt.Println(str)
// }

// func util.Log.Comment(msg ...interface{}) {
// 	str := "// " + sprint(msg...)
// 	log2GUI(str)
// 	log2File(str)
// 	color.Green(str)
// }

// func LogErr(msg ...interface{}) {
// 	str := "// " + sprint(msg...)
// 	log2GUI(str)
// 	log2File(str)
// 	color.New(color.FgRed).Fprintln(os.Stderr, str)
// }

// func log2File(msg string) {
// 	if logfile != nil {
// 		_, err := logfile.Write([]byte(msg + "\n"))
// 		if err != nil {
// 			LogErr("Error writing to log file:", err)
// 		}
// 	}
// }

// func initLog() {
// 	f, err := httpfs.Create(OD() + "log.txt")
// 	util.Log.PanicIfError(err)
// 	logfile = f // otherwise f gets dropped
// 	_, err = logfile.Write([]byte(Hist))
// 	if err != nil {
// 		LogErr("Error writing to log file:", err)
// 	}
// }

// func AutoFlushLog2File() {
// 	for {
// 		logfile.Flush()
// 		time.Sleep(5 * time.Second)
// 	}
// }

// func log2GUI(msg string) {
// 	if len(msg) > 1000 {
// 		msg = msg[:1000-len("...")] + "..."
// 	}
// 	if Hist != "" { // prepend newline
// 		Hist += "\n"
// 	}
// 	Hist += msg
// }

// // like fmt.Sprint but with spaces between args
// func sprint(msg ...interface{}) string {
// 	str := fmt.Sprintln(msg...)
// 	str = str[:len(str)-1] // strip newline
// 	return str
// }
