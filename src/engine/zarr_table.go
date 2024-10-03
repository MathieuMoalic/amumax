package engine

import (
	"sync"
	"time"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/httpfs"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/script"
	"github.com/MathieuMoalic/amumax/src/zarr"
)

func init() {
	DeclFunc("TableSave", TableSave, "Save the data table right now.")
	DeclFunc("TableAdd", TableAdd, "Save the data table periodically.")
	DeclFunc("TableAddVar", TableAddVar, "Save the data table periodically.")
	DeclFunc("TableAddAs", TableAddAs, "Save the data table periodically.")
	DeclFunc("TableAutoSave", TableAutoSave, "Save the data table periodically.")
	Table = TableStruct{
		Data:           make(map[string][]float64),
		Step:           -1,
		AutoSavePeriod: 0.0,
		FlushInterval:  5 * time.Second,
	}
}

var Table TableStruct

// the Table is kept in RAM and used for the API
type TableStruct struct {
	quantities     []Quantity
	Columns        []Column
	Data           map[string][]float64 `json:"data"`
	AutoSavePeriod float64              `json:"autoSavePeriod"`
	AutoSaveStart  float64              `json:"autoSaveStart"`
	Step           int                  `json:"step"`
	FlushInterval  time.Duration        `json:"flushInterval"`
	Mu             sync.Mutex
}

type Column struct {
	Name   string
	Unit   string
	buffer []byte
	io     httpfs.WriteCloseFlusher
}

func (ts *TableStruct) WriteToBuffer() {
	buf := []float64{}
	// always save the current time
	buf = append(buf, Time)
	// for each quantity we append each component to the buffer
	for _, q := range ts.quantities {
		buf = append(buf, AverageOf(q)...)
	}
	// size of buf should be same as size of []Ztable
	ts.Mu.Lock() // Lock the mutex before modifying the map
	defer ts.Mu.Unlock()
	for i, b := range buf {
		ts.Columns[i].buffer = append(ts.Columns[i].buffer, zarr.Float64ToBytes(b)...)
		ts.Data[ts.Columns[i].Name] = append(ts.Data[ts.Columns[i].Name], b)
	}
}

func (ts *TableStruct) Flush() {
	for i := range ts.Columns {
		_, err := ts.Columns[i].io.Write(ts.Columns[i].buffer)
		log.Log.PanicIfError(err)
		ts.Columns[i].buffer = []byte{}
		// saving .zarray before the data might help resolve some unsync
		// errors when the simulation is running and the user loads data
		zarr.SaveFileTableZarray(OD()+"table/"+ts.Columns[i].Name, ts.Step)
		ts.Columns[i].io.Flush()
	}
}

func (ts *TableStruct) NeedSave() bool {
	return ts.AutoSavePeriod != 0 && (Time-ts.AutoSaveStart)-float64(ts.Step)*ts.AutoSavePeriod >= ts.AutoSavePeriod
}

func (ts *TableStruct) Exists(q Quantity, name string) bool {
	suffixes := []string{"x", "y", "z"}
	for _, i := range Table.Columns {
		if q.NComp() == 1 {
			if i.Name == name {
				return true
			}
		} else {
			for comp := 0; comp < q.NComp(); comp++ {
				if name+suffixes[comp] == i.Name {
					return true
				}
			}
		}
	}
	return false
}

func (ts *TableStruct) AddColumn(name, unit string) {
	err := httpfs.Mkdir(OD() + "table/" + name)
	log.Log.PanicIfError(err)
	f, err := httpfs.Create(OD() + "table/" + name + "/0")
	log.Log.PanicIfError(err)
	ts.Columns = append(ts.Columns, Column{Name: name, Unit: unit, buffer: []byte{}, io: f})
}

func TableInit() {
	err := httpfs.Remove(OD() + "table")
	log.Log.PanicIfError(err)
	zarr.MakeZgroup("table", OD(), &zGroups)
	err = httpfs.Mkdir(OD() + "table/t")
	log.Log.PanicIfError(err)
	f, err := httpfs.Create(OD() + "table/t/0")
	log.Log.PanicIfError(err)
	Table.Columns = append(Table.Columns, Column{"t", "s", []byte{}, f})
	TableAdd(&M)
	go TablesAutoFlush()

}

func TablesAutoFlush() {
	for {
		Table.Flush()
		time.Sleep(Table.FlushInterval)
	}
}

func TableSave() {
	if len(Table.Columns) == 0 {
		TableInit()
	}
	Table.Step += 1
	Table.WriteToBuffer()
}

func TableAdd(q Quantity) {
	TableAddAs(q, NameOf(q))
}

func TableAddAs(q Quantity, name string) {
	suffixes := []string{"x", "y", "z"}
	if Table.Step != -1 {
		log.Log.Warn("You cannot add a new quantity to the table after the simulation has started. Ignoring.")
	}
	if len(Table.Columns) == 0 {
		TableInit()
	}

	if Table.Exists(q, name) {
		log.Log.Warn("%s is already in the table. Ignoring.", name)
		return
	}
	Table.quantities = append(Table.quantities, q)
	if q.NComp() == 1 {
		Table.AddColumn(name, UnitOf(q))
	} else {
		for comp := 0; comp < q.NComp(); comp++ {
			Table.AddColumn(name+suffixes[comp], UnitOf(q))
		}
	}
}

func TableAutoSave(period float64) {
	Table.AutoSaveStart = Time
	Table.AutoSavePeriod = period
}

func TableAddVar(customvar script.ScalarFunction, name, unit string) {
	TableAdd(&userVar{customvar, name, unit})
}

type userVar struct {
	value      script.ScalarFunction
	name, unit string
}

func (x *userVar) Name() string       { return x.name }
func (x *userVar) NComp() int         { return 1 }
func (x *userVar) Unit() string       { return x.unit }
func (x *userVar) average() []float64 { return []float64{x.value.Float()} }
func (x *userVar) EvalTo(dst *data.Slice) {
	avg := x.average()
	for c := 0; c < x.NComp(); c++ {
		cuda.Memset(dst.Comp(c), float32(avg[c]))
	}
}
