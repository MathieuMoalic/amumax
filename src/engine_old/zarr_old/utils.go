package zarr_old

import (
	"strings"

	"github.com/MathieuMoalic/amumax/src/engine_old/fsutil_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
)

type Zattrs struct {
	Buffer []float64 `json:"t"`
}

func InitZgroup(name string, od string) {
	err := fsutil_old.Mkdir(od + name)
	if err != nil && !strings.Contains(err.Error(), "file exists") {
		log_old.Log.PanicIfError(err)
	}
	path := ""
	if name == "" {
		path = od + ".zgroup"
	} else {
		path = od + name + "/.zgroup"
	}
	zgroup, err := fsutil_old.Create(path)
	log_old.Log.PanicIfError(err)
	defer zgroup.Close()
	_, err = zgroup.Write([]byte("{\"zarr_format\": 2}"))
	log_old.Log.PanicIfError(err)
}
