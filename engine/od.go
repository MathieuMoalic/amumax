package engine

// Management of output directory.

import (
	"fmt"
	"os"
	"strings"

	"github.com/MathieuMoalic/amumax/httpfs"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/MathieuMoalic/amumax/zarr"
)

var (
	outputdir string // Output directory
	InputFile string
)

func OD() string {
	if outputdir == "" {
		panic("output not yet initialized")
	}
	return outputdir
}

// SetOD sets the output directory where auto-saved files will be stored.
// The -o flag can also be used for this purpose.
func InitIO(inputfile, od string) {
	if outputdir != "" {
		panic("output directory already set")
	}
	InputFile = inputfile
	if !strings.HasSuffix(od, "/") {
		od += "/"
	}
	outputdir = od
	if strings.HasPrefix(outputdir, "http://") {
		httpfs.SetWD(outputdir + "/../")
	}
	if httpfs.Exists(od) {
		// if directory exists and --skip-exist flag is set, skip the directory
		if *Flag_skip_exists {
			util.LogWarn(fmt.Sprintf("Directory `%s` exists, skipping `%s` because of --skip-exist flag.", od, inputfile))
			os.Exit(0)
			// if directory exists and --force-clean flag is set, remove the directory
		} else if *Flag_forceclean {
			util.LogWarn(fmt.Sprintf("Cleaning `%s`", od))
			util.FatalErr(httpfs.Remove(od))
			util.FatalErr(httpfs.Mkdir(od))
		}
	} else {
		util.FatalErr(httpfs.Mkdir(od))
	}
	LogOut("Output directory:", outputdir)
	initLog()
	zarr.InitZgroup(OD())
}
