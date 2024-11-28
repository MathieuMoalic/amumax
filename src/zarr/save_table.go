package zarr

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/MathieuMoalic/amumax/src/fsutil_old"
	"github.com/MathieuMoalic/amumax/src/log_old"
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
	if !pathExists(path) {
		log_old.Log.PanicIfError(errors.New("error: `%s` does not exist"))
	}
	z := ztableFile{}
	z.Dtype = `<f8`
	z.FillValue = 0.0
	z.Order = "C"
	z.ZarrFormat = 2
	z.Chunks = [1]int{zTableAutoSaveStep + 1}
	z.Shape = [1]int{zTableAutoSaveStep + 1}

	f, err := fsutil_old.Create(path + "/.zarray")
	log_old.Log.PanicIfError(err)

	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	log_old.Log.PanicIfError(enc.Encode(z))
	f.Flush()
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
