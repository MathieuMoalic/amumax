package engine_old

import (
	"os/exec"

	"github.com/MathieuMoalic/amumax/src/log_old"
)

func init() {
	Insecure = false
}

var Insecure bool

// runPython executes a Python script with optional arguments after the simulation completes.
func runShell(cmd_str string) {
	if !Insecure {
		log_old.Log.Err("Insecure mode is disabled. To run shell commands, use the --insecure flag.")
		return
	}

	output, err := exec.Command(cmd_str).CombinedOutput()
	if err != nil {
		log_old.Log.Err("Error running shell commands: %v\nOutput: %s", err, output)
	}

	log_old.Log.Info("%s", output)
}
