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

func IntToByteArray(num int32) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}

// synchronous chunky save
func zSyncSave(array *data.Slice, qname string, time int) {
	data := array.Tensors()
	size := array.Size()
	ncomp := array.NComp()

	var icc_max int
	var ic_max int

	if ncomp == 1 {
		icc_max = 1
		ic_max = 1
	} else {
		if chunks.c.len == 1 {
			icc_max = 3
			ic_max = 1
		} else {
			icc_max = 1
			ic_max = 3

		}

	}
	count := 0
	for icx := 0; icx < chunks.x.nb; icx++ {
		for icy := 0; icy < chunks.y.nb; icy++ {
			for icz := 0; icz < chunks.z.nb; icz++ {
				for icc := 0; icc < icc_max; icc++ {
					bdata := []byte{}
					var bytes []byte
					f, err := httpfs.Create(fmt.Sprintf(OD()+"%s/%d.%d.%d.%d.%d", qname, time+1, icz, icy, icx, icc))
					util.FatalErr(err)
					defer f.Close()
					for iz := 0; iz < chunks.z.len; iz++ {
						z := icz*chunks.z.len + iz
						for iy := 0; iy < chunks.y.len; iy++ {
							y := icy*chunks.y.len + iy
							for ix := 0; ix < chunks.x.len; ix++ {
								x := icx*chunks.x.len + ix
								for ic := 0; ic < ic_max; ic++ {
									c := icc*chunks.c.len + ic
									bytes = (*[4]byte)(unsafe.Pointer(&data[c][z][y][x]))[:]
									for k := 0; k < 4; k++ {
										bdata = append(bdata, bytes[k])
									}
								}
								count++
							}
						}
					}
					compressedData, err := zstd.Compress(nil, bdata)
					util.FatalErr(err)
					f.Write(compressedData)
				}
			}
		}
	}
	//.zarray file
	zarr.SaveFileZarray(fmt.Sprintf(OD()+"%s/.zarray", qname), size, ncomp, time+1, chunks.z.len, chunks.y.len, chunks.x.len, ncomp)
}
