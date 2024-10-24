package engine

import (
	"encoding/json"
	"fmt"
	"unsafe"

	"github.com/DataDog/zstd"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/fsutil"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/zarr"
)

func mx3AutoSave(q Quantity, period float64) {
	zVerifyAndSave(q, nameOf(q), requestedChunking{1, 1, 1, 1}, period)
}

func mx3AutoSaveAs(q Quantity, name string, period float64) {
	zVerifyAndSave(q, name, requestedChunking{1, 1, 1, 1}, period)
}

func mx3AutoSaveAsChunk(q Quantity, name string, period float64, rchunks requestedChunking) {
	zVerifyAndSave(q, name, rchunks, period)
}

func mx3SaveAs(q Quantity, name string) {
	zVerifyAndSave(q, name, requestedChunking{1, 1, 1, 1}, 0)
}

func mx3zSave(q Quantity) {
	zVerifyAndSave(q, nameOf(q), requestedChunking{1, 1, 1, 1}, 0)
}

func mx3SaveAsChunk(q Quantity, name string, rchunks requestedChunking) {
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
	chunks  chunks
	rchunks requestedChunking
}

// returns true when the time is right to save.
func (a *zArray) needSave() bool {
	return a.period != 0 && (Time-a.start)-float64(a.count)*a.period >= a.period
}

func (a *zArray) SaveAttrs() {
	// it's stupid and wasteful but it works
	// keeping the whole array of times wastes a few MB of RAM
	u, err := json.Marshal(zarr.Zattrs{Buffer: a.times})
	log.Log.PanicIfError(err)
	err = fsutil.Remove(OD() + a.name + "/.zattrs")
	log.Log.PanicIfError(err)
	err = fsutil.Put(OD()+a.name+"/.zattrs", []byte(string(u)))
	log.Log.PanicIfError(err)
}

func (a *zArray) Save() {
	a.times = append(a.times, Time)
	a.SaveAttrs()
	buffer := ValueOf(a.q)
	defer cuda.Recycle(buffer)
	data := buffer.HostCopy() // must be copy (async io)
	t := a.count              // no desync this way
	queOutput(func() { syncSave(data, a.name, t, a.chunks) })
	a.count++
}

// entrypoint of all the user facing save functions
func zVerifyAndSave(q Quantity, name string, rchunks requestedChunking, period float64) {
	if zArrayExists(q, name, rchunks) {
		for _, z := range zArrays {
			if z.name == name {
				z.period = period
				z.Save()
			}
		}
	} else {
		if fsutil.Exists(OD() + name) {
			err := fsutil.Remove(OD() + name)
			log.Log.PanicIfError(err)
		}
		err := fsutil.Mkdir(OD() + name)
		log.Log.PanicIfError(err)
		var b []float64
		a := zArray{name, q, period, Time, -1, b, newChunks(q, rchunks), rchunks}
		zArrays = append(zArrays, &a)
		a.Save()
	}
}

func zArrayExists(q Quantity, name string, rchunks requestedChunking) bool {
	for _, z := range zArrays {
		if z.name == name {
			if z.rchunks != rchunks {
				log.Log.ErrAndExit("Error: The dataset %v has already been initialized with different chunks.", name)
			} else if z.q != q {
				log.Log.ErrAndExit("Error: The dataset %v has already been initialized with a different quantity.", name)
			} else {
				return true
			}
		}
	}
	return false
}

// synchronous chunky save
func syncSave(array *data.Slice, qname string, time int, chunks chunks) {
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
					f, err := fsutil.Create(fmt.Sprintf(OD()+"%s/%d.%d.%d.%d.%d", qname, time+1, icz, icy, icx, icc))
					log.Log.PanicIfError(err)
					defer f.Close()
					for iz := 0; iz < chunks.z.len; iz++ {
						z := icz*chunks.z.len + iz
						for iy := 0; iy < chunks.y.len; iy++ {
							y := icy*chunks.y.len + iy
							for ix := 0; ix < chunks.x.len; ix++ {
								x := icx*chunks.x.len + ix
								for ic := 0; ic < chunks.c.len; ic++ {
									// log.Log.Comment(ic,icc)
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
					log.Log.PanicIfError(err)
					_, err = f.Write(compressedData)
					log.Log.PanicIfError(err)
				}
			}
		}
	}
}
