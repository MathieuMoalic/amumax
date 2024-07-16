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
	if *Flag_skip_exists {
		err := httpfs.Mkdir(od)
		if err != nil {
			// cursed error check, not sure how to do better
			if fmt.Sprint(err) == fmt.Sprintf("mkdir %s: file exists", od) {
				LogErr(fmt.Sprintf("Directory `%s` exists, skipping `%s` because of --skip-exist flag.", od, inputfile))
				os.Exit(0)
			} else {
				util.FatalErr(err)
			}
		}
	}
	LogOut("output directory:", outputdir)

	if *Flag_forceclean && !*Flag_skip_exists {
		err := httpfs.Remove(od)
		if err != nil {
			util.FatalErr(err)
		}
	}

	err := httpfs.Mkdir(od)
	if err != nil {
		util.FatalErr(err)
	}
	// util.FatalErr(err)
	initLog()
	zarr.InitZgroup(OD())
}
