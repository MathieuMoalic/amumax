package engine_old

import (
	"sync"
	"time"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/fsutil_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/script_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/zarr_old"
)

func init() {
	Table = tableStruct{
		Data:           make(map[string][]float64),
		Step:           -1,
		AutoSavePeriod: 0.0,
		FlushInterval:  5 * time.Second,
	}
}

var Table tableStruct

// the Table is kept in RAM and used for the API
type tableStruct struct {
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
	io     fsutil_old.WriteCloseFlusher
}

func (ts *tableStruct) WriteToBuffer() {
	buf := []float64{}
	buf = append(buf, float64(ts.Step))
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
		ts.Columns[i].buffer = append(ts.Columns[i].buffer, zarr_old.Float64ToBytes(b)...)
		ts.Data[ts.Columns[i].Name] = append(ts.Data[ts.Columns[i].Name], b)
	}
}

func (ts *tableStruct) Flush() {
	for i := range ts.Columns {
		_, err := ts.Columns[i].io.Write(ts.Columns[i].buffer)
		log_old.Log.PanicIfError(err)
		ts.Columns[i].buffer = []byte{}
		// saving .zarray before the data might help resolve some unsync
		// errors when the simulation is running and the user loads data
		zarr_old.SaveFileTableZarray(OD()+"table/"+ts.Columns[i].Name, ts.Step)
		ts.Columns[i].io.Flush()
	}
}

func (ts *tableStruct) NeedSave() bool {
	return ts.AutoSavePeriod != 0 && (Time-ts.AutoSaveStart)-float64(ts.Step)*ts.AutoSavePeriod >= ts.AutoSavePeriod
}

func (ts *tableStruct) Exists(q Quantity, name string) bool {
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

func (ts *tableStruct) AddColumn(name, unit string) {
	err := fsutil_old.Mkdir(OD() + "table/" + name)
	log_old.Log.PanicIfError(err)
	f, err := fsutil_old.Create(OD() + "table/" + name + "/0")
	log_old.Log.PanicIfError(err)
	ts.Columns = append(ts.Columns, column{Name: name, Unit: unit, buffer: []byte{}, io: f})
}

func tableInit() {
	err := fsutil_old.Remove(OD() + "table")
	log_old.Log.PanicIfError(err)
	zarr_old.InitZgroup("table", OD())
	Table.AddColumn("step", "")
	Table.AddColumn("t", "s")
	tableAdd(&NormMag)
	go tablesAutoFlush()

}

func tablesAutoFlush() {
	for {
		Table.Flush()
		time.Sleep(Table.FlushInterval)
	}
}

func tableSave() {
	if len(Table.Columns) == 0 {
		tableInit()
	}
	Table.Step += 1
	Table.WriteToBuffer()
}

func tableAdd(q Quantity) {
	tableAddAs(q, nameOf(q))
}

func tableAddAs(q Quantity, name string) {
	suffixes := []string{"x", "y", "z"}
	if Table.Step != -1 {
		log_old.Log.Warn("You cannot add a new quantity to the table after the simulation has started. Ignoring.")
	}
	if len(Table.Columns) == 0 {
		tableInit()
	}

	if Table.Exists(q, name) {
		log_old.Log.Warn("%s is already in the table. Ignoring.", name)
		return
	}
	Table.quantities = append(Table.quantities, q)
	if q.NComp() == 1 {
		Table.AddColumn(name, unitOf(q))
	} else {
		for comp := 0; comp < q.NComp(); comp++ {
			Table.AddColumn(name+suffixes[comp], unitOf(q))
		}
	}
}

func tableAutoSave(period float64) {
	Table.AutoSaveStart = Time
	Table.AutoSavePeriod = period
}

func tableAddVar(customvar script_old.ScalarFunction, name, unit string) {
	tableAdd(&userVar{customvar, name, unit})
}

type userVar struct {
	value      script_old.ScalarFunction
	name, unit string
}

func (x *userVar) Name() string       { return x.name }
func (x *userVar) NComp() int         { return 1 }
func (x *userVar) Unit() string       { return x.unit }
func (x *userVar) average() []float64 { return []float64{x.value.Float()} }
func (x *userVar) EvalTo(dst *data_old.Slice) {
	avg := x.average()
	for c := 0; c < x.NComp(); c++ {
		cuda.Memset(dst.Comp(c), float32(avg[c]))
	}
}
