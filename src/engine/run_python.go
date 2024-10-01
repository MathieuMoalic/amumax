package engine

import (
	"os/exec"

	"github.com/MathieuMoalic/amumax/src/util"
)

func init() {
	DeclFunc("RunShell", runShell, "Run a shell command")
	Insecure = false
}

var Insecure bool

// runPython executes a Python script with optional arguments after the simulation completes.
func runShell(cmd_str string) {
	if !Insecure {
		util.Log.Err("Insecure mode is disabled. To run shell commands, use the --insecure flag.")
		return
	}

	output, err := exec.Command(cmd_str).CombinedOutput()
	if err != nil {
		util.Log.Err("Error running shell commands: %v\nOutput: %s", err, output)
	}

	util.Log.Info("%s", output)
}
