package engine

import (
	"fmt"
	"time"

	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/httpfs"
	"github.com/MathieuMoalic/amumax/script"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/MathieuMoalic/amumax/zarr"
)

func init() {
	DeclFunc("TableSave", TableSave, "Save the data table right now.")
	DeclFunc("TableAdd", TableAdd, "Save the data table periodically.")
	DeclFunc("TableAddVar", TableAddVar, "Save the data table periodically.")
	DeclFunc("TableAddAs", TableAddAs, "Save the data table periodically.")
	DeclFunc("TableAutoSave", TableAutoSave, "Save the data table periodically.")
	Table = TableStruct{Data: make(map[string][]float64), Step: -1, AutoSavePeriod: 0.0, FlushInterval: 5 * time.Second}
}

var Table TableStruct

// the Table is kept in RAM and used for the API
type TableStruct struct {
	quantities     []Quantity
	columns        []Column
	Data           map[string][]float64 `json:"data"`
	AutoSavePeriod float64              `json:"autoSavePeriod"`
	AutoSaveStart  float64              `json:"autoSaveStart"`
	Step           int                  `json:"step"`
	FlushInterval  time.Duration        `json:"flushInterval"`
}

type Column struct {
	Name   string
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
	for i, b := range buf {
		ts.columns[i].buffer = append(ts.columns[i].buffer, zarr.Float64ToByte(b)...)
		ts.Data[ts.columns[i].Name] = append(ts.Data[ts.columns[i].Name], b)
	}
}

func (ts *TableStruct) Flush() {
	for i := range ts.columns {
		ts.columns[i].io.Write(ts.columns[i].buffer)
		ts.columns[i].buffer = []byte{}
		// saving .zarray before the data might help resolve some unsync
		// errors when the simulation is running and the user loads data
		zarr.SaveFileTableZarray(OD()+"table/"+ts.columns[i].Name, ts.Step)
		ts.columns[i].io.Flush()
	}
}

func (ts *TableStruct) NeedSave() bool {
	return ts.AutoSavePeriod != 0 && (Time-ts.AutoSaveStart)-float64(ts.Step)*ts.AutoSavePeriod >= ts.AutoSavePeriod
}

func (ts *TableStruct) Exists(q Quantity, name string) bool {
	suffixes := []string{"x", "y", "z"}
	for _, i := range Table.columns {
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

func (ts *TableStruct) GetTableNames() []string {
	names := []string{}
	for _, i := range ts.columns {
		names = append(names, i.Name)
	}
	return names
}

func TableInit() {
	httpfs.Remove(OD() + "table")
	zarr.MakeZgroup("table", OD(), &zGroups)
	err := httpfs.Mkdir(OD() + "table/t")
	util.FatalErr(err)
	f, err := httpfs.Create(OD() + "table/t/0")
	util.FatalErr(err)
	Table.columns = append(Table.columns, Column{"t", []byte{}, f})
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
	if len(Table.columns) == 0 {
		TableInit()
	}
	Table.Step += 1
	Table.WriteToBuffer()
}

func CreateTable(name string) Column {
	err := httpfs.Mkdir(OD() + "table/" + name)
	util.FatalErr(err)
	f, err := httpfs.Create(OD() + "table/" + name + "/0")
	util.FatalErr(err)
	return Column{Name: name, buffer: []byte{}, io: f}
}

func TableAdd(q Quantity) {
	TableAddAs(q, NameOf(q))
}

func TableAddAs(q Quantity, name string) {
	suffixes := []string{"x", "y", "z"}
	if Table.Step != -1 {
		util.LogWarn("You cannot add a new quantity to the table after the simulation has started. Ignoring.")
	}
	if len(Table.columns) == 0 {
		TableInit()
	}

	if Table.Exists(q, name) {
		util.LogWarn(fmt.Sprint(name, " is already in the table. Ignoring."))
		return
	}
	Table.quantities = append(Table.quantities, q)
	if q.NComp() == 1 {
		Table.columns = append(Table.columns, CreateTable(name))
	} else {
		for comp := 0; comp < q.NComp(); comp++ {
			Table.columns = append(Table.columns, CreateTable(name+suffixes[comp]))
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
