package engine

// Bookkeeping for auto-saving quantities at given intervals.

import (
	"encoding/json"
	"fmt"
	"unsafe"

	"github.com/DataDog/zstd"

	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/httpfs"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/MathieuMoalic/amumax/zarr"
)

func init() {
	DeclFunc("AutoSaveAs", zAutoSaveAs, "Auto save space-dependent quantity every period (s) as the zarr standard.")
	DeclFunc("AutoSave", zAutoSave, "Auto save space-dependent quantity every period (s) as the zarr standard.")
	DeclFunc("SaveAs", zSaveAs, "Save space-dependent quantity as the zarr standard.")
	DeclFunc("Save", zSave, "Save space-dependent quantity as the zarr standard.")
	// DeclFunc("Savec", zSavec, "Save space-dependent quantity as the zarr standard.")
}

var zGroups []string
var zArrays = make(map[string]*zArray)

type zArray struct {
	name   string
	q      Quantity
	period float64                // How often to save
	start  float64                // Starting point
	count  int                    // Number of times it has been autosaved
	save   func(Quantity, string) // called to do the actual save
	times  []float64
}

// returns true when the time is right to save.
func (a *zArray) needSave() bool {
	return a.period != 0 && (Time-a.start)-float64(a.count)*a.period >= a.period
}
func (a *zArray) SaveAttrs() {
	// it's stupid and wasteful but it works
	// keeping the whole array of times wastes a few MB of RAM
	u, err := json.Marshal(zarr.Zattrs{Buffer: a.times})
	util.FatalErr(err)
	httpfs.Remove(OD() + a.name + "/.zattrs")
	httpfs.Put(OD()+a.name+"/.zattrs", []byte(string(u)))

}

func zAutoSave(q Quantity, period float64) {
	zAutoSaveAs(q, NameOf(q), period)
}

func zInitArray(name string, q Quantity, period float64) {
	httpfs.Mkdir(OD() + name)
	var b []float64
	zArrays[name] = &zArray{name, q, period, Time, -1, zSaveAs, b} // init count to -1 allows save at t=0

}

// period == 0 stops autosaving.
func zAutoSaveAs(q Quantity, name string, period float64) {
	if period == 0 {
		delete(zArrays, name)
	} else {
		zInitArray(name, q, period)
	}
}

func zSaveAs(q Quantity, name string) {
	if _, ok := zArrays[name]; !ok {
		zInitArray(name, q, 0)
	}
	zArrays[name].times = append(zArrays[name].times, Time)
	zArrays[name].SaveAttrs()
	buffer := ValueOf(q)
	defer cuda.Recycle(buffer)
	data := buffer.HostCopy() // must be copy (async io)
	t := zArrays[name].count  // no desync this way
	queOutput(func() { zSyncSave(data, name, t) })
	zArrays[name].count++

}

// Save once, with auto file name
func zSave(q Quantity) {
	zSaveAs(q, NameOf(q))
}

// synchronous save
func zSyncSave(array *data.Slice, qname string, time int) {
	f, err := httpfs.Create(fmt.Sprintf(OD()+"%s/%d.0.0.0.0", qname, time+1))
	util.FatalErr(err)
	defer f.Close()

	data := array.Tensors()
	size := array.Size()

	var bdata []byte
	var bytes []byte

	ncomp := array.NComp()
	for iz := 0; iz < size[Z]; iz++ {
		for iy := 0; iy < size[Y]; iy++ {
			for ix := 0; ix < size[X]; ix++ {
				for c := 0; c < ncomp; c++ {
					bytes = (*[4]byte)(unsafe.Pointer(&data[c][iz][iy][ix]))[:]
					for k := 0; k < 4; k++ {
						bdata = append(bdata, bytes[k])
					}
				}
			}
		}
	}
	// CompressLevel(dst, src []byte, level int) // alternative with compress levels
	compressedData, err := zstd.Compress(nil, bdata)
	util.FatalErr(err)
	f.Write(compressedData)

	//.zarray file
	zarr.SaveFileZarray(fmt.Sprintf(OD()+"%s/.zarray", qname), size, ncomp, time+1)
}

// func zSavec(q Quantity) {
// 	name := NameOf(q)
// 	if _, ok := zArrays[name]; !ok {
// 		zInitArray(name, q, 0)
// 	}
// 	zArrays[name].times = append(zArrays[name].times, Time)
// 	zArrays[name].SaveAttrs()
// 	buffer := ValueOf(q)
// 	defer cuda.Recycle(buffer)
// 	data := buffer.HostCopy() // must be copy (async io)
// 	t := zArrays[name].count  // no desync this way
// 	queOutput(func() { zSyncSavec(data, name, t) })
// 	zArrays[name].count++

// }

// func zSyncSavec(array *data.Slice, qname string, time int) {
// 	// f, err := httpfs.Create(fmt.Sprintf(OD()+"%s/%d.0.0.0.0", qname, time+1))
// 	// util.FatalErr(err)
// 	// defer f.Close()
// 	data := array.Tensors()
// 	size := array.Size()
// 	ncomp := array.NComp()

// 	tchunk := 5
// 	// zchunk := size[Z]
// 	ychunk := size[Y] / 32
// 	xchunk := size[X] / 32

// 	for iy := 0; iy < ychunk; iy++ {
// 		for ix := 0; ix < xchunk; ix++ {
// 			if (time % tchunk) == 0 {
// 				f, err := httpfs.Create(fmt.Sprintf(OD()+"%s/%d.0.0.0.0", qname, time+1))
// 				util.FatalErr(err)

// 			}
// 		}
// 	}

// 	var bdata []byte
// 	var bytes []byte

// 	for iz := 0; iz < size[Z]; iz++ {
// 		for iy := 0; iy < size[Y]; iy++ {
// 			for ix := 0; ix < size[X]; ix++ {
// 				for c := 0; c < ncomp; c++ {
// 					bytes = (*[4]byte)(unsafe.Pointer(&data[c][iz][iy][ix]))[:]
// 					for k := 0; k < 4; k++ {
// 						bdata = append(bdata, bytes[k])
// 					}
// 				}
// 			}
// 		}
// 	}
// 	// CompressLevel(dst, src []byte, level int) // alternative with compress levels
// 	// compressedData, err := zstd.Compress(nil, bdata)
// 	// util.FatalErr(err)
// 	// f.Write(compressedData)

// 	//.zarray file
// 	// zarr.SaveFileZarray(fmt.Sprintf(OD()+"%s/.zarray", qname), size, ncomp, time+1)
// }
