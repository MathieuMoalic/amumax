package fsutil

import (
	"encoding/json"
	"errors"
	"strings"
)

func (fs *FileSystem) CreateZarrGroup(name string) error {
	err := fs.Mkdir(name)
	if err != nil && !strings.Contains(err.Error(), "file exists") {
		return err
	}
	writer, file, err := fs.Create(name + ".zgroup")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = writer.WriteString("{\"zarr_format\": 2}")
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

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

func (fs *FileSystem) SaveFileTableZarray(path string, zTableAutoSaveStep int) error {
	if !fs.Exists(path) {
		return errors.New("error: `%s` does not exist")
	}
	z := ztableFile{}
	z.Dtype = `<f8`
	z.FillValue = 0.0
	z.Order = "C"
	z.ZarrFormat = 2
	z.Chunks = [1]int{zTableAutoSaveStep + 1}
	z.Shape = [1]int{zTableAutoSaveStep + 1}

	writer, file, err := fs.Create(path + "/.zarray")
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(writer)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	err = enc.Encode(z)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

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

func (fs *FileSystem) SaveFileZarray(path string, size [3]int, ncomp int, step int, cz int, cy int, cx int, cc int) error {
	IsSaving = true
	defer func() { IsSaving = false }()
	z := zarrayFile{}
	z.Compressor = ZstdCompressor{"zstd", 1}
	z.Dtype = `<f4`
	z.FillValue = 0.0
	z.Order = "C"
	z.ZarrFormat = 2
	z.Chunks = [5]int{1, cz, cy, cx, cc}
	z.Shape = [5]int{step + 1, size[2], size[1], size[0], ncomp}

	writer, file, err := fs.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	err = enc.Encode(z)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

func (fs *FileSystem) SaveZattrs(path string, data map[string]interface{}) error {
	writer, file, err := fs.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(writer)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	err = enc.Encode(data)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}
