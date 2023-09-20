package zarr

import (
	"encoding/json"
	"fmt"

	"github.com/MathieuMoalic/amumax/httpfs"
	"github.com/MathieuMoalic/amumax/util"
)

type ztableFile struct {
	Chunks     [1]int  `json:"chunks"`
	Compressor []int   `json:"compressor"`
	Dtype      string  `json:"dtype"`
	FillValue  float64 `json:"fill_value"`
	Filters    []int   `json:"filters"`
	Order      string  `json:"order"`
	Shape      [1]int  `json:"shape"`
	ZarrFormat int     `json:"zarr_format"`
}

func SaveFileTableZarray(path string, zTableAutoSaveStep int) {
	if !util.PathExists(path) {
		util.FatalErr(fmt.Errorf("error: `%s` does not exist", path))
	}
	z := ztableFile{}
	z.Dtype = `<f8`
	z.FillValue = 0.0
	z.Order = "C"
	z.ZarrFormat = 2
	z.Chunks = [1]int{zTableAutoSaveStep + 1}
	z.Shape = [1]int{zTableAutoSaveStep + 1}

	f, err := httpfs.Create(path + "/.zarray")
	util.FatalErr(err)
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	enc.Encode(z)
	f.Flush()
}
