package engine

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/DataDog/zstd"

	"github.com/MathieuMoalic/amumax/src/chunk"
	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/quantity"
)

type savedQuantity struct {
	e        *engineState
	name     string
	q        quantity.Quantity
	period   float64
	times    []float64
	chunks   chunk.Chunks
	rchunks  chunk.RequestedChunking
	nextTime float64 // Next time when autosave should trigger
}

func newSavedQuantity(engineState *engineState, q quantity.Quantity, name string, rchunks chunk.RequestedChunking, period float64) *savedQuantity {
	return &savedQuantity{
		e:       engineState,
		name:    name,
		q:       q,
		period:  period,
		times:   []float64{},
		chunks:  chunk.NewChunks(engineState.log, q, rchunks),
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
	err := sq.e.fs.SaveZattrs(sq.name+"/.zattrs", map[string]interface{}{"t": sq.times})
	sq.e.log.PanicIfError(err)
}

// func valueOf(q quantity.Quantity) *data.Slice {
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
		sq.e.log.PanicIfError(err)
	})
}

// syncSave writes the data slice into chunked, compressed files compatible with the Zarr format.
func (sq *savedQuantity) syncSave(array *data.Slice, qname string, step int, chunks chunk.Chunks) error {
	data := array.Tensors()
	size := array.Size()
	ncomp := array.NComp()

	// Save .zarray metadata
	err := sq.e.fs.SaveFileZarray(
		fmt.Sprintf("%s/.zarray", qname),
		size,
		ncomp,
		step,
		chunks.Z.Len, chunks.Y.Len, chunks.X.Len, chunks.C.Len,
	)
	if err != nil {
		return err
	}

	// Iterate over chunks and save data
	for icx := 0; icx < chunks.X.Count; icx++ {
		for icy := 0; icy < chunks.Y.Count; icy++ {
			for icz := 0; icz < chunks.Z.Count; icz++ {
				for icc := 0; icc < chunks.C.Count; icc++ {
					var bdata bytes.Buffer
					for iz := 0; iz < chunks.Z.Len; iz++ {
						z := icz*chunks.Z.Len + iz
						for iy := 0; iy < chunks.Y.Len; iy++ {
							y := icy*chunks.Y.Len + iy
							for ix := 0; ix < chunks.X.Len; ix++ {
								x := icx*chunks.X.Len + ix
								for ic := 0; ic < chunks.C.Len; ic++ {
									c := icc*chunks.C.Len + ic
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
	e          *engineState
	Quantities []savedQuantity
}

func newSavedQuantities(engineState *engineState) *savedQuantities {
	sqs := &savedQuantities{e: engineState}
	engineState.script.RegisterFunction("Save", sqs.save)
	engineState.script.RegisterFunction("SaveAs", sqs.saveAs)
	engineState.script.RegisterFunction("SaveAsChunks", sqs.saveAsChunk)
	engineState.script.RegisterFunction("AutoSave", sqs.autoSave)
	engineState.script.RegisterFunction("AutoSaveAs", sqs.autoSaveAs)
	engineState.script.RegisterFunction("AutoSaveAsChunk", sqs.autoSaveAsChunk)
	engineState.script.RegisterFunction("Chunks", chunk.CreateRequestedChunk)
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

func (sqs *savedQuantities) createSavedQuantity(q quantity.Quantity, name string, rchunks chunk.RequestedChunking, period float64) *savedQuantity {
	if sqs.e.fs.Exists(name) {
		err := sqs.e.fs.Remove(name)
		sqs.e.log.PanicIfError(err)
	}
	err := sqs.e.fs.Mkdir(name)
	sqs.e.log.PanicIfError(err)
	sq := newSavedQuantity(sqs.e, q, name, rchunks, period)
	sqs.Quantities = append(sqs.Quantities, *sq)
	return sq

}

func (sqs *savedQuantities) updateSavedQuantity(q quantity.Quantity, name string, rchunks chunk.RequestedChunking, period float64) {
	sq := sqs.getSavedQuantity(name)
	if sq.rchunks != rchunks {
		sq.e.log.ErrAndExit("Error: The dataset %v has already been initialized with different chunks.", name)
	} else if sq.q != q {
		sq.e.log.ErrAndExit("Error: The dataset %v has already been initialized with a different quantity.Quantity.", name)
	} else if sq.period != period {
		if sq.period == 0 && period != 0 {
			// enable autosave
			sq.period = period
			sq.nextTime = sqs.e.solver.time + period
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
	sqs.e.log.ErrAndExit("Error: The dataset %v has not been initialized.", name)
	return nil
}

// createOrUpdateSavedQuantity is the unified function for saving quantities.
func (sqs *savedQuantities) createOrUpdateSavedQuantity(q quantity.Quantity, name string, period float64, rchunks chunk.RequestedChunking) {
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

func (sqs *savedQuantities) autoSaveInner(q quantity.Quantity, name string, period float64, rchunks chunk.RequestedChunking) {
	if period == 0 {
		sq := sqs.getSavedQuantity(name)
		sq.period = 0
	} else {
		sqs.createOrUpdateSavedQuantity(q, name, period, rchunks)
	}
}

// User-facing save functions (function signatures cannot change)
func (sqs *savedQuantities) autoSave(q quantity.Quantity, period float64) {
	sqs.autoSaveInner(q, q.Name(), period, chunk.RequestedChunking{X: 1, Y: 1, Z: 1, C: 1})
}

func (sqs *savedQuantities) autoSaveAs(q quantity.Quantity, name string, period float64) {
	sqs.autoSaveInner(q, name, period, chunk.RequestedChunking{X: 1, Y: 1, Z: 1, C: 1})
}

func (sqs *savedQuantities) autoSaveAsChunk(q quantity.Quantity, name string, period float64, rchunks chunk.RequestedChunking) {
	sqs.autoSaveInner(q, name, period, rchunks)
}

func (sqs *savedQuantities) saveAsInner(q quantity.Quantity, name string, rchunks chunk.RequestedChunking) {
	if !sqs.savedQuandtityExists(name) {
		sqs.createSavedQuantity(q, name, rchunks, 0)
	}
	sqs.getSavedQuantity(name).save()
}
func (sqs *savedQuantities) saveAs(q quantity.Quantity, name string) {
	sqs.saveAsInner(q, name, chunk.RequestedChunking{X: 1, Y: 1, Z: 1, C: 1})
}

func (sqs *savedQuantities) save(q quantity.Quantity) {
	sqs.saveAsInner(q, q.Name(), chunk.RequestedChunking{X: 1, Y: 1, Z: 1, C: 1})
}

func (sqs *savedQuantities) saveAsChunk(q quantity.Quantity, name string, rchunks chunk.RequestedChunking) {
	sqs.saveAsInner(q, name, rchunks)
}
