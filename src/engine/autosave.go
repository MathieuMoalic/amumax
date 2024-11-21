package engine

// Bookkeeping for auto-saving quantities at given intervals.

import (
	"fmt"
)

var (
	output  = make(map[Quantity]*autosave) // when to save quantities
	autonum = make(map[string]int)         // auto number for out file
)

// Periodically called by run loop to save everything that's needed at this time.
func saveIfNeeded() {
	for q, a := range output {
		if a.needSave() {
			a.save(q)
			a.count++
		}
	}
	saveZarrArrays()
	if Table.NeedSave() {
		tableSave()
	}
	if EngineState.Metadata.NeedSave() {
		EngineState.Metadata.Save()
	}
}

// Register quant to be auto-saved every period.
// period == 0 stops autosaving.
func autoSaveOvf(q Quantity, period float64) {
	autoSaveInner(q, period, saveovf)
}

// Register quant to be auto-saved as image, every period.
func autoSnapshot(q Quantity, period float64) {
	autoSaveInner(q, period, snapshot)
}

// register save(q) to be called every period
func autoSaveInner(q Quantity, period float64, save func(Quantity)) {
	if period == 0 {
		delete(output, q)
	} else {
		output[q] = &autosave{period, Time, -1, save} // init count to -1 allows save at t=0
	}
}

// generate auto file name based on save count and FilenameFormat. E.g.:
//
//	m000001.ovf
func autoFname(name string, format outputFormatType, num int) string {
	return fmt.Sprintf(OD()+filenameFormat+"."+StringFromOutputFormat[format], name, num)
}

// keeps info needed to decide when a quantity needs to be periodically saved
type autosave struct {
	period float64        // How often to save
	start  float64        // Starting point
	count  int            // Number of times it has been autosaved
	save   func(Quantity) // called to do the actual save
}

// returns true when the time is right to save.
func (a *autosave) needSave() bool {
	t := Time - a.start
	return a.period != 0 && t-float64(a.count)*a.period >= a.period
}
