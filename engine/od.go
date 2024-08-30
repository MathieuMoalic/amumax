package engine

// Management of output directory.

import (
	"os"
	"strings"

	"github.com/MathieuMoalic/amumax/httpfs"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/MathieuMoalic/amumax/zarr"
)

var (
	outputdir string // Output directory
	InputFile string

	CacheDir       string
	SkipExists     bool
	ForceClean     bool
	ShowProgresBar bool
	SelfTest       bool
	SyncAndLog     bool
)

func OD() string {
	if outputdir == "" {
		panic("output not yet initialized")
	}
	return outputdir
}

// SetOD sets the output directory where auto-saved files will be stored.
// The -o flag can also be used for this purpose.
func InitIO(mx3Path, od, cachedir string, skipexists, forceclean, showprogressbar, selftest, syncandlog bool) {
	CacheDir = cachedir
	SkipExists = skipexists
	ForceClean = forceclean
	ShowProgresBar = showprogressbar
	SelfTest = selftest
	SyncAndLog = syncandlog

	if outputdir != "" {
		panic("output directory already set")
	}
	InputFile = mx3Path
	if !strings.HasSuffix(od, "/") {
		od += "/"
	}
	outputdir = od
	if strings.HasPrefix(outputdir, "http://") {
		httpfs.SetWD(outputdir + "/../")
	}
	if httpfs.IsDir(od) {
		// if directory exists and --skip-exist flag is set, skip the directory
		if SkipExists {
			util.Log.Warn("Directory `%s` exists, skipping `%s` because of --skip-exist flag.", od, mx3Path)
			os.Exit(0)
			// if directory exists and --force-clean flag is set, remove the directory
		} else if ForceClean {
			util.Log.Warn("Cleaning `%s`", od)
			util.Log.PanicIfError(httpfs.Remove(od))
			util.Log.PanicIfError(httpfs.Mkdir(od))
		}
	} else {
		util.Log.PanicIfError(httpfs.Mkdir(od))
	}
	zarr.InitZgroup(OD())
}
