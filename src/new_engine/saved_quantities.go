package new_engine

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

type savedQuantity struct {
	engineState *EngineStateStruct
	name        string
	q           Quantity
	period      float64
	times       []float64
	chunks      chunks
	rchunks     requestedChunking
	nextTime    float64 // Next time when autosave should trigger
}

// needSave returns true when it's time to save based on the period.
func (sq *savedQuantity) needSave() bool {
	if sq.period == 0 {
		return false
	}
	return sq.engineState.Solver.Time >= sq.nextTime
}

// SaveAttrs updates the .zattrs file with the times data.
func (sq *savedQuantity) SaveAttrs() {
	u, err := json.Marshal(zarr.Zattrs{Buffer: sq.times})
	log.Log.PanicIfError(err)
	err = fsutil.Remove(sq.engineState.ZarrPath + sq.name + "/.zattrs")
	log.Log.PanicIfError(err)
	err = fsutil.Put(sq.engineState.ZarrPath+sq.name+"/.zattrs", u)
	log.Log.PanicIfError(err)
}

func ValueOf(q Quantity) *data.Slice {
	// TODO: check for Buffered() implementation
	buf := cuda.Buffer(q.NComp(), q.Size())
	q.EvalTo(buf)
	return buf
}

// Save writes the data to disk and updates the times.
func (sq *savedQuantity) Save() {
	sq.times = append(sq.times, sq.engineState.Solver.Time)
	sq.SaveAttrs()
	buffer := cuda.Buffer(sq.q.NComp(), sq.q.Size())
	sq.q.EvalTo(buffer)
	defer cuda.Recycle(buffer)
	dataSlice := buffer.HostCopy()
	queOutput(func() {
		err := sq.syncSave(dataSlice, sq.name, len(sq.times), sq.chunks)
		log.Log.PanicIfError(err)
	})
}

// syncSave writes the data slice into chunked, compressed files compatible with the Zarr format.
func (sq *savedQuantity) syncSave(array *data.Slice, qname string, steps int, chunks chunks) error {
	data := array.Tensors()
	size := array.Size()
	ncomp := array.NComp()

	// Save .zarray metadata
	zarr.SaveFileZarray(
		fmt.Sprintf(sq.engineState.ZarrPath+"%s/.zarray", qname),
		size,
		ncomp,
		steps+1,
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
					filename := fmt.Sprintf(sq.engineState.ZarrPath+"%s/%d.%d.%d.%d.%d", qname, steps+1, icz, icy, icx, icc)
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

type savedQuantitiesType struct {
	EngineState *EngineStateStruct
	Quantities  []savedQuantity
}

func NewSavedQuantities(engineState *EngineStateStruct) *savedQuantitiesType {
	return &savedQuantitiesType{EngineState: engineState}
}

// saveZarrArrays is called periodically to save arrays when needed.
func (sqs *savedQuantitiesType) SaveIfNeeded() {
	for i := range sqs.Quantities {
		if sqs.Quantities[i].needSave() {
			sqs.Quantities[i].Save()
			sqs.Quantities[i].nextTime = sqs.Quantities[i].nextTime + sqs.Quantities[i].period
		}
	}
}

func (sqs *savedQuantitiesType) savedQuandtityExists(name string) bool {
	for _, z := range sqs.Quantities {
		if z.name == name {
			return true
		}
	}
	return false
}

func (sqs *savedQuantitiesType) createSavedQuantity(q Quantity, name string, rchunks requestedChunking, period float64) *savedQuantity {
	if fsutil.Exists(sqs.EngineState.ZarrPath + name) {
		err := fsutil.Remove(sqs.EngineState.ZarrPath + name)
		log.Log.PanicIfError(err)
	}
	err := fsutil.Mkdir(sqs.EngineState.ZarrPath + name)
	log.Log.PanicIfError(err)
	newZArray := &savedQuantity{
		name:    name,
		q:       q,
		period:  period,
		times:   []float64{},
		chunks:  newChunks(q, rchunks),
		rchunks: rchunks,
	}
	sqs.Quantities = append(sqs.Quantities, *newZArray)
	return newZArray

}

func (sqs *savedQuantitiesType) updateSavedQuantity(q Quantity, name string, rchunks requestedChunking, period float64) {
	sq := sqs.getSavedQuantity(name)
	if sq.rchunks != rchunks {
		log.Log.ErrAndExit("Error: The dataset %v has already been initialized with different chunks.", name)
	} else if sq.q != q {
		log.Log.ErrAndExit("Error: The dataset %v has already been initialized with a different quantity.", name)
	} else if sq.period != period {
		if sq.period == 0 && period != 0 {
			// enable autosave
			sq.period = period
			sq.nextTime = sqs.EngineState.Solver.Time + period
		} else if sq.period != 0 && period == 0 {
			// disable autosave
			sq.period = period
		}
	}
}

func (sqs *savedQuantitiesType) getSavedQuantity(name string) *savedQuantity {
	for i := range sqs.Quantities {
		z := &sqs.Quantities[i]
		if z.name == name {
			return z
		}
	}
	log.Log.ErrAndExit("Error: The dataset %v has not been initialized.", name)
	return nil
}

// createOrUpdateSavedQuantity is the unified function for saving quantities.
func (sqs *savedQuantitiesType) createOrUpdateSavedQuantity(q Quantity, name string, period float64, rchunks requestedChunking) {
	if !sqs.savedQuandtityExists(name) {
		sqs.createSavedQuantity(q, name, rchunks, period)
	} else {
		sqs.updateSavedQuantity(q, name, rchunks, period)
	}
	sq := sqs.getSavedQuantity(name)
	if period == 0 {
		sq.Save()
	}
}

func (sqs *savedQuantitiesType) autoSaveInner(q Quantity, name string, period float64, rchunks requestedChunking) {
	if period == 0 {
		sq := sqs.getSavedQuantity(name)
		sq.period = 0
	} else {
		sqs.createOrUpdateSavedQuantity(q, name, period, rchunks)
	}
}

// User-facing save functions (function signatures cannot change)
func (sqs *savedQuantitiesType) autoSave(q Quantity, period float64) {
	sqs.autoSaveInner(q, q.Name(), period, requestedChunking{1, 1, 1, 1})
}

func (sqs *savedQuantitiesType) autoSaveAs(q Quantity, name string, period float64) {
	sqs.autoSaveInner(q, name, period, requestedChunking{1, 1, 1, 1})
}

func (sqs *savedQuantitiesType) autoSaveAsChunk(q Quantity, name string, period float64, rchunks requestedChunking) {
	sqs.autoSaveInner(q, name, period, rchunks)
}

func (sqs *savedQuantitiesType) saveAsInner(q Quantity, name string, rchunks requestedChunking) {
	if !sqs.savedQuandtityExists(name) {
		sqs.createSavedQuantity(q, name, rchunks, 0)
	}
	sqs.getSavedQuantity(name).Save()
}
func (sqs *savedQuantitiesType) saveAs(q Quantity, name string) {
	sqs.saveAsInner(q, name, requestedChunking{1, 1, 1, 1})
}

func (sqs *savedQuantitiesType) save(q Quantity) {
	sqs.saveAsInner(q, q.Name(), requestedChunking{1, 1, 1, 1})
}

func (sqs *savedQuantitiesType) saveAsChunk(q Quantity, name string, rchunks requestedChunking) {
	sqs.saveAsInner(q, name, rchunks)
}
