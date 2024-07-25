package engine

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
	DeclFunc("AutoSaveAs", Mx3AutoSaveAs, "Auto save space-dependent quantity every period (s) as the zarr standard.")
	DeclFunc("AutoSaveAsChunk", Mx3AutoSaveAsChunk, "Auto save space-dependent quantity every period (s) as the zarr standard.")
	DeclFunc("AutoSave", Mx3AutoSave, "Auto save space-dependent quantity every period (s) as the zarr standard.")
	DeclFunc("SaveAs", Mx3SaveAs, "Save space-dependent quantity as the zarr standard.")
	DeclFunc("SaveAsChunk", Mx3SaveAsChunk, "")
	DeclFunc("Save", Mx3zSave, "Save space-dependent quantity as the zarr standard.")
}

func Mx3AutoSave(q Quantity, period float64) {
	zVerifyAndSave(q, NameOf(q), RequestedChunking{1, 1, 1, 1}, period)
}

func Mx3AutoSaveAs(q Quantity, name string, period float64) {
	zVerifyAndSave(q, name, RequestedChunking{1, 1, 1, 1}, period)
}

func Mx3AutoSaveAsChunk(q Quantity, name string, period float64, rchunks RequestedChunking) {
	zVerifyAndSave(q, name, rchunks, period)
}

func Mx3SaveAs(q Quantity, name string) {
	zVerifyAndSave(q, name, RequestedChunking{1, 1, 1, 1}, 0)
}

func Mx3zSave(q Quantity) {
	zVerifyAndSave(q, NameOf(q), RequestedChunking{1, 1, 1, 1}, 0)
}

func Mx3SaveAsChunk(q Quantity, name string, rchunks RequestedChunking) {
	zVerifyAndSave(q, name, rchunks, 0)
}

var zGroups []string
var zArrays []*zArray

type zArray struct {
	name    string
	q       Quantity
	period  float64 // How often to save
	start   float64 // Starting point
	count   int     // Number of times it has been autosaved
	times   []float64
	chunks  Chunks
	rchunks RequestedChunking
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
	err = httpfs.Remove(OD() + a.name + "/.zattrs")
	util.FatalErr(err)
	err = httpfs.Put(OD()+a.name+"/.zattrs", []byte(string(u)))
	util.FatalErr(err)
}

func (a *zArray) Save() {
	a.times = append(a.times, Time)
	a.SaveAttrs()
	buffer := ValueOf(a.q)
	defer cuda.Recycle(buffer)
	data := buffer.HostCopy() // must be copy (async io)
	t := a.count              // no desync this way
	queOutput(func() { SyncSave(data, a.name, t, a.chunks) })
	a.count++
}

// entrypoint of all the user facing save functions
func zVerifyAndSave(q Quantity, name string, rchunks RequestedChunking, period float64) {
	if zArrayExists(q, name, rchunks) {
		for _, z := range zArrays {
			if z.name == name {
				z.period = period
				z.Save()
			}
		}
	} else {
		if httpfs.Exists(OD() + name) {
			err := httpfs.Remove(OD() + name)
			util.FatalErr(err)
		}
		err := httpfs.Mkdir(OD() + name)
		util.FatalErr(err)
		var b []float64
		a := zArray{name, q, period, Time, -1, b, NewChunks(q, rchunks), rchunks}
		zArrays = append(zArrays, &a)
		a.Save()
	}
}

func zArrayExists(q Quantity, name string, rchunks RequestedChunking) bool {
	for _, z := range zArrays {
		if z.name == name {
			if z.rchunks != rchunks {
				util.Fatal("Error: The dataset `", name, "` has already been initialized with different chunks.")
			} else if z.q != q {
				util.Fatal("Error: The dataset `", name, "` has already been initialized with a different quantity.")
			} else {
				return true
			}
		}
	}
	return false
}

// synchronous chunky save
func SyncSave(array *data.Slice, qname string, time int, chunks Chunks) {
	data := array.Tensors()
	size := array.Size()
	ncomp := array.NComp()
	// saving .zarray before the data might help resolve some unsync
	// errors when the simulation is running and the user loads data
	zarr.SaveFileZarray(fmt.Sprintf(OD()+"%s/.zarray", qname), size, ncomp, time+1, chunks.z.len, chunks.y.len, chunks.x.len, chunks.c.len)
	var bytes []byte
	for icx := 0; icx < chunks.x.nb; icx++ {
		for icy := 0; icy < chunks.y.nb; icy++ {
			for icz := 0; icz < chunks.z.nb; icz++ {
				for icc := 0; icc < chunks.c.nb; icc++ {
					bdata := []byte{}
					f, err := httpfs.Create(fmt.Sprintf(OD()+"%s/%d.%d.%d.%d.%d", qname, time+1, icz, icy, icx, icc))
					util.FatalErr(err)
					defer f.Close()
					for iz := 0; iz < chunks.z.len; iz++ {
						z := icz*chunks.z.len + iz
						for iy := 0; iy < chunks.y.len; iy++ {
							y := icy*chunks.y.len + iy
							for ix := 0; ix < chunks.x.len; ix++ {
								x := icx*chunks.x.len + ix
								for ic := 0; ic < chunks.c.len; ic++ {
									// LogOut(ic,icc)
									c := icc*chunks.c.len + ic
									bytes = (*[4]byte)(unsafe.Pointer(&data[c][z][y][x]))[:]
									for k := 0; k < 4; k++ {
										bdata = append(bdata, bytes[k])
									}
								}
							}
						}
					}
					compressedData, err := zstd.Compress(nil, bdata)
					util.FatalErr(err)
					_, err = f.Write(compressedData)
					util.FatalErr(err)
				}
			}
		}
	}
}
