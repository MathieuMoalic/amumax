package engine

import (
	"os/exec"

	"github.com/MathieuMoalic/amumax/src/log"
)

func init() {
	Insecure = false
}

var Insecure bool

// runPython executes a Python script with optional arguments after the simulation completes.
func runShell(cmdStr string) {
	if !Insecure {
		log.Log.Err("Insecure mode is disabled. To run shell commands, use the --insecure flag.")
		return
	}

	output, err := exec.Command(cmdStr).CombinedOutput()
	if err != nil {
		log.Log.Err("Error running shell commands: %v\nOutput: %s", err, output)
	}

	log.Log.Info("%s", output)
}
