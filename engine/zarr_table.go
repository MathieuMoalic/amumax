package engine

import (
	"time"

	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/httpfs"
	"github.com/MathieuMoalic/amumax/script"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/MathieuMoalic/amumax/zarr"
)

func init() {
	DeclFunc("TableSave", ZTableSave, "Save the data table right now.")
	DeclFunc("TableAdd", ZTableAdd, "Save the data table periodically.")
	DeclFunc("TableAddVar", ZTableAddVar, "Save the data table periodically.")
	DeclFunc("TableAddAs", ZTableAddAs, "Save the data table periodically.")
	DeclFunc("TableAutoSave", ZTableAutoSave, "Save the data table periodically.")
	ZTables = ZTablesStruct{Data: make(map[string][]float64), Step: -1, AutoSavePeriod: 0.0, FlushInterval: 5 * time.Second}
}

var ZTables ZTablesStruct

// the Table is kept in RAM and used for the API
type ZTablesStruct struct {
	quantities     []Quantity
	tables         []ZTable
	Data           map[string][]float64 `json:"data"`
	AutoSavePeriod float64              `json:"autoSavePeriod"`
	AutoSaveStart  float64              `json:"autoSaveStart"`
	Step           int                  `json:"step"`
	FlushInterval  time.Duration        `json:"flushInterval"`
}

type ZTable struct {
	Name   string
	buffer []byte
	io     httpfs.WriteCloseFlusher
}

func (ts *ZTablesStruct) WriteToBuffer() {
	buf := []float64{}
	// always save the current time
	buf = append(buf, Time)
	// for each quantity we append each component to the buffer
	for _, q := range ts.quantities {
		buf = append(buf, AverageOf(q)...)
	}
	// size of buf should be same as size of []Ztable
	for i, b := range buf {
		ts.tables[i].buffer = append(ts.tables[i].buffer, zarr.Float64ToByte(b)...)
		ts.Data[ts.tables[i].Name] = append(ts.Data[ts.tables[i].Name], b)
		// ts.Tables[i].Data = append(ts.Tables[i].Data, b)
	}
}

func (ts *ZTablesStruct) Flush() {
	for i := range ts.tables {
		ts.tables[i].io.Write(ts.tables[i].buffer)
		ts.tables[i].buffer = []byte{}
		// saving .zarray before the data might help resolve some unsync
		// errors when the simulation is running and the user loads data
		zarr.SaveFileTableZarray(OD()+"table/"+ts.tables[i].Name, ts.Step)
		ts.tables[i].io.Flush()
	}
}

func (ts *ZTablesStruct) NeedSave() bool {
	return ts.AutoSavePeriod != 0 && (Time-ts.AutoSaveStart)-float64(ts.Step)*ts.AutoSavePeriod >= ts.AutoSavePeriod
}

func (ts *ZTablesStruct) Exists(q Quantity, name string) bool {
	suffixes := []string{"x", "y", "z"}
	for _, i := range ZTables.tables {
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

func TableInit() {
	httpfs.Remove(OD() + "table")
	zarr.MakeZgroup("table", OD(), &zGroups)
	err := httpfs.Mkdir(OD() + "table/t")
	util.FatalErr(err)
	f, err := httpfs.Create(OD() + "table/t/0")
	util.FatalErr(err)
	ZTables.tables = append(ZTables.tables, ZTable{"t", []byte{}, f})
	ZTableAdd(&M)
	go ZTablesAutoFlush()

}

func ZTablesAutoFlush() {
	for {
		ZTables.Flush()
		time.Sleep(ZTables.FlushInterval)
	}
}

func ZTableSave() {
	if len(ZTables.tables) == 0 {
		TableInit()
	}
	ZTables.Step += 1
	ZTables.WriteToBuffer()
}

func CreateTable(name string) ZTable {
	err := httpfs.Mkdir(OD() + "table/" + name)
	util.FatalErr(err)
	f, err := httpfs.Create(OD() + "table/" + name + "/0")
	util.FatalErr(err)
	return ZTable{Name: name, buffer: []byte{}, io: f}
}

func ZTableAdd(q Quantity) {
	ZTableAddAs(q, NameOf(q))
}

func ZTableAddAs(q Quantity, name string) {
	suffixes := []string{"x", "y", "z"}
	if ZTables.Step != -1 {
		util.Fatal("Add Table Quantity BEFORE you save the table for the first time")
	}
	if len(ZTables.tables) == 0 {
		TableInit()
	}

	if ZTables.Exists(q, name) {
		LogOut(name, " is already in the table. Ignoring.")
		return
	}
	ZTables.quantities = append(ZTables.quantities, q)
	if q.NComp() == 1 {
		ZTables.tables = append(ZTables.tables, CreateTable(name))
	} else {
		for comp := 0; comp < q.NComp(); comp++ {
			ZTables.tables = append(ZTables.tables, CreateTable(name+suffixes[comp]))
		}
	}
}

func ZTableAutoSave(period float64) {
	ZTables.AutoSaveStart = Time
	ZTables.AutoSavePeriod = period
}

func ZTableAddVar(customvar script.ScalarFunction, name, unit string) {
	ZTableAdd(&userVar{customvar, name, unit})
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
