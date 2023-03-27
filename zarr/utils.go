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
		util.FatalErr(err)
		InitZgroup(od + name + "/")
		*zGroups = append(*zGroups, name)
	}
}

func InitZgroup(path string) {
	zgroup, err := httpfs.Create(path + ".zgroup")
	util.FatalErr(err)
	defer zgroup.Close()
	zgroup.Write([]byte("{\"zarr_format\": 2}"))
}
