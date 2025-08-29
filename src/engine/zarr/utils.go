package zarr

import (
	"strings"

	"github.com/MathieuMoalic/amumax/src/engine/fsutil"
	"github.com/MathieuMoalic/amumax/src/engine/log"
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
	defer zgroup.Close()
	_, err = zgroup.Write([]byte("{\"zarr_format\": 2}"))
	log.Log.PanicIfError(err)
}
