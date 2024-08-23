package zarr

import (
	"github.com/MathieuMoalic/amumax/httpfs"
	"github.com/MathieuMoalic/amumax/util"
)

type Zattrs struct {
	Buffer []float64 `json:"t"`
}

func MakeZgroup(name string, od string, zGroups *[]string) {
	exists := false
	for _, v := range *zGroups {
		if name == v {
			exists = true
			*zGroups = append(*zGroups, name)
		}
	}
	if !exists {
		err := httpfs.Mkdir(od + name)
		util.Log.PanicIfError(err)
		InitZgroup(od + name + "/")
		*zGroups = append(*zGroups, name)
	}
}

func InitZgroup(path string) {
	zgroup, err := httpfs.Create(path + ".zgroup")
	util.Log.PanicIfError(err)
	defer zgroup.Close()
	_, err = zgroup.Write([]byte("{\"zarr_format\": 2}"))
	util.Log.PanicIfError(err)
}
