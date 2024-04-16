package zarr

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/MathieuMoalic/amumax/httpfs"
	"github.com/MathieuMoalic/amumax/util"
)

type Metadata struct {
	Fields    map[string]interface{}
	Path      string
	startTime time.Time
}

func (m *Metadata) Init(currentDir string, StartTime time.Time, dx, dy, dz float64, Nx, Ny, Nz int, Tx, Ty, Tz float64, PBCx, PBCy, PBCz int, GPUInfo string) {
	m.Fields["start_time"] = StartTime
	m.Fields["dx"] = dx
	m.Fields["dy"] = dy
	m.Fields["dz"] = dz
	m.Fields["Nx"] = Nx
	m.Fields["Ny"] = Ny
	m.Fields["Nz"] = Nz
	m.Fields["Tx"] = Tx
	m.Fields["Ty"] = Ty
	m.Fields["Tz"] = Tz
	m.Fields["PBCx"] = PBCx
	m.Fields["PBCy"] = PBCy
	m.Fields["PBCz"] = PBCz
	m.Fields["gpu"] = GPUInfo
	m.Path = currentDir + "/.zattrs"
	m.startTime = StartTime
	m.Save()
}

func (m *Metadata) Add(key string, val interface{}) {
	m.Fields[key] = val
	m.Save()
}
func (m *Metadata) End() {
	m.Fields["end_time"] = time.Now().Format(time.UnixDate)
	m.Fields["total_time"] = fmt.Sprint(time.Since(m.startTime))
	m.Save()
}

func (m *Metadata) Save() {
	util.Log(m.Path)
	zattrs, err := httpfs.Create(m.Path)
	util.FatalErr(err)
	defer zattrs.Close()
	json_meta, err := json.MarshalIndent(m.Fields, "", "\t")
	util.FatalErr(err)
	zattrs.Write([]byte(json_meta))
}
