package zarr

import (
	"encoding/json"
	"io"
	"os"
	"path"

	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/httpfs"
	"github.com/MathieuMoalic/amumax/util"

	"github.com/DataDog/zstd"
)

func Read(fname string) (s *data.Slice, err error) {
	basedir := path.Dir(fname)
	content, err := os.ReadFile(basedir + "/.zarray")
	util.LogErr(err)
	var zarray zarrayFile
	json.Unmarshal([]byte(content), &zarray)
	if zarray.Compressor.ID != "zstd" {
		util.LogThenExit("Error: LoadFile: Only the Zstd compressor is supported")
	}
	sizez := zarray.Chunks[1]
	sizey := zarray.Chunks[2]
	sizex := zarray.Chunks[3]
	sizec := zarray.Chunks[4]

	util.Log("// chunks:", zarray.Chunks)
	util.Log("// size:", sizez, sizey, sizex, sizec)
	util.Log("// compressor:", zarray.Compressor.ID)
	array := data.NewSlice(sizec, [3]int{sizex, sizey, sizez})
	tensors := array.Tensors()
	ncomp := array.NComp()
	io_reader, err := httpfs.Open(fname)
	if err != nil {
		panic(err)
	}
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
