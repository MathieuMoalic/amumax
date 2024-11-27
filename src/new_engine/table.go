package new_engine

import (
	"sync"
	"time"

	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/new_fsutil"
	"github.com/MathieuMoalic/amumax/src/zarr"
)

// the Table is kept in RAM and used for the API
type Table struct {
	EngineState    *EngineStateStruct
	quantities     []Quantity
	Columns        []column
	Data           map[string][]float64 `json:"data"`
	AutoSavePeriod float64              `json:"autoSavePeriod"`
	AutoSaveStart  float64              `json:"autoSaveStart"`
	Step           int                  `json:"step"`
	FlushInterval  time.Duration        `json:"flushInterval"`
	Mu             sync.Mutex
}

type column struct {
	Name   string
	Unit   string
	buffer []byte
	io     new_fsutil.WriteCloseFlusher
}

func (ts *Table) WriteToBuffer() {
	buf := []float64{}
	buf = append(buf, float64(ts.Step))
	// always save the current time
	buf = append(buf, ts.EngineState.solver.Time)
	// for each quantity we append each component to the buffer
	for _, q := range ts.quantities {
		buf = append(buf, q.Average()...)
	}
	// size of buf should be same as size of []Ztable
	ts.Mu.Lock() // Lock the mutex before modifying the map
	defer ts.Mu.Unlock()
	for i, b := range buf {
		ts.Columns[i].buffer = append(ts.Columns[i].buffer, zarr.Float64ToBytes(b)...)
		ts.Data[ts.Columns[i].Name] = append(ts.Data[ts.Columns[i].Name], b)
	}
}

func (ts *Table) Flush() {
	for i := range ts.Columns {
		_, err := ts.Columns[i].io.Write(ts.Columns[i].buffer)
		log.Log.PanicIfError(err)
		ts.Columns[i].buffer = []byte{}
		// saving .zarray before the data might help resolve some unsync
		// errors when the simulation is running and the user loads data
		zarr.SaveFileTableZarray(ts.EngineState.zarrPath+"table/"+ts.Columns[i].Name, ts.Step)
		ts.Columns[i].io.Flush()
	}
}

func (ts *Table) NeedSave() bool {
	return ts.AutoSavePeriod != 0 && (ts.EngineState.solver.Time-ts.AutoSaveStart)-float64(ts.Step)*ts.AutoSavePeriod >= ts.AutoSavePeriod
}

func (ts *Table) Exists(q Quantity, name string) bool {
	suffixes := []string{"x", "y", "z"}
	for _, col := range ts.Columns {
		if q.NComp() == 1 {
			if col.Name == name {
				return true
			}
		} else {
			for comp := 0; comp < q.NComp(); comp++ {
				if name+suffixes[comp] == col.Name {
					return true
				}
			}
		}
	}
	return false
}

func (ts *Table) AddColumn(name, unit string) {
	err := ts.EngineState.fs.Mkdir("table/" + name)
	log.Log.PanicIfError(err)
	f, err := ts.EngineState.fs.Create("table/" + name + "/0")
	log.Log.PanicIfError(err)
	ts.Columns = append(ts.Columns, column{Name: name, Unit: unit, buffer: []byte{}, io: f})
}

func (ts *Table) tablesAutoFlush() {
	for {
		ts.Flush()
		time.Sleep(ts.FlushInterval)
	}
}

func (ts *Table) tableSave() {
	if len(ts.Columns) == 0 {
		ts.EngineState.log.Warn("No columns in table, not saving.")
	}
	ts.Step += 1
	ts.WriteToBuffer()
}

func (ts *Table) tableAdd(q Quantity) {
	ts.tableAddAs(q, q.Name())
}

func (ts *Table) tableAddAs(q Quantity, name string) {
	suffixes := []string{"x", "y", "z"}
	if ts.Step != -1 {
		log.Log.Warn("You cannot add a new quantity to the table after the simulation has started. Ignoring.")
	}
	if len(ts.Columns) == 0 {
		ts.EngineState.log.Warn("No columns in table, not saving.")
	}

	if ts.Exists(q, name) {
		log.Log.Warn("%s is already in the table. Ignoring.", name)
		return
	}
	ts.quantities = append(ts.quantities, q)
	if q.NComp() == 1 {
		ts.AddColumn(name, q.Unit())
	} else {
		for comp := 0; comp < q.NComp(); comp++ {
			ts.AddColumn(name+suffixes[comp], q.Unit())
		}
	}
}

func (ts *Table) TableAutoSave(period float64) {
	ts.AutoSaveStart = ts.EngineState.solver.Time
	ts.AutoSavePeriod = period
}

// func (ts *TableStruct) tableAddVar(customvar script.ScalarFunction, name, unit string) {
// 	ts.tableAdd(&userVar{customvar, name, unit})
// }

// type userVar struct {
// 	value      script.ScalarFunction
// 	name, unit string
// }

// func (x *userVar) Name() string       { return x.name }
// func (x *userVar) NComp() int         { return 1 }
// func (x *userVar) Unit() string       { return x.unit }
// func (x *userVar) average() []float64 { return []float64{x.value.Float()} }
// func (x *userVar) EvalTo(dst *data.Slice) {
// 	avg := x.average()
// 	for c := 0; c < x.NComp(); c++ {
// 		cuda.Memset(dst.Comp(c), float32(avg[c]))
// 	}
// }
