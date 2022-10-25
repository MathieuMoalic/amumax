package zarr

import (
	"encoding/json"

	"github.com/MathieuMoalic/amumax/httpfs"
	"github.com/MathieuMoalic/amumax/util"
)

type compressorStruc struct {
	ID    string `json:"id"`
	Level int    `json:"level"`
}
type zarrayFile struct {
	Chunks     [5]int          `json:"chunks"`
	Compressor compressorStruc `json:"compressor"`
	Dtype      string          `json:"dtype"`
	FillValue  float64         `json:"fill_value"`
	Filters    []int           `json:"filters"`
	Order      string          `json:"order"`
	Shape      [5]int          `json:"shape"`
	ZarrFormat int             `json:"zarr_format"`
}

func SaveFileZarray(path string, size [3]int, ncomp int, time int, chunky bool) {
	z := zarrayFile{}
	z.Compressor = compressorStruc{"zstd", 1}
	z.Dtype = `<f4`
	z.FillValue = 0.0
	z.Order = "C"
	z.ZarrFormat = 2
	if chunky {
		z.Chunks = [5]int{1, size[2], 1, size[0], ncomp}
	} else {
		z.Chunks = [5]int{1, size[2], size[1], size[0], ncomp}
	}
	z.Shape = [5]int{time + 1, size[2], size[1], size[0], ncomp}

	f, err := httpfs.Create(path)
	util.FatalErr(err)
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	enc.Encode(z)
	f.Flush()
}
