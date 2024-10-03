package zarr

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"
	"time"

	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/httpfs"
	"github.com/MathieuMoalic/amumax/src/log"

	"github.com/DataDog/zstd"
)

func Read(binary_path string, pwd string) (s *data.Slice, err error) {
	if !path.IsAbs(binary_path) {
		wd := ""
		wd, err = os.Getwd()
		log.Log.PanicIfError(err)
		binary_path = wd + "/" + path.Dir(pwd) + "/" + binary_path
	}
	binary_path = path.Clean(binary_path)

	// loop and wait until the file is saved
	for IsSaving {
		log.Log.Info("Waiting for all the files to be saved before reading...")
		time.Sleep(1 * time.Second)
	}

	zarray_path := path.Dir(binary_path) + "/.zarray"
	log.Log.Info("Reading:  %v", binary_path)
	io_reader, err := httpfs.Open(binary_path)
	log.Log.PanicIfError(err)
	content, err := os.ReadFile(zarray_path)
	log.Log.PanicIfError(err)
	var zarray zarrayFile
	err = json.Unmarshal([]byte(content), &zarray)
	log.Log.PanicIfError(err)
	if zarray.Compressor.ID != "zstd" {
		log.Log.PanicIfError(errors.New("LoadFile: Only the Zstd compressor is supported"))
	}
	sizez := zarray.Chunks[1]
	sizey := zarray.Chunks[2]
	sizex := zarray.Chunks[3]
	sizec := zarray.Chunks[4]

	array := data.NewSlice(sizec, [3]int{sizex, sizey, sizez})
	tensors := array.Tensors()
	ncomp := array.NComp()
	compressedData, err := io.ReadAll(io_reader)
	if err != nil {
		panic(err)
	}
	data, err := zstd.Decompress(nil, compressedData)
	if err != nil {
		panic(err)
	}
	count := 0
	for iz := 0; iz < sizez; iz++ {
		for iy := 0; iy < sizey; iy++ {
			for ix := 0; ix < sizex; ix++ {
				for c := 0; c < ncomp; c++ {
					tensors[c][iz][iy][ix] = BytesToFloat32(data[count*4 : (count+1)*4])
					count++
				}
			}
		}
	}

	return array, nil
}
