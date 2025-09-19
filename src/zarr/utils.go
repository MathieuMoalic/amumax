// Package zarr provides utilities to create zarr files.
package zarr

import (
	"strings"

	"github.com/MathieuMoalic/amumax/src/fsutil"
	"github.com/MathieuMoalic/amumax/src/log"
)

type Zattrs struct {
	Buffer []float64 `json:"t"`
}

func InitZgroup(name string, od string) {
	err := fsutil.Mkdir(od + name)
	if err != nil && !strings.Contains(err.Error(), "file exists") {
		log.Log.PanicIfError(err)
	}
	path := ""
	if name == "" {
		path = od + ".zgroup"
	} else {
		path = od + name + "/.zgroup"
	}
	zgroup, err := fsutil.Create(path)
	log.Log.PanicIfError(err)
	defer func() {
		cerr := zgroup.Close()
		if cerr != nil {
			log.Log.Err("Error closing zgroup file: %v", cerr)
		}
	}()
	_, err = zgroup.Write([]byte("{\"zarr_format\": 2}"))
	log.Log.PanicIfError(err)
}
