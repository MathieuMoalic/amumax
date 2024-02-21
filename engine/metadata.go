package engine

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/httpfs"
	"github.com/MathieuMoalic/amumax/util"
)

var Metadata map[string]interface{}

func init() {
	DeclFunc("Metadata", AddMetadata, "")
	Metadata = make(map[string]interface{})
}

func InitMetadata() {
	Metadata["start_time"] = StartTime
	Metadata["dx"] = dx
	Metadata["dy"] = dy
	Metadata["dz"] = dz
	Metadata["Nx"] = Nx
	Metadata["Ny"] = Ny
	Metadata["Nz"] = Nz
	Metadata["Tx"] = Tx
	Metadata["Ty"] = Ty
	Metadata["Tz"] = Tz
	Metadata["PBCx"] = PBCx
	Metadata["PBCy"] = PBCy
	Metadata["PBCz"] = PBCz
	Metadata["gpu"] = cuda.GPUInfo
	SaveMetadata()
}

func AddMetadata(key string, val interface{}) {
	Metadata[key] = val
	SaveMetadata()

}
func EndMetadata() {
	Metadata["end_time"] = time.Now().Format(time.UnixDate)
	Metadata["total_time"] = fmt.Sprint(time.Since(StartTime))
	SaveMetadata()
}
func SaveMetadata() {
	zattrs, err := httpfs.Create(OD() + ".zattrs")
	util.FatalErr(err)
	defer zattrs.Close()
	json_meta, err := json.MarshalIndent(Metadata, "", "\t")
	util.FatalErr(err)
	zattrs.Write([]byte(json_meta))
}
