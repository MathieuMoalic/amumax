package zarr

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/MathieuMoalic/amumax/src/fsutil"
	"github.com/MathieuMoalic/amumax/src/log"
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
		log.Log.PanicIfError(errors.New("error: `%s` does not exist"))
	}
	z := ztableFile{}
	z.Dtype = `<f8`
	z.FillValue = 0.0
	z.Order = "C"
	z.ZarrFormat = 2
	z.Chunks = [1]int{zTableAutoSaveStep + 1}
	z.Shape = [1]int{zTableAutoSaveStep + 1}

	f, err := fsutil.Create(path + "/.zarray")
	log.Log.PanicIfError(err)

	defer func() {
		cerr := f.Close()
		if cerr != nil {
			log.Log.Err("Error closing zarray file: %v", cerr)
		}
	}()
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	log.Log.PanicIfError(enc.Encode(z))
	err = f.Flush()
	log.Log.PanicIfError(err)
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
