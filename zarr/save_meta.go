package zarr

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/httpfs"
	"github.com/MathieuMoalic/amumax/util"
)

type MetaStruct struct {
	root_path  string
	start_time time.Time
	Dx         float64 `json:"dx"`
	Dy         float64 `json:"dy"`
	Dz         float64 `json:"dz"`
	Nx         int     `json:"Nx"`
	Ny         int     `json:"Ny"`
	Nz         int     `json:"Nz"`
	Tx         float64 `json:"Tx"`
	Ty         float64 `json:"Ty"`
	Tz         float64 `json:"Tz"`
	StartTime  string  `json:"start_time"`
	EndTime    string  `json:"end_time"`
	TotalTime  string  `json:"total_time"`
	PBC        [3]int  `json:"PBC"`
	Gpu        string  `json:"gpu"`
	Host       string  `json:"host"`
}

func (meta *MetaStruct) Init(m data.Mesh, savepath string, gpu string) {
	meta.root_path = savepath
	meta.start_time = time.Now()
	D, N := m.CellSize(), m.Size()
	meta.Dx = D[0]
	meta.Dy = D[1]
	meta.Dz = D[2]
	meta.Nx = N[0]
	meta.Ny = N[1]
	meta.Nz = N[2]
	meta.Tx = float64(N[0]) * D[0]
	meta.Ty = float64(N[1]) * D[1]
	meta.Tz = float64(N[2]) * D[2]
	meta.PBC = m.PBC()
	meta.Gpu = gpu
	meta.StartTime = time.Now().Format(time.UnixDate)
	node := os.Getenv("SLURM_NODELIST")
	if node == "" {
		hostname, _ := os.Hostname()
		meta.Host = hostname
	} else {
		meta.Host = node
	}
	meta.Save()
}

func (meta *MetaStruct) End() {
	meta.EndTime = time.Now().Format(time.UnixDate)
	meta.TotalTime = fmt.Sprint(time.Since(meta.start_time))
	meta.Save()
}

func (meta *MetaStruct) Save() {
	zattrs, err := httpfs.Create(meta.root_path + ".zattrs")
	util.FatalErr(err)
	defer zattrs.Close()
	json_meta, err := json.MarshalIndent(meta, "", "\t")
	util.FatalErr(err)
	zattrs.Write([]byte(json_meta))

}
