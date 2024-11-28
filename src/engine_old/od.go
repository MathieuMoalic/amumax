package engine_old

// Management of output directory.

import (
	"os"
	"strings"

	"github.com/MathieuMoalic/amumax/src/fsutil_old"
	"github.com/MathieuMoalic/amumax/src/log_old"
	"github.com/MathieuMoalic/amumax/src/zarr"
)

var (
	outputdir string // Output directory
	InputFile string

	CacheDir       string
	SkipExists     bool
	ForceClean     bool
	HideProgresBar bool
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
func InitIO(mx3Path, od, cachedir string, skipexists, forceclean, hideprogressbar, selftest, syncandlog bool) {
	CacheDir = cachedir
	SkipExists = skipexists
	ForceClean = forceclean
	HideProgresBar = hideprogressbar
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
		fsutil_old.SetWD(outputdir + "/../")
	}
	if fsutil_old.IsDir(od) {
		// if directory exists and --skip-exist flag is set, skip the directory
		if SkipExists {
			log_old.Log.Warn("Directory `%s` exists, skipping `%s` because of --skip-exist flag.", od, mx3Path)
			os.Exit(0)
			// if directory exists and --force-clean flag is set, remove the directory
		} else if ForceClean {
			log_old.Log.Warn("Cleaning `%s`", od)
			log_old.Log.PanicIfError(fsutil_old.Remove(od))
			log_old.Log.PanicIfError(fsutil_old.Mkdir(od))
		}
	} else {
		log_old.Log.PanicIfError(fsutil_old.Mkdir(od))
	}
	zarr.InitZgroup("", OD())
}
