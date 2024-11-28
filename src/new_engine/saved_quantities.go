package new_engine

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/DataDog/zstd"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/zarr"
)

type savedQuantity struct {
	e        *engineState
	name     string
	q        quantity
	period   float64
	times    []float64
	chunks   chunks
	rchunks  requestedChunking
	nextTime float64 // Next time when autosave should trigger
}

func newSavedQuantity(engineState *engineState, q quantity, name string, rchunks requestedChunking, period float64) *savedQuantity {
	return &savedQuantity{
		e:       engineState,
		name:    name,
		q:       q,
		period:  period,
		times:   []float64{},
		chunks:  newChunks(q, rchunks),
		rchunks: rchunks,
	}
}

// // needSave returns true when it's time to save based on the period.
// func (sq *savedQuantity) needSave() bool {
// 	if sq.period == 0 {
// 		return false
// 	}
// 	return sq.e.solver.time >= sq.nextTime
// }

// saveAttrs updates the .zattrs file with the times data.
func (sq *savedQuantity) saveAttrs() {
	u, err := json.Marshal(zarr.Zattrs{Buffer: sq.times})
	log.Log.PanicIfError(err)
	err = sq.e.fs.Remove(sq.name + "/.zattrs")
	log.Log.PanicIfError(err)
	err = sq.e.fs.Put(sq.name+"/.zattrs", u)
	log.Log.PanicIfError(err)
}

// func valueOf(q quantity) *data.Slice {
// 	// TODO: check for Buffered() implementation
// 	buf := cuda.Buffer(q.NComp(), q.Size())
// 	q.EvalTo(buf)
// 	return buf
// }

// save writes the data to disk and updates the times.
func (sq *savedQuantity) save() {
	sq.times = append(sq.times, sq.e.solver.time)
	sq.saveAttrs()
	buffer := cuda.Buffer(sq.q.NComp(), sq.q.Size())
	sq.q.EvalTo(buffer)
	defer cuda.Recycle(buffer)
	dataSlice := buffer.HostCopy()
	sq.e.fs.QueueOutput(func() {
		err := sq.syncSave(dataSlice, sq.name, len(sq.times), sq.chunks)
		log.Log.PanicIfError(err)
	})
}

// syncSave writes the data slice into chunked, compressed files compatible with the Zarr format.
func (sq *savedQuantity) syncSave(array *data.Slice, qname string, step int, chunks chunks) error {
	data := array.Tensors()
	size := array.Size()
	ncomp := array.NComp()

	// Save .zarray metadata
	zarr.SaveFileZarray(
		fmt.Sprintf(sq.e.zarrPath+"%s/.zarray", qname),
		size,
		ncomp,
		step,
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
					filename := fmt.Sprintf("%s/%d.%d.%d.%d.%d", qname, step+1, icz, icy, icx, icc)
					err = sq.e.fs.Put(filename, compressedData)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

type savedQuantities struct {
	EngineState *engineState
	Quantities  []savedQuantity
}

func newSavedQuantities(engineState *engineState) *savedQuantities {
	sqs := &savedQuantities{EngineState: engineState}
	engineState.world.registerFunction("Save", sqs.save)
	engineState.world.registerFunction("SaveAs", sqs.saveAs)
	engineState.world.registerFunction("SaveAsChunks", sqs.saveAsChunk)
	engineState.world.registerFunction("AutoSave", sqs.autoSave)
	engineState.world.registerFunction("AutoSaveAs", sqs.autoSaveAs)
	engineState.world.registerFunction("AutoSaveAsChunk", sqs.autoSaveAsChunk)
	engineState.world.registerFunction("Chunks", createRequestedChunk)
	return sqs
}

// // saveZarrArrays is called periodically to save arrays when needed.
// func (sqs *savedQuantities) saveIfNeeded() {
// 	for i := range sqs.Quantities {
// 		if sqs.Quantities[i].needSave() {
// 			sqs.Quantities[i].save()
// 			sqs.Quantities[i].nextTime = sqs.Quantities[i].nextTime + sqs.Quantities[i].period
// 		}
// 	}
// }

func (sqs *savedQuantities) savedQuandtityExists(name string) bool {
	for _, z := range sqs.Quantities {
		if z.name == name {
			return true
		}
	}
	return false
}

func (sqs *savedQuantities) createSavedQuantity(q quantity, name string, rchunks requestedChunking, period float64) *savedQuantity {
	if sqs.EngineState.fs.Exists(name) {
		err := sqs.EngineState.fs.Remove(name)
		log.Log.PanicIfError(err)
	}
	err := sqs.EngineState.fs.Mkdir(name)
	log.Log.PanicIfError(err)
	sq := newSavedQuantity(sqs.EngineState, q, name, rchunks, period)
	sqs.Quantities = append(sqs.Quantities, *sq)
	return sq

}

func (sqs *savedQuantities) updateSavedQuantity(q quantity, name string, rchunks requestedChunking, period float64) {
	sq := sqs.getSavedQuantity(name)
	if sq.rchunks != rchunks {
		log.Log.ErrAndExit("Error: The dataset %v has already been initialized with different chunks.", name)
	} else if sq.q != q {
		log.Log.ErrAndExit("Error: The dataset %v has already been initialized with a different quantity.", name)
	} else if sq.period != period {
		if sq.period == 0 && period != 0 {
			// enable autosave
			sq.period = period
			sq.nextTime = sqs.EngineState.solver.time + period
		} else if sq.period != 0 && period == 0 {
			// disable autosave
			sq.period = period
		}
	}
}

func (sqs *savedQuantities) getSavedQuantity(name string) *savedQuantity {
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
func (sqs *savedQuantities) createOrUpdateSavedQuantity(q quantity, name string, period float64, rchunks requestedChunking) {
	if !sqs.savedQuandtityExists(name) {
		sqs.createSavedQuantity(q, name, rchunks, period)
	} else {
		sqs.updateSavedQuantity(q, name, rchunks, period)
	}
	sq := sqs.getSavedQuantity(name)
	if period == 0 {
		sq.save()
	}
}

func (sqs *savedQuantities) autoSaveInner(q quantity, name string, period float64, rchunks requestedChunking) {
	if period == 0 {
		sq := sqs.getSavedQuantity(name)
		sq.period = 0
	} else {
		sqs.createOrUpdateSavedQuantity(q, name, period, rchunks)
	}
}

// User-facing save functions (function signatures cannot change)
func (sqs *savedQuantities) autoSave(q quantity, period float64) {
	sqs.autoSaveInner(q, q.Name(), period, requestedChunking{1, 1, 1, 1})
}

func (sqs *savedQuantities) autoSaveAs(q quantity, name string, period float64) {
	sqs.autoSaveInner(q, name, period, requestedChunking{1, 1, 1, 1})
}

func (sqs *savedQuantities) autoSaveAsChunk(q quantity, name string, period float64, rchunks requestedChunking) {
	sqs.autoSaveInner(q, name, period, rchunks)
}

func (sqs *savedQuantities) saveAsInner(q quantity, name string, rchunks requestedChunking) {
	if !sqs.savedQuandtityExists(name) {
		sqs.createSavedQuantity(q, name, rchunks, 0)
	}
	sqs.getSavedQuantity(name).save()
}
func (sqs *savedQuantities) saveAs(q quantity, name string) {
	sqs.saveAsInner(q, name, requestedChunking{1, 1, 1, 1})
}

func (sqs *savedQuantities) save(q quantity) {
	sqs.saveAsInner(q, q.Name(), requestedChunking{1, 1, 1, 1})
}

func (sqs *savedQuantities) saveAsChunk(q quantity, name string, rchunks requestedChunking) {
	sqs.saveAsInner(q, name, rchunks)
}
