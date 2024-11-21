package engine

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/DataDog/zstd"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/fsutil"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/zarr"
)

// Global slice to keep track of all zArrays
var zArrays []*zArray
var zGroups []string

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

// needSave returns true when it's time to save based on the period.
func (a *zArray) needSave() bool {
	t := Time - a.start
	return a.period != 0 && t-float64(a.count)*a.period >= a.period
}

// SaveAttrs updates the .zattrs file with the times data.
func (a *zArray) SaveAttrs() {
	u, err := json.Marshal(zarr.Zattrs{Buffer: a.times})
	log.Log.PanicIfError(err)
	err = fsutil.Remove(OD() + a.name + "/.zattrs")
	log.Log.PanicIfError(err)
	err = fsutil.Put(OD()+a.name+"/.zattrs", u)
	log.Log.PanicIfError(err)
}

// Save writes the data to disk and updates the times.
func (a *zArray) Save() {
	a.times = append(a.times, Time)
	a.SaveAttrs()
	buffer := ValueOf(a.q)
	defer cuda.Recycle(buffer)
	dataSlice := buffer.HostCopy() // Must be a copy (async IO)
	t := a.count                   // Prevent desync
	queOutput(func() {
		err := syncSave(dataSlice, a.name, t, a.chunks)
		log.Log.PanicIfError(err)
	})
	a.count++
}

// saveZarrArrays is called periodically to save arrays when needed.
func saveZarrArrays() {
	for _, z := range zArrays {
		if z.needSave() {
			z.Save()
		}
	}
}

// getOrCreateZArray retrieves an existing zArray or creates a new one.
func getOrCreateZArray(q Quantity, name string, rchunks requestedChunking, period float64) *zArray {
	for _, z := range zArrays {
		if z.name == name {
			if z.rchunks != rchunks {
				log.Log.ErrAndExit("Error: The dataset %v has already been initialized with different chunks.", name)
			}
			if z.q != q {
				log.Log.ErrAndExit("Error: The dataset %v has already been initialized with a different quantity.", name)
			}
			// Update period if a new non-zero period is provided
			if period > 0 {
				z.period = period
			}
			return z
		}
	}
	// Create a new zArray
	if fsutil.Exists(OD() + name) {
		err := fsutil.Remove(OD() + name)
		log.Log.PanicIfError(err)
	}
	err := fsutil.Mkdir(OD() + name)
	log.Log.PanicIfError(err)
	newZArray := &zArray{
		name:    name,
		q:       q,
		period:  period,
		start:   Time,
		count:   -1,
		times:   []float64{},
		chunks:  newChunks(q, rchunks),
		rchunks: rchunks,
	}
	zArrays = append(zArrays, newZArray)
	return newZArray
}

// zVerifyAndSave is the unified function for saving quantities.
func zVerifyAndSave(q Quantity, name string, rchunks requestedChunking, period float64) {
	z := getOrCreateZArray(q, name, rchunks, period)
	z.Save()
}

// User-facing save functions (function signatures cannot change)
func autoSave(q Quantity, period float64) {
	zVerifyAndSave(q, nameOf(q), requestedChunking{1, 1, 1, 1}, period)
}

func autoSaveAs(q Quantity, name string, period float64) {
	zVerifyAndSave(q, name, requestedChunking{1, 1, 1, 1}, period)
}

func autoSaveAsChunk(q Quantity, name string, period float64, rchunks requestedChunking) {
	zVerifyAndSave(q, name, rchunks, period)
}

func saveAs(q Quantity, name string) {
	zVerifyAndSave(q, name, requestedChunking{1, 1, 1, 1}, 0)
}

func save(q Quantity) {
	zVerifyAndSave(q, nameOf(q), requestedChunking{1, 1, 1, 1}, 0)
}

func saveAsChunk(q Quantity, name string, rchunks requestedChunking) {
	zVerifyAndSave(q, name, rchunks, 0)
}

// syncSave writes the data slice into chunked, compressed files compatible with the Zarr format.
func syncSave(array *data.Slice, qname string, time int, chunks chunks) error {
	data := array.Tensors()
	size := array.Size()
	ncomp := array.NComp()

	// Save .zarray metadata
	zarr.SaveFileZarray(
		fmt.Sprintf(OD()+"%s/.zarray", qname),
		size,
		ncomp,
		time+1,
		chunks.z.len, chunks.y.len, chunks.x.len, chunks.c.len,
	)

	// Iterate over chunks and save data
	for icx := 0; icx < chunks.x.nb; icx++ {
		for icy := 0; icy < chunks.y.nb; icy++ {
			for icz := 0; icz < chunks.z.nb; icz++ {
				for icc := 0; icc < chunks.c.nb; icc++ {
					var bdata bytes.Buffer
					for iz := 0; iz < chunks.z.len; iz++ {
						z := icz*chunks.z.len + iz
						for iy := 0; iy < chunks.y.len; iy++ {
							y := icy*chunks.y.len + iy
							for ix := 0; ix < chunks.x.len; ix++ {
								x := icx*chunks.x.len + ix
								for ic := 0; ic < chunks.c.len; ic++ {
									c := icc*chunks.c.len + ic
									value := data[c][z][y][x]
									err := binary.Write(&bdata, binary.LittleEndian, value)
									if err != nil {
										return err
									}
								}
							}
						}
					}
					compressedData, err := zstd.Compress(nil, bdata.Bytes())
					if err != nil {
						return err
					}
					filename := fmt.Sprintf(OD()+"%s/%d.%d.%d.%d.%d", qname, time+1, icz, icy, icx, icc)
					err = fsutil.Put(filename, compressedData)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
