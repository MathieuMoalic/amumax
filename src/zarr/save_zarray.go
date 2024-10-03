package zarr

import (
	"encoding/json"

	"github.com/MathieuMoalic/amumax/src/httpfs"
	"github.com/MathieuMoalic/amumax/src/log"
)

var IsSaving bool

func init() {
	IsSaving = false
}

type ZstdCompressor struct {
	ID    string `json:"id"`
	Level int    `json:"level"`
}
type zarrayFile struct {
	Chunks     [5]int         `json:"chunks"`
	Compressor ZstdCompressor `json:"compressor"`
	Dtype      string         `json:"dtype"`
	FillValue  float64        `json:"fill_value"`
	Filters    []int          `json:"filters"`
	Order      string         `json:"order"`
	Shape      [5]int         `json:"shape"`
	ZarrFormat int            `json:"zarr_format"`
}

func SaveFileZarray(path string, size [3]int, ncomp int, time int, cz int, cy int, cx int, cc int) {
	IsSaving = true
	defer func() { IsSaving = false }()
	z := zarrayFile{}
	z.Compressor = ZstdCompressor{"zstd", 1}
	z.Dtype = `<f4`
	z.FillValue = 0.0
	z.Order = "C"
	z.ZarrFormat = 2
	z.Chunks = [5]int{1, cz, cy, cx, cc}
	z.Shape = [5]int{time + 1, size[2], size[1], size[0], ncomp}

	f, err := httpfs.Create(path)
	log.Log.PanicIfError(err)
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	log.Log.PanicIfError(enc.Encode(z))
	f.Flush()
}
