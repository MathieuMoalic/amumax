package zarr

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/httpfs"

	"github.com/DataDog/zstd"
)

type JsonHell struct {
	Chunks     [5]int
	Compressor string
	Dtype      string
	FillValue  string
	Filters    string
	Order      string
	Shape      string
	ZFormat    string
}

func Read(fname string) (s *data.Slice, err error) {
	basedir := path.Dir(fname)
	content, _ := os.ReadFile(basedir + "/.zarray")
	var zarray JsonHell
	json.Unmarshal([]byte(content), &zarray)

	sizez := zarray.Chunks[1]
	sizey := zarray.Chunks[2]
	sizex := zarray.Chunks[3]
	sizec := zarray.Chunks[4]

	fmt.Println(zarray.Chunks)
	fmt.Println(sizez, sizey, sizex, sizec)
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
