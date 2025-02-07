package engine_old

import (
	"flag"
	"fmt"
	"os"

	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
)

// check all input files for errors, don't run.
func Vet() {
	status := 0
	for _, f := range flag.Args() {
		src, ioerr := os.ReadFile(f)
		log_old.Log.PanicIfError(ioerr)
		World.EnterScope() // avoid name collisions between separate files
		_, err := World.Compile(string(src))
		World.ExitScope()
		if err != nil {
			fmt.Println(f, ":", err)
			status = 1
		} else {
			fmt.Println(f, ":", "OK")
		}
	}
	os.Exit(status)
}
