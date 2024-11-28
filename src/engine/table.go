package engine

import (
	"bufio"
	"os"
	"sync"
	"time"

	"github.com/MathieuMoalic/amumax/src/log_old"
	"github.com/MathieuMoalic/amumax/src/zarr"
)

// the table is kept in RAM and used for the API
type table struct {
	e              *engineState
	quantities     []quantity
	columns        []column
	Data           map[string][]float64 `json:"data"`
	AutoSavePeriod float64              `json:"autoSavePeriod"`
	AutoSaveStart  float64              `json:"autoSaveStart"`
	Step           int                  `json:"step"`
	FlushInterval  time.Duration        `json:"flushInterval"`
	mu             sync.Mutex
}

type column struct {
	name   string
	unit   string
	writer *bufio.Writer
	file   *os.File
}

func newTable(e *engineState) *table {
	t := &table{
		e:          e,
		Data:       make(map[string][]float64),
		Step:       -1,
		quantities: []quantity{},
		columns:    []column{},
	}
	err := e.fs.Remove("table")
	log_old.Log.PanicIfError(err)
	zarr.InitZgroup("table", e.zarrPath)
	t.addColumn("step", "")
	t.addColumn("t", "s")
	e.world.registerFunction("TableAutoSave", t.tableAutoSave)
	e.world.registerFunction("TableAdd", t.tableAdd)
	e.world.registerFunction("TableAddAs", t.tableAddAs)
	e.world.registerFunction("TableSave", t.tableSave)
	return t
}

func (ts *table) writeToBuffer() {
	buf := []float64{}
	buf = append(buf, float64(ts.Step))
	buf = append(buf, ts.e.solver.time)
	for _, q := range ts.quantities {
		buf = append(buf, q.Average()...)
	}

	ts.mu.Lock() // Lock before modifying shared data
	defer ts.mu.Unlock()

	for i, b := range buf {
		// Convert float64 to bytes
		data := zarr.Float64ToBytes(b)

		// Write directly to the buffered writer
		_, err := ts.columns[i].writer.Write(data)
		if err != nil {
			log_old.Log.PanicIfError(err)
		}

		// Update in-memory data
		ts.Data[ts.columns[i].name] = append(ts.Data[ts.columns[i].name], b)
	}
}

func (ts *table) flush() {
	for i := range ts.columns {
		// Update zarray if necessary
		zarr.SaveFileTableZarray(ts.e.zarrPath+"table/"+ts.columns[i].name, ts.Step)
		err := ts.columns[i].writer.Flush()
		if err != nil {
			log_old.Log.PanicIfError(err)
		}
	}
}

// func (ts *table) needSave() bool {
// 	return ts.AutoSavePeriod != 0 && (ts.e.solver.time-ts.AutoSaveStart)-float64(ts.Step)*ts.AutoSavePeriod >= ts.AutoSavePeriod
// }

func (ts *table) exists(q quantity, name string) bool {
	suffixes := []string{"x", "y", "z"}
	for _, col := range ts.columns {
		if q.NComp() == 1 {
			if col.name == name {
				return true
			}
		} else {
			for comp := 0; comp < q.NComp(); comp++ {
				if name+suffixes[comp] == col.name {
					return true
				}
			}
		}
	}
	return false
}

func (ts *table) addColumn(name, unit string) {
	err := ts.e.fs.Mkdir("table/" + name)
	log_old.Log.PanicIfError(err)
	writer, file, err := ts.e.fs.Create("table/" + name + "/0")
	log_old.Log.PanicIfError(err)
	ts.columns = append(ts.columns, column{name: name, unit: unit, writer: writer, file: file})
}

// func (ts *table) tablesAutoFlush() {
// 	for {
// 		ts.flush()
// 		time.Sleep(ts.FlushInterval)
// 	}
// }

func (ts *table) tableSave() {
	if len(ts.columns) == 0 {
		ts.e.log.Warn("No columns in table, not saving.")
	}
	ts.Step += 1
	ts.writeToBuffer()
}

func (ts *table) tableAdd(q quantity) {
	ts.tableAddAs(q, q.Name())
}

func (ts *table) tableAddAs(q quantity, name string) {
	suffixes := []string{"x", "y", "z"}
	if ts.Step != -1 {
		log_old.Log.Warn("You cannot add a new quantity to the table after the simulation has started. Ignoring.")
	}
	if len(ts.columns) == 0 {
		ts.e.log.Warn("No columns in table, not saving.")
	}

	if ts.exists(q, name) {
		log_old.Log.Warn("%s is already in the table. Ignoring.", name)
		return
	}
	ts.quantities = append(ts.quantities, q)
	if q.NComp() == 1 {
		ts.addColumn(name, q.Unit())
	} else {
		for comp := 0; comp < q.NComp(); comp++ {
			ts.addColumn(name+suffixes[comp], q.Unit())
		}
	}
}

func (ts *table) tableAutoSave(period float64) {
	ts.AutoSaveStart = ts.e.solver.time
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
