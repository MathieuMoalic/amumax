package new_engine

import (
	"bufio"
	"os"
	"sync"
	"time"

	"github.com/MathieuMoalic/amumax/src/log"
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
	writer *bufio.Writer
	file   *os.File
}

func (ts *Table) WriteToBuffer() {
	buf := []float64{}
	buf = append(buf, float64(ts.Step))
	buf = append(buf, ts.EngineState.solver.Time)
	for _, q := range ts.quantities {
		buf = append(buf, q.Average()...)
	}

	ts.Mu.Lock() // Lock before modifying shared data
	defer ts.Mu.Unlock()

	for i, b := range buf {
		// Convert float64 to bytes
		data := zarr.Float64ToBytes(b)

		// Write directly to the buffered writer
		_, err := ts.Columns[i].writer.Write(data)
		if err != nil {
			log.Log.PanicIfError(err)
		}

		// Update in-memory data
		ts.Data[ts.Columns[i].Name] = append(ts.Data[ts.Columns[i].Name], b)
	}
}

func (ts *Table) Flush() {
	for i := range ts.Columns {
		// Update zarray if necessary
		zarr.SaveFileTableZarray(ts.EngineState.zarrPath+"table/"+ts.Columns[i].Name, ts.Step)
		err := ts.Columns[i].writer.Flush()
		if err != nil {
			log.Log.PanicIfError(err)
		}
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
	writer, file, err := ts.EngineState.fs.Create("table/" + name + "/0")
	log.Log.PanicIfError(err)
	ts.Columns = append(ts.Columns, column{Name: name, Unit: unit, writer: writer, file: file})
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
